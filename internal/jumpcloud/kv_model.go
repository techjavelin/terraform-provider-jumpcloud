package jumpcloud

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type KVItemModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}
