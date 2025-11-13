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

package version

import "fmt"

// The following variables should be set during the build, for example by GoReleaser.
// This can be accomplished using `ldflags`, for example `-X github.com/dynatrace-oss/terraform-provider-dynatrace/provider/version.Version=1.2.3`.
var (
	// Version is a string representing the current version of the Dynatrace Terraform provider.
	// This should be in `major.minor.patch` format, with prereleases having an optional suffix introduced by a dash, for example, `1.2.0-beta`.
	// The value ~>1.0 should allow the latest 1.x version during development.
	Version = "~>1.0"

	// OperatingSystem is a string representing the operating system the provider was built for.
	OperatingSystem = "(unknown-os)"

	// CPUArchitecture is a string representing the CPU archictecture the provider was built for.
	CPUArchitecture = "(unknown-architecture)"
)

// UserAgent returns a string identifying the current version of the Dynatrace Terraform provider.
// This is intended for use in HTTP requests.
func UserAgent() string {
	return fmt.Sprintf("Dynatrace Terraform Provider/%s %s %s", Version, OperatingSystem, CPUArchitecture)
}
