package activedirectory

import (
    "testing"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccActiveDirectoryResource(t *testing.T) {
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
					resource.TestCheckResourceAttr("jumpcloud_ad.test", "domain", "DC=test,DC=com"),
					resource.TestCheckResourceAttrSet("jumpcloud_ad.test", "id"),
				),
			},
		},
	})
}