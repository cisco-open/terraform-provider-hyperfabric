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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Jeffail/gabs/v2"
	"github.com/cisco-open/terraform-provider-hyperfabric/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &NodePortResource{}
var _ resource.ResourceWithImportState = &NodePortResource{}

func NewNodePortResource() resource.Resource {
	return &NodePortResource{}
}

// NodePortResource defines the resource implementation.
type NodePortResource struct {
	client *client.Client
}

// NodePortResourceModel describes the resource data model.
type NodePortResourceModel struct {
	Id          types.String `tfsdk:"id"`
	NodeId      types.String `tfsdk:"node_id"`
	PortId      types.String `tfsdk:"port_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	// Breakout           types.Bool    `tfsdk:"breakout"`
	// BreakoutIndex      types.Float64 `tfsdk:"breakout_index"`
	Index              types.Float64 `tfsdk:"index"`
	Ipv4Addresses      types.Set     `tfsdk:"ipv4_addresses"`
	Ipv6Addresses      types.Set     `tfsdk:"ipv6_addresses"`
	Linecard           types.Float64 `tfsdk:"linecard"`
	PreventForwarding  types.Bool    `tfsdk:"prevent_forwarding"`
	LldpHost           types.String  `tfsdk:"lldp_host"`
	LldpInfo           types.String  `tfsdk:"lldp_info"`
	LldpPort           types.String  `tfsdk:"lldp_port"`
	MaxSpeed           types.String  `tfsdk:"max_speed"`
	Mtu                types.Float64 `tfsdk:"mtu"`
	Roles              types.Set     `tfsdk:"roles"`
	Speed              types.String  `tfsdk:"speed"`
	SubInterfacesCount types.Float64 `tfsdk:"sub_interfaces_count"`
	VlanIds            types.Set     `tfsdk:"vlan_ids"`
	Vnis               types.Set     `tfsdk:"vnis"`
	VrfId              types.String  `tfsdk:"vrf_id"`
	Metadata           types.Object  `tfsdk:"metadata"`
	Labels             types.Set     `tfsdk:"labels"`
	Annotations        types.Set     `tfsdk:"annotations"`
}

func getEmptyNodePortResourceModel() *NodePortResourceModel {
	return &NodePortResourceModel{
		Id:          basetypes.NewStringNull(),
		NodeId:      basetypes.NewStringNull(),
		PortId:      basetypes.NewStringNull(),
		Name:        basetypes.NewStringNull(),
		Description: basetypes.NewStringNull(),
		Enabled:     basetypes.NewBoolValue(false),
		// Breakout:           basetypes.NewBoolValue(false),
		// BreakoutIndex:      basetypes.NewFloat64Null(),
		Index:              basetypes.NewFloat64Null(),
		Ipv4Addresses:      basetypes.NewSetNull(SetStringResourceModelAttributeType()),
		Ipv6Addresses:      basetypes.NewSetNull(SetStringResourceModelAttributeType()),
		Linecard:           basetypes.NewFloat64Null(),
		PreventForwarding:  basetypes.NewBoolValue(false),
		LldpHost:           basetypes.NewStringNull(),
		LldpInfo:           basetypes.NewStringNull(),
		LldpPort:           basetypes.NewStringNull(),
		MaxSpeed:           basetypes.NewStringNull(),
		Mtu:                basetypes.NewFloat64Null(),
		Roles:              basetypes.NewSetNull(SetStringResourceModelAttributeType()),
		Speed:              basetypes.NewStringNull(),
		SubInterfacesCount: basetypes.NewFloat64Null(),
		VlanIds:            basetypes.NewSetNull(SetStringResourceModelAttributeType()),
		Vnis:               basetypes.NewSetNull(SetStringResourceModelAttributeType()),
		VrfId:              basetypes.NewStringNull(),
		Metadata:           basetypes.NewObjectNull(MetadataResourceModelAttributeType()),
		Labels:             basetypes.NewSetNull(SetStringResourceModelAttributeType()),
		Annotations:        basetypes.NewSetNull(AnnotationResourceModelAttributeType()),
	}
}

type NodePortIdentifier struct {
	Id types.String
}

func (r *NodePortResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Debug(ctx, "Start metadata of resource: hyperfabric_node_port")
	resp.TypeName = req.ProviderTypeName + "_node_port"
	tflog.Debug(ctx, "End metadata of resource: hyperfabric_node_port")
}

func (r *NodePortResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Debug(ctx, "Start schema of resource: hyperfabric_node_port")
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Node Port resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "`id` defines the unique identifier of the Port of a Node in a Fabric.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"port_id": schema.StringAttribute{
				MarkdownDescription: "`port_id` defines the unique identifier of a Port of a Node.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"node_id": schema.StringAttribute{
				MarkdownDescription: "`node_id` defines the unique identifier of a Node in a Fabric.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Port of the Node.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description is a user defined field to store notes about the Port of the Node.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "The enabled admin state of the Port of the Node.",
				Optional:            true,
				Computed:            true,
				// Default:             booldefault.StaticBool(true),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			// "breakout": schema.BoolAttribute{
			// 	MarkdownDescription: "The breakout state of the Port of the Node.",
			// 	Optional:            true,
			// 	Computed:            true,
			// 	// Default:             booldefault.StaticBool(true),
			// 	PlanModifiers: []planmodifier.Bool{
			// 		boolplanmodifier.UseStateForUnknown(),
			// 	},
			// },
			// "breakout_index": schema.Float64Attribute{
			// 	MarkdownDescription: "The index of the sub-port on the breakout Port.",
			// 	Computed:            true,
			// 	PlanModifiers: []planmodifier.Float64{
			// 		float64planmodifier.UseStateForUnknown(),
			// 	},
			// },
			"index": schema.Float64Attribute{
				MarkdownDescription: "The index number of the Port of the Node.",
				Computed:            true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"ipv4_addresses": getIpv4AddressesSchemaAttribute(),
			"ipv6_addresses": getIpv6AddressesSchemaAttribute(),
			"linecard": schema.Float64Attribute{
				MarkdownDescription: "The linecard index number of the Port of the Node.",
				Computed:            true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"prevent_forwarding": schema.BoolAttribute{
				MarkdownDescription: "Prevent traffic from being forwarded by the Port.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"lldp_host": schema.StringAttribute{
				MarkdownDescription: "The name of host reported by LLDP connected to the Port of the Node.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"lldp_info": schema.StringAttribute{
				MarkdownDescription: "The info about the host reported by LLDP connected to the Port of the Node.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"lldp_port": schema.StringAttribute{
				MarkdownDescription: "The port of host reported by LLDP connected to the Port of the Node.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"max_speed": schema.StringAttribute{
				MarkdownDescription: "The maximum speed of the Port of the Node.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mtu": schema.Float64Attribute{
				MarkdownDescription: "The MTU of the Port of the Node.",
				// Optional:            true,
				Computed: true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"roles": getPortRolesSchemaAttribute(),
			"speed": schema.StringAttribute{
				MarkdownDescription: "The configured speed of the Port of the Node.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"sub_interfaces_count": schema.Float64Attribute{
				MarkdownDescription: "The number of sub-interfaces of the Port of the Node.",
				Computed:            true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"vlan_ids": getVlanIdsSchemaAttribute(),
			"vnis":     getVnisSchemaAttribute(),
			"vrf_id": schema.StringAttribute{
				MarkdownDescription: "The `vrf_id` to associate with the Port of the Node. Required when the Port roles include `ROUTED_PORT`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"metadata":    getMetadataSchemaAttribute(),
			"labels":      getLabelsSchemaAttribute(),
			"annotations": getAnnotationsSchemaAttribute(),
		},
	}
	tflog.Debug(ctx, "End schema of resource: hyperfabric_node_port")
}

func getIpv4AddressesSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A set of IPv4 addresses to be configured on the Port of the Node.`,
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.Set{
			// Validate this attribute must be configured with other_attr.
			setvalidator.AlsoRequires(path.Expressions{
				path.MatchRoot("vrf_id"),
			}...),
		},
		ElementType: types.StringType,
	}
}

func getIpv4AddressesDataSourceSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A set of IPv4 addresses configured on the Port of the Node.`,
		Computed:            true,
		ElementType:         types.StringType,
	}
}

func getIpv6AddressesSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A set of IPv6 addresses to be configured on the Port of the Node.`,
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.Set{
			// Validate this attribute must be configured with other_attr.
			setvalidator.AlsoRequires(path.Expressions{
				path.MatchRoot("vrf_id"),
			}...),
		},
		ElementType: types.StringType,
	}
}

func getIpv6AddressesDataSourceSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A set of IPv6 addresses configured on the Port of the Node.`,
		Computed:            true,
		ElementType:         types.StringType,
	}
}

func getPortRolesSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A set of roles used for the Port of the Node.`,
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.Set{
			setvalidator.ValueStringsAre(stringvalidator.OneOf([]string{"PORT_ROLE_UNSPECIFIED", "UNUSED_PORT", "FABRIC_PORT", "HOST_PORT", "ROUTED_PORT", "LAG_PORT"}...)),
		},
		ElementType: types.StringType,
	}
}

func getPortRolesDataSourceSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A set of roles used for the Port of the Node.`,
		Computed:            true,
		ElementType:         types.StringType,
	}
}

func getVlanIdsSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A set of VLAN IDs used by the Port of the Node.`,
		Computed:            true,
		ElementType:         types.StringType,
	}
}

func getVnisSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A set of VNIs used by the Port of the Node.`,
		Computed:            true,
		ElementType:         types.StringType,
	}
}

func (r *NodePortResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of resource: hyperfabric_node_port")
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
	tflog.Debug(ctx, "End configure of resource: hyperfabric_node_port")
}

func (r *NodePortResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Start create of resource: hyperfabric_node_port")

	var data *NodePortResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Create of resource hyperfabric_node_port with name '%s'", data.Name.ValueString()))

	jsonPayload := getNodePortJsonPayload(ctx, &resp.Diagnostics, data, "update")
	if resp.Diagnostics.HasError() {
		return
	}

	container := DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/ports/%s", data.NodeId.ValueString(), data.Name.ValueString()), "PUT", jsonPayload)
	if resp.Diagnostics.HasError() {
		return
	}

	portContainer, err := container.ArrayElement(0, "ports")
	if err != nil {
		return
	}

	portId := StripQuotes(portContainer.Search("id").String())
	if portId != "" {
		data.Id = basetypes.NewStringValue(fmt.Sprintf("%s/ports/%s", data.NodeId.ValueString(), portId))
		data.PortId = basetypes.NewStringValue(portId)
		getAndSetNodePortAttributes(ctx, &resp.Diagnostics, r.client, data)
	} else {
		data.Id = basetypes.NewStringNull()
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End create of resource hyperfabric_node_port with id '%s'", data.Id.ValueString()))
}

func (r *NodePortResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Start read of resource: hyperfabric_node_port")
	var data *NodePortResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Read of resource hyperfabric_node_port with id '%s'", data.Id.ValueString()))
	checkAndSetNodePortIds(data)
	getAndSetNodePortAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	if data.Id.IsNull() {
		var emptyData *NodePortResourceModel
		resp.Diagnostics.Append(resp.State.Set(ctx, &emptyData)...)
	} else {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
	tflog.Debug(ctx, fmt.Sprintf("End read of resource hyperfabric_node_port with id '%s'", data.Id.ValueString()))
}

func (r *NodePortResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Start update of resource: hyperfabric_node_port")
	var data *NodePortResourceModel
	var stateData *NodePortResourceModel

	// Read Terraform plan and state data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Update of resource hyperfabric_node_port with id '%s'", data.Id.ValueString()))

	jsonPayload := getNodePortJsonPayload(ctx, &resp.Diagnostics, data, "update")

	if resp.Diagnostics.HasError() {
		return
	}

	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/ports/%s", data.NodeId.ValueString(), data.PortId.ValueString()), "PUT", jsonPayload)

	if resp.Diagnostics.HasError() {
		return
	}

	getAndSetNodePortAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End update of resource hyperfabric_node_port with id '%s'", data.Id.ValueString()))
}

func (r *NodePortResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Start delete of resource: hyperfabric_node_port")
	var data *NodePortResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Delete of resource hyperfabric_node_port with id '%s'", data.Id.ValueString()))
	checkAndSetNodePortIds(data)
	// DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/ports/%s", data.NodeId.ValueString(), data.PortId.ValueString()), "DELETE", nil)
	jsonPayload := getNodePortJsonPayload(ctx, &resp.Diagnostics, data, "delete")
	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/ports/%s", data.NodeId.ValueString(), data.PortId.ValueString()), "PUT", jsonPayload)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("End delete of resource hyperfabric_node_port with id '%s'", data.Id.ValueString()))
}

func (r *NodePortResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Start import state of resource: hyperfabric_node_port")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	var stateData *NodePortResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &stateData)...)
	tflog.Debug(ctx, fmt.Sprintf("Import state of resource hyperfabric_node_port with id '%s'", stateData.Id.ValueString()))
	tflog.Debug(ctx, "End import of state resource: hyperfabric_node_port")
}

func getAndSetNodePortAttributes(ctx context.Context, diags *diag.Diagnostics, client *client.Client, data *NodePortResourceModel) {
	requestData := DoRestRequest(ctx, diags, client, fmt.Sprintf("/api/v1/fabrics/%s/ports/%s", data.NodeId.ValueString(), data.PortId.ValueString()), "GET", nil)
	if diags.HasError() {
		return
	}

	newNodePort := *getEmptyNodePortResourceModel()
	newNodePort.NodeId = data.NodeId
	newNodePort.PortId = data.PortId
	newNodePort.Id = data.Id

	if requestData.Data() != nil {
		for attributeName, attributeValue := range requestData.Data().(map[string]interface{}) {
			if attributeName == "id" && (data.PortId.IsNull() || data.PortId.IsUnknown() || data.PortId.ValueString() == "" || data.PortId.ValueString() != attributeValue.(string)) {
				newNodePort.PortId = basetypes.NewStringValue(attributeValue.(string))
				newNodePort.Id = basetypes.NewStringValue(fmt.Sprintf("%s/ports/%s", newNodePort.NodeId.ValueString(), newNodePort.PortId.ValueString()))
			} else if attributeName == "name" {
				newNodePort.Name = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "description" {
				newNodePort.Description = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "enabled" {
				newNodePort.Enabled = basetypes.NewBoolValue(attributeValue.(bool))
				// } else if attributeName == "breakout" {
				// 	newNodePort.Breakout = basetypes.NewBoolValue(attributeValue.(bool))
				// } else if attributeName == "breakoutIndex" {
				// 	newNodePort.BreakoutIndex = basetypes.NewFloat64Value(attributeValue.(float64))
			} else if attributeName == "index" {
				newNodePort.Index = basetypes.NewFloat64Value(attributeValue.(float64))
			} else if attributeName == "ipv4Addresses" {
				newNodePort.Ipv4Addresses = NewSetString(ctx, attributeValue.([]interface{}))
			} else if attributeName == "ipv6Addresses" {
				newNodePort.Ipv6Addresses = NewSetString(ctx, attributeValue.([]interface{}))
			} else if attributeName == "linecard" {
				newNodePort.Linecard = basetypes.NewFloat64Value(attributeValue.(float64))
			} else if attributeName == "linkDown" {
				newNodePort.PreventForwarding = basetypes.NewBoolValue(attributeValue.(bool))
			} else if attributeName == "lldpHost" {
				newNodePort.LldpHost = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "lldpInfo" {
				newNodePort.LldpInfo = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "lldpPort" {
				newNodePort.LldpPort = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "maxSpeed" {
				newNodePort.MaxSpeed = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "mtu" {
				newNodePort.Mtu = basetypes.NewFloat64Value(attributeValue.(float64))
			} else if attributeName == "roles" {
				newNodePort.Roles = NewSetString(ctx, attributeValue.([]interface{}))
			} else if attributeName == "speed" {
				newNodePort.Speed = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "subInfCount" {
				newNodePort.SubInterfacesCount = basetypes.NewFloat64Value(attributeValue.(float64))
			} else if attributeName == "vlanIds" {
				newNodePort.VlanIds = NewSetString(ctx, attributeValue.([]interface{}))
			} else if attributeName == "vnis" {
				newNodePort.Vnis = NewSetString(ctx, attributeValue.([]interface{}))
			} else if attributeName == "vrfId" {
				newNodePort.VrfId = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "metadata" {
				newNodePort.Metadata = NewMetadataObject(ctx, attributeValue.(map[string]interface{}))
			} else if attributeName == "labels" {
				newNodePort.Labels = NewSetString(ctx, attributeValue.([]interface{}))
			} else if attributeName == "annotations" {
				newNodePort.Annotations = NewAnnotationsSet(ctx, attributeValue.([]interface{}))
			}
		}
	} else {
		newNodePort.Id = basetypes.NewStringNull()
	}
	*data = newNodePort
}

func getNodePortJsonPayload(ctx context.Context, diags *diag.Diagnostics, data *NodePortResourceModel, action string) *gabs.Container {
	payloadMap := map[string]interface{}{}
	payloadList := []map[string]interface{}{}

	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		payloadMap["name"] = data.Name.ValueString()
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		payloadMap["description"] = data.Description.ValueString()
	}

	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		payloadMap["enabled"] = data.Enabled.ValueBool()
	}

	if !data.Ipv4Addresses.IsNull() && !data.Ipv4Addresses.IsUnknown() {
		payloadMap["ipv4Addresses"] = getSetStringJsonPayload(ctx, data.Ipv4Addresses)
	}

	if !data.Ipv6Addresses.IsNull() && !data.Ipv6Addresses.IsUnknown() {
		payloadMap["ipv6Addresses"] = getSetStringJsonPayload(ctx, data.Ipv4Addresses)
	}

	if !data.PreventForwarding.IsNull() && !data.PreventForwarding.IsUnknown() {
		payloadMap["linkDown"] = data.PreventForwarding.ValueBool()
	}

	if !data.Roles.IsNull() && !data.Roles.IsUnknown() {
		payloadMap["roles"] = getSetStringJsonPayload(ctx, data.Roles)
	}

	if !data.VrfId.IsNull() && !data.VrfId.IsUnknown() {
		payloadMap["vrfId"] = data.VrfId.ValueString()
	}

	if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
		payloadMap["labels"] = getSetStringJsonPayload(ctx, data.Labels)
	}

	if !data.Annotations.IsNull() && !data.Annotations.IsUnknown() {
		payloadMap["annotations"] = getAnnotationsJsonPayload(ctx, data.Annotations)
	}

	var payload map[string]interface{}
	if action == "create" {
		payloadList = append(payloadList, payloadMap)
		payload = map[string]interface{}{"ports": payloadList}
	} else if action == "delete" {
		payload = map[string]interface{}{
			"name":    data.Name.ValueString(),
			"enabled": true,
			"roles":   []string{"FABRIC_PORT"},
		}
	} else {
		payload = payloadMap
	}

	marshalPayload, err := json.Marshal(payload)
	if err != nil {
		diags.AddError(
			"Marshalling of JSON payload failed",
			fmt.Sprintf("Err: %s. Please report this issue to the provider developers.", err),
		)
		return nil
	}

	jsonPayload, err := gabs.ParseJSON(marshalPayload)
	if err != nil {
		diags.AddError(
			"Construction of JSON payload failed",
			fmt.Sprintf("Err: %s. Please report this issue to the provider developers.", err),
		)
		return nil
	}
	return jsonPayload
}

func checkAndSetNodePortIds(data *NodePortResourceModel) {
	if strings.Contains(data.Id.ValueString(), "/ports/") {
		if data.NodeId.IsNull() || data.NodeId.IsUnknown() || data.NodeId.ValueString() == "" ||
			data.PortId.IsNull() || data.PortId.IsUnknown() || data.PortId.ValueString() == "" {
			splitId := strings.Split(data.Id.ValueString(), "/ports/")
			data.NodeId = basetypes.NewStringValue(splitId[0])
			data.PortId = basetypes.NewStringValue(splitId[1])
		}
	}
}
