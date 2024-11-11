package config

type Config struct {
	OutputDir string
	Sites     map[string]site `koanf:"site"`
}

type site struct {
	Session *string
}
