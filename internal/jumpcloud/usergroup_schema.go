package jumpcloud

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/techjavelin/terraform-provider-jumpcloud/internal/pkg/planmodifiers"
)

var UserGroupSchema = tfsdk.Schema{
	MarkdownDescription: "JumpCloud User Group",
	Description:         "JumpCloud User Group",
	Version:             0,

	Attributes: map[string]tfsdk.Attribute{
		"id": {
			Computed:            true,
			MarkdownDescription: "Resource ID (Computed / Read-Only)",
			PlanModifiers: tfsdk.AttributePlanModifiers{
				resource.UseStateForUnknown(),
			},
			Type: types.StringType,
		},
		"name": {
			MarkdownDescription: "Name for the User Group",
			Type:                types.StringType,
			Required:            true,
		},
		"sudo": {
			MarkdownDescription: "Sudo configuration for the user-group",
			Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
				"enabled": {
					MarkdownDescription: "Whether this user-group will allowed to use sudo",
					Type:                types.BoolType,
					Required:            true,
				},
				"passwordless": {
					MarkdownDescription: "Whether members of this user-group will be able to use sudo without entering a password",
					Type:                types.BoolType,
					Required:            true,
				},
			}),
			Optional: true,
		},
		"ldap": {
			MarkdownDescription: "List of LDAP Groups the user-group is mapped to",
			Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
				"groups": {
					Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
						"name": {
							MarkdownDescription: "The LDAP Group Name",
							Type:                types.StringType,
							Required:            true,
						},
					}),
					Optional: true,
					Computed: true,
				},
			}),
			Optional: true,
			Computed: true,
		},
		"posix": {
			MarkdownDescription: "List of POSIX Groups the user-group is mapped to",
			Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
				"id": {
					MarkdownDescription: "The posix group id",
					Type:                types.Int64Type,
					Required:            true,
				},
				"name": {
					MarkdownDescription: "The posix group name",
					Type:                types.StringType,
					Required:            true,
				},
			}),
			Optional: true,
		},
		"radius": {
			MarkdownDescription: "List of RADIUS Replies to associate with the user-group",
			Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
				"name": {
					MarkdownDescription: "The reply name",
					Type:                types.StringType,
					Required:            true,
				},
				"value": {
					MarkdownDescription: "The reply value",
					Type:                types.StringType,
					Required:            true,
				},
			}),
			Optional: true,
		},
		"samba": {
			MarkdownDescription: "Whether samba propogation is enabled for this user-group",
			Type:                types.BoolType,
			Optional:            true,
		},
		"properties": {
			MarkdownDescription: "List of attribute properties to set on the user-group",
			Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
				"name": {
					MarkdownDescription: "The property name",
					Type:                types.StringType,
					Required:            true,
				},
				"value": {
					MarkdownDescription: "The property value",
					Type:                types.StringType,
					Required:            true,
				},
			}),
			Optional: true,
		},
		"description": {
			MarkdownDescription: "Description for the User Group",
			Type:                types.StringType,
			Optional:            true,
			Computed:            true,
			PlanModifiers: []tfsdk.AttributePlanModifier{
				planmodifiers.StringDefaultModifier{
					Default: "",
				},
			},
		},
		"email": {
			MarkdownDescription: "E-Mail Address for the User Group (Mailing List Group)",
			Type:                types.StringType,
			Optional:            true,
			Computed:            true,
			PlanModifiers: []tfsdk.AttributePlanModifier{
				planmodifiers.StringDefaultModifier{
					Default: "",
				},
			},
		},
		"member_queries": {
			MarkdownDescription: "Query using a sequence of field filters.",
			Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
				"query": {
					Required: true,
					Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
						"field": {
							MarkdownDescription: "The name of the field to query",
							Type:                types.StringType,
							Required:            true,
						},
						"operator": {
							MarkdownDescription: "The operator to use for the query",
							Type:                types.StringType,
							Required:            true,
							Validators: []tfsdk.AttributeValidator{
								stringvalidator.OneOf([]string{"eq", "ne", "gt", "lt", "ge", "le", "between", "search", "in"}...),
							},
						},
						"value": {
							MarkdownDescription: "The value for the filter expression",
							Type:                types.StringType,
							Required:            true,
						},
					}),
				},
			}),
			Optional: true,
		},
		"notify": {
			MarkdownDescription: "Whether to send notifications for new member suggestions that match member-query-filters",
			Type:                types.BoolType,
			Optional:            true,
			Computed:            true,
		},
		"auto": {
			MarkdownDescription: "Whether users matching member-query-filters should be automatically added to the user-group",
			Type:                types.BoolType,
			Optional:            true,
			Computed:            true,
		},
	},
}
