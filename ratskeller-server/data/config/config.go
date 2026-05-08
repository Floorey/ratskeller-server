package config

import "os"

type Config struct {
	Port          string
	AppEnv        string
	SessionSecret string

	PublicDir string
	AdminDir  string
	ConfigDir string
	UploadDir string
}

func (c Config) IsProduction() bool {
	return c.AppEnv == "production"
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
