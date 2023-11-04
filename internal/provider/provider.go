// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/resendlabs/resend-go"
)

// Ensure ResendProvider satisfies various provider interfaces.
var _ provider.Provider = &ResendProvider{}

// ResendProvider defines the provider implementation.
type ResendProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ResendProviderModel describes the provider data model.
type ResendProviderModel struct {
	ApiKey types.String `tfsdk:"api_key"`
}

func (p *ResendProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "resend"
	resp.Version = p.version
}

func (p *ResendProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "A resend API key",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *ResendProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config ResendProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown API key",
			"The provider cannot create the Resend API client if the key is unknown. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the RESEND_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}
	log.Println(os.Environ())

	apiKey := os.Getenv("RESEND_API_KEY")
	if !config.ApiKey.IsNull() {
		apiKey = config.ApiKey.ValueString()
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing api key",
			"The provider cannot create the Resend API client as there is a missing or empty value for the Resend api key. "+
				"Set the api_key value in the configuration or use the RESEND_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)

		return
	}
	tflog.Info(ctx, fmt.Sprintf("Creating Resend API client %s", apiKey))
	client := resend.NewClient(config.ApiKey.ValueString())
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *ResendProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDomainResource,
		NewApiKeyResource,
	}
}

func (p *ResendProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ResendProvider{
			version: version,
		}
	}
}
