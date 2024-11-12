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

func AssertMangaID(t *testing.T, extractedURL manga.ExtractedURL, mangaID any) {
	extractedMangaID, err := extractedURL.MangaID()
	assert.NoError(t, err)
	assert.Equal(t, mangaID, extractedMangaID)
}

func AssertMangaAndChapterID(t *testing.T, extractedURL manga.ExtractedURL, mangaID, chapterID any) {
	extractedMangaID, err := extractedURL.MangaID()
	assert.NoError(t, err)
	assert.Equal(t, mangaID, extractedMangaID)

	extractedChapterID, err := extractedURL.ChapterID()
	assert.NoError(t, err)
	assert.Equal(t, chapterID, extractedChapterID)
}
