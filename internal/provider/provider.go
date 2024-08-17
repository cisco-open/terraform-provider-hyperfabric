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
	"os"
	"strconv"
	"strings"

	"github.com/cisco-open/terraform-provider-hyperfabric/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var globalLabel string

// Ensure HyperfabricProvider satisfies various provider interfaces.
var _ provider.Provider = &HyperfabricProvider{}
var _ provider.ProviderWithFunctions = &HyperfabricProvider{}

// HyperfabricProvider defines the provider implementation.
type HyperfabricProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// HyperfabricProviderModel describes the provider data model.
type HyperfabricProviderModel struct {
	IsInsecure types.Bool   `tfsdk:"insecure"`
	Label      types.String `tfsdk:"label"`
	MaxRetries types.Int32  `tfsdk:"retries"`
	ProxyUrl   types.String `tfsdk:"proxy_url"`
	ProxyCreds types.String `tfsdk:"proxy_creds"`
	Token      types.String `tfsdk:"token"`
	URL        types.String `tfsdk:"url"`
}

func (p *HyperfabricProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "hyperfabric"
	resp.Version = p.version
}

func (p *HyperfabricProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"insecure": schema.BoolAttribute{
				MarkdownDescription: "Allow insecure HTTPS client. This can also be set as the HYPERFABRIC_INSECURE environment variable. Defaults to `true`.",
				Optional:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Global label for the provider. This can also be set as the HYPERFABRIC_LABEL environment variable. Defaults to `terraform`.",
				Optional:            true,
			},
			"proxy_url": schema.StringAttribute{
				MarkdownDescription: "Proxy Server URL with port number. This can also be set as the HYPERFABRIC_PROXY_URL environment variable.",
				Optional:            true,
			},
			"proxy_creds": schema.StringAttribute{
				MarkdownDescription: "Proxy server credentials in the form of username:password. This can also be set as the HYPERFABRIC_PROXY_CREDS environment variable.",
				Optional:            true,
			},
			"retries": schema.Int32Attribute{
				MarkdownDescription: "Number of retries for REST API calls. This can also be set as the HYPERFABRIC_RETRIES environment variable. Defaults to `2`.",
				Optional:            true,
				Validators: []validator.Int32{
					int32validator.Between(0, 10),
				},
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "API token of user in a the Hyperfabric service organization. This can also be set as the HYPERFABRIC_TOKEN environment variable.",
				Sensitive:           true,
				Optional:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL of the Hyperfabric service. This can also be set as the HYPERFABRIC_URL environment variable. Defaults to `https://hyperfabric.cisco.com`.",
				Optional:            true,
			},
		},
	}
}

func getStringAttribute(attribute basetypes.StringValue, envKey string, defaultValue string) string {
	if attribute.IsNull() {
		envValue, found := os.LookupEnv(envKey)
		if found {
			return envValue
		} else {
			return defaultValue
		}
	}
	return attribute.ValueString()
}

func getBoolAttribute(attribute basetypes.BoolValue, envKey string, defaultValue bool) bool {
	if attribute.IsNull() {
		envValue, err := strconv.ParseBool(os.Getenv(envKey))
		if err != nil {
			return defaultValue
		}
		return envValue
	}
	return attribute.ValueBool()
}

func getIntAttribute(attribute basetypes.Int32Value, envKey string, defaultValue int) int {
	if attribute.IsNull() {
		envValue, err := strconv.ParseInt(os.Getenv(envKey), 10, 32)
		if err != nil {
			return defaultValue
		}
		return int(envValue)
	}
	return int(attribute.ValueInt32())
}

func (p *HyperfabricProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data HyperfabricProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := getStringAttribute(data.URL, "HYPERFABRIC_URL", "https://hyperfabric.cisco.com")
	if !strings.HasPrefix(url, "https://") {
		resp.Diagnostics.AddError(
			"Incorrect URL prefix",
			fmt.Sprintf("URL '%s' must start with 'https://'", url),
		)
	}

	token := getStringAttribute(data.Token, "HYPERFABRIC_TOKEN", "")
	insecure := getBoolAttribute(data.IsInsecure, "HYPERFABRIC_INSECURE", false)
	maxRetries := getIntAttribute(data.MaxRetries, "HYPERFABRIC_RETRIES", 2)
	proxyCreds := getStringAttribute(data.ProxyCreds, "HYPERFABRIC_PROXY_CREDS", "")
	proxyUrl := getStringAttribute(data.ProxyUrl, "HYPERFABRIC_PROXY_URL", "")
	globalLabel = getStringAttribute(data.Label, "HYPERFABRIC_LABEL", "terraform")

	// Client configuration for data sources and resources
	hyperfabricClient := client.GetClient(url, token, client.Insecure(insecure), client.ProxyUrl(proxyUrl), client.ProxyCreds(proxyCreds), client.MaxRetries(maxRetries))
	resp.DataSourceData = hyperfabricClient
	resp.ResourceData = hyperfabricClient
}

func (p *HyperfabricProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewBearerTokenResource,
		NewFabricResource,
		NewNodeResource,
		NewNodeManagementPortResource,
		NewNodePortResource,
		NewConnectionResource,
		NewBindToNodeResource,
		NewUserResource,
		NewVrfResource,
		NewVniResource,
	}
}

func (p *HyperfabricProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewBearerTokenDataSource,
		NewDeviceDataSource,
		NewFabricDataSource,
		NewNodeDataSource,
		NewNodeManagementPortDataSource,
		NewNodePortDataSource,
		NewUserDataSource,
		NewVrfDataSource,
		NewVniDataSource,
	}
}

func (p *HyperfabricProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		// NewExampleFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &HyperfabricProvider{
			version: version,
		}
	}
}
