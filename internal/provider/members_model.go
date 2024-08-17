// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

func getEmptyMemberResourceModel() *MemberResourceModel {
	return &MemberResourceModel{
		PortName: basetypes.NewStringNull(),
		NodeId:   basetypes.NewStringNull(),
		NodeName: basetypes.NewStringNull(),
		VlanId:   basetypes.NewFloat64Null(),
		// Untagged: basetypes.NewBoolNull(),
	}
}

func getNewMemberResourceModelFromData(data *MemberResourceModel) *MemberResourceModel {
	newMember := getEmptyMemberResourceModel()

	if !data.PortName.IsNull() && !data.PortName.IsUnknown() {
		newMember.PortName = data.PortName
	}

	if !data.NodeId.IsNull() && !data.NodeId.IsUnknown() {
		newMember.NodeId = data.NodeId
	}

	if !data.NodeName.IsNull() && !data.NodeName.IsUnknown() {
		newMember.NodeName = data.NodeName
	}

	if !data.VlanId.IsNull() && !data.VlanId.IsUnknown() {
		newMember.VlanId = data.VlanId
	}

	// if !data.Untagged.IsNull() && !data.Untagged.IsUnknown() {
	// 	newMember.Untagged = data.Untagged
	// }

	return newMember
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
						SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
					},
					// Default:             stringdefault.StaticString("*"),
					MarkdownDescription: `The name of the port or "*" for all ports on a node or all nodes.`,
				},
				"node_id": schema.StringAttribute{
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
						SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
					},
					// Default:             stringdefault.StaticString("*"),
					MarkdownDescription: `The unique Id of a node in the Fabric.`,
				},
				"node_name": schema.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
						// SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
					},
					MarkdownDescription: `The name of a node in the Fabric.`,
				},
				"vlan_id": schema.Float64Attribute{
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.Float64{
						float64planmodifier.UseStateForUnknown(),
						SetToFloat64NullWhenStateIsNullPlanIsUnknownDuringUpdate(),
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

func NewMemberResourceModel(ctx context.Context, data *MemberResourceModel, attributes map[string]interface{}) MemberResourceModel {
	tflog.Debug(ctx, fmt.Sprintf("MEMBER LHLOG '%v' & '%v'", data, attributes))
	member := *getNewMemberResourceModelFromData(data)
	isUntagged := false
	for attributeName, attributeValue := range attributes {
		if attributeName == "port" && attributeValue != nil {
			for portAttributeName, portAttributeValue := range attributeValue.(map[string]interface{}) {
				if portAttributeName == "portName" {
					member.PortName = basetypes.NewStringValue(portAttributeValue.(string))
				} else if portAttributeName == "nodeId" {
					member.NodeId = basetypes.NewStringValue(portAttributeValue.(string))
				} else if portAttributeName == "nodeName" {
					member.NodeName = basetypes.NewStringValue(portAttributeValue.(string))
				}
			}
		} else if attributeName == "vlanId" {
			member.VlanId = basetypes.NewFloat64Value(attributeValue.(float64))
		} else if attributeName == "untagged" {
			isUntagged = attributeValue.(bool)
		}
	}
	if isUntagged {
		member.VlanId = basetypes.NewFloat64Null()
	}
	tflog.Debug(ctx, fmt.Sprintf("MEMBER LHLOG2 '%v' & '%v' & '%v'", data, attributes, member))
	return member
}

func NewMembersSet(ctx context.Context, data *[]MemberResourceModel, requestData []interface{}) basetypes.SetValue {
	members := make([]MemberResourceModel, 0)
	for _, member := range requestData {
		newMember := NewMemberResourceModel(ctx, getEmptyMemberResourceModel(), member.(map[string]interface{}))
		members = append(members, newMember)
	}
	membersSet, _ := types.SetValueFrom(ctx, MemberResourceModelAttributeType(), members)
	return membersSet
}

func NewMembersSetFromSetValue(ctx context.Context, data *[]MemberResourceModel, requestData []interface{}) basetypes.SetValue {
	members := make([]MemberResourceModel, 0)
	for index, member := range requestData {
		newMember := NewMemberResourceModel(ctx, &(*data)[index], member.(map[string]interface{}))
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
