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
