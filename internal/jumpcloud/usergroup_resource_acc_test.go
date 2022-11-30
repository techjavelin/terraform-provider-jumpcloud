package jumpcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUserGroupResource(t *testing.T) {
	test_env := GetTestEnv()
	group_name := fmt.Sprintf("terraform-test-usergroup-%s", test_env)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig() + `
resource "jumpcloud_usergroup" "test" {
	name = "` + group_name + `"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jumpcloud_usergroup.test", "name", group_name),
					resource.TestCheckResourceAttrSet("jumpcloud_usergroup.test", "id"),
				),
			},
			{
				ResourceName:      "jumpcloud_usergroup.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: ProviderConfig() + `
resource "jumpcloud_usergroup" "test" {
	name = "` + group_name + `-updated"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jumpcloud_usergroup.test", "name", group_name+"-updated"),
				),
			},
		},
	})
}
