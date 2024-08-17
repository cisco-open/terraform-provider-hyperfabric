// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/Jeffail/gabs/v2"
	"github.com/cisco-open/terraform-provider-hyperfabric/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func ContainsString(strings []string, matchString string) bool {
	for _, stringValue := range strings {
		if stringValue == matchString {
			return true
		}
	}
	return false
}

func toStringMap(intf interface{}) map[string]string {
	result := make(map[string]string)
	temp := intf.(map[string]interface{})

	for key, value := range temp {
		A(result, key, value.(string))

	}

	return result
}

func StripQuotes(word string) string {
	if strings.HasPrefix(word, "\"") && strings.HasSuffix(word, "\"") {
		return strings.TrimSuffix(strings.TrimPrefix(word, "\""), "\"")
	}
	return word
}

func StripSquareBrackets(word string) string {
	if strings.HasPrefix(word, "[") && strings.HasSuffix(word, "]") {
		return strings.TrimSuffix(strings.TrimPrefix(word, "["), "]")
	}
	return word
}

func A(data map[string]string, key, value string) {
	if value != "" {
		data[key] = value
	}

	if value == "{}" {
		data[key] = ""
	}
}

func G(cont *gabs.Container, key string) string {
	return StripQuotes(cont.S(key).String())
}

func DoRestRequest(ctx context.Context, diags *diag.Diagnostics, restClient *client.Client, path, method string, payload *gabs.Container) *gabs.Container {
	container, err := restClient.DoRestRequest(path, method, payload)
	if err != nil {
		diags.AddError(err.Summary, err.Detail)
		return nil
	}
	return container
}

type setToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate struct{}

func SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate() planmodifier.String {
	return setToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate{}
}

func (m setToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate) Description(_ context.Context) string {
	return "During the update phase, set the value of this attribute to StringNull when the state value is null and the plan value is unknown."
}

func (m setToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate) MarkdownDescription(_ context.Context) string {
	return "During the update phase, set the value of this attribute to StringNull when the state value is null and the plan value is unknown."
}

// Custom plan modifier to set the plan value to null under certain conditions
func (m setToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Set the plan value to StringType null when state value is null and plan value is unknown during an Update
	if !req.State.Raw.IsNull() && req.StateValue.IsNull() && req.PlanValue.IsUnknown() {
		resp.PlanValue = types.StringNull()
	}
}

type setToSetNullWhenStateIsNullPlanIsUnknownDuringUpdate struct{}

func SetToSetNullWhenStateIsNullPlanIsUnknownDuringUpdate() planmodifier.Set {
	return setToSetNullWhenStateIsNullPlanIsUnknownDuringUpdate{}
}

func (m setToSetNullWhenStateIsNullPlanIsUnknownDuringUpdate) Description(_ context.Context) string {
	return "During the update phase, set the value of this attribute to SetNull when the state value is null and the plan value is unknown."
}

func (m setToSetNullWhenStateIsNullPlanIsUnknownDuringUpdate) MarkdownDescription(_ context.Context) string {
	return "During the update phase, set the value of this attribute to SetNull when the state value is null and the plan value is unknown."
}

// Custom plan modifier to set the plan value to null under certain conditions
func (m setToSetNullWhenStateIsNullPlanIsUnknownDuringUpdate) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	// Set the plan value to SetType null when state value is null and plan value is unknown during an Update
	if !req.State.Raw.IsNull() && req.StateValue.IsNull() && req.PlanValue.IsUnknown() {
		resp.PlanValue = types.SetNull(req.StateValue.ElementType(ctx))
	}
}

type setToBoolNullWhenStateIsNullPlanIsUnknownDuringUpdate struct{}

func SetToBoolNullWhenStateIsNullPlanIsUnknownDuringUpdate() planmodifier.Bool {
	return setToBoolNullWhenStateIsNullPlanIsUnknownDuringUpdate{}
}

func (m setToBoolNullWhenStateIsNullPlanIsUnknownDuringUpdate) Description(_ context.Context) string {
	return "During the update phase, set the value of this attribute to BoolNull when the state value is null and the plan value is unknown."
}

func (m setToBoolNullWhenStateIsNullPlanIsUnknownDuringUpdate) MarkdownDescription(_ context.Context) string {
	return "During the update phase, set the value of this attribute to BoolNull when the state value is null and the plan value is unknown."
}

// Custom plan modifier to set the plan value to null under certain conditions
func (m setToBoolNullWhenStateIsNullPlanIsUnknownDuringUpdate) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	// Set the plan value to BoolType null when state value is null and plan value is unknown during an Update
	if !req.State.Raw.IsNull() && req.StateValue.IsNull() && req.PlanValue.IsUnknown() {
		resp.PlanValue = types.BoolNull()
	}
}

type setToFloat64NullWhenStateIsNullPlanIsUnknownDuringUpdate struct{}

func SetToFloat64NullWhenStateIsNullPlanIsUnknownDuringUpdate() planmodifier.Float64 {
	return setToFloat64NullWhenStateIsNullPlanIsUnknownDuringUpdate{}
}

func (m setToFloat64NullWhenStateIsNullPlanIsUnknownDuringUpdate) Description(_ context.Context) string {
	return "During the update phase, set the value of this attribute to Float64Null when the state value is null and the plan value is unknown."
}

func (m setToFloat64NullWhenStateIsNullPlanIsUnknownDuringUpdate) MarkdownDescription(_ context.Context) string {
	return "During the update phase, set the value of this attribute to Float64Null when the state value is null and the plan value is unknown."
}

// Custom plan modifier to set the plan value to null under certain conditions
func (m setToFloat64NullWhenStateIsNullPlanIsUnknownDuringUpdate) PlanModifyFloat64(ctx context.Context, req planmodifier.Float64Request, resp *planmodifier.Float64Response) {
	// Set the plan value to Float64Type null when state value is null and plan value is unknown during an Update
	if !req.State.Raw.IsNull() && req.StateValue.IsNull() && req.PlanValue.IsUnknown() {
		resp.PlanValue = types.Float64Null()
	}
}

type setToObjectNullWhenStateIsNullPlanIsUnknownDuringUpdate struct{}

func SetToObjectNullWhenStateIsNullPlanIsUnknownDuringUpdate() planmodifier.Object {
	return setToObjectNullWhenStateIsNullPlanIsUnknownDuringUpdate{}
}

func (m setToObjectNullWhenStateIsNullPlanIsUnknownDuringUpdate) Description(_ context.Context) string {
	return "During the update phase, set the value of this attribute to ObjectNull when the state value is null and the plan value is unknown."
}

func (m setToObjectNullWhenStateIsNullPlanIsUnknownDuringUpdate) MarkdownDescription(_ context.Context) string {
	return "During the update phase, set the value of this attribute to ObjectNull when the state value is null and the plan value is unknown."
}

// Custom plan modifier to set the plan value to null under certain conditions
func (m setToObjectNullWhenStateIsNullPlanIsUnknownDuringUpdate) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	// Set the plan value to SetType null when state value is null and plan value is unknown during an Update
	if !req.State.Raw.IsNull() && req.StateValue.IsNull() && req.PlanValue.IsUnknown() {
		resp.PlanValue = types.ObjectNull(req.StateValue.AttributeTypes(ctx))
	}
}

// MakeStringRequiredValidator validates that an attribute is not null, as a workaround for when a resource has read-only attributes and nested sets with required attributes
// https://github.com/hashicorp/terraform-plugin-framework/issues/898

var _ validator.String = MakeStringRequiredValidator{}

// MakeStringRequiredValidator validates that an attribute is not null. Most
// attributes should set Required: true instead, however in certain scenarios,
// such as a computed nested attribute, all underlying attributes must also be
// computed for planning to not show unexpected differences.
type MakeStringRequiredValidator struct{}

// Description describes the validation in plain text formatting.
func (v MakeStringRequiredValidator) Description(_ context.Context) string {
	return "is required"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v MakeStringRequiredValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v MakeStringRequiredValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if !req.ConfigValue.IsNull() {
		return
	}

	setName := req.Path.String()[0:strings.Index(req.Path.String(), "[")]
	attributeName := req.Path.String()[strings.Index(req.Path.String(), "]")+2:]

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Incorrect attribute value type",
		fmt.Sprintf("Inappropriate value for attribute \"%s\": attribute \"%s\" is required.", setName, attributeName),
	)
}

// StringNotNull returns an validator which ensures that the string attribute is
// configured. Most attributes should set Required: true instead, however in
// certain scenarios, such as a computed nested attribute, all underlying
// attributes must also be computed for planning to not show unexpected
// differences.
func MakeStringRequired() validator.String {
	return MakeStringRequiredValidator{}
}

// Generic type for Set of Strings and companion functions
func SetStringResourceModelAttributeType() attr.Type {
	return types.StringType
}

func getSetStringJsonPayload(ctx context.Context, data basetypes.SetValue) []string {
	strings := []string{}
	data.ElementsAs(ctx, &strings, false)
	return strings
}

func NewSetString(ctx context.Context, data []interface{}) basetypes.SetValue {
	stringsSet, _ := types.SetValueFrom(ctx, SetStringResourceModelAttributeType(), data)
	return stringsSet
}

func IsEmptySingleNestedAttribute(attributes map[string]attr.Value) bool {
	for _, value := range attributes {
		if !value.IsNull() {
			return false
		}
	}
	return true
}
