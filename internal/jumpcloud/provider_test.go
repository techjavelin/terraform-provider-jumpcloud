package jumpcloud

import (
	"strings"
	"testing"

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
	out = strings.ReplaceAll(name, ".", "-")
	out = strings.ReplaceAll(out, "*", "x")
	return out
}

func TestMakeFriendlyName(t *testing.T) {
	expect := "1-2-x"
	test := makeFriendlyName("1.2.*")
	if test != expect {
		t.Fatalf("Expected %s but got %s", expect, test)
	}
}
