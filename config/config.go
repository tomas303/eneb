package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type DataConfig struct {
	Storepath string
}

type ServerConfig struct {
	Port int
}

type Config struct {
	Data   DataConfig
	Server ServerConfig
}

func (cfg *Config) defaultValues() {
	cfg.Data = DataConfig{Storepath: "."}
	cfg.Server = ServerConfig{Port: 8080}
}

func New(filename string) (*Config, error) {
	var c Config
	c.defaultValues()
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return &c, nil
		}
		return &c, err
	}
	_, err = toml.DecodeFile(filename, &c)
	if err != nil {
		return &c, err
	}
	return &c, nil
}
