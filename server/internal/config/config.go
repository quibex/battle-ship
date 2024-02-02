package config

import (

	"github.com/go-playground/validator"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string         `yaml:"env" validate:"required,oneof=local dev prod"`
	RabbitMQ RabbitMQConfig `yaml:"rabbitmq" validate:"required"`
	Postgres PostgresConfig `yaml:"postgres" validate:"required"`
}

type RabbitMQConfig struct {
	Host     string `yaml:"host" validate:"required,hostname_rfc1123"`
	Port     string `yaml:"port" validate:"required,numeric,gte=0,lte=65535"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
}

type PostgresConfig struct {
	Host     string `yaml:"host" validate:"required,hostname_rfc1123"`
	Port     string `yaml:"port" validate:"required,numeric,gte=0,lte=65535"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	DBName   string `yaml:"dbname" validate:"required"`
}

func MustLoad(configPath string) *Config {
	if configPath == "" {
		panic("config path is empty")
	}
	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		panic(err)
	}

	// cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	// cfg.RabbitMQ.Password = os.Getenv("RABBITMQ_PASSWORD")
	// fmt.Println(cfg)

	err = validator.New().Struct(cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}
