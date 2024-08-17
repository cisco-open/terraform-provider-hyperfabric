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
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccConnectionResource(t *testing.T) {
	fabricName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify provided and default Hyperfabric values.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Connection - Create with minimum config and verify provided and default Hyperfabric values.")
				},
				Config:             testConnectionResourceHclConfig(fabricName, "minimal"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "local.port_name", "Ethernet1_1"),
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "local.node_name", "node1"),
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "remote.port_name", "Ethernet1_1"),
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "remote.node_name", "node2"),
				),
			},
			// Update with all config and verify provided values.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Connection - Update with all config and verify provided values.")
				},
				Config:             testConnectionResourceHclConfig(fabricName, "full"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "local.port_name", "Ethernet1_1"),
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "local.node_name", "node1"),
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "remote.port_name", "Ethernet1_1"),
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "remote.node_name", "node2"),
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "description", "This connection is powered by Cisco Nexus Hyperfabric"),
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "cable_type", "DAC"),
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "cable_length", "7"),
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "pluggable", "QDD-400-AOC7M"),
				),
			},
			// Update with minimum config and verify config is unchanged.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Connection - Update with minimum config and verify config is unchanged.")
				},
				Config:             testConnectionResourceHclConfig(fabricName, "minimal"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "local.port_name", "Ethernet1_1"),
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "local.node_name", "node1"),
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "remote.port_name", "Ethernet1_1"),
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "remote.node_name", "node2"),
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "description", "This connection is powered by Cisco Nexus Hyperfabric"),
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "cable_type", "DAC"),
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "cable_length", "7"),
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "pluggable", "QDD-400-AOC7M"),
				),
			},
			// ImportState testing with pre-existing Id.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Connection - ImportState testing with pre-existing Id.")
				},
				ResourceName:      "hyperfabric_connection.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// ImportState testing with name.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Connection - ImportState testing with name.")
				},
				ResourceName:      "hyperfabric_connection.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: generateConnectionIdFromState,
			},
			// Update with config containing all optional attributes with empty values and verify config is cleared.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Connection - Update with config containing all optional attributes with empty values and verify config is cleared.")
				},
				Config:             testConnectionResourceHclConfig(fabricName, "clear"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "local.port_name", "Ethernet1_1"),
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "local.node_name", "node1"),
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "remote.port_name", "Ethernet1_1"),
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "remote.node_name", "node2"),
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "description", ""),
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "cable_type", "DAC"),
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "cable_length", "0"),
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "pluggable", ""),
				),
			},
			// Run Plan Only with minimal config and check that plan is empty.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Connection - Run Plan Only with minimal config and check that plan is empty.")
				},
				Config:             testConnectionResourceHclConfig(fabricName, "minimal"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "local.port_name", "Ethernet1_1"),
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "local.node_name", "node1"),
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "remote.port_name", "Ethernet1_1"),
					resource.TestCheckResourceAttr("hyperfabric_connection.test", "remote.node_name", "node2"),
				),
			},
		},
	})
}

func testConnectionResourceHclConfig(fabricName string, configType string) string {
	if configType == "full" {
		return fmt.Sprintf(`
resource "hyperfabric_fabric" "test" {
    name = "%[1]s"
}
resource "hyperfabric_node" "node1" {
    fabric_id = hyperfabric_fabric.test.id
    name = "node1"
    model_name = "HF6100-32D"
    roles = ["LEAF"]
}
resource "hyperfabric_node" "node2" {
    fabric_id = hyperfabric_fabric.test.id
    name = "node2"
    model_name = "HF6100-32D"
    roles = ["LEAF"]
}
resource "hyperfabric_connection" "test" {
    fabric_id = hyperfabric_fabric.test.id
    local = {
        node_id = hyperfabric_node.node1.node_id
        port_name = "Ethernet1_1"
    }
    remote = {
        node_id = hyperfabric_node.node2.node_id
        port_name = "Ethernet1_1"
    }
    description = "This connection is powered by Cisco Nexus Hyperfabric"
    cable_type = "DAC"
    cable_length = 7
    pluggable = "QDD-400-AOC7M"
}
`, fabricName)
	} else if configType == "clear" {
		return fmt.Sprintf(`
resource "hyperfabric_fabric" "test" {
    name = "%[1]s"
}
resource "hyperfabric_node" "node1" {
    fabric_id = hyperfabric_fabric.test.id
    name = "node1"
    model_name = "HF6100-32D"
    roles = ["LEAF"]
}
resource "hyperfabric_node" "node2" {
    fabric_id = hyperfabric_fabric.test.id
    name = "node2"
    model_name = "HF6100-32D"
    roles = ["LEAF"]
}
resource "hyperfabric_connection" "test" {
    fabric_id = hyperfabric_fabric.test.id
    local = {
        node_id = hyperfabric_node.node1.node_id
        port_name = "Ethernet1_1"
    }
    remote = {
        node_id = hyperfabric_node.node2.node_id
        port_name = "Ethernet1_1"
    }
    description = ""
    cable_type = "DAC"
    cable_length = 0
    pluggable = ""
}
`, fabricName)
	} else {
		return fmt.Sprintf(`
resource "hyperfabric_fabric" "test" {
    name = "%[1]s"
}
resource "hyperfabric_node" "node1" {
    fabric_id = hyperfabric_fabric.test.id
    name = "node1"
    model_name = "HF6100-32D"
    roles = ["LEAF"]
}
resource "hyperfabric_node" "node2" {
    fabric_id = hyperfabric_fabric.test.id
    name = "node2"
    model_name = "HF6100-32D"
    roles = ["LEAF"]
}
resource "hyperfabric_connection" "test" {
    fabric_id = hyperfabric_fabric.test.id
    local = {
        node_id = hyperfabric_node.node1.node_id
        port_name = "Ethernet1_1"
    }
    remote = {
        node_id = hyperfabric_node.node2.node_id
        port_name = "Ethernet1_1"
    }
}
`, fabricName)
	}
}

func generateConnectionIdFromState(state *terraform.State) (string, error) {
	var fabricName, connectionId string
	for _, m := range state.Modules {
		if len(m.Resources) > 0 {
			if v, ok := m.Resources["hyperfabric_connection.test"]; ok {
				connectionId = v.Primary.Attributes["connection_id"]
			}
			if w, ok := m.Resources["hyperfabric_fabric.test"]; ok {
				fabricName = w.Primary.Attributes["name"]
			}
		}
	}
	importId := fmt.Sprintf("%s/connections/%s", fabricName, connectionId)
	return importId, nil
}
