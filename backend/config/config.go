package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	Port int `mapstructure:"port"`
}

var (
	cfg  *Config
	once sync.Once
)

func Load() (*Config, error) {
	var err error
	once.Do(func() {
		cfg = &Config{
			Port: 8080,
		}

		if err = viper.Unmarshal(cfg); err != nil {
			err = fmt.Errorf("unmarshal config: %w", err)
			return
		}

		if err = cfg.validate(); err != nil {
			err = fmt.Errorf("validate config: %w", err)
		}
	})
	return cfg, err
}

func (c *Config) validate() error {
	if c.Port == 0 {
		return fmt.Errorf("Port is required")
	}
	return nil
}

func Get() *Config {
	return cfg
}
