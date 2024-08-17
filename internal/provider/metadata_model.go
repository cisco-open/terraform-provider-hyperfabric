// Copyright 2024 Cisco Systems, Inc. and its affiliates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type MetadataResourceModel struct {
	CreatedAt  types.String `tfsdk:"created_at"`
	CreatedBy  types.String `tfsdk:"created_by"`
	ModifiedAt types.String `tfsdk:"modified_at"`
	ModifiedBy types.String `tfsdk:"modified_by"`
	RevisionId types.String `tfsdk:"revision_id"`
}

func MetadataResourceModelAttributeType() map[string]attr.Type {
	return map[string]attr.Type{
		"created_at":  types.StringType,
		"created_by":  types.StringType,
		"modified_at": types.StringType,
		"modified_by": types.StringType,
		"revision_id": types.StringType,
	}
}

func getEmptyMetadataResourceModel() MetadataResourceModel {
	return MetadataResourceModel{
		CreatedAt:  basetypes.NewStringNull(),
		CreatedBy:  basetypes.NewStringNull(),
		ModifiedAt: basetypes.NewStringNull(),
		ModifiedBy: basetypes.NewStringNull(),
		RevisionId: basetypes.NewStringNull(),
	}
}

func getMetadataSchemaAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: `The metadata information for an object.`,
		Computed:            true,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.UseStateForUnknown(),
			SetToObjectNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
		},
		Attributes: map[string]schema.Attribute{
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: `The timestamp when this object was created in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.`,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"created_by": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: `The user that created this object.`,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"modified_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: `The timestamp when this object was last modified in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.`,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"modified_by": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: `The user that modified this object last.`,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"revision_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: `An integer that represent the current revision of the object.`,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
		},
	}
}

func NewMetadataResourceModel(data map[string]interface{}) MetadataResourceModel {
	metadata := getEmptyMetadataResourceModel()
	for attributeName, attributeValue := range data {
		if attributeName == "createdAt" && attributeValue != nil {
			stringAttr := attributeValue.(string)
			if stringAttr != "" {
				metadata.CreatedAt = basetypes.NewStringValue(stringAttr)
			}
		} else if attributeName == "createdBy" && attributeValue != nil {
			stringAttr := attributeValue.(string)
			if stringAttr != "" {
				metadata.CreatedBy = basetypes.NewStringValue(stringAttr)
			}
		} else if attributeName == "modifiedAt" && attributeValue != nil {
			stringAttr := attributeValue.(string)
			if stringAttr != "" {
				metadata.ModifiedAt = basetypes.NewStringValue(stringAttr)
			}
		} else if attributeName == "modifiedBy" && attributeValue != nil {
			stringAttr := attributeValue.(string)
			if stringAttr != "" {
				metadata.ModifiedBy = basetypes.NewStringValue(stringAttr)
			}
		} else if attributeName == "revisionId" && attributeValue != nil {
			stringAttr := attributeValue.(string)
			if stringAttr != "" {
				metadata.RevisionId = basetypes.NewStringValue(stringAttr)
			}
		}
	}
	return metadata
}

func NewMetadataObject(ctx context.Context, data map[string]interface{}) basetypes.ObjectValue {
	metadata := NewMetadataResourceModel(data)
	metadataObject, _ := types.ObjectValueFrom(ctx, MetadataResourceModelAttributeType(), metadata)
	return metadataObject
}
