package config

import (
	"log/slog"
	"os"

	"github.com/ghodss/yaml"
	"github.com/joho/godotenv"
	env "github.com/qiangxue/go-env"
)

var Config = new(BotConfig)

const (
	yamlFilePath = "config.yml"
	envFilePath  = ".env"
	envPrefix    = "BOT_"
)

func Load() error {
	if _, err := os.Stat(yamlFilePath); err == nil {
		slog.Info("Yaml file located, loading", "path", yamlFilePath)
		err = loadConfigFromYaml()
		if err == nil {
			slog.Info("Config loaded")
		}
		return err
	}

	slog.Warn("Yaml config file not found, using environment...")

	if _, err := os.Stat(envFilePath); err == nil {
		slog.Info("Found env file, loading...", "path", envFilePath)
		err := godotenv.Load(envFilePath)
		if err != nil {
			slog.Error("Cannot load env file", err)
			os.Exit(1)
		}
	}

	loader := env.New(envPrefix, nil)

	return loader.Load(Config)
}

func loadConfigFromYaml() error {
	file, err := os.ReadFile(yamlFilePath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(file, Config)
}
