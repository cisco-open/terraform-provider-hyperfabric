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
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &BearerTokenResource{}
var _ resource.ResourceWithImportState = &BearerTokenResource{}

func NewBearerTokenResource() resource.Resource {
	return &BearerTokenResource{}
}

// BearerTokenResource defines the resource implementation.
type BearerTokenResource struct {
	client *client.Client
}

// BearerTokenResourceModel describes the resource data model.
type BearerTokenResourceModel struct {
	Id          types.String      `tfsdk:"id"`
	TokenId     types.String      `tfsdk:"token_id"`
	Name        types.String      `tfsdk:"name"`
	Description types.String      `tfsdk:"description"`
	NotAfter    timetypes.RFC3339 `tfsdk:"not_after"`
	NotBefore   timetypes.RFC3339 `tfsdk:"not_before"`
	Scope       types.String      `tfsdk:"scope"`
	Token       types.String      `tfsdk:"token"`
	Metadata    types.Object      `tfsdk:"metadata"`
	// Labels      types.Set    `tfsdk:"labels"`
	// Annotations types.Set    `tfsdk:"annotations"`
}

func getEmptyBearerTokenResourceModel() *BearerTokenResourceModel {
	return &BearerTokenResourceModel{
		Id:          basetypes.NewStringNull(),
		TokenId:     basetypes.NewStringNull(),
		Name:        basetypes.NewStringNull(),
		Description: basetypes.NewStringNull(),
		NotAfter:    timetypes.NewRFC3339Null(),
		NotBefore:   timetypes.NewRFC3339Null(),
		Scope:       basetypes.NewStringValue("ADMIN"),
		Token:       basetypes.NewStringNull(),
		Metadata:    basetypes.NewObjectNull(MetadataResourceModelAttributeType()),
		// Labels:      basetypes.NewSetNull(SetStringResourceModelAttributeType()),
		// Annotations: basetypes.NewSetNull(AnnotationResourceModelAttributeType()),
	}
}

func getNewBearerTokenResourceModelFromData(data *BearerTokenResourceModel) *BearerTokenResourceModel {
	newBearerToken := getEmptyBearerTokenResourceModel()

	if !data.Id.IsNull() && !data.Id.IsUnknown() {
		newBearerToken.Id = data.Id
	}

	if !data.TokenId.IsNull() && !data.TokenId.IsUnknown() {
		newBearerToken.TokenId = data.TokenId
	}

	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		newBearerToken.Name = data.Name
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		newBearerToken.Description = data.Description
	}

	if !data.NotAfter.IsNull() && !data.NotAfter.IsUnknown() {
		newBearerToken.NotAfter = data.NotAfter
	}

	if !data.NotBefore.IsNull() && !data.NotBefore.IsUnknown() {
		newBearerToken.NotBefore = data.NotBefore
	}

	if !data.Scope.IsNull() && !data.Scope.IsUnknown() {
		newBearerToken.Scope = data.Scope
	}

	if !data.Token.IsNull() && !data.Token.IsUnknown() {
		newBearerToken.Token = data.Token
	}

	if !data.Metadata.IsNull() && !data.Metadata.IsUnknown() {
		newBearerToken.Metadata = data.Metadata
	}

	return newBearerToken
}

type BearerTokenIdentifier struct {
	Id types.String
}

func (r *BearerTokenResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Debug(ctx, "Start metadata of resource: hyperfabric_bearer_token")
	resp.TypeName = req.ProviderTypeName + "_bearer_token"
	tflog.Debug(ctx, "End metadata of resource: hyperfabric_bearer_token")
}

func (r *BearerTokenResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Debug(ctx, "Start schema of resource: hyperfabric_bearer_token")
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Bearer Token resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "`id` defines the unique identifier of a Bearer Token.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"token_id": schema.StringAttribute{
				MarkdownDescription: "`id` defines the unique identifier of a Bearer Token.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "The JWT token that represent the Bearer Token.",
				Computed:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Bearer Token.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description is a user defined field to store notes about the Bearer Token.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"not_after": schema.StringAttribute{
				MarkdownDescription: "The end date for the validity of the Bearer Token.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
				CustomType: timetypes.RFC3339Type{},
			},
			"not_before": schema.StringAttribute{
				MarkdownDescription: "The start date for the validity of the Bearer Token.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
				CustomType: timetypes.RFC3339Type{},
			},
			"scope": schema.StringAttribute{
				MarkdownDescription: "The scope assigned to the Bearer Token.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
				Default: stringdefault.StaticString("READ_ONLY"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"ADMIN", "READ_WRITE", "READ_ONLY"}...),
				},
			},
			"metadata": getMetadataSchemaAttribute(),
			// "labels":   getLabelsSchemaAttribute(),
			// "annotations": getAnnotationsSchemaAttribute(),
		},
	}
	tflog.Debug(ctx, "End schema of resource: hyperfabric_bearer_token")
}

func (r *BearerTokenResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of resource: hyperfabric_bearer_token")
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
	tflog.Debug(ctx, "End configure of resource: hyperfabric_bearer_token")
}

func (r *BearerTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Start create of resource: hyperfabric_bearer_token")

	var data *BearerTokenResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Create of resource hyperfabric_bearer_token with name '%s'", data.Name.ValueString()))

	jsonPayload := getBearerTokenJsonPayload(ctx, &resp.Diagnostics, data, "create")
	if resp.Diagnostics.HasError() {
		return
	}

	container := DoRestRequest(ctx, &resp.Diagnostics, r.client, "/api/v1/bearerTokens", "POST", jsonPayload)
	if resp.Diagnostics.HasError() {
		return
	}

	bearerTokensContainer, err := container.ArrayElement(0, "tokens")
	if err != nil {
		return
	}

	bearertokenId := StripQuotes(bearerTokensContainer.Search("tokenId").String())
	if bearertokenId != "" {
		data.Id = basetypes.NewStringValue(bearertokenId)
		data.TokenId = basetypes.NewStringValue(bearertokenId)
		token := StripQuotes(container.Search("token").String())
		if token != "" {
			data.Token = basetypes.NewStringValue(token)
		}
		getAndSetBearerTokenAttributes(ctx, &resp.Diagnostics, r.client, data)
	} else {
		data.Id = basetypes.NewStringNull()
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End create of resource hyperfabric_bearer_token with name '%s'", data.Name.ValueString()))
}

func (r *BearerTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Start read of resource: hyperfabric_bearer_token")
	var data *BearerTokenResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Read of resource hyperfabric_bearer_token with id '%s'", data.Id.ValueString()))

	getAndSetBearerTokenAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	if data.Id.IsNull() {
		var emptyData *BearerTokenResourceModel
		resp.Diagnostics.Append(resp.State.Set(ctx, &emptyData)...)
	} else {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}

	tflog.Debug(ctx, fmt.Sprintf("End read of resource hyperfabric_bearer_token with id '%s'", data.Id.ValueString()))
}

func (r *BearerTokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Start update of resource: hyperfabric_bearer_token")
	var data *BearerTokenResourceModel
	var stateData *BearerTokenResourceModel

	// Read Terraform plan and state data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Update of resource hyperfabric_bearer_token with id '%s'", data.Id.ValueString()))

	// jsonPayload := getBearerTokenJsonPayload(ctx, &resp.Diagnostics, data, "update")

	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/bearerTokens/%s", data.Id.ValueString()), "PUT", jsonPayload)

	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// getAndSetBearerTokenAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End update of resource hyperfabric_bearer_token with id '%s'", data.Id.ValueString()))
}

func (r *BearerTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Start delete of resource: hyperfabric_bearer_token")
	var data *BearerTokenResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Delete of resource hyperfabric_bearer_token with id '%s'", data.Id.ValueString()))
	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/bearerTokens/%s", data.Id.ValueString()), "DELETE", nil)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("End delete of resource hyperfabric_bearer_token with id '%s'", data.Id.ValueString()))
}

func (r *BearerTokenResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Start import state of resource: hyperfabric_bearer_token")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	var stateData *BearerTokenResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &stateData)...)
	tflog.Debug(ctx, fmt.Sprintf("Import state of resource hyperfabric_bearer_token with id '%s'", stateData.Id.ValueString()))
	tflog.Debug(ctx, "End import of state resource: hyperfabric_bearer_token")
}

func getAndSetBearerTokenAttributes(ctx context.Context, diags *diag.Diagnostics, client *client.Client, data *BearerTokenResourceModel) {
	requestData := DoRestRequest(ctx, diags, client, fmt.Sprintf("/api/v1/bearerTokens/%s", data.Id.ValueString()), "GET", nil)
	if diags.HasError() {
		return
	}

	newBearerToken := *getNewBearerTokenResourceModelFromData(data)
	// newBearerToken.Id = data.Id
	// newBearerToken.Name = data.Name
	// newBearerToken.Token = data.Token
	// newBearerToken.TokenId = data.TokenId

	if requestData.Data() != nil {
		attributes := requestData.Data().(map[string]interface{})
		for attributeName, attributeValue := range attributes {
			if attributeName == "tokenId" && (data.Id.IsNull() || data.Id.IsUnknown() || data.Id.ValueString() == "" || data.Id.ValueString() == attributeValue.(string)) {
				newBearerToken.Id = basetypes.NewStringValue(attributeValue.(string))
				newBearerToken.TokenId = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "name" {
				newBearerToken.Name = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "description" {
				newBearerToken.Description = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "notAfter" {
				timeValue, err := timetypes.NewRFC3339Value(attributeValue.(string))
				if err == nil {
					newBearerToken.NotAfter = timeValue
				}
			} else if attributeName == "notBefore" {
				timeValue, err := timetypes.NewRFC3339Value(attributeValue.(string))
				if err == nil {
					newBearerToken.NotBefore = timeValue
				}
			} else if attributeName == "scope" {
				newBearerToken.Scope = basetypes.NewStringValue(strings.Split(attributeValue.(string), "TOKEN_SCOPE_")[1])
			} else if attributeName == "token" {
				newBearerToken.Token = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "metadata" {
				newBearerToken.Metadata = NewMetadataObject(ctx, attributeValue.(map[string]interface{}))
				// } else if attributeName == "labels" {
				// 	newBearerToken.Labels = NewSetString(ctx, attributeValue.([]interface{}))
				// } else if attributeName == "annotations" {
				// 	newBearerToken.Annotations = NewAnnotationsSet(ctx, attributeValue.([]interface{}))
			}
		}
	} else {
		newBearerToken.Id = basetypes.NewStringNull()
	}
	*data = newBearerToken
}

func getBearerTokenJsonPayload(ctx context.Context, diags *diag.Diagnostics, data *BearerTokenResourceModel, action string) *gabs.Container {
	payloadMap := map[string]interface{}{}
	payloadList := []map[string]interface{}{}

	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		payloadMap["name"] = data.Name.ValueString()
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		payloadMap["description"] = data.Description.ValueString()
	}

	if !data.NotAfter.IsNull() && !data.NotAfter.IsUnknown() {
		payloadMap["notAfter"] = data.NotAfter.ValueString()
	}

	if !data.NotBefore.IsNull() && !data.NotBefore.IsUnknown() {
		payloadMap["notBefore"] = data.NotBefore.ValueString()
	}

	if !data.Scope.IsNull() && !data.Scope.IsUnknown() {
		payloadMap["scope"] = fmt.Sprintf("TOKEN_SCOPE_%s", data.Scope.ValueString())
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
		payload = map[string]interface{}{"tokens": payloadList}
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
