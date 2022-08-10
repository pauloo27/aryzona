package config

import (
	"os"

	"github.com/ghodss/yaml"
)

var Config = new(BotConfig)

func Load() (*BotConfig, error) {
	file, err := os.ReadFile("config.yml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, Config)
	return Config, err
}
