package jumpcloud

import (
	"strings"
	"testing"

	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func ProviderConfig() string {
	return "provider \"jumpcloud\" { api_key = \"" + os.Getenv("JUMPCLOUD_API_KEY") + "\" }\n"
}

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"jumpcloud": providerserver.NewProtocol6WithError(New("dev")()),
}

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
