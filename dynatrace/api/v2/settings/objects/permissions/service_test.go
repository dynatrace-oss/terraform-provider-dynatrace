//go:build integration

/*
 * @license
 * Copyright 2025 Dynatrace LLC
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

package permissions_test

import (
	"testing"

	settingsapi "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/api"
)

func TestAC(t *testing.T) {
	// for the example that we show in the documentation we may not want to add a lifecycle to ignore the default group changes. Therefore, setting ExpectNonEmptyPlan to true
	settingsapi.TestAcc(t, settingsapi.TestAccOptions{ExpectNonEmptyPlan: true})
}
