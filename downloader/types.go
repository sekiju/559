package downloader

import (
	"github.com/sekiju/mdl/sdk/manga"
	"sync"
)

type (
	downloadPageFunc func(dir string, page *manga.Page) error

	Downloader struct {
		ch           chan *queueInfo
		wg           sync.WaitGroup
		downloadPage downloadPageFunc
	}

	queueInfo struct {
		URL       string
		ChapterID string
		Pages     []*manga.Page
	}
)
