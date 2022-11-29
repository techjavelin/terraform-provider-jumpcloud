package jumpcloud

import (
	"regexp"
	"strings"

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

func makeFriendlyName(name string) (out string) {
	m := regexp.MustCompile(`[a-zA-Z0-9_-]`)
	out = strings.ReplaceAll(name, ".", "-")
	out = strings.ReplaceAll(out, "*", "x")
	return m.ReplaceAllString(out, "_")
}
