package jumpcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const ResourceConfig = `
resource "jumpcloud_ad" "test" {
	domain = "DC=test,DC=com"
}
`

func TestAccActiveDirectoryResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read Testing
			{
				Config: ProviderConfig() + ResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jumpcloud_ad.test", "domain", "DC=test,DC=com"),
					resource.TestCheckResourceAttrSet("jumpcloud_ad.test", "id"),
				),
			},
			// ImportState Testing
			{
				ResourceName:      "jumpcloud_ad.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update is not supported -- how do we test for expecting an error tho?
			{
				Config: ProviderConfig() + `resource "jumpcloud_ad" "test" { domain = "DC=update,DC=test,DC=com" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jumpcloud_ad.test", "domain", "DC=update,DC=test,DC=com"),
				),
			},
			// Delete Testing happens automatically
		},
	})
}
