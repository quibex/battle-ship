package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	RmqURL string `yaml:"rmq_url"`
}

func MustLoad(configPath string) *Config { 
	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		panic(err)
	}
	return &cfg
}