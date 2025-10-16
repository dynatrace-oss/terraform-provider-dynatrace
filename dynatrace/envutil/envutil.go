/**
* @license
* Copyright 2020 Dynatrace LLC
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

package envutil

import (
"os"
"strconv"
)

// Environment variable name constants
const (
// Export package environment variables
EnvNoRefreshOnImport              = "DYNATRACE_NO_REFRESH_ON_IMPORT"
EnvQuickInit                      = "DYNATRACE_QUICK_INIT"
EnvUltraParallel                  = "DYNATRACE_ULTRA_PARALLEL"
EnvParallel                       = "DYNATRACE_PARALLEL"
EnvAtomicDependencies             = "DYNATRACE_ATOMIC_DEPENDENCIES"
EnvHCLNoFormat                    = "DYNATRACE_HCL_NO_FORMAT"
EnvNameReplaceDash                = "DYNATRACE_NAME_REPLACE_DASH"
EnvFormatHCLFiles                 = "DYNATRACE_FORMAT_HCL_FILES"
EnvCleanTargetFolder              = "DYNATRACE_CLEAN_TARGET_FOLDER"
EnvShorterNames                   = "DYNATRACE_SHORTER_NAMES"
EnvEnableExportDashboard          = "DYNATRACE_ENABLE_EXPORT_DASHBOARD"
EnvPrevStateOn                    = "DYNATRACE_PREV_STATE_ON"
EnvIgnoreChangesRequiresAttention = "DYNATRACE_IGNORE_CHANGES_REQUIRES_ATTENTION"
)

// GetBoolEnv reads a boolean environment variable with a default value.
// It uses strconv.ParseBool which accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.
func GetBoolEnv(key string, defaultValue bool) bool {
if value := os.Getenv(key); value != "" {
if parsed, err := strconv.ParseBool(value); err == nil {
return parsed
}
}
return defaultValue
}
