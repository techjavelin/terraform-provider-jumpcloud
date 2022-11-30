package jumpcloud

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type UserGroupResourceModel struct {
	Id                      types.String       `tfsdk:"id"`
	Name                    types.String       `tfsdk:"name"`
	Sudo                    *SudoConfigModel   `tfsdk:"sudo"`
	Ldap                    types.Object       `tfsdk:"ldap"`
	PosixGroups             []PosixGroupModel  `tfsdk:"posix"`
	RadiusReplies           []KVItemModel      `tfsdk:"radius"`
	Samba                   *SambaConfig       `tfsdk:"samba"`
	Properties              []KVItemModel      `tfsdk:"properties"`
	Description             types.String       `tfsdk:"description"`
	Email                   types.String       `tfsdk:"email"`
	MemberQuery             []MemberQueryModel `tfsdk:"member_queries"`
	MemberSuggestionsNotify types.Bool         `tfsdk:"notify"`
	MembershipAutomated     types.Bool         `tfsdk:"auto"`
}

type SambaConfig struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

type LdapInfo struct {
	LdapGroups []LdapGroupModel `tfsdk:"groups"`
}

func (l LdapInfo) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"groups": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"name": types.StringType,
				},
			},
		},
	}
}

type LdapGroupModel struct {
	Name types.String `tfsdk:"name"`
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

type SudoConfigModel struct {
	Enabled      types.Bool `tfsdk:"enabled"`
	Passwordless types.Bool `tfsdk:"passwordless"`
}
