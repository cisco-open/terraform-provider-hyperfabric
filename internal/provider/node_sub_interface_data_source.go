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
var _ datasource.DataSource = &NodeSubInterfaceDataSource{}

func NewNodeSubInterfaceDataSource() datasource.DataSource {
	return &NodeSubInterfaceDataSource{}
}

// NodeSubInterfaceDataSource defines the data source implementation.
type NodeSubInterfaceDataSource struct {
	client *client.Client
}

func (d *NodeSubInterfaceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	tflog.Debug(ctx, "Start metadata of datasource: hyperfabric_node_sub_interface")
	resp.TypeName = req.ProviderTypeName + "_node_sub_interface"
	tflog.Debug(ctx, "End metadata of datasource: hyperfabric_node_sub_interface")
}

func (d *NodeSubInterfaceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	tflog.Debug(ctx, "Start schema of datasource: hyperfabric_node_sub_interface")
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Node Sub-Interface data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "`id` defines the unique identifier of the Sub-Interface of a Node in a Fabric.",
				Computed:            true,
			},
			"sub_interface_id": schema.StringAttribute{
				MarkdownDescription: "`sub_interface_id` defines the unique identifier of a Sub-Interface of a Node in a Fabric.",
				Computed:            true,
			},
			"node_id": schema.StringAttribute{
				MarkdownDescription: "`node_id` defines the unique identifier of a Node in a Fabric.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Sub-Interface of the Node.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description is a user defined field to store notes about the Sub-Interface of the Node.",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "The enabled admin state of the Sub-Interface of the Node.",
				Computed:            true,
			},
			"ipv4_addresses": getSubInterfaceIpv4AddressesDataSourceSchemaAttribute(),
			"ipv6_addresses": getSubInterfaceIpv6AddressesDataSourceSchemaAttribute(),
			"vlan_id": schema.StringAttribute{
				MarkdownDescription: "The VLAN ID to use as encapsulation for the Sub-Interface of the Node.",
				Computed:            true,
			},
			"vrf_id": schema.StringAttribute{
				MarkdownDescription: "The `vrf_id` of a VRF to associate with the Sub-Interface of the Node.",
				Computed:            true,
			},
			"parent": schema.StringAttribute{
				MarkdownDescription: "The name of the `parent` Port of the Sub-Interface of the Node.",
				Computed:            true,
			},
			"metadata":    getMetadataSchemaAttribute(),
			"labels":      getLabelsDataSourceSchemaAttribute(),
			"annotations": getAnnotationsDataSourceSchemaAttribute(),
		},
	}
	tflog.Debug(ctx, "End schema of datasource: hyperfabric_node_sub_interface")
}

func (d *NodeSubInterfaceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of datasource: hyperfabric_node_sub_interface")
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
	tflog.Debug(ctx, "End configure of datasource: hyperfabric_node_sub_interface")
}

func (d *NodeSubInterfaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Start read of datasource: hyperfabric_node_sub_interface")
	var data *NodeSubInterfaceResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create a copy of the Id for when not found during getAndSetNodeSubInterfaceAttributes
	cachedId := data.Id.ValueString()
	if cachedId == "" && data.Name.ValueString() != "" {
		data.SubInterfaceId = data.Name
	}

	tflog.Debug(ctx, fmt.Sprintf("Read of datasource hyperfabric_node_sub_interface with id '%s'", data.Id.ValueString()))

	getAndSetNodeSubInterfaceAttributes(ctx, &resp.Diagnostics, d.client, data)

	if data.Id.IsNull() {
		resp.Diagnostics.AddError(
			"Failed to read hyperfabric_node_sub_interface data source",
			fmt.Sprintf("The hyperfabric_node_sub_interface data source with id '%s' has not been found", cachedId),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End read of datasource hyperfabric_node_sub_interface with id '%s'", data.Id.ValueString()))
}
