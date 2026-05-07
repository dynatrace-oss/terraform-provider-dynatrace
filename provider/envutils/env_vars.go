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

package envutils

// --- General Provider ---

// DynatraceDebug enables debug logging for the provider.
var DynatraceDebug = BoolEnvVar{
	Key:          "DYNATRACE_DEBUG",
	DefaultValue: false,
}

// DynatraceLogDebugPrefix sets the prefix filter for debug log entries.
var DynatraceLogDebugPrefix = StringEnvVar{
	Key:          "DYNATRACE_LOG_DEBUG_PREFIX",
	DefaultValue: "",
}

// DynatraceLogHTTP sets the file path for HTTP request/response logging.
// If set to "false" (note case) or an empty value, HTTP logging is disabled.
// If set to "true" (note case), HTTP logs will be written to standard error with a "[HTTP]" prefix.
// If set to "stdout" (note case), HTTP logs will be written to standard output with a "[HTTP]" prefix.
// If set to any other non-empty value, it is treated as a file path where HTTP logs will be written.
var DynatraceLogHTTP = StringEnvVar{
	Key:          "DYNATRACE_LOG_HTTP",
	DefaultValue: "",
}

// --- HTTP ---

// DynatraceHTTPInsecure disables TLS certificate verification for HTTP requests.
var DynatraceHTTPInsecure = BoolEnvVar{
	Key:          "DYNATRACE_HTTP_INSECURE",
	DefaultValue: false,
}

// DynatraceHTTPLegacy uses the legacy HTTP client.
var DynatraceHTTPLegacy = BoolEnvVar{
	Key:          "DYNATRACE_HTTP_LEGACY",
	DefaultValue: false,
}

// DynatraceHTTPOAuthPreference prefers OAuth over API token authentication when both are available.
var DynatraceHTTPOAuthPreference = BoolEnvVar{
	Key:          "DYNATRACE_HTTP_OAUTH_PREFERENCE",
	DefaultValue: false,
}

// DynatraceHTTPResponse logs HTTP responses.
var DynatraceHTTPResponse = BoolEnvVar{
	Key:          "DYNATRACE_HTTP_RESPONSE",
	DefaultValue: false,
}

// DTDebugGetOk enables debug logging for successful GET requests.
var DTDebugGetOk = BoolEnvVar{
	Key:          "DT_DEBUG_GET_OK",
	DefaultValue: false,
}

// DynatraceMaxHTTPWorkers sets the maximum number of concurrent HTTP workers. Values outside [Min, Max] are clamped to the nearest boundary.
var DynatraceMaxHTTPWorkers = ClampedIntEnvVar{
	Key:          "DYNATRACE_MAX_HTTP_WORKERS",
	DefaultValue: 20,
	Min:          1,
	Max:          50,
}

// DTRestDebugLog sets the file path for REST debug logging.
var DTRestDebugLog = StringEnvVar{
	Key:          "DT_REST_DEBUG_LOG",
	DefaultValue: "",
}

// --- Settings 2.0 ---

// DTNoRepairInput disables automatic repair of invalid settings input.
var DTNoRepairInput = BoolEnvVar{
	Key:          "DT_NO_REPAIR_INPUT",
	DefaultValue: false,
}

// DynatraceDisableOrderingSupport disables ordering support for settings resources.
var DynatraceDisableOrderingSupport = BoolEnvVar{
	Key:          "DYNATRACE_DISABLE_ORDERING_SUPPORT",
	DefaultValue: false,
}

// --- Cache ---

// DTCacheFolder sets the path to the cache folder.
var DTCacheFolder = StringEnvVar{
	Key:          "DT_CACHE_FOLDER",
	DefaultValue: "",
}

// CacheOfflineMode enables offline mode using cached data.
var CacheOfflineMode = BoolEnvVar{
	Key:          "CACHE_OFFLINE_MODE",
	DefaultValue: false,
}

// DTCacheDeleteOnLaunch enables deleting the cache on provider launch.
var DTCacheDeleteOnLaunch = BoolEnvVar{
	Key:          "DT_CACHE_DELETE_ON_LAUNCH",
	DefaultValue: false,
}

// DTNoCacheCleanup specifies whether to skip automatic cache cleanup on provider shutdown.
var DTNoCacheCleanup = BoolEnvVar{
	Key:          "DT_NO_CACHE_CLEANUP",
	DefaultValue: false,
}

// DynatraceInMemoryTarFolders uses in-memory tar archives for cache folders.
var DynatraceInMemoryTarFolders = BoolEnvVar{
	Key:          "DYNATRACE_IN_MEMORY_TAR_FOLDERS",
	DefaultValue: false,
}

// --- Backwards Compatibility ---

// DTBackwardsCompatibility enables backwards compatibility mode.
var DTBackwardsCompatibility = BoolEnvVar{
	Key:          "DT_BACKWARDS_COMPATIBILITY",
	DefaultValue: false,
}

// --- Buckets ---

// DTBucketsIgnoreUnexpectedEOF ignores unexpected EOF errors when managing buckets.
var DTBucketsIgnoreUnexpectedEOF = BoolEnvVar{
	Key:          "DT_BUCKETS_IGNORE_UNEXPECTED_EOF",
	DefaultValue: false,
}

// DTBucketsRetries is the number of retries when waiting for bucket readiness. Values outside [Min, Max] are clamped to the nearest boundary.
var DTBucketsRetries = ClampedIntEnvVar{
	Key:          "DT_BUCKETS_RETRIES",
	DefaultValue: 180,
	Min:          180,
	Max:          360,
}

// DTBucketsNumSuccesses is the number of consecutive successes required for bucket readiness. Values outside [Min, Max] are clamped to the nearest boundary.
var DTBucketsNumSuccesses = ClampedIntEnvVar{
	Key:          "DT_BUCKETS_NUM_SUCCESSES",
	DefaultValue: 10,
	Min:          10,
	Max:          50,
}

// --- Management Zones ---

// DTMgmzRetries is the number of retries when waiting for management zone readiness. Values outside [Min, Max] are clamped to the nearest boundary.
var DTMgmzRetries = ClampedIntEnvVar{
	Key:          "DT_MGMZ_RETRIES",
	DefaultValue: 50,
	Min:          50,
	Max:          600,
}
