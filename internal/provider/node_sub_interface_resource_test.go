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

func TestAccNodeSubInterfaceResource(t *testing.T) {
	fabricName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify provided and default Hyperfabric values.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Sub-Interface - Create with minimum config and verify provided and default Hyperfabric values.")
				},
				Config:             testNodeSubInterfaceResourceHclConfig(fabricName, "", "minimal"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "name", "Ethernet1_1.100"),
				),
			},
			// Update with all config and verify provided values.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Sub-Interface - Update with all config and verify provided values.")
				},
				Config:             testNodeSubInterfaceResourceHclConfig(fabricName, "test", "full"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "name", "Ethernet1_1.100"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "description", "Loopback for BGP peering"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "enabled", "true"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "ipv4_addresses.#", "1"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "ipv4_addresses.0", "10.1.0.1/24"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "ipv6_addresses.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "ipv6_addresses.0", "2001:1::1/64"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "ipv6_addresses.1", "2002:1::1/64"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "parent", "Ethernet1_1"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "labels.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "annotations.#", "2"),
					resource.TestCheckResourceAttrSet("hyperfabric_node_sub_interface.test", "vlan_id"),
					resource.TestCheckResourceAttrSet("hyperfabric_node_sub_interface.test", "vrf_id"),
				),
			},
			// Update with minimum config and verify config is unchanged.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Sub-Interface - Update with minimum config and verify config is unchanged.")
				},
				Config:             testNodeSubInterfaceResourceHclConfig(fabricName, "test", "minimal"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "name", "Ethernet1_1.100"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "description", "Loopback for BGP peering"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "enabled", "true"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "ipv4_addresses.#", "1"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "ipv4_addresses.0", "10.1.0.1/24"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "ipv6_addresses.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "ipv6_addresses.0", "2001:1::1/64"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "ipv6_addresses.1", "2002:1::1/64"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "parent", "Ethernet1_1"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "labels.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "annotations.#", "2"),
					resource.TestCheckResourceAttrSet("hyperfabric_node_sub_interface.test", "vlan_id"),
					resource.TestCheckResourceAttrSet("hyperfabric_node_sub_interface.test", "vrf_id"),
				),
			},
			// ImportState testing with pre-existing Id.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Sub-Interface - ImportState testing with pre-existing Id.")
				},
				ResourceName:      "hyperfabric_node_sub_interface.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// ImportState testing with fabric and node name.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Sub-Interface - ImportState testing with fabric, node and interface name.")
				},
				ResourceName:      "hyperfabric_node_sub_interface.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fabricName + "/nodes/node1/subInterfaces/Ethernet1_1.100",
			},
			// Update with all config and verify provided values but reset vrf_id
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Sub-Interface - Update with all config and verify provided values but change vrf_id")
				},
				Config:             testNodeSubInterfaceResourceHclConfig(fabricName, "test2", "full"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "name", "Ethernet1_1.100"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "description", "Loopback for BGP peering"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "enabled", "true"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "ipv4_addresses.#", "1"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "ipv4_addresses.0", "10.1.0.1/24"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "ipv6_addresses.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "ipv6_addresses.0", "2001:1::1/64"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "ipv6_addresses.1", "2002:1::1/64"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "parent", "Ethernet1_1"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "labels.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "annotations.#", "2"),
					resource.TestCheckResourceAttrSet("hyperfabric_node_sub_interface.test", "vlan_id"),
					resource.TestCheckResourceAttrSet("hyperfabric_node_sub_interface.test", "vrf_id"),
				),
			},
			// Update with config containing all optional attributes with empty values and verify config is cleared.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Sub-Interface - Update with config containing all optional attributes with empty values and verify config is cleared.")
				},
				Config:             testNodeSubInterfaceResourceHclConfig(fabricName, "test2", "clear"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "name", "Ethernet1_1.100"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "description", ""),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "enabled", "false"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "ipv4_addresses.#", "0"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "ipv6_addresses.#", "0"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "parent", "Ethernet1_1"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "labels.#", "0"),
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "annotations.#", "0"),
					resource.TestCheckResourceAttrSet("hyperfabric_node_sub_interface.test", "vlan_id"),
					resource.TestCheckResourceAttrSet("hyperfabric_node_sub_interface.test", "vrf_id"),
				),
			},
			// Run Plan Only with minimal config and check that plan is empty.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Sub-Interface - Run Plan Only with minimal config and check that plan is empty.")
				},
				Config:             testNodeSubInterfaceResourceHclConfig(fabricName, "test2", "minimal"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node_sub_interface.test", "name", "Ethernet1_1.100"),
				),
			},
		},
	})
}

func testNodeSubInterfaceResourceHclConfig(fabricName string, vrf string, configType string) string {
	vrfConfigLine := ""
	if vrf != "" {
		vrfConfigLine = fmt.Sprintf("vrf_id = hyperfabric_vrf.%s.vrf_id", vrf)
	} else if vrf == "clear" {
		vrfConfigLine = "vrf_id = \"\""
	}

	if configType == "full" {
		return fmt.Sprintf(`
resource "hyperfabric_fabric" "test" {
	name = "%[1]s"
}

resource "hyperfabric_vrf" "test" {
    fabric_id = hyperfabric_fabric.test.id
	name      = "Vrf1"
}

resource "hyperfabric_vrf" "test2" {
    fabric_id = hyperfabric_fabric.test.id
	name      = "Vrf2"
}

resource "hyperfabric_node" "test" {
    fabric_id   = hyperfabric_fabric.test.id
	name        = "node1"
	model_name  = "HF6100-32D"
}

resource "hyperfabric_node_sub_interface" "test" {
	node_id        = hyperfabric_node.test.id
	name           = "Ethernet1_1.100"
	description    = "Loopback for BGP peering"
	enabled        = true
	ipv4_addresses = ["10.1.0.1/24"]
	ipv6_addresses = ["2001:1::1/64", "2002:1::1/64"]
	vlan_id        = 100
	%[2]s
	labels         = [
		"sj01-1-101-AAA01",
		"blue"
	]
	annotations    = [
		{
			name      = "color"
			value     = "blue"
		},
		{
			data_type = "UINT32"
			name      = "rack"
			value     = "1"
		}
	]
}
`, fabricName, vrfConfigLine)
	} else if configType == "clear" {
		return fmt.Sprintf(`
resource "hyperfabric_fabric" "test" {
	name = "%[1]s"
}

resource "hyperfabric_vrf" "test" {
    fabric_id = hyperfabric_fabric.test.id
	name      = "Vrf1"
}

resource "hyperfabric_vrf" "test2" {
    fabric_id = hyperfabric_fabric.test.id
	name      = "Vrf2"
}

resource "hyperfabric_node" "test" {
	fabric_id  = hyperfabric_fabric.test.id
	name       = "node1"
	model_name = "HF6100-32D"
}

resource "hyperfabric_node_sub_interface" "test" {
	node_id        = hyperfabric_node.test.id
	name           = "Ethernet1_1.100"
	description    = ""
	enabled        = false
	ipv4_addresses = []
	ipv6_addresses = []
	labels         = []
	annotations    = []
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

resource "hyperfabric_vrf" "test2" {
    fabric_id = hyperfabric_fabric.test.id
	name      = "Vrf2"
}

resource "hyperfabric_node" "test" {
	fabric_id  = hyperfabric_fabric.test.id
	name       = "node1"
	model_name = "HF6100-32D"
}

resource "hyperfabric_node_sub_interface" "test" {
	node_id = hyperfabric_node.test.id
	name    = "Ethernet1_1.100"
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

resource "hyperfabric_vrf" "test2" {
    fabric_id = hyperfabric_fabric.test.id
	name      = "Vrf2"
}

resource "hyperfabric_node" "test" {
    fabric_id  = hyperfabric_fabric.test.id
	name       = "node1"
	model_name = "HF6100-32D"
}

resource "hyperfabric_node_sub_interface" "test" {
    node_id = hyperfabric_node.test.id
	name    = "Ethernet1_1.100"
}
`, fabricName)
	}
}
