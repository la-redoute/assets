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
	_ datasource.DataSource              = &objectTypeDataSource{}
	_ datasource.DataSourceWithConfigure = &objectTypeDataSource{}
)

type objectTypeDataSource struct {
	client       *assets.Client
	workspace_id string
}

// NewObjectDataSource is a helper function to simplify the provider implementation.
func NewObjectTypeDataSource() datasource.DataSource {
	return &objectTypeDataSource{}
}

// Metadata returns the data source type name.
func (d *objectTypeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objecttype"
}

// Configure adds the provider configured client to the resource.
func (r *objectTypeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *objectTypeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"workspace_id": schema.StringAttribute{
				Computed: true,
			},
			"global_id": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"icon_id": schema.StringAttribute{
				Computed: true,
			},
			"position": schema.Int64Attribute{
				Computed: true,
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
			"parent_object_type_id": schema.StringAttribute{
				Computed:    true,
				Description: "The id of the parent object type",
			},
			"object_schema_id": schema.StringAttribute{
				Computed: true,
			},
			"inherited": schema.BoolAttribute{
				Computed:    true,
				Description: "Describes if this object type is configured for inheritance i.e. it's children inherits the attributes of this object type",
			},
			"abstract_object_type": schema.BoolAttribute{
				Computed: true,
			},
			"parent_object_type_inherited": schema.BoolAttribute{
				Computed:    true,
				Description: "Describes if this object types parent is inherited i.e. this object type has attributes that are inherited from one or more parents",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *objectTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state objectTypeResourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	objecttype, _, err := d.client.ObjectType.Get(ctx, d.workspace_id, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading objecttype",
			"Could not read objecttype, unexpected error: "+err.Error(),
		)
		return
	}

	FillInformationsForObjectType(&state, objecttype)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
