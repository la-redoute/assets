// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"
	"strconv"

	"github.com/ctreminiom/go-atlassian/assets"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure AssetsProvider satisfies various provider interfaces.
var _ provider.Provider = &AssetsProvider{}

// AssetsProvider defines the provider implementation.
type AssetsProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// AssetsProviderModel describes the provider data model.
type AssetsProviderModel struct {
	Host        types.String `tfsdk:"host"`
	Token       types.String `tfsdk:"token"`
	Mail        types.String `tfsdk:"mail"`
	WorkspaceId types.String `tfsdk:"workspace_id"`
	Features    types.Object `tfsdk:"features"` // <<featuresModel
}

type featuresModel struct {
	DestroyObject                 types.Bool   `tfsdk:"destroy_object"`
	ObsoleteObjectTypeAttributeId types.String `tfsdk:"obsolete_objecttypeattribute_id"`
}

type features struct {
	DestroyObject                 bool
	ObsoleteObjectTypeAttributeId string
}

// Custom client to store the workspace ID
type AssetsProviderClient struct {
	Client      *assets.Client
	WorkspaceId string
	Features    *features
}

func (p *AssetsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "assets"
	resp.Version = p.version
}

func (p *AssetsProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional:    true,
				Description: "URL for your Jira/Confluence instance.",
			},
			"token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "A personal access token for authentication.",
			},
			"mail": schema.StringAttribute{
				Optional:    true,
				Description: "The mail of the PAT account.",
			},
			"workspace_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the Atlassian workspace.",
			},
			"features": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"destroy_object": schema.BoolAttribute{
						Optional:    true,
						Description: "Destroy object ? If false, obsolete_objecttypeattribute_id must be defined. Defaults to true.",
					},
					"obsolete_objecttypeattribute_id": schema.StringAttribute{
						Optional:    true,
						Description: "The objecttypeattribute ID of the obsolete attribute.",
					},
				},
			},
		},
	}
}

func (p *AssetsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config AssetsProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Assets API Token",
			"The provider cannot create the Assets API client as there is an unknown configuration value for the Assets API Token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ATLASSIAN_TOKEN environment variable.",
		)
	}

	if config.Mail.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("mail"),
			"Unknown Assets API Mail",
			"The provider cannot create the Assets API client as there is an unknown configuration value for the Assets API Mail. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ATLASSIAN_MAIL environment variable.",
		)
	}

	if config.WorkspaceId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("workspace_id"),
			"Unknown Assets Workspace ID",
			"The provider cannot create the Assets API client as there is an unknown configuration value for the Assets Workspace ID. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ASSETS_WORKSPACE_ID environment variable.",
		)
	}

	if config.Features.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("features"),
			"Unknown Features",
			"The provider cannot create the Assets API client as there is an unknown configuration value for the Features.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	host := os.Getenv("ATLASSIAN_HOST")
	token := os.Getenv("ATLASSIAN_TOKEN")
	mail := os.Getenv("ATLASSIAN_MAIL")
	workspace_id := os.Getenv("ASSETS_WORKSPACE_ID")
	destroy_object_env := os.Getenv("ASSETS_DESTROY_OBJECT")
	obsolete_objecttypeattribute_id := os.Getenv("ASSETS_OBJECTTYPEATTRIBUTE_ID")

	destroy_object, err := strconv.ParseBool(destroy_object_env)
	if err != nil {
		destroy_object = true
	}

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	if !config.Mail.IsNull() {
		mail = config.Mail.ValueString()
	}

	if !config.WorkspaceId.IsNull() {
		workspace_id = config.WorkspaceId.ValueString()
	}

	var feats featuresModel
	if !config.Features.IsNull() {
		diags = config.Features.As(ctx, &feats, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		if !feats.DestroyObject.IsNull() {
			destroy_object = feats.DestroyObject.ValueBool()
		}

		if !feats.ObsoleteObjectTypeAttributeId.IsNull() {
			obsolete_objecttypeattribute_id = feats.ObsoleteObjectTypeAttributeId.ValueString()
		}
	}

	features := features{
		DestroyObject:                 destroy_object,
		ObsoleteObjectTypeAttributeId: obsolete_objecttypeattribute_id,
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Assets API Token",
			"The provider cannot create the Assets API client as there is a missing or empty value for the Assets API Token. "+
				"Set the pat value in the configuration or use the ATLASSIAN_PAT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if mail == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("mail"),
			"Missing Assets API Mail",
			"The provider cannot create the Assets API client as there is a missing or empty value for the Assets API Mail. "+
				"Set the pat value in the configuration or use the ATLASSIAN_MAIL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if workspace_id == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("workspace_id"),
			"Missing Assets Workspace ID",
			"The provider cannot create the Assets API client as there is a missing or empty value for the Assets Workspace ID. "+
				"Set the pat value in the configuration or use the ASSETS_WORKSPACE_ID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if !destroy_object && (obsolete_objecttypeattribute_id == "") {
		resp.Diagnostics.AddAttributeError(
			path.Root("features"),
			"Unknown obsolete_objecttypeattribute_id",
			"The provider cannot create the Assets API client as there is an unknown configuration value for the Features.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new Assets client using the configuration values
	client, err := assets.New(nil, host)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Assets API Client",
			"An unexpected error occurred when creating the Assets API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Assets Client Error: "+err.Error(),
		)
		return
	}
	client.Auth.SetBasicAuth(mail, token)

	assetsClient := AssetsProviderClient{
		Client:      client,
		WorkspaceId: workspace_id,
		Features:    &features,
	}

	// Make the Assets client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = assetsClient
	resp.ResourceData = assetsClient
}

func (p *AssetsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewObjectResource,
		NewObjectTypeResource,
		NewObjectTypeAttributeResource,
		NewObjectSchemaResource,
	}
}

func (p *AssetsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewObjectDataSource,
		NewIconDataSource,
		NewGlobalIconsDataSource,
		NewObjectTypeDataSource,
		NewObjectTypeAttributesDataSource,
		NewObjectSchemaDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AssetsProvider{
			version: version,
		}
	}
}
