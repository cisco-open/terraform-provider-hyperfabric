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

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type AnnotationResourceModel struct {
	DataType types.String `tfsdk:"data_type"`
	Name     types.String `tfsdk:"name"`
	Value    types.String `tfsdk:"value"`
}

func AnnotationResourceModelAttributeType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"data_type": types.StringType,
			"name":      types.StringType,
			"value":     types.StringType,
		},
	}
}

func getEmptyAnnotationResourceModel() AnnotationResourceModel {
	return AnnotationResourceModel{
		DataType: basetypes.NewStringValue("STRING"),
		Name:     basetypes.NewStringNull(),
		Value:    basetypes.NewStringNull(),
	}
}

func getAnnotationsSchemaAttribute() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: `A set of annotations to store user-defined data.`,
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.UseStateForUnknown(),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"data_type": schema.StringAttribute{
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
					Default: stringdefault.StaticString("STRING"),
					Validators: []validator.String{
						stringvalidator.OneOf([]string{"STRING", "INT32", "UINT32", "INT64", "UINT64", "BOOL", "TIME", "UUID", "DURATION", "JSON"}...),
					},
					MarkdownDescription: `The type of data stored in the value of the annotation.`,
				},
				"name": schema.StringAttribute{
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
					Validators: []validator.String{
						MakeStringRequired(),
					},
					MarkdownDescription: `The name used to uniquely identify the annotation.`,
				},
				"value": schema.StringAttribute{
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
					Validators: []validator.String{
						MakeStringRequired(),
					},
					MarkdownDescription: `The value of the annotation.`,
				},
			},
		},
	}
}

func getAnnotationsDataSourceSchemaAttribute() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: `A set of annotations to store user-defined data.`,
		Computed:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"data_type": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: `The type of data stored in the value of the annotation.`,
				},
				"name": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: `The name used to uniquely identify the annotation.`,
				},
				"value": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: `The value of the annotation.`,
				},
			},
		},
	}
}

func NewAnnotationResourceModel(data map[string]interface{}) AnnotationResourceModel {
	annotation := getEmptyAnnotationResourceModel()
	for attributeName, attributeValue := range data {
		if attributeName == "dataType" && attributeValue != nil {
			stringAttr := attributeValue.(string)
			if stringAttr != "" {
				annotation.DataType = basetypes.NewStringValue(stringAttr)
			}
		} else if attributeName == "name" && attributeValue != nil {
			stringAttr := attributeValue.(string)
			if stringAttr != "" {
				annotation.Name = basetypes.NewStringValue(stringAttr)
			}
		} else if attributeName == "value" && attributeValue != nil {
			stringAttr := attributeValue.(string)
			if stringAttr != "" {
				annotation.Value = basetypes.NewStringValue(stringAttr)
			}
		}
	}
	return annotation
}

func NewAnnotationsSet(ctx context.Context, data []interface{}) basetypes.SetValue {
	annotations := make([]AnnotationResourceModel, 0)
	for _, annotation := range data {
		newAnnotation := NewAnnotationResourceModel(annotation.(map[string]interface{}))
		annotations = append(annotations, newAnnotation)
	}
	annotationsSet, _ := types.SetValueFrom(ctx, AnnotationResourceModelAttributeType(), annotations)
	return annotationsSet
}

func NewNodeAnnotationsSet(ctx context.Context, data []interface{}) basetypes.SetValue {
	annotations := make([]AnnotationResourceModel, 0)
	for _, annotation := range data {
		newAnnotation := NewAnnotationResourceModel(annotation.(map[string]interface{}))
		if newAnnotation.Name.ValueString() != "position" {
			annotations = append(annotations, newAnnotation)
		}
	}
	annotationsSet, _ := types.SetValueFrom(ctx, AnnotationResourceModelAttributeType(), annotations)
	return annotationsSet
}

func getAnnotationsJsonPayload(ctx context.Context, data basetypes.SetValue) []map[string]string {
	annotations := []AnnotationResourceModel{}
	data.ElementsAs(ctx, &annotations, false)
	annotationPayloads := []map[string]string{}
	for _, annotation := range annotations {
		annotationPayloads = append(annotationPayloads, map[string]string{
			"dataType": StripQuotes(annotation.DataType.String()),
			"name":     StripQuotes(annotation.Name.String()),
			"value":    StripQuotes(annotation.Value.String()),
		})
	}
	return annotationPayloads
}
