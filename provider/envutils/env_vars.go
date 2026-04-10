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

// DynatraceMaxHTTPWorkers sets the maximum number of concurrent HTTP workers.
var DynatraceMaxHTTPWorkers = StringEnvVar{
	Key:          "DYNATRACE_MAX_HTTP_WORKERS",
	DefaultValue: "",
}

// DTRestDebugLog sets the file path for REST debug logging.
var DTRestDebugLog = StringEnvVar{
	Key:          "DT_REST_DEBUG_LOG",
	DefaultValue: "",
}

// --- IAM ---

// DynatraceDisableIAMRateLimiter disables the IAM rate limiter.
var DynatraceDisableIAMRateLimiter = BoolEnvVar{
	Key:          "DYNATRACE_DISABLE_IAM_RATE_LIMITER",
	DefaultValue: false,
}

// DynatraceIAMRateLimiterRate sets the IAM rate limiter rate in requests per second.
var DynatraceIAMRateLimiterRate = StringEnvVar{
	Key:          "DYNATRACE_IAM_RATE_LIMITER_RATE",
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

// DTCacheDeleteOnLaunch specifies cache entries to delete on provider launch.
var DTCacheDeleteOnLaunch = StringEnvVar{
	Key:          "DT_CACHE_DELETE_ON_LAUNCH",
	DefaultValue: "",
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

// DTMgmzSuccesses is the number of consecutive successes required for management zone readiness. Values outside [Min, Max] are clamped to the nearest boundary.
var DTMgmzSuccesses = ClampedIntEnvVar{
	Key:          "DT_MGMZ_SUCCESSES",
	DefaultValue: 5,
	Min:          5,
	Max:          100,
}

// --- Documents ---

// DTDocumentsIgnoreUnexpectedEOF ignores unexpected EOF errors when managing documents.
var DTDocumentsIgnoreUnexpectedEOF = BoolEnvVar{
	Key:          "DT_DOCUMENTS_IGNORE_UNEXPECTED_EOF",
	DefaultValue: false,
}

// --- Custom Device ---

// DTCustomDeviceApplyTimeout is the timeout in seconds for custom device apply operations. Values outside [Min, Max] are clamped to the nearest boundary.
var DTCustomDeviceApplyTimeout = ClampedIntEnvVar{
	Key:          "DT_CUSTOM_DEVICE_APPLY_TIMEOUT",
	DefaultValue: 100,
	Min:          100,
	Max:          500,
}

// --- Entity ---

// DynatraceDisableEntityCache disables the entity cache.
var DynatraceDisableEntityCache = BoolEnvVar{
	Key:          "DYNATRACE_DISABLE_ENTITY_CACHE",
	DefaultValue: false,
}

// --- Duplicates ---

// DynatraceDuplicateReject specifies resource types for which duplicate detection rejects the import.
var DynatraceDuplicateReject = StringEnvVar{
	Key:          "DYNATRACE_DUPLICATE_REJECT",
	DefaultValue: "",
}

// DynatraceDuplicateHijack specifies resource types for which duplicate detection hijacks the existing resource.
var DynatraceDuplicateHijack = StringEnvVar{
	Key:          "DYNATRACE_DUPLICATE_HIJACK",
	DefaultValue: "",
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

// --- Dashboards ---

// DynatraceDashboardTests configures dashboard test behavior.
var DynatraceDashboardTests = StringEnvVar{
	Key:          "DYNATRACE_DASHBOARD_TESTS",
	DefaultValue: "",
}

// --- DQL ---

// DynatraceDQLPollSleepDuration is the sleep duration in milliseconds between DQL poll attempts. Values outside [Min, Max] are clamped to the nearest boundary.
var DynatraceDQLPollSleepDuration = ClampedIntEnvVar{
	Key:          "DYNATRACE_DQL_POLL_SLEEP_DURATION",
	DefaultValue: 5000,
	Min:          0,
	Max:          60000,
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

// DynatraceHostMonitoringStrictUpdateRetries sets the number of retries for strict host monitoring updates.
var DynatraceHostMonitoringStrictUpdateRetries = StringEnvVar{
	Key:          "DYNATRACE_HOST_MONITORING_STRICT_UPDATE_RETRIES",
	DefaultValue: "",
}

// --- Workflows ---

// DynatraceWorkflowTasksUseTypeList uses a type list for workflow task filtering.
var DynatraceWorkflowTasksUseTypeList = BoolEnvVar{
	Key:          "DYNATRACE_WORKFLOW_TASKS_USE_TYPE_LIST",
	DefaultValue: false,
}

// --- HCL / Terraform Generation ---

// DynatraceHeredoc configures heredoc usage in generated HCL.
var DynatraceHeredoc = StringEnvVar{
	Key:          "DYNATRACE_HEREDOC",
	DefaultValue: "",
}

// --- Export ---

// DynatraceTargetFolder sets the target folder for export output.
var DynatraceTargetFolder = StringEnvVar{
	Key:          "DYNATRACE_TARGET_FOLDER",
	DefaultValue: "",
}

// DynatraceCleanTargetFolder enables cleaning of the target folder before export.
var DynatraceCleanTargetFolder = BoolEnvVar{
	Key:          "DYNATRACE_CLEAN_TARGET_FOLDER",
	DefaultValue: false,
}

// DynatraceProviderSource sets the provider source for generated Terraform files.
var DynatraceProviderSource = StringEnvVar{
	Key:          "DYNATRACE_PROVIDER_SOURCE",
	DefaultValue: "",
}

// DynatraceProviderVersion sets the provider version for generated Terraform files.
var DynatraceProviderVersion = StringEnvVar{
	Key:          "DYNATRACE_PROVIDER_VERSION",
	DefaultValue: "",
}

// DynatraceCustomProviderLocation sets the custom provider binary location.
var DynatraceCustomProviderLocation = StringEnvVar{
	Key:          "DYNATRACE_CUSTOM_PROVIDER_LOCATION",
	DefaultValue: "",
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

// DynatraceParallel sets the number of parallel export workers.
var DynatraceParallel = StringEnvVar{
	Key:          "DYNATRACE_PARALLEL",
	DefaultValue: "",
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

// DynatraceMigrationCacheFolder sets the migration cache folder path.
var DynatraceMigrationCacheFolder = StringEnvVar{
	Key:          "DYNATRACE_MIGRATION_CACHE_FOLDER",
	DefaultValue: "",
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

// DynatraceExportIgnoreResources is a comma-separated list of resource types to exclude from export.
var DynatraceExportIgnoreResources = StringEnvVar{
	Key:          "DYNATRACE_EXPORT_IGNORE_RESOURCES",
	DefaultValue: "",
}

// DynatraceIgnoreChangesRequiresAttention enables adding a lifecycle { ignore_changes = [...] } block
// to exported resources that contain sensitive fields (passwords, secrets, etc.) — listing those sensitive
// attributes so Terraform won't overwrite values the user may have set manually after the export.
var DynatraceIgnoreChangesRequiresAttention = BoolEnvVar{
	Key:          "DYNATRACE_IGNORE_CHANGES_REQUIRES_ATTENTION",
	DefaultValue: false,
}

// DynatraceImportStatePath sets the path for importing Terraform state.
var DynatraceImportStatePath = StringEnvVar{
	Key:          "DYNATRACE_IMPORT_STATE_PATH",
	DefaultValue: "",
}

// DynatracePrevStateOn enables keeping resource identifiers stable between runs so downstream Terraform references don't break.
var DynatracePrevStateOn = BoolEnvVar{
	Key:          "DYNATRACE_PREV_STATE_ON",
	DefaultValue: false,
}

// DynatracePrevStatePathThis sets the path to the previous state for the current environment.
var DynatracePrevStatePathThis = StringEnvVar{
	Key:          "DYNATRACE_PREV_STATE_PATH_THIS",
	DefaultValue: "",
}

// DynatracePrevStatePathLinked sets the path to the previous state for linked environments.
var DynatracePrevStatePathLinked = StringEnvVar{
	Key:          "DYNATRACE_PREV_STATE_PATH_LINKED",
	DefaultValue: "",
}

// --- Source Environment (multi-key fallback) ---

// DynatraceSourceEnvURL is the URL of the source Dynatrace environment for migration.
var DynatraceSourceEnvURL = MultiStringEnvVar{
	Keys: []string{
		"DYNATRACE_SOURCE_ENV_URL",
		"DT_SOURCE_ENV_URL",
		"DYNATRACE_SOURCE_ENVIRONMENT_URL",
		"DT_SOURCE_ENVIRONMENT_URL",
	},
}

// DynatraceSourceAPIToken is the API token for the source Dynatrace environment.
var DynatraceSourceAPIToken = MultiStringEnvVar{
	Keys: []string{
		"DYNATRACE_SOURCE_API_TOKEN",
		"DT_SOURCE_API_TOKEN",
	},
}

// DynatraceSourceClientID is the OAuth client ID for the source Dynatrace environment.
var DynatraceSourceClientID = MultiStringEnvVar{
	Keys: []string{
		"DT_SOURCE_CLIENT_ID",
		"DYNATRACE_SOURCE_CLIENT_ID",
	},
}

// DynatraceSourceAccountID is the account ID for the source Dynatrace environment.
var DynatraceSourceAccountID = MultiStringEnvVar{
	Keys: []string{
		"DT_SOURCE_ACCOUNT_ID",
		"DYNATRACE_SOURCE_ACCOUNT_ID",
	},
}

// DynatraceSourceClientSecret is the OAuth client secret for the source Dynatrace environment.
var DynatraceSourceClientSecret = MultiStringEnvVar{
	Keys: []string{
		"DT_SOURCE_CLIENT_SECRET",
		"DYNATRACE_SOURCE_CLIENT_SECRET",
	},
}

// DynatraceSourcePlatformToken is the platform token for the source Dynatrace environment.
var DynatraceSourcePlatformToken = MultiStringEnvVar{
	Keys: []string{
		"DYNATRACE_SOURCE_PLATFORM_TOKEN",
		"DT_SOURCE_PLATFORM_TOKEN",
	},
}

// --- Testing ---

// DynatraceEnvURL sets the Dynatrace environment URL for acceptance tests.
var DynatraceEnvURL = StringEnvVar{
	Key:          "DYNATRACE_ENV_URL",
	DefaultValue: "",
}

// DynatraceAPIToken sets the Dynatrace API token for acceptance tests.
var DynatraceAPIToken = StringEnvVar{
	Key:          "DYNATRACE_API_TOKEN",
	DefaultValue: "",
}

// TFAcc enables Terraform acceptance tests when set to a non-empty value.
var TFAcc = StringEnvVar{
	Key:          "TF_ACC",
	DefaultValue: "",
}

// --- Migration ---

// Migration enables migration mode.
var Migration = BoolEnvVar{
	Key:          "MIGRATION",
	DefaultValue: false,
}

// --- Synthetic Monitors ---

// DynatraceCreateConfirmSyntheticMonitorsV2 is the number of confirmation retries when creating synthetic monitors v2. Falls back to the default if the value is outside [Min, Max].
var DynatraceCreateConfirmSyntheticMonitorsV2 = BoundedIntEnvVar{
	Key:          "DYNATRACE_CREATE_CONFIRM_SYNTHETIC_MONITORS_V2",
	DefaultValue: 8,
	Min:          1,
	Max:          50,
}

// --- Web Application ---

// DynatraceCreateConfirmWebApplication is the number of confirmation retries when creating web applications. Falls back to the default if the value is outside [Min, Max].
var DynatraceCreateConfirmWebApplication = BoundedIntEnvVar{
	Key:          "DYNATRACE_CREATE_CONFIRM_WEB_APPLICATION",
	DefaultValue: 280,
	Min:          20,
	Max:          500,
}

// --- Custom Tags ---

// DynatraceMaxConcurrentCustomTagListRequests is the maximum number of concurrent requests for listing custom tags. Falls back to the default if the value is outside [Min, Max].
var DynatraceMaxConcurrentCustomTagListRequests = BoundedIntEnvVar{
	Key:          "DYNATRACE_MAX_CONCURRENT_CUSTOM_TAG_LIST_REQUESTS",
	DefaultValue: 4,
	Min:          1,
	Max:          20,
}
