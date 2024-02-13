package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

//Plan Modifier for the Avatar Object in the object resource
func SyncAvatarPlanModifier() planmodifier.String {
	return &syncAvatarPlanModifier{}
}

type syncAvatarPlanModifier struct {
}

func (d *syncAvatarPlanModifier) Description(ctx context.Context) string {
	return "Ensures that avatarUUID and avatar attributes are kept synchronized."
}

func (d *syncAvatarPlanModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d *syncAvatarPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	var plan objectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !req.StateValue.IsNull() {
		var state objectResourceModel
		diags = req.State.Get(ctx, &state)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		avatar_uuid_from_state := ""
		avatar_uuid_from_plan := ""

		if !state.Avatar.IsNull() && !state.Avatar.IsUnknown() {
			var avatarState avatarModel
			diags = state.Avatar.As(ctx, &avatarState, basetypes.ObjectAsOptions{})
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			avatar_uuid_from_state = avatarState.AvatarUuid.ValueString()
		}

		if !plan.Avatar.IsNull() && !plan.Avatar.IsUnknown() {
			var avatarPlan avatarModel
			diags = plan.Avatar.As(ctx, &avatarPlan, basetypes.ObjectAsOptions{})
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			avatar_uuid_from_plan = avatarPlan.AvatarUuid.ValueString()
		}

		if avatar_uuid_from_plan != avatar_uuid_from_state {
			resp.PlanValue = types.StringUnknown()
			return
		}
		resp.PlanValue = req.StateValue
		return
	}
	resp.PlanValue = types.StringUnknown()
}

//Plan modifier for the label parameter in the object resource
func SyncLabelPlanModifier() planmodifier.String {
	return &syncLabelPlanModifier{}
}

type syncLabelPlanModifier struct {
}

func (d *syncLabelPlanModifier) Description(ctx context.Context) string {
	return "Ensures that label and attributes are kept synchronized."
}

func (d *syncLabelPlanModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d *syncLabelPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	var plan objectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !req.StateValue.IsNull() {
		var state objectResourceModel
		diags = req.State.Get(ctx, &state)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		if state.Attributes.IsNull() {
			return
		}

		elements := make([]objectAttributeModel, 0, len(state.Attributes.Elements()))
		diags = state.Attributes.ElementsAs(ctx, &elements, false)

		id := "0"
		for _, element := range elements {
			if element.ObjectTypeAttributeLabel.ValueBool() {
				id = element.ObjectTypeAttributeId.ValueString()
			}
		}

		if id == "0" {
			resp.Diagnostics.AddError(
				"Error in object attribute for the label.",
				"Object attribute for the label not found in the object schema.",
			)
			return
		}

		for _, att := range plan.AttributesIn {
			if att.ObjectTypeAttributeId.ValueString() == id {
				if len(att.ObjectAttributeValuesIn) > 1 {
					resp.Diagnostics.AddError(
						"Error in object attribute for the label.",
						"Only one value expected for the label attribute.",
					)
					return
				}
				resp.PlanValue = att.ObjectAttributeValuesIn[0].Value
				return
			}
		}
		resp.Diagnostics.AddError(
			"Error in object attribute for the label.",
			"Object attribute for the label not found.",
		)
		return
	}

	resp.PlanValue = types.StringUnknown()
}
