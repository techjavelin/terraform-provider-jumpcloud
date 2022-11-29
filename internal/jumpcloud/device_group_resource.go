package jumpcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	// "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/davecgh/go-spew/spew"

	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &DeviceGroupResource{}
var _ resource.ResourceWithImportState = &DeviceGroupResource{}

func NewDeviceGroupResource() resource.Resource {
	return &DeviceGroupResource{}
}

type DeviceGroupResource struct {
	api *JumpCloudClientApi
}

type DeviceGroupResourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (r *DeviceGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devicegroup"
}

func (r *DeviceGroupResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Device Group",

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
				MarkdownDescription: "Name for the Device Group",
				Type:                types.StringType,
				Required:            true,
			},
		},
	}, nil
}

func (r *DeviceGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	api, ok := req.ProviderData.(*JumpCloudClientApi)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *JumpCloudClientApi, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.api = api
}

func (r *DeviceGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *DeviceGroupResourceModel

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

	group, response, error := r.api.client.SystemGroupsApi.GroupsSystemPost(r.api.auth, API_ACCEPT_TYPE, API_CONTENT_TYPE, options)

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

func (r *DeviceGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DeviceGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Refreshing Device Group State from JumpCloud")

	group, response, error := r.api.client.SystemGroupsApi.GroupsSystemGet(r.api.auth, state.Id.ValueString(), API_CONTENT_TYPE, API_ACCEPT_TYPE, nil)
	if error != nil {
		resp.Diagnostics.AddError(
			"Error retreiving Active Directory from JumpCloud",
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

func (r *DeviceGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *DeviceGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DeviceGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, error := r.api.client.SystemGroupsApi.GroupsSystemDelete(r.api.auth, state.Id.ValueString(), API_CONTENT_TYPE, API_ACCEPT_TYPE, nil)
	if error != nil {
		resp.Diagnostics.AddError(
			"Error deleting Device from JumpCloud",
			fmt.Sprintf("API Error: %s", spew.Sdump(error)),
		)

		return
	}
	tflog.Trace(ctx, "JumpCloud API Response: \n"+spew.Sdump(response))

}

func (r *DeviceGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
}
