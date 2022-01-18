package config

import (
	"os"

	yaml "gopkg.in/yaml.v2"
)

func New(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	err = yaml.NewDecoder(f).Decode(&cfg)

	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

type Config struct {
	Server    Server    `yaml:"server"`
	Database  DBConfig  `yaml:"database"`
	LogConfig LogConfig `yaml:"log_config"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type Server struct {
	HTTP HTTP `yaml:"http"`
}

// HTTP defines server config for http server
type HTTP struct {
	Address        string `yaml:"address"`
	WriteTimeout   string `yaml:"write_timeout"`
	ReadTimeout    string `yaml:"read_timeout"`
	MaxHeaderBytes int    `yaml:"max_header_bytes"`
}

type LogConfig struct {
	LogFilename string `yaml:"log_filename"`
}
