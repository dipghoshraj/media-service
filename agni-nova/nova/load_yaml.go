package nova

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Seeder struct {
	Address      string `yaml:"address"`
	Fingureprint string `yaml:"fingureprint"`
}

type Nova struct {
	Seeder Seeder `yaml:"Seeder"`
	Name   string `yaml:"name"`
	Port   string `yaml:"port"`
}

type Config struct {
	Version string `yaml:"version"`
	Nova    Nova   `yaml:"Nova"`
}

var YamlConfig *Config

func init() {
	data, err := os.ReadFile("agni-config.yaml")
	if err != nil {
		Logger.Error("failed to read config file", "file", "agni-config.yaml", "error", err)
		os.Exit(1)
	}

	var config Config
	if err = yaml.Unmarshal(data, &config); err != nil {
		Logger.Error("failed to parse config file", "file", "agni-config.yaml", "error", err)
		os.Exit(1)
	}

	YamlConfig = &config
	Logger.Info("config loaded", "file", "agni-config.yaml", "version", config.Version)
}
