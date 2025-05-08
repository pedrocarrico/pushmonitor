package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/pedrocarrico/pushmonitor/internal/pushtest"
	"gopkg.in/yaml.v2"
)

type Config struct {
	PushTests []pushtest.PushTest `yaml:"push_tests"`
	Logging   LogConfig           `yaml:"logging"`
	PIDFile   string              `yaml:"pid_file"`
	Timeout   int                 `yaml:"timeout"`
	mutex     sync.RWMutex
}

type LogConfig struct {
	File  string `yaml:"file"`
	Level string `yaml:"level"`
}

func (c *Config) Load() error {
	configLocations := []string{
		"/etc/pushmonitor/config.yaml",
		"config/config.yaml",
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

func (c *Config) Reload() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err := c.Load(); err != nil {
		return fmt.Errorf("failed to reload configuration: %v", err)
	}

	return nil
}
