package giga_viewer

import (
	"fmt"
	"github.com/sekiju/mary/pkg/manga/internal/renamer"
	"github.com/sekiju/mary/pkg/manga/internal/util"
	"github.com/sekiju/mary/pkg/sdk/extractor/manga"
	"github.com/sekiju/rq"
	"regexp"
	"strconv"
)

type Provider struct {
	hostname  string
	sessionID *string
}

var re = regexp.MustCompile(`(magazine|episode)/(\d+)`)

func (p *Provider) ExtractMangaID(URL string) (manga.ExtractedURL, error) {
	matches := re.FindStringSubmatch(URL)
	if len(matches) != 3 {
		return nil, manga.ErrInvalidURLFormat
	}

	return manga.ExtractedURL{"url": URL}, nil
}

func (p *Provider) FindManga(mangaID manga.ID) (*manga.Manga, error) {
	URL, err := extractURL(mangaID)
	if err != nil {
		return nil, err
	}

	res, err := rq.New().Get(URL)
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

	data := &manga.Manga{
		Provider: ProviderID,
		ID:       episodeResult.ReadableProduct.Id,
		URL:      episodeResult.ReadableProduct.Permalink,
	}

	if episodeResult.ReadableProduct.Series != nil {
		data.Title = episodeResult.ReadableProduct.Series.Title
		data.Cover = &episodeResult.ReadableProduct.Series.ThumbnailUri
	} else if episodeResult.ReadableProduct.Toc != nil {
		data.Title = episodeResult.ReadableProduct.Toc.Title
		data.Cover = &episodeResult.ReadableProduct.Toc.ThumbnailUrl
	}

	return data, nil
}

type searchFn func(URL string) ([]*manga.Chapter, error)

func (p *Provider) FindChapters(mangaID manga.ID) ([]*manga.Chapter, error) {
	URL, err := extractURL(mangaID)
	if err != nil {
		return nil, err
	}

	chapters := make([]*manga.Chapter, 0)
	visitedIDs := make(map[string]bool)

	var fn searchFn
	fn = func(episodeURL string) ([]*manga.Chapter, error) {
		if visitedIDs[episodeURL] {
			return nil, nil
		}

		visitedIDs[episodeURL] = true

		res, err := rq.New().Get(episodeURL)
		if err != nil {
			return nil, err
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
			Provider: ProviderID,
			ID:       episodeResult.ReadableProduct.Id,
			Number:   strconv.Itoa(episodeResult.ReadableProduct.Number),
			Title:    episodeResult.ReadableProduct.Title,
			Index:    episodeResult.ReadableProduct.Number - 1,
			URL:      episodeResult.ReadableProduct.Permalink,
			MangaID:  episodeResult.ReadableProduct.Id,
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

func (p *Provider) FindChapter(chapterID manga.ID) (*manga.Chapter, error) {
	URL, err := extractURL(chapterID)
	if err != nil {
		return nil, err
	}

	res, err := rq.New().Get(URL)
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
		Provider: ProviderID,
		ID:       episodeResult.ReadableProduct.Id,
		Number:   strconv.Itoa(episodeResult.ReadableProduct.Number),
		Title:    episodeResult.ReadableProduct.Title,
		Index:    episodeResult.ReadableProduct.Number - 1,
		URL:      episodeResult.ReadableProduct.Permalink,
		MangaID:  episodeResult.ReadableProduct.Id,
	}, nil
}

func (p *Provider) ExtractPages(chapter *manga.Chapter) ([]*manga.Page, error) {
	req := rq.New(rq.SetHeader("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1"))
	if p.sessionID != nil {
		req.Config.Set(rq.SetHeader("Cookie", fmt.Sprintf("glsc=%s", *p.sessionID)))
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

	mainPages := make([]*EpisodeResultPage, 0)
	for _, page := range episodeResult.ReadableProduct.PageStructure.Pages {
		if page.Type != "main" {
			continue
		}

		mainPages = append(mainPages, &page)
	}

	chapterPages := make([]*manga.Page, 0)
	padRenamer := renamer.NewPadRenamer(len(mainPages))

	for index, page := range mainPages {
		chapterPages = append(chapterPages, &manga.Page{
			Provider:     ProviderID,
			Index:        index,
			URL:          page.Src,
			Filename:     padRenamer.NewName(index, ".jpg"),
			DescrambleFn: nil,
		})
	}

	return chapterPages, nil
}

func extractURL(ID manga.ID) (string, error) {
	var URL string

	switch v := ID.(type) {
	case manga.ExtractedURL:
		strURL, err := v.URL()
		if err != nil {
			return "", err
		}

		URL = strURL
	case string:
		URL = v
	default:
		return "", manga.ErrInvalidID
	}

	return URL, nil
}

const ProviderID = "giga_viewer"

func New(hostname string) manga.Provider {
	return &Provider{hostname: hostname}
}

func NewWithSession(hostname, sessionID string) manga.Provider {
	return &Provider{hostname: hostname, sessionID: &sessionID}
}
