package bridge

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Seeder struct {
	Address      string `yaml:"address"`
	Fingureprint string `yaml:"fingureprint"`
}

type Agent struct {
	Name    string `yaml:"name"`
	Domain  string `yaml:"domain"`
	Forward int    `yaml:"forward"`
	Seeder  Seeder `yaml:"Seeder"`
	Region  string `yaml:"region"`
	Certs   string `yaml:"certs"`
}

type Config struct {
	Version string `yaml:"version"`
	Agent   Agent  `yaml:"Agent"`
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
