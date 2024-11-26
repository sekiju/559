package downloader

import (
	"github.com/rs/zerolog/log"
	"github.com/sekiju/mdl/internal/config"
	"github.com/sekiju/mdl/sdk/manga"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func (d *Downloader) Queue(URL string) {
	qi := &QueueInfo{URL: URL}

	log.Debug().Msgf("New queue item with URL: %s", URL)

	d.wg.Add(1)
	d.ch <- qi
}

func (d *Downloader) Stop() {
	close(d.ch)
	d.wg.Wait()
}

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
	semaphore := make(chan struct{}, d.batchSize)

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
				return
			}
		}(page)
	}

	wg.Wait()

	return nil
}

func (d *Downloader) run() {
	semaphore := make(chan struct{}, 1)

	for queueInfo := range d.ch {
		semaphore <- struct{}{}
		go func(queueInfo *QueueInfo) {
			defer func() {
				<-semaphore
				d.wg.Done()
			}()

			log.Info().Str("url", queueInfo.URL).Msg("Downloading next chapter in queue")
			start := time.Now()

			parsedURL, err := url.Parse(queueInfo.URL)
			if err != nil {
				log.Error().Str("url", queueInfo.URL).Err(err).Msg("Invalid chapter URL")
				return
			}

			ext, err := d.newExtractor(parsedURL.Hostname())
			if err != nil {
				log.Error().Err(err).Str("url", queueInfo.URL).Send()
				return
			}

			chapter, err := ext.FindChapter(queueInfo.URL)
			if err != nil {
				log.Error().Str("url", queueInfo.URL).Err(err).Msg("Failed to find chapter")
				return
			}

			queueInfo.ChapterID = chapter.ID

			pages, err := ext.FindChapterPages(chapter)
			if err != nil {
				log.Error().Str("chapterId", chapter.ID).Err(err).Msg("Failed to find chapter pages")
				return
			}

			queueInfo.Pages = pages

			if err = d.downloadImages(queueInfo); err != nil {
				log.Error().Str("chapterId", queueInfo.ChapterID).Err(err).Msg("Failed to download chapter")
				return
			}

			log.Info().Str("chapterId", queueInfo.ChapterID).Str("duration", time.Since(start).String()).Msg("Download complete")

		}(queueInfo)
	}
}

func NewDownloader(opts *NewDownloaderOptions) *Downloader {
	var downloadFunc DownloadPageFunc
	if opts.OutputFileFormat == config.AutoOutputFormat {
		downloadFunc = func(dir string, page *manga.Page) error {
			r, err := getReader(page)
			if err != nil {
				return err
			}

			return saveFile(dir, page.Filename, r)
		}
	} else {
		downloadFunc = func(dir string, page *manga.Page) error {
			r, err := getReader(page)
			if err != nil {
				return err
			}

			return saveEncodedImage(dir, page.Filename, opts.OutputFileFormat, r)
		}
	}

	d := &Downloader{
		ch:               make(chan *QueueInfo),
		downloadPage:     downloadFunc,
		newExtractor:     opts.NewExtractor,
		batchSize:        opts.BatchSize,
		cleanDestination: opts.CleanDestination,
		downloadDir:      opts.Directory,
	}

	go d.run()

	return d
}
