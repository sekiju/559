package util

import (
	"bytes"
	"github.com/sekiju/mary/pkg/sdk/extractor/manga"
	"github.com/sekiju/rq"
	"github.com/stretchr/testify/assert"
	"testing"
)

func AssertImage(t *testing.T, decodedImageURL string, page *manga.Page) {
	res, err := rq.New().Get(page.URL)
	if err != nil {
		t.Error(err)
	}

	pageBytes, err := res.Bytes()
	if err != nil {
		t.Error(err)
	}

	if page.DescrambleFn != nil {
		pageBytes, err = (*page.DescrambleFn)(pageBytes)
		if err != nil {
			t.Error(err)
		}
	}

	res, err = rq.New().Get(decodedImageURL)
	if err != nil {
		t.Error(err)
	}

	decodedBytes, err := res.Bytes()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 0, bytes.Compare(pageBytes, decodedBytes))
}

type TestExtractedURL struct {
	extractedURL manga.ExtractedURL
	mangaID      *string
	chapterID    *string
	url          *string
}

type Opt func(*TestExtractedURL)

func NewTestExtractedURL(extractedURL manga.ExtractedURL, opts ...Opt) *TestExtractedURL {
	config := TestExtractedURL{extractedURL: extractedURL}
	for _, fn := range opts {
		fn(&config)
	}
	return &config
}

func ValidateMangaID(mangaID string) Opt {
	return func(c *TestExtractedURL) {
		c.mangaID = &mangaID
	}
}

func ValidateChapterID(chapterID string) Opt {
	return func(c *TestExtractedURL) {
		c.chapterID = &chapterID
	}
}

func ValidateURL(URL string) Opt {
	return func(c *TestExtractedURL) {
		c.url = &URL
	}
}

func (u *TestExtractedURL) Assert(t *testing.T) {
	if v := u.mangaID; v != nil {
		mangaID, err := u.extractedURL.MangaID()
		assert.NoError(t, err)
		assert.Equal(t, *v, mangaID)
	}

	if v := u.chapterID; v != nil {
		chapterID, err := u.extractedURL.ChapterID()
		assert.NoError(t, err)
		assert.Equal(t, *v, chapterID)
	}

	if v := u.url; v != nil {
		url, err := u.extractedURL.URL()
		assert.NoError(t, err)
		assert.Equal(t, *v, url)
	}
}
