package jumpcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/techjavelin/terraform-provider-jumpcloud/internal/pkg/jumpcloud/apiclient"

	"github.com/davecgh/go-spew/spew"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &UserGroupResource{}
var _ resource.ResourceWithImportState = &UserGroupResource{}
var _ resource.ResourceWithImportState = &UserGroupResource{}

func NewUserGroupResource() resource.Resource {
	return &UserGroupResource{}
}

type UserGroupResource struct {
	api *apiclient.Client
}

func (r *UserGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_usergroup"
}

func (r *UserGroupResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return UserGroupSchema, nil
}

func (r *UserGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	api, ok := req.ProviderData.(JumpCloudApi)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *api.JumpCloudClientApi, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.api = &api.Internal
}

func (r *UserGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *UserGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	tflog.Trace(ctx, fmt.Sprintf("Called UserGroupResource.Create with:\n%s", spew.Sdump(plan)))

	if resp.Diagnostics.HasError() {
		return
	}

	usergroup := convertResourceToUserGroup(ctx, plan)

	tflog.Debug(ctx, fmt.Sprintf("Calling GroupsSystemPost with\n%s", spew.Sdump(usergroup)))

	group, response, error := r.api.CreateUserGroup(&usergroup)
	tflog.Info(ctx, fmt.Sprintf("JumpCloud API Response: %s\n", r.api.ReadBody(response.Body)))

	if error != nil {
		resp.Diagnostics.AddError(
			"Error creating User Group",
			fmt.Sprintf("API Error: %s", spew.Sdump(error)),
		)
		return
	}

	var created *UserGroupResourceModel = &UserGroupResourceModel{}

	r.convertApiResponseToResource(ctx, created, &group)
	tflog.Info(ctx, "Created new User Group", map[string]interface{}{
		"api_response":   "\n\n" + spew.Sdump(group) + "\n\n",
		"resource_model": "\n\n" + spew.Sdump(created) + "\n\n",
	})

	var diags = resp.State.Set(ctx, created)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *UserGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Refreshing User Group State from JumpCloud")

	var plan *UserGroupResourceModel

	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, _, error := r.api.GetUserGroupDetails(plan.Id.ValueString())

	if error != nil {
		resp.Diagnostics.AddError(
			"Error retreiving User Group from JumpCloud",
			fmt.Sprintf("API Error: %s", spew.Sdump(error)),
		)

		return
	}

	// tflog.Trace(ctx, "JumpCloud API Response: \n"+spew.Sdump(response.Body))

	r.convertApiResponseToResource(ctx, plan, &group)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *UserGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating UserGroupResource")
	var plan *UserGroupResourceModel

	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Got Error while trying to set plan")
		return
	}

	tflog.Trace(ctx, "Evaluated Plan", map[string]interface{}{
		"UserGroupResourceModel": spew.Sdump(plan),
	})

	tflog.Info(ctx, "Refreshing User Group State from JumpCloud")
	usergroup := convertResourceToUserGroup(ctx, plan)

	group, response, error := r.api.UpdateUserGroup(&usergroup)

	tflog.Trace(ctx, "Got response from JumpCloud API", map[string]interface{}{
		"Response":  r.api.ReadBody(response.Body),
		"UserGroup": spew.Sdump(group),
		"Error":     spew.Sdump(error),
	})

	if error != nil {
		resp.Diagnostics.AddError(
			"Error updating User Group on JumpCloud",
			fmt.Sprintf("API Error: %s", spew.Sdump(error)),
		)

		return
	}

	r.convertApiResponseToResource(ctx, plan, &group)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *UserGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *UserGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, error := r.api.DeleteUserGroup(state.Id.ValueString())
	tflog.Info(ctx, fmt.Sprintf("JumpCloud API Response: %s\n", r.api.ReadBody(response.Body)))

	if error != nil {
		resp.Diagnostics.AddError(
			"Error deleting User Group from JumpCloud",
			fmt.Sprintf("API Error: %s", spew.Sdump(error)),
		)

		return
	}

	tflog.Trace(ctx, "JumpCloud API Response: \n"+spew.Sdump(response.Body))
}

func (r *UserGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getStringValueIfNotNil(val *string) types.String {
	if val != nil {
		return types.StringValue(*val)
	}

	return types.StringNull()
}

func (r *UserGroupResource) convertApiResponseToResource(ctx context.Context, plan *UserGroupResourceModel, group *apiclient.UserGroup) (diags diag.Diagnostics) {
	plan.Id = types.StringValue(group.Id)
	plan.Name = types.StringValue(group.Name)
	plan.Description = types.StringValue(group.Description)
	plan.Email = types.StringValue(group.Email)
	plan.MemberSuggestionsNotify = types.BoolValue(group.MemberSuggestionsNotify)
	plan.MembershipAutomated = types.BoolValue(group.MembershipAutomated)

	if group.Attributes != nil {
		if group.Attributes.SambaEnabled {
			plan.Samba.Enabled = types.BoolValue(group.Attributes.SambaEnabled)
		}

		if group.Attributes.Sudo != nil {
			plan.Sudo.Enabled = types.BoolValue(group.Attributes.Sudo.Enabled)
			plan.Sudo.Passwordless = types.BoolValue(group.Attributes.Sudo.WithoutPassword)
		}

		for _, ldapGroup := range group.Attributes.LdapGroups {
			var ldapInfo LdapInfo
			if plan.Ldap.IsNull() || plan.Ldap.IsUnknown() {
				ldapInfo = LdapInfo{}
			} else {
				diags := plan.Ldap.As(ctx, ldapInfo, types.ObjectAsOptions{
					UnhandledNullAsEmpty:    true,
					UnhandledUnknownAsEmpty: true,
				})

				if diags.HasError() {
					tflog.Error(ctx, "Encountered error attempting to retrieve LdapInfo as Object from Type", map[string]interface{}{
						"IsNull":    plan.Ldap.IsNull(),
						"IsUnknown": plan.Ldap.IsUnknown(),
					})
					return diags
				}
			}

			ldapInfo.LdapGroups = append(ldapInfo.LdapGroups, LdapGroupModel{
				Name: types.StringValue(ldapGroup.Name),
			})

			tflog.Trace(ctx, "Setting LdapInfo as types.Object on plan", map[string]interface{}{
				"LdapInfo": spew.Sdump(ldapInfo),
			})

			ldap, d := types.ObjectValueFrom(ctx, ldapInfo.AttrTypes(), ldapInfo)

			diags.Append(d...)
			if d.HasError() {
				return diags
			}

			plan.Ldap = ldap
		}

		for _, posixGroup := range group.Attributes.PosixGroups {
			plan.PosixGroups = append(plan.PosixGroups, PosixGroupModel{
				Id:   types.Int64Value(posixGroup.Id),
				Name: types.StringValue(posixGroup.Name),
			})
		}

		if group.Attributes.Radius != nil {
			for _, radiusReply := range group.Attributes.Radius.Reply {
				plan.RadiusReplies = append(plan.RadiusReplies, KVItemModel{
					Name:  types.StringValue(radiusReply.Name),
					Value: types.StringValue(radiusReply.Value),
				})
			}
		}
	}

	if group.MemberQuery != nil && len(group.MemberQuery.Filters) > 0 {
		for _, filter := range group.MemberQuery.Filters {
			plan.MemberQuery = append(plan.MemberQuery, MemberQueryModel{
				Field:    types.StringValue(filter.Field),
				Operator: types.StringValue(filter.Operator),
				Value:    types.StringValue(filter.Value),
			})
		}
	}

	tflog.Trace(ctx, fmt.Sprintf("Converted UserGroup to UserGroupResourceModel:\n\tUserGroup: %s\n\tUserGroupResourceModel: %s", spew.Sdump(group), spew.Sdump(plan)))

	return diags
}

func getStringIfAttributeNotNil(in types.String) *string {
	if in.IsNull() {
		return nil
	}

	out := in.ValueString()
	return &out
}

func convertResourceToUserGroup(ctx context.Context, plan *UserGroupResourceModel) apiclient.UserGroup {
	var sudoConfig *apiclient.UserGroupSudoConfig
	if plan.Sudo != nil {
		sudoConfig = &apiclient.UserGroupSudoConfig{
			Enabled:         plan.Sudo.Enabled.ValueBool(),
			WithoutPassword: plan.Sudo.Passwordless.ValueBool(),
		}
	}

	var ldapGroups []apiclient.LdapGroup

	if !plan.Ldap.IsNull() {
		var ldapInfo LdapInfo
		plan.Ldap.As(ctx, ldapInfo, types.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		})

		for _, ldap_group := range ldapInfo.LdapGroups {
			ldapGroups = append(ldapGroups, apiclient.LdapGroup{
				Name: ldap_group.Name.ValueString(),
			})
		}
	}

	var posixGroups []apiclient.PosixGroup
	for _, posix_group := range plan.PosixGroups {
		posixGroups = append(posixGroups, apiclient.PosixGroup{
			Id:   posix_group.Id.ValueInt64(),
			Name: posix_group.Name.ValueString(),
		})
	}

	var radiusConfig *apiclient.UserGroupRadiusConfig
	var radiusReplies []apiclient.RadiusReply
	for _, reply := range plan.RadiusReplies {
		radiusReplies = append(radiusReplies, apiclient.RadiusReply{
			Name:  reply.Name.ValueString(),
			Value: reply.Value.ValueString(),
		})
	}

	if len(radiusReplies) > 0 {
		radiusConfig = &apiclient.UserGroupRadiusConfig{
			Reply: radiusReplies,
		}
	}

	var sambaEnabled bool
	if plan.Samba != nil {
		sambaEnabled = plan.Samba.Enabled.Equal(types.BoolValue(true))
	}

	var attributes *apiclient.UserGroupAttributes
	if sudoConfig != nil || len(ldapGroups) > 0 || len(posixGroups) > 0 || radiusConfig != nil || sambaEnabled {
		attributes = &apiclient.UserGroupAttributes{
			Sudo:         sudoConfig,
			LdapGroups:   ldapGroups,
			PosixGroups:  posixGroups,
			Radius:       radiusConfig,
			SambaEnabled: sambaEnabled,
		}
	}

	var memberQuery *apiclient.UserGroupMemberQuery
	var filterQueries []apiclient.QueryFilter
	if len(plan.MemberQuery) > 0 {
		for _, query := range plan.MemberQuery {
			filterQueries = append(filterQueries, apiclient.QueryFilter{
				Field:    query.Field.ValueString(),
				Operator: query.Operator.ValueString(),
				Value:    query.Value.ValueString(),
			})
		}

		memberQuery = &apiclient.UserGroupMemberQuery{
			QueryType: "FilterQuery",
			Filters:   filterQueries,
		}
	}

	usergroup := apiclient.UserGroup{
		Id:          plan.Id.ValueString(),
		Name:        plan.Name.ValueString(),
		MemberQuery: memberQuery,
		Attributes:  attributes,
	}

	usergroup.Description = plan.Description.ValueString()
	usergroup.Email = plan.Description.ValueString()

	usergroup.MemberSuggestionsNotify = plan.MemberSuggestionsNotify.ValueBool()
	usergroup.MembershipAutomated = plan.MembershipAutomated.ValueBool()

	tflog.Info(ctx, fmt.Sprintf("Converted UserGroupResourceModel to UserGroup:\n\tUserGroup: %s\n\tUserGroupResourceModel: %s", spew.Sdump(usergroup), spew.Sdump(plan)))
	return usergroup
}
