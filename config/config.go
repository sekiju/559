package config

import (
	"errors"
	"github.com/knadh/koanf/parsers/hcl"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"
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

func Load(filepath string) {
	k := koanf.NewWithConf(koanf.Conf{
		Delim:       ".",
		StrictMerge: false,
	})

	if err := k.Load(file.Provider(filepath), hcl.Parser(true)); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Fatal().Err(err).Send()
		} else {
			log.Info().Str("filepath", filepath).Msg("Config doesn't exist. Using default configuration")
		}
	}

	if err := k.Unmarshal("", &Params); err != nil {
		log.Fatal().Err(err).Send()
	}
}
