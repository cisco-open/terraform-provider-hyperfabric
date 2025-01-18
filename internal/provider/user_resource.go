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

	"github.com/Jeffail/gabs/v2"
	"github.com/cisco-open/terraform-provider-hyperfabric/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &UserResource{}
var _ resource.ResourceWithImportState = &UserResource{}

func NewUserResource() resource.Resource {
	return &UserResource{}
}

// UserResource defines the resource implementation.
type UserResource struct {
	client *client.Client
}

// UserResourceModel describes the resource data model.
type UserResourceModel struct {
	Id        types.String `tfsdk:"id"`
	Email     types.String `tfsdk:"email"`
	LastLogin types.String `tfsdk:"last_login"`
	Enabled   types.Bool   `tfsdk:"enabled"`
	Provider  types.String `tfsdk:"auth_provider"`
	Role      types.String `tfsdk:"role"`
	Metadata  types.Object `tfsdk:"metadata"`
	Labels    types.Set    `tfsdk:"labels"`
	// Annotations types.Set    `tfsdk:"annotations"`
}

func getEmptyUserResourceModel() *UserResourceModel {
	return &UserResourceModel{
		Id:        basetypes.NewStringNull(),
		Email:     basetypes.NewStringNull(),
		LastLogin: basetypes.NewStringNull(),
		Enabled:   basetypes.NewBoolValue(false),
		Provider:  basetypes.NewStringNull(),
		Role:      basetypes.NewStringNull(),
		Metadata:  basetypes.NewObjectNull(MetadataResourceModelAttributeType()),
		Labels:    basetypes.NewSetNull(SetStringResourceModelAttributeType()),
		// Annotations: basetypes.NewSetNull(AnnotationResourceModelAttributeType()),
	}
}

func getNewUserResourceModelFromData(data *UserResourceModel) *UserResourceModel {
	newUser := getEmptyUserResourceModel()

	if !data.Id.IsNull() && !data.Id.IsUnknown() {
		newUser.Id = data.Id
	}

	if !data.Email.IsNull() && !data.Email.IsUnknown() {
		newUser.Email = data.Email
	}

	if !data.LastLogin.IsNull() && !data.LastLogin.IsUnknown() {
		newUser.LastLogin = data.LastLogin
	}

	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		newUser.Enabled = data.Enabled
	}

	if !data.Provider.IsNull() && !data.Provider.IsUnknown() {
		newUser.Provider = data.Provider
	}

	if !data.Role.IsNull() && !data.Role.IsUnknown() {
		newUser.Role = data.Role
	}

	if !data.Metadata.IsNull() && !data.Metadata.IsUnknown() {
		newUser.Metadata = data.Metadata
	}

	if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
		newUser.Labels = data.Labels
	}
	return newUser
}

type UserIdentifier struct {
	Id types.String
}

func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Debug(ctx, "Start metadata of resource: hyperfabric_user")
	resp.TypeName = req.ProviderTypeName + "_user"
	tflog.Debug(ctx, "End metadata of resource: hyperfabric_user")
}

func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Debug(ctx, "Start schema of resource: hyperfabric_user")
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "User resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "`id` defines the unique identifier of a User.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The email of the User.",
				Required:            true,
			},
			"last_login": schema.StringAttribute{
				MarkdownDescription: "The last time the User logged into the application.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "The enabled state of the User.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
				Default: booldefault.StaticBool(true),
			},
			"auth_provider": schema.StringAttribute{
				MarkdownDescription: "The authentication provider for the User.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "The role assigned to the User.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
				Default: stringdefault.StaticString("READ_ONLY"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"ADMIN", "READ_WRITE", "READ_ONLY"}...),
				},
			},
			"metadata": getMetadataSchemaAttribute(),
			"labels":   getLabelsSchemaAttribute(),
			// "annotations": getAnnotationsSchemaAttribute(),
		},
	}
	tflog.Debug(ctx, "End schema of resource: hyperfabric_user")
}

func (r *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of resource: hyperfabric_user")
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
	tflog.Debug(ctx, "End configure of resource: hyperfabric_user")
}

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Start create of resource: hyperfabric_user")

	var data *UserResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Create of resource hyperfabric_user with email '%s'", data.Email.ValueString()))

	jsonPayload := getUserJsonPayload(ctx, &resp.Diagnostics, data, "create")
	if resp.Diagnostics.HasError() {
		return
	}

	container := DoRestRequest(ctx, &resp.Diagnostics, r.client, "/api/v1/users", "POST", jsonPayload)
	if resp.Diagnostics.HasError() {
		return
	}

	userContainer, err := container.ArrayElement(0, "users")
	if err != nil {
		return
	}

	userId := StripQuotes(userContainer.Search("id").String())
	if userId != "" {
		data.Id = basetypes.NewStringValue(userId)
		getAndSetUserAttributes(ctx, &resp.Diagnostics, r.client, data)
	} else {
		data.Id = basetypes.NewStringNull()
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End create of resource hyperfabric_user with email '%s'", data.Email.ValueString()))
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Start read of resource: hyperfabric_user")
	var data *UserResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Read of resource hyperfabric_user with id '%s'", data.Id.ValueString()))

	getAndSetUserAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	if data.Id.IsNull() {
		var emptyData *UserResourceModel
		resp.Diagnostics.Append(resp.State.Set(ctx, &emptyData)...)
	} else {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}

	tflog.Debug(ctx, fmt.Sprintf("End read of resource hyperfabric_user with id '%s'", data.Id.ValueString()))
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Start update of resource: hyperfabric_user")
	var data *UserResourceModel
	var stateData *UserResourceModel

	// Read Terraform plan and state data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Update of resource hyperfabric_user with id '%s'", data.Id.ValueString()))

	jsonPayload := getUserJsonPayload(ctx, &resp.Diagnostics, data, "update")

	if resp.Diagnostics.HasError() {
		return
	}

	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/users/%s", data.Id.ValueString()), "PUT", jsonPayload)

	if resp.Diagnostics.HasError() {
		return
	}

	getAndSetUserAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End update of resource hyperfabric_user with id '%s'", data.Id.ValueString()))
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Start delete of resource: hyperfabric_user")
	var data *UserResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Delete of resource hyperfabric_user with id '%s'", data.Id.ValueString()))
	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/users/%s", data.Id.ValueString()), "DELETE", nil)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("End delete of resource hyperfabric_user with id '%s'", data.Id.ValueString()))
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Start import state of resource: hyperfabric_user")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	var stateData *UserResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &stateData)...)
	tflog.Debug(ctx, fmt.Sprintf("Import state of resource hyperfabric_user with id '%s'", stateData.Id.ValueString()))
	tflog.Debug(ctx, "End import of state resource: hyperfabric_user")
}

func getAndSetUserAttributes(ctx context.Context, diags *diag.Diagnostics, client *client.Client, data *UserResourceModel) {
	requestData := DoRestRequest(ctx, diags, client, fmt.Sprintf("/api/v1/users/%s", data.Id.ValueString()), "GET", nil)
	// requestData := DoRestRequest(ctx, diags, client, "/api/v1/users", "GET", nil)
	if diags.HasError() {
		return
	}

	newUser := *getNewUserResourceModelFromData(data)
	// newUser.Id = data.Id
	// newUser.Email = data.Email

	if requestData.Data() != nil {
		attributes := requestData.Data().(map[string]interface{})
		for attributeName, attributeValue := range attributes {
			if attributeName == "id" && (data.Id.IsNull() || data.Id.IsUnknown() || data.Id.ValueString() == "" || data.Id.ValueString() != attributeValue.(string)) {
				newUser.Id = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "email" {
				newUser.Email = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "provider" {
				newUser.Provider = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "lastLogin" {
				newUser.LastLogin = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "enabled" {
				newUser.Enabled = basetypes.NewBoolValue(attributeValue.(bool))
			} else if attributeName == "role" {
				newUser.Role = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "metadata" {
				newUser.Metadata = NewMetadataObject(ctx, attributeValue.(map[string]interface{}))
			} else if attributeName == "labels" {
				newUser.Labels = NewSetString(ctx, attributeValue.([]interface{}))
				// } else if attributeName == "annotations" {
				// 	newUser.Annotations = NewAnnotationsSet(ctx, attributeValue.([]interface{}))
			}
		}
	} else {
		newUser.Id = basetypes.NewStringNull()
	}
	*data = newUser
}

func getUserJsonPayload(ctx context.Context, diags *diag.Diagnostics, data *UserResourceModel, action string) *gabs.Container {
	payloadMap := map[string]interface{}{}
	payloadList := []map[string]interface{}{}

	if !data.Email.IsNull() && !data.Email.IsUnknown() && action == "create" {
		payloadMap["email"] = data.Email.ValueString()
	}

	if !data.Role.IsNull() && !data.Role.IsUnknown() {
		payloadMap["role"] = data.Role.ValueString()
	}

	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		payloadMap["enabled"] = data.Enabled.ValueBool()
	}

	if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
		payloadMap["labels"] = getSetStringJsonPayload(ctx, data.Labels)
	}

	// if !data.Annotations.IsNull() && !data.Annotations.IsUnknown() {
	// 	payloadMap["annotations"] = getAnnotationsJsonPayload(ctx, data.Annotations)
	// }

	var payload map[string]interface{}
	if action == "create" {
		payloadList = append(payloadList, payloadMap)
		payload = map[string]interface{}{"users": payloadList}
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
