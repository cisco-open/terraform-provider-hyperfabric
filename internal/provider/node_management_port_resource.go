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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &NodeManagementPortResource{}
var _ resource.ResourceWithImportState = &NodeManagementPortResource{}

func NewNodeManagementPortResource() resource.Resource {
	return &NodeManagementPortResource{}
}

// NodeManagementPortResource defines the resource implementation.
type NodeManagementPortResource struct {
	client *client.Client
}

// NodeManagementPortResourceModel describes the resource data model.
type NodeManagementPortResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	NodeId               types.String `tfsdk:"node_id"`
	NodeManagementPortId types.String `tfsdk:"node_management_port_id"`
	// FabricId             types.String `tfsdk:"fabric_id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	Enabled           types.Bool   `tfsdk:"enabled"`
	CloudUrls         types.Set    `tfsdk:"cloud_urls"`
	Ipv4ConfigType    types.String `tfsdk:"ipv4_config_type"`
	Ipv4Address       types.String `tfsdk:"ipv4_address"`
	Ipv4Gateway       types.String `tfsdk:"ipv4_gateway"`
	Ipv6ConfigType    types.String `tfsdk:"ipv6_config_type"`
	Ipv6Address       types.String `tfsdk:"ipv6_address"`
	Ipv6Gateway       types.String `tfsdk:"ipv6_gateway"`
	DnsAddresses      types.Set    `tfsdk:"dns_addresses"`
	NtpAddresses      types.Set    `tfsdk:"ntp_addresses"`
	NoProxy           types.Set    `tfsdk:"no_proxy"`
	ProxyAddress      types.String `tfsdk:"proxy_address"`
	ProxyCredentialId types.String `tfsdk:"proxy_credential_id"`
	ProxyUsername     types.String `tfsdk:"proxy_username"`
	ProxyPassword     types.String `tfsdk:"proxy_password"`
	// SetProxyPassword  types.Bool   `tfsdk:"set_proxy_password"`
	ConfigOrigin   types.String `tfsdk:"config_origin"`
	ConnectedState types.String `tfsdk:"connected_state"`
	Metadata       types.Object `tfsdk:"metadata"`
	// Labels            types.Set    `tfsdk:"labels"`
	// Annotations       types.Set    `tfsdk:"annotations"`
}

func getEmptyNodeManagementPortResourceModel() *NodeManagementPortResourceModel {
	return &NodeManagementPortResourceModel{
		Id:                   basetypes.NewStringNull(),
		NodeId:               basetypes.NewStringNull(),
		NodeManagementPortId: basetypes.NewStringNull(),
		// FabricId:             basetypes.NewStringNull(),
		Name:              basetypes.NewStringNull(),
		Description:       basetypes.NewStringNull(),
		Enabled:           basetypes.NewBoolValue(false),
		CloudUrls:         basetypes.NewSetNull(CloudUrlsResourceModelAttributeType()),
		Ipv4ConfigType:    basetypes.NewStringNull(),
		Ipv4Address:       basetypes.NewStringNull(),
		Ipv4Gateway:       basetypes.NewStringNull(),
		Ipv6ConfigType:    basetypes.NewStringNull(),
		Ipv6Address:       basetypes.NewStringNull(),
		Ipv6Gateway:       basetypes.NewStringNull(),
		DnsAddresses:      basetypes.NewSetNull(DnsAddressesResourceModelAttributeType()),
		NtpAddresses:      basetypes.NewSetNull(NtpAddressesResourceModelAttributeType()),
		NoProxy:           basetypes.NewSetNull(NoProxyResourceModelAttributeType()),
		ProxyAddress:      basetypes.NewStringNull(),
		ProxyCredentialId: basetypes.NewStringNull(),
		ProxyUsername:     basetypes.NewStringNull(),
		ProxyPassword:     basetypes.NewStringNull(),
		// SetProxyPassword:  basetypes.NewBoolValue(false),
		ConfigOrigin:   basetypes.NewStringNull(),
		ConnectedState: basetypes.NewStringNull(),
		Metadata:       basetypes.NewObjectNull(MetadataResourceModelAttributeType()),
		// Labels:            basetypes.NewSetNull(SetStringResourceModelAttributeType()),
		// Annotations:       basetypes.NewSetNull(AnnotationResourceModelAttributeType()),
	}
}

func CloudUrlsResourceModelAttributeType() attr.Type {
	return types.StringType
}

func DnsAddressesResourceModelAttributeType() attr.Type {
	return types.StringType
}

func NtpAddressesResourceModelAttributeType() attr.Type {
	return types.StringType
}

func NoProxyResourceModelAttributeType() attr.Type {
	return types.StringType
}

type NodeManagementPortIdentifier struct {
	Id types.String
}

func (r *NodeManagementPortResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Debug(ctx, "Start metadata of resource: hyperfabric_node_management_port")
	resp.TypeName = req.ProviderTypeName + "_node_management_port"
	tflog.Debug(ctx, "End metadata of resource: hyperfabric_node_management_port")
}

func (r *NodeManagementPortResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Debug(ctx, "Start schema of resource: hyperfabric_node_management_port")
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Node Management Port resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "`id` defines the unique identifier of the Management Port of a Node in a Fabric.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"node_management_port_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "`node_management_port_id` defines the unique identifier of a Management Port of a Node.",
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
			// "fabric_id": schema.StringAttribute{
			// 	MarkdownDescription: "`fabric_id` defines the unique identifier of a Fabric.",
			// 	Required:            true,
			// 	PlanModifiers: []planmodifier.String{
			// 		stringplanmodifier.RequiresReplace(),
			// 	},
			// },
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Management Port of the Node.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("eth0"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description is a user defined field to store notes about the Management Port of the Node.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "The enabled state of the Management Port of the Node.",
				Computed:            true,
				// Default:             booldefault.StaticBool(true),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"cloud_urls": getCloudUrlsSchemaAttribute(),
			"ipv4_config_type": schema.StringAttribute{
				MarkdownDescription: "Determines if the IPv4 configuration is static or from DHCP",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
				Default: stringdefault.StaticString("CONFIG_TYPE_DHCP"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"CONFIG_TYPE_UNSPECIFIED", "CONFIG_TYPE_STATIC", "CONFIG_TYPE_DHCP"}...),
				},
			},
			"ipv4_address": schema.StringAttribute{
				MarkdownDescription: "The IPv4 address for the Management Port of the Node.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
				Validators: []validator.String{
					// Validate this attribute must be configured with other_attr.
					stringvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("ipv4_gateway"),
						path.MatchRoot("dns_addresses"),
					}...),
				},
			},
			"ipv4_gateway": schema.StringAttribute{
				MarkdownDescription: "The IPv4 gateway address for the Management Port of the Node.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
				Validators: []validator.String{
					// Validate this attribute must be configured with other_attr.
					stringvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("ipv4_address"),
						path.MatchRoot("dns_addresses"),
					}...),
				},
			},
			"ipv6_config_type": schema.StringAttribute{
				MarkdownDescription: "Determines if the IPv6 configuration is static or from DHCP",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
				Default: stringdefault.StaticString("CONFIG_TYPE_DHCP"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"CONFIG_TYPE_UNSPECIFIED", "CONFIG_TYPE_STATIC", "CONFIG_TYPE_DHCP"}...),
				},
			},
			"ipv6_address": schema.StringAttribute{
				MarkdownDescription: "The IPv6 address for the Management Port of the Node.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
				Validators: []validator.String{
					// Validate this attribute must be configured with other_attr.
					stringvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("ipv6_gateway"),
						path.MatchRoot("dns_addresses"),
					}...),
				},
			},
			"ipv6_gateway": schema.StringAttribute{
				MarkdownDescription: "The IPv6 gateway address for the Management Port of the Node.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
				Validators: []validator.String{
					// Validate this attribute must be configured with other_attr.
					stringvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("ipv6_address"),
						path.MatchRoot("dns_addresses"),
					}...),
				},
			},
			"dns_addresses": getDnsAddressesSchemaAttribute(),
			"ntp_addresses": getNtpAddressesSchemaAttribute(),
			"no_proxy":      getNoProxySchemaAttribute(),

			"proxy_address": schema.StringAttribute{
				MarkdownDescription: "The URL for a configured HTTPs proxy for the Node.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"proxy_username": schema.StringAttribute{
				MarkdownDescription: "A username to be used to authenticate to the proxy.",
				Optional:            true,
				// Computed:            true,
				PlanModifiers: []planmodifier.String{
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
				Validators: []validator.String{
					// Validate this attribute must be configured with other_attr.
					stringvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("proxy_password"),
					}...),
				},
			},
			"proxy_password": schema.StringAttribute{
				MarkdownDescription: "A password to be used to authenticate to the proxy.",
				Optional:            true,
				// Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"proxy_credential_id": schema.StringAttribute{
				MarkdownDescription: "`proxy_credential_id` defines the unique identifier of a set of credentials for the proxy.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"config_origin": schema.StringAttribute{
				MarkdownDescription: "The source of the configuration, either from the cloud or the device.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					// SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
				// Default: stringdefault.StaticString("CONFIG_ORIGIN_CLOUD"),
				// Validators: []validator.String{
				// 	stringvalidator.OneOf([]string{"CONFIG_ORIGIN_UNSPECIFIED", "CONFIG_ORIGIN_CLOUD", "CONFIG_ORIGIN_DEVICE"}...),
				// },
			},
			"connected_state": schema.StringAttribute{
				MarkdownDescription: "The connected state denoting if the port has ever successfully connected to the service.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"metadata": getMetadataSchemaAttribute(),
			// "labels":      getLabelsSchemaAttribute(),
			// "annotations": getAnnotationsSchemaAttribute(),
		},
	}
	tflog.Debug(ctx, "End schema of resource: hyperfabric_node_management_port")
}

func getCloudUrlsSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A set of Cloud URLs used by a Node.`,
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.UseStateForUnknown(),
		},
		ElementType: types.StringType,
	}
}

func getCloudUrlsDataSourceSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A set of Cloud URLs used by a Node.`,
		Computed:            true,
		ElementType:         types.StringType,
	}
}

func getDnsAddressesSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A set of DNS IP addresses used by a Node.`,
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.UseStateForUnknown(),
		},
		ElementType: types.StringType,
	}
}

func getDnsAddressesDataSourceSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A set of DNS IP addresses used by a Node.`,
		Computed:            true,
		ElementType:         types.StringType,
	}
}

func getNtpAddressesSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A set of NTP Server IP addresses used by a Node.`,
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.UseStateForUnknown(),
		},
		ElementType: types.StringType,
	}
}

func getNtpAddressesDataSourceSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A set of NTP Server IP addresses used by a Node.`,
		Computed:            true,
		ElementType:         types.StringType,
	}
}

func getNoProxySchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A list of IP addresses or domain names that should not be proxied.`,
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.UseStateForUnknown(),
		},
		ElementType: types.StringType,
	}
}

func getNoProxyDataSourceSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A list of IP addresses or domain names that should not be proxied.`,
		Computed:            true,
		ElementType:         types.StringType,
	}
}

func (r *NodeManagementPortResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of resource: hyperfabric_node_management_port")
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
	tflog.Debug(ctx, "End configure of resource: hyperfabric_node_management_port")
}

func (r *NodeManagementPortResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Start create of resource: hyperfabric_node_management_port")

	var data *NodeManagementPortResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Create of resource hyperfabric_node_management_port with name '%s'", data.Name.ValueString()))

	jsonPayload := getNodeManagementPortJsonPayload(ctx, &resp.Diagnostics, data, "create")
	if resp.Diagnostics.HasError() {
		return
	}

	container := DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/managementPorts", data.NodeId.ValueString()), "POST", jsonPayload)
	if resp.Diagnostics.HasError() {
		return
	}

	managementPortContainer, err := container.ArrayElement(0, "ports")
	if err != nil {
		return
	}

	managementPortId := StripQuotes(managementPortContainer.Search("id").String())
	if managementPortId != "" {
		data.Id = basetypes.NewStringValue(fmt.Sprintf("%s/managementPorts/%s", data.NodeId.ValueString(), managementPortId))
		data.NodeManagementPortId = basetypes.NewStringValue(managementPortId)
		getAndSetNodeManagementPortAttributes(ctx, &resp.Diagnostics, r.client, data)
	} else {
		data.Id = basetypes.NewStringNull()
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End create of resource hyperfabric_node_management_port with id '%s'", data.Id.ValueString()))
}

func (r *NodeManagementPortResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Start read of resource: hyperfabric_node_management_port")
	var data *NodeManagementPortResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Read of resource hyperfabric_node_management_port with id '%s'", data.Id.ValueString()))
	checkAndSetNodeManagementPortIds(data)
	getAndSetNodeManagementPortAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	if data.Id.IsNull() {
		var emptyData *NodeManagementPortResourceModel
		resp.Diagnostics.Append(resp.State.Set(ctx, &emptyData)...)
	} else {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
	tflog.Debug(ctx, fmt.Sprintf("End read of resource hyperfabric_node_management_port with id '%s'", data.Id.ValueString()))
}

func (r *NodeManagementPortResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Start update of resource: hyperfabric_node_management_port")
	var data *NodeManagementPortResourceModel
	var stateData *NodeManagementPortResourceModel

	// Read Terraform plan and state data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Update of resource hyperfabric_node_management_port with id '%s'", data.Id.ValueString()))

	jsonPayload := getNodeManagementPortJsonPayload(ctx, &resp.Diagnostics, data, "update")

	if resp.Diagnostics.HasError() {
		return
	}

	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/managementPorts/%s", data.NodeId.ValueString(), data.NodeManagementPortId.ValueString()), "PUT", jsonPayload)

	if resp.Diagnostics.HasError() {
		return
	}

	getAndSetNodeManagementPortAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End update of resource hyperfabric_node_management_port with id '%s'", data.Id.ValueString()))
}

func (r *NodeManagementPortResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Start delete of resource: hyperfabric_node_management_port")
	var data *NodeManagementPortResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Delete of resource hyperfabric_node_management_port with id '%s'", data.Id.ValueString()))
	// checkAndSetNodeManagementPortIds(data)
	// DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/managementPorts/%s", data.NodeId.ValueString(), data.NodeManagementPortId.ValueString()), "DELETE", nil)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }
	tflog.Debug(ctx, fmt.Sprintf("End delete of resource hyperfabric_node_management_port with id '%s'", data.Id.ValueString()))
}

func (r *NodeManagementPortResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Start import state of resource: hyperfabric_node_management_port")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	var stateData *NodeManagementPortResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &stateData)...)
	tflog.Debug(ctx, fmt.Sprintf("Import state of resource hyperfabric_node_management_port with id '%s'", stateData.Id.ValueString()))
	tflog.Debug(ctx, "End import of state resource: hyperfabric_node_management_port")
}

func getAndSetNodeManagementPortAttributes(ctx context.Context, diags *diag.Diagnostics, client *client.Client, data *NodeManagementPortResourceModel) {
	// requestData := DoRestRequest(ctx, diags, client, fmt.Sprintf("/api/v1/fabrics/%s/managementPorts/%s", data.NodeId.ValueString(), data.NodeManagementPortId.ValueString()), "GET", nil)
	requestData := DoRestRequest(ctx, diags, client, fmt.Sprintf("/api/v1/fabrics/%s/managementPorts", data.NodeId.ValueString()), "GET", nil)
	if diags.HasError() {
		return
	}

	newNodeManagementPort := *getEmptyNodeManagementPortResourceModel()
	// newNodeManagementPort.FabricId = data.FabricId
	newNodeManagementPort.NodeId = data.NodeId
	newNodeManagementPort.NodeManagementPortId = data.NodeManagementPortId
	newNodeManagementPort.Id = data.Id
	// Keeping Proxy password as set in state.
	newNodeManagementPort.ProxyPassword = data.ProxyPassword

	if requestData.Data() != nil {
		requestMap := requestData.Data().(map[string]interface{})
		for _, ports := range requestMap {
			listPorts := ports.([]interface{})
			if len(listPorts) == 1 {
				for attributeName, attributeValue := range listPorts[0].(map[string]interface{}) {
					// if attributeName == "nodeId" && (data.NodeId.IsNull() || data.NodeId.IsUnknown() || data.NodeId.ValueString() == "" || data.NodeId.ValueString() != attributeValue.(string)) {
					// 	newNodeManagementPort.NodeId = basetypes.NewStringValue(attributeValue.(string))
					// 	newNodeManagementPort.Id = basetypes.NewStringValue(fmt.Sprintf("%s/nodes/%s/managementPorts/%s", newNodeManagementPort.FabricId.ValueString(), newNodeManagementPort.NodeId.ValueString(), newNodeManagementPort.NodeManagementPortId.ValueString()))
					// } else
					if attributeName == "id" && (data.NodeManagementPortId.IsNull() || data.NodeManagementPortId.IsUnknown() || data.NodeManagementPortId.ValueString() == "" || data.NodeManagementPortId.ValueString() != attributeValue.(string)) {
						newNodeManagementPort.NodeManagementPortId = basetypes.NewStringValue(attributeValue.(string))
						newNodeManagementPort.Id = basetypes.NewStringValue(fmt.Sprintf("%s/managementPorts/%s", newNodeManagementPort.NodeId.ValueString(), newNodeManagementPort.NodeManagementPortId.ValueString()))
					} else if attributeName == "name" {
						newNodeManagementPort.Name = basetypes.NewStringValue(attributeValue.(string))
					} else if attributeName == "description" {
						newNodeManagementPort.Description = basetypes.NewStringValue(attributeValue.(string))
					} else if attributeName == "enabled" {
						newNodeManagementPort.Enabled = basetypes.NewBoolValue(attributeValue.(bool))
					} else if attributeName == "cloudUrls" {
						newNodeManagementPort.CloudUrls = NewCloudUrlsSet(ctx, attributeValue.([]interface{}))
					} else if attributeName == "ipv4ConfigType" {
						newNodeManagementPort.Ipv4ConfigType = basetypes.NewStringValue(attributeValue.(string))
					} else if attributeName == "ipv4Address" {
						newNodeManagementPort.Ipv4Address = basetypes.NewStringValue(attributeValue.(string))
					} else if attributeName == "ipv4Gateway" {
						newNodeManagementPort.Ipv4Gateway = basetypes.NewStringValue(attributeValue.(string))
					} else if attributeName == "ipv6ConfigType" {
						newNodeManagementPort.Ipv6ConfigType = basetypes.NewStringValue(attributeValue.(string))
					} else if attributeName == "ipv6Address" {
						newNodeManagementPort.Ipv6Address = basetypes.NewStringValue(attributeValue.(string))
					} else if attributeName == "ipv6Gateway" {
						newNodeManagementPort.Ipv6Gateway = basetypes.NewStringValue(attributeValue.(string))
					} else if attributeName == "dnsAddresses" {
						newNodeManagementPort.DnsAddresses = NewDnsAddressesSet(ctx, attributeValue.([]interface{}))
					} else if attributeName == "ntpAddresses" {
						newNodeManagementPort.NtpAddresses = NewNtpAddressesSet(ctx, attributeValue.([]interface{}))
					} else if attributeName == "noProxy" {
						newNodeManagementPort.NoProxy = NewNoProxySet(ctx, attributeValue.([]interface{}))
					} else if attributeName == "proxyAddress" {
						newNodeManagementPort.ProxyAddress = basetypes.NewStringValue(attributeValue.(string))
					} else if attributeName == "proxyCredentialId" {
						newNodeManagementPort.ProxyCredentialId = basetypes.NewStringValue(attributeValue.(string))
					} else if attributeName == "proxyUsername" {
						newNodeManagementPort.ProxyUsername = basetypes.NewStringValue(attributeValue.(string))
						// Not setting password as it is not returned and want to keep state intact
						// } else if attributeName == "proxyPassword" {
						// 	newNodeManagementPort.ProxyPassword = basetypes.NewStringValue(attributeValue.(string))
						// } else if attributeName == "setProxyPassword" {
						// 	newNodeManagementPort.SetProxyPassword = basetypes.NewBoolValue(attributeValue.(bool))
					} else if attributeName == "connectedState" {
						newNodeManagementPort.ConnectedState = basetypes.NewStringValue(attributeValue.(string))
					} else if attributeName == "configOrigin" {
						newNodeManagementPort.ConfigOrigin = basetypes.NewStringValue(attributeValue.(string))
					} else if attributeName == "metadata" {
						newNodeManagementPort.Metadata = NewMetadataObject(ctx, attributeValue.(map[string]interface{}))
						// } else if attributeName == "labels" {
						// 	newNodeManagementPort.Labels = NewSetString(ctx, attributeValue.([]interface{}))
						// } else if attributeName == "annotations" {
						// 	newNodeManagementPort.Annotations = NewAnnotationsSet(ctx, attributeValue.([]interface{}))
					}
				}
			} else {
				tflog.Debug(ctx, fmt.Sprintf("Wrong number of management ports in hyperfabric_node_management_port with id '%s", data.Id.ValueString()))
				newNodeManagementPort.Id = basetypes.NewStringNull()
			}
		}
	} else {
		newNodeManagementPort.Id = basetypes.NewStringNull()
	}
	*data = newNodeManagementPort
}

func getNodeManagementPortJsonPayload(ctx context.Context, diags *diag.Diagnostics, data *NodeManagementPortResourceModel, action string) *gabs.Container {
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

	if !data.CloudUrls.IsNull() && !data.CloudUrls.IsUnknown() {
		payloadMap["cloud_urls"] = getCloudUrlsJsonPayload(ctx, data.CloudUrls)
	}

	if !data.Ipv4ConfigType.IsNull() && !data.Ipv4ConfigType.IsUnknown() {
		payloadMap["ipv4ConfigType"] = data.Ipv4ConfigType.ValueString()
	}

	if !data.Ipv4Address.IsNull() && !data.Ipv4Address.IsUnknown() {
		payloadMap["ipv4Address"] = data.Ipv4Address.ValueString()
	}

	if !data.Ipv4Gateway.IsNull() && !data.Ipv4Gateway.IsUnknown() {
		payloadMap["ipv4Gateway"] = data.Ipv4Gateway.ValueString()
	}

	if !data.Ipv6ConfigType.IsNull() && !data.Ipv6ConfigType.IsUnknown() {
		payloadMap["ipv6ConfigType"] = data.Ipv6ConfigType.ValueString()
	}

	if !data.Ipv6Address.IsNull() && !data.Ipv6Address.IsUnknown() {
		payloadMap["ipv6Address"] = data.Ipv6Address.ValueString()
	}

	if !data.Ipv6Gateway.IsNull() && !data.Ipv6Gateway.IsUnknown() {
		payloadMap["ipv6Gateway"] = data.Ipv6Gateway.ValueString()
	}

	if !data.DnsAddresses.IsNull() && !data.DnsAddresses.IsUnknown() {
		payloadMap["dnsAddresses"] = getDnsAddressesJsonPayload(ctx, data.DnsAddresses)
	}

	if !data.NtpAddresses.IsNull() && !data.NtpAddresses.IsUnknown() {
		payloadMap["ntpAddresses"] = getNtpAddressesJsonPayload(ctx, data.NtpAddresses)
	}

	if !data.NoProxy.IsNull() && !data.NoProxy.IsUnknown() {
		payloadMap["noProxy"] = getNoProxyJsonPayload(ctx, data.NoProxy)
	}

	if !data.ProxyAddress.IsNull() && !data.ProxyAddress.IsUnknown() {
		payloadMap["proxyAddress"] = data.ProxyAddress.ValueString()
	}

	if !data.ProxyPassword.IsNull() && !data.ProxyPassword.IsUnknown() {
		payloadMap["proxyPassword"] = data.ProxyPassword.ValueString()
		payloadMap["setProxyPassword"] = true
	}

	if !data.ProxyUsername.IsNull() && !data.ProxyUsername.IsUnknown() {
		payloadMap["proxyUsername"] = data.ProxyUsername.ValueString()
	}

	if (data.ProxyPassword.IsNull() && data.ProxyPassword.IsUnknown() && data.ProxyUsername.IsNull() && data.ProxyUsername.IsUnknown()) ||
		(!data.ProxyUsername.IsNull() && !data.ProxyUsername.IsUnknown() && data.ProxyUsername.ValueString() == "") {
		payloadMap["setProxyPassword"] = false
		delete(payloadMap, "proxyUsername")
		delete(payloadMap, "proxyPassword")
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
		payload = map[string]interface{}{"ports": payloadList}
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

func getCloudUrlsJsonPayload(ctx context.Context, data basetypes.SetValue) []string {
	cloudUrls := []string{}
	data.ElementsAs(ctx, &cloudUrls, false)
	return cloudUrls
}

func NewCloudUrlsSet(ctx context.Context, data []interface{}) basetypes.SetValue {
	cloudUrlsSet, _ := types.SetValueFrom(ctx, CloudUrlsResourceModelAttributeType(), data)
	return cloudUrlsSet
}

func getDnsAddressesJsonPayload(ctx context.Context, data basetypes.SetValue) []string {
	dnsAddresses := []string{}
	data.ElementsAs(ctx, &dnsAddresses, false)
	return dnsAddresses
}

func NewDnsAddressesSet(ctx context.Context, data []interface{}) basetypes.SetValue {
	dnsAddressesSet, _ := types.SetValueFrom(ctx, DnsAddressesResourceModelAttributeType(), data)
	return dnsAddressesSet
}

func getNtpAddressesJsonPayload(ctx context.Context, data basetypes.SetValue) []string {
	ntpAddresses := []string{}
	data.ElementsAs(ctx, &ntpAddresses, false)
	return ntpAddresses
}

func NewNtpAddressesSet(ctx context.Context, data []interface{}) basetypes.SetValue {
	ntpAddressesSet, _ := types.SetValueFrom(ctx, NtpAddressesResourceModelAttributeType(), data)
	return ntpAddressesSet
}

func getNoProxyJsonPayload(ctx context.Context, data basetypes.SetValue) []string {
	noProxies := []string{}
	data.ElementsAs(ctx, &noProxies, false)
	return noProxies
}

func NewNoProxySet(ctx context.Context, data []interface{}) basetypes.SetValue {
	noProxiesSet, _ := types.SetValueFrom(ctx, NoProxyResourceModelAttributeType(), data)
	return noProxiesSet
}

func checkAndSetNodeManagementPortIds(data *NodeManagementPortResourceModel) {
	if strings.Contains(data.Id.ValueString(), "/managementPorts/") {
		if data.NodeId.IsNull() || data.NodeId.IsUnknown() || data.NodeId.ValueString() == "" ||
			data.NodeManagementPortId.IsNull() || data.NodeManagementPortId.IsUnknown() || data.NodeManagementPortId.ValueString() == "" {
			splitId := strings.Split(data.Id.ValueString(), "/managementPorts/")
			data.NodeId = basetypes.NewStringValue(splitId[0])
			data.NodeManagementPortId = basetypes.NewStringValue(splitId[1])
		}
	}
}
