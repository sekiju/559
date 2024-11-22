package util

import (
	"bytes"
	"github.com/sekiju/htt"
	"github.com/sekiju/mdl/internal/sdk/extractor/manga"
	"github.com/stretchr/testify/assert"
	"testing"
)

func AssertImage(t *testing.T, expectedImageURL string, page *manga.Page) {
	res, err := htt.New().Get(page.URL)
	if err != nil {
		t.Error(err)
	}

	actual, err := res.Bytes()
	if err != nil {
		t.Error(err)
	}

	if page.Decode != nil {
		actual, err = page.Decode(actual)
		if err != nil {
			t.Error(err)
		}
	}

	res, err = htt.New().Get(expectedImageURL)
	if err != nil {
		t.Error(err)
	}

	expected, err := res.Bytes()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 0, bytes.Compare(actual, expected))
}
