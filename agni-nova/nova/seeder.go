package nova

import (
	"fmt"
	"time"

	mp "github.com/odio4u/mem-sdk/memsdk/maps"
	"github.com/odio4u/mem-sdk/memsdk/pkg"
)

func SeederClient() (*mp.Client, error) {

	config := pkg.Config{
		Address:     YamlConfig.Nova.Seeder.Address,
		Fingerprint: YamlConfig.Nova.Seeder.Fingureprint,
		Timeout:     5 * time.Second,
	}

	client, err := mp.NewSdkOperation(config)
	if err != nil {
		return nil, fmt.Errorf("Seeder connection error %v", err)
	}
	return client, nil
}
