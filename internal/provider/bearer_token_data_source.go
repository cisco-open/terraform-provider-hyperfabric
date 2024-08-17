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
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &BearerTokenDataSource{}

func NewBearerTokenDataSource() datasource.DataSource {
	return &BearerTokenDataSource{}
}

// BearerTokenDataSource defines the data source implementation.
type BearerTokenDataSource struct {
	client *client.Client
}

func (d *BearerTokenDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	tflog.Debug(ctx, "Start metadata of datasource: hyperfabric_bearer_token")
	resp.TypeName = req.ProviderTypeName + "_bearer_token"
	tflog.Debug(ctx, "End metadata of datasource: hyperfabric_bearer_token")
}

func (d *BearerTokenDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	tflog.Debug(ctx, "Start schema of datasource: hyperfabric_bearer_token")
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Bearer Token data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "`id` defines the unique identifier of a Bearer Token.",
				Computed:            true,
			},
			"token_id": schema.StringAttribute{
				MarkdownDescription: "`id` defines the unique identifier of a Bearer Token.",
				Computed:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "The JWT token that represent the Bearer Token.",
				Computed:            true,
				Sensitive:           true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Bearer Token.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description is a user defined field to store notes about the Bearer Token.",
				Computed:            true,
			},
			"not_after": schema.StringAttribute{
				MarkdownDescription: "The end date for the validity of the Bearer Token.",
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
			},
			"not_before": schema.StringAttribute{
				MarkdownDescription: "The start date for the validity of the Bearer Token.",
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
			},
			"scope": schema.StringAttribute{
				MarkdownDescription: "The scope assigned to the Bearer Token.",
				Computed:            true,
			},
			"metadata": getMetadataSchemaAttribute(),
			// "labels":   getLabelsDataSourceSchemaAttribute(),
			// "annotations": getAnnotationsDataSourceSchemaAttribute(),
		},
	}
	tflog.Debug(ctx, "End schema of datasource: hyperfabric_bearer_token")
}

func (d *BearerTokenDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of datasource: hyperfabric_bearer_token")
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
	tflog.Debug(ctx, "End configure of datasource: hyperfabric_bearer_token")
}

func (d *BearerTokenDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Start read of datasource: hyperfabric_bearer_token")
	var data *BearerTokenResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create a copy of the Id for when not found during getAndSetBearerTokenAttributes
	cachedId := data.Id.ValueString()
	if cachedId == "" && data.Name.ValueString() != "" {
		data.Id = data.Name
	}

	tflog.Debug(ctx, fmt.Sprintf("Read of datasource hyperfabric_bearer_token with id '%s'", data.Id.ValueString()))

	getAndSetBearerTokenAttributes(ctx, &resp.Diagnostics, d.client, data)

	if data.Id.IsNull() {
		resp.Diagnostics.AddError(
			"Failed to read hyperfabric_bearer_token data source",
			fmt.Sprintf("The hyperfabric_bearer_token data source with id '%s' has not been found", cachedId),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End read of datasource hyperfabric_bearer_token with id '%s'", data.Id.ValueString()))
}
