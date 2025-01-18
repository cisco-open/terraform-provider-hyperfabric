// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Jeffail/gabs/v2"
	"github.com/cisco-open/terraform-provider-hyperfabric/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &NodeSubInterfaceResource{}
var _ resource.ResourceWithImportState = &NodeSubInterfaceResource{}

func NewNodeSubInterfaceResource() resource.Resource {
	return &NodeSubInterfaceResource{}
}

// NodeSubInterfaceResource defines the resource implementation.
type NodeSubInterfaceResource struct {
	client *client.Client
}

// NodeSubInterfaceResourceModel describes the resource data model.
type NodeSubInterfaceResourceModel struct {
	Id             types.String  `tfsdk:"id"`
	SubInterfaceId types.String  `tfsdk:"sub_interface_id"`
	NodeId         types.String  `tfsdk:"node_id"`
	Name           types.String  `tfsdk:"name"`
	Description    types.String  `tfsdk:"description"`
	Enabled        types.Bool    `tfsdk:"enabled"`
	Ipv4Addresses  types.Set     `tfsdk:"ipv4_addresses"`
	Ipv6Addresses  types.Set     `tfsdk:"ipv6_addresses"`
	VlanId         types.Float64 `tfsdk:"vlan_id"`
	VrfId          types.String  `tfsdk:"vrf_id"`
	Parent         types.String  `tfsdk:"parent"`
	Metadata       types.Object  `tfsdk:"metadata"`
	Labels         types.Set     `tfsdk:"labels"`
	Annotations    types.Set     `tfsdk:"annotations"`
}

func getEmptyNodeSubInterfaceResourceModel() *NodeSubInterfaceResourceModel {
	return &NodeSubInterfaceResourceModel{
		Id:             basetypes.NewStringNull(),
		SubInterfaceId: basetypes.NewStringNull(),
		NodeId:         basetypes.NewStringNull(),
		Name:           basetypes.NewStringNull(),
		Description:    basetypes.NewStringNull(),
		Enabled:        basetypes.NewBoolValue(false),
		Ipv4Addresses:  basetypes.NewSetNull(SetStringResourceModelAttributeType()),
		Ipv6Addresses:  basetypes.NewSetNull(SetStringResourceModelAttributeType()),
		VlanId:         basetypes.NewFloat64Null(),
		VrfId:          basetypes.NewStringNull(),
		Parent:         basetypes.NewStringNull(),
		Metadata:       basetypes.NewObjectNull(MetadataResourceModelAttributeType()),
		Labels:         basetypes.NewSetNull(SetStringResourceModelAttributeType()),
		Annotations:    basetypes.NewSetNull(AnnotationResourceModelAttributeType()),
	}
}

func getNewNodeSubInterfaceResourceModelFromData(data *NodeSubInterfaceResourceModel) *NodeSubInterfaceResourceModel {
	newNodeSubInterface := getEmptyNodeSubInterfaceResourceModel()

	if !data.Id.IsNull() && !data.Id.IsUnknown() {
		newNodeSubInterface.Id = data.Id
	}

	if !data.SubInterfaceId.IsNull() && !data.SubInterfaceId.IsUnknown() {
		newNodeSubInterface.SubInterfaceId = data.SubInterfaceId
	}

	if !data.NodeId.IsNull() && !data.NodeId.IsUnknown() {
		newNodeSubInterface.NodeId = data.NodeId
	}

	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		newNodeSubInterface.Name = data.Name
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		newNodeSubInterface.Description = data.Description
	}

	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		newNodeSubInterface.Enabled = data.Enabled
	}

	if !data.Ipv4Addresses.IsNull() && !data.Ipv4Addresses.IsUnknown() {
		newNodeSubInterface.Ipv4Addresses = data.Ipv4Addresses
	}

	if !data.Ipv6Addresses.IsNull() && !data.Ipv6Addresses.IsUnknown() {
		newNodeSubInterface.Ipv6Addresses = data.Ipv6Addresses
	}

	if !data.VlanId.IsNull() && !data.VlanId.IsUnknown() {
		newNodeSubInterface.VlanId = data.VlanId
	}

	if !data.VrfId.IsNull() && !data.VrfId.IsUnknown() {
		newNodeSubInterface.VrfId = data.VrfId
	}

	if !data.Parent.IsNull() && !data.Parent.IsUnknown() {
		newNodeSubInterface.Parent = data.Parent
	}

	if !data.Metadata.IsNull() && !data.Metadata.IsUnknown() {
		newNodeSubInterface.Metadata = data.Metadata
	}

	if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
		newNodeSubInterface.Labels = data.Labels
	}

	if !data.Annotations.IsNull() && !data.Annotations.IsUnknown() {
		newNodeSubInterface.Annotations = data.Annotations
	}

	return newNodeSubInterface
}

type NodeSubInterfaceIdentifier struct {
	Id types.String
}

func (r *NodeSubInterfaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Debug(ctx, "Start metadata of resource: hyperfabric_node_sub_interface")
	resp.TypeName = req.ProviderTypeName + "_node_sub_interface"
	tflog.Debug(ctx, "End metadata of resource: hyperfabric_node_sub_interface")
}

func (r *NodeSubInterfaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Debug(ctx, "Start schema of resource: hyperfabric_node_sub_interface")
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Node Sub-Interface resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "`id` defines the unique identifier of the Sub-Interface of a Node in a Fabric.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"sub_interface_id": schema.StringAttribute{
				MarkdownDescription: "`sub_interface_id` defines the unique identifier of a Sub-Interface of a Node in a Fabric.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
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
				MarkdownDescription: "The name of the Sub-Interface of the Node. The name should be in the `<Port Name>.<Integer>` format (i.e. `Ethernet1_1.100`). If `vlan_id` attribute is not provided, the integer in the Sub-Interface name will be used as the encapsulation VLAN ID.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description is a user defined field to store notes about the Sub-Interface of the Node.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "The enabled admin state of the Sub-Interface of the Node.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ipv4_addresses": getSubInterfaceIpv4AddressesSchemaAttribute(),
			"ipv6_addresses": getSubInterfaceIpv6AddressesSchemaAttribute(),
			"vlan_id": schema.Float64Attribute{
				MarkdownDescription: "The VLAN ID to use as encapsulation for the Sub-Interface of the Node. If not provided, the integer in the Sub-Interface `name` will be used as the encapsulation VLAN ID.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.RequiresReplace(),
					float64planmodifier.UseStateForUnknown(),
					SetToFloat64NullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"vrf_id": schema.StringAttribute{
				MarkdownDescription: "The `vrf_id` of a VRF to associate with the Sub-Interface of the Node.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"parent": schema.StringAttribute{
				MarkdownDescription: "The name of the `parent` Port of the Sub-Interface of the Node.",
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
	tflog.Debug(ctx, "End schema of resource: hyperfabric_node_sub_interface")
}

func getSubInterfaceIpv4AddressesSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A set of IPv4 addresses to be configured on the Sub-Interface of the Node.`,
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.UseStateForUnknown(),
		},
		ElementType: types.StringType,
	}
}

func getSubInterfaceIpv4AddressesDataSourceSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A set of IPv4 addresses configured on the Sub-Interface of the Node.`,
		Computed:            true,
		ElementType:         types.StringType,
	}
}

func getSubInterfaceIpv6AddressesSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A set of IPv6 addresses to be configured on the Sub-Interface of the Node.`,
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.UseStateForUnknown(),
		},
		ElementType: types.StringType,
	}
}

func getSubInterfaceIpv6AddressesDataSourceSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A set of IPv6 addresses configured on the Sub-Interface of the Node.`,
		Computed:            true,
		ElementType:         types.StringType,
	}
}

func (r *NodeSubInterfaceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of resource: hyperfabric_node_sub_interface")
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
	tflog.Debug(ctx, "End configure of resource: hyperfabric_node_sub_interface")
}

func (r *NodeSubInterfaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Start create of resource: hyperfabric_node_sub_interface")

	var data *NodeSubInterfaceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Create of resource hyperfabric_node_sub_interface with name '%s'", data.Name.ValueString()))

	jsonPayload := getNodeSubInterfaceJsonPayload(ctx, &resp.Diagnostics, data, "create")
	if resp.Diagnostics.HasError() {
		return
	}

	container := DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/subInterfaces", data.NodeId.ValueString()), "POST", jsonPayload)
	if resp.Diagnostics.HasError() {
		return
	}

	subInterfaceContainer, err := container.ArrayElement(0, "subInterfaces")
	if err != nil {
		return
	}

	subInterfaceId := StripQuotes(subInterfaceContainer.Search("id").String())
	if subInterfaceId != "" {
		data.Id = basetypes.NewStringValue(fmt.Sprintf("%s/subInterfaces/%s", data.NodeId.ValueString(), subInterfaceId))
		data.SubInterfaceId = basetypes.NewStringValue(subInterfaceId)
		getAndSetNodeSubInterfaceAttributes(ctx, &resp.Diagnostics, r.client, data)
	} else {
		data.Id = basetypes.NewStringNull()
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End create of resource hyperfabric_node_sub_interface with id '%s'", data.Id.ValueString()))
}

func (r *NodeSubInterfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Start read of resource: hyperfabric_node_sub_interface")
	var data *NodeSubInterfaceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Read of resource hyperfabric_node_sub_interface with id '%s'", data.Id.ValueString()))
	checkAndSetNodeSubInterfaceIds(data)
	getAndSetNodeSubInterfaceAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	if data.Id.IsNull() {
		var emptyData *NodeSubInterfaceResourceModel
		resp.Diagnostics.Append(resp.State.Set(ctx, &emptyData)...)
	} else {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
	tflog.Debug(ctx, fmt.Sprintf("End read of resource hyperfabric_node_sub_interface with id '%s'", data.Id.ValueString()))
}

func (r *NodeSubInterfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Start update of resource: hyperfabric_node_sub_interface")
	var data *NodeSubInterfaceResourceModel
	var stateData *NodeSubInterfaceResourceModel

	// Read Terraform plan and state data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Update of resource hyperfabric_node_sub_interface with id '%s'", data.Id.ValueString()))

	jsonPayload := getNodeSubInterfaceJsonPayload(ctx, &resp.Diagnostics, data, "update")

	if resp.Diagnostics.HasError() {
		return
	}

	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/subInterfaces/%s", data.NodeId.ValueString(), data.SubInterfaceId.ValueString()), "PUT", jsonPayload)

	if resp.Diagnostics.HasError() {
		return
	}

	getAndSetNodeSubInterfaceAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End update of resource hyperfabric_node_sub_interface with id '%s'", data.Id.ValueString()))
}

func (r *NodeSubInterfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Start delete of resource: hyperfabric_node_sub_interface")
	var data *NodeSubInterfaceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Delete of resource hyperfabric_node_sub_interface with id '%s'", data.Id.ValueString()))
	checkAndSetNodeSubInterfaceIds(data)
	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/subInterfaces/%s", data.NodeId.ValueString(), data.SubInterfaceId.ValueString()), "DELETE", nil)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("End delete of resource hyperfabric_node_sub_interface with id '%s'", data.Id.ValueString()))
}

func (r *NodeSubInterfaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Start import state of resource: hyperfabric_node_sub_interface")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	var stateData *NodeSubInterfaceResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &stateData)...)
	tflog.Debug(ctx, fmt.Sprintf("Import state of resource hyperfabric_node_sub_interface with id '%s'", stateData.Id.ValueString()))
	tflog.Debug(ctx, "End import of state resource: hyperfabric_node_sub_interface")
}

func getAndSetNodeSubInterfaceAttributes(ctx context.Context, diags *diag.Diagnostics, client *client.Client, data *NodeSubInterfaceResourceModel) {
	requestData := DoRestRequest(ctx, diags, client, fmt.Sprintf("/api/v1/fabrics/%s/subInterfaces/%s", data.NodeId.ValueString(), data.SubInterfaceId.ValueString()), "GET", nil)
	if diags.HasError() {
		return
	}

	newNodeSubInterface := *getNewNodeSubInterfaceResourceModelFromData(data)
	node := getEmptyNodeResourceModel()
	node.Id = newNodeSubInterface.NodeId
	checkAndSetNodeIds(node)

	if requestData.Data() != nil {
		for attributeName, attributeValue := range requestData.Data().(map[string]interface{}) {
			if attributeName == "id" && (data.SubInterfaceId.IsNull() || data.SubInterfaceId.IsUnknown() || data.SubInterfaceId.ValueString() == "" || data.SubInterfaceId.ValueString() != attributeValue.(string)) {
				newNodeSubInterface.SubInterfaceId = basetypes.NewStringValue(attributeValue.(string))
				newNodeSubInterface.Id = basetypes.NewStringValue(fmt.Sprintf("%s/subInterfaces/%s", newNodeSubInterface.NodeId.ValueString(), newNodeSubInterface.SubInterfaceId.ValueString()))
			} else if attributeName == "fabricId" && (node.FabricId.IsNull() || node.FabricId.IsUnknown() || node.FabricId.ValueString() == "" || node.FabricId.ValueString() != attributeValue.(string)) {
				node.FabricId = basetypes.NewStringValue(attributeValue.(string))
				newNodeSubInterface.NodeId = basetypes.NewStringValue(fmt.Sprintf("%s/nodes/%s", node.FabricId.ValueString(), node.NodeId.ValueString()))
				newNodeSubInterface.Id = basetypes.NewStringValue(fmt.Sprintf("%s/subInterfaces/%s", newNodeSubInterface.NodeId.ValueString(), newNodeSubInterface.SubInterfaceId.ValueString()))
			} else if attributeName == "nodeId" && (node.NodeId.IsNull() || node.NodeId.IsUnknown() || node.NodeId.ValueString() == "" || node.NodeId.ValueString() != attributeValue.(string)) {
				node.NodeId = basetypes.NewStringValue(attributeValue.(string))
				newNodeSubInterface.NodeId = basetypes.NewStringValue(fmt.Sprintf("%s/nodes/%s", node.FabricId.ValueString(), node.NodeId.ValueString()))
				newNodeSubInterface.Id = basetypes.NewStringValue(fmt.Sprintf("%s/subInterfaces/%s", newNodeSubInterface.NodeId.ValueString(), newNodeSubInterface.SubInterfaceId.ValueString()))
			} else if attributeName == "name" {
				newNodeSubInterface.Name = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "description" {
				newNodeSubInterface.Description = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "enabled" {
				newNodeSubInterface.Enabled = basetypes.NewBoolValue(attributeValue.(bool))
			} else if attributeName == "ipv4Addresses" {
				newNodeSubInterface.Ipv4Addresses = NewSetString(ctx, attributeValue.([]interface{}))
			} else if attributeName == "ipv6Addresses" {
				newNodeSubInterface.Ipv6Addresses = NewSetString(ctx, attributeValue.([]interface{}))
			} else if attributeName == "vlanId" {
				newNodeSubInterface.VlanId = basetypes.NewFloat64Value(attributeValue.(float64))
			} else if attributeName == "vrfId" {
				newNodeSubInterface.VrfId = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "parent" {
				newNodeSubInterface.Parent = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "metadata" {
				newNodeSubInterface.Metadata = NewMetadataObject(ctx, attributeValue.(map[string]interface{}))
			} else if attributeName == "labels" {
				newNodeSubInterface.Labels = NewSetString(ctx, attributeValue.([]interface{}))
			} else if attributeName == "annotations" {
				newNodeSubInterface.Annotations = NewAnnotationsSet(ctx, attributeValue.([]interface{}))
			}
		}
	} else {
		newNodeSubInterface.Id = basetypes.NewStringNull()
	}
	*data = newNodeSubInterface
}

func getNodeSubInterfaceJsonPayload(ctx context.Context, diags *diag.Diagnostics, data *NodeSubInterfaceResourceModel, action string) *gabs.Container {
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
		payloadMap["ipv6Addresses"] = getSetStringJsonPayload(ctx, data.Ipv6Addresses)
	}

	if !data.VlanId.IsNull() && !data.VlanId.IsUnknown() {
		payloadMap["vlanId"] = data.VlanId.ValueFloat64()
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
		payload = map[string]interface{}{"subInterfaces": payloadList}
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

func checkAndSetNodeSubInterfaceIds(data *NodeSubInterfaceResourceModel) {
	if strings.Contains(data.Id.ValueString(), "/subInterfaces/") {
		if data.NodeId.IsNull() || data.NodeId.IsUnknown() || data.NodeId.ValueString() == "" ||
			data.SubInterfaceId.IsNull() || data.SubInterfaceId.IsUnknown() || data.SubInterfaceId.ValueString() == "" {
			splitId := strings.Split(data.Id.ValueString(), "/subInterfaces/")
			data.NodeId = basetypes.NewStringValue(splitId[0])
			data.SubInterfaceId = basetypes.NewStringValue(splitId[1])
		}
	}
}
