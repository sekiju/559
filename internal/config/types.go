package config

type arguments struct {
	Cookie       string
	ConfigPath   string
	ListChapters bool
}

type Config struct {
	PrimaryCookie    *string
	ListChaptersMode bool

	Application application     `koanf:"application"`
	Output      output          `koanf:"output"`
	Sites       map[string]site `koanf:"site"`
}

type application struct {
	CheckUpdates         bool `koanf:"check_updates"`
	MaxParallelDownloads int  `koanf:"max_parallel_downloads"`
}

type output struct {
	Directory    string           `koanf:"directory"`
	CleanOnStart bool             `koanf:"clean_on_start"`
	FileFormat   OutputFileFormat `koanf:"file_format"`
}

type site struct {
	Cookie *string `koanf:"cookie"`
}

type OutputFileFormat string

const (
	AutoOutputFormat OutputFileFormat = "auto"
	PngOutputFormat  OutputFileFormat = "png"
	JpegOutputFormat OutputFileFormat = "jpeg"
	AvifOutputFormat OutputFileFormat = "avif"
	WebpOutputFormat OutputFileFormat = "webp"
)
