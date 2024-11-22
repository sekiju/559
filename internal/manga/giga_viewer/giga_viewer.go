package giga_viewer

import (
	"github.com/sekiju/htt"
	"github.com/sekiju/mdl/internal/renamer"
	"github.com/sekiju/mdl/internal/sdk/extractor/manga"
	"github.com/sekiju/mdl/internal/util"
	"strconv"
)

type Extractor struct {
	hostname     string
	CookieString *string
}

type searchFn func(URL string) ([]*manga.Chapter, error)

func (e *Extractor) FindChapters(URL string) ([]*manga.Chapter, error) {
	chapters := make([]*manga.Chapter, 0)
	visitedIDs := make(map[string]bool)

	var fn searchFn
	fn = func(episodeURL string) ([]*manga.Chapter, error) {
		if visitedIDs[episodeURL] {
			return nil, nil
		}

		visitedIDs[episodeURL] = true

		res, err := htt.New().Get(episodeURL)
		if err != nil {
			return nil, err
		}

		if res.StatusCode == 404 {
			return nil, manga.ErrMangaNotFound
		}

		html, err := res.Text()
		if err != nil {
			return nil, err
		}

		episodeResult, err := util.ExtractJSONFromHTML[EpisodeResult](html, `<script id='episode-json' type='text/json' data-value='`, `'></script>`)
		if err != nil {
			return nil, err
		}

		chapters = append(chapters, &manga.Chapter{
			ID:      episodeResult.ReadableProduct.Id,
			Number:  strconv.Itoa(episodeResult.ReadableProduct.Number),
			Title:   episodeResult.ReadableProduct.Title,
			Index:   uint(episodeResult.ReadableProduct.Number - 1),
			URL:     episodeResult.ReadableProduct.Permalink,
			MangaID: episodeResult.ReadableProduct.Id,
		})

		if prevURI := episodeResult.ReadableProduct.PrevReadableProductUri; prevURI != nil {
			_, err = fn(*prevURI)
			if err != nil {
				return nil, err
			}
		}

		if nextURI := episodeResult.ReadableProduct.NextReadableProductUri; nextURI != nil {
			_, err = fn(*nextURI)
			if err != nil {
				return nil, err
			}
		}

		return chapters, nil
	}

	return fn(URL)
}

func (e *Extractor) FindChapter(URL string) (*manga.Chapter, error) {
	res, err := htt.New().Get(URL)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 404 {
		return nil, manga.ErrChapterNotFound
	}

	html, err := res.Text()
	if err != nil {
		return nil, err
	}

	episodeResult, err := util.ExtractJSONFromHTML[EpisodeResult](html, `<script id='episode-json' type='text/json' data-value='`, `'></script>`)
	if err != nil {
		return nil, err
	}

	return &manga.Chapter{
		ID:      episodeResult.ReadableProduct.Id,
		Number:  strconv.Itoa(episodeResult.ReadableProduct.Number),
		Title:   episodeResult.ReadableProduct.Title,
		Index:   uint(episodeResult.ReadableProduct.Number - 1),
		URL:     episodeResult.ReadableProduct.Permalink,
		MangaID: episodeResult.ReadableProduct.Id,
	}, nil
}

func (e *Extractor) FindChapterPages(chapter *manga.Chapter) ([]*manga.Page, error) {
	req := htt.New().SetHeader("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1")
	if e.CookieString != nil {
		req.SetHeader("Cookie", *e.CookieString)
	}

	res, err := req.Get(chapter.URL)

	if err != nil {
		return nil, err
	}

	if res.StatusCode == 404 {
		return nil, manga.ErrChapterNotFound
	}

	html, err := res.Text()
	if err != nil {
		return nil, err
	}

	episodeResult, err := util.ExtractJSONFromHTML[EpisodeResult](html, `<script id='episode-json' type='text/json' data-value='`, `'></script>`)
	if err != nil {
		return nil, err
	}

	if !episodeResult.ReadableProduct.IsPublic && !episodeResult.ReadableProduct.HasPurchased {
		return nil, manga.ErrPaidChapter
	}

	var mainPages []*EpisodeResultPage
	for _, page := range episodeResult.ReadableProduct.PageStructure.Pages {
		if page.Type != "main" {
			continue
		}

		mainPages = append(mainPages, &page)
	}

	chapterPages := make([]*manga.Page, len(mainPages))
	padRenamer := renamer.New(len(mainPages))

	for index, page := range mainPages {
		chapterPages[index] = &manga.Page{
			Index:    uint(index),
			URL:      page.Src,
			Filename: padRenamer.Name(index, ".jpg"),
		}
	}

	return chapterPages, nil
}

func New(hostname string) manga.Extractor {
	return &Extractor{hostname: hostname}
}

func NewAuthorized(hostname, cookieString string) manga.Extractor {
	return &Extractor{hostname: hostname, CookieString: &cookieString}
}
