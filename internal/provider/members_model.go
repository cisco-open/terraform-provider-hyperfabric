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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type MemberResourceModel struct {
	PortName types.String  `tfsdk:"port_name"`
	NodeId   types.String  `tfsdk:"node_id"`
	NodeName types.String  `tfsdk:"node_name"`
	VlanId   types.Float64 `tfsdk:"vlan_id"`
	// Untagged types.Bool	`tfsdk:"untagged"`
}

// {
// 	"vlanId": 2,
// 	"port": {
// 		"portName": "*",
// 		"chassisId": "603ce8f2-2e10-409b-9ffe-f19378d46423",
// 		"chassisName": "fab1-leaf1"
// 	},
// 	"untagged": false
// }

func MemberResourceModelAttributeType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"port_name": types.StringType,
			"node_id":   types.StringType,
			"node_name": types.StringType,
			"vlan_id":   types.Float64Type,
			// "untagged":  types.BoolType,
		},
	}
}

func getEmptyMemberResourceModel() MemberResourceModel {
	return MemberResourceModel{
		PortName: basetypes.NewStringNull(),
		NodeId:   basetypes.NewStringNull(),
		NodeName: basetypes.NewStringNull(),
		VlanId:   basetypes.NewFloat64Null(),
		// Untagged: basetypes.NewBoolNull(),
	}
}

func getMembersSchemaAttribute() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: `A set of member ports assigning a specific VLAN to the VNI.`,
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.UseStateForUnknown(),
			// SetToSetNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"port_name": schema.StringAttribute{
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
					// Default:             stringdefault.StaticString("*"),
					MarkdownDescription: `The name of the port or "*" for all ports on a node or all nodes.`,
				},
				"node_id": schema.StringAttribute{
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
					// Default:             stringdefault.StaticString("*"),
					MarkdownDescription: `The unique Id of a node in the Fabric.`,
				},
				"node_name": schema.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
						SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
					},
					MarkdownDescription: `The name of a node in the Fabric.`,
				},
				"vlan_id": schema.Float64Attribute{
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.Float64{
						float64planmodifier.UseStateForUnknown(),
					},
					MarkdownDescription: `The VLAN ID used as encapsulation for the traffic on this port for this VNI.`,
				},
				// "untagged": schema.BoolAttribute{
				// 	Optional: true,
				// 	Computed: true,
				// 	PlanModifiers: []planmodifier.Bool{
				// 		boolplanmodifier.UseStateForUnknown(),
				// 	},
				// 	// Validators: []validator.String{
				// 	// 	MakeStringRequired(),
				// 	// },
				// 	MarkdownDescription: `The untagged state for the traffic on this port for this VNI.`,
				// },
			},
		},
	}
}

func getMembersDataSourceSchemaAttribute() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: `A set of member ports assigning a specific VLAN to the VNI.`,
		Computed:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"port_name": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: `The name of the port or "*" for all ports on a node or all nodes.`,
				},
				"node_id": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: `The unique Id of a node in the Fabric.`,
				},
				"node_name": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: `The name of a node in the Fabric.`,
				},
				"vlan_id": schema.Float64Attribute{
					Computed:            true,
					MarkdownDescription: `The VLAN ID used as encapsulation for the traffic on this port for this VNI.`,
				},
				// "untagged": schema.BoolAttribute{
				// 	Computed:            true,
				// 	MarkdownDescription: `The untagged state for the traffic on this port for this VNI.`,
				// },
			},
		},
	}
}

func NewMemberResourceModel(data map[string]interface{}) MemberResourceModel {
	member := getEmptyMemberResourceModel()
	isUntagged := false
	for attributeName, attributeValue := range data {
		if attributeName == "port" && attributeValue != nil {
			for portAttributeName, portAttributeValue := range attributeValue.(map[string]interface{}) {
				if portAttributeName == "portName" && portAttributeValue != nil {
					stringAttr := portAttributeValue.(string)
					if stringAttr != "" {
						member.PortName = basetypes.NewStringValue(stringAttr)
					}
				} else if portAttributeName == "nodeId" && portAttributeValue != nil {
					stringAttr := portAttributeValue.(string)
					if stringAttr != "" {
						member.NodeId = basetypes.NewStringValue(stringAttr)
					}
				} else if portAttributeName == "nodeName" && portAttributeValue != nil {
					stringAttr := portAttributeValue.(string)
					if stringAttr != "" {
						member.NodeName = basetypes.NewStringValue(stringAttr)
					}
				}
			}
		} else if attributeName == "vlanId" && attributeValue != nil {
			float64Attr := attributeValue.(float64)
			member.VlanId = basetypes.NewFloat64Value(float64Attr)
		} else if attributeName == "untagged" && attributeValue != nil {
			isUntagged = attributeValue.(bool)
		}
	}
	if isUntagged {
		member.VlanId = basetypes.NewFloat64Null()
	}
	return member
}

func NewMembersSet(ctx context.Context, data []interface{}) basetypes.SetValue {
	members := make([]MemberResourceModel, 0)
	for _, member := range data {
		newMember := NewMemberResourceModel(member.(map[string]interface{}))
		members = append(members, newMember)
	}
	membersSet, _ := types.SetValueFrom(ctx, MemberResourceModelAttributeType(), members)
	return membersSet
}

func getMembersJsonPayload(ctx context.Context, data basetypes.SetValue) []map[string]interface{} {
	members := []MemberResourceModel{}
	data.ElementsAs(ctx, &members, false)
	memberPayloads := make([]map[string]interface{}, 0)
	for _, member := range members {
		memberPayload := map[string]interface{}{
			"port": map[string]string{
				"portName": StripQuotes(member.PortName.String()),
				"nodeId":   StripQuotes(member.NodeId.String()),
			},
		}
		if !member.VlanId.IsNull() && !member.VlanId.IsUnknown() {
			memberPayload["vlanId"] = member.VlanId.ValueFloat64()
		} else {
			memberPayload["untagged"] = true
		}
		memberPayloads = append(memberPayloads, memberPayload)
	}
	return memberPayloads
}
