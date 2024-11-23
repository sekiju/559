package config

type arguments struct {
	Cookie     string
	ConfigPath string
}

type Config struct {
	PrimaryCookie *string

	Application application     `koanf:"application"`
	Output      output          `koanf:"output"`
	Download    download        `koanf:"download"`
	Sites       map[string]site `koanf:"site"`
}

type application struct {
	CheckForUpdates bool `koanf:"check_for_updates"`
}

type output struct {
	Dir      string       `koanf:"dir"`
	CleanDir bool         `koanf:"clean_dir"`
	Format   OutputFormat `koanf:"format"`
}

type download struct {
	PreloadNextChapters int `koanf:"preload_next_chapters"`
	PageBatchSize       int `koanf:"page_batch_size"`
}

type site struct {
	CookieString *string `koanf:"cookie_string"`
}

type OutputFormat string

const (
	AutoOutputFormat OutputFormat = "auto"
	PngOutputFormat  OutputFormat = "png"
	JpegOutputFormat OutputFormat = "jpeg"
	AvifOutputFormat OutputFormat = "avif"
	WebpOutputFormat OutputFormat = "webp"
)
