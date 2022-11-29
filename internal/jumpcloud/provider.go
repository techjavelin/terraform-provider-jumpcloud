package jumpcloud

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/techjavelin/terraform-provider-jumpcloud/internal/pkg/jumpcloud/api"
	"github.com/techjavelin/terraform-provider-jumpcloud/internal/pkg/jumpcloud/apiclient"

	jcapiv1 "github.com/TheJumpCloud/jcapi-go/v1"
	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
)

var _ provider.Provider = &JumpCloudProvider{}
var _ provider.ProviderWithMetadata = &JumpCloudProvider{}

type JumpCloudProvider struct {
	version string
}

type JumpCloudProviderModel struct {
	APIKey types.String `tfsdk:"api_key"`
}

type JumpCloudApi struct {
	V1       api.JumpCloudClientApiV1
	V2       api.JumpCloudClientApiV2
	Internal apiclient.Client
}

func (p *JumpCloudProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "jumpcloud"
	resp.Version = p.version
}

func (p *JumpCloudProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"api_key": {
				MarkdownDescription: "API Key Used to connect to the JumpCloud API",
				Required:            true,
				Type:                types.StringType,
				Sensitive:           true,
			},
		},
	}, nil
}

func (p *JumpCloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config JumpCloudProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if config.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing API Key",
			"The provider cannot create the JumpCloud API client due to a missing API Key",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	api_key := os.Getenv("JUMPCLOUD_API_KEY")

	if !config.APIKey.IsNull() {
		api_key = config.APIKey.ValueString()
	}

	if api_key == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing JumpCloud API Key",
			"The provider cannot create the JumpCloud API client due to a missing or empty value for the JumpCloud API Key. "+
				"Set the api_key value in the configuration or use the JUMPCLOUD_API_KEY environment varibale. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	api := JumpCloudApi{
		V1: api.JumpCloudClientApiV1{
			Client: jcapiv1.NewAPIClient(jcapiv1.NewConfiguration()),
			Auth:   context.WithValue(context.TODO(), jcapiv1.ContextAPIKey, jcapiv1.APIKey{Key: api_key}),
		},
		V2: api.JumpCloudClientApiV2{
			Client: jcapiv2.NewAPIClient(jcapiv2.NewConfiguration()),
			Auth:   context.WithValue(context.TODO(), jcapiv2.ContextAPIKey, jcapiv2.APIKey{Key: api_key}),
		},
		Internal: apiclient.New(ctx, api_key, p.version),
	}

	resp.DataSourceData = api
	resp.ResourceData = api
}

func (p *JumpCloudProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewActiveDirectoryResource,
		NewDeviceGroupResource,
		NewUserGroupResource,
	}
}

func (p *JumpCloudProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &JumpCloudProvider{
			version: version,
		}
	}
}
