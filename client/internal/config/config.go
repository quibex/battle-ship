package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	RmqURL   string        `yaml:"rmq_url"`
	TimeOut  time.Duration `yaml:"timeout"`
	IconPath string        `yaml:"icon_path"`
}

func MustLoad(configPath string) *Config {
	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		panic(err)
	}
	return &cfg
}
