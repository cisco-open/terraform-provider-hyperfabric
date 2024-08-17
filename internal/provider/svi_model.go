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

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type SviResourceModel struct {
	Enabled       types.Bool `tfsdk:"enabled"`
	Ipv4Addresses types.Set  `tfsdk:"ipv4_addresses"`
	Ipv6Addresses types.Set  `tfsdk:"ipv6_addresses"`
	// VlanId        types.Int32 `tfsdk:"vlan_id"`
}

func SviResourceModelAttributeType() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled":        types.BoolType,
		"ipv4_addresses": types.SetType{}.WithElementType(types.StringType),
		"ipv6_addresses": types.SetType{}.WithElementType(types.StringType),
		// "vlanId":        types.StringType,
	}
}

func getEmptySviResourceModel() SviResourceModel {
	return SviResourceModel{
		Enabled:       basetypes.NewBoolNull(),
		Ipv4Addresses: basetypes.NewSetNull(types.StringType),
		Ipv6Addresses: basetypes.NewSetNull(types.StringType),
		// VlanId:        basetypes.NewInt32Null(),
	}
}

func getSviSchemaAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: `The SVI / Distributed GW object for a specific VNI.`,
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.Object{
			// Validate this attribute must be configured with other_attr.
			objectvalidator.AlsoRequires(path.Expressions{
				path.MatchRoot("vrf_id"),
			}...),
		},
		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: `"The enabled state of the SVI.`,
			},
			"ipv4_addresses": schema.SetAttribute{
				MarkdownDescription: `A set of IPv4 addresses for the SVI.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"ipv6_addresses": schema.SetAttribute{
				MarkdownDescription: `A set of IPv6 addresses for the SVI.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			// "vlan_id": schema.Int32Attribute{
			//	Optional: true,
			// 	Computed: true,
			// 	PlanModifiers: []planmodifier.Int32{
			// 		int32planmodifier.UseStateForUnknown(),
			// 	},
			// 	MarkdownDescription: `The VLAN ID used as encapsulation for this SVI.`,
			// }
		},
	}
}

func getSviDataSourceSchemaAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: `A set of annotations to store user-defined data.`,
		Computed:            true,
		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: `"The enabled state of the SVI.`,
			},
			"ipv4_addresses": schema.SetAttribute{
				MarkdownDescription: `A set of IPv4 addresses for the SVI.`,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"ipv6_addresses": schema.SetAttribute{
				MarkdownDescription: `A set of IPv6 addresses for the SVI.`,
				Computed:            true,
				ElementType:         types.StringType,
			},
			// "vlan_id": schema.Int32Attribute{
			// 	Computed: true,
			// 	MarkdownDescription: `The VLAN ID used as encapsulation for this SVI.`,
			// }
		},
	}
}

func NewSviResourceModel(ctx context.Context, data map[string]interface{}) SviResourceModel {
	svi := getEmptySviResourceModel()
	for attributeName, attributeValue := range data {
		if attributeName == "enabled" && attributeValue != nil {
			boolAttr := attributeValue.(bool)
			svi.Enabled = basetypes.NewBoolValue(boolAttr)
		} else if attributeName == "ipv4Addresses" && attributeValue != nil {
			stringListAttr := attributeValue.([]interface{})
			if len(stringListAttr) > 0 {
				ipv4AddressesSet, _ := types.SetValueFrom(ctx, types.StringType, stringListAttr)
				svi.Ipv4Addresses = ipv4AddressesSet
			}
		} else if attributeName == "ipv6Addresses" && attributeValue != nil {
			stringListAttr := attributeValue.([]interface{})
			if len(stringListAttr) > 0 {
				ipv6AddressesSet, _ := types.SetValueFrom(ctx, types.StringType, stringListAttr)
				svi.Ipv6Addresses = ipv6AddressesSet
			}
			// } else if attributeName == "vlanId" && attributeValue != nil {
			// 	int32Attr := attributeValue.(int32)
			// 	if int32Attr != 0 {
			// 		svi.VlanId = basetypes.NewInt32Value(int32Attr)
			// 	}
		}
	}
	return svi
}

func NewSviObject(ctx context.Context, data []interface{}) basetypes.ObjectValue {
	var sviObject basetypes.ObjectValue
	if len(data) > 0 {
		svi := NewSviResourceModel(ctx, data[0].(map[string]interface{}))
		sviObject, _ = types.ObjectValueFrom(ctx, SviResourceModelAttributeType(), svi)
	} else {
		sviObject = basetypes.NewObjectNull(SviResourceModelAttributeType())
	}
	return sviObject
}

func getSviJsonPayload(ctx context.Context, data basetypes.ObjectValue) []map[string]interface{} {
	svi := SviResourceModel{}
	data.As(ctx, &svi, basetypes.ObjectAsOptions{})
	ipv4Addresses := make([]string, 0)
	ipv6Addresses := make([]string, 0)
	svi.Ipv4Addresses.ElementsAs(ctx, &ipv4Addresses, false)
	svi.Ipv6Addresses.ElementsAs(ctx, &ipv6Addresses, false)
	sviPayload := map[string]interface{}{
		"enabled":       svi.Enabled.ValueBool(),
		"ipv4Addresses": ipv4Addresses,
		"ipv6Addresses": ipv6Addresses,
		// "vlanId":        svi.VlanId,
	}
	svisPayloads := make([]map[string]interface{}, 0)
	svisPayloads = append(svisPayloads, sviPayload)
	return svisPayloads
}
