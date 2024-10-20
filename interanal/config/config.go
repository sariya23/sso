package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type ConfigDataBase struct {
	Port     string `yaml:"port"`
	Host     string `yaml:"host"`
	DBName   string `yaml:"db_name"`
	User     string `yaml:"username"`
	Password string `yaml:"password"`
}

type Config struct {
	Env      string        `yaml:"env" env-required:"true"`
	TokenTTL time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC     GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env-default:"8080"`
	Timeout time.Duration `yaml:"timeout" env-default:"1h"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}
	if _, err := os.Stat(path); err != nil {
		panic("file does not exist: " + path)
	}

	var config Config

	if err := cleanenv.ReadConfig(path, &config); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &config
}

func MustLoadDBConfig(path string) *ConfigDataBase {
	var cfg ConfigDataBase
	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		panic(err)
	}
	return &cfg
}

// Priority: flag > env > default
func fetchConfigPath() string {
	var path string
	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}
	return path
}
