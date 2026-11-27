//go:build integration

/**
* @license
* Copyright 2026 Dynatrace LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package credentialsv2_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceCredentialsV2(t *testing.T) {
	api.TestAcc(t, api.TestAccOptions{Check: matchedIDsEqualExpectedIDs})
}

func matchedIDsEqualExpectedIDs(state *terraform.State) error {
	outputs := state.RootModule().Outputs
	matched, found := outputs["matched_ids"]
	if !found {
		return fmt.Errorf("output matched_ids is missing")
	}
	expected, found := outputs["expected_ids"]
	if !found {
		return fmt.Errorf("output expected_ids is missing")
	}
	if !reflect.DeepEqual(matched.Value, expected.Value) {
		return fmt.Errorf("expected credentials %v, got %v", expected.Value, matched.Value)
	}
	return nil
}
