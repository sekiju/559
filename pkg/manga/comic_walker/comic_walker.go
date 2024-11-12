package comic_walker

import (
	"encoding/hex"
	"fmt"
	"github.com/sekiju/mary/pkg/manga/internal/renamer"
	"github.com/sekiju/mary/pkg/sdk/extractor/manga"
	"github.com/sekiju/rq"
	"regexp"
	"strconv"
)

type Provider struct{}

const ProviderID = "comic_walker"

func New() manga.Provider {
	return &Provider{}
}

var re = regexp.MustCompile("https://comic-walker.com/detail/(KC_[a-zA-Z0-9_]*)(/episodes/(KC_[a-zA-Z0-9_]*))?")

func (p *Provider) ExtractMangaID(URL string) (manga.ExtractedURL, error) {
	matches := re.FindStringSubmatch(URL)
	if len(matches) < 2 {
		return nil, manga.ErrInvalidURLFormat
	}

	res := manga.ExtractedURL{"manga": matches[1]}

	if len(matches) == 4 {
		res["chapter"] = matches[3]
	}

	return res, nil
}

func (p *Provider) FindManga(ID manga.ID) (*manga.Manga, error) {
	extractedURL, err := extractURL(ID)
	if err != nil {
		return nil, err
	}

	res, err := rq.New().Getf("https://comic-walker.com/api/contents/details/work?workCode=%s", extractedURL.MustMangaID())
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 404 {
		return nil, manga.ErrMangaNotFound
	}

	var workResult WorkResult
	if err = res.JSON(&workResult); err != nil {
		return nil, err
	}

	return &manga.Manga{
		Provider: ProviderID,
		ID:       workResult.Work.Code,
		Title:    workResult.Work.Title,
		Cover:    &workResult.Work.BookCover,
		URL:      fmt.Sprintf("https://comic-walker.com/detail/%s", workResult.Work.Code),
	}, nil
}

type searchFn func(episodeID string) ([]*manga.Chapter, error)

func (p *Provider) FindChapters(ID manga.ID) ([]*manga.Chapter, error) {
	extractedURL, err := extractURL(ID)
	if err != nil {
		return nil, err
	}

	workCode := extractedURL.MustMangaID().(string)

	res, err := rq.New().Getf("https://comic-walker.com/api/contents/details/episode?workCode=%s&episodeType=first", workCode)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 404 {
		return nil, manga.ErrMangaNotFound
	}

	var episodeResult EpisodeResult
	if err = res.JSON(&episodeResult); err != nil {
		return nil, err
	}

	chapters := []*manga.Chapter{
		{
			Provider: ProviderID,
			ID:       episodeResult.Episode.Id,
			Number:   strconv.Itoa(episodeResult.Episode.Internal.EpisodeNo),
			Title:    episodeResult.Episode.Title,
			Index:    episodeResult.Episode.Internal.EpisodeNo - 1,
			URL:      fmt.Sprintf("https://comic-walker.com/detail/%s/episodes/%s", workCode, episodeResult.Episode.Code),
			MangaID:  workCode,
		},
	}

	var fn searchFn
	fn = func(episodeID string) ([]*manga.Chapter, error) {
		res, err = rq.New().Getf("https://comic-walker.com/api/contents/viewer-jump-forward?episodeId=%s", episodeID)
		if err != nil {
			return nil, err
		}

		var viewerJumpForwardResult ViewerJumpForwardResult
		if err = res.JSON(&viewerJumpForwardResult); err != nil {
			return nil, err
		}

		if viewerJumpForwardResult.Episode != nil {
			chapters = append(chapters, &manga.Chapter{
				Provider: ProviderID,
				ID:       viewerJumpForwardResult.Episode.Id,
				Number:   strconv.Itoa(viewerJumpForwardResult.Episode.Internal.EpisodeNo),
				Title:    viewerJumpForwardResult.Episode.Title,
				Index:    viewerJumpForwardResult.Episode.Internal.EpisodeNo - 1,
				URL:      fmt.Sprintf("https://comic-walker.com/detail/%s/episodes/%s", workCode, viewerJumpForwardResult.Episode.Code),
				MangaID:  workCode,
			})

			return fn(viewerJumpForwardResult.Episode.Id)
		}

		return chapters, nil
	}

	return fn(episodeResult.Episode.Id)
}

func (p *Provider) FindChapter(ID manga.ID) (*manga.Chapter, error) {
	extractedURL, err := extractURL(ID)
	if err != nil {
		return nil, err
	}

	workCode := extractedURL.MustMangaID().(string)
	episodeCode, err := extractedURL.ChapterID()
	if err != nil {
		return nil, err
	}

	res, err := rq.New().Getf("https://comic-walker.com/api/contents/details/episode?workCode=%s&episodeCode=%s&episodeType=first", workCode, episodeCode)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 404 {
		return nil, manga.ErrChapterNotFound
	}

	var episodeResult EpisodeResult
	if err = res.JSON(&episodeResult); err != nil {
		return nil, err
	}

	return &manga.Chapter{
		Provider: ProviderID,
		ID:       episodeResult.Episode.Id,
		Number:   strconv.Itoa(episodeResult.Episode.Internal.EpisodeNo),
		Title:    episodeResult.Episode.Title,
		Index:    episodeResult.Episode.Internal.EpisodeNo - 1,
		URL:      fmt.Sprintf("https://comic-walker.com/detail/%s/episodes/%s", workCode, episodeResult.Episode.Code),
		MangaID:  workCode,
	}, nil
}

func (p *Provider) ExtractPages(chapter *manga.Chapter) ([]*manga.Page, error) {
	res, err := rq.New().Getf("https://comic-walker.com/api/contents/viewer?episodeId=%s&imageSizeType=width%%3A1284", chapter.ID)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 404 {
		return nil, manga.ErrChapterNotFound
	}

	var viewerResult ViewerResult
	if err = res.JSON(&viewerResult); err != nil {
		return nil, err
	}

	pages := make([]*manga.Page, 0)
	padRenamer := renamer.NewPadRenamer(len(viewerResult.Manuscripts))

	for index, page := range viewerResult.Manuscripts {
		var decodeFn manga.DescrambleFn = func(bytes []byte) ([]byte, error) {
			keyBytes, err := hex.DecodeString(page.DrmHash[:16])
			if err != nil {
				return nil, err
			}

			r, i := len(bytes), len(keyBytes)
			decodedBytes := make([]byte, r)

			for a := 0; a < r; a++ {
				decodedBytes[a] = bytes[a] ^ keyBytes[a%i]
			}

			return decodedBytes, nil
		}

		pages = append(pages, &manga.Page{
			Provider:     ProviderID,
			Index:        index,
			URL:          page.DrmImageUrl,
			Filename:     padRenamer.NewName(index, ".webp"),
			DescrambleFn: &decodeFn,
		})
	}

	return pages, nil
}

func extractURL(ID manga.ID) (manga.ExtractedURL, error) {
	switch v := ID.(type) {
	case manga.ExtractedURL:
		return v, nil
	default:
		return nil, manga.ErrStringAsIDUnsupported
	}
}
