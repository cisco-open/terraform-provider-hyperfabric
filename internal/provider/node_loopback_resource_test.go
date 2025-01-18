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

func TestAccNodeLoopbackResource(t *testing.T) {
	fabricName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	loopbackName := "Loopback10"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify provided and default Hyperfabric values.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Loopback - Create with minimum config and verify provided and default Hyperfabric values.")
				},
				Config:             testNodeLoopbackResourceHclConfig(fabricName, loopbackName, "", "minimal"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node_loopback.test", "name", loopbackName),
				),
			},
			// Update with all config and verify provided values.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Loopback - Update with all config and verify provided values.")
				},
				Config:             testNodeLoopbackResourceHclConfig(fabricName, loopbackName, "test", "full"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node_loopback.test", "name", loopbackName),
					resource.TestCheckResourceAttr("hyperfabric_node_loopback.test", "description", "Loopback for BGP Peering"),
					resource.TestCheckResourceAttr("hyperfabric_node_loopback.test", "ipv4_address", "10.1.0.1"),
					resource.TestCheckResourceAttr("hyperfabric_node_loopback.test", "ipv6_address", "2001:1::1"),
					resource.TestCheckResourceAttrSet("hyperfabric_node_loopback.test", "vrf_id"),
					resource.TestCheckResourceAttr("hyperfabric_node_loopback.test", "labels.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_node_loopback.test", "annotations.#", "2"),
				),
			},
			// Update with minimum config and verify config is unchanged.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Loopback - Update with minimum config and verify config is unchanged.")
				},
				Config:             testNodeLoopbackResourceHclConfig(fabricName, loopbackName, "test", "minimal"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node_loopback.test", "name", loopbackName),
					resource.TestCheckResourceAttr("hyperfabric_node_loopback.test", "description", "Loopback for BGP Peering"),
					resource.TestCheckResourceAttr("hyperfabric_node_loopback.test", "ipv4_address", "10.1.0.1"),
					resource.TestCheckResourceAttr("hyperfabric_node_loopback.test", "ipv6_address", "2001:1::1"),
					resource.TestCheckResourceAttrSet("hyperfabric_node_loopback.test", "vrf_id"),
					resource.TestCheckResourceAttr("hyperfabric_node_loopback.test", "labels.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_node_loopback.test", "annotations.#", "2"),
				),
			},
			// ImportState testing with pre-existing Id.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Loopback - ImportState testing with pre-existing Id.")
				},
				ResourceName:      "hyperfabric_node_loopback.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// ImportState testing with fabric and node name.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Loopback - ImportState testing with fabric, node and interface name.")
				},
				ResourceName:      "hyperfabric_node_loopback.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fabricName + "/nodes/node1/loopbacks/" + loopbackName,
			},
			// Update with all config and verify provided values but reset vrf_id
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Loopback - Update with all config and verify provided values but change vrf_id")
				},
				Config:             testNodeLoopbackResourceHclConfig(fabricName, loopbackName, "test2", "full"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node_loopback.test", "name", loopbackName),
					resource.TestCheckResourceAttr("hyperfabric_node_loopback.test", "description", "Loopback for BGP Peering"),
					resource.TestCheckResourceAttr("hyperfabric_node_loopback.test", "ipv4_address", "10.1.0.1"),
					resource.TestCheckResourceAttr("hyperfabric_node_loopback.test", "ipv6_address", "2001:1::1"),
					resource.TestCheckResourceAttrSet("hyperfabric_node_loopback.test", "vrf_id"),
					resource.TestCheckResourceAttr("hyperfabric_node_loopback.test", "labels.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_node_loopback.test", "annotations.#", "2"),
				),
			},
			// Update with config containing all optional attributes with empty values and verify config is cleared.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Loopback - Update with config containing all optional attributes with empty values and verify config is cleared.")
				},
				Config:             testNodeLoopbackResourceHclConfig(fabricName, loopbackName, "test2", "clear"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node_loopback.test", "name", loopbackName),
					resource.TestCheckResourceAttr("hyperfabric_node_loopback.test", "description", ""),
					resource.TestCheckResourceAttr("hyperfabric_node_loopback.test", "ipv4_address", "10.1.0.1"),
					resource.TestCheckResourceAttr("hyperfabric_node_loopback.test", "ipv6_address", ""),
					resource.TestCheckResourceAttr("hyperfabric_node_loopback.test", "labels.#", "0"),
					resource.TestCheckResourceAttr("hyperfabric_node_loopback.test", "annotations.#", "0"),
				),
			},
			// Run Plan Only with minimal config and check that plan is empty.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Loopback - Run Plan Only with minimal config and check that plan is empty.")
				},
				Config:             testNodeLoopbackResourceHclConfig(fabricName, loopbackName, "", "minimal"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node_loopback.test", "name", loopbackName),
				),
			},
		},
	})
}

func testNodeLoopbackResourceHclConfig(fabricName string, loopbackName string, vrf string, configType string) string {
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

resource "hyperfabric_node_loopback" "test" {
	node_id      = hyperfabric_node.test.id
	name         = "%[2]s"
	description  = "Loopback for BGP Peering"
	ipv4_address = "10.1.0.1"
	ipv6_address = "2001:1::1"
	%[3]s
	labels       = [
		"sj01-1-101-AAA01",
		"blue"
	]
	annotations  = [
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
`, fabricName, loopbackName, vrfConfigLine)
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
	fabric_id   = hyperfabric_fabric.test.id
	name        = "node1"
	model_name  = "HF6100-32D"
}

resource "hyperfabric_node_loopback" "test" {
	node_id      = hyperfabric_node.test.id
	name         = "%[2]s"
	description  = ""
	ipv4_address = "10.1.0.1"
	ipv6_address = ""
	%[3]s
	labels       = []
	annotations  = []
}
`, fabricName, loopbackName, vrfConfigLine)
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

resource "hyperfabric_node_loopback" "test" {
	node_id = hyperfabric_node.test.id
	name    = "%[2]s"
	ipv4_address = "10.1.0.1"
}
`, fabricName, loopbackName)
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

resource "hyperfabric_node_loopback" "test" {
    node_id = hyperfabric_node.test.id
	name    = "%[2]s"
	ipv4_address = "10.1.0.1"
}
`, fabricName, loopbackName)
	}
}
