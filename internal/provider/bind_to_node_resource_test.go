// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBindToNodeResource(t *testing.T) {
	fabricName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	deviceId := getStringAttribute(basetypes.NewStringNull(), "TF_ACC_HYPERFABRIC_DEVICE_ID", "")
	if deviceId == "" {
		t.Fatalf("ERROR: Missing deviceId for test. Please configure environment variable TF_ACC_HYPERFABRIC_DEVICE_ID with a valid deviceId.")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify provided and default Hyperfabric values.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Port - Create with minimum config and verify provided and default Hyperfabric values.")
				},
				Config:             testBindToNodeResourceHclConfig(fabricName, deviceId, "minimal"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hyperfabric_bind_to_node.test", "node_id"),
					resource.TestCheckResourceAttrSet("hyperfabric_bind_to_node.test", "device_id"),
				),
			},
			// Update with all config and verify provided values.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Port - Update with all config and verify provided values.")
				},
				Config:             testBindToNodeResourceHclConfig(fabricName, deviceId, "minimal"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hyperfabric_bind_to_node.test", "node_id"),
					resource.TestCheckResourceAttrSet("hyperfabric_bind_to_node.test", "device_id"),
				),
			},
			// Update with minimum config and verify config is unchanged.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Port - Update with minimum config and verify config is unchanged.")
				},
				Config:             testBindToNodeResourceHclConfig(fabricName, deviceId, "minimal"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hyperfabric_bind_to_node.test", "node_id"),
					resource.TestCheckResourceAttrSet("hyperfabric_bind_to_node.test", "device_id"),
				),
			},
			// ImportState testing with pre-existing Id.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Port - ImportState testing with pre-existing Id.")
				},
				ResourceName:      "hyperfabric_bind_to_node.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// ImportState testing with fabric and node name.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Port - ImportState testing with fabric, node and interface name.")
				},
				ResourceName:      "hyperfabric_bind_to_node.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fabricName + "/nodes/node1",
			},
			// Update with config containing all optional attributes with empty values and verify config is cleared.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Port - Update with config containing all optional attributes with empty values and verify config is cleared.")
				},
				Config:             testBindToNodeResourceHclConfig(fabricName, deviceId, "minimal"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hyperfabric_bind_to_node.test", "node_id"),
					resource.TestCheckResourceAttrSet("hyperfabric_bind_to_node.test", "device_id"),
				),
			},
			// Run Plan Only with minimal config and check that plan is empty.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Port - Run Plan Only with minimal config and check that plan is empty.")
				},
				Config:             testBindToNodeResourceHclConfig(fabricName, deviceId, "minimal"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hyperfabric_bind_to_node.test", "node_id"),
					resource.TestCheckResourceAttrSet("hyperfabric_bind_to_node.test", "device_id"),
				),
			},
		},
	})
}

func testBindToNodeResourceHclConfig(fabricName string, deviceId string, configType string) string {
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

resource "hyperfabric_bind_to_node" "test" {
    node_id = hyperfabric_node.test.id
	device_id = "%[2]s"
}
`, fabricName, deviceId)
}
