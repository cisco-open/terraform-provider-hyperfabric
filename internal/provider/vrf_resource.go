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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &VrfResource{}
var _ resource.ResourceWithImportState = &VrfResource{}

func NewVrfResource() resource.Resource {
	return &VrfResource{}
}

// VrfResource defines the resource implementation.
type VrfResource struct {
	client *client.Client
}

// VrfResourceModel describes the resource data model.
type VrfResourceModel struct {
	Id          types.String  `tfsdk:"id"`
	VrfId       types.String  `tfsdk:"vrf_id"`
	FabricId    types.String  `tfsdk:"fabric_id"`
	Name        types.String  `tfsdk:"name"`
	Description types.String  `tfsdk:"description"`
	Enabled     types.Bool    `tfsdk:"enabled"`
	IsDefault   types.Bool    `tfsdk:"is_default"`
	Asn         types.Float64 `tfsdk:"asn"`
	Vni         types.Float64 `tfsdk:"vni"`
	RouteTarget types.String  `tfsdk:"route_target"`
	Metadata    types.Object  `tfsdk:"metadata"`
	Labels      types.Set     `tfsdk:"labels"`
	Annotations types.Set     `tfsdk:"annotations"`
}

func getEmptyVrfResourceModel() *VrfResourceModel {
	return &VrfResourceModel{
		Id:          basetypes.NewStringNull(),
		VrfId:       basetypes.NewStringNull(),
		FabricId:    basetypes.NewStringNull(),
		Name:        basetypes.NewStringNull(),
		Description: basetypes.NewStringNull(),
		Enabled:     basetypes.NewBoolNull(),
		IsDefault:   basetypes.NewBoolNull(),
		Asn:         basetypes.NewFloat64Null(),
		Vni:         basetypes.NewFloat64Null(),
		RouteTarget: basetypes.NewStringNull(),
		Metadata:    basetypes.NewObjectNull(MetadataResourceModelAttributeType()),
		Labels:      basetypes.NewSetNull(SetStringResourceModelAttributeType()),
		Annotations: basetypes.NewSetNull(AnnotationResourceModelAttributeType()),
	}
}

func getNewVrfResourceModelFromData(data *VrfResourceModel) *VrfResourceModel {
	newVrf := getEmptyVrfResourceModel()

	if !data.Id.IsNull() && !data.Id.IsUnknown() {
		newVrf.Id = data.Id
	}

	if !data.VrfId.IsNull() && !data.VrfId.IsUnknown() {
		newVrf.VrfId = data.VrfId
	}

	if !data.FabricId.IsNull() && !data.FabricId.IsUnknown() {
		newVrf.FabricId = data.FabricId
	}

	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		newVrf.Name = data.Name
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		newVrf.Description = data.Description
	}

	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		newVrf.Enabled = data.Enabled
	}

	if !data.IsDefault.IsNull() && !data.IsDefault.IsUnknown() {
		newVrf.IsDefault = data.IsDefault
	}

	if !data.Asn.IsNull() && !data.Asn.IsUnknown() {
		newVrf.Asn = data.Asn
	}

	if !data.Vni.IsNull() && !data.Vni.IsUnknown() {
		newVrf.Vni = data.Vni
	}

	if !data.RouteTarget.IsNull() && !data.RouteTarget.IsUnknown() {
		newVrf.RouteTarget = data.RouteTarget
	}

	if !data.Metadata.IsNull() && !data.Metadata.IsUnknown() {
		newVrf.Metadata = data.Metadata
	}

	if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
		newVrf.Labels = data.Labels
	}

	if !data.Annotations.IsNull() && !data.Annotations.IsUnknown() {
		newVrf.Annotations = data.Annotations
	}

	return newVrf
}

type VrfIdentifier struct {
	Id types.String
}

func (r *VrfResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Debug(ctx, "Start metadata of resource: hyperfabric_vrf")
	resp.TypeName = req.ProviderTypeName + "_vrf"
	tflog.Debug(ctx, "End metadata of resource: hyperfabric_vrf")
}

func (r *VrfResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Debug(ctx, "Start schema of resource: hyperfabric_vrf")
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "VRF resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "`id` defines the unique identifier of a VRF in a Fabric.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vrf_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "`vrf_id` defines the unique identifier of a VRF.",
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
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the VRF.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description is a user defined field to store notes about the VRF.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "The enabled state of the VRF.",
				// Optional:            true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					SetToBoolNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"is_default": schema.BoolAttribute{
				MarkdownDescription: "The flag that denote if the VRF is the default VRF or not.",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					SetToBoolNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"asn": schema.Float64Attribute{
				MarkdownDescription: "The Autonomous System Number (ASN) used for the VRF external connections.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
					SetToFloat64NullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"vni": schema.Float64Attribute{
				MarkdownDescription: "The VXLAN Network Identifier (VNI) used for the VRF.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
					SetToFloat64NullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"route_target": schema.StringAttribute{
				MarkdownDescription: "The route target associated with the VRF.",
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
	tflog.Debug(ctx, "End schema of resource: hyperfabric_vrf")
}

func (r *VrfResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of resource: hyperfabric_vrf")
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
	tflog.Debug(ctx, "End configure of resource: hyperfabric_vrf")
}

func (r *VrfResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Start create of resource: hyperfabric_vrf")

	var data *VrfResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Create of resource hyperfabric_vrf in fabric '%s' with VRF name '%s'", data.FabricId.ValueString(), data.Name.ValueString()))

	jsonPayload := getVrfJsonPayload(ctx, &resp.Diagnostics, data, "create")
	if resp.Diagnostics.HasError() {
		return
	}

	container := DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/vrfs", data.FabricId.ValueString()), "POST", jsonPayload)
	if resp.Diagnostics.HasError() {
		return
	}

	vrfContainer, err := container.ArrayElement(0, "vrfs")
	if err != nil {
		return
	}

	vrfId := StripQuotes(vrfContainer.Search("id").String())
	if vrfId != "" {
		data.Id = basetypes.NewStringValue(fmt.Sprintf("%s/vrfs/%s", data.FabricId.ValueString(), vrfId))
		data.VrfId = basetypes.NewStringValue(vrfId)
		getAndSetVrfAttributes(ctx, &resp.Diagnostics, r.client, data)
	} else {
		data.Id = basetypes.NewStringNull()
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End create of resource hyperfabric_vrf with id '%s'", data.Id.ValueString()))
}

func (r *VrfResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Start read of resource: hyperfabric_vrf")
	var data *VrfResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Read of resource hyperfabric_vrf with id '%s'", data.Id.ValueString()))
	checkAndSetVrfIds(data)
	getAndSetVrfAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	if data.Id.IsNull() {
		var emptyData *VrfResourceModel
		resp.Diagnostics.Append(resp.State.Set(ctx, &emptyData)...)
	} else {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
	tflog.Debug(ctx, fmt.Sprintf("End read of resource hyperfabric_vrf with id '%s'", data.Id.ValueString()))
}

func (r *VrfResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Start update of resource: hyperfabric_vrf")
	var data *VrfResourceModel
	var stateData *VrfResourceModel

	// Read Terraform plan and state data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Update of resource hyperfabric_vrf with id '%s'", data.Id.ValueString()))

	jsonPayload := getVrfJsonPayload(ctx, &resp.Diagnostics, data, "update")

	if resp.Diagnostics.HasError() {
		return
	}

	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/vrfs/%s", data.FabricId.ValueString(), data.VrfId.ValueString()), "PUT", jsonPayload)

	if resp.Diagnostics.HasError() {
		return
	}

	getAndSetVrfAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End update of resource hyperfabric_vrf with id '%s'", data.Id.ValueString()))
}

func (r *VrfResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Start delete of resource: hyperfabric_vrf")
	var data *VrfResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Delete of resource hyperfabric_vrf with id '%s'", data.Id.ValueString()))
	checkAndSetVrfIds(data)
	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/vrfs/%s", data.FabricId.ValueString(), data.VrfId.ValueString()), "DELETE", nil)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("End delete of resource hyperfabric_vrf with id '%s'", data.Id.ValueString()))
}

func (r *VrfResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Start import state of resource: hyperfabric_vrf")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	var stateData *VrfResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &stateData)...)
	tflog.Debug(ctx, fmt.Sprintf("Import state of resource hyperfabric_vrf with id '%s'", stateData.Id.ValueString()))
	tflog.Debug(ctx, "End import of state resource: hyperfabric_vrf with id")
}

func getAndSetVrfAttributes(ctx context.Context, diags *diag.Diagnostics, client *client.Client, data *VrfResourceModel) {
	requestData := DoRestRequest(ctx, diags, client, fmt.Sprintf("/api/v1/fabrics/%s/vrfs/%s", data.FabricId.ValueString(), data.VrfId.ValueString()), "GET", nil)
	if diags.HasError() {
		return
	}

	newVrf := *getNewVrfResourceModelFromData(data)

	if requestData.Data() != nil {
		attributes := requestData.Data().(map[string]interface{})
		for attributeName, attributeValue := range attributes {
			if attributeName == "fabricId" && (data.FabricId.IsNull() || data.FabricId.IsUnknown() || data.FabricId.ValueString() == "" || data.FabricId.ValueString() != attributeValue.(string)) {
				newVrf.FabricId = basetypes.NewStringValue(attributeValue.(string))
				newVrf.Id = basetypes.NewStringValue(fmt.Sprintf("%s/vrfs/%s", newVrf.FabricId.ValueString(), newVrf.VrfId.ValueString()))
			} else if attributeName == "id" && (data.VrfId.IsNull() || data.VrfId.IsUnknown() || data.VrfId.ValueString() == "" || data.VrfId.ValueString() != attributeValue.(string)) {
				newVrf.VrfId = basetypes.NewStringValue(attributeValue.(string))
				newVrf.Id = basetypes.NewStringValue(fmt.Sprintf("%s/vrfs/%s", newVrf.FabricId.ValueString(), newVrf.VrfId.ValueString()))
			} else if attributeName == "name" {
				newVrf.Name = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "description" {
				newVrf.Description = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "enabled" {
				newVrf.Enabled = basetypes.NewBoolValue(attributeValue.(bool))
			} else if attributeName == "isDefault" {
				newVrf.IsDefault = basetypes.NewBoolValue(attributeValue.(bool))
			} else if attributeName == "asn" {
				newVrf.Asn = basetypes.NewFloat64Value(attributeValue.(float64))
			} else if attributeName == "vni" {
				newVrf.Vni = basetypes.NewFloat64Value(attributeValue.(float64))
			} else if attributeName == "rt" {
				newVrf.RouteTarget = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "metadata" {
				newVrf.Metadata = NewMetadataObject(ctx, attributeValue.(map[string]interface{}))
			} else if attributeName == "labels" {
				newVrf.Labels = NewSetString(ctx, attributeValue.([]interface{}))
			} else if attributeName == "annotations" {
				newVrf.Annotations = NewAnnotationsSet(ctx, attributeValue.([]interface{}))
			}
		}
	} else {
		newVrf.Id = basetypes.NewStringNull()
	}
	*data = newVrf
}

func getVrfJsonPayload(ctx context.Context, diags *diag.Diagnostics, data *VrfResourceModel, action string) *gabs.Container {
	payloadMap := map[string]interface{}{}
	payloadList := []map[string]interface{}{}

	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		payloadMap["name"] = data.Name.ValueString()
	}

	// TO FIX: NEED TO BE REMOVE WHEN FIXED
	if !data.FabricId.IsNull() && !data.FabricId.IsUnknown() {
		payloadMap["fabricId"] = data.FabricId.ValueString()
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		payloadMap["description"] = data.Description.ValueString()
	}

	// TO FIX
	payloadMap["enabled"] = true
	// if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
	// 	payloadMap["enabled"] = data.Enabled.ValueBool()
	// }

	if !data.Asn.IsNull() && !data.Asn.IsUnknown() {
		payloadMap["asn"] = data.Asn.ValueFloat64()
	}

	if !data.Vni.IsNull() && !data.Vni.IsUnknown() {
		payloadMap["vni"] = data.Vni.ValueFloat64()
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
		payload = map[string]interface{}{"vrfs": payloadList}
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

func checkAndSetVrfIds(data *VrfResourceModel) {
	if strings.Contains(data.Id.ValueString(), "/vrfs/") {
		if data.FabricId.IsNull() || data.FabricId.IsUnknown() || data.FabricId.ValueString() == "" || data.VrfId.IsNull() || data.VrfId.IsUnknown() || data.VrfId.ValueString() == "" {
			splitId := strings.Split(data.Id.ValueString(), "/vrfs/")
			data.FabricId = basetypes.NewStringValue(splitId[0])
			data.VrfId = basetypes.NewStringValue(splitId[1])
		}
	}
}
