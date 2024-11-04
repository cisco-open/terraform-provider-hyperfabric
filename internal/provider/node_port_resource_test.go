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

func TestAccNodePortResource(t *testing.T) {
	fabricName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify provided and default Hyperfabric values.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Port - Create with minimum config and verify provided and default Hyperfabric values.")
				},
				Config:             testNodePortResourceHclConfig(fabricName, "minimal"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "name", "Ethernet1_1"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "roles.#", "1"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "roles.0", "ROUTED_PORT"),
				),
			},
			// Update with all config and verify provided values.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Port - Update with all config and verify provided values.")
				},
				Config:             testNodePortResourceHclConfig(fabricName, "full"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "name", "Ethernet1_1"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "description", "Connected to server01"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "enabled", "true"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "ipv4_addresses.#", "1"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "ipv4_addresses.0", "10.1.0.1/24"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "ipv6_addresses.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "ipv6_addresses.0", "2001:1::1/64"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "ipv6_addresses.1", "2002:1::1/64"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "prevent_forwarding", "true"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "roles.#", "1"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "roles.0", "ROUTED_PORT"),
					resource.TestCheckResourceAttrSet("hyperfabric_node_port.test", "vrf_id"),
				),
			},
			// Update with minimum config and verify config is unchanged.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Port - Update with minimum config and verify config is unchanged.")
				},
				Config:             testNodePortResourceHclConfig(fabricName, "minimal"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "name", "Ethernet1_1"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "description", "Connected to server01"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "enabled", "true"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "ipv4_addresses.#", "1"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "ipv4_addresses.0", "10.1.0.1/24"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "ipv6_addresses.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "ipv6_addresses.0", "2001:1::1/64"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "ipv6_addresses.1", "2002:1::1/64"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "prevent_forwarding", "true"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "roles.#", "1"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "roles.0", "ROUTED_PORT"),
					resource.TestCheckResourceAttrSet("hyperfabric_node_port.test", "vrf_id"),
				),
			},
			// ImportState testing with pre-existing Id.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Port - ImportState testing with pre-existing Id.")
				},
				ResourceName:      "hyperfabric_node_port.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// ImportState testing with fabric and node name.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Port - ImportState testing with fabric, node and interface name.")
				},
				ResourceName:      "hyperfabric_node_port.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fabricName + "/nodes/node1/ports/Ethernet1_1",
			},
			// Update with config containing all optional attributes with empty values and verify config is cleared.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Port - Update with config containing all optional attributes with empty values and verify config is cleared.")
				},
				Config:             testNodePortResourceHclConfig(fabricName, "clear"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "name", "Ethernet1_1"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "roles.#", "1"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "roles.0", "UNUSED_PORT"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "description", ""),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "enabled", ""),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "ipv4_addresses.#", "0"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "ipv6_addresses.#", "0"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "prevent_forwarding", ""),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "vrf_id", ""),
				),
			},
			// Run Plan Only with minimal config and check that plan is empty.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Port - Run Plan Only with minimal config and check that plan is empty.")
				},
				Config:             testNodePortResourceHclConfig(fabricName, "minimal"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "name", "Ethernet1_1"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "roles.#", "1"),
					resource.TestCheckResourceAttr("hyperfabric_node_port.test", "roles.0", "ROUTED_PORT"),
				),
			},
		},
	})
}

func testNodePortResourceHclConfig(fabricName string, configType string) string {
	if configType == "full" {
		return fmt.Sprintf(`
resource "hyperfabric_fabric" "test" {
	name = "%[1]s"
}

resource "hyperfabric_vrf" "test" {
    fabric_id = hyperfabric_fabric.test.id
	name      = "Vrf1"
}

resource "hyperfabric_node" "test" {
    fabric_id   = hyperfabric_fabric.test.id
	name        = "node1"
	model_name  = "HF6100-32D"
    roles       = ["LEAF"]
}

resource "hyperfabric_node_port" "test" {
	node_id            = hyperfabric_node.test.id
	name               = "Ethernet1_1"
	description        = "Connected to server01"
	enabled            = true
	ipv4_addresses     = ["10.1.0.1/24"]
	ipv6_addresses     = ["2001:1::1/64", "2002:1::1/64"]
	prevent_forwarding = true
	roles              = ["ROUTED_PORT"]
	vrf_id             = hyperfabric_vrf.test.vrf_id
}
`, fabricName)
	} else if configType == "clear" {
		return fmt.Sprintf(`
resource "hyperfabric_fabric" "test" {
	name = "%[1]s"
}

resource "hyperfabric_vrf" "test" {
    fabric_id = hyperfabric_fabric.test.id
	name      = "Vrf1"
}

resource "hyperfabric_node" "test" {
	fabric_id   = hyperfabric_fabric.test.id
	name        = "node1"
	model_name  = "HF6100-32D"
    roles       = ["LEAF"]
}

resource "hyperfabric_node_port" "test" {
	node_id            = hyperfabric_node.test.id
	name               = "Ethernet1_1"
	description        = "Connected to server01"
	enabled            = true
	ipv4_addresses     = []
	ipv6_addresses     = []
	prevent_forwarding = false
	roles              = []
	vrf_id             = ""
}
}
`, fabricName)
	} else if configType == "minimal+" {
		return fmt.Sprintf(`
resource "hyperfabric_fabric" "test" {
	name = "%[1]s"
}

resource "hyperfabric_vrf" "test" {
    fabric_id = hyperfabric_fabric.test.id
	name      = "Vrf1"
}

resource "hyperfabric_node" "test" {
	fabric_id  = hyperfabric_fabric.test.id
	name       = "node1"
	model_name = "HF6100-32D"
	roles      = ["LEAF"]
}

resource "hyperfabric_node_port" "test" {
	node_id = hyperfabric_node.test.id
	name    = "Ethernet1_1"
	roles   = ["ROUTED_PORT"]
}
`, fabricName)
	} else {
		return fmt.Sprintf(`
resource "hyperfabric_fabric" "test" {
	name = "%[1]s"
}

resource "hyperfabric_vrf" "test" {
    fabric_id = hyperfabric_fabric.test.id
	name      = "Vrf1"
}

resource "hyperfabric_node" "test" {
    fabric_id  = hyperfabric_fabric.test.id
	name       = "node1"
	model_name = "HF6100-32D"
    roles      = ["LEAF"]
}

resource "hyperfabric_node_port" "test" {
    node_id = hyperfabric_node.test.id
	name    = "Ethernet1_1"
	roles   = ["ROUTED_PORT"]
}
`, fabricName)
	}
}
