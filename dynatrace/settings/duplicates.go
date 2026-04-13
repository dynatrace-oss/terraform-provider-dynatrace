/**
* @license
* Copyright 2025 Dynatrace LLC
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

package settings

import (
	"slices"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/envutils"
)

const VALUE_ALL = "ALL"

func RejectDuplicate(resourceNames ...string) bool {
	return envVarContains(envutils.DynatraceDuplicateReject, resourceNames...)
}

func HijackDuplicate(resourceNames ...string) bool {
	return envVarContains(envutils.DynatraceDuplicateHijack, resourceNames...)
}

func envVarContains(envVar envutils.StringEnvVar, search ...string) bool {
	svalues := envVar.Get()
	if len(svalues) == 0 {
		return false
	}
	if svalues == VALUE_ALL {
		return true
	}
	values := strings.SplitSeq(svalues, ",")
	for value := range values {
		value = strings.TrimSpace(value)
		if slices.Contains(search, value) {
			return true
		}
	}
	return false
}
