package provider

import (
	"context"
	"fmt"

	"github.com/ctreminiom/go-atlassian/assets"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &objectTypeAttributesDataSource{}
	_ datasource.DataSourceWithConfigure = &objectTypeAttributesDataSource{}
)

type objectTypeAttributesDataSource struct {
	client       *assets.Client
	workspace_id string
}

type objectTypeAttributesDataSourceModel struct {
	ObjectTypeId   types.String                   `tfsdk:"objecttype_id"`
	ObjectSchemaId types.String                   `tfsdk:"objectschema_id"`
	Attributes     []objectTypeAttributeDataModel `tfsdk:"attributes"`
}

type objectTypeAttributeDataModel struct {
	WorkspaceId             types.String `tfsdk:"workspace_id"`
	GlobalId                types.String `tfsdk:"global_id"`
	Id                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Label                   types.Bool   `tfsdk:"label"`
	Type                    types.Int64  `tfsdk:"type"`
	Description             types.String `tfsdk:"description"`
	DefaultType             types.Object `tfsdk:"default_type"` //defaultTypeModel
	TypeValue               types.String `tfsdk:"type_value"`
	TypeValueMulti          types.List   `tfsdk:"type_value_multi"` //String List
	AdditionalValue         types.String `tfsdk:"additional_value"`
	ReferenceType           types.Object `tfsdk:"reference_type"` //referenceTypeModel
	ReferenceObjectTypeId   types.String `tfsdk:"reference_object_type_id"`
	Editable                types.Bool   `tfsdk:"editable"`
	System                  types.Bool   `tfsdk:"system"`
	Indexed                 types.Bool   `tfsdk:"indexed"`
	Sortable                types.Bool   `tfsdk:"sortable"`
	Summable                types.Bool   `tfsdk:"summable"`
	MinimumCardinality      types.Int64  `tfsdk:"minimum_cardinality"`
	MaximumCardinality      types.Int64  `tfsdk:"maximum_cardinality"`
	Suffix                  types.String `tfsdk:"suffix"`
	Removable               types.Bool   `tfsdk:"removable"`
	ObjectAttributeExists   types.Bool   `tfsdk:"object_attribute_exists"`
	Hidden                  types.Bool   `tfsdk:"hidden"`
	IncludeChildObjectTypes types.Bool   `tfsdk:"include_child_object_types"`
	UniqueAttribute         types.Bool   `tfsdk:"unique_attribute"`
	RegexValidation         types.String `tfsdk:"regex_validation"`
	QlQuery                 types.String `tfsdk:"ql_query"`
	Options                 types.String `tfsdk:"options"`
	Position                types.Int64  `tfsdk:"position"`
}

// NewObjectDataSource is a helper function to simplify the provider implementation.
func NewObjectTypeAttributesDataSource() datasource.DataSource {
	return &objectTypeAttributesDataSource{}
}

// Metadata returns the data source type name.
func (d *objectTypeAttributesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objecttypeattributes"
}

// Configure adds the provider configured client to the resource.
func (r *objectTypeAttributesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *objectTypeAttributesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"objecttype_id": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("objecttype_id"),
						path.MatchRoot("objectschema_id"),
					),
				},
			},
			"objectschema_id": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("objecttype_id"),
						path.MatchRoot("objectschema_id"),
					),
				},
			},
			"attributes": schema.SetNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"workspace_id": schema.StringAttribute{
							Computed: true,
						},
						"global_id": schema.StringAttribute{
							Computed: true,
						},
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"label": schema.BoolAttribute{
							Computed: true,
						},
						"type": schema.Int64Attribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"default_type": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"id": schema.Int64Attribute{
									Computed: true,
								},
								"name": schema.StringAttribute{
									Computed: true,
								},
							},
						},
						"type_value": schema.StringAttribute{
							Computed: true,
						},
						"type_value_multi": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"additional_value": schema.StringAttribute{
							Computed: true,
						},
						"reference_type": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"workspace_id": schema.StringAttribute{
									Computed: true,
								},
								"global_id": schema.StringAttribute{
									Computed: true,
								},
								"name": schema.StringAttribute{
									Computed: true,
								},
							},
						},
						"reference_object_type_id": schema.StringAttribute{
							Computed: true,
						},
						"editable": schema.BoolAttribute{
							Computed: true,
						},
						"system": schema.BoolAttribute{
							Computed: true,
						},
						"indexed": schema.BoolAttribute{
							Computed:    true,
							Description: "Describes if this object type attribute is indexed. For an indexed attribute the AQL search will be faster, but this will affect memory consumption.",
						},
						"sortable": schema.BoolAttribute{
							Computed: true,
						},
						"summable": schema.BoolAttribute{
							Computed: true,
						},
						"minimum_cardinality": schema.Int64Attribute{
							Computed: true,
						},
						"maximum_cardinality": schema.Int64Attribute{
							Computed: true,
						},
						"suffix": schema.StringAttribute{
							Computed: true,
						},
						"removable": schema.BoolAttribute{
							Computed: true,
						},
						"object_attribute_exists": schema.BoolAttribute{
							Computed: true,
						},
						"hidden": schema.BoolAttribute{
							Computed: true,
						},
						"include_child_object_types": schema.BoolAttribute{
							Computed: true,
						},
						"unique_attribute": schema.BoolAttribute{
							Computed: true,
						},
						"regex_validation": schema.StringAttribute{
							Computed: true,
						},
						"ql_query": schema.StringAttribute{
							Computed: true,
						},
						"options": schema.StringAttribute{
							Computed: true,
						},
						"position": schema.Int64Attribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *objectTypeAttributesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state objectTypeAttributesDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var objectTypeAttributes []*models.ObjectTypeAttributeScheme
	var err error
	if state.ObjectTypeId.ValueString() != "" {
		objectTypeAttributes, _, err = d.client.ObjectType.Attributes(ctx, d.workspace_id, state.ObjectTypeId.ValueString(), nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading objecttypeattributes",
				"Could not read objecttypeattributes, unexpected error: "+err.Error(),
			)
			return
		}
	} else if state.ObjectSchemaId.ValueString() != "" {
		var payload *models.ObjectSchemaAttributesParamsScheme = &models.ObjectSchemaAttributesParamsScheme{
			Extended: true,
		}
		objectTypeAttributes, _, err = d.client.ObjectSchema.Attributes(ctx, d.workspace_id, state.ObjectSchemaId.ValueString(), payload)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading objectschemaattributes",
				"Could not read objectschemaattributes, unexpected error: "+err.Error(),
			)
			return
		}
	} else {
		resp.Diagnostics.AddError(
			"Error Reading attributes",
			"objecttype_id and objectschema_id are null.",
		)
		return
	}

	state.Attributes = make([]objectTypeAttributeDataModel, len(objectTypeAttributes))

	for index, attr := range objectTypeAttributes {
		diags = FillInformationsForDataObjectTypeAttribute(ctx, &state.Attributes[index], attr)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
