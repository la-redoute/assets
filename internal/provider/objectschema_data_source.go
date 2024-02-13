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
	_ datasource.DataSource              = &objectSchemaDataSource{}
	_ datasource.DataSourceWithConfigure = &objectSchemaDataSource{}
)

type objectSchemaDataSource struct {
	client       *assets.Client
	workspace_id string
}

// NewObjectDataSource is a helper function to simplify the provider implementation.
func NewObjectSchemaDataSource() datasource.DataSource {
	return &objectSchemaDataSource{}
}

// Metadata returns the data source type name.
func (d *objectSchemaDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objectschema"
}

// Configure adds the provider configured client to the resource.
func (r *objectSchemaDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *objectSchemaDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"workspace_id": schema.StringAttribute{
				Computed: true,
			},
			"global_id": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The object schema id",
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"object_schema_key": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Always 'Ok'",
			},
			"created": schema.StringAttribute{
				Computed: true,
			},
			"updated": schema.StringAttribute{
				Computed: true,
			},
			"object_count": schema.Int64Attribute{
				Computed: true,
			},
			"object_type_count": schema.Int64Attribute{
				Computed: true,
			},
			"can_manage": schema.BoolAttribute{
				Computed: true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *objectSchemaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state objectSchemaResourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	objectschema, _, err := d.client.ObjectSchema.Get(ctx, d.workspace_id, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading objectschema",
			"Could not read objectschema, unexpected error: "+err.Error(),
		)
		return
	}

	FillInformationsForObjectSchema(&state, objectschema)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
