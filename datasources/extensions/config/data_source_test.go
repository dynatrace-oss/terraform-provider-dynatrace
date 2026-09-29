//go:build integration

/*
 * @license
 * Copyright 2026 Dynatrace LLC
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package extension_config_test

import (
	"fmt"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceExtensionConfigurations(t *testing.T) {
	identifier := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	identifier2 := identifier + "-2"

	api.TestAcc(t, api.TestAccOptions{
		Identifier: identifier,
		Check: func(s *terraform.State) error {
			desc1 := s.RootModule().Outputs["description1"].Value
			desc2 := s.RootModule().Outputs["description2"].Value

			if desc1 != identifier {
				return fmt.Errorf("expected %s, got %s", identifier, desc1)
			}
			if desc2 != identifier2 {
				return fmt.Errorf("expected %s, got %s", identifier2, desc2)
			}
			return nil
		},
	})
}
