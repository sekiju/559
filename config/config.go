package config

import (
	"errors"
	"flag"
	"github.com/knadh/koanf/parsers/hcl"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"
	"github.com/sekiju/mdl/constant"
	"os"
)

var Params = config{
	Application: application{
		CheckUpdates:         true,
		MaxParallelDownloads: 4,
	},
	Output: output{
		Directory:    "downloads",
		CleanOnStart: false,
		FileFormat:   AutoOutputFormat,
	},
}

func init() {
	rootFlags := flag.NewFlagSet(constant.MDL, flag.ExitOnError)
	primaryCookie := rootFlags.String("cookie", "", "Cookie string for the current session")
	configPath := rootFlags.String("config", "config.hcl", "Path to the config file")

	if len(os.Args) > 1 && os.Args[1] == "chapters" {
		Params.ListChaptersMode = true
		os.Args = append(os.Args[:1], os.Args[2:]...)
	}

	if err := rootFlags.Parse(os.Args[1:]); err != nil {
		log.Fatal().Err(err).Send()
	}

	Params.DownloadChapters = rootFlags.Args()
	Params.PrimaryCookie = primaryCookie

	k := koanf.NewWithConf(koanf.Conf{
		Delim:       ".",
		StrictMerge: false,
	})

	if err := k.Load(file.Provider(*configPath), hcl.Parser(true)); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Fatal().Err(err).Send()
		} else {
			log.Info().Str("filename", *configPath).Msg("Config doesn't exist. Using default configuration")
		}
	}

	if err := k.Unmarshal("", &Params); err != nil {
		log.Fatal().Err(err).Send()
	}
}
