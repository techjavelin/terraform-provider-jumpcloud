package jumpcloud

import (
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func ProviderConfig() string {
	return `
provider "jumpcloud" {
	api_key = "${ ` + os.Getenv("JUMPCLOUD_API_KEY") + ` }"
}
`
}

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"jumpcloud": providerserver.NewProtocol6WithError(New("dev")()),
}
