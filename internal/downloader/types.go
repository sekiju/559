package downloader

import (
	"github.com/sekiju/mdl/internal/config"
	"github.com/sekiju/mdl/sdk/manga"
	"sync"
)

type (
	DownloadPageFunc func(dir string, page *manga.Page) error

	NewExtractorFunc func(hostname string) (manga.Extractor, error)

	Downloader struct {
		ch               chan *QueueInfo
		wg               sync.WaitGroup
		cleanDestination bool
		downloadDir      string
		batchSize        int
		downloadPage     DownloadPageFunc
		newExtractor     NewExtractorFunc
	}

	NewDownloaderOptions struct {
		BatchSize        int
		Directory        string
		CleanDestination bool
		OutputFileFormat config.OutputFileFormat
		NewExtractor     NewExtractorFunc
	}

	QueueInfo struct {
		URL       string
		ChapterID string
		Pages     []*manga.Page
	}
)
