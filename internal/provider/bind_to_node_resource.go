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
	"fmt"
	"strings"

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
var _ resource.Resource = &BindToNodeResource{}
var _ resource.ResourceWithImportState = &BindToNodeResource{}

func NewBindToNodeResource() resource.Resource {
	return &BindToNodeResource{}
}

// BindToNodeResource defines the resource implementation.
type BindToNodeResource struct {
	client *client.Client
}

// BindToNodeResourceModel describes the resource data model.
type BindToNodeResourceModel struct {
	Id       types.String `tfsdk:"id"`
	NodeId   types.String `tfsdk:"node_id"`
	DeviceId types.String `tfsdk:"device_id"`
	// FabricId     types.String `tfsdk:"fabric_id"`
}

type BindToNodeIdentifier struct {
	Id types.String
}

func (r *BindToNodeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Debug(ctx, "Start metadata of resource: hyperfabric_bind_to_node")
	resp.TypeName = req.ProviderTypeName + "_bind_to_node"
	tflog.Debug(ctx, "End metadata of resource: hyperfabric_bind_to_node")
}

func (r *BindToNodeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Debug(ctx, "Start schema of resource: hyperfabric_bind_to_node")
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Bind to Node resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "`id` defines the unique identifier of a node bound to a node.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"node_id": schema.StringAttribute{
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"device_id": schema.StringAttribute{
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
	tflog.Debug(ctx, "End schema of resource: hyperfabric_bind_to_node")
}

func (r *BindToNodeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of resource: hyperfabric_bind_to_node")
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
	tflog.Debug(ctx, "End configure of resource: hyperfabric_bind_to_node")
}

func (r *BindToNodeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Start create of resource: hyperfabric_bind_to_node")

	var data *BindToNodeResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Create of resource hyperfabric_bind_to_node with NodeId '%s' and SwichId '%s'", data.NodeId.ValueString(), data.DeviceId.ValueString()))

	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/chassis/%s/switches/%s", data.NodeId.ValueString(), data.DeviceId.ValueString()), "PUT", nil)
	if resp.Diagnostics.HasError() {
		return
	}

	// TO FIX
	// nodeContainer, err := container.ArrayElement(0, "chassis")
	// if err != nil {
	// 	return
	// }

	data.Id = basetypes.NewStringValue(fmt.Sprintf("%s/%s", data.NodeId.ValueString(), data.DeviceId.ValueString()))
	// getAndSetBindToNodeAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End create of resource hyperfabric_bind_to_node with NodeId '%s' and DeviceId '%s'", data.NodeId.ValueString(), data.DeviceId.ValueString()))
}

func (r *BindToNodeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Start read of resource: hyperfabric_bind_to_node")
	var data *BindToNodeResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Read of resource hyperfabric_bind_to_node with id '%s'", data.Id.ValueString()))

	getAndSetBindToNodeAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	if data.Id.IsNull() {
		var emptyData *BindToNodeResourceModel
		resp.Diagnostics.Append(resp.State.Set(ctx, &emptyData)...)
	} else {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}

	tflog.Debug(ctx, fmt.Sprintf("End read of resource hyperfabric_bind_to_node with id '%s'", data.Id.ValueString()))
}

func (r *BindToNodeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Start update of resource: hyperfabric_bind_to_node")
	var data *BindToNodeResourceModel
	var stateData *BindToNodeResourceModel

	// Read Terraform plan and state data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Update of resource hyperfabric_bind_to_node with id '%s'", data.Id.ValueString()))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End update of resource hyperfabric_bind_to_node with id '%s'", data.Id.ValueString()))
}

func (r *BindToNodeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Start delete of resource: hyperfabric_bind_to_node")
	var data *BindToNodeResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Delete of resource hyperfabric_bind_to_node with id '%s'", data.Id.ValueString()))
	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/chassis/%s/switches", data.NodeId.ValueString()), "DELETE", nil)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("End delete of resource hyperfabric_bind_to_node with id '%s'", data.Id.ValueString()))
}

func (r *BindToNodeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Start import state of resource: hyperfabric_bind_to_node")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	var stateData *BindToNodeResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &stateData)...)
	tflog.Debug(ctx, fmt.Sprintf("Import state of resource hyperfabric_bind_to_node with id '%s'", stateData.Id.ValueString()))
	tflog.Debug(ctx, "End import of state resource: hyperfabric_bind_to_node")
}

func getAndSetBindToNodeAttributes(ctx context.Context, diags *diag.Diagnostics, client *client.Client, data *BindToNodeResourceModel) {

	idSplit := strings.Split(data.Id.ValueString(), "/")
	data.NodeId = basetypes.NewStringValue(idSplit[0])

	requestData := DoRestRequest(ctx, diags, client, fmt.Sprintf("/api/v1/chassis/%s", data.NodeId.ValueString()), "GET", nil)
	if diags.HasError() {
		return
	}

	data.DeviceId = basetypes.NewStringValue("")

	if requestData.Data() != nil {
		attributes := requestData.Data().(map[string]interface{})
		for attributeName, attributeValue := range attributes {
			if attributeName == "id" {
				data.Id = basetypes.NewStringValue(attributeValue.(string))
				// setBindToNodeParentId(ctx, attributeValue.(string), data)
			} else if attributeName == "switchId" {
				data.DeviceId = basetypes.NewStringValue(attributeValue.(string))
			}
		}
	} else {
		data.Id = basetypes.NewStringNull()
	}
}
