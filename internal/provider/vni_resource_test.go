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

func TestAccVniResource(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	fabricName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify provided and default Hyperfabric values.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: VNI - Create with minimum config and verify provided and default Hyperfabric values.")
				},
				Config:             testVniResourceHclConfig(fabricName, name, 0, "true", "", "minimal"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "name", name),
					resource.TestCheckResourceAttrSet("hyperfabric_vni.test", "vni"),
				),
			},
			// Update with all config except `vni` and verify provided values.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: VNI - Update with all config except `vni` and verify provided values.")
				},
				Config:             testVniResourceHclConfig(fabricName, name, 0, "true", "test", "full"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "name", name),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "description", "This VNI is powered by Cisco Nexus Hyperfabric"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.enabled", "true"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.ipv4_addresses.#", "1"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.ipv4_addresses.0", "192.168.0.254/24"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.ipv6_addresses.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.ipv6_addresses.0", "2001::1/64"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.ipv6_addresses.1", "2002::1/64"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.#", "3"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.0.node_id", "*"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.0.port_name", "Ethernet1_10"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.0.vlan_id", "103"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.1.node_id", "*"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.1.port_name", "Ethernet1_11"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.1.vlan_id", "103"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.2.node_name", "node1"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.2.port_name", "Ethernet1_9"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.2.vlan_id", "103"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "labels.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "annotations.#", "2"),
					resource.TestCheckResourceAttrSet("hyperfabric_vni.test", "vrf_id"),
				),
			},
			// Update with all config and verify provided values.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: VNI - Update with all config and verify provided values.")
				},
				Config:             testVniResourceHclConfig(fabricName, name, 169, "true", "test", "full"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "name", name),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "description", "This VNI is powered by Cisco Nexus Hyperfabric"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "vni", "169"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.enabled", "true"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.ipv4_addresses.#", "1"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.ipv4_addresses.0", "192.168.0.254/24"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.ipv6_addresses.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.ipv6_addresses.0", "2001::1/64"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.ipv6_addresses.1", "2002::1/64"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.#", "3"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.0.node_id", "*"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.0.port_name", "Ethernet1_10"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.0.vlan_id", "103"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.1.node_id", "*"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.1.port_name", "Ethernet1_11"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.1.vlan_id", "103"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.2.node_name", "node1"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.2.port_name", "Ethernet1_9"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.2.vlan_id", "103"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "labels.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "annotations.#", "2"),
					resource.TestCheckResourceAttrSet("hyperfabric_vni.test", "vrf_id"),
				),
			},
			// Update with minimum config and verify config is unchanged.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: VNI - Update with minimum config and verify config is unchanged.")
				},
				Config:             testVniResourceHclConfig(fabricName, name, 0, "true", "", "minimal"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "name", name),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "description", "This VNI is powered by Cisco Nexus Hyperfabric"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "vni", "169"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.enabled", "true"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.ipv4_addresses.#", "1"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.ipv4_addresses.0", "192.168.0.254/24"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.ipv6_addresses.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.ipv6_addresses.0", "2001::1/64"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.ipv6_addresses.1", "2002::1/64"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.#", "3"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.0.node_id", "*"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.0.port_name", "Ethernet1_10"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.0.vlan_id", "103"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.1.node_id", "*"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.1.port_name", "Ethernet1_11"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.1.vlan_id", "103"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.2.node_name", "node1"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.2.port_name", "Ethernet1_9"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.2.vlan_id", "103"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "labels.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "annotations.#", "2"),
					resource.TestCheckResourceAttrSet("hyperfabric_vni.test", "vrf_id"),
				),
			},
			// ImportState testing with pre-existing Id.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: VNI - ImportState testing with pre-existing Id.")
				},
				ResourceName:      "hyperfabric_vni.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// ImportState testing with name.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: VNI - ImportState testing with name.")
				},
				ResourceName:      "hyperfabric_vni.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fabricName + "/vnis/" + name,
			},
			// Update with all config and verify provided values but reset vrf_id
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: VNI - Update with all config and verify provided values but change vrf_id")
				},
				Config:             testVniResourceHclConfig(fabricName, name, 169, "false", "test2", "full"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "name", name),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "description", "This VNI is powered by Cisco Nexus Hyperfabric"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "vni", "169"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.enabled", "false"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.ipv4_addresses.#", "1"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.ipv4_addresses.0", "192.168.0.254/24"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.ipv6_addresses.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.ipv6_addresses.0", "2001::1/64"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.ipv6_addresses.1", "2002::1/64"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.#", "3"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.0.node_id", "*"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.0.port_name", "Ethernet1_10"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.0.vlan_id", "103"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.1.node_id", "*"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.1.port_name", "Ethernet1_11"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.1.vlan_id", "103"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.2.node_name", "node1"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.2.port_name", "Ethernet1_9"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.2.vlan_id", "103"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "labels.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "annotations.#", "2"),
					resource.TestCheckResourceAttrSet("hyperfabric_vni.test", "vrf_id"),
				),
			},
			// Update with all config and verify provided values but svi.enabled = false
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: VNI - Update with all config and verify provided values but svi.enabled set to false")
				},
				Config:             testVniResourceHclConfig(fabricName, name, 169, "false", "test2", "full"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "name", name),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "description", "This VNI is powered by Cisco Nexus Hyperfabric"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "vni", "169"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.enabled", "false"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.ipv4_addresses.#", "1"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.ipv4_addresses.0", "192.168.0.254/24"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.ipv6_addresses.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.ipv6_addresses.0", "2001::1/64"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "svi.ipv6_addresses.1", "2002::1/64"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.#", "3"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.0.node_id", "*"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.0.port_name", "Ethernet1_10"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.0.vlan_id", "103"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.1.node_id", "*"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.1.port_name", "Ethernet1_11"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.1.vlan_id", "103"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.2.node_name", "node1"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.2.port_name", "Ethernet1_9"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.2.vlan_id", "103"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "labels.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "annotations.#", "2"),
					resource.TestCheckResourceAttrSet("hyperfabric_vni.test", "vrf_id"),
				),
			},
			// Update with config containing all optional attributes with empty values and verify config is cleared.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: VNI - Update with config containing all optional attributes with empty values and verify config is cleared.")
				},
				Config:             testVniResourceHclConfig(fabricName, name, 169, "true", "clear", "clear"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "name", name),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "description", ""),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "vni", "169"),
					resource.TestCheckNoResourceAttr("hyperfabric_vni.test", "svi.enabled"),
					resource.TestCheckNoResourceAttr("hyperfabric_vni.test", "svi.ipv4_addresses"),
					resource.TestCheckNoResourceAttr("hyperfabric_vni.test", "svi.ipv6_addresses"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "members.#", "0"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "labels.#", "0"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "annotations.#", "0"),
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "vrf_id", ""),
				),
			},
			// Run Plan Only with minimal config and check that plan is empty.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: VNI - Run Plan Only with minimal config and check that plan is empty.")
				},
				Config:             testVniResourceHclConfig(fabricName, name, 0, "true", "", "minimal"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_vni.test", "name", name),
				),
			},
		},
	})
}

func testVniResourceHclConfig(fabricName string, name string, vni int64, sviEnabled string, vrf string, configType string) string {
	vniConfigLine := ""
	if vni != 0 {
		vniConfigLine = fmt.Sprintf("vni = %d", vni)
	}

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

resource "hyperfabric_node" "test" {
    fabric_id = hyperfabric_fabric.test.id
	name = "node1"
	model_name = "HF6100-32D"
    roles = ["LEAF"]
}

resource "hyperfabric_vrf" "test" {
    fabric_id = hyperfabric_fabric.test.id
	name = "Vrf1"
}

resource "hyperfabric_vrf" "test2" {
    fabric_id = hyperfabric_fabric.test.id
	name = "Vrf2"
}

resource "hyperfabric_vni" "test" {
    fabric_id   = hyperfabric_fabric.test.id
	name        = "%[2]s"
	description = "This VNI is powered by Cisco Nexus Hyperfabric"
	%[3]s
	svi = {
		enabled        = %[4]s
		ipv4_addresses = ["192.168.0.254/24"]
		ipv6_addresses = ["2001::1/64", "2002::1/64"]
	}
	members = [
		{
		node_id = "*"
		port_name = "Ethernet1_11"
		vlan_id = 103
		},
		{
		node_id = "*"
		port_name = "Ethernet1_10"
		vlan_id = 103
		},
		{
		node_id   = hyperfabric_node.test.node_id
		port_name = "Ethernet1_9"
		vlan_id   = 103
		}
	]
	%[5]s
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
`, fabricName, name, vniConfigLine, sviEnabled, vrfConfigLine)
	} else if configType == "clear" {
		return fmt.Sprintf(`
resource "hyperfabric_fabric" "test" {
	name = "%[1]s"
}

resource "hyperfabric_node" "test" {
    fabric_id = hyperfabric_fabric.test.id
	name = "node1"
	model_name = "HF6100-32D"
    roles = ["LEAF"]
}

resource "hyperfabric_vrf" "test" {
    fabric_id = hyperfabric_fabric.test.id
	name = "Vrf1"
}

resource "hyperfabric_vrf" "test2" {
    fabric_id = hyperfabric_fabric.test.id
	name = "Vrf2"
}

resource "hyperfabric_vni" "test" {
	fabric_id   = hyperfabric_fabric.test.id
	name        = "%[2]s"
	description = ""
	%[3]s
	svi = {}
	members = []
	vrf_id = ""
	labels = []
	annotations = []
}
`, fabricName, name, vniConfigLine)
	} else {
		return fmt.Sprintf(`
resource "hyperfabric_fabric" "test" {
	name = "%[1]s"
}

resource "hyperfabric_node" "test" {
    fabric_id = hyperfabric_fabric.test.id
	name = "node1"
	model_name = "HF6100-32D"
    roles = ["LEAF"]
}

resource "hyperfabric_vrf" "test" {
    fabric_id = hyperfabric_fabric.test.id
	name = "Vrf1"
}

resource "hyperfabric_vrf" "test2" {
    fabric_id = hyperfabric_fabric.test.id
	name = "Vrf2"
}

resource "hyperfabric_vni" "test" {
    fabric_id = hyperfabric_fabric.test.id
	name = "%[2]s"
	%[3]s
}
`, fabricName, name, vniConfigLine)
	}
}
