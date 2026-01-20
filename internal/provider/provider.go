package provider

import (
	"context"
	"terraform-provider-hrobot/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const defaultURL string = "https://robot-ws.your-server.de"

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &hrobotProvider{}
)

type hrobotProviderModel struct {
	BaseURL  types.String `tfsdk:"url"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &hrobotProvider{
			version: version,
		}
	}
}

// hrobotProvider is the provider implementation.
type hrobotProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *hrobotProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "hrobot"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *hrobotProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Optional: true,
			},
			"username": schema.StringAttribute{
				Required: true,
			},
			"password": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

// Configure prepares a hrobot API client for data sources and resources.
func (p *hrobotProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config hrobotProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var baseURL, username, password string
	if !config.BaseURL.IsNull() {
		baseURL = config.BaseURL.ValueString()
	}
	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}
	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	if baseURL == "" {
		baseURL = defaultURL
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Hetzner Robot username",
			"The provider cannot create the Hetzner Robot HTTP client as there is a missing or empty value for the Robot webservice username. "+
				"Set the username value in the configuration. ",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Hetzner Robot password",
			"The provider cannot create the Hetzner Robot HTTP client as there is a missing or empty value for the Robot webservice password. "+
				"Set the password value in the configuration. ",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	c := client.NewHetznerRobotClient(baseURL, username, password)
	resp.DataSourceData = c
	resp.ResourceData = c
}

// DataSources defines the data sources implemented in the provider.
func (p *hrobotProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewVSwitchDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *hrobotProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
