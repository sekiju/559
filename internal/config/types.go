package config

type Arguments struct {
	DownloadChapter string
	Session         string
	ConfigPath      string
}

type Config struct {
	OutputDir string          `koanf:"output_dir"`
	Sites     map[string]site `koanf:"site"`
}

type site struct {
	Session *string
}
