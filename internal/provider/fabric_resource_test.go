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

func TestAccFabricResource(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify provided and default Hyperfabric values.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Fabric - Create with minimum config and verify provided and default Hyperfabric values.")
				},
				Config:             testFabricResourceHclConfig(name, "", "minimal"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "name", name),
				),
			},
			// Update with all config and verify provided values.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Fabric - Update with all config and verify provided values.")
				},
				Config:             testFabricResourceHclConfig(name, "", "full"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "name", name),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "description", "This fabric is powered by Cisco Nexus Hyperfabric"),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "address", "170 West Tasman Dr."),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "city", "San Jose"),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "country", "USA"),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "location", "sj01-1-101-AAA01"),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "labels.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "annotations.#", "2"),
				),
			},
			// Update with minimum config and verify config is unchanged.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Fabric - Update with minimum config and verify config is unchanged.")
				},
				Config:             testFabricResourceHclConfig(name, "", "minimal"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "name", name),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "description", "This fabric is powered by Cisco Nexus Hyperfabric"),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "topology", "MESH"),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "address", "170 West Tasman Dr."),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "city", "San Jose"),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "country", "USA"),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "location", "sj01-1-101-AAA01"),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "labels.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "annotations.#", "2"),
				),
			},
			// Update with minimum config and change topology to SPINE_LEAF.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Fabric - Update with minimum config and change topology to SPINE_LEAF.")
				},
				Config:             testFabricResourceHclConfig(name, "SPINE_LEAF", "minimal"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "name", name),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "description", "This fabric is powered by Cisco Nexus Hyperfabric"),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "topology", "SPINE_LEAF"),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "address", "170 West Tasman Dr."),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "city", "San Jose"),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "country", "USA"),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "location", "sj01-1-101-AAA01"),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "labels.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "annotations.#", "2"),
				),
			},
			// ImportState testing with pre-existing Id.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Fabric - ImportState testing with pre-existing Id.")
				},
				ResourceName:      "hyperfabric_fabric.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// ImportState testing with name.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Fabric - ImportState testing with name.")
				},
				ResourceName:      "hyperfabric_fabric.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     name,
			},
			// Update with config containing all optional attributes with empty values and verify config is cleared.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Fabric - Update with config containing all optional attributes with empty values and verify config is cleared.")
				},
				Config:             testFabricResourceHclConfig(name, "", "clear"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "name", name),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "description", ""),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "address", ""),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "city", ""),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "country", ""),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "location", ""),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "labels.#", "0"),
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "annotations.#", "0"),
				),
			},
			// Run Plan Only with minimal config and check that plan is empty.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Fabric - Run Plan Only with minimal config and check that plan is empty.")
				},
				Config:             testFabricResourceHclConfig(name, "", "minimal"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_fabric.test", "name", name),
				),
			},
		},
	})
}

func testFabricResourceHclConfig(name string, topology string, configType string) string {
	topologyConfigLine := ""
	if topology != "" {
		topologyConfigLine = fmt.Sprintf("topology = %v", topology)
	}

	if configType == "full" {
		return fmt.Sprintf(`
resource "hyperfabric_fabric" "test" {
	name = "%[1]s"
	description = "This fabric is powered by Cisco Nexus Hyperfabric"
	%[2]s
	address     = "170 West Tasman Dr."
	city        = "San Jose"
	country     = "USA"
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
`, name, topologyConfigLine)
	} else if configType == "clear" {
		return fmt.Sprintf(`
resource "hyperfabric_fabric" "test" {
	name = "%[1]s"
	description = ""
	address     = ""
	city        = ""
	country     = ""
	location    = ""
	labels = []
	annotations = []
}
`, name)
	} else {
		return fmt.Sprintf(`
resource "hyperfabric_fabric" "test" {
	name = "%[1]s"
}
`, name)
	}
}
