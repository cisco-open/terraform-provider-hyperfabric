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

func TestAccBearerTokenResource(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify provided and default Hyperfabric values.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: BearerToken - Create with minimum config and verify provided and default Hyperfabric values.")
				},
				Config:             testBearerTokenResourceHclConfig(name, "minimal"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_bearer_token.test", "name", name),
				),
			},
			// Update with all config and verify provided values.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: BearerToken - Update with all config and verify provided values.")
				},
				Config:             testBearerTokenResourceHclConfig(name, "full"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_bearer_token.test", "name", name),
					resource.TestCheckResourceAttr("hyperfabric_bearer_token.test", "description", "This bearer token is powered by Cisco Nexus Hyperfabric"),
					resource.TestCheckResourceAttr("hyperfabric_bearer_token.test", "not_after", "2025-09-03T08:00:00.000Z"),
					resource.TestCheckResourceAttr("hyperfabric_bearer_token.test", "not_before", "2024-09-03T08:00:00.000Z"),
					resource.TestCheckResourceAttr("hyperfabric_bearer_token.test", "scope", "ADMIN"),
				),
			},
			// Update with minimum config and verify config is unchanged.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: BearerToken - Update with minimum config and verify config is unchanged.")
				},
				Config:             testBearerTokenResourceHclConfig(name, "minimal+"),
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_bearer_token.test", "name", name),
					resource.TestCheckResourceAttr("hyperfabric_bearer_token.test", "description", "This bearer token is powered by Cisco Nexus Hyperfabric"),
					resource.TestCheckResourceAttr("hyperfabric_bearer_token.test", "not_after", "2025-09-03T08:00:00.000Z"),
					resource.TestCheckResourceAttr("hyperfabric_bearer_token.test", "not_before", "2024-09-03T08:00:00.000Z"),
					resource.TestCheckResourceAttr("hyperfabric_bearer_token.test", "scope", "ADMIN"),
				),
			},
			// ImportState testing with pre-existing Id.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: BearerToken - ImportState testing with pre-existing Id.")
				},
				ResourceName:      "hyperfabric_bearer_token.test",
				ImportState:       true,
				ImportStateVerify: true,
				// TODO implement ImportStateCheck ImportStateCheckFunc for not_before, not_after
				ImportStateVerifyIgnore: []string{"token", "not_before", "not_after"},
			},
			// ImportState testing with name.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: BearerToken - ImportState testing with name.")
				},
				ResourceName:      "hyperfabric_bearer_token.test",
				ImportState:       true,
				ImportStateVerify: true,
				// TODO implement ImportStateCheck ImportStateCheckFunc for not_before, not_after
				ImportStateVerifyIgnore: []string{"token", "not_before", "not_after"},
				ImportStateId:           name,
			},
			// Update with config containing all optional attributes with empty values and verify config is cleared.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: BearerToken - Update with config containing all optional attributes with empty values and verify config is cleared.")
				},
				Config:             testBearerTokenResourceHclConfig(name, "clear"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_bearer_token.test", "name", name),
					resource.TestCheckResourceAttr("hyperfabric_bearer_token.test", "description", ""),
					resource.TestCheckResourceAttr("hyperfabric_bearer_token.test", "not_after", "2025-09-03T08:00:00.000Z"),
					resource.TestCheckResourceAttr("hyperfabric_bearer_token.test", "not_before", "2024-09-03T08:00:00.000Z"),
					resource.TestCheckResourceAttr("hyperfabric_bearer_token.test", "scope", "ADMIN"),
				),
			},
			// Run Plan Only with minimal config and check that plan is empty.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: BearerToken - Run Plan Only with minimal config and check that plan is empty.")
				},
				Config:             testBearerTokenResourceHclConfig(name, "minimal+"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_bearer_token.test", "name", name),
				),
			},
		},
	})
}

func testBearerTokenResourceHclConfig(name string, configType string) string {
	if configType == "full" {
		return fmt.Sprintf(`
resource "hyperfabric_bearer_token" "test" {
	name = "%[1]s"
	description = "This bearer token is powered by Cisco Nexus Hyperfabric"
	not_after   = "2025-09-03T08:00:00.000Z"
    not_before  = "2024-09-03T08:00:00.000Z"
    scope       = "ADMIN"
}
`, name)
	} else if configType == "clear" {
		return fmt.Sprintf(`
resource "hyperfabric_bearer_token" "test" {
	name = "%[1]s"
	description = ""
	not_after   = "2025-09-03T08:00:00.000Z"
    not_before  = "2024-09-03T08:00:00.000Z"
    scope       = "ADMIN"
}
`, name)
	} else if configType == "minimal+" {
		return fmt.Sprintf(`
resource "hyperfabric_bearer_token" "test" {
name = "%[1]s"
scope       = "ADMIN"
}
`, name)
	} else {
		return fmt.Sprintf(`
resource "hyperfabric_bearer_token" "test" {
	name = "%[1]s"
}
`, name)
	}
}
