package config

import (
	"github.com/knadh/koanf/parsers/hcl"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func New(path string) (*Config, error) {
	k := koanf.NewWithConf(koanf.Conf{
		Delim:       ".",
		StrictMerge: true,
	})

	c := Config{OutputDir: "downloads"}

	if err := k.Load(file.Provider(path), hcl.Parser(true)); err != nil {
		return nil, err
	}

	if err := k.Unmarshal("", &c); err != nil {
		return nil, err
	}

	return &c, nil
}
