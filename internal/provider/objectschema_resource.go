package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/ctreminiom/go-atlassian/assets"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &objectSchemaResource{}
	_ resource.ResourceWithConfigure   = &objectSchemaResource{}
	_ resource.ResourceWithImportState = &objectSchemaResource{}
)

// NewObjectResource is a helper function to simplify the provider implementation.
func NewObjectSchemaResource() resource.Resource {
	return &objectSchemaResource{}
}

// objectResource is the resource implementation.
type objectSchemaResource struct {
	client       *assets.Client
	workspace_id string
}

type objectSchemaResourceModel struct {
	WorkspaceId     types.String `tfsdk:"workspace_id"`
	GlobalId        types.String `tfsdk:"global_id"`
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	ObjectSchemaKey types.String `tfsdk:"object_schema_key"`
	Description     types.String `tfsdk:"description"`
	Status          types.String `tfsdk:"status"`
	Created         types.String `tfsdk:"created"`
	Updated         types.String `tfsdk:"updated"`
	ObjectCount     types.Int64  `tfsdk:"object_count"`
	ObjectTypeCount types.Int64  `tfsdk:"object_type_count"`
	CanManage       types.Bool   `tfsdk:"can_manage"`
}

// Configure adds the provider configured client to the resource.
func (r *objectSchemaResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *objectSchemaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objectschema"
}

// Schema defines the schema for the resource.
func (r *objectSchemaResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 50),
				},
			},
			"object_schema_key": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 10),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[A-Z0-9]*$`),
						"must contain only uppercase alphanumeric characters",
					),
				},
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"status": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Always 'Ok'",
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
			"object_type_count": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"can_manage": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func createObjectSchemaPayload(objectSchema objectSchemaResourceModel, payload *models.ObjectSchemaPayloadScheme) {
	payload.Name = objectSchema.Name.ValueString()
	payload.Description = objectSchema.Description.ValueString()
	payload.ObjectSchemaKey = objectSchema.ObjectSchemaKey.ValueString()
}

// Create creates the resource and sets the initial Terraform state.
func (r *objectSchemaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan objectSchemaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var payload models.ObjectSchemaPayloadScheme

	createObjectSchemaPayload(plan, &payload)

	objectSchema, _, err := r.client.ObjectSchema.Create(ctx, r.workspace_id, &payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating objectschema",
			"Could not create objectschema, unexpected error: "+err.Error(),
		)
		return
	}

	FillInformationsForObjectSchema(&plan, objectSchema)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *objectSchemaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state objectSchemaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	objectSchema, _, err := r.client.ObjectSchema.Get(ctx, r.workspace_id, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading objectschema",
			"Could not read objectschema, unexpected error: "+err.Error(),
		)
		return
	}

	FillInformationsForObjectSchema(&state, objectSchema)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *objectSchemaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan objectSchemaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var payload models.ObjectSchemaPayloadScheme

	createObjectSchemaPayload(plan, &payload)

	objectSchema, _, err := r.client.ObjectSchema.Update(ctx, r.workspace_id, plan.Id.ValueString(), &payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating objectschema",
			"Could not update objectschema, unexpected error: "+err.Error(),
		)
		return
	}

	FillInformationsForObjectSchema(&plan, objectSchema)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *objectSchemaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state objectSchemaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing object
	_, _, err := r.client.ObjectSchema.Delete(ctx, r.workspace_id, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting objectschema",
			"Could not delete objectschema, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *objectSchemaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
