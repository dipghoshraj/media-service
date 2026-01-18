package nova

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Seeder struct {
	Address      string `yaml:"address"`
	Fingureprint string `yaml:"fingureprint"`
}

type Nova struct {
	Seeder Seeder `yaml:"Seeder"`
}

type Config struct {
	Version string `yaml:"version"`
	Nova    Nova   `yaml:"Nova"`
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
