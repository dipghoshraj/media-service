package bridge

import (
	"github.com/odio4u/mem-sdk/certengine/pkg"
)

func BuildCreds(dns, name string) error {
	err := pkg.GenerateSelfSignedAgent(
		name,
		[]string{dns},
	)
	pkg.Must(err)
	return nil
}
