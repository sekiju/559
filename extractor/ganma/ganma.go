package ganma

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/sekiju/htt"
	"github.com/sekiju/mdl/internal/renamer"
	"github.com/sekiju/mdl/internal/util"
	"github.com/sekiju/mdl/sdk/manga"
	"net/url"
	"regexp"
	"strconv"
)

type Extractor struct {
	cookie string
}

const (
	sha256hash = "75f44fb799c1d505ae245b52633b59e3db9be1e2fe90b4c766eb5d96a86d5be7"
)

func (e *Extractor) FindChapters(URL string) ([]*manga.Chapter, error) {
	res, err := htt.New().Get(URL)
	if err != nil {
		return nil, err
	}

	text, err := res.Text()
	if err != nil {
		return nil, err
	}

	mangaID, err := util.ExtractStringFromHTML(text, `\"magazineId\":\"`, `\"`)
	if err != nil {
		return nil, err
	}

	res, err = htt.New().
		SetHeader("Cookie", e.cookie).
		SetHeader("X-From", "https://reader.ganma.jp/api/").
		Getf("https://reader.ganma.jp/api/3.2/magazines/%s", mangaID)
	if err != nil {
		return nil, err
	}

	var magazine magazineResult
	if err = res.JSON(&magazine); err != nil {
		return nil, err
	}

	chapters := make([]*manga.Chapter, len(magazine.Root.Items))
	for i, item := range magazine.Root.Items {
		chapters[i] = &manga.Chapter{
			ID:      item.StoryId,
			Number:  strconv.Itoa(item.Number),
			Title:   item.Title,
			Index:   uint(i),
			URL:     fmt.Sprintf("https://ganma.jp/web/reader/%s/%s/0", magazine.Root.Id, item.StoryId),
			MangaID: magazine.Root.Id,
		}
	}

	return chapters, nil
}

var re = regexp.MustCompile(`https://ganma.jp/web/reader/([a-zA-Z0-9_-]*)/([a-zA-Z0-9_-]*)`)

func (e *Extractor) FindChapter(URL string) (*manga.Chapter, error) {
	matches := re.FindStringSubmatch(URL)
	if len(matches) != 3 {
		return nil, manga.ErrInvalidURLFormat
	}

	if util.IsValidUUID(matches[1]) {
		res, err := htt.New().
			SetHeader("Cookie", e.cookie).
			SetHeader("X-From", "https://reader.ganma.jp/api/").
			Getf("https://reader.ganma.jp/api/3.2/magazines/%s", matches[1])
		if err != nil {
			return nil, err
		}

		var magazine magazineResult
		if err = res.JSON(&magazine); err != nil {
			return nil, err
		}

		for _, item := range magazine.Root.Items {
			if item.StoryId == matches[2] {
				return &manga.Chapter{
					ID:      item.StoryId,
					Number:  strconv.Itoa(item.Number),
					Title:   item.Title,
					Index:   0,
					URL:     fmt.Sprintf("https://ganma.jp/web/reader/%s/%s/0", magazine.Root.Id, item.StoryId),
					MangaID: magazine.Root.Id,
				}, nil
			}
		}
	}

	res, err := htt.New().Get(URL)
	if err != nil {
		return nil, err
	}

	text, err := res.Text()
	if err != nil {
		return nil, err
	}

	mangaID, err := util.ExtractStringFromHTML(text, `\"magazineId\":\"`, `\"`)
	if err != nil {
		return nil, err
	}

	return &manga.Chapter{
		ID:      matches[2],
		Number:  "",
		Title:   "",
		Index:   0,
		URL:     URL,
		MangaID: mangaID,
	}, nil
}

func (e *Extractor) FindChapterPages(chapter *manga.Chapter) ([]*manga.Page, error) {
	fmt.Println()

	res, err := htt.New().
		SetHeader("Cookie", e.cookie).
		SetHeader("X-From", "https://reader.ganma.jp/api/").
		Getf(
			"https://ganma.jp/api/graphql?operationName=MagazineStoryReaderQuery&variables=%s&extensions=%s",
			url.QueryEscape(fmt.Sprintf(`{"magazineIdOrAlias":%q,"storyId":%q,"publicKey":null}`, chapter.MangaID, chapter.ID)),
			url.QueryEscape(fmt.Sprintf(`{"persistedQuery":{"version":1,"sha256Hash":%q}}`, sha256hash)),
		)
	if err != nil {
		return nil, err
	}

	// todo: handle STORY_COUNT_LIMITED == PaidChapter

	var reader readerResult
	if err = res.JSON(&reader); err != nil {
		return nil, err
	}

	pages := make([]*manga.Page, reader.Data.Magazine.StoryContents.PageImages.PageCount)
	padRenamer := renamer.New(reader.Data.Magazine.StoryContents.PageImages.PageCount)

	for index := range reader.Data.Magazine.StoryContents.PageImages.PageCount {
		pages[index] = &manga.Page{
			URL:      fmt.Sprintf("%s%d.jpg?%s&w=4000", reader.Data.Magazine.StoryContents.PageImages.PageImageBaseURL, index+1, reader.Data.Magazine.StoryContents.PageImages.PageImageSign),
			Filename: padRenamer.Name(index, ".jpeg"),
			Index:    uint(index),
		}
	}

	return pages, nil
}

func New(cookie string) (manga.Extractor, error) {
	if cookie == "CREATE" {
		res, err := htt.New().SetHeader("X-From", "https://reader.ganma.jp/api/").Post("https://reader.ganma.jp/api/1.0/account")
		if err != nil {
			return nil, err
		}

		if res.StatusCode != 200 {
			return nil, errors.New("failed to create account")
		}

		var createAccount createAccountResponse
		if err = res.JSON(&createAccount); err != nil {
			return nil, err
		}

		res, err = htt.New().SetHeader("X-From", "https://reader.ganma.jp/api/").
			Body(createAccount.Root).
			Post("https://reader.ganma.jp/api/3.0/session")
		if err != nil {
			return nil, err
		}

		if res.StatusCode != 200 {
			return nil, errors.New("failed to login with generated account")
		}

		cookie = res.Header.Get("Set-Cookie")

		log.Info().Msgf("New cookie created for ganma.jp >>> %s", cookie)
	}

	return &Extractor{cookie}, nil
}
