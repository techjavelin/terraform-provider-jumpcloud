package jumpcloud

import (
	"context"
	"fmt"

	"github.com/techjavelin/terraform-provider-jumpcloud/internal/pkg/jumpcloud/api"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/davecgh/go-spew/spew"

	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &ActiveDirectoryResource{}
	_ resource.ResourceWithConfigure   = &ActiveDirectoryResource{}
	_ resource.ResourceWithImportState = &ActiveDirectoryResource{}
)

func NewActiveDirectoryResource() resource.Resource {
	return &ActiveDirectoryResource{}
}

type ActiveDirectoryResource struct {
	api *api.JumpCloudClientApi
}

type ActiveDirectoryResourceModel struct {
	Domain types.String `tfsdk:"domain"`
	Id     types.String `tfsdk:"id"`
}

func (r *ActiveDirectoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ad"
}

func (r *ActiveDirectoryResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Active Directory",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed:            true,
				MarkdownDescription: "Resource ID (Computed / Read-Only)",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
			"domain": {
				MarkdownDescription: "The Active Directory Domain (eg DC=mydomain,DC=com}",
				Required:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (r *ActiveDirectoryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	api, ok := req.ProviderData.(*api.JumpCloudClientApi)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *JumpCloudClientApi, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.api = api
}

func (r *ActiveDirectoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *ActiveDirectoryResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var domain = plan.Domain.ValueString()

	var options = make(map[string]interface{})
	options["body"] = jcapiv2.ActiveDirectoryInput{
		Domain: domain,
	}

	tflog.Info(ctx, fmt.Sprintf("Calling ActiveDirectoriesPost with\n%s", spew.Sdump(options)))

	ad, response, error := r.api.Client.ActiveDirectoryApi.ActivedirectoriesPost(r.api.Auth, api.API_ACCEPT_TYPE, api.API_CONTENT_TYPE, options)

	if error != nil {
		resp.Diagnostics.AddError(
			"Error creating Active Directory",
			fmt.Sprintf("API Error: %s", spew.Sdump(error)),
		)
		return
	}

	tflog.Trace(ctx, "JumpCloud API Response: \n"+spew.Sdump(response))
	tflog.Info(ctx, fmt.Sprintf("Created new Active Directory\n%s", spew.Sdump(ad)))

	plan.Id = types.StringValue(ad.Id)
	plan.Domain = types.StringValue(ad.Domain)

	var diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ActiveDirectoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ActiveDirectoryResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Refreshing Active Directory State from JumpCloud")

	ad, response, error := r.api.Client.ActiveDirectoryApi.ActivedirectoriesGet(r.api.Auth, state.Id.ValueString(), api.API_CONTENT_TYPE, api.API_ACCEPT_TYPE, nil)
	if error != nil {
		resp.Diagnostics.AddError(
			"Error retreiving Active Directory from JumpCloud",
			fmt.Sprintf("API Error: %s", spew.Sdump(error)),
		)

		return
	}

	tflog.Trace(ctx, "JumpCloud API Response: \n"+spew.Sdump(response))

	state.Domain = types.StringValue(ad.Domain)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ActiveDirectoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update is unsupported for Active Directories",
		"Update is unsupported for Active Directories",
	)
}

func (r *ActiveDirectoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ActiveDirectoryResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, error := r.api.Client.ActiveDirectoryApi.ActivedirectoriesDelete(r.api.Auth, state.Id.ValueString(), api.API_CONTENT_TYPE, api.API_ACCEPT_TYPE, nil)
	if error != nil {
		resp.Diagnostics.AddError(
			"Error deleting Active Directory from JumpCloud",
			fmt.Sprintf("API Error: %s", spew.Sdump(error)),
		)

		return
	}
	tflog.Trace(ctx, "JumpCloud API Response: \n"+spew.Sdump(response))

}

func (r *ActiveDirectoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
