package model

import "github.com/hashicorp/terraform-plugin-framework/attr"

type ModelWithAttributeTypes interface {
	AttrTypes() map[string]attr.Type
}
