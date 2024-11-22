package comic_walker

import (
	"encoding/hex"
	"fmt"
	"github.com/sekiju/htt"
	"github.com/sekiju/mdl/internal/renamer"
	manga "github.com/sekiju/mdl/internal/sdk/extractor/manga"
	"regexp"
	"strconv"
)

type Extractor struct{}

type searchFn func(episodeID string) ([]*manga.Chapter, error)

func (p *Extractor) FindChapters(URL string) ([]*manga.Chapter, error) {
	parsedURL, err := parseURL(URL)
	if err != nil {
		return nil, err
	}

	res, err := htt.New().Getf("https://comic-walker.com/api/contents/details/episode?workCode=%s&episodeType=first", parsedURL.WorkCode)
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
			ID:      episodeResult.Episode.Id,
			Number:  strconv.Itoa(episodeResult.Episode.Internal.EpisodeNo),
			Title:   episodeResult.Episode.Title,
			Index:   uint(episodeResult.Episode.Internal.EpisodeNo - 1),
			URL:     fmt.Sprintf("https://comic-walker.com/detail/%s/episodes/%s", parsedURL.WorkCode, episodeResult.Episode.Code),
			MangaID: parsedURL.WorkCode,
		},
	}

	var fn searchFn
	fn = func(episodeID string) ([]*manga.Chapter, error) {
		res, err = htt.New().Getf("https://comic-walker.com/api/contents/viewer-jump-forward?episodeId=%s", episodeID)
		if err != nil {
			return nil, err
		}

		var viewerJumpForwardResult ViewerJumpForwardResult
		if err = res.JSON(&viewerJumpForwardResult); err != nil {
			return nil, err
		}

		if viewerJumpForwardResult.Episode != nil {
			chapters = append(chapters, &manga.Chapter{
				ID:      viewerJumpForwardResult.Episode.Id,
				Number:  strconv.Itoa(viewerJumpForwardResult.Episode.Internal.EpisodeNo),
				Title:   viewerJumpForwardResult.Episode.Title,
				Index:   uint(viewerJumpForwardResult.Episode.Internal.EpisodeNo - 1),
				URL:     fmt.Sprintf("https://comic-walker.com/detail/%s/episodes/%s", parsedURL.WorkCode, viewerJumpForwardResult.Episode.Code),
				MangaID: parsedURL.WorkCode,
			})

			return fn(viewerJumpForwardResult.Episode.Id)
		}

		return chapters, nil
	}

	return fn(episodeResult.Episode.Id)
}

func (p *Extractor) FindChapter(URL string) (*manga.Chapter, error) {
	parsedURL, err := parseURL(URL)
	if err != nil {
		return nil, err
	}

	res, err := htt.New().Getf("https://comic-walker.com/api/contents/details/episode?workCode=%s&episodeCode=%s&episodeType=first", parsedURL.WorkCode, parsedURL.EpisodeCode)
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
		ID:      episodeResult.Episode.Id,
		Number:  strconv.Itoa(episodeResult.Episode.Internal.EpisodeNo),
		Title:   episodeResult.Episode.Title,
		Index:   uint(episodeResult.Episode.Internal.EpisodeNo - 1),
		URL:     fmt.Sprintf("https://comic-walker.com/detail/%s/episodes/%s", parsedURL.WorkCode, episodeResult.Episode.Code),
		MangaID: parsedURL.WorkCode,
	}, nil
}

func (p *Extractor) FindChapterPages(chapter *manga.Chapter) ([]*manga.Page, error) {
	res, err := htt.New().Getf("https://comic-walker.com/api/contents/viewer?episodeId=%s&imageSizeType=width%%3A1284", chapter.ID)
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

	pages := make([]*manga.Page, len(viewerResult.Manuscripts))
	padRenamer := renamer.New(len(viewerResult.Manuscripts))

	for index, page := range viewerResult.Manuscripts {
		pages[index] = &manga.Page{
			Index:    uint(index),
			URL:      page.DrmImageUrl,
			Filename: padRenamer.Name(index, ".webp"),
			Decode: func(b []byte) ([]byte, error) {
				keyBytes, err := hex.DecodeString(page.DrmHash[:16])
				if err != nil {
					return nil, err
				}

				r, i := len(b), len(keyBytes)
				decodedBytes := make([]byte, r)

				for a := 0; a < r; a++ {
					decodedBytes[a] = b[a] ^ keyBytes[a%i]
				}

				return decodedBytes, nil
			},
		}
	}

	return pages, nil
}

func New() manga.Extractor {
	return &Extractor{}
}

var re = regexp.MustCompile("https://comic-walker.com/detail/(KC_[a-zA-Z0-9_]*)(/episodes/(KC_[a-zA-Z0-9_]*))?")

type extractorURL struct {
	WorkCode    string
	EpisodeCode string
}

func parseURL(URL string) (*extractorURL, error) {
	matches := re.FindStringSubmatch(URL)
	if len(matches) < 2 {
		return nil, manga.ErrInvalidURLFormat
	}

	res := &extractorURL{WorkCode: matches[1]}

	if len(matches) == 4 {
		res.EpisodeCode = matches[3]
	}

	return res, nil
}
