package config

type Arguments struct {
	DownloadChapter string
	Session         string
	ConfigPath      string
}

type Config struct {
	Output output          `koanf:"output"`
	Sites  map[string]site `koanf:"site"`
}

type output struct {
	Dir    string       `koanf:"dir"`
	Format OutputFormat `koanf:"format"`
}

type site struct {
	Session *string
}

type OutputFormat string

const (
	AutoOutputFormat OutputFormat = "auto"
	PngOutputFormat  OutputFormat = "png"
	JpegOutputFormat OutputFormat = "jpeg"
	AvifOutputFormat OutputFormat = "avif"
	WebpOutputFormat OutputFormat = "webp"
)
