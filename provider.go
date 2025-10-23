package main

import (
	"context"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &clerkProvider{}
)

// clerkProvider is the provider implementation
type clerkProvider struct {
	version string
}

// clerkProviderModel describes the provider data model
type clerkProviderModel struct {
	APIKey types.String `tfsdk:"api_key"`
}

// New returns a new provider instance
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &clerkProvider{
			version: version,
		}
	}
}

// Metadata returns the provider type name
func (p *clerkProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "clerk"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data
func (p *clerkProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for managing Clerk organizations and resources.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Description: "Clerk API Key. Can also be set via CLERK_API_KEY environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares the provider for data operations
func (p *clerkProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config clerkProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for API key in configuration or environment variable
	apiKey := os.Getenv("CLERK_API_KEY")
	if !config.APIKey.IsNull() {
		apiKey = config.APIKey.ValueString()
	}

	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing API Key Configuration",
			"While configuring the provider, the API key was not found in "+
				"the CLERK_API_KEY environment variable or provider "+
				"configuration block api_key attribute.",
		)
		return
	}

	// Set the global Clerk API key
	clerk.SetKey(apiKey)

	client := &ClerkClient{
		APIKey: apiKey,
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

// Resources defines the resources implemented in the provider
func (p *clerkProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewOrganizationResource,
	}
}

// DataSources defines the data sources implemented in the provider
func (p *clerkProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
