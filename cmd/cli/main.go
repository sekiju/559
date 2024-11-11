package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sekiju/mary/internal/config"
	"github.com/sekiju/mary/internal/download"
	"github.com/sekiju/mary/pkg/manga/comic_walker"
	"github.com/sekiju/mary/pkg/manga/giga_viewer"
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

type Arguments struct {
	DownloadChapter string
	Session         string
	ConfigPath      string
}

type ProviderFactory func(session *string) manga.Provider

func parseArgs() Arguments {
	args := Arguments{}

	flag.StringVar(&args.DownloadChapter, "download-chapter", "", "URL of the chapter to download")
	flag.StringVar(&args.Session, "session", "", "Session token for the current service")
	flag.StringVar(&args.ConfigPath, "config", "config.hcl", "Path to the config file (default: config.yaml)")
	flag.Parse()

	return args
}

func getProviderFactories() map[string]ProviderFactory {
	return map[string]ProviderFactory{
		"comic-walker.com": func(session *string) manga.Provider {
			return comic_walker.New()
		},
		"shonenjumpplus.com": func(session *string) manga.Provider {
			if session != nil {
				return giga_viewer.NewWithSession("shonenjumpplus.com", *session)
			}
			return giga_viewer.New("shonenjumpplus.com")
		},
	}
}

func createExtractors(cfg *config.Config, args Arguments, hostname string) (manga.Provider, error) {
	factories := getProviderFactories()
	factory, exists := factories[hostname]
	if !exists {
		return nil, fmt.Errorf("extractor not found for hostname: %s", hostname)
	}

	var session *string
	if args.Session != "" {
		session = &args.Session
	} else if site, exists := cfg.Sites[hostname]; exists && site.Session != nil {
		session = site.Session
	}

	return factory(session), nil
}

func downloadChapter(ext manga.Provider, outputDir, chapterURL string) error {
	mangaID, err := ext.ExtractMangaID(chapterURL)
	if err != nil && !errors.Is(err, manga.ErrURLeqID) {
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

	dir := fmt.Sprintf(outputDir)
	if err = os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	for _, page := range pages {
		if err = download.Bytes(dir, page); err != nil {
			return err
		}
	}

	return nil
}

func run() error {
	args := parseArgs()

	if args.DownloadChapter == "" {
		return errors.New("no chapter URL provided. Use --download-chapter flag")
	}

	cfg, err := config.New(args.ConfigPath)
	if err != nil {
		return err
	}

	chapterURL, err := url.Parse(args.DownloadChapter)
	if err != nil {
		return fmt.Errorf("invalid chapter URL: %v", err)
	}

	ext, err := createExtractors(cfg, args, chapterURL.Hostname())
	if err != nil {
		return err
	}

	if err = downloadChapter(ext, cfg.OutputDir, args.DownloadChapter); err != nil {
		return err
	}

	return nil
}
