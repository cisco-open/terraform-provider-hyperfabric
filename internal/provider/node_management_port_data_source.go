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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &NodeManagementPortDataSource{}

func NewNodeManagementPortDataSource() datasource.DataSource {
	return &NodeManagementPortDataSource{}
}

// NodeManagementPortDataSource defines the data source implementation.
type NodeManagementPortDataSource struct {
	client *client.Client
}

func (d *NodeManagementPortDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	tflog.Debug(ctx, "Start metadata of datasource: hyperfabric_node_management_port")
	resp.TypeName = req.ProviderTypeName + "_node_management_port"
	tflog.Debug(ctx, "End metadata of datasource: hyperfabric_node_management_port")
}

func (d *NodeManagementPortDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	tflog.Debug(ctx, "Start schema of datasource: hyperfabric_node_management_port")
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Node Management Port data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "`id` defines the unique identifier of the Management Port of a Node in a Fabric.",
				Computed:            true,
			},
			"node_management_port_id": schema.StringAttribute{
				MarkdownDescription: "`node_management_port_id` defines the unique identifier of a Management Port of a Node.",
				Computed:            true,
			},
			"node_id": schema.StringAttribute{
				MarkdownDescription: "`node_id` defines the unique identifier of a Node in a Fabric.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Management Port of the Node.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description is a user defined field to store notes about the Management Port of the Node.",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "The enabled state of the Management Port of the Node.",
				Computed:            true,
			},
			"cloud_urls": getCloudUrlsDataSourceSchemaAttribute(),
			"ipv4_config_type": schema.StringAttribute{
				MarkdownDescription: "Determines if the IPv4 configuration is static or from DHCP",
				Computed:            true,
			},
			"ipv4_address": schema.StringAttribute{
				MarkdownDescription: "The IPv4 address for the Management Port of the Node.",
				Computed:            true,
			},
			"ipv4_gateway": schema.StringAttribute{
				MarkdownDescription: "The IPv4 gateway address for the Management Port of the Node.",
				Computed:            true,
			},
			"ipv6_config_type": schema.StringAttribute{
				MarkdownDescription: "Determines if the IPv6 configuration is static or from DHCP",
				Computed:            true,
			},
			"ipv6_address": schema.StringAttribute{
				MarkdownDescription: "The IPv6 address for the Management Port of the Node.",
				Computed:            true,
			},
			"ipv6_gateway": schema.StringAttribute{
				MarkdownDescription: "The IPv6 gateway address for the Management Port of the Node.",
				Computed:            true,
			},
			"dns_addresses": getDnsAddressesDataSourceSchemaAttribute(),
			"ntp_addresses": getNtpAddressesDataSourceSchemaAttribute(),
			"no_proxy":      getNoProxyDataSourceSchemaAttribute(),
			"proxy_address": schema.StringAttribute{
				MarkdownDescription: "The URL for a configured HTTPs proxy for the Node.",
				Computed:            true,
			},
			"proxy_username": schema.StringAttribute{
				MarkdownDescription: "A username to be used to authenticate to the proxy.",
				Computed:            true,
			},
			"proxy_password": schema.StringAttribute{
				MarkdownDescription: "A password to be used to authenticate to the proxy.",
				Computed:            true,
			},
			"proxy_credential_id": schema.StringAttribute{
				MarkdownDescription: "`proxy_credential_id` defines the unique identifier of a set of credentials for the proxy.",
				Computed:            true,
			},
			"config_origin": schema.StringAttribute{
				MarkdownDescription: "The source of the configuration, either from the cloud or the device.",
				Computed:            true,
			},
			"connected_state": schema.StringAttribute{
				MarkdownDescription: "The connected state denoting if the port has ever successfully connected to the service.",
				Computed:            true,
			},
			"metadata": getMetadataSchemaAttribute(),
			// "labels":      getLabelsDataSourceSchemaAttribute(),
			// "annotations": getAnnotationsDataSourceSchemaAttribute(),
		},
	}
	tflog.Debug(ctx, "End schema of datasource: hyperfabric_node_management_port")
}

func (d *NodeManagementPortDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of datasource: hyperfabric_node_management_port")
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
	tflog.Debug(ctx, "End configure of datasource: hyperfabric_node_management_port")
}

func (d *NodeManagementPortDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Start read of datasource: hyperfabric_node_management_port")
	var data *NodeManagementPortResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Name.IsNull() || data.Name.IsUnknown() {
		data.Name = basetypes.NewStringValue("eth0")
	}

	// Create a copy of the Id for when not found during getAndSetNodeManagementPortAttributes
	cachedId := data.Id.ValueString()
	if cachedId == "" && data.Name.ValueString() != "" {
		data.NodeManagementPortId = data.Name
	}

	tflog.Debug(ctx, fmt.Sprintf("Read of datasource hyperfabric_node_management_port with id '%s'", data.Id.ValueString()))

	getAndSetNodeManagementPortAttributes(ctx, &resp.Diagnostics, d.client, data)

	if data.Id.IsNull() {
		resp.Diagnostics.AddError(
			"Failed to read hyperfabric_node_management_port data source",
			fmt.Sprintf("The hyperfabric_node_management_port data source with id '%s' has not been found", cachedId),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End read of datasource hyperfabric_node_management_port with id '%s'", data.Id.ValueString()))
}
