// Copyright 2024 Cisco Systems, Inc. and its affiliates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserResource(t *testing.T) {
	emailName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	emailDomain := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	email := emailName + "@" + emailDomain + ".tld"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify provided and default Hyperfabric values.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: User - Create with minimum config and verify provided and default Hyperfabric values.")
				},
				Config:             testUserResourceHclConfig(email, "minimal"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_user.test", "email", email),
				),
			},
			// Update with all config and verify provided values.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: User - Update with all config and verify provided values.")
				},
				Config:             testUserResourceHclConfig(email, "full"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_user.test", "email", email),
					resource.TestCheckResourceAttr("hyperfabric_user.test", "enabled", "true"),
					resource.TestCheckResourceAttr("hyperfabric_user.test", "role", "ADMIN"),
					resource.TestCheckResourceAttr("hyperfabric_user.test", "labels.#", "2"),
				),
			},
			// Update with minimum config and verify config is unchanged.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: User - Update with minimum config and verify config is unchanged.")
				},
				Config:             testUserResourceHclConfig(email, "minimal"),
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_user.test", "email", email),
					resource.TestCheckResourceAttr("hyperfabric_user.test", "enabled", "true"),
					resource.TestCheckResourceAttr("hyperfabric_user.test", "role", "ADMIN"),
					resource.TestCheckResourceAttr("hyperfabric_user.test", "labels.#", "2"),
				),
			},
			// ImportState testing with pre-existing Id.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: User - ImportState testing with pre-existing Id.")
				},
				ResourceName:      "hyperfabric_user.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// ImportState testing with email.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: User - ImportState testing with email.")
				},
				ResourceName:      "hyperfabric_user.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     email,
			},
			// Update with config containing all optional attributes with empty values and verify config is cleared.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: User - Update with config containing all optional attributes with empty values and verify config is cleared.")
				},
				Config:             testUserResourceHclConfig(email, "clear"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_user.test", "email", email),
					resource.TestCheckResourceAttr("hyperfabric_user.test", "enabled", "true"),
					resource.TestCheckResourceAttr("hyperfabric_user.test", "role", "ADMIN"),
					resource.TestCheckResourceAttr("hyperfabric_user.test", "labels.#", "0"),
				),
			},
			// Run Plan Only with minimal config and check that plan is empty.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: User - Run Plan Only with minimal config and check that plan is empty.")
				},
				Config:             testUserResourceHclConfig(email, "minimal"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_user.test", "email", email),
				),
			},
		},
	})
}

func testUserResourceHclConfig(email string, configType string) string {
	if configType == "full" {
		return fmt.Sprintf(`
resource "hyperfabric_user" "test" {
	email = "%[1]s"
	enabled     = true
	role    = "ADMIN"
	labels = [
		"sj01-1-101-AAA01",
		"blue"
	]
}
`, email)
	} else if configType == "clear" {
		return fmt.Sprintf(`
resource "hyperfabric_user" "test" {
	email = "%[1]s"
	enabled     = true
	role    = "ADMIN"
	labels = []
}
`, email)
	} else {
		return fmt.Sprintf(`
resource "hyperfabric_user" "test" {
	email = "%[1]s"
}
`, email)
	}
}
