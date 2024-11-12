package comic_walker

import (
	"encoding/hex"
	"fmt"
	"github.com/sekiju/mary/pkg/manga/internal/renamer"
	"github.com/sekiju/mary/pkg/sdk/extractor/manga"
	"github.com/sekiju/rq"
	"regexp"
	"strconv"
	"strings"
)

type Provider struct{}

const ProviderID = "comic_walker"

func New() manga.Provider {
	return &Provider{}
}

var urlRegex = regexp.MustCompile("https://comic-walker.com/detail/(KC_[a-zA-Z0-9_]*)(/episodes/(KC_[a-zA-Z0-9_]*))?")

func (p *Provider) ExtractMangaID(URL string) (string, error) {
	matches := urlRegex.FindStringSubmatch(URL)
	if len(matches) < 2 {
		return "", manga.ErrInvalidURLFormat
	}

	if len(matches) == 4 {
		matches = append(matches[:2], matches[3:]...)
	}

	return strings.Join(matches[1:], "$"), nil
}

func (p *Provider) FindManga(mangaID string) (*manga.Manga, error) {
	IDs := strings.Split(mangaID, "$")
	res, err := rq.New().Getf("https://comic-walker.com/api/contents/details/work?workCode=%s", IDs[0])
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

func (p *Provider) FindChapters(mangaID string) ([]*manga.Chapter, error) {
	IDs := strings.Split(mangaID, "$")
	res, err := rq.New().Getf("https://comic-walker.com/api/contents/details/episode?workCode=%s&episodeType=first", IDs[0])
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
			URL:      fmt.Sprintf("https://comic-walker.com/detail/%s/episodes/%s", mangaID, episodeResult.Episode.Code),
			MangaID:  mangaID,
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
				URL:      fmt.Sprintf("https://comic-walker.com/detail/%s/episodes/%s", mangaID, viewerJumpForwardResult.Episode.Code),
				MangaID:  mangaID,
			})

			return fn(viewerJumpForwardResult.Episode.Id)
		}

		return chapters, nil
	}

	return fn(episodeResult.Episode.Id)
}

func (p *Provider) FindChapter(chapterID string) (*manga.Chapter, error) {
	IDs := strings.Split(chapterID, "$")
	res, err := rq.New().Getf("https://comic-walker.com/api/contents/details/episode?workCode=%s&episodeCode=%s&episodeType=first", IDs[0], IDs[1])
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
		URL:      fmt.Sprintf("https://comic-walker.com/detail/%s/episodes/%s", IDs[0], episodeResult.Episode.Code),
		MangaID:  IDs[0],
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
