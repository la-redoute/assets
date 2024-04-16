package provider

import (
	"context"
	"fmt"

	"github.com/ctreminiom/go-atlassian/assets"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	//    "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &objectTypeAttributeResource{}
	_ resource.ResourceWithConfigure   = &objectTypeAttributeResource{}
	_ resource.ResourceWithImportState = &objectTypeAttributeResource{}
)

// NewObjectResource is a helper function to simplify the provider implementation.
func NewObjectTypeAttributeResource() resource.Resource {
	return &objectTypeAttributeResource{}
}

// objectResource is the resource implementation.
type objectTypeAttributeResource struct {
	client       *assets.Client
	workspace_id string
}

type objectTypeAttributeResourceModel struct {
	WorkspaceId             types.String `tfsdk:"workspace_id"`
	GlobalId                types.String `tfsdk:"global_id"`
	Id                      types.String `tfsdk:"id"`
	ObjectTypeId            types.String `tfsdk:"object_type_id"`
	Name                    types.String `tfsdk:"name"`
	Label                   types.Bool   `tfsdk:"label"`
	Type                    types.Int64  `tfsdk:"type"`
	Description             types.String `tfsdk:"description"`
	DefaultType             types.Object `tfsdk:"default_type"` //defaultTypeModel
	DefaultTypeId           types.Int64  `tfsdk:"default_type_id"`
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

type defaultTypeModel struct {
	Id   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type referenceTypeModel struct {
	WorkspaceId types.String `tfsdk:"workspace_id"`
	GlobalId    types.String `tfsdk:"global_id"`
	Name        types.String `tfsdk:"name"`
}

// Configure adds the provider configured client to the resource.
func (r *objectTypeAttributeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Metadata returns the resource type name.
func (r *objectTypeAttributeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objecttypeattribute"
}

// Schema defines the schema for the resource.
func (r *objectTypeAttributeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"workspace_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"global_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"object_type_id": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"label": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					int64validator.Any(
						int64validator.OneOf([]int64{2, 4, 7}...),
						int64validator.All(
							int64validator.OneOf([]int64{1}...),
							int64validator.AlsoRequires(
								path.MatchRoot("type_value"),
								path.MatchRoot("additional_value"),
							),
						),
						int64validator.All(
							int64validator.OneOf([]int64{0}...),
							int64validator.AlsoRequires(
								path.MatchRoot("default_type_id"),
							),
						),
					),
				},
			},
			"default_type_id": schema.Int64Attribute{
				Optional: true,
				Validators: []validator.Int64{
					int64validator.Between(-1, 11),
				},
			},
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"default_type": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed: true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"type_value": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type_value_multi": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"additional_value": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"editable": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"system": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"indexed": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"sortable": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"summable": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"minimum_cardinality": schema.Int64Attribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"maximum_cardinality": schema.Int64Attribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"suffix": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"removable": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"object_attribute_exists": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"hidden": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_child_object_types": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"unique_attribute": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"regex_validation": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ql_query": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"options": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"position": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func createObjectTypeAttributePayload(ctx context.Context, objectTypeAttribute objectTypeAttributeResourceModel, payload *models.ObjectTypeAttributePayloadScheme) diag.Diagnostics {
	var diags diag.Diagnostics

	payload.Name = objectTypeAttribute.Name.ValueString()
	payload.Label = objectTypeAttribute.Label.ValueBool()

	if !objectTypeAttribute.Type.IsNull() && !objectTypeAttribute.Type.IsUnknown() {
		Type := int(objectTypeAttribute.Type.ValueInt64())
		payload.Type = &Type
	}

	payload.Description = objectTypeAttribute.Description.ValueString()

	if !objectTypeAttribute.DefaultTypeId.IsNull() && !objectTypeAttribute.DefaultTypeId.IsUnknown() {
		DefaultTypeId := int(objectTypeAttribute.DefaultTypeId.ValueInt64())
		payload.DefaultTypeId = &DefaultTypeId
	}

	payload.TypeValue = objectTypeAttribute.TypeValue.ValueString()

	if !objectTypeAttribute.TypeValueMulti.IsNull() && !objectTypeAttribute.TypeValueMulti.IsUnknown() {
		values := make([]types.String, 0, len(objectTypeAttribute.TypeValueMulti.Elements()))
		diags = objectTypeAttribute.TypeValueMulti.ElementsAs(ctx, &values, false)
		if diags.HasError() {
			return diags
		}

		typeValueMulti := make([]string, 0, len(values))
		for _, value := range values {
			typeValueMulti = append(typeValueMulti, value.ValueString())
		}
		payload.TypeValueMulti = typeValueMulti
	}

	payload.AdditionalValue = objectTypeAttribute.AdditionalValue.ValueString()
	payload.Summable = objectTypeAttribute.Summable.ValueBool()

	if !objectTypeAttribute.MinimumCardinality.IsNull() && !objectTypeAttribute.MinimumCardinality.IsUnknown() {
		MinimumCardinality := int(objectTypeAttribute.MinimumCardinality.ValueInt64())
		payload.MinimumCardinality = &MinimumCardinality
	}

	if !objectTypeAttribute.MaximumCardinality.IsNull() && !objectTypeAttribute.MaximumCardinality.IsUnknown() {
		MaximumCardinality := int(objectTypeAttribute.MaximumCardinality.ValueInt64())
		payload.MaximumCardinality = &MaximumCardinality
	}

	payload.Suffix = objectTypeAttribute.Suffix.ValueString()
	payload.Hidden = objectTypeAttribute.Hidden.ValueBool()
	payload.IncludeChildObjectTypes = objectTypeAttribute.IncludeChildObjectTypes.ValueBool()
	payload.UniqueAttribute = objectTypeAttribute.UniqueAttribute.ValueBool()
	payload.RegexValidation = objectTypeAttribute.RegexValidation.ValueString()
	payload.QlQuery = objectTypeAttribute.RegexValidation.ValueString()
	payload.Options = objectTypeAttribute.Options.ValueString()

	return diags
}

// Create creates the resource and sets the initial Terraform state.
func (r *objectTypeAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan objectTypeAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var payload models.ObjectTypeAttributePayloadScheme

	diags = createObjectTypeAttributePayload(ctx, plan, &payload)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	objectTypeAttribute, _, err := r.client.ObjectTypeAttribute.Create(ctx, r.workspace_id, plan.ObjectTypeId.ValueString(), &payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating objecttypeattribute",
			"Could not create objecttypeattribute, unexpected error: "+err.Error(),
		)
		return
	}

	diags = FillInformationsForObjectTypeAttribute(ctx, &plan, objectTypeAttribute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *objectTypeAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state objectTypeAttributeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	objectTypeAttributes, response, err := r.client.ObjectType.Attributes(ctx, r.workspace_id, state.ObjectTypeId.ValueString(), nil)
	if err != nil {
		if response.Code != 404 {
			resp.Diagnostics.AddError(
				"Error Reading objecttype",
				"Could not read objecttype, unexpected error: "+err.Error(),
			)
		} else {
			resp.State.RemoveResource(ctx)
		}
		return
	}

	var objectTypeAttribute *models.ObjectTypeAttributeScheme

	for _, attr := range objectTypeAttributes {
		if attr.ID == state.Id.ValueString() {
			objectTypeAttribute = attr
			break
		}
	}

	if objectTypeAttribute == nil {
		resp.Diagnostics.AddError(
			"Error Reading objecttypeattribute",
			"Could not read objecttypeattribute, unexpected error: objecttypeattribute not found.",
		)
		return
	}

	diags = FillInformationsForObjectTypeAttribute(ctx, &state, objectTypeAttribute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *objectTypeAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan objectTypeAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var payload models.ObjectTypeAttributePayloadScheme

	diags = createObjectTypeAttributePayload(ctx, plan, &payload)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	objectTypeAttribute, _, err := r.client.ObjectTypeAttribute.Update(ctx, r.workspace_id, plan.ObjectTypeId.ValueString(), plan.Id.ValueString(), &payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating objecttype",
			"Could not update objecttype, unexpected error: "+err.Error(),
		)
		return
	}

	diags = FillInformationsForObjectTypeAttribute(ctx, &plan, objectTypeAttribute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *objectTypeAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state objectTypeAttributeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing object
	_, err := r.client.ObjectTypeAttribute.Delete(ctx, r.workspace_id, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting objecttypeattribute",
			"Could not delete objecttypeattribute, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *objectTypeAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func defaultTypeAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.Int64Type,
		"name": types.StringType,
	}
}

func referenceTypeAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"workspace_id": types.StringType,
		"global_id":    types.StringType,
		"name":         types.StringType,
	}
}
