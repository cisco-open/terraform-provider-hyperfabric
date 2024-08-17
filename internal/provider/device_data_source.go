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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DeviceDataSource{}

func NewDeviceDataSource() datasource.DataSource {
	return &DeviceDataSource{}
}

// DeviceDataSource defines the data source implementation.
type DeviceDataSource struct {
	client *client.Client
}

// DeviceDataSourceModel describes the data source data model.
type DeviceDataSourceModel struct {
	Id           types.String `tfsdk:"id"`
	DeviceId     types.String `tfsdk:"device_id"`
	ModelName    types.String `tfsdk:"model_name"`
	FabricId     types.String `tfsdk:"fabric_id"`
	NodeId       types.String `tfsdk:"node_id"`
	OsType       types.String `tfsdk:"os_type"`
	RackId       types.String `tfsdk:"rack_id"`
	Roles        types.Set    `tfsdk:"roles"`
	SerialNumber types.String `tfsdk:"serial_number"`
}

func getEmptyDeviceDataSourceModel() *DeviceDataSourceModel {
	return &DeviceDataSourceModel{
		Id:           basetypes.NewStringNull(),
		DeviceId:     basetypes.NewStringNull(),
		ModelName:    basetypes.NewStringNull(),
		FabricId:     basetypes.NewStringNull(),
		NodeId:       basetypes.NewStringNull(),
		OsType:       basetypes.NewStringNull(),
		RackId:       basetypes.NewStringNull(),
		Roles:        basetypes.NewSetNull(SetStringResourceModelAttributeType()),
		SerialNumber: basetypes.NewStringNull(),
	}
}

func (d *DeviceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	tflog.Debug(ctx, "Start metadata of datasource: hyperfabric_device")
	resp.TypeName = req.ProviderTypeName + "_device"
	tflog.Debug(ctx, "End metadata of datasource: hyperfabric_device")
}

func (d *DeviceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	tflog.Debug(ctx, "Start schema of datasource: hyperfabric_device")
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Device data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "`id` defines the unique identifier of a Device.",
				Computed:            true,
			},
			"device_id": schema.StringAttribute{
				MarkdownDescription: "`device_id` defines the unique identifier of a Device.",
				Optional:            true,
				Computed:            true,
			},
			"model_name": schema.StringAttribute{
				MarkdownDescription: "The model name of the Device.",
				Computed:            true,
			},
			"fabric_id": schema.StringAttribute{
				MarkdownDescription: "`fabric_id` defines the unique identifier of a Fabric.",
				Computed:            true,
			},
			"node_id": schema.StringAttribute{
				MarkdownDescription: "`node_id` defines the unique identifier of a Node.",
				Computed:            true,
			},
			"serial_number": schema.StringAttribute{
				MarkdownDescription: "The serial number of the Device.",
				Optional:            true,
				Computed:            true,
			},
			"os_type": schema.StringAttribute{
				MarkdownDescription: "The operating system type of the Device.",
				Computed:            true,
			},
			"rack_id": schema.StringAttribute{
				MarkdownDescription: "`rack_id` defines the unique identifier of a Rack.",
				Computed:            true,
			},
			"roles": getDeviceRolesSchemaAttribute(),
		},
	}
	tflog.Debug(ctx, "End schema of datasource: hyperfabric_device")
}

func getDeviceRolesSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A set of roles used by the Device.`,
		Computed:            true,
		ElementType:         types.StringType,
	}
}

func (d *DeviceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of datasource: hyperfabric_device")
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
	tflog.Debug(ctx, "End configure of datasource: hyperfabric_device")
}

func (d *DeviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Start read of datasource: hyperfabric_device")
	var data *DeviceDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create a copy of the Id for when not found during getAndSetNodePortAttributes
	cachedId := data.Id.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Read of datasource hyperfabric_device '%v'", data.SerialNumber))

	getAndSetDeviceAttributes(ctx, &resp.Diagnostics, d.client, data)

	if data.Id.IsNull() {
		resp.Diagnostics.AddError(
			"Failed to read hyperfabric_device data source",
			fmt.Sprintf("The hyperfabric_device data source with id '%s' has not been found", cachedId),
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End read of datasource hyperfabric_device with id '%s'", data.Id.ValueString()))
}

func getAndSetDeviceAttributes(ctx context.Context, diags *diag.Diagnostics, client *client.Client, data *DeviceDataSourceModel) {
	requestData := DoRestRequest(ctx, diags, client, "/api/v1/devices", "GET", nil)
	if diags.HasError() {
		return
	}

	if requestData.Data() != nil {
		responseMap := requestData.Data().(map[string]interface{})
		for responseKey, responseValue := range responseMap {
			if responseKey == "devices" {
				devices := responseValue.([]interface{})
				for _, device := range devices {
					newDevice := *getEmptyDeviceDataSourceModel()
					for attributeName, attributeValue := range device.(map[string]interface{}) {
						if attributeName == "deviceId" {
							newDevice.Id = basetypes.NewStringValue(attributeValue.(string))
							newDevice.DeviceId = basetypes.NewStringValue(attributeValue.(string))
						} else if attributeName == "fabricId" {
							newDevice.FabricId = basetypes.NewStringValue(attributeValue.(string))
						} else if attributeName == "modelName" {
							newDevice.ModelName = basetypes.NewStringValue(attributeValue.(string))
						} else if attributeName == "nodeId" {
							newDevice.NodeId = basetypes.NewStringValue(attributeValue.(string))
						} else if attributeName == "osType" {
							newDevice.OsType = basetypes.NewStringValue(attributeValue.(string))
						} else if attributeName == "rackId" {
							newDevice.RackId = basetypes.NewStringValue(attributeValue.(string))
						} else if attributeName == "roles" {
							newDevice.Roles = NewSetString(ctx, attributeValue.([]interface{}))
						} else if attributeName == "serialNumber" {
							newDevice.SerialNumber = basetypes.NewStringValue(attributeValue.(string))
						}
					}
					if (!data.SerialNumber.IsNull() && !data.SerialNumber.IsUnknown() && data.SerialNumber.ValueString() != "" && newDevice.SerialNumber == data.SerialNumber) ||
						(!data.DeviceId.IsNull() && !data.DeviceId.IsUnknown() && data.DeviceId.ValueString() != "" && newDevice.DeviceId == data.DeviceId) {
						*data = newDevice
					}
				}
			}
		}
	}
}
