package provider

import (
	"context"

	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func FillInformationsForDataObject(ctx context.Context, object *objectDataResourceModel, assetsObject *models.ObjectScheme) diag.Diagnostics {
	var diags diag.Diagnostics
	// Map response body to schema and populate Computed attribute values
	object.WorkspaceId = types.StringValue(assetsObject.WorkspaceId)
	object.GlobalId = types.StringValue(assetsObject.GlobalId)
	object.Id = types.StringValue(assetsObject.ID)
	object.Label = types.StringValue(assetsObject.Label)
	object.ObjectKey = types.StringValue(assetsObject.ObjectKey)
	object.ObjectTypeId = types.StringValue(assetsObject.ObjectType.Id)
	object.Created = types.StringValue(assetsObject.Created)
	object.Updated = types.StringValue(assetsObject.Updated)
	object.HasAvatar = types.BoolValue(assetsObject.HasAvatar)

	avatar := types.ObjectNull(avatarAttrTypes())
	if assetsObject.Avatar != nil {
		avatarElements := avatarModel{
			WorkspaceId: types.StringValue(assetsObject.Avatar.WorkspaceId),
			GlobalId:    types.StringValue(assetsObject.Avatar.GlobalId),
			Id:          types.StringValue(assetsObject.Avatar.AvatarUUID),
			AvatarUuid:  types.StringValue(assetsObject.Avatar.AvatarUUID),
			Url16:       types.StringValue(assetsObject.Avatar.Url16),
			Url48:       types.StringValue(assetsObject.Avatar.Url48),
			Url72:       types.StringValue(assetsObject.Avatar.Url72),
			Url144:      types.StringValue(assetsObject.Avatar.Url144),
			Url288:      types.StringValue(assetsObject.Avatar.Url288),
			ObjectId:    types.StringValue(assetsObject.Avatar.ObjectId),
		}
		avatar, diags = types.ObjectValueFrom(ctx, avatarAttrTypes(), avatarElements)
		if diags.HasError() {
			return diags
		}
	}

	object.Avatar = avatar
	linksElements := objectModel{
		Self: types.StringValue(assetsObject.Links.Self),
	}

	object.Links, diags = types.ObjectValueFrom(ctx, objectAttrTypes(), linksElements)
	if diags.HasError() {
		return diags
	}

	var attributes []attr.Value
	for _, att := range assetsObject.Attributes {
		var values []attr.Value
		for _, value := range att.ObjectAttributeValues {
			group := types.ObjectNull(groupAttrTypes())
			if value.Group != nil {
				groupValue := groupModel{
					AvatarUrl: types.StringValue(value.Group.AvatarUrl),
					Name:      types.StringValue(value.Group.Name),
				}
				group, diags = types.ObjectValueFrom(ctx, groupAttrTypes(), groupValue)
				if diags.HasError() {
					return diags
				}
			}
			status := types.ObjectNull(statusAttrTypes())
			if value.Status != nil {
				statusValue := statusModel{
					Id:       types.StringValue(value.Status.ID),
					Name:     types.StringValue(value.Status.Name),
					Category: types.Int64Value(int64(value.Status.Category)),
				}
				status, diags = types.ObjectValueFrom(ctx, statusAttrTypes(), statusValue)
				if diags.HasError() {
					return diags
				}
			}

			valValue := objectAttributeValueModel{
				Value:           types.StringValue(value.Value),
				DisplayValue:    types.StringValue(value.DisplayValue),
				SearchValue:     types.StringValue(value.SearchValue),
				Group:           group,
				Status:          status,
				AdditionalValue: types.StringValue(value.AdditionalValue),
			}
			var val basetypes.ObjectValue
			val, diags = types.ObjectValueFrom(ctx, objectValueAttrTypes(), valValue)
			if diags.HasError() {
				return diags
			}
			values = append(values, val)
		}

		var valueSet basetypes.SetValue
		valueSet, diags = types.SetValueFrom(ctx, types.ObjectType{objectValueAttrTypes()}, values)
		if diags.HasError() {
			return diags
		}

		attributeValue := objectAttributeModel{
			WorkspaceId:              types.StringValue(att.WorkspaceId),
			GlobalId:                 types.StringValue(att.GlobalId),
			Id:                       types.StringValue(att.ID),
			ObjectTypeAttributeLabel: types.BoolValue(att.ObjectTypeAttribute.Label),
			ObjectTypeAttributeId:    types.StringValue(att.ObjectTypeAttributeId),
			ObjectAttributeValues:    valueSet,
		}

		var attribut basetypes.ObjectValue
		attribut, diags = types.ObjectValueFrom(ctx, objectAttributeAttrTypes(), attributeValue)
		if diags.HasError() {
			return diags
		}
		attributes = append(attributes, attribut)
	}

	object.Attributes, diags = types.SetValueFrom(ctx, types.ObjectType{objectAttributeAttrTypes()}, attributes)
	return diags
}

func FillInformationsForObject(ctx context.Context, object *objectResourceModel, assetsObject *models.ObjectScheme) diag.Diagnostics {
	var diags diag.Diagnostics
	// Map response body to schema and populate Computed attribute values
	object.WorkspaceId = types.StringValue(assetsObject.WorkspaceId)
	object.GlobalId = types.StringValue(assetsObject.GlobalId)
	object.Id = types.StringValue(assetsObject.ID)
	object.Label = types.StringValue(assetsObject.Label)
	object.ObjectKey = types.StringValue(assetsObject.ObjectKey)
	object.ObjectTypeId = types.StringValue(assetsObject.ObjectType.Id)
	object.Created = types.StringValue(assetsObject.Created)
	object.Updated = types.StringValue(assetsObject.Updated)
	object.HasAvatar = types.BoolValue(assetsObject.HasAvatar)

	avatar := types.ObjectNull(avatarAttrTypes())
	if assetsObject.Avatar != nil {
		avatarElements := avatarModel{
			WorkspaceId: types.StringValue(assetsObject.Avatar.WorkspaceId),
			GlobalId:    types.StringValue(assetsObject.Avatar.GlobalId),
			Id:          types.StringValue(assetsObject.Avatar.AvatarUUID),
			AvatarUuid:  types.StringValue(assetsObject.Avatar.AvatarUUID),
			Url16:       types.StringValue(assetsObject.Avatar.Url16),
			Url48:       types.StringValue(assetsObject.Avatar.Url48),
			Url72:       types.StringValue(assetsObject.Avatar.Url72),
			Url144:      types.StringValue(assetsObject.Avatar.Url144),
			Url288:      types.StringValue(assetsObject.Avatar.Url288),
			ObjectId:    types.StringValue(assetsObject.Avatar.ObjectId),
		}
		avatar, diags = types.ObjectValueFrom(ctx, avatarAttrTypes(), avatarElements)
		if diags.HasError() {
			return diags
		}
	}

	object.Avatar = avatar
	linksElements := objectModel{
		Self: types.StringValue(assetsObject.Links.Self),
	}

	object.Links, diags = types.ObjectValueFrom(ctx, objectAttrTypes(), linksElements)
	if diags.HasError() {
		return diags
	}

	var attributes []attr.Value
	for _, att := range assetsObject.Attributes {
		var values []attr.Value
		for _, value := range att.ObjectAttributeValues {
			group := types.ObjectNull(groupAttrTypes())
			if value.Group != nil {
				groupValue := groupModel{
					AvatarUrl: types.StringValue(value.Group.AvatarUrl),
					Name:      types.StringValue(value.Group.Name),
				}
				group, diags = types.ObjectValueFrom(ctx, groupAttrTypes(), groupValue)
				if diags.HasError() {
					return diags
				}
			}
			status := types.ObjectNull(statusAttrTypes())
			if value.Status != nil {
				statusValue := statusModel{
					Id:       types.StringValue(value.Status.ID),
					Name:     types.StringValue(value.Status.Name),
					Category: types.Int64Value(int64(value.Status.Category)),
				}
				status, diags = types.ObjectValueFrom(ctx, statusAttrTypes(), statusValue)
				if diags.HasError() {
					return diags
				}
			}

			valValue := objectAttributeValueModel{
				Value:           types.StringValue(value.Value),
				DisplayValue:    types.StringValue(value.DisplayValue),
				SearchValue:     types.StringValue(value.SearchValue),
				Group:           group,
				Status:          status,
				AdditionalValue: types.StringValue(value.AdditionalValue),
			}
			var val basetypes.ObjectValue
			val, diags = types.ObjectValueFrom(ctx, objectValueAttrTypes(), valValue)
			if diags.HasError() {
				return diags
			}
			values = append(values, val)
		}

		var valueSet basetypes.SetValue
		valueSet, diags = types.SetValueFrom(ctx, types.ObjectType{objectValueAttrTypes()}, values)
		if diags.HasError() {
			return diags
		}

		attributeValue := objectAttributeModel{
			WorkspaceId:              types.StringValue(att.WorkspaceId),
			GlobalId:                 types.StringValue(att.GlobalId),
			Id:                       types.StringValue(att.ID),
			ObjectTypeAttributeLabel: types.BoolValue(att.ObjectTypeAttribute.Label),
			ObjectTypeAttributeId:    types.StringValue(att.ObjectTypeAttributeId),
			ObjectAttributeValues:    valueSet,
		}

		var attribut basetypes.ObjectValue
		attribut, diags = types.ObjectValueFrom(ctx, objectAttributeAttrTypes(), attributeValue)
		if diags.HasError() {
			return diags
		}
		attributes = append(attributes, attribut)
	}

	object.Attributes, diags = types.SetValueFrom(ctx, types.ObjectType{objectAttributeAttrTypes()}, attributes)
	return diags
}

func FillInformationForIcon(icon *IconModel, assetsIcon *models.IconScheme) {
	icon.Id = types.StringValue(assetsIcon.ID)
	icon.Name = types.StringValue(assetsIcon.Name)
	icon.Url16 = types.StringValue(assetsIcon.URL16)
	icon.Url48 = types.StringValue(assetsIcon.URL48)
}

func FillInformationsForObjectType(objectType *objectTypeResourceModel, assetsObjectType *models.ObjectTypeScheme) {
	objectType.WorkspaceId = types.StringValue(assetsObjectType.WorkspaceId)
	objectType.GlobalId = types.StringValue(assetsObjectType.GlobalId)
	objectType.Id = types.StringValue(assetsObjectType.Id)
	objectType.Name = types.StringValue(assetsObjectType.Name)
	objectType.Description = types.StringValue(assetsObjectType.Description)
	objectType.IconId = types.StringValue(assetsObjectType.Icon.ID)
	objectType.Position = types.Int64Value(int64(assetsObjectType.Position))
	objectType.Created = types.StringValue(assetsObjectType.Created)
	objectType.Updated = types.StringValue(assetsObjectType.Updated)
	objectType.ObjectCount = types.Int64Value(int64(assetsObjectType.ObjectCount))
	objectType.ParentObjectTypeId = types.StringValue(assetsObjectType.ParentObjectTypeId)
	objectType.ObjectSchemaId = types.StringValue(assetsObjectType.ObjectSchemaId)
	objectType.Inherited = types.BoolValue(assetsObjectType.Inherited)
	objectType.AbstractObjectType = types.BoolValue(assetsObjectType.AbstractObjectType)
	objectType.ParentObjectTypeInherited = types.BoolValue(assetsObjectType.ParentObjectTypeInherited)
}

func FillInformationsForObjectTypeAttribute(ctx context.Context, objectTypeAttribute *objectTypeAttributeResourceModel, assetsObjectTypeAttribute *models.ObjectTypeAttributeScheme) diag.Diagnostics {
	var diags diag.Diagnostics

	objectTypeAttribute.WorkspaceId = types.StringValue(assetsObjectTypeAttribute.WorkspaceId)
	objectTypeAttribute.GlobalId = types.StringValue(assetsObjectTypeAttribute.GlobalId)
	objectTypeAttribute.Id = types.StringValue(assetsObjectTypeAttribute.ID)
	objectTypeAttribute.Name = types.StringValue(assetsObjectTypeAttribute.Name)
	objectTypeAttribute.Label = types.BoolValue(assetsObjectTypeAttribute.Label)
	objectTypeAttribute.Type = types.Int64Value(int64(assetsObjectTypeAttribute.Type))
	objectTypeAttribute.Description = types.StringValue(assetsObjectTypeAttribute.Description)

	defaultType := types.ObjectNull(defaultTypeAttrTypes())
	if assetsObjectTypeAttribute.DefaultType != nil {
		defaultTypeElement := defaultTypeModel{
			Id:   types.Int64Value(int64(assetsObjectTypeAttribute.DefaultType.ID)),
			Name: types.StringValue(assetsObjectTypeAttribute.DefaultType.Name),
		}
		defaultType, diags = types.ObjectValueFrom(ctx, defaultTypeAttrTypes(), defaultTypeElement)
		if diags.HasError() {
			return diags
		}
	}
	objectTypeAttribute.DefaultType = defaultType

	objectTypeAttribute.TypeValue = types.StringValue(assetsObjectTypeAttribute.TypeValue)
	objectTypeAttribute.TypeValueMulti, diags = types.ListValueFrom(ctx, types.StringType, assetsObjectTypeAttribute.TypeValueMulti)
	if diags.HasError() {
		return diags
	}
	objectTypeAttribute.AdditionalValue = types.StringValue(assetsObjectTypeAttribute.AdditionalValue)

	referenceType := types.ObjectNull(referenceTypeAttrTypes())
	if assetsObjectTypeAttribute.ReferenceType != nil {
		referenceTypeElement := referenceTypeModel{
			WorkspaceId: types.StringValue(assetsObjectTypeAttribute.ReferenceType.WorkspaceId),
			GlobalId:    types.StringValue(assetsObjectTypeAttribute.ReferenceType.GlobalId),
			Name:        types.StringValue(assetsObjectTypeAttribute.ReferenceType.Name),
		}
		referenceType, diags = types.ObjectValueFrom(ctx, referenceTypeAttrTypes(), referenceTypeElement)
		if diags.HasError() {
			return diags
		}
	}

	objectTypeAttribute.ReferenceType = referenceType

	objectTypeAttribute.ReferenceObjectTypeId = types.StringValue(assetsObjectTypeAttribute.ReferenceObjectTypeId)
	objectTypeAttribute.Editable = types.BoolValue(assetsObjectTypeAttribute.Editable)
	objectTypeAttribute.System = types.BoolValue(assetsObjectTypeAttribute.System)
	objectTypeAttribute.Indexed = types.BoolValue(assetsObjectTypeAttribute.Indexed)
	objectTypeAttribute.Sortable = types.BoolValue(assetsObjectTypeAttribute.Sortable)
	objectTypeAttribute.Summable = types.BoolValue(assetsObjectTypeAttribute.Summable)
	objectTypeAttribute.MinimumCardinality = types.Int64Value(int64(assetsObjectTypeAttribute.MinimumCardinality))
	objectTypeAttribute.MaximumCardinality = types.Int64Value(int64(assetsObjectTypeAttribute.MaximumCardinality))
	objectTypeAttribute.Suffix = types.StringValue(assetsObjectTypeAttribute.Suffix)
	objectTypeAttribute.Removable = types.BoolValue(assetsObjectTypeAttribute.Removable)
	objectTypeAttribute.ObjectAttributeExists = types.BoolValue(assetsObjectTypeAttribute.ObjectAttributeExists)
	objectTypeAttribute.Hidden = types.BoolValue(assetsObjectTypeAttribute.Hidden)
	objectTypeAttribute.IncludeChildObjectTypes = types.BoolValue(assetsObjectTypeAttribute.IncludeChildObjectTypes)
	objectTypeAttribute.UniqueAttribute = types.BoolValue(assetsObjectTypeAttribute.UniqueAttribute)
	objectTypeAttribute.RegexValidation = types.StringValue(assetsObjectTypeAttribute.RegexValidation)
	objectTypeAttribute.QlQuery = types.StringValue(assetsObjectTypeAttribute.QlQuery)
	objectTypeAttribute.Options = types.StringValue(assetsObjectTypeAttribute.Options)
	objectTypeAttribute.Position = types.Int64Value(int64(assetsObjectTypeAttribute.Position))

	return diags
}

func FillInformationsForDataObjectTypeAttribute(ctx context.Context, objectTypeAttribute *objectTypeAttributeDataModel, assetsObjectTypeAttribute *models.ObjectTypeAttributeScheme) diag.Diagnostics {
	var diags diag.Diagnostics

	objectTypeAttribute.WorkspaceId = types.StringValue(assetsObjectTypeAttribute.WorkspaceId)
	objectTypeAttribute.GlobalId = types.StringValue(assetsObjectTypeAttribute.GlobalId)
	objectTypeAttribute.Id = types.StringValue(assetsObjectTypeAttribute.ID)
	objectTypeAttribute.Name = types.StringValue(assetsObjectTypeAttribute.Name)
	objectTypeAttribute.Label = types.BoolValue(assetsObjectTypeAttribute.Label)
	objectTypeAttribute.Type = types.Int64Value(int64(assetsObjectTypeAttribute.Type))
	objectTypeAttribute.Description = types.StringValue(assetsObjectTypeAttribute.Description)

	defaultType := types.ObjectNull(defaultTypeAttrTypes())
	if assetsObjectTypeAttribute.DefaultType != nil {
		defaultTypeElement := defaultTypeModel{
			Id:   types.Int64Value(int64(assetsObjectTypeAttribute.DefaultType.ID)),
			Name: types.StringValue(assetsObjectTypeAttribute.DefaultType.Name),
		}
		defaultType, diags = types.ObjectValueFrom(ctx, defaultTypeAttrTypes(), defaultTypeElement)
		if diags.HasError() {
			return diags
		}
	}
	objectTypeAttribute.DefaultType = defaultType

	objectTypeAttribute.TypeValue = types.StringValue(assetsObjectTypeAttribute.TypeValue)
	objectTypeAttribute.TypeValueMulti, diags = types.ListValueFrom(ctx, types.StringType, assetsObjectTypeAttribute.TypeValueMulti)
	if diags.HasError() {
		return diags
	}
	objectTypeAttribute.AdditionalValue = types.StringValue(assetsObjectTypeAttribute.AdditionalValue)

	referenceType := types.ObjectNull(referenceTypeAttrTypes())
	if assetsObjectTypeAttribute.ReferenceType != nil {
		referenceTypeElement := referenceTypeModel{
			WorkspaceId: types.StringValue(assetsObjectTypeAttribute.ReferenceType.WorkspaceId),
			GlobalId:    types.StringValue(assetsObjectTypeAttribute.ReferenceType.GlobalId),
			Name:        types.StringValue(assetsObjectTypeAttribute.ReferenceType.Name),
		}
		referenceType, diags = types.ObjectValueFrom(ctx, referenceTypeAttrTypes(), referenceTypeElement)
		if diags.HasError() {
			return diags
		}
	}

	objectTypeAttribute.ReferenceType = referenceType

	objectTypeAttribute.ReferenceObjectTypeId = types.StringValue(assetsObjectTypeAttribute.ReferenceObjectTypeId)
	objectTypeAttribute.Editable = types.BoolValue(assetsObjectTypeAttribute.Editable)
	objectTypeAttribute.System = types.BoolValue(assetsObjectTypeAttribute.System)
	objectTypeAttribute.Indexed = types.BoolValue(assetsObjectTypeAttribute.Indexed)
	objectTypeAttribute.Sortable = types.BoolValue(assetsObjectTypeAttribute.Sortable)
	objectTypeAttribute.Summable = types.BoolValue(assetsObjectTypeAttribute.Summable)
	objectTypeAttribute.MinimumCardinality = types.Int64Value(int64(assetsObjectTypeAttribute.MinimumCardinality))
	objectTypeAttribute.MaximumCardinality = types.Int64Value(int64(assetsObjectTypeAttribute.MaximumCardinality))
	objectTypeAttribute.Suffix = types.StringValue(assetsObjectTypeAttribute.Suffix)
	objectTypeAttribute.Removable = types.BoolValue(assetsObjectTypeAttribute.Removable)
	objectTypeAttribute.ObjectAttributeExists = types.BoolValue(assetsObjectTypeAttribute.ObjectAttributeExists)
	objectTypeAttribute.Hidden = types.BoolValue(assetsObjectTypeAttribute.Hidden)
	objectTypeAttribute.IncludeChildObjectTypes = types.BoolValue(assetsObjectTypeAttribute.IncludeChildObjectTypes)
	objectTypeAttribute.UniqueAttribute = types.BoolValue(assetsObjectTypeAttribute.UniqueAttribute)
	objectTypeAttribute.RegexValidation = types.StringValue(assetsObjectTypeAttribute.RegexValidation)
	objectTypeAttribute.QlQuery = types.StringValue(assetsObjectTypeAttribute.QlQuery)
	objectTypeAttribute.Options = types.StringValue(assetsObjectTypeAttribute.Options)
	objectTypeAttribute.Position = types.Int64Value(int64(assetsObjectTypeAttribute.Position))

	return diags
}

func FillInformationsForObjectSchema(objectSchema *objectSchemaResourceModel, assetsObjectSchema *models.ObjectSchemaScheme) {
	objectSchema.WorkspaceId = types.StringValue(assetsObjectSchema.WorkspaceId)
	objectSchema.GlobalId = types.StringValue(assetsObjectSchema.GlobalId)
	objectSchema.Id = types.StringValue(assetsObjectSchema.Id)
	objectSchema.Name = types.StringValue(assetsObjectSchema.Name)
	objectSchema.ObjectSchemaKey = types.StringValue(assetsObjectSchema.ObjectSchemaKey)
	objectSchema.Description = types.StringValue(assetsObjectSchema.Description)
	objectSchema.Status = types.StringValue(assetsObjectSchema.Status)
	objectSchema.Created = types.StringValue(assetsObjectSchema.Created)
	objectSchema.Updated = types.StringValue(assetsObjectSchema.Updated)
	objectSchema.ObjectCount = types.Int64Value(int64(assetsObjectSchema.ObjectCount))
	objectSchema.ObjectTypeCount = types.Int64Value(int64(assetsObjectSchema.ObjectTypeCount))
	objectSchema.CanManage = types.BoolValue(assetsObjectSchema.CanManage)
}
