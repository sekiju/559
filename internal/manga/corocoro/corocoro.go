package corocoro

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/sekiju/htt"
	"github.com/sekiju/mdl/internal/renamer"
	"github.com/sekiju/mdl/internal/sdk/extractor/manga"
	"github.com/sekiju/mdl/internal/util"
	"regexp"
	"strings"
)

type Extractor struct {
	cookieString *string
}

func (e *Extractor) FindChapters(URL string) ([]*manga.Chapter, error) {
	//TODO implement me
	panic("implement me")
}

var re = regexp.MustCompile(`https://www.corocoro.jp/chapter/(\d*)/viewer`)

func (e *Extractor) FindChapter(URL string) (*manga.Chapter, error) {
	matches := re.FindStringSubmatch(URL)
	if len(matches) != 2 {
		return nil, manga.ErrInvalidChapterURL
	}

	req := htt.New()
	if e.cookieString != nil {
		req.SetHeader("Cookie", *e.cookieString)
	}

	res, err := req.Get(URL)
	if err != nil {
		return nil, err
	}

	html, err := res.Text()
	if err != nil {
		return nil, err
	}

	chapterMainName, err := util.ExtractStringFromHTML(html, `\"chapterMainName\":\"`, `\"`)
	if err != nil {
		return nil, err
	}

	return &manga.Chapter{
		ID:      matches[1],
		Number:  "",
		Title:   chapterMainName,
		Index:   0,
		URL:     URL,
		MangaID: "",
	}, nil
}

func (e *Extractor) FindChapterPages(chapter *manga.Chapter) ([]*manga.Page, error) {
	req := htt.New()
	if e.cookieString != nil {
		req.SetHeader("Cookie", *e.cookieString)
	}

	res, err := req.Get(chapter.URL)
	if err != nil {
		return nil, err
	}

	html, err := res.Text()
	if err != nil {
		return nil, err
	}

	jsonStr, err := util.ExtractStringFromHTML(html, `,\"pages\":`, `,\"directionRightToLeft\"`)
	if err != nil {
		return nil, err
	}

	jsonStr = strings.Replace(jsonStr, `\"`, `"`, -1)

	var result pagesResult
	if err = json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, err
	}

	pages := make([]*manga.Page, len(result))
	padRenamer := renamer.New(len(result))

	for index, page := range result {
		pages[index] = &manga.Page{
			Index:    uint(index),
			URL:      page.Src,
			Filename: padRenamer.Name(index, ".webp"),
			Decode: func(b []byte) ([]byte, error) {
				key, err := hex.DecodeString(page.Crypto.Key)
				if err != nil {
					return nil, fmt.Errorf("invalid key hex: %w", err)
				}

				iv, err := hex.DecodeString(page.Crypto.Iv)
				if err != nil {
					return nil, fmt.Errorf("invalid IV hex: %w", err)
				}

				block, err := aes.NewCipher(key)
				if err != nil {
					return nil, err
				}

				mode := cipher.NewCBCDecrypter(block, iv)
				mode.CryptBlocks(b, b)

				return b[:(len(b) - int(b[len(b)-1]))], nil
			},
		}
	}

	return pages, nil
}

func New() manga.Extractor {
	return new(Extractor)
}

func NewWithCookieString(cookieString string) manga.Extractor {
	return &Extractor{cookieString: &cookieString}
}
