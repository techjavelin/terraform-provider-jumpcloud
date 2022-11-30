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
	var updatePlan *UserGroupResourceModel

	diags := req.Plan.Get(ctx, &updatePlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Got Error while trying to set plan")
		return
	}

	tflog.Trace(ctx, "Evaluated Plan", map[string]interface{}{
		"UserGroupResourceModel": spew.Sdump(updatePlan),
	})

	tflog.Info(ctx, "Refreshing User Group State from JumpCloud")
	apiModel := convertResourceToUserGroup(ctx, updatePlan)

	updatedApiModel, response, error := r.api.UpdateUserGroup(&apiModel)

	tflog.Trace(ctx, "Got response from JumpCloud API", map[string]interface{}{
		"Response":  r.api.ReadBody(response.Body),
		"UserGroup": spew.Sdump(updatedApiModel),
		"Error":     spew.Sdump(error),
	})

	if error != nil {
		resp.Diagnostics.AddError(
			"Error updating User Group on JumpCloud",
			fmt.Sprintf("API Error: %s", spew.Sdump(error)),
		)

		return
	}

	r.convertApiResponseToResource(ctx, updatePlan, &updatedApiModel)

	diags = resp.State.Set(ctx, updatePlan)
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

func (r *UserGroupResource) convertApiResponseToResource(ctx context.Context, resourceModel *UserGroupResourceModel, apiModel *apiclient.UserGroup) (diags diag.Diagnostics) {
	resourceModel.Id = types.StringValue(apiModel.Id)
	resourceModel.Name = types.StringValue(apiModel.Name)
	resourceModel.Description = types.StringValue(apiModel.Description)
	resourceModel.Email = types.StringValue(apiModel.Email)
	resourceModel.MemberSuggestionsNotify = types.BoolValue(apiModel.MemberSuggestionsNotify)
	resourceModel.MembershipAutomated = types.BoolValue(apiModel.MembershipAutomated)

	if apiModel.Attributes != nil {
		if apiModel.Attributes.SambaEnabled {
			resourceModel.Samba.Enabled = types.BoolValue(apiModel.Attributes.SambaEnabled)
		}

		if apiModel.Attributes.Sudo != nil {
			resourceModel.Sudo.Enabled = types.BoolValue(apiModel.Attributes.Sudo.Enabled)
			resourceModel.Sudo.Passwordless = types.BoolValue(apiModel.Attributes.Sudo.WithoutPassword)
		}

		for _, ldapGroup := range apiModel.Attributes.LdapGroups {
			var ldapInfo LdapInfo
			if resourceModel.Ldap.IsNull() || resourceModel.Ldap.IsUnknown() {
				ldapInfo = LdapInfo{}
			} else {
				diags := resourceModel.Ldap.As(ctx, ldapInfo, types.ObjectAsOptions{
					UnhandledNullAsEmpty:    true,
					UnhandledUnknownAsEmpty: true,
				})

				if diags.HasError() {
					tflog.Error(ctx, "Encountered error attempting to retrieve LdapInfo as Object from Type", map[string]interface{}{
						"IsNull":    resourceModel.Ldap.IsNull(),
						"IsUnknown": resourceModel.Ldap.IsUnknown(),
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

			resourceModel.Ldap = ldap
		}

		for _, posixGroup := range apiModel.Attributes.PosixGroups {
			resourceModel.PosixGroups = append(resourceModel.PosixGroups, PosixGroupModel{
				Id:   types.Int64Value(posixGroup.Id),
				Name: types.StringValue(posixGroup.Name),
			})
		}

		if apiModel.Attributes.Radius != nil {
			for _, radiusReply := range apiModel.Attributes.Radius.Reply {
				resourceModel.RadiusReplies = append(resourceModel.RadiusReplies, KVItemModel{
					Name:  types.StringValue(radiusReply.Name),
					Value: types.StringValue(radiusReply.Value),
				})
			}
		}
	}

	if apiModel.MemberQuery != nil && len(apiModel.MemberQuery.Filters) > 0 {
		for _, filter := range apiModel.MemberQuery.Filters {
			resourceModel.MemberQuery = append(resourceModel.MemberQuery, MemberQueryModel{
				Field:    types.StringValue(filter.Field),
				Operator: types.StringValue(filter.Operator),
				Value:    types.StringValue(filter.Value),
			})
		}
	}

	tflog.Trace(ctx, fmt.Sprintf("Converted UserGroup to UserGroupResourceModel:\n\tUserGroup: %s\n\tUserGroupResourceModel: %s", spew.Sdump(apiModel), spew.Sdump(resourceModel)))

	return diags
}

func convertResourceToUserGroup(ctx context.Context, resourceModel *UserGroupResourceModel) apiclient.UserGroup {
	var sudoConfig *apiclient.UserGroupSudoConfig
	if resourceModel.Sudo != nil {
		sudoConfig = &apiclient.UserGroupSudoConfig{
			Enabled:         resourceModel.Sudo.Enabled.ValueBool(),
			WithoutPassword: resourceModel.Sudo.Passwordless.ValueBool(),
		}
	}

	var ldapGroups []apiclient.LdapGroup

	if !resourceModel.Ldap.IsNull() {
		var ldapInfo LdapInfo
		resourceModel.Ldap.As(ctx, ldapInfo, types.ObjectAsOptions{
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
	for _, posix_group := range resourceModel.PosixGroups {
		posixGroups = append(posixGroups, apiclient.PosixGroup{
			Id:   posix_group.Id.ValueInt64(),
			Name: posix_group.Name.ValueString(),
		})
	}

	var radiusConfig *apiclient.UserGroupRadiusConfig
	var radiusReplies []apiclient.RadiusReply
	for _, reply := range resourceModel.RadiusReplies {
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
	if resourceModel.Samba != nil {
		sambaEnabled = resourceModel.Samba.Enabled.Equal(types.BoolValue(true))
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
	if len(resourceModel.MemberQuery) > 0 {
		for _, query := range resourceModel.MemberQuery {
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

	apiModel := apiclient.UserGroup{
		Id:          resourceModel.Id.ValueString(),
		Name:        resourceModel.Name.ValueString(),
		MemberQuery: memberQuery,
		Attributes:  attributes,
	}

	apiModel.Description = resourceModel.Description.ValueString()
	apiModel.Email = resourceModel.Email.ValueString()

	apiModel.MemberSuggestionsNotify = resourceModel.MemberSuggestionsNotify.ValueBool()
	apiModel.MembershipAutomated = resourceModel.MembershipAutomated.ValueBool()

	tflog.Info(ctx, fmt.Sprintf("Converted UserGroupResourceModel to UserGroup:\n\tUserGroup: %s\n\tUserGroupResourceModel: %s", spew.Sdump(apiModel), spew.Sdump(resourceModel)))
	return apiModel
}
