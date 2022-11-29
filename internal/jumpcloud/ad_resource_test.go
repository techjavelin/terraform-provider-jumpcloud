package jumpcloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccActiveDirectoryResource(t *testing.T) {

	test_env := os.Getenv("TF_VAR_test_env")
	if len(test_env) == 0 {
		test_env = "default"
	}

	test_env = makeFriendlyName(test_env)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "jumpcloud_ad" "test" {
	domain = "DC=test,DC=com"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jumpcloud_ad.test", "domain", fmt.Sprintf("DC=%sDC=test,DC=com", test_env)),
					resource.TestCheckResourceAttrSet("jumpcloud_ad.test", "id"),
				),
			},
		},
	})
}
