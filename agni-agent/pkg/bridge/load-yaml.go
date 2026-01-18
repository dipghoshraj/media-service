package bridge

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Registry struct {
	Address      string `yaml:"address"`
	Fingureprint string `yaml:"fingureprint"`
}

type Agent struct {
	Name     string   `yaml:"name"`
	Domain   string   `yaml:"domain"`
	Forward  int      `yaml:"forward"`
	Registry Registry `yaml:"registry"`
	Region   string   `yaml:"region"`
	Certs    string   `yaml:"certs"`
}

type Config struct {
	Version string `yaml:"version"`
	Agent   Agent  `yaml:"Agent"`
}

var YamlConfig *Config

func init() {

	data, err := os.ReadFile("agni-config.yaml")
	log.Println("load the config")

	if err != nil {
		log.Fatalf("error reading YAML file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error parsing YAML file: %v", err)
	}

	YamlConfig = &config
}
