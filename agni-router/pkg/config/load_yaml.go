package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Registry struct {
	Address      string `yaml:"address"`
	Fingureprint string `yaml:"fingureprint"`
}

type Router struct {
	Name       string   `yaml:"name"`
	Registry   Registry `yaml:"registry"`
	Region     string   `yaml:"region"`
	Certs      string   `yaml:"certs"`
	RouterIP   string   `yaml:"router_ip"`
	RouterPort string   `yaml:"rpc_port"`
	Dns        string   `yaml:"dns"`
	ProxtPort  string   `yaml:"proxy_port"`
}

type Config struct {
	Version string `yaml:"version"`
	Router  Router `yaml:"Router"`
}

var YamlConfig *Config

func init() {

	data, err := os.ReadFile("agent-config.yaml")
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
