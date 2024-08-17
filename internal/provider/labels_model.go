// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getLabelsSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A set of user-defined labels for searching and locating objects.`,
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.UseStateForUnknown(),
		},
		ElementType: types.StringType,
	}
}

func getLabelsDataSourceSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A set of user-defined labels for searching and locating objects.`,
		Computed:            true,
		ElementType:         types.StringType,
	}
}
