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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &NodeLoopbackResource{}
var _ resource.ResourceWithImportState = &NodeLoopbackResource{}

func NewNodeLoopbackResource() resource.Resource {
	return &NodeLoopbackResource{}
}

// NodeLoopbackResource defines the resource implementation.
type NodeLoopbackResource struct {
	client *client.Client
}

// NodeLoopbackResourceModel describes the resource data model.
type NodeLoopbackResourceModel struct {
	Id          types.String `tfsdk:"id"`
	NodeId      types.String `tfsdk:"node_id"`
	LoopbackId  types.String `tfsdk:"loopback_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	// Enabled     types.Bool   `tfsdk:"enabled"`
	Ipv4Address types.String `tfsdk:"ipv4_address"`
	Ipv6Address types.String `tfsdk:"ipv6_address"`
	VrfId       types.String `tfsdk:"vrf_id"`
	Metadata    types.Object `tfsdk:"metadata"`
	Labels      types.Set    `tfsdk:"labels"`
	Annotations types.Set    `tfsdk:"annotations"`
}

func getEmptyNodeLoopbackResourceModel() *NodeLoopbackResourceModel {
	return &NodeLoopbackResourceModel{
		Id:          basetypes.NewStringNull(),
		NodeId:      basetypes.NewStringNull(),
		LoopbackId:  basetypes.NewStringNull(),
		Name:        basetypes.NewStringNull(),
		Description: basetypes.NewStringNull(),
		// Enabled:     basetypes.NewBoolValue(false),
		Ipv4Address: basetypes.NewStringNull(),
		Ipv6Address: basetypes.NewStringNull(),
		VrfId:       basetypes.NewStringNull(),
		Metadata:    basetypes.NewObjectNull(MetadataResourceModelAttributeType()),
		Labels:      basetypes.NewSetNull(SetStringResourceModelAttributeType()),
		Annotations: basetypes.NewSetNull(AnnotationResourceModelAttributeType()),
	}
}

func getNewNodeLoopbackResourceModelFromData(data *NodeLoopbackResourceModel) *NodeLoopbackResourceModel {
	newNodeLoopback := getEmptyNodeLoopbackResourceModel()

	if !data.Id.IsNull() && !data.Id.IsUnknown() {
		newNodeLoopback.Id = data.Id
	}

	if !data.LoopbackId.IsNull() && !data.LoopbackId.IsUnknown() {
		newNodeLoopback.LoopbackId = data.LoopbackId
	}

	if !data.NodeId.IsNull() && !data.NodeId.IsUnknown() {
		newNodeLoopback.NodeId = data.NodeId
	}

	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		newNodeLoopback.Name = data.Name
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		newNodeLoopback.Description = data.Description
	}

	// if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
	// 	newNodeLoopback.Enabled = data.Enabled
	// }

	if !data.Ipv4Address.IsNull() && !data.Ipv4Address.IsUnknown() {
		newNodeLoopback.Ipv4Address = data.Ipv4Address
	}

	if !data.Ipv6Address.IsNull() && !data.Ipv6Address.IsUnknown() {
		newNodeLoopback.Ipv6Address = data.Ipv6Address
	}

	if !data.VrfId.IsNull() && !data.VrfId.IsUnknown() {
		newNodeLoopback.VrfId = data.VrfId
	}

	if !data.Metadata.IsNull() && !data.Metadata.IsUnknown() {
		newNodeLoopback.Metadata = data.Metadata
	}

	if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
		newNodeLoopback.Labels = data.Labels
	}

	if !data.Annotations.IsNull() && !data.Annotations.IsUnknown() {
		newNodeLoopback.Annotations = data.Annotations
	}

	return newNodeLoopback
}

type NodeLoopbackIdentifier struct {
	Id types.String
}

func (r *NodeLoopbackResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Debug(ctx, "Start metadata of resource: hyperfabric_node_loopback")
	resp.TypeName = req.ProviderTypeName + "_node_loopback"
	tflog.Debug(ctx, "End metadata of resource: hyperfabric_node_loopback")
}

func (r *NodeLoopbackResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Debug(ctx, "Start schema of resource: hyperfabric_node_loopback")
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Node Loopback resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "`id` defines the unique identifier of the Loopback of a Node in a Fabric.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"loopback_id": schema.StringAttribute{
				MarkdownDescription: "`loopback_id` defines the unique identifier of a Loopback of a Node.",
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
				MarkdownDescription: "The name of the Loopback of the Node.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description is a user defined field to store notes about the Loopback of the Node.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			// "enabled": schema.BoolAttribute{
			// 	MarkdownDescription: "The enabled admin state of the Loopback of the Node.",
			// 	Optional:            true,
			// 	Computed:            true,
			// 	// Default:             booldefault.StaticBool(true),
			// 	PlanModifiers: []planmodifier.Bool{
			// 		boolplanmodifier.UseStateForUnknown(),
			// 	},
			// },
			"ipv4_address": schema.StringAttribute{
				MarkdownDescription: "The IPv4 address configured on the Loopback of the Node.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"ipv6_address": schema.StringAttribute{
				MarkdownDescription: "The IPv6 address configured on the Loopback of the Node.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"vrf_id": schema.StringAttribute{
				MarkdownDescription: "The `vrf_id` of a VRF to associate with the Loopback of the Node.",
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
	tflog.Debug(ctx, "End schema of resource: hyperfabric_node_loopback")
}

func (r *NodeLoopbackResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of resource: hyperfabric_node_loopback")
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
	tflog.Debug(ctx, "End configure of resource: hyperfabric_node_loopback")
}

func (r *NodeLoopbackResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Start create of resource: hyperfabric_node_loopback")

	var data *NodeLoopbackResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Create of resource hyperfabric_node_loopback with name '%s'", data.Name.ValueString()))

	jsonPayload := getNodeLoopbackJsonPayload(ctx, &resp.Diagnostics, data, "create")
	if resp.Diagnostics.HasError() {
		return
	}

	container := DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/loopbacks", data.NodeId.ValueString()), "POST", jsonPayload)
	if resp.Diagnostics.HasError() {
		return
	}

	loopbackContainer, err := container.ArrayElement(0, "loopbacks")
	if err != nil {
		return
	}

	loopbackId := StripQuotes(loopbackContainer.Search("id").String())
	if loopbackId != "" {
		data.Id = basetypes.NewStringValue(fmt.Sprintf("%s/loopbacks/%s", data.NodeId.ValueString(), loopbackId))
		data.LoopbackId = basetypes.NewStringValue(loopbackId)
		getAndSetNodeLoopbackAttributes(ctx, &resp.Diagnostics, r.client, data)
	} else {
		data.Id = basetypes.NewStringNull()
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End create of resource hyperfabric_node_loopback with id '%s'", data.Id.ValueString()))
}

func (r *NodeLoopbackResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Start read of resource: hyperfabric_node_loopback")
	var data *NodeLoopbackResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Read of resource hyperfabric_node_loopback with id '%s'", data.Id.ValueString()))
	checkAndSetNodeLoopbackIds(data)
	getAndSetNodeLoopbackAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	if data.Id.IsNull() {
		var emptyData *NodeLoopbackResourceModel
		resp.Diagnostics.Append(resp.State.Set(ctx, &emptyData)...)
	} else {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
	tflog.Debug(ctx, fmt.Sprintf("End read of resource hyperfabric_node_loopback with id '%s'", data.Id.ValueString()))
}

func (r *NodeLoopbackResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Start update of resource: hyperfabric_node_loopback")
	var data *NodeLoopbackResourceModel
	var stateData *NodeLoopbackResourceModel

	// Read Terraform plan and state data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Update of resource hyperfabric_node_loopback with id '%s'", data.Id.ValueString()))

	jsonPayload := getNodeLoopbackJsonPayload(ctx, &resp.Diagnostics, data, "update")

	if resp.Diagnostics.HasError() {
		return
	}

	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/loopbacks/%s", data.NodeId.ValueString(), data.LoopbackId.ValueString()), "PUT", jsonPayload)

	if resp.Diagnostics.HasError() {
		return
	}

	getAndSetNodeLoopbackAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End update of resource hyperfabric_node_loopback with id '%s'", data.Id.ValueString()))
}

func (r *NodeLoopbackResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Start delete of resource: hyperfabric_node_loopback")
	var data *NodeLoopbackResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Delete of resource hyperfabric_node_loopback with id '%s'", data.Id.ValueString()))
	checkAndSetNodeLoopbackIds(data)
	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/loopbacks/%s", data.NodeId.ValueString(), data.LoopbackId.ValueString()), "DELETE", nil)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("End delete of resource hyperfabric_node_loopback with id '%s'", data.Id.ValueString()))
}

func (r *NodeLoopbackResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Start import state of resource: hyperfabric_node_loopback")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	var stateData *NodeLoopbackResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &stateData)...)
	tflog.Debug(ctx, fmt.Sprintf("Import state of resource hyperfabric_node_loopback with id '%s'", stateData.Id.ValueString()))
	tflog.Debug(ctx, "End import of state resource: hyperfabric_node_loopback")
}

func getAndSetNodeLoopbackAttributes(ctx context.Context, diags *diag.Diagnostics, client *client.Client, data *NodeLoopbackResourceModel) {
	requestData := DoRestRequest(ctx, diags, client, fmt.Sprintf("/api/v1/fabrics/%s/loopbacks/%s", data.NodeId.ValueString(), data.LoopbackId.ValueString()), "GET", nil)
	if diags.HasError() {
		return
	}

	newNodeLoopback := *getNewNodeLoopbackResourceModelFromData(data)
	node := getEmptyNodeResourceModel()
	node.Id = newNodeLoopback.NodeId
	checkAndSetNodeIds(node)

	if requestData.Data() != nil {
		for attributeName, attributeValue := range requestData.Data().(map[string]interface{}) {
			if attributeName == "id" && (data.LoopbackId.IsNull() || data.LoopbackId.IsUnknown() || data.LoopbackId.ValueString() == "" || data.LoopbackId.ValueString() != attributeValue.(string)) {
				newNodeLoopback.LoopbackId = basetypes.NewStringValue(attributeValue.(string))
				newNodeLoopback.Id = basetypes.NewStringValue(fmt.Sprintf("%s/loopbacks/%s", newNodeLoopback.NodeId.ValueString(), newNodeLoopback.LoopbackId.ValueString()))
			} else if attributeName == "fabricId" && (node.FabricId.IsNull() || node.FabricId.IsUnknown() || node.FabricId.ValueString() == "" || node.FabricId.ValueString() != attributeValue.(string)) {
				node.FabricId = basetypes.NewStringValue(attributeValue.(string))
				newNodeLoopback.NodeId = basetypes.NewStringValue(fmt.Sprintf("%s/nodes/%s", node.FabricId.ValueString(), node.NodeId.ValueString()))
				newNodeLoopback.Id = basetypes.NewStringValue(fmt.Sprintf("%s/loopbacks/%s", newNodeLoopback.NodeId.ValueString(), newNodeLoopback.LoopbackId.ValueString()))
			} else if attributeName == "nodeId" && (node.NodeId.IsNull() || node.NodeId.IsUnknown() || node.NodeId.ValueString() == "" || node.NodeId.ValueString() != attributeValue.(string)) {
				node.NodeId = basetypes.NewStringValue(attributeValue.(string))
				newNodeLoopback.NodeId = basetypes.NewStringValue(fmt.Sprintf("%s/nodes/%s", node.FabricId.ValueString(), node.NodeId.ValueString()))
				newNodeLoopback.Id = basetypes.NewStringValue(fmt.Sprintf("%s/loopbacks/%s", newNodeLoopback.NodeId.ValueString(), newNodeLoopback.LoopbackId.ValueString()))
			} else if attributeName == "name" {
				newNodeLoopback.Name = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "description" {
				newNodeLoopback.Description = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "ipv4Address" {
				newNodeLoopback.Ipv4Address = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "ipv6Address" {
				newNodeLoopback.Ipv6Address = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "vrfId" {
				newNodeLoopback.VrfId = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "metadata" {
				newNodeLoopback.Metadata = NewMetadataObject(ctx, attributeValue.(map[string]interface{}))
			} else if attributeName == "labels" {
				newNodeLoopback.Labels = NewSetString(ctx, attributeValue.([]interface{}))
			} else if attributeName == "annotations" {
				newNodeLoopback.Annotations = NewAnnotationsSet(ctx, attributeValue.([]interface{}))
			}
		}
	} else {
		newNodeLoopback.Id = basetypes.NewStringNull()
	}
	*data = newNodeLoopback
}

func getNodeLoopbackJsonPayload(ctx context.Context, diags *diag.Diagnostics, data *NodeLoopbackResourceModel, action string) *gabs.Container {
	payloadMap := map[string]interface{}{}
	payloadList := []map[string]interface{}{}

	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		payloadMap["name"] = data.Name.ValueString()
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		payloadMap["description"] = data.Description.ValueString()
	}

	// if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
	// 	payloadMap["enabled"] = data.Enabled.ValueBool()
	// }
	payloadMap["enabled"] = true

	if !data.Ipv4Address.IsNull() && !data.Ipv4Address.IsUnknown() {
		payloadMap["ipv4Address"] = data.Ipv4Address.ValueString()
	}

	if !data.Ipv6Address.IsNull() && !data.Ipv6Address.IsUnknown() {
		payloadMap["ipv6Address"] = data.Ipv6Address.ValueString()
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
		payload = map[string]interface{}{"loopbacks": payloadList}
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

func checkAndSetNodeLoopbackIds(data *NodeLoopbackResourceModel) {
	if strings.Contains(data.Id.ValueString(), "/loopbacks/") {
		if data.NodeId.IsNull() || data.NodeId.IsUnknown() || data.NodeId.ValueString() == "" ||
			data.LoopbackId.IsNull() || data.LoopbackId.IsUnknown() || data.LoopbackId.ValueString() == "" {
			splitId := strings.Split(data.Id.ValueString(), "/loopbacks/")
			data.NodeId = basetypes.NewStringValue(splitId[0])
			data.LoopbackId = basetypes.NewStringValue(splitId[1])
		}
	}
}
