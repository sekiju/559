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
	"github.com/sekiju/mdl/internal/config"
	"github.com/sekiju/mdl/internal/download"
	"github.com/sekiju/mdl/internal/extractor"
	"github.com/sekiju/mdl/internal/sdk/extractor/manga"
	"github.com/sekiju/mdl/internal/util"
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
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadBytes('\n')
}

func getChapterURLsFromUser() []string {
	args := flag.Args()
	if len(args) == 0 {
		fmt.Print("Please enter the chapter URL: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		return strings.Split(strings.TrimSpace(input), " ")
	}

	return args
}

func run() error {
	cfg, err := config.New()
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		log.Info().Msg("No configuration file found, using defaults")
	}

	if cfg.Application.CheckForUpdates {
		if err = isOutdatedVersion(); err != nil {
			return err
		}
	}

	chapterURLs := getChapterURLsFromUser()

	for _, chapterURL := range chapterURLs {
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

		if err = os.MkdirAll(outputDir, 0755); err != nil {
			return err
		}

		log.Info().Msgf("Output directory: %s", outputDir)
		log.Info().Msgf("Total pages: %d", len(pages))

		start := time.Now()

		var downloadFn func(page *manga.Page) error
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

func isOutdatedVersion() error {
	log.Trace().Msgf("Current version: %s", version)

	res, err := htt.New().Get("https://api.github.com/repos/sekiju/mdl/tags")
	if err != nil {
		return err
	}

	var tagsResult []map[string]interface{}
	if err = res.JSON(&tagsResult); err != nil {
		return err
	}

	currentVersion := semver.MustParse(version)

	var versions []*semver.Version
	for _, tag := range tagsResult {
		v, err := semver.NewVersion(tag["name"].(string))
		if err != nil {
			return err
		}

		if v.GreaterThan(currentVersion) {
			versions = append(versions, v)
		}
	}

	if len(versions) == 0 {
		return nil
	}

	sort.Sort(semver.Collection(versions))

	latestVersion := versions[len(versions)-1]

	log.Info().Msgf("New version available: %s - download release from: https://github.com/sekiju/mdl/releases", latestVersion.String())

	return nil
}

func downloadPages(pages []*manga.Page, workers int, downloadFn func(page *manga.Page) error) error {
	ch := make(chan *manga.Page, len(pages))
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for page := range ch {
				if err := downloadFn(page); err != nil {
					log.Error().Err(err).Msgf("downloading page %d", page.Index)
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
