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
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ConnectionResource{}
var _ resource.ResourceWithImportState = &ConnectionResource{}

func NewConnectionResource() resource.Resource {
	return &ConnectionResource{}
}

// ConnectionResource defines the resource implementation.
type ConnectionResource struct {
	client *client.Client
}

// ConnectionResourceModel describes the resource data model.
type ConnectionResourceModel struct {
	Id           types.String  `tfsdk:"id"`
	ConnectionId types.String  `tfsdk:"connection_id"`
	FabricId     types.String  `tfsdk:"fabric_id"`
	Description  types.String  `tfsdk:"description"`
	CableType    types.String  `tfsdk:"cable_type"`
	CableLength  types.Float64 `tfsdk:"cable_length"`
	Pluggable    types.String  `tfsdk:"pluggable"`
	Local        types.Object  `tfsdk:"local"`
	Remote       types.Object  `tfsdk:"remote"`
	OsType       types.String  `tfsdk:"os_type"`
	Unrecognized types.Bool    `tfsdk:"unrecognized"`
	// Metadata    types.Object `tfsdk:"metadata"`
	// Labels        types.Set    `tfsdk:"labels"`
	// Annotations types.Set    `tfsdk:"annotations"`
}

func getEmptyConnectionResourceModel() *ConnectionResourceModel {
	return &ConnectionResourceModel{
		Id:           basetypes.NewStringNull(),
		ConnectionId: basetypes.NewStringNull(),
		FabricId:     basetypes.NewStringNull(),
		Description:  basetypes.NewStringNull(),
		CableType:    basetypes.NewStringValue("CABLE_TYPE_UNSPECIFIED"),
		CableLength:  basetypes.NewFloat64Null(),
		Pluggable:    basetypes.NewStringNull(),
		Local:        basetypes.NewObjectNull(LocalRemoteConnectionResourceModelAttributeType()),
		Remote:       basetypes.NewObjectNull(LocalRemoteConnectionResourceModelAttributeType()),
		OsType:       basetypes.NewStringNull(),
		Unrecognized: basetypes.NewBoolValue(false),
		// Metadata:    basetypes.NewObjectNull(MetadataResourceModelAttributeType()),
		// Labels:        basetypes.NewSetNull(SetStringResourceModelAttributeType()),
		// Annotations: basetypes.NewSetNull(AnnotationResourceModelAttributeType()),
	}
}

func getNewConnectionResourceModelFromData(data *ConnectionResourceModel) *ConnectionResourceModel {
	newConnection := getEmptyConnectionResourceModel()

	if !data.Id.IsNull() && !data.Id.IsUnknown() {
		newConnection.Id = data.Id
	}

	if !data.ConnectionId.IsNull() && !data.ConnectionId.IsUnknown() {
		newConnection.ConnectionId = data.ConnectionId
	}

	if !data.FabricId.IsNull() && !data.FabricId.IsUnknown() {
		newConnection.FabricId = data.FabricId
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		newConnection.Description = data.Description
	}

	if !data.CableType.IsNull() && !data.CableType.IsUnknown() {
		newConnection.CableType = data.CableType
	}

	if !data.CableLength.IsNull() && !data.CableLength.IsUnknown() {
		newConnection.CableLength = data.CableLength
	}

	if !data.Pluggable.IsNull() && !data.Pluggable.IsUnknown() {
		newConnection.Pluggable = data.Pluggable
	}

	if !data.Local.IsNull() && !data.Local.IsUnknown() {
		newConnection.Local = data.Local
	}

	if !data.Remote.IsNull() && !data.Remote.IsUnknown() {
		newConnection.Remote = data.Remote
	}

	if !data.OsType.IsNull() && !data.OsType.IsUnknown() {
		newConnection.OsType = data.OsType
	}

	if !data.Unrecognized.IsNull() && !data.Unrecognized.IsUnknown() {
		newConnection.Unrecognized = data.Unrecognized
	}

	// if !data.Metadata.IsNull() && !data.Metadata.IsUnknown() {
	// 	newConnection.Metadata = data.Metadata
	// }

	// if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
	// 	newConnection.Labels = data.Labels
	// }

	// if !data.Annotations.IsNull() && !data.Annotations.IsUnknown() {
	// 	newConnection.Annotations = data.Annotations
	// }

	return newConnection
}

type LocalRemoteConnectionResourceModel struct {
	NodeId   types.String `tfsdk:"node_id"`
	NodeName types.String `tfsdk:"node_name"`
	PortName types.String `tfsdk:"port_name"`
}

func LocalRemoteConnectionResourceModelAttributeType() map[string]attr.Type {
	return map[string]attr.Type{
		"node_id":   types.StringType,
		"node_name": types.StringType,
		"port_name": types.StringType,
	}
}

func getEmptyLocalRemoteConnectionResourceModel() LocalRemoteConnectionResourceModel {
	return LocalRemoteConnectionResourceModel{
		NodeId:   basetypes.NewStringNull(),
		NodeName: basetypes.NewStringNull(),
		PortName: basetypes.NewStringNull(),
	}
}

func getLocalRemoteConnectionSchemaAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: `An object that represents the local side of the Connection.`,
		Required:            true,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplace(),
		},
		Attributes: map[string]schema.Attribute{
			"node_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: `The ID of the Node used as as local/remote side of this Connection.`,
			},
			"port_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: `The name of the port on the Node used as local/remote side of this Connection.`,
			},
			"node_name": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: `The name of the Node used as local/remote side of this Connection.`,
			},
		},
	}
}

func NewLocalRemoteConnectionResourceModel(data map[string]interface{}) LocalRemoteConnectionResourceModel {
	localRemoteConnection := getEmptyLocalRemoteConnectionResourceModel()
	for attributeName, attributeValue := range data {
		if attributeName == "nodeId" && attributeValue != nil {
			stringAttr := attributeValue.(string)
			if stringAttr != "" {
				localRemoteConnection.NodeId = basetypes.NewStringValue(stringAttr)
			}
		} else if attributeName == "nodeName" && attributeValue != nil {
			stringAttr := attributeValue.(string)
			if stringAttr != "" {
				localRemoteConnection.NodeName = basetypes.NewStringValue(stringAttr)
			}
		} else if attributeName == "portName" && attributeValue != nil {
			stringAttr := attributeValue.(string)
			if stringAttr != "" {
				localRemoteConnection.PortName = basetypes.NewStringValue(stringAttr)
			}
		}
	}
	return localRemoteConnection
}

func NewLocalRemoteConnectionObject(ctx context.Context, data map[string]interface{}) basetypes.ObjectValue {
	localRemoteConnection := NewLocalRemoteConnectionResourceModel(data)
	localRemoteConnectionObject, _ := types.ObjectValueFrom(ctx, LocalRemoteConnectionResourceModelAttributeType(), localRemoteConnection)
	return localRemoteConnectionObject
}

func getLocalRemoteConnectionJsonPayload(ctx context.Context, data basetypes.ObjectValue) map[string]string {
	localRemoteConnection := LocalRemoteConnectionResourceModel{}
	data.As(ctx, &localRemoteConnection, basetypes.ObjectAsOptions{})
	localRemoteConnectionPayload := map[string]string{
		"nodeId":   StripQuotes(localRemoteConnection.NodeId.String()),
		"portName": StripQuotes(localRemoteConnection.PortName.String()),
	}
	return localRemoteConnectionPayload
}

type ConnectionIdentifier struct {
	Id types.String
}

func (r *ConnectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Debug(ctx, "Start metadata of resource: hyperfabric_connection")
	resp.TypeName = req.ProviderTypeName + "_connection"
	tflog.Debug(ctx, "End metadata of resource: hyperfabric_connection")
}

func (r *ConnectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Debug(ctx, "Start schema of resource: hyperfabric_connection")
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Connection resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "`id` defines the unique identifier of a Connection in a Fabric.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connection_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "`connection_id` defines the unique identifier of a Connection.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"fabric_id": schema.StringAttribute{
				MarkdownDescription: "`fabric_id` defines the unique identifier of a Fabric.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description is a user defined field to store notes about the Connection.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cable_type": schema.StringAttribute{
				MarkdownDescription: "The type of cable used for the Connection.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"CABLE_TYPE_UNSPECIFIED", "DAC", "FIBER"}...),
				},
			},
			"cable_length": schema.Float64Attribute{
				MarkdownDescription: "The length of the cable used for the Connection.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
					float64planmodifier.RequiresReplace(),
				},
			},
			"pluggable": schema.StringAttribute{
				MarkdownDescription: "The type of pluggable used for the Connection.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"local":  getLocalRemoteConnectionSchemaAttribute(),
			"remote": getLocalRemoteConnectionSchemaAttribute(),
			"os_type": schema.StringAttribute{
				MarkdownDescription: "The operating system type of the remote side of the Connection.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"unrecognized": schema.BoolAttribute{
				MarkdownDescription: "If the remote side of the Connection is recognized or not.",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			// "metadata":    getMetadataSchemaAttribute(),
			// "labels":        getLabelsSchemaAttribute(),
			// "annotations": getAnnotationsSchemaAttribute(),
		},
	}
	tflog.Debug(ctx, "End schema of resource: hyperfabric_connection")
}

func (r *ConnectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of resource: hyperfabric_connection")
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
	tflog.Debug(ctx, "End configure of resource: hyperfabric_connection")
}

func (r *ConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Start create of resource: hyperfabric_connection")

	var data *ConnectionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	local := data.Local.Attributes()
	remote := data.Remote.Attributes()
	tflog.Debug(ctx, fmt.Sprintf("Create of resource hyperfabric_connection in fabric '%s' with local node '%s' interface '%s' and remote node '%s' interface '%s'", data.FabricId.ValueString(), local["node_id"].String(), local["port_name"].String(), remote["node_id"].String(), remote["port_name"].String()))

	jsonPayload := getConnectionJsonPayload(ctx, &resp.Diagnostics, data, "create")
	if resp.Diagnostics.HasError() {
		return
	}

	container := DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/connections", data.FabricId.ValueString()), "POST", jsonPayload)
	if resp.Diagnostics.HasError() {
		return
	}

	connectionContainer, err := container.ArrayElement(0, "connections")
	if err != nil {
		return
	}
	connectionId := StripQuotes(connectionContainer.Search("id").String())
	if connectionId != "" {
		data.Id = basetypes.NewStringValue(fmt.Sprintf("%s/connections/%s", data.FabricId.ValueString(), connectionId))
		data.ConnectionId = basetypes.NewStringValue(connectionId)
		getAndSetConnectionAttributes(ctx, &resp.Diagnostics, r.client, data)
	} else {
		data.Id = basetypes.NewStringNull()
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End create of resource hyperfabric_connection with id '%s'", data.Id.ValueString()))
}

func (r *ConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Start read of resource: hyperfabric_connection")
	var data *ConnectionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Read of resource hyperfabric_connection with id '%s'", data.Id.ValueString()))
	checkAndSetConnectionIds(data)
	getAndSetConnectionAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	if data.Id.IsNull() {
		var emptyData *ConnectionResourceModel
		resp.Diagnostics.Append(resp.State.Set(ctx, &emptyData)...)
	} else {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
	tflog.Debug(ctx, fmt.Sprintf("End read of resource hyperfabric_connection with id '%s'", data.Id.ValueString()))
}

func (r *ConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Start update of resource: hyperfabric_connection")
	var data *ConnectionResourceModel
	var stateData *ConnectionResourceModel

	// Read Terraform plan and state data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Update of resource hyperfabric_connection with id '%s'", data.Id.ValueString()))

	jsonPayload := getConnectionJsonPayload(ctx, &resp.Diagnostics, data, "update")

	if resp.Diagnostics.HasError() {
		return
	}

	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/connections/%s", data.FabricId.ValueString(), data.ConnectionId.ValueString()), "PUT", jsonPayload)

	if resp.Diagnostics.HasError() {
		return
	}

	getAndSetConnectionAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End update of resource hyperfabric_connection with id '%s'", data.Id.ValueString()))
}

func (r *ConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Start delete of resource: hyperfabric_connection")
	var data *ConnectionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Delete of resource hyperfabric_connection with id '%s'", data.Id.ValueString()))
	checkAndSetConnectionIds(data)
	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/connections/%s", data.FabricId.ValueString(), data.ConnectionId.ValueString()), "DELETE", nil)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("End delete of resource hyperfabric_connection with id '%s'", data.Id.ValueString()))
}

func (r *ConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Start import state of resource: hyperfabric_connection")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	var stateData *ConnectionResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &stateData)...)
	tflog.Debug(ctx, fmt.Sprintf("Import state of resource hyperfabric_connection with id '%s'", stateData.Id.ValueString()))
	tflog.Debug(ctx, "End import of state resource: hyperfabric_connection with id")
}

func getAndSetConnectionAttributes(ctx context.Context, diags *diag.Diagnostics, client *client.Client, data *ConnectionResourceModel) {
	requestData := DoRestRequest(ctx, diags, client, fmt.Sprintf("/api/v1/fabrics/%s/connections/%s", data.FabricId.ValueString(), data.ConnectionId.ValueString()), "GET", nil)
	if diags.HasError() {
		return
	}

	newConnection := *getNewConnectionResourceModelFromData(data)

	if requestData.Data() != nil {
		attributes := requestData.Data().(map[string]interface{})
		for attributeName, attributeValue := range attributes {
			if attributeName == "id" && (data.ConnectionId.IsNull() || data.ConnectionId.IsUnknown() || data.ConnectionId.ValueString() == "" || data.ConnectionId.ValueString() != attributeValue.(string)) {
				newConnection.ConnectionId = basetypes.NewStringValue(attributeValue.(string))
				newConnection.Id = basetypes.NewStringValue(fmt.Sprintf("%s/connections/%s", newConnection.FabricId.ValueString(), newConnection.ConnectionId.ValueString()))
			} else if attributeName == "fabricId" && (data.FabricId.IsNull() || data.FabricId.IsUnknown() || data.FabricId.ValueString() == "" || data.FabricId.ValueString() != attributeValue.(string)) {
				newConnection.FabricId = basetypes.NewStringValue(attributeValue.(string))
				newConnection.Id = basetypes.NewStringValue(fmt.Sprintf("%s/connections/%s", newConnection.FabricId.ValueString(), newConnection.ConnectionId.ValueString()))
			} else if attributeName == "description" {
				newConnection.Description = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "cableType" {
				newConnection.CableType = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "cableLength" {
				newConnection.CableLength = basetypes.NewFloat64Value(attributeValue.(float64))
			} else if attributeName == "pluggable" {
				newConnection.Pluggable = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "local" {
				newConnection.Local = NewLocalRemoteConnectionObject(ctx, attributeValue.(map[string]interface{}))
			} else if attributeName == "remote" {
				newConnection.Remote = NewLocalRemoteConnectionObject(ctx, attributeValue.(map[string]interface{}))
			} else if attributeName == "osType" {
				newConnection.OsType = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "unrecognized" {
				newConnection.Unrecognized = basetypes.NewBoolValue(attributeValue.(bool))
				// } else if attributeName == "metadata" {
				// 	newConnection.Metadata = NewMetadataObject(ctx, attributeValue.(map[string]interface{}))
				// } else if attributeName == "labels" {
				// 	newConnection.Labels = NewSetString(ctx, attributeValue.([]interface{}))
				// } else if attributeName == "annotations" {
				// 	newConnection.Annotations = NewAnnotationsSet(ctx, attributeValue.([]interface{}))
			}
		}
	} else {
		data.Id = basetypes.NewStringNull()
	}
	*data = newConnection
}

func getConnectionJsonPayload(ctx context.Context, diags *diag.Diagnostics, data *ConnectionResourceModel, action string) *gabs.Container {
	payloadMap := map[string]interface{}{}
	payloadList := []map[string]interface{}{}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		payloadMap["description"] = data.Description.ValueString()
	}

	if !data.CableType.IsNull() && !data.CableType.IsUnknown() {
		payloadMap["cableType"] = data.CableType.ValueString()
	}

	if !data.CableLength.IsNull() && !data.CableLength.IsUnknown() {
		payloadMap["cableLength"] = data.CableLength.ValueFloat64()
	}

	if !data.Pluggable.IsNull() && !data.Pluggable.IsUnknown() {
		payloadMap["pluggable"] = data.Pluggable.ValueString()
	}

	if !data.Local.IsNull() && !data.Local.IsUnknown() {
		payloadMap["local"] = getLocalRemoteConnectionJsonPayload(ctx, data.Local)
	}

	if !data.Remote.IsNull() && !data.Remote.IsUnknown() {
		payloadMap["remote"] = getLocalRemoteConnectionJsonPayload(ctx, data.Remote)
	}

	// if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
	// 	payloadMap["labels"] = getSetStringJsonPayload(ctx, data.Labels)
	// }

	// if !data.Annotations.IsNull() && !data.Annotations.IsUnknown() {
	// 	payloadMap["annotations"] = getAnnotationsJsonPayload(ctx, data.Annotations)
	// }

	var payload map[string]interface{}
	if action == "create" {
		payloadList = append(payloadList, payloadMap)
		payload = map[string]interface{}{"connections": payloadList}
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

func checkAndSetConnectionIds(data *ConnectionResourceModel) {
	if strings.Contains(data.Id.ValueString(), "/connections/") {
		if data.FabricId.IsNull() || data.FabricId.IsUnknown() || data.FabricId.ValueString() == "" || data.ConnectionId.IsNull() || data.ConnectionId.IsUnknown() || data.ConnectionId.ValueString() == "" {
			splitId := strings.Split(data.Id.ValueString(), "/connections/")
			data.FabricId = basetypes.NewStringValue(splitId[0])
			data.ConnectionId = basetypes.NewStringValue(splitId[1])
		}
	}
}
