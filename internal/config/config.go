package config

import (
	"flag"
	"fmt"
	"github.com/knadh/koanf/parsers/hcl"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/sekiju/mdl/internal/ptr"
)

func New() (*Config, error) {
	args, err := parseArguments()
	if err != nil {
		return nil, err
	}

	k := koanf.NewWithConf(koanf.Conf{
		Delim:       ".",
		StrictMerge: true,
	})

	c := Config{
		PrimaryCookie: ptr.String(args.Cookie),
		Application: application{
			CheckForUpdates: true,
		},
		Output: output{
			Dir:      "downloads",
			CleanDir: false,
			Format:   AutoOutputFormat,
		},
		Download: download{
			PreloadNextChapters: 2,
			PageBatchSize:       4,
		},
	}

	if err = k.Load(file.Provider(args.ConfigPath), hcl.Parser(true)); err != nil {
		return &c, err
	}

	if err = k.Unmarshal("", &c); err != nil {
		return &c, err
	}

	return &c, nil
}

func parseArguments() (*arguments, error) {
	args := arguments{}

	flag.StringVar(&args.Cookie, "cookie", "", "Cookie string for the current session")
	flag.StringVar(&args.ConfigPath, "config", "config.hcl", "Path to the config file (default: config.yaml)")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: mdl [OPTIONS] chapterURL [chapterURLs...]\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	return &args, nil
}
