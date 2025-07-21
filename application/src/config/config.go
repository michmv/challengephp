package config

import (
	"challengephp/lib"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Debug   bool   `yaml:"debug"`
	LogFile string `yaml:"logFile"`
	DB      DB     `yaml:"db"`
}

type DB struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
}

func LoadConfig(path string) (Config, lib.Error) {
	data, err := os.ReadFile(path)
	conf := Config{}
	if err != nil {
		return conf, lib.Err(err)
	}
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return conf, lib.Err(err)
	}

	return setDefaultValues(conf), nil
}

func setDefaultValues(conf Config) Config {
	if conf.LogFile == "" {
		conf.LogFile = "app.log"
	}
	return conf
}
