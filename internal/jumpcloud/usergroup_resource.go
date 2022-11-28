package jumpcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/techjavelin/terraform-provider-jumpcloud/internal/pkg/jumpcloud/api"

	"github.com/davecgh/go-spew/spew"

	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &UserGroupResource{}
var _ resource.ResourceWithImportState = &UserGroupResource{}

func NewUserGroupResource() resource.Resource {
	return &UserGroupResource{}
}

type UserGroupResource struct {
	api *api.JumpCloudClientApi
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

	api, ok := req.ProviderData.(*api.JumpCloudClientApi)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *api.JumpCloudClientApi, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.api = api
}

func (r *UserGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *UserGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var name = plan.Name.ValueString()

	var options = make(map[string]interface{})
	options["body"] = jcapiv2.SystemGroupData{
		Name: name,
	}

	tflog.Info(ctx, fmt.Sprintf("Calling GroupsSystemPost with\n%s", spew.Sdump(options)))

	group, response, error := r.api.Client.SystemGroupsApi.GroupsSystemPost(r.api.Auth, api.API_ACCEPT_TYPE, api.API_CONTENT_TYPE, options)

	if error != nil {
		resp.Diagnostics.AddError(
			"Error creating Device Group",
			fmt.Sprintf("API Error: %s", spew.Sdump(error)),
		)
		return
	}

	tflog.Trace(ctx, "JumpCloud API Response: \n"+spew.Sdump(response))
	tflog.Info(ctx, fmt.Sprintf("Created new Device Group\n%s", spew.Sdump(group)))

	plan.Id = types.StringValue(group.Id)
	plan.Name = types.StringValue(group.Name)

	var diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *UserGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Refreshing User Group State from JumpCloud")

	group, response, error := r.api.Client.SystemGroupsApi.GroupsSystemGet(r.api.Auth, state.Id.ValueString(), api.API_CONTENT_TYPE, api.API_ACCEPT_TYPE, nil)
	if error != nil {
		resp.Diagnostics.AddError(
			"Error retreiving User Group from JumpCloud",
			fmt.Sprintf("API Error: %s", spew.Sdump(error)),
		)

		return
	}

	tflog.Trace(ctx, "JumpCloud API Response: \n"+spew.Sdump(response))

	state.Name = types.StringValue(group.Name)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *UserGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *UserGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, error := r.api.Client.SystemGroupsApi.GroupsSystemDelete(r.api.Auth, state.Id.ValueString(), api.API_CONTENT_TYPE, api.API_ACCEPT_TYPE, nil)
	if error != nil {
		resp.Diagnostics.AddError(
			"Error deleting User Group from JumpCloud",
			fmt.Sprintf("API Error: %s", spew.Sdump(error)),
		)

		return
	}
	tflog.Trace(ctx, "JumpCloud API Response: \n"+spew.Sdump(response))

}

func (r *UserGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
}
