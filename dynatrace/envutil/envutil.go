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

// API package environment variables
EnvDocumentsIgnoreUnexpectedEOF = "DT_DOCUMENTS_IGNORE_UNEXPECTED_EOF"
EnvBucketsIgnoreUnexpectedEOF   = "DT_BUCKETS_IGNORE_UNEXPECTED_EOF"
EnvDebugIAMBearer               = "DT_DEBUG_IAM_BEARER"
EnvDisableIAMRateLimiter        = "DYNATRACE_DISABLE_IAM_RATE_LIMITER"
EnvWorkflowTasksUseTypeList     = "DYNATRACE_WORKFLOW_TASKS_USE_TYPE_LIST"
EnvTagsErrZeroMatched           = "DYNATRACE_TAGS_ERR_ZERO_MATCHED"
EnvDisableEntityCache           = "DYNATRACE_DISABLE_ENTITY_CACHE"
EnvForceNewOnHeaders            = "DYNATRACE_FORCE_NEW_ON_HEADERS"
EnvHostMonitoringWarnings       = "DYNATRACE_HOST_MONITORING_WARNINGS"
EnvHostMonitoringOffline        = "DYNATRACE_HOST_MONITORING_OFFLINE"
EnvBackwardsCompatibility       = "DT_BACKWARDS_COMPATIBILITY"

// Address package environment variables
EnvBuildAddressFiles = "DYNATRACE_BUILD_ADDRESS_FILES"

// REST package environment variables
EnvHTTPLegacy          = "DYNATRACE_HTTP_LEGACY"
EnvHTTPOAuthPreference = "DYNATRACE_HTTP_OAUTH_PREFERENCE"
EnvHTTPInsecure        = "DYNATRACE_HTTP_INSECURE"
EnvHTTPResponse        = "DYNATRACE_HTTP_RESPONSE"
EnvHTTPOAuth           = "DYNATRACE_HTTP_OAUTH"

// Settings services cache environment variables
EnvInMemoryTarFolders = "DYNATRACE_IN_MEMORY_TAR_FOLDERS"

// Settings20 service environment variables
EnvDisableOrderingSupport = "DYNATRACE_DISABLE_ORDERING_SUPPORT"
EnvNoRepairInput          = "DT_NO_REPAIR_INPUT"

// Terraform/Resources environment variables
EnvHeredoc            = "DYNATRACE_HEREDOC"
EnvDebugGetOk         = "DT_DEBUG_GET_OK"
EnvGoldenStateEnabled = "DYNATRACE_GOLDEN_STATE_ENABLED"
EnvTerraformImport    = "DT_TERRAFORM_IMPORT"

// Provider/Logging environment variables
EnvDebug = "DYNATRACE_DEBUG"
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
