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

func TestAccVrfResource(t *testing.T) {
	name := "Vrf" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	fabricName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify provided and default Hyperfabric values.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: VRF - Create with minimum config and verify provided and default Hyperfabric values.")
				},
				Config:             testVrfResourceHclConfig(fabricName, name, "minimal"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_vrf.test", "name", name),
				),
			},
			// Update with all config and verify provided values.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: VRF - Update with all config and verify provided values.")
				},
				Config:             testVrfResourceHclConfig(fabricName, name, "full"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_vrf.test", "name", name),
					resource.TestCheckResourceAttr("hyperfabric_vrf.test", "description", "This VRF is powered by Cisco Nexus Hyperfabric"),
					// resource.TestCheckResourceAttr("hyperfabric_vrf.test", "asn", "65002"),
					// resource.TestCheckResourceAttr("hyperfabric_vrf.test", "vni", "169"),
					resource.TestCheckResourceAttr("hyperfabric_vrf.test", "labels.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_vrf.test", "annotations.#", "2"),
				),
			},
			// Update with minimum config and verify config is unchanged.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: VRF - Update with minimum config and verify config is unchanged.")
				},
				Config:             testVrfResourceHclConfig(fabricName, name, "minimal"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_vrf.test", "name", name),
					resource.TestCheckResourceAttr("hyperfabric_vrf.test", "description", "This VRF is powered by Cisco Nexus Hyperfabric"),
					// resource.TestCheckResourceAttr("hyperfabric_vrf.test", "asn", "65002"),
					// resource.TestCheckResourceAttr("hyperfabric_vrf.test", "vni", "169"),
					resource.TestCheckResourceAttr("hyperfabric_vrf.test", "labels.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_vrf.test", "annotations.#", "2"),
				),
			},
			// ImportState testing with pre-existing Id.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: VRF - ImportState testing with pre-existing Id.")
				},
				ResourceName:      "hyperfabric_vrf.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// ImportState testing with name.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: VRF - ImportState testing with name.")
				},
				ResourceName:      "hyperfabric_vrf.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fabricName + "/vrfs/" + name,
			},
			// Update with config containing all optional attributes with empty values and verify config is cleared.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: VRF - Update with config containing all optional attributes with empty values and verify config is cleared.")
				},
				Config:             testVrfResourceHclConfig(fabricName, name, "clear"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_vrf.test", "name", name),
					// resource.TestCheckResourceAttr("hyperfabric_vrf.test", "asn", "65002"),
					// resource.TestCheckResourceAttr("hyperfabric_vrf.test", "vni", "169"),
					resource.TestCheckResourceAttr("hyperfabric_vrf.test", "description", ""),
					resource.TestCheckResourceAttr("hyperfabric_vrf.test", "labels.#", "0"),
					resource.TestCheckResourceAttr("hyperfabric_vrf.test", "annotations.#", "0"),
				),
			},
			// Run Plan Only with minimal config and check that plan is empty.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: VRF - Run Plan Only with minimal config and check that plan is empty.")
				},
				Config:             testVrfResourceHclConfig(fabricName, name, "minimal"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_vrf.test", "name", name),
				),
			},
		},
	})
}

func testVrfResourceHclConfig(fabricName string, name string, configType string) string {
	if configType == "full" {
		return fmt.Sprintf(`
resource "hyperfabric_fabric" "test" {
	name = "%[1]s"
}

resource "hyperfabric_vrf" "test" {
    fabric_id   = hyperfabric_fabric.test.id
	name        = "%[2]s"
	description = "This VRF is powered by Cisco Nexus Hyperfabric"
	// asn         = 65002
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

resource "hyperfabric_vrf" "test" {
	fabric_id   = hyperfabric_fabric.test.id
	name        = "%[2]s"
	// asn         = 65002
	description = ""
	labels = []
	annotations = []
}
`, fabricName, name)
	} else {
		return fmt.Sprintf(`
resource "hyperfabric_fabric" "test" {
	name = "%[1]s"
}

resource "hyperfabric_vrf" "test" {
    fabric_id = hyperfabric_fabric.test.id
	name = "%[2]s"
}
`, fabricName, name)
	}
}
