package downloader

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/sekiju/mdl/sdk/manga"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type (
	DownloadPageFunc func(dir string, page *manga.Page) error

	NewExtractorFunc func(hostname string) (manga.Extractor, error)

	Downloader struct {
		ch chan *QueueInfo
		wg sync.WaitGroup

		cleanDestination      bool
		downloadDir           string
		maxPageBatchSize      int
		maxPrefetchedChapters int

		downloadPage DownloadPageFunc
		newExtractor NewExtractorFunc
	}

	NewDownloaderOptions struct {
		MaximumPrefetchedChapters int
		MaxPageBatchSize          int
		DownloadDir               string
		CleanDestination          bool

		DownloadPage DownloadPageFunc
		NewExtractor NewExtractorFunc
	}

	QueueInfo struct {
		URL       string
		ChapterID string
		Pages     []*manga.Page
		Status    QueueStatus
	}

	QueueStatus int
)

const (
	QueueStatusNotStarted QueueStatus = iota
	QueueStatusPrefetched
)

func (d *Downloader) downloadImages(qi *QueueInfo) error {
	destination := filepath.Join(d.downloadDir, qi.ChapterID)

	if _, err := os.Stat(destination); err == nil && d.cleanDestination {
		if err = os.RemoveAll(destination); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(destination, os.ModePerm); err != nil {
		log.Error().Err(err).Msgf("Failed to create download directory for chapter %s", qi.ChapterID)
		return err
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, d.maxPageBatchSize)
	errorCount := 0
	var errorMu sync.Mutex

	for _, page := range qi.Pages {
		semaphore <- struct{}{}
		wg.Add(1)
		go func(page *manga.Page) {
			defer func() {
				<-semaphore
				wg.Done()
			}()

			if err := d.downloadPage(destination, page); err != nil {
				log.Error().Err(err).Msgf("Failed to download page #%d", page.Index)
				errorMu.Lock()
				errorCount++
				errorMu.Unlock()
				return
			}
		}(page)
	}

	wg.Wait()

	if errorCount > 0 {
		return fmt.Errorf("failed to download %d pages", errorCount)
	}

	return nil
}

func (d *Downloader) Queue(URL string) {
	qi := &QueueInfo{
		URL:    URL,
		Status: QueueStatusNotStarted,
	}

	log.Debug().Msgf("New queue item for URL: %s", URL)

	d.wg.Add(1)
	d.ch <- qi
}

func (d *Downloader) GracefulStop() {
	d.wg.Wait()
}

func (d *Downloader) run() {
	prefetchSemaphore := make(chan struct{}, d.maxPrefetchedChapters)
	downloadSemaphore := make(chan struct{}, 1)

	for queueInfo := range d.ch {
		prefetchSemaphore <- struct{}{}
		go func(queueInfo *QueueInfo) {
			defer func() { <-prefetchSemaphore }()

			if queueInfo.Status == QueueStatusPrefetched {
				downloadSemaphore <- struct{}{}
				defer func() {
					<-downloadSemaphore
					d.wg.Done()
				}()

				log.Info().Str("url", queueInfo.URL).Msg("Downloading next chapter in queue")
				start := time.Now()

				if err := d.downloadImages(queueInfo); err != nil {
					log.Error().Str("chapterId", queueInfo.ChapterID).Err(err).Msg("Failed to download chapter")
					return
				}

				log.Info().Str("chapterId", queueInfo.ChapterID).Str("duration", time.Since(start).String()).Msg("Download complete")
			}

			if queueInfo.Status == QueueStatusNotStarted {
				log.Info().Str("url", queueInfo.URL).Msg("Preloading next chapter pages")

				parsedURL, err := url.Parse(queueInfo.URL)
				if err != nil {
					log.Error().Str("url", queueInfo.URL).Err(err).Msg("Invalid chapter URL")
					return
				}

				ext, err := d.newExtractor(parsedURL.Hostname())
				if err != nil {
					log.Error().Str("url", queueInfo.URL).Msg("Website unsupported")
					return
				}

				chapter, err := ext.FindChapter(queueInfo.URL)
				if err != nil {
					log.Error().Str("url", queueInfo.URL).Err(err).Msg("Failed to find chapter")
					return
				}

				pages, err := ext.FindChapterPages(chapter)
				if err != nil {
					log.Error().Str("chapterId", chapter.ID).Err(err).Msg("Failed to find chapter pages")
					return
				}

				queueInfo.Status = QueueStatusPrefetched
				queueInfo.Pages = pages
				queueInfo.ChapterID = chapter.ID

				d.ch <- queueInfo
			}
		}(queueInfo)
	}
}

func NewDownloader(opts *NewDownloaderOptions) *Downloader {
	d := &Downloader{
		ch:                    make(chan *QueueInfo),
		downloadPage:          opts.DownloadPage,
		newExtractor:          opts.NewExtractor,
		maxPageBatchSize:      opts.MaxPageBatchSize,
		maxPrefetchedChapters: opts.MaximumPrefetchedChapters,
		cleanDestination:      opts.CleanDestination,
		downloadDir:           opts.DownloadDir,
	}

	go d.run()

	return d
}
