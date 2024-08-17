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

	"github.com/cisco-open/terraform-provider-hyperfabric/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &VniDataSource{}

func NewVniDataSource() datasource.DataSource {
	return &VniDataSource{}
}

// VniDataSource defines the data source implementation.
type VniDataSource struct {
	client *client.Client
}

func (r *VniDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	tflog.Debug(ctx, "Start metadata of datasource: hyperfabric_vni")
	resp.TypeName = req.ProviderTypeName + "_vni"
	tflog.Debug(ctx, "End metadata of datasource: hyperfabric_vni")
}

func (r *VniDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	tflog.Debug(ctx, "Start schema of datasource: hyperfabric_vni")
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "VNI data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "`id` defines the unique identifier of a VNI in a Fabric.",
			},
			"vni_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "`vni_id` defines the unique identifier of a VNI.",
			},
			"fabric_id": schema.StringAttribute{
				MarkdownDescription: "`fabric_id` defines the unique identifier of a Fabric.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the VNI.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description is a user defined field to store notes about the VNI.",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "The enabled state of the VNI.",
				Computed:            true,
			},
			"is_default": schema.BoolAttribute{
				MarkdownDescription: "The flag that denote if the VNI is the default VNI or not.",
				Computed:            true,
			},
			"vrf_id": schema.StringAttribute{
				MarkdownDescription: "The Id of the VRF associated with the VNI.",
				Computed:            true,
			},
			"mtu": schema.Float64Attribute{
				MarkdownDescription: "The MTU of the SVI of the VNI.",
				Computed:            true,
			},
			"vni": schema.Float64Attribute{
				MarkdownDescription: "The VXLAN Network Identifier (VNI) used for the VNI.",
				Computed:            true,
			},
			"members":     getMembersDataSourceSchemaAttribute(),
			"svi":         getSviDataSourceSchemaAttribute(),
			"metadata":    getMetadataSchemaAttribute(),
			"labels":      getLabelsDataSourceSchemaAttribute(),
			"annotations": getAnnotationsDataSourceSchemaAttribute(),
		},
	}
	tflog.Debug(ctx, "End schema of datasource: hyperfabric_vni")
}

func (r *VniDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of datasource: hyperfabric_vni")
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
	tflog.Debug(ctx, "End configure of datasource: hyperfabric_vni")
}

func (r *VniDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Start read of datasource: hyperfabric_vni")
	var data *VniResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create a copy of the Id for when not found during getAndSetNodeAttributes
	cachedId := data.Id.ValueString()
	if cachedId == "" && data.Name.ValueString() != "" {
		data.VniId = data.Name
	}

	tflog.Debug(ctx, fmt.Sprintf("Read of datasource hyperfabric_vni with id '%s'", data.Id.ValueString()))

	getAndSetVniAttributes(ctx, &resp.Diagnostics, r.client, data)

	if data.Id.IsNull() {
		resp.Diagnostics.AddError(
			"Failed to read hyperfabric_node data source",
			fmt.Sprintf("The hyperfabric_node data source with id '%s' has not been found", cachedId),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End read of datasource hyperfabric_vni with id '%s'", data.Id.ValueString()))
}
