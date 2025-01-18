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
var _ datasource.DataSource = &NodePortDataSource{}

func NewNodePortDataSource() datasource.DataSource {
	return &NodePortDataSource{}
}

// NodePortDataSource defines the data source implementation.
type NodePortDataSource struct {
	client *client.Client
}

func (d *NodePortDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	tflog.Debug(ctx, "Start metadata of datasource: hyperfabric_node_port")
	resp.TypeName = req.ProviderTypeName + "_node_port"
	tflog.Debug(ctx, "End metadata of datasource: hyperfabric_node_port")
}

func (d *NodePortDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	tflog.Debug(ctx, "Start schema of datasource: hyperfabric_node_port")
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Node Port data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "`id` defines the unique identifier of the Port of a Node in a Fabric.",
				Computed:            true,
			},
			"port_id": schema.StringAttribute{
				MarkdownDescription: "`port_id` defines the unique identifier of a Port of a Node.",
				Computed:            true,
			},
			"node_id": schema.StringAttribute{
				MarkdownDescription: "`node_id` defines the unique identifier of a Node in a Fabric.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Port of the Node.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description is a user defined field to store notes about the Port of the Node.",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "The enabled admin state of the Port of the Node.",
				Computed:            true,
			},
			// "breakout": schema.BoolAttribute{
			// 	MarkdownDescription: "The breakout state of the Port of the Node.",
			// 	Computed:            true,
			// },
			// "breakout_index": schema.Float64Attribute{
			// 	MarkdownDescription: "The index of the sub-port on the breakout Port.",
			// 	Computed:            true,
			// },
			"index": schema.Float64Attribute{
				MarkdownDescription: "The index number of the Port of the Node.",
				Computed:            true,
			},
			"ipv4_addresses": getIpv4AddressesDataSourceSchemaAttribute(),
			"ipv6_addresses": getIpv6AddressesDataSourceSchemaAttribute(),
			"linecard": schema.Float64Attribute{
				MarkdownDescription: "The linecard index number of the Port of the Node.",
				Computed:            true,
			},
			"prevent_forwarding": schema.BoolAttribute{
				MarkdownDescription: "Prevent traffic from being forwarded by the Port.",
				Optional:            true,
				Computed:            true,
			},
			"lldp_host": schema.StringAttribute{
				MarkdownDescription: "The name of host reported by LLDP connected to the Port of the Node.",
				Computed:            true,
			},
			"lldp_info": schema.StringAttribute{
				MarkdownDescription: "The info about the host reported by LLDP connected to the Port of the Node.",
				Computed:            true,
			},
			"lldp_port": schema.StringAttribute{
				MarkdownDescription: "The port of host reported by LLDP connected to the Port of the Node.",
				Computed:            true,
			},
			"max_speed": schema.StringAttribute{
				MarkdownDescription: "The maximum speed of the Port of the Node.",
				Computed:            true,
			},
			"mtu": schema.Float64Attribute{
				MarkdownDescription: "The MTU of the Port of the Node.",
				Computed:            true,
			},
			"roles": getPortRolesDataSourceSchemaAttribute(),
			"speed": schema.StringAttribute{
				MarkdownDescription: "The configured speed of the Port of the Node.",
				Computed:            true,
			},
			"sub_interfaces_count": schema.Float64Attribute{
				MarkdownDescription: "The number of sub-interfaces of the Port of the Node.",
				Computed:            true,
			},
			"vlan_ids": getVlanIdsSchemaAttribute(),
			"vnis":     getVnisSchemaAttribute(),
			"vrf_id": schema.StringAttribute{
				MarkdownDescription: "The `vrf_id` of a VRF to associate with the Port of the Node. Required when the Port roles include `ROUTED_PORT`.",
				Computed:            true,
			},
			"metadata":    getMetadataSchemaAttribute(),
			"labels":      getLabelsDataSourceSchemaAttribute(),
			"annotations": getAnnotationsDataSourceSchemaAttribute(),
		},
	}
	tflog.Debug(ctx, "End schema of datasource: hyperfabric_node_port")
}

func (d *NodePortDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of datasource: hyperfabric_node_port")
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
	tflog.Debug(ctx, "End configure of datasource: hyperfabric_node_port")
}

func (d *NodePortDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Start read of datasource: hyperfabric_node_port")
	var data *NodePortResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create a copy of the Id for when not found during getAndSetNodePortAttributes
	cachedId := data.Id.ValueString()
	if cachedId == "" && data.Name.ValueString() != "" {
		data.PortId = data.Name
	}

	tflog.Debug(ctx, fmt.Sprintf("Read of datasource hyperfabric_node_port with id '%s'", data.Id.ValueString()))

	getAndSetNodePortAttributes(ctx, &resp.Diagnostics, d.client, data)

	if data.Id.IsNull() {
		resp.Diagnostics.AddError(
			"Failed to read hyperfabric_node_port data source",
			fmt.Sprintf("The hyperfabric_node_port data source with id '%s' has not been found", cachedId),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End read of datasource hyperfabric_node_port with id '%s'", data.Id.ValueString()))
}
