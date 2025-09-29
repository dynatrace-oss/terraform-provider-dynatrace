/**
* @license
* Copyright 2024 Dynatrace LLC
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

package openpipeline_test

import (
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/api"
)

func TestAccOpenPipeline(t *testing.T) {
	// Notes:
	// - Seems like no kind currently supports the "azure_log_forwarding_processor" and "security_event_extraction_processor"
	t.Skip("Tests skipped until REST API doesn't respond back with version conflicts any more")
	api.TestAcc(t)
}
