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
var _ resource.Resource = &VniResource{}
var _ resource.ResourceWithImportState = &VniResource{}

func NewVniResource() resource.Resource {
	return &VniResource{}
}

// VniResource defines the resource implementation.
type VniResource struct {
	client *client.Client
}

// VniResourceModel describes the resource data model.
type VniResourceModel struct {
	Id          types.String `tfsdk:"id"`
	VniId       types.String `tfsdk:"vni_id"`
	FabricId    types.String `tfsdk:"fabric_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	IsDefault   types.Bool   `tfsdk:"is_default"`
	// IsL3        types.Bool    `tfsdk:"is_l3"`
	VrfId       types.String  `tfsdk:"vrf_id"`
	Vni         types.Float64 `tfsdk:"vni"`
	Mtu         types.Float64 `tfsdk:"mtu"`
	Members     types.Set     `tfsdk:"members"`
	Svi         types.Object  `tfsdk:"svi"`
	Metadata    types.Object  `tfsdk:"metadata"`
	Labels      types.Set     `tfsdk:"labels"`
	Annotations types.Set     `tfsdk:"annotations"`
}

func getEmptyVniResourceModel() *VniResourceModel {
	return &VniResourceModel{
		Id:          basetypes.NewStringNull(),
		VniId:       basetypes.NewStringNull(),
		FabricId:    basetypes.NewStringNull(),
		Name:        basetypes.NewStringNull(),
		Description: basetypes.NewStringNull(),
		Enabled:     basetypes.NewBoolNull(),
		IsDefault:   basetypes.NewBoolNull(),
		// IsL3:        basetypes.NewBoolNull(),
		VrfId:       basetypes.NewStringNull(),
		Vni:         basetypes.NewFloat64Null(),
		Mtu:         basetypes.NewFloat64Null(),
		Members:     basetypes.NewSetNull(MemberResourceModelAttributeType()),
		Svi:         basetypes.NewObjectNull(SviResourceModelAttributeType()),
		Metadata:    basetypes.NewObjectNull(MetadataResourceModelAttributeType()),
		Labels:      basetypes.NewSetNull(SetStringResourceModelAttributeType()),
		Annotations: basetypes.NewSetNull(AnnotationResourceModelAttributeType()),
	}
}

func getNewVniResourceModelFromData(data *VniResourceModel) *VniResourceModel {
	newVni := getEmptyVniResourceModel()

	if !data.Id.IsNull() && !data.Id.IsUnknown() {
		newVni.Id = data.Id
	}

	if !data.VniId.IsNull() && !data.VniId.IsUnknown() {
		newVni.VniId = data.VniId
	}

	if !data.FabricId.IsNull() && !data.FabricId.IsUnknown() {
		newVni.FabricId = data.FabricId
	}

	if !data.VrfId.IsNull() && !data.VrfId.IsUnknown() {
		newVni.VrfId = data.VrfId
	}

	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		newVni.Name = data.Name
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		newVni.Description = data.Description
	}

	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		newVni.Enabled = data.Enabled
	}

	if !data.IsDefault.IsNull() && !data.IsDefault.IsUnknown() {
		newVni.IsDefault = data.IsDefault
	}

	if !data.Vni.IsNull() && !data.Vni.IsUnknown() {
		newVni.Vni = data.Vni
	}

	if !data.Mtu.IsNull() && !data.Mtu.IsUnknown() {
		newVni.Mtu = data.Mtu
	}

	if !data.Members.IsNull() && !data.Members.IsUnknown() {
		newVni.Members = data.Members
	}

	if !data.Svi.IsNull() && !data.Svi.IsUnknown() {
		newVni.Svi = data.Svi
	}

	if !data.Metadata.IsNull() && !data.Metadata.IsUnknown() {
		newVni.Metadata = data.Metadata
	}

	if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
		newVni.Labels = data.Labels
	}

	if !data.Annotations.IsNull() && !data.Annotations.IsUnknown() {
		newVni.Annotations = data.Annotations
	}

	return newVni
}

type VniIdentifier struct {
	Id types.String
}

func (r *VniResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if !req.Plan.Raw.IsNull() {
		var planData, stateData, configData *VniResourceModel
		resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
		resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
		resp.Diagnostics.Append(req.Config.Get(ctx, &configData)...)

		if resp.Diagnostics.HasError() {
			return
		}

		// if (planData.Id.IsUnknown() || planData.Id.IsNull()) && !planData.ParentDn.IsUnknown() && !planData.Name.IsUnknown() {
		// 	setNetflowMonitorPolId(ctx, planData)
		// }

		// if stateData == nil && !globalAllowExistingOnCreate && !planData.Id.IsUnknown() && !planData.Id.IsNull() {
		// 	CheckDn(ctx, &resp.Diagnostics, r.client, "netflowMonitorPol", planData.Id.ValueString())
		// 	if resp.Diagnostics.HasError() {
		// 		return
		// 	}
		// }
		if !configData.Svi.IsNull() && stateData != nil {
			if IsEmptySingleNestedAttribute(configData.Svi.Attributes()) {
				svi, _ := types.ObjectValueFrom(ctx, SviResourceModelAttributeType(), getEmptySviResourceModel())
				planData.Svi = svi
			}
		}

		resp.Diagnostics.Append(resp.Plan.Set(ctx, &planData)...)
	}
}

func (r *VniResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Debug(ctx, "Start metadata of resource: hyperfabric_vni")
	resp.TypeName = req.ProviderTypeName + "_vni"
	tflog.Debug(ctx, "End metadata of resource: hyperfabric_vni")
}

func (r *VniResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Debug(ctx, "Start schema of resource: hyperfabric_vni")
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "VNI resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "`id` defines the unique identifier of a VNI in a Fabric.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vni_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "`vni_id` defines the unique identifier of a VNI.",
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
				MarkdownDescription: "The name of the VNI.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description is a user defined field to store notes about the VNI.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "The enabled state of the VNI.",
				// Optional:            true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					SetToBoolNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"is_default": schema.BoolAttribute{
				MarkdownDescription: "The flag that denote if the VNI is the default VNI or not.",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					SetToBoolNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			// "is_l3": schema.BoolAttribute{
			// 	MarkdownDescription: "The flag that denote if the VNI is L3 only.",
			// 	// Optional:            true,
			// 	Computed: true,
			// 	PlanModifiers: []planmodifier.Bool{
			// 		boolplanmodifier.UseStateForUnknown(),
			// 		SetToBoolNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
			// 	},
			// },
			"vrf_id": schema.StringAttribute{
				MarkdownDescription: "The Id of the VRF associated with the VNI. Used when VNI is L3 (l2_only=false).",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					// SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"vni": schema.Float64Attribute{
				MarkdownDescription: "The VXLAN Network Identifier (VNI) used for the VRF.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.RequiresReplace(),
					float64planmodifier.UseStateForUnknown(),
					SetToFloat64NullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"mtu": schema.Float64Attribute{
				MarkdownDescription: "The MTU of the SVI of the VNI.",
				// Optional:            true,
				Computed: true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
					SetToFloat64NullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"members":     getMembersSchemaAttribute(),
			"svi":         getSviSchemaAttribute(),
			"metadata":    getMetadataSchemaAttribute(),
			"labels":      getLabelsSchemaAttribute(),
			"annotations": getAnnotationsSchemaAttribute(),
		},
	}
	tflog.Debug(ctx, "End schema of resource: hyperfabric_vni")
}

func (r *VniResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of resource: hyperfabric_vni")
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
	tflog.Debug(ctx, "End configure of resource: hyperfabric_vni")
}

func (r *VniResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Start create of resource: hyperfabric_vni")

	var data *VniResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Create of resource hyperfabric_vni in Fabric '%s' with name '%s'", data.FabricId.ValueString(), data.Name.ValueString()))

	jsonPayload := getVniJsonPayload(ctx, &resp.Diagnostics, data, "create")
	if resp.Diagnostics.HasError() {
		return
	}

	container := DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/vnis", data.FabricId.ValueString()), "POST", jsonPayload)
	if resp.Diagnostics.HasError() {
		return
	}

	vniContainer, err := container.ArrayElement(0, "vnis")
	if err != nil {
		return
	}
	vniId := StripQuotes(vniContainer.Search("id").String())
	if vniId != "" {
		data.Id = basetypes.NewStringValue(fmt.Sprintf("%s/vnis/%s", data.FabricId.ValueString(), vniId))
		data.VniId = basetypes.NewStringValue(vniId)
		getAndSetVniAttributes(ctx, &resp.Diagnostics, r.client, data)
	} else {
		data.Id = basetypes.NewStringNull()
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End create of resource hyperfabric_vni with id '%s'", data.Id.ValueString()))
}

func (r *VniResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Start read of resource: hyperfabric_vni")
	var data *VniResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Read of resource hyperfabric_vni with id '%s'", data.Id.ValueString()))
	checkAndSetVniIds(data)
	getAndSetVniAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	if data.Id.IsNull() {
		var emptyData *VniResourceModel
		resp.Diagnostics.Append(resp.State.Set(ctx, &emptyData)...)
	} else {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}

	tflog.Debug(ctx, fmt.Sprintf("End read of resource hyperfabric_vni with id '%s'", data.Id.ValueString()))
}

func (r *VniResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Start update of resource: hyperfabric_vni")
	var data *VniResourceModel
	var stateData *VniResourceModel

	// Read Terraform plan and state data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Update of resource hyperfabric_vni with id '%s'", data.Id.ValueString()))

	jsonPayload := getVniJsonPayload(ctx, &resp.Diagnostics, data, "update")

	if resp.Diagnostics.HasError() {
		return
	}

	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/vnis/%s", data.FabricId.ValueString(), data.VniId.ValueString()), "PUT", jsonPayload)

	if resp.Diagnostics.HasError() {
		return
	}

	getAndSetVniAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End update of resource hyperfabric_vni with id '%s'", data.Id.ValueString()))
}

func (r *VniResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Start delete of resource: hyperfabric_vni")
	var data *VniResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Delete of resource hyperfabric_vni with id '%s'", data.Id.ValueString()))
	checkAndSetVniIds(data)
	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/vnis/%s", data.FabricId.ValueString(), data.VniId.ValueString()), "DELETE", nil)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("End delete of resource hyperfabric_vni with id '%s'", data.Id.ValueString()))
}

func (r *VniResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Start import state of resource: hyperfabric_vni")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	var stateData *VniResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &stateData)...)
	tflog.Debug(ctx, fmt.Sprintf("Import state of resource hyperfabric_vni with id '%s'", stateData.Id.ValueString()))
	tflog.Debug(ctx, "End import of state resource: hyperfabric_vni with id")
}

func getAndSetVniAttributes(ctx context.Context, diags *diag.Diagnostics, client *client.Client, data *VniResourceModel) {
	requestData := DoRestRequest(ctx, diags, client, fmt.Sprintf("/api/v1/fabrics/%s/vnis/%s", data.FabricId.ValueString(), data.VniId.ValueString()), "GET", nil)
	if diags.HasError() {
		return
	}

	newVni := *getNewVniResourceModelFromData(data)

	if requestData.Data() != nil {
		attributes := requestData.Data().(map[string]interface{})
		for attributeName, attributeValue := range attributes {
			if attributeName == "fabricId" && (data.FabricId.IsNull() || data.FabricId.IsUnknown() || data.FabricId.ValueString() == "" || data.FabricId.ValueString() != attributeValue.(string)) {
				newVni.FabricId = basetypes.NewStringValue(attributeValue.(string))
				newVni.Id = basetypes.NewStringValue(fmt.Sprintf("%s/vnis/%s", newVni.FabricId.ValueString(), newVni.VniId.ValueString()))
			} else if attributeName == "id" && (data.VniId.IsNull() || data.VniId.IsUnknown() || data.VniId.ValueString() == "" || data.VniId.ValueString() != attributeValue.(string)) {
				newVni.VniId = basetypes.NewStringValue(attributeValue.(string))
				newVni.Id = basetypes.NewStringValue(fmt.Sprintf("%s/vnis/%s", newVni.FabricId.ValueString(), newVni.VniId.ValueString()))
			} else if attributeName == "name" {
				newVni.Name = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "description" {
				newVni.Description = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "enabled" {
				newVni.Enabled = basetypes.NewBoolValue(attributeValue.(bool))
			} else if attributeName == "isDefault" {
				newVni.IsDefault = basetypes.NewBoolValue(attributeValue.(bool))
				// } else if attributeName == "isL3" {
				// 	newVni.IsL3 = basetypes.NewBoolValue(attributeValue.(bool))
			} else if attributeName == "vrfId" && attributeValue.(string) != "" {
				newVni.VrfId = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "vni" {
				newVni.Vni = basetypes.NewFloat64Value(attributeValue.(float64))
			} else if attributeName == "mtu" {
				newVni.Mtu = basetypes.NewFloat64Value(attributeValue.(float64))
			} else if attributeName == "members" {
				stateMembers := make([]MemberResourceModel, 0)
				data.Members.ElementsAs(ctx, &stateMembers, false)
				newVni.Members = NewMembersSet(ctx, &stateMembers, attributeValue.([]interface{}))
			} else if attributeName == "svis" {
				newVni.Svi = NewSviObject(ctx, attributeValue.([]interface{}))
			} else if attributeName == "metadata" {
				newVni.Metadata = NewMetadataObject(ctx, attributeValue.(map[string]interface{}))
			} else if attributeName == "labels" {
				newVni.Labels = NewSetString(ctx, attributeValue.([]interface{}))
			} else if attributeName == "annotations" {
				newVni.Annotations = NewAnnotationsSet(ctx, attributeValue.([]interface{}))
			}
		}
	} else {
		newVni.Id = basetypes.NewStringNull()
	}
	*data = newVni
}

func getVniJsonPayload(ctx context.Context, diags *diag.Diagnostics, data *VniResourceModel, action string) *gabs.Container {
	payloadMap := map[string]interface{}{}
	payloadList := []map[string]interface{}{}

	if !data.FabricId.IsNull() && !data.FabricId.IsUnknown() {
		payloadMap["fabricId"] = data.FabricId.ValueString()
	}

	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		payloadMap["name"] = data.Name.ValueString()
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		payloadMap["description"] = data.Description.ValueString()
	}

	payloadMap["enabled"] = true
	// if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
	// 	payloadMap["enabled"] = data.Enabled.ValueBool()
	// }

	// if !data.IsDefault.IsNull() && !data.IsDefault.IsUnknown() {
	// 	payloadMap["enabled"] = data.IsDefault.ValueBool()
	// }

	// if !data.IsL3.IsNull() && !data.IsL3.IsUnknown() {
	// 	payloadMap["isL3"] = data.IsL3.ValueBool()
	// }

	if !data.VrfId.IsNull() && !data.VrfId.IsUnknown() {
		payloadMap["vrfId"] = data.VrfId.ValueString()
	}

	if !data.Vni.IsNull() && !data.Vni.IsUnknown() && action == "create" {
		payloadMap["vni"] = data.Vni.ValueFloat64()
	}

	if !data.Mtu.IsNull() && !data.Mtu.IsUnknown() {
		payloadMap["mtu"] = data.Mtu.ValueFloat64()
	}

	if !data.Members.IsNull() && !data.Members.IsUnknown() {
		payloadMap["members"] = getMembersJsonPayload(ctx, data.Members)
	}

	if !data.Svi.IsNull() && !data.Svi.IsUnknown() && !IsEmptySingleNestedAttribute(data.Svi.Attributes()) {
		payloadMap["svis"] = getSviJsonPayload(ctx, data.Svi)
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
		payload = map[string]interface{}{"vnis": payloadList}
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

func checkAndSetVniIds(data *VniResourceModel) {
	if strings.Contains(data.Id.ValueString(), "/vnis/") {
		if data.FabricId.IsNull() || data.FabricId.IsUnknown() || data.FabricId.ValueString() == "" || data.VniId.IsNull() || data.VniId.IsUnknown() || data.VniId.ValueString() == "" {
			splitId := strings.Split(data.Id.ValueString(), "/vnis/")
			data.FabricId = basetypes.NewStringValue(splitId[0])
			data.VniId = basetypes.NewStringValue(splitId[1])
		}
	}
}
