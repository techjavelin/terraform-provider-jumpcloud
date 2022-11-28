package jumpcloud

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type UserGroupResourceModel struct {
	Id                      types.String       `tfsdk:"id"`
	Name                    types.String       `tfsdk:"name"`
	SudoEnabled             types.Bool         `tfsdk:"sudo-enabled"`
	SudoPasswordless        types.Bool         `tfsdk:"sudo-passwordless"`
	LdapGroups              []types.String     `tfsdk:"ldap-groups"`
	PosixGroups             []PosixGroupModel  `tfsdk:"posix-groups"`
	RadiusReplies           []KVItemModel      `tfsdk:"radius-replies"`
	SambaEnabled            types.Bool         `tfsdk:"samba-enabled"`
	Properties              []KVItemModel      `tfsdk:"attribute-properties"`
	Description             types.String       `tfsdk:"description"`
	Email                   types.String       `tfsdk:"email"`
	MemberQuery             []MemberQueryModel `tfsdk:"member-query"`
	MemberSuggestionsNotify types.Bool         `tfsdk:"member-suggestions-notify"`
	MembershipAutomated     types.Bool         `tfsdk:"membership-automated"`
}

type MemberQueryModel struct {
	Field    types.String `tfsdk:"field"`
	Operator types.String `tfsdk:"operator"`
	Value    types.String `tfsdk:"value"`
}

type PosixGroupModel struct {
	Id   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}
