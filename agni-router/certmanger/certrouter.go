package main

import (
	"log"

	"github.com/Purple-House/agni-tunnels/agni-router/pkg/config"
	"github.com/Purple-House/mem-sdk/certengine/pkg"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// config.YamlConfig.Router.RouterPort

	routerIps := []string{config.YamlConfig.Router.RouterIP}
	dns := []string{config.YamlConfig.Router.Dns}

	_, err := pkg.GenerateSelfSignedGPR(config.YamlConfig.Router.Name, routerIps, dns)
	pkg.Must(err)

	log.Println("All operations completed successfully.")
}
