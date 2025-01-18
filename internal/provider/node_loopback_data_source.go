// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

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
var _ datasource.DataSource = &NodeLoopbackDataSource{}

func NewNodeLoopbackDataSource() datasource.DataSource {
	return &NodeLoopbackDataSource{}
}

// NodeLoopbackDataSource defines the data source implementation.
type NodeLoopbackDataSource struct {
	client *client.Client
}

func (d *NodeLoopbackDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	tflog.Debug(ctx, "Start metadata of datasource: hyperfabric_node_loopback")
	resp.TypeName = req.ProviderTypeName + "_node_loopback"
	tflog.Debug(ctx, "End metadata of datasource: hyperfabric_node_loopback")
}

func (d *NodeLoopbackDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	tflog.Debug(ctx, "Start schema of datasource: hyperfabric_node_loopback")
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Node Loopback data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "`id` defines the unique identifier of the Loopback of a Node in a Fabric.",
				Computed:            true,
			},
			"loopback_id": schema.StringAttribute{
				MarkdownDescription: "`loopback_id` defines the unique identifier of a Loopback of a Node.",
				Computed:            true,
			},
			"node_id": schema.StringAttribute{
				MarkdownDescription: "`node_id` defines the unique identifier of a Node in a Fabric.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Loopback of the Node.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description is a user defined field to store notes about the Loopback of the Node.",
				Computed:            true,
			},
			// "enabled": schema.BoolAttribute{
			// 	MarkdownDescription: "The enabled admin state of the Loopback of the Node.",
			// 	Computed:            true,
			// },
			"ipv4_address": schema.StringAttribute{
				MarkdownDescription: "The IPv4 address configured on the Loopback of the Node.",
				Computed:            true,
			},
			"ipv6_address": schema.StringAttribute{
				MarkdownDescription: "The IPv6 address configured on the Loopback of the Node.",
				Computed:            true,
			},
			"vrf_id": schema.StringAttribute{
				MarkdownDescription: "The `vrf_id` of a VRF to associate with the Loopback of the Node. Required when the Loopback roles include `ROUTED_PORT`.",
				Computed:            true,
			},
			"metadata":    getMetadataSchemaAttribute(),
			"labels":      getLabelsDataSourceSchemaAttribute(),
			"annotations": getAnnotationsDataSourceSchemaAttribute(),
		},
	}
	tflog.Debug(ctx, "End schema of datasource: hyperfabric_node_loopback")
}

func (d *NodeLoopbackDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of datasource: hyperfabric_node_loopback")
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
	tflog.Debug(ctx, "End configure of datasource: hyperfabric_node_loopback")
}

func (d *NodeLoopbackDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Start read of datasource: hyperfabric_node_loopback")
	var data *NodeLoopbackResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create a copy of the Id for when not found during getAndSetNodeLoopbackAttributes
	cachedId := data.Id.ValueString()
	if cachedId == "" && data.Name.ValueString() != "" {
		data.LoopbackId = data.Name
	}

	tflog.Debug(ctx, fmt.Sprintf("Read of datasource hyperfabric_node_loopback with id '%s'", data.Id.ValueString()))

	getAndSetNodeLoopbackAttributes(ctx, &resp.Diagnostics, d.client, data)

	if data.Id.IsNull() {
		resp.Diagnostics.AddError(
			"Failed to read hyperfabric_node_loopback data source",
			fmt.Sprintf("The hyperfabric_node_loopback data source with id '%s' has not been found", cachedId),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End read of datasource hyperfabric_node_loopback with id '%s'", data.Id.ValueString()))
}
