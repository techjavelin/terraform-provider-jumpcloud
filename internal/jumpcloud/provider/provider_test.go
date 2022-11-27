package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	providerConfig = `
variable "jumpcloud_api_key" {}
provider "jumpcloud" {
	api_key = "${ var.jumpcloud_api_key }"
}
`
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"jumpcloud": providerserver.NewProtocol6WithError(New("dev")()),
	}
)
