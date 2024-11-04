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

func TestAccNodeManagementPortResource(t *testing.T) {
	fabricName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify provided and default Hyperfabric values.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Management Port - Create with minimum config and verify provided and default Hyperfabric values.")
				},
				Config:             testNodeManagementPortResourceHclConfig(fabricName, "minimal"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "name", "eth0"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ipv4_config_type", "CONFIG_TYPE_DHCP"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ipv6_config_type", "CONFIG_TYPE_DHCP"),
				),
			},
			// Update with all config and verify provided values.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Management Port - Update with all config and verify provided values.")
				},
				Config:             testNodeManagementPortResourceHclConfig(fabricName, "full"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "name", "eth0"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ipv4_config_type", "CONFIG_TYPE_STATIC"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ipv4_address", "10.0.0.3/24"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ipv4_gateway", "10.0.0.254"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ipv6_config_type", "CONFIG_TYPE_STATIC"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ipv6_address", "2001::3/64"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ipv6_gateway", "2001::254"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "dns_addresses.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "dns_addresses.0", "1.1.1.1"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "dns_addresses.1", "8.8.8.8"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "cloud_urls.#", "1"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "cloud_urls.0", "https://hyperfabric.cisco.com"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ntp_addresses.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ntp_addresses.0", "be.pool.ntp.org"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ntp_addresses.1", "us.pool.ntp.org"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "no_proxy.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "no_proxy.0", "10.0.0.1"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "no_proxy.1", "server.local"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "proxy_address", "http://proxy.mycompany.com:80"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "proxy_username", "my_proxy_user"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "proxy_password", "my_super_secret_password"),
				),
			},
			// Update with minimum config and verify config is unchanged.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Management Port - Update with minimum config and verify config is unchanged.")
				},
				Config:             testNodeManagementPortResourceHclConfig(fabricName, "minimal+"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "name", "eth0"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ipv4_config_type", "CONFIG_TYPE_STATIC"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ipv4_address", "10.0.0.3/24"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ipv4_gateway", "10.0.0.254"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ipv6_config_type", "CONFIG_TYPE_STATIC"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ipv6_address", "2001::3/64"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ipv6_gateway", "2001::254"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "dns_addresses.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "dns_addresses.0", "8.8.8.8"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "dns_addresses.1", "1.1.1.1"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "cloud_urls.#", "1"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "cloud_urls.0", "https://hyperfabric.cisco.com"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ntp_addresses.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ntp_addresses.0", "be.pool.ntp.org"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ntp_addresses.1", "us.pool.ntp.org"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "no_proxy.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "no_proxy.0", "10.0.0.1"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "no_proxy.1", "server.local"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "proxy_address", "http://proxy.mycompany.com:80"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "proxy_username", "my_proxy_user"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "proxy_password", "my_super_secret_password"),
				),
			},
			// ImportState testing with pre-existing Id.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Management Port - ImportState testing with pre-existing Id.")
				},
				ResourceName:            "hyperfabric_node_management_port.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"proxy_password"},
			},
			// ImportState testing with fabric and node name.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Management Port - ImportState testing with fabric and node name.")
				},
				ResourceName:            "hyperfabric_node_management_port.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"proxy_password"},
				ImportStateId:           fabricName + "/nodes/node1",
			},
			// ImportState testing with fabric, node and interface name.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Management Port - ImportState testing with fabric, node and interface name.")
				},
				ResourceName:            "hyperfabric_node_management_port.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"proxy_password"},
				ImportStateId:           fabricName + "/nodes/node1/managementPorts/eth0",
			},
			// Update with config containing all optional attributes with empty values and verify config is cleared.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Management Port - Update with config containing all optional attributes with empty values and verify config is cleared.")
				},
				Config:             testNodeManagementPortResourceHclConfig(fabricName, "clear"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "name", "eth0"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ipv4_config_type", "CONFIG_TYPE_DHCP"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ipv4_address", ""),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ipv4_gateway", ""),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ipv6_config_type", "CONFIG_TYPE_DHCP"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ipv6_address", ""),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ipv6_gateway", ""),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "dns_addresses.#", "0"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "cloud_urls.#", "0"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ntp_addresses.#", "0"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "no_proxy.#", "0"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "proxy_address", ""),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "proxy_username", ""),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "proxy_password", ""),
				),
			},
			// Run Plan Only with minimal config and check that plan is empty.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Management Port - Run Plan Only with minimal config and check that plan is empty.")
				},
				Config:             testNodeManagementPortResourceHclConfig(fabricName, "minimal"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "name", "eth0"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ipv4_config_type", "CONFIG_TYPE_DHCP"),
					resource.TestCheckResourceAttr("hyperfabric_node_management_port.test", "ipv6_config_type", "CONFIG_TYPE_DHCP"),
				),
			},
		},
	})
}

func testNodeManagementPortResourceHclConfig(fabricName string, configType string) string {
	if configType == "full" {
		return fmt.Sprintf(`
resource "hyperfabric_fabric" "test" {
	name = "%[1]s"
}

resource "hyperfabric_node" "test" {
    fabric_id   = hyperfabric_fabric.test.id
	name        = "node1"
	model_name  = "HF6100-32D"
    roles       = ["LEAF"]
}

resource "hyperfabric_node_management_port" "test" {
	node_id          = hyperfabric_node.test.id
	name             = "eth0"
	ipv4_config_type = "CONFIG_TYPE_STATIC"
	ipv4_address     = "10.0.0.3/24"
	ipv4_gateway     = "10.0.0.254"
	ipv6_config_type = "CONFIG_TYPE_STATIC"
	ipv6_address     = "2001::3/64"
	ipv6_gateway     = "2001::254"
	dns_addresses    = ["8.8.8.8", "1.1.1.1"]
	cloud_urls       = ["https://hyperfabric.cisco.com"]
	ntp_addresses    = ["be.pool.ntp.org", "us.pool.ntp.org"]
	no_proxy         = ["10.0.0.1", "server.local"]
	proxy_address    = "http://proxy.mycompany.com:80"
	proxy_username   = "my_proxy_user"
	proxy_password   = "my_super_secret_password"
}
`, fabricName)
	} else if configType == "clear" {
		return fmt.Sprintf(`
resource "hyperfabric_fabric" "test" {
	name = "%[1]s"
}

resource "hyperfabric_node" "test" {
	fabric_id   = hyperfabric_fabric.test.id
	name        = "node1"
	model_name  = "HF6100-32D"
    roles       = ["LEAF"]
}

resource "hyperfabric_node_management_port" "test" {
	node_id          = hyperfabric_node.test.id
	name             = "eth0"
	ipv4_config_type = "CONFIG_TYPE_DHCP"
	ipv4_address     = ""
	ipv4_gateway     = ""
	ipv6_config_type = "CONFIG_TYPE_DHCP"
	ipv6_address     = ""
	ipv6_gateway     = ""
	dns_addresses    = []
	cloud_urls       = []
	ntp_addresses    = []
	no_proxy         = []
	proxy_address    = ""
	proxy_username   = ""
	proxy_password   = ""
}
`, fabricName)
	} else if configType == "minimal+" {
		return fmt.Sprintf(`
resource "hyperfabric_fabric" "test" {
	name = "%[1]s"
}

resource "hyperfabric_node" "test" {
	fabric_id  = hyperfabric_fabric.test.id
	name       = "node1"
	model_name = "HF6100-32D"
	roles      = ["LEAF"]
}

resource "hyperfabric_node_management_port" "test" {
	node_id = hyperfabric_node.test.id
	ipv4_config_type = "CONFIG_TYPE_STATIC"
	ipv6_config_type = "CONFIG_TYPE_STATIC"
}
`, fabricName)
	} else {
		return fmt.Sprintf(`
resource "hyperfabric_fabric" "test" {
	name = "%[1]s"
}

resource "hyperfabric_node" "test" {
    fabric_id  = hyperfabric_fabric.test.id
	name       = "node1"
	model_name = "HF6100-32D"
    roles      = ["LEAF"]
}

resource "hyperfabric_node_management_port" "test" {
    node_id = hyperfabric_node.test.id
}
`, fabricName)
	}
}
