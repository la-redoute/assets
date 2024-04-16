package provider

import (
	"context"
	"fmt"

	"github.com/ctreminiom/go-atlassian/assets"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &objectTypeResource{}
	_ resource.ResourceWithConfigure   = &objectTypeResource{}
	_ resource.ResourceWithImportState = &objectTypeResource{}
)

// NewObjectResource is a helper function to simplify the provider implementation.
func NewObjectTypeResource() resource.Resource {
	return &objectTypeResource{}
}

// objectResource is the resource implementation.
type objectTypeResource struct {
	client       *assets.Client
	workspace_id string
}

type objectTypeResourceModel struct {
	WorkspaceId               types.String `tfsdk:"workspace_id"`
	GlobalId                  types.String `tfsdk:"global_id"`
	Id                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	Description               types.String `tfsdk:"description"`
	IconId                    types.String `tfsdk:"icon_id"`
	Position                  types.Int64  `tfsdk:"position"`
	Created                   types.String `tfsdk:"created"`
	Updated                   types.String `tfsdk:"updated"`
	ObjectCount               types.Int64  `tfsdk:"object_count"`
	ParentObjectTypeId        types.String `tfsdk:"parent_object_type_id"`
	ObjectSchemaId            types.String `tfsdk:"object_schema_id"`
	Inherited                 types.Bool   `tfsdk:"inherited"`
	AbstractObjectType        types.Bool   `tfsdk:"abstract_object_type"`
	ParentObjectTypeInherited types.Bool   `tfsdk:"parent_object_type_inherited"`
}

// Configure adds the provider configured client to the resource.
func (r *objectTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *objectTypeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objecttype"
}

// Schema defines the schema for the resource.
func (r *objectTypeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"icon_id": schema.StringAttribute{
				Required: true,
			},
			"position": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"created": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated": schema.StringAttribute{
				Computed: true,
			},
			"object_count": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"parent_object_type_id": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The id of the parent object type",
			},
			"object_schema_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"inherited": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
				Description: "Describes if this object type is configured for inheritance i.e. it's children inherits the attributes of this object type",
			},
			"abstract_object_type": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"parent_object_type_inherited": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
				Description: "Describes if this object types parent is inherited i.e. this object type has attributes that are inherited from one or more parents",
			},
		},
	}
}

func createObjectTypePayload(objectType objectTypeResourceModel, payload *models.ObjectTypePayloadScheme) {
	payload.Name = objectType.Name.ValueString()
	payload.Description = objectType.Description.ValueString()
	payload.IconId = objectType.IconId.ValueString()
	payload.ObjectSchemaId = objectType.ObjectSchemaId.ValueString()
	payload.ParentObjectTypeId = objectType.ParentObjectTypeId.ValueString()
	payload.Inherited = objectType.Inherited.ValueBool()
	payload.AbstractObjectType = objectType.AbstractObjectType.ValueBool()
}

// Create creates the resource and sets the initial Terraform state.
func (r *objectTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan objectTypeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var payload models.ObjectTypePayloadScheme

	createObjectTypePayload(plan, &payload)

	objectType, _, err := r.client.ObjectType.Create(ctx, r.workspace_id, &payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating objecttype",
			"Could not create objecttype, unexpected error: "+err.Error(),
		)
		return
	}

	FillInformationsForObjectType(&plan, objectType)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *objectTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state objectTypeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	objectType, response, err := r.client.ObjectType.Get(ctx, r.workspace_id, state.Id.ValueString())
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

	FillInformationsForObjectType(&state, objectType)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *objectTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan objectTypeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var payload models.ObjectTypePayloadScheme

	createObjectTypePayload(plan, &payload)

	objectType, _, err := r.client.ObjectType.Update(ctx, r.workspace_id, plan.Id.ValueString(), &payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating objecttype",
			"Could not update objecttype, unexpected error: "+err.Error(),
		)
		return
	}

	FillInformationsForObjectType(&plan, objectType)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *objectTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state objectTypeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing object
	_, _, err := r.client.ObjectType.Delete(ctx, r.workspace_id, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting objecttype",
			"Could not delete objecttype, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *objectTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
