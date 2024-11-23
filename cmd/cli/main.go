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
	"github.com/sekiju/mdl/internal/downloader"
	"github.com/sekiju/mdl/internal/util"
	"github.com/sekiju/mdl/sdk/manga"
	"os"
	"sort"
	"strings"
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

func run() error {
	cfg, err := config.New()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if cfg.Application.CheckUpdates && checkForUpdates() != nil {
		return err
	}

	var downloadFunc downloader.DownloadPageFunc
	if cfg.Output.FileFormat == "auto" {
		downloadFunc = func(dir string, page *manga.Page) error {
			return downloader.Bytes(dir, page)
		}
	} else {
		downloadFunc = func(dir string, page *manga.Page) error {
			return downloader.WithEncode(dir, cfg.Output.FileFormat, page)
		}
	}

	loader := downloader.NewDownloader(&downloader.NewDownloaderOptions{
		BatchSize:        cfg.Application.MaxParallelDownloads,
		Directory:        cfg.Output.Directory,
		DownloadPage:     downloadFunc,
		CleanDestination: cfg.Output.CleanOnStart,
		NewExtractor: func(hostname string) (manga.Extractor, error) {
			return extractor.NewExtractor(cfg, hostname)
		},
	})

	chapterURLs := getChapterURLs()

	for _, chapterURL := range chapterURLs {
		loader.Queue(chapterURL)
	}

	loader.GracefulStop()

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

	log.Info().Msgf("NewDownloader version available: %s - download release from: https://github.com/sekiju/mdl/releases", versions[len(versions)-1].String())

	return nil
}
