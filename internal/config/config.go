package config

import (
	"errors"
	"flag"
	"github.com/knadh/koanf/parsers/hcl"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func New(path string) (*Config, error) {
	k := koanf.NewWithConf(koanf.Conf{
		Delim:       ".",
		StrictMerge: true,
	})

	c := Config{Output: output{
		Dir:    "downloads",
		Format: AutoOutputFormat,
	}}

	if err := k.Load(file.Provider(path), hcl.Parser(true)); err != nil {
		return nil, err
	}

	if err := k.Unmarshal("", &c); err != nil {
		return nil, err
	}

	return &c, nil
}

func ParseArguments() (*Arguments, error) {
	args := Arguments{}

	flag.StringVar(&args.DownloadChapter, "download-chapter", "", "URL of the chapter to download")
	flag.StringVar(&args.Session, "session", "", "Session token for the current service")
	flag.StringVar(&args.ConfigPath, "config", "config.hcl", "Path to the config file (default: config.yaml)")
	flag.Parse()

	if args.DownloadChapter == "" {
		return nil, errors.New("no chapter URL provided. Use --download-chapter flag")
	}

	return &args, nil
}
