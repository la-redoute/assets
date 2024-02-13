package provider

import (
	"context"
	"fmt"

	"github.com/ctreminiom/go-atlassian/assets"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &iconDataSource{}
	_ datasource.DataSourceWithConfigure = &iconDataSource{}
)

type iconDataSource struct {
	client       *assets.Client
	workspace_id string
}

type IconModel struct {
	Id    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Url16 types.String `tfsdk:"url16"`
	Url48 types.String `tfsdk:"url48"`
}

// NewObjectDataSource is a helper function to simplify the provider implementation.
func NewIconDataSource() datasource.DataSource {
	return &iconDataSource{}
}

// Metadata returns the data source type name.
func (d *iconDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_icon"
}

// Configure adds the provider configured client to the resource.
func (r *iconDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	assetsClient, ok := req.ProviderData.(AssetsProviderClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *assets.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = assetsClient.Client
	r.workspace_id = assetsClient.WorkspaceId
}

// Schema defines the schema for the data source.
func (d *iconDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"url16": schema.StringAttribute{
				Computed:    true,
				Description: "A url to the icon to display with small resolution",
			},
			"url48": schema.StringAttribute{
				Computed:    true,
				Description: "A url to the icon to display with large resolution",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *iconDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state IconModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	icon, _, err := d.client.Icon.Get(ctx, d.workspace_id, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading icon",
			"Could not read icon, unexpected error: "+err.Error(),
		)
		return
	}

	FillInformationForIcon(&state, icon)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
