package provider

import (
	"context"
	"fmt"

	"github.com/ctreminiom/go-atlassian/assets"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	//    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	//    "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &objectResource{}
	_ resource.ResourceWithConfigure   = &objectResource{}
	_ resource.ResourceWithImportState = &objectResource{}
)

// NewObjectResource is a helper function to simplify the provider implementation.
func NewObjectResource() resource.Resource {
	return &objectResource{}
}

// objectResource is the resource implementation.
type objectResource struct {
	client       *assets.Client
	workspace_id string
	features     *features
}

type objectResourceModel struct {
	WorkspaceId types.String `tfsdk:"workspace_id"`
	GlobalId    types.String `tfsdk:"global_id"`
	Id          types.String `tfsdk:"id"`
	Label       types.String `tfsdk:"label"`
	ObjectKey   types.String `tfsdk:"object_key"`
	//AvatarUuid   types.String `tfsdk:"avatar_uuid"`
	ObjectTypeId types.String              `tfsdk:"object_type_id"`
	Created      types.String              `tfsdk:"created"`
	Updated      types.String              `tfsdk:"updated"`
	HasAvatar    types.Bool                `tfsdk:"has_avatar"`
	Attributes   types.Set                 `tfsdk:"attributes"` //<<[]objectAttributeModel
	Links        types.Object              `tfsdk:"links"`      //<<objectModel
	AttributesIn []*objectAttributeInModel `tfsdk:"attributes_in"`
	Avatar       types.Object              `tfsdk:"avatar"` //<<avatarModel
}

type avatarModel struct {
	WorkspaceId types.String `tfsdk:"workspace_id"`
	GlobalId    types.String `tfsdk:"global_id"`
	Id          types.String `tfsdk:"id"`
	AvatarUuid  types.String `tfsdk:"avatar_uuid"`
	Url16       types.String `tfsdk:"url16"`
	Url48       types.String `tfsdk:"url48"`
	Url72       types.String `tfsdk:"url72"`
	Url144      types.String `tfsdk:"url144"`
	Url288      types.String `tfsdk:"url288"`
	ObjectId    types.String `tfsdk:"object_id"`
}

type objectAttributeInModel struct {
	ObjectTypeAttributeId   types.String                   `tfsdk:"object_type_attribute_id"`
	ObjectAttributeValuesIn []*objectAttributeValueInModel `tfsdk:"object_attribute_values_in"`
}

type objectAttributeValueInModel struct {
	Value types.String `tfsdk:"value"`
}

type objectAttributeModel struct {
	WorkspaceId              types.String `tfsdk:"workspace_id"`
	GlobalId                 types.String `tfsdk:"global_id"`
	Id                       types.String `tfsdk:"id"`
	ObjectTypeAttributeId    types.String `tfsdk:"object_type_attribute_id"`
	ObjectAttributeValues    types.Set    `tfsdk:"object_attribute_values"`
	ObjectTypeAttributeLabel types.Bool   `tfsdk:"object_type_attribute_label"`
}

type objectAttributeValueModel struct {
	Value           types.String `tfsdk:"value"`
	DisplayValue    types.String `tfsdk:"display_value"`
	SearchValue     types.String `tfsdk:"search_value"`
	Group           types.Object `tfsdk:"group"`
	Status          types.Object `tfsdk:"status"`
	AdditionalValue types.String `tfsdk:"additional_value"`
}

type groupModel struct {
	AvatarUrl types.String `tfsdk:"avatar_url"`
	Name      types.String `tfsdk:"name"`
}

type statusModel struct {
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Category types.Int64  `tfsdk:"category"`
}

type objectModel struct {
	Self types.String `tfsdk:"self"`
}

// Configure adds the provider configured client to the resource.
func (r *objectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.features = assetsClient.Features
}

// Metadata returns the resource type name.
func (r *objectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_object"
}

// Schema defines the schema for the resource.
func (r *objectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"label": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					SyncLabelPlanModifier(),
				},
				Description: "The name of the object. This value is fetched from the attribute that is currently marked as label for the object type of this object",
			},
			"object_key": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The external identifier for this object",
			},
			"object_type_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The object type determines where the object should be stored and which attributes are available",
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
			"has_avatar": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"links": schema.ObjectAttribute{
				AttributeTypes: objectAttrTypes(),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Computed: true,
			},
			"attributes_in": schema.SetNestedAttribute{
				Required: true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"object_type_attribute_id": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
							Description: "The type of the attribute. The type decides how this value should be interpreted",
						},
						"object_attribute_values_in": schema.SetNestedAttribute{
							Required: true,
							PlanModifiers: []planmodifier.Set{
								setplanmodifier.UseStateForUnknown(),
							},
							Description: "The value(s)",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"value": schema.StringAttribute{
										Required: true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
										},
									},
								},
							},
						},
					},
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
						"object_type_attribute_id": schema.StringAttribute{
							Computed: true,
						},
						"object_type_attribute_label": schema.BoolAttribute{
							Computed: true,
						},
						"object_attribute_values": schema.SetNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"value": schema.StringAttribute{
										Computed: true,
									},
									"display_value": schema.StringAttribute{
										Computed: true,
									},
									"search_value": schema.StringAttribute{
										Computed: true,
									},
									"group": schema.SingleNestedAttribute{
										Computed: true,
										Attributes: map[string]schema.Attribute{
											"avatar_url": schema.StringAttribute{
												Computed: true,
											},
											"name": schema.StringAttribute{
												Computed: true,
											},
										},
									},
									"status": schema.SingleNestedAttribute{
										Computed: true,
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Computed: true,
											},
											"name": schema.StringAttribute{
												Computed: true,
											},
											"category": schema.Int64Attribute{
												Computed: true,
											},
										},
									},
									"additional_value": schema.StringAttribute{
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"avatar": schema.SingleNestedAttribute{
				Computed: true,
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"workspace_id": schema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							SyncAvatarPlanModifier(),
						},
					},
					"global_id": schema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							SyncAvatarPlanModifier(),
						},
					},
					"id": schema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							SyncAvatarPlanModifier(),
						},
					},
					"avatar_uuid": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"url16": schema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							SyncAvatarPlanModifier(),
						},
					},
					"url48": schema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							SyncAvatarPlanModifier(),
						},
					},
					"url72": schema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							SyncAvatarPlanModifier(),
						},
					},
					"url144": schema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							SyncAvatarPlanModifier(),
						},
					},
					"url288": schema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							SyncAvatarPlanModifier(),
						},
					},
					"object_id": schema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							SyncAvatarPlanModifier(),
						},
						Description: "A reference to the object that this avatar is associated with",
					},
				},
			},
		},
	}
}

func createObjectPayload(ctx context.Context, object objectResourceModel, payload *models.ObjectPayloadScheme) diag.Diagnostics {
	var diags diag.Diagnostics

	var payloadAttributes []*models.ObjectPayloadAttributeScheme

	for _, attr := range object.AttributesIn {
		var values []*models.ObjectPayloadAttributeValueScheme

		for _, value := range attr.ObjectAttributeValuesIn {
			values = append(values, &models.ObjectPayloadAttributeValueScheme{
				Value: value.Value.ValueString(),
			})
		}

		payloadAttributes = append(payloadAttributes, &models.ObjectPayloadAttributeScheme{
			ObjectTypeAttributeID: attr.ObjectTypeAttributeId.ValueString(),
			ObjectAttributeValues: values,
		})
	}

	payload.ObjectTypeID = object.ObjectTypeId.ValueString()
	payload.Attributes = payloadAttributes
	payload.HasAvatar = object.HasAvatar.ValueBool()

	if !object.Avatar.IsNull() && !object.Avatar.IsUnknown() {
		var avatar avatarModel
		diags = object.Avatar.As(ctx, &avatar, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return diags
		}
		payload.AvatarUUID = avatar.AvatarUuid.ValueString()
	}

	return diags
}

// Create creates the resource and sets the initial Terraform state.
func (r *objectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan objectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var payload models.ObjectPayloadScheme

	diags = createObjectPayload(ctx, plan, &payload)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	object_without_attributes, _, err := r.client.Object.Create(ctx, r.workspace_id, &payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating object",
			"Could not create object, unexpected error: "+err.Error(),
		)
		return
	}

	object, _, err := r.client.Object.Get(ctx, r.workspace_id, object_without_attributes.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error get attributes",
			"Could not get attributes, unexpected error: "+err.Error(),
		)
		return
	}

	diags = FillInformationsForObject(ctx, &plan, object)
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
func (r *objectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state objectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, _, err := r.client.Object.Get(ctx, r.workspace_id, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading object",
			"Could not read object, unexpected error: "+err.Error(),
		)
		return
	}

	diags = FillInformationsForObject(ctx, &state, object)
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
func (r *objectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan objectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var payload models.ObjectPayloadScheme

	diags = createObjectPayload(ctx, plan, &payload)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, _, err := r.client.Object.Update(ctx, r.workspace_id, plan.Id.ValueString(), &payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating object",
			"Could not update object, unexpected error: "+err.Error(),
		)
		return
	}

	diags = FillInformationsForObject(ctx, &plan, object)
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
func (r *objectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state objectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !r.features.DestroyObject {
		payload := models.ObjectPayloadScheme{
			Attributes: []*models.ObjectPayloadAttributeScheme{
				&models.ObjectPayloadAttributeScheme{
					ObjectTypeAttributeID: r.features.ObsoleteObjectTypeAttributeId,
					ObjectAttributeValues: []*models.ObjectPayloadAttributeValueScheme{
						&models.ObjectPayloadAttributeValueScheme{
							Value: "Obsolete",
						},
					},
				},
			},
			ObjectTypeID: state.ObjectTypeId.ValueString(),
		}

		_, _, err := r.client.Object.Update(ctx, r.workspace_id, state.Id.ValueString(), &payload)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating object",
				"Could not update object, unexpected error: "+err.Error(),
			)
		}
		return
	}

	// Delete existing object
	_, err := r.client.Object.Delete(ctx, r.workspace_id, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting object",
			"Could not delete object, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *objectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func avatarAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"workspace_id": types.StringType,
		"global_id":    types.StringType,
		"id":           types.StringType,
		"avatar_uuid":  types.StringType,
		"url16":        types.StringType,
		"url48":        types.StringType,
		"url72":        types.StringType,
		"url144":       types.StringType,
		"url288":       types.StringType,
		"object_id":    types.StringType,
	}
}

func objectAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"self": types.StringType,
	}
}

func statusAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":       types.StringType,
		"name":     types.StringType,
		"category": types.Int64Type,
	}
}

func groupAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"avatar_url": types.StringType,
		"name":       types.StringType,
	}
}

func objectValueAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"value":            types.StringType,
		"display_value":    types.StringType,
		"search_value":     types.StringType,
		"group":            types.ObjectType{groupAttrTypes()},
		"status":           types.ObjectType{statusAttrTypes()},
		"additional_value": types.StringType,
	}
}

func objectAttributeAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"workspace_id":                types.StringType,
		"global_id":                   types.StringType,
		"id":                          types.StringType,
		"object_type_attribute_id":    types.StringType,
		"object_type_attribute_label": types.BoolType,
		"object_attribute_values":     types.SetType{types.ObjectType{objectValueAttrTypes()}},
	}
}
