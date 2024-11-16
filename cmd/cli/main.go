package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sekiju/mary/internal/config"
	"github.com/sekiju/mary/internal/download"
	"github.com/sekiju/mary/internal/extractor"
	"github.com/sekiju/mary/internal/util"
	"github.com/sekiju/mary/pkg/sdk/extractor/manga"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	isGUI := util.IsRunningFromExplorer()

	if err := run(isGUI); err != nil {
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
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func getChapterURLFromUser() string {
	fmt.Print("Please enter the chapter URL: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func run(isGUI bool) error {
	args, err := config.ParseArguments()
	if err != nil {
		return err
	}

	if args.DownloadChapter == "" {
		if isGUI {
			args.DownloadChapter = getChapterURLFromUser()
		}
		if args.DownloadChapter == "" {
			return errors.New("no chapter URL provided. Use --download-chapter flag or enter URL when prompted")
		}
	}

	cfg, err := config.New(args.ConfigPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}

		log.Info().Msg("No configuration file found, using defaults")
	}

	chapterURL, err := url.Parse(args.DownloadChapter)
	if err != nil {
		return fmt.Errorf("invalid chapter URL: %v", err)
	}

	ext, err := extractor.CreateProvider(cfg, args, chapterURL.Hostname())
	if err != nil {
		return err
	}

	if err = downloadChapter(ext, cfg, args.DownloadChapter); err != nil {
		return err
	}

	return nil
}

func downloadChapter(ext manga.Provider, cfg *config.Config, chapterURL string) error {
	mangaID, err := ext.ExtractMangaID(chapterURL)
	if err != nil {
		return err
	}

	chapter, err := ext.FindChapter(mangaID)
	if err != nil {
		return err
	}

	pages, err := ext.ExtractPages(chapter)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(cfg.Output.Dir, 0755); err != nil {
		return err
	}

	log.Info().Msgf("Output directory: %s", cfg.Output.Dir)
	log.Info().Msgf("Total pages: %d", len(pages))

	start := time.Now()

	var downloadFn = func(page *manga.Page) error {
		return download.Bytes(cfg.Output.Dir, page)
	}

	if cfg.Output.Format != "auto" {
		downloadFn = func(page *manga.Page) error {
			return download.WithEncode(cfg.Output.Dir, cfg.Output.Format, page)
		}
	}

	for _, page := range pages {
		if err = downloadFn(page); err != nil {
			return err
		}
	}

	timeElapsed := time.Since(start)
	log.Info().Msgf("Download completed in %s", timeElapsed)

	return nil
}
