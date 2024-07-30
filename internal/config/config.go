package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env           string     `yaml:  env-default:"local"`
	GRPC          GRPCConfig `yaml:"grpc" env-required:"true"`
	MigrationPath string
	TokenTTL      time.Duration `yaml:tokenTTL env-default:"1h"`
}

type GRPCConfig struct {
	Port    int           `yaml:port`
	Timeout time.Duration `yaml:timeout`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config path does not exist:" + configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("config file is empty:" + configPath)
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
