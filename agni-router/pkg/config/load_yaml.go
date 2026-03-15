package config

import (
	"os"

	"github.com/odio4u/agni-tunnels/agni-router/pkg/logger"
	"gopkg.in/yaml.v3"
)

type Seeder struct {
	Address      string `yaml:"address"`
	Fingureprint string `yaml:"fingureprint"`
}

type Router struct {
	Name       string `yaml:"name"`
	Seeder     Seeder `yaml:"Seeder"`
	Region     string `yaml:"region"`
	Certs      string `yaml:"certs"`
	RouterIP   string `yaml:"router_ip"`
	RouterPort string `yaml:"rpc_port"`
	Dns        string `yaml:"dns"`
	ProxtPort  string `yaml:"proxy_port"`
}

type Config struct {
	Version string `yaml:"version"`
	Router  Router `yaml:"Router"`
}

var YamlConfig *Config

func init() {
	data, err := os.ReadFile("agni-config.yaml")
	if err != nil {
		logger.Logger.Error("failed to read config file", "file", "agni-config.yaml", "error", err)
		os.Exit(1)
	}

	var config Config
	if err = yaml.Unmarshal(data, &config); err != nil {
		logger.Logger.Error("failed to parse config file", "file", "agni-config.yaml", "error", err)
		os.Exit(1)
	}

	YamlConfig = &config
	logger.Logger.Info("config loaded", "file", "agni-config.yaml", "version", config.Version)
}
