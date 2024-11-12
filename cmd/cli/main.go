package main

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sekiju/mary/internal/config"
	"github.com/sekiju/mary/internal/download"
	"github.com/sekiju/mary/internal/extractor"
	"github.com/sekiju/mary/pkg/sdk/extractor/manga"
	"net/url"
	"os"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	if err := run(); err != nil {
		log.Fatal().Err(err).Send()
	}
}

func downloadChapter(ext manga.Provider, cfg *config.Config, chapterURL string) error {
	mangaID, err := ext.ExtractMangaID(chapterURL)
	if err != nil && !errors.Is(err, manga.ErrURLIsID) {
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

	return nil
}

func run() error {
	args, err := config.ParseArguments()
	if err != nil {
		return err
	}

	cfg, err := config.New(args.ConfigPath)
	if err != nil {
		return err
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
