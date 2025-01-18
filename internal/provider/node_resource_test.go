// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNodeResource(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	fabricName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify provided and default Hyperfabric values.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node - Create with minimum config and verify provided and default Hyperfabric values.")
				},
				Config:             testNodeResourceHclConfig(fabricName, name, "minimal"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node.test", "name", name),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "model_name", "HF6100-32D"),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "roles.#", "1"),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "roles.0", "LEAF"),
				),
			},
			// Update with all config and verify provided values.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node - Update with all config and verify provided values.")
				},
				Config:             testNodeResourceHclConfig(fabricName, name, "full"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node.test", "name", name),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "model_name", "HF6100-32D"),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "roles.#", "1"),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "roles.0", "LEAF"),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "description", "This node is powered by Cisco Nexus Hyperfabric"),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "enabled", "true"),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "location", "sj01-1-101-AAA01"),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "labels.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "annotations.#", "2"),
				),
			},
			// Update with minimum config and verify config is unchanged.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node - Update with minimum config and verify config is unchanged.")
				},
				Config:             testNodeResourceHclConfig(fabricName, name, "minimal"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node.test", "name", name),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "model_name", "HF6100-32D"),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "roles.#", "1"),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "roles.0", "LEAF"),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "description", "This node is powered by Cisco Nexus Hyperfabric"),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "enabled", "true"),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "location", "sj01-1-101-AAA01"),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "labels.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "annotations.#", "2"),
				),
			},
			// ImportState testing with pre-existing Id.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node - ImportState testing with pre-existing Id.")
				},
				ResourceName:      "hyperfabric_node.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// ImportState testing with name.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node - ImportState testing with name.")
				},
				ResourceName:      "hyperfabric_node.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fabricName + "/nodes/" + name,
			},
			// Update with config containing all optional attributes with empty values and verify config is cleared.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node - Update with config containing all optional attributes with empty values and verify config is cleared.")
				},
				Config:             testNodeResourceHclConfig(fabricName, name, "clear"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node.test", "name", name),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "model_name", "HF6100-32D"),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "roles.#", "1"),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "roles.0", "LEAF"),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "enabled", "true"),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "description", ""),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "location", ""),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "labels.#", "0"),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "annotations.#", "0"),
				),
			},
			// Run Plan Only with minimal config and check that plan is empty.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node - Run Plan Only with minimal config and check that plan is empty.")
				},
				Config:             testNodeResourceHclConfig(fabricName, name, "minimal"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node.test", "name", name),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "model_name", "HF6100-32D"),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "roles.#", "1"),
					resource.TestCheckResourceAttr("hyperfabric_node.test", "roles.0", "LEAF"),
				),
			},
		},
	})
}

func testNodeResourceHclConfig(fabricName string, name string, configType string) string {
	if configType == "full" {
		return fmt.Sprintf(`
resource "hyperfabric_fabric" "test" {
	name = "%[1]s"
}

resource "hyperfabric_node" "test" {
    fabric_id   = hyperfabric_fabric.test.id
	name        = "%[2]s"
	model_name  = "HF6100-32D"
    roles       = ["LEAF"]
	description = "This node is powered by Cisco Nexus Hyperfabric"
	location    = "sj01-1-101-AAA01"
	labels = [
		"sj01-1-101-AAA01",
		"blue"
	]
	annotations = [
		{
			name      = "color"
			value     = "blue"
		},
		{
			data_type = "UINT32"
			name  = "rack"
			value = "1"
		}
	]
}
`, fabricName, name)
	} else if configType == "clear" {
		return fmt.Sprintf(`
resource "hyperfabric_fabric" "test" {
	name = "%[1]s"
}

resource "hyperfabric_node" "test" {
	fabric_id   = hyperfabric_fabric.test.id
	name        = "%[2]s"
	model_name  = "HF6100-32D"
    roles       = ["LEAF"]
	description = ""
	location    = ""
	labels = []
	annotations = []
}
`, fabricName, name)
	} else {
		return fmt.Sprintf(`
resource "hyperfabric_fabric" "test" {
	name = "%[1]s"
}
resource "hyperfabric_node" "test" {
    fabric_id = hyperfabric_fabric.test.id
	name = "%[2]s"
	model_name = "HF6100-32D"
}
`, fabricName, name)
	}
}
