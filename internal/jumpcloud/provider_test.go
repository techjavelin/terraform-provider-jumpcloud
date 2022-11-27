package jumpcloud

import (
    "github.com/hashicorp/terraform-plugin-framework/providerserver"
    "github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	providerConfig = `
provider "jumpcloud" {
	api_key = "test"
}
`
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error) {
		"jumpcloud": providerserver.NewProtocol6WithError(New("dev")()),
	}
)