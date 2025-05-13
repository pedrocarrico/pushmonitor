package config

import (
	"fmt"
	"os"

	"github.com/pedrocarrico/pushmonitor/internal/pushtest"
	"gopkg.in/yaml.v2"
)

type Config struct {
	PushTests []pushtest.PushTest `yaml:"push_tests"`
	Logging   LogConfig           `yaml:"logging"`
	Timeout   int                 `yaml:"timeout"`
}

type LogConfig struct {
	File  string `yaml:"file"`
	Level string `yaml:"level"`
}

func (c *Config) Load(configLocations ...string) error {
	if configLocations == nil {
		configLocations = []string{
			"/etc/pushmonitor/config.yaml",
			"config/config.yaml",
		}
	}

	var configData []byte
	var err error

	for _, location := range configLocations {
		configData, err = os.ReadFile(location)
		if err == nil {
			break
		}
	}

	if err != nil {
		return fmt.Errorf("error reading config file from any location: %v", err)
	}

	err = yaml.Unmarshal(configData, c)
	if err != nil {
		return fmt.Errorf("error parsing config file: %v", err)
	}

	return nil
}
