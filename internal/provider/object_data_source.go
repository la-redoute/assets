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
	_ datasource.DataSource              = &objectDataSource{}
	_ datasource.DataSourceWithConfigure = &objectDataSource{}
)

type objectDataSource struct {
	client       *assets.Client
	workspace_id string
}

type objectDataResourceModel struct {
	WorkspaceId  types.String `tfsdk:"workspace_id"`
	GlobalId     types.String `tfsdk:"global_id"`
	Id           types.String `tfsdk:"id"`
	Label        types.String `tfsdk:"label"`
	ObjectKey    types.String `tfsdk:"object_key"`
	ObjectTypeId types.String `tfsdk:"object_type_id"`
	Created      types.String `tfsdk:"created"`
	Updated      types.String `tfsdk:"updated"`
	HasAvatar    types.Bool   `tfsdk:"has_avatar"`
	Attributes   types.Set    `tfsdk:"attributes"` //<<[]objectAttributeModel
	Links        types.Object `tfsdk:"links"`      //<<objectModel
	Avatar       types.Object `tfsdk:"avatar"`     //<<avatarModel
}

// NewObjectDataSource is a helper function to simplify the provider implementation.
func NewObjectDataSource() datasource.DataSource {
	return &objectDataSource{}
}

// Metadata returns the data source type name.
func (d *objectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_object"
}

// Configure adds the provider configured client to the resource.
func (r *objectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *objectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
				Description: "The object id to operate on",
			},
			"label": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the object. This value is fetched from the attribute that is currently marked as label for the object type of this object",
			},
			"object_key": schema.StringAttribute{
				Computed:    true,
				Description: "The external identifier for this object",
			},
			"avatar": schema.ObjectAttribute{
				AttributeTypes: avatarAttrTypes(),
				Computed:       true,
				Description:    "The object avatar is a custom image that represents an object. If the object has no avatar the icon for the object type will be used",
			},
			"object_type_id": schema.StringAttribute{
				Computed:    true,
				Description: "The Assets object type",
			},
			"created": schema.StringAttribute{
				Computed: true,
			},
			"updated": schema.StringAttribute{
				Computed: true,
			},
			"has_avatar": schema.BoolAttribute{
				Computed: true,
			},
			"attributes": schema.SetAttribute{
				Computed: true,
				ElementType: types.ObjectType{
					AttrTypes: objectAttributeAttrTypes(),
				},
			},
			"links": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"self": schema.StringAttribute{
						Computed: true,
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *objectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state objectDataResourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, _, err := d.client.Object.Get(ctx, d.workspace_id, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading object",
			"Could not read object, unexpected error: "+err.Error(),
		)
		return
	}

	FillInformationsForDataObject(ctx, &state, object)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
