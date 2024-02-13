package provider

import (
	"context"
	"fmt"

	"github.com/ctreminiom/go-atlassian/assets"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &globalIconsDataSource{}
	_ datasource.DataSourceWithConfigure = &globalIconsDataSource{}
)

type globalIconsDataSource struct {
	client       *assets.Client
	workspace_id string
}

type GlobalIconModel struct {
	Icons []IconModel `tfsdk:"icons"`
}

// NewObjectDataSource is a helper function to simplify the provider implementation.
func NewGlobalIconsDataSource() datasource.DataSource {
	return &globalIconsDataSource{}
}

// Metadata returns the data source type name.
func (d *globalIconsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_global_icons"
}

// Configure adds the provider configured client to the resource.
func (r *globalIconsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *globalIconsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"icons": schema.SetNestedAttribute{
				Computed:    true,
				Description: "All existing global icons",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"url16": schema.StringAttribute{
							Computed:    true,
							Description: "A url to the icon to display with small resolution.",
						},
						"url48": schema.StringAttribute{
							Computed:    true,
							Description: "A url to the icon to display with large resolution",
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *globalIconsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	// Get current state
	var state GlobalIconModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	icons, _, err := d.client.Icon.Global(ctx, d.workspace_id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Global icons",
			"Could not read Global icons, unexpected error: "+err.Error(),
		)
		return
	}

	var global_icons = make([]IconModel, len(icons))

	for index, icon := range icons {
		FillInformationForIcon(&global_icons[index], icon)
	}

	state.Icons = global_icons

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
