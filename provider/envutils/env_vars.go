package envutils

// --- General Provider ---

// DynatraceDebug enables debug logging for the provider.
var DynatraceDebug = BoolEnvVar{
	Key:          "DYNATRACE_DEBUG",
	DefaultValue: false,
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

// DynatraceHTTPOAuth enables OAuth authentication for HTTP requests.
var DynatraceHTTPOAuth = BoolEnvVar{
	Key:          "DYNATRACE_HTTP_OAUTH",
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

// DTDebugIAMBearer enables debug logging for IAM bearer tokens.
var DTDebugIAMBearer = BoolEnvVar{
	Key:          "DT_DEBUG_IAM_BEARER",
	DefaultValue: false,
}

// DTDebugGetOk enables debug logging for successful GET requests.
var DTDebugGetOk = BoolEnvVar{
	Key:          "DT_DEBUG_GET_OK",
	DefaultValue: false,
}

// --- IAM ---

// DynatraceDisableIAMRateLimiter disables the IAM rate limiter.
var DynatraceDisableIAMRateLimiter = BoolEnvVar{
	Key:          "DYNATRACE_DISABLE_IAM_RATE_LIMITER",
	DefaultValue: false,
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

// CacheOfflineMode enables offline mode using cached data.
var CacheOfflineMode = BoolEnvVar{
	Key:          "CACHE_OFFLINE_MODE",
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

// --- Documents ---

// DTDocumentsIgnoreUnexpectedEOF ignores unexpected EOF errors when managing documents.
var DTDocumentsIgnoreUnexpectedEOF = BoolEnvVar{
	Key:          "DT_DOCUMENTS_IGNORE_UNEXPECTED_EOF",
	DefaultValue: false,
}

// --- Entity ---

// DynatraceDisableEntityCache disables the entity cache.
var DynatraceDisableEntityCache = BoolEnvVar{
	Key:          "DYNATRACE_DISABLE_ENTITY_CACHE",
	DefaultValue: false,
}

// --- Tags ---

// DynatraceTagsErrZeroMatched returns an error when a tag query matches zero entities.
var DynatraceTagsErrZeroMatched = BoolEnvVar{
	Key:          "DYNATRACE_TAGS_ERR_ZERO_MATCHED",
	DefaultValue: false,
}

// --- Resources ---

// DTTerraformImport enables Terraform import mode.
var DTTerraformImport = BoolEnvVar{
	Key:          "DT_TERRAFORM_IMPORT",
	DefaultValue: false,
}

// DynatraceForceNewOnHeaders forces resource recreation when HTTP headers change.
var DynatraceForceNewOnHeaders = BoolEnvVar{
	Key:          "DYNATRACE_FORCE_NEW_ON_HEADERS",
	DefaultValue: false,
}

// DynatraceGoldenStateEnabled enables golden state tracking for resources.
var DynatraceGoldenStateEnabled = BoolEnvVar{
	Key:          "DYNATRACE_GOLDEN_STATE_ENABLED",
	DefaultValue: false,
}

// --- Host Monitoring ---

// DynatraceHostMonitoringOffline allows managing host monitoring for offline hosts.
var DynatraceHostMonitoringOffline = BoolEnvVar{
	Key:          "DYNATRACE_HOST_MONITORING_OFFLINE",
	DefaultValue: false,
}

// DynatraceHostMonitoringWarnings enables warnings for host monitoring issues.
var DynatraceHostMonitoringWarnings = BoolEnvVar{
	Key:          "DYNATRACE_HOST_MONITORING_WARNINGS",
	DefaultValue: false,
}

// --- Workflows ---

// DynatraceWorkflowTasksUseTypeList uses a type list for workflow task filtering.
var DynatraceWorkflowTasksUseTypeList = BoolEnvVar{
	Key:          "DYNATRACE_WORKFLOW_TASKS_USE_TYPE_LIST",
	DefaultValue: false,
}

// --- Export ---

// DynatraceCleanTargetFolder enables cleaning of the target folder before export.
var DynatraceCleanTargetFolder = BoolEnvVar{
	Key:          "DYNATRACE_CLEAN_TARGET_FOLDER",
	DefaultValue: false,
}

// DynatraceNoRefreshOnImport skips refreshing resources after export.
var DynatraceNoRefreshOnImport = BoolEnvVar{
	Key:          "DYNATRACE_NO_REFRESH_ON_IMPORT",
	DefaultValue: false,
}

// DynatraceQuickInit enables quick initialization during export.
var DynatraceQuickInit = BoolEnvVar{
	Key:          "DYNATRACE_QUICK_INIT",
	DefaultValue: false,
}

// DynatraceUltraParallel enables ultra-parallel export mode.
var DynatraceUltraParallel = BoolEnvVar{
	Key:          "DYNATRACE_ULTRA_PARALLEL",
	DefaultValue: false,
}

// DynatraceShorterNames enables shortening resource names longer than would exceed 240 characters in length in generated HCL.
var DynatraceShorterNames = BoolEnvVar{
	Key:          "DYNATRACE_SHORTER_NAMES",
	DefaultValue: false,
}

// DynatraceEnableExportDashboard enables export of dashboard resources.
var DynatraceEnableExportDashboard = BoolEnvVar{
	Key:          "DYNATRACE_ENABLE_EXPORT_DASHBOARD",
	DefaultValue: false,
}

// DynatraceAtomicDependencies exports dependencies atomically.
var DynatraceAtomicDependencies = BoolEnvVar{
	Key:          "DYNATRACE_ATOMIC_DEPENDENCIES",
	DefaultValue: false,
}

// DynatraceFormatHCLFiles enables formatting of generated HCL files.
var DynatraceFormatHCLFiles = BoolEnvVar{
	Key:          "DYNATRACE_FORMAT_HCL_FILES",
	DefaultValue: false,
}

// DynatraceHCLNoFormat disables HCL formatting during export.
var DynatraceHCLNoFormat = BoolEnvVar{
	Key:          "DYNATRACE_HCL_NO_FORMAT",
	DefaultValue: false,
}

// DynatraceNameReplaceDash replaces dashes with underscores in resource names.
var DynatraceNameReplaceDash = BoolEnvVar{
	Key:          "DYNATRACE_NAME_REPLACE_DASH",
	DefaultValue: false,
}

// DynatraceBuildAddressFiles generates address files during export.
var DynatraceBuildAddressFiles = BoolEnvVar{
	Key:          "DYNATRACE_BUILD_ADDRESS_FILES",
	DefaultValue: false,
}

// DynatraceIgnoreChangesRequiresAttention enables adding a lifecycle { ignore_changes = [...] } block
// to exported resources that contain sensitive fields (passwords, secrets, etc.) — listing those sensitive
// attributes so Terraform won't overwrite values the user may have set manually after the export.
var DynatraceIgnoreChangesRequiresAttention = BoolEnvVar{
	Key:          "DYNATRACE_IGNORE_CHANGES_REQUIRES_ATTENTION",
	DefaultValue: false,
}

// DynatracePrevStateOn enables keeping resource identifiers stable between runs so downstream Terraform references don't break.
var DynatracePrevStateOn = BoolEnvVar{
	Key:          "DYNATRACE_PREV_STATE_ON",
	DefaultValue: false,
}

// --- Migration ---

// Migration enables migration mode.
var Migration = BoolEnvVar{
	Key:          "MIGRATION",
	DefaultValue: false,
}
