package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sekiju/htt"
	"github.com/sekiju/mdl/extractor"
	"github.com/sekiju/mdl/internal/config"
	"github.com/sekiju/mdl/internal/download"
	"github.com/sekiju/mdl/internal/util"
	"github.com/sekiju/mdl/sdk/manga"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

var version = "1.0.0"

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	isGUI := !util.IsRunningFromCLI()

	if err := run(); err != nil {
		log.Error().Err(err).Send()
		if isGUI {
			waitForInput()
		}
		os.Exit(1)
	}

	if isGUI {
		waitForInput()
	}
}

func waitForInput() {
	fmt.Println("\nPress Enter to exit...")
	_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func getChapterURLs() []string {
	if args := flag.Args(); len(args) > 0 {
		return args
	}

	fmt.Print("Please enter the chapter URL: ")
	input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.Split(strings.TrimSpace(input), " ")
}

type DownloadFunc func(page *manga.Page) error

func run() error {
	cfg, err := config.New()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if cfg.Application.CheckForUpdates && checkForUpdates() != nil {
		return err
	}

	for _, chapterURL := range getChapterURLs() {
		parsedURL, err := url.Parse(chapterURL)
		if err != nil {
			return fmt.Errorf("invalid chapter URL: %v", err)
		}

		ext, err := extractor.NewExtractor(cfg, parsedURL.Hostname())
		if err != nil {
			return err
		}

		chapter, err := ext.FindChapter(chapterURL)
		if err != nil {
			return err
		}

		log.Info().Msg("Extracting pages...")

		pages, err := ext.FindChapterPages(chapter)
		if err != nil {
			return err
		}

		if len(pages) == 0 {
			return errors.New("no pages found to download")
		}

		outputDir := filepath.Join(cfg.Output.Dir, chapter.ID)

		if cfg.Output.CleanDir {
			if stat, err := os.Stat(outputDir); err == nil && stat.IsDir() {
				if err = os.RemoveAll(outputDir); err != nil {
					return err
				}
			}
		}

		if err = os.MkdirAll(outputDir, os.ModePerm); err != nil {
			return err
		}

		log.Info().Msgf("Output: %s | Pages: %d", outputDir, len(pages))

		start := time.Now()

		var downloadFn DownloadFunc
		if cfg.Output.Format == "auto" {
			downloadFn = func(page *manga.Page) error {
				return download.Bytes(outputDir, page)
			}
		} else {
			downloadFn = func(page *manga.Page) error {
				return download.WithEncode(outputDir, cfg.Output.Format, page)
			}
		}

		if err = downloadPages(pages, cfg.Download.ConcurrentProcesses, downloadFn); err != nil {
			return fmt.Errorf("downloading pages: %w", err)
		}

		log.Info().Msgf("Download completed in %s", time.Since(start))
	}

	return nil
}

func checkForUpdates() error {
	log.Trace().Msgf("Current version: %s | Checking for updates...", version)

	res, err := htt.New().Get("https://api.github.com/repos/sekiju/mdl/tags")
	if err != nil {
		return err
	}

	var tags []map[string]interface{}
	if err = res.JSON(&tags); err != nil {
		return err
	}

	currentVersion := semver.MustParse(version)

	var versions []*semver.Version
	for _, tag := range tags {
		if v, err := semver.NewVersion(tag["name"].(string)); err == nil && v.GreaterThan(currentVersion) {
			versions = append(versions, v)
		}
	}

	if len(versions) == 0 {
		return nil
	}

	sort.Sort(semver.Collection(versions))

	log.Info().Msgf("New version available: %s - download release from: https://github.com/sekiju/mdl/releases", versions[len(versions)-1].String())

	return nil
}

func downloadPages(pages []*manga.Page, workers int, download DownloadFunc) error {
	ch := make(chan *manga.Page, len(pages))
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for page := range ch {
				if err := download(page); err != nil {
					log.Error().Err(err).Msgf("failed to download #%d page", page.Index)
				}
			}
		}()
	}

	go func() {
		for _, page := range pages {
			ch <- page
		}
		close(ch)
	}()

	wg.Wait()
	return nil
}
