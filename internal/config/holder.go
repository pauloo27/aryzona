package config

import (
	"os"

	"github.com/Pauloo27/logger"
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
		logger.Infof("Found %s, loading...", yamlFilePath)
		err = loadConfigFromYaml()
		if err == nil {
			logger.Success("Config loaded")
		}
		return err
	}

	logger.Warnf("Config file %s not found, loading from environment...", yamlFilePath)

	if _, err := os.Stat(envFilePath); err == nil {
		logger.Infof("Found %s file, loading...", envFilePath)
		err := godotenv.Load(envFilePath)
		if err != nil {
			logger.Fatal(err)
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
