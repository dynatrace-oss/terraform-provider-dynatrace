---
layout: ""
page_title: "Environment Variables"
description: |-
  Environment variables for configuring the Dynatrace Terraform provider.
---

# Environment Variables

The Dynatrace Terraform provider can be configured using the following environment variables.

## General Provider

| Variable | Default | Description |
|---|---|---|
| `DYNATRACE_DEBUG` | `false` | Enables debug logging for the provider. |
| `DYNATRACE_LOG_DEBUG_PREFIX` | `""` | Sets the prefix filter for debug log entries. |
| `DYNATRACE_LOG_HTTP` | `""` | Sets the file path for HTTP request/response logging. |

## HTTP

| Variable | Default | Description |
|---|---|---|
| `DYNATRACE_HTTP_INSECURE` | `false` | Disables TLS certificate verification for HTTP requests. |
| `DYNATRACE_HTTP_LEGACY` | `false` | Uses the legacy HTTP client. |
| `DYNATRACE_HTTP_OAUTH` | `false` | Enables OAuth authentication for HTTP requests. |
| `DYNATRACE_HTTP_OAUTH_PREFERENCE` | `false` | Prefers OAuth over API token authentication when both are available. |
| `DYNATRACE_HTTP_RESPONSE` | `false` | Logs HTTP responses. |
| `DT_DEBUG_IAM_BEARER` | `false` | Enables debug logging for IAM bearer tokens. |
| `DT_DEBUG_GET_OK` | `false` | Enables debug logging for successful GET requests. |
| `DYNATRACE_MAX_HTTP_WORKERS` | `""` | Sets the maximum number of concurrent HTTP workers. |
| `DT_REST_DEBUG_LOG` | `""` | Sets the file path for REST debug logging. |

## IAM

| Variable | Default | Description |
|---|---|---|
| `DYNATRACE_DISABLE_IAM_RATE_LIMITER` | `false` | Disables the IAM rate limiter. |
| `DYNATRACE_IAM_RATE_LIMITER_RATE` | `""` | Sets the IAM rate limiter rate in requests per second. |

## Settings 2.0

| Variable | Default | Description |
|---|---|---|
| `DT_NO_REPAIR_INPUT` | `false` | Disables automatic repair of invalid settings input. |
| `DYNATRACE_DISABLE_ORDERING_SUPPORT` | `false` | Disables ordering support for settings resources. |

## Cache

| Variable | Default | Description |
|---|---|---|
| `DT_CACHE_FOLDER` | `""` | Sets the path to the cache folder. |
| `CACHE_OFFLINE_MODE` | `false` | Enables offline mode using cached data. |
| `DT_CACHE_DELETE_ON_LAUNCH` | `""` | Specifies cache entries to delete on provider launch. |
| `DT_NO_CACHE_CLEANUP` | `false` | Disables cache cleanup on provider shutdown. |
| `DYNATRACE_IN_MEMORY_TAR_FOLDERS` | `false` | Uses in-memory tar archives for cache folders. |

## Backwards Compatibility

| Variable | Default | Description |
|---|---|---|
| `DT_BACKWARDS_COMPATIBILITY` | `false` | Enables backwards compatibility mode. |

## Buckets

| Variable | Default | Description |
|---|---|---|
| `DT_BUCKETS_IGNORE_UNEXPECTED_EOF` | `false` | Ignores unexpected EOF errors when managing buckets. |

| Variable | Default | Min | Max | Description |
|---|---|---|---|---|
| `DT_BUCKETS_RETRIES` | `180` | `180` | `360` | Number of retries when waiting for bucket readiness. Values outside [Min, Max] are clamped to the nearest boundary. |
| `DT_BUCKETS_NUM_SUCCESSES` | `10` | `10` | `50` | Number of consecutive successes required for bucket readiness. Values outside [Min, Max] are clamped to the nearest boundary. |

## Documents

| Variable | Default | Description |
|---|---|---|
| `DT_DOCUMENTS_IGNORE_UNEXPECTED_EOF` | `false` | Ignores unexpected EOF errors when managing documents. |

## Entity

| Variable | Default | Description |
|---|---|---|
| `DYNATRACE_DISABLE_ENTITY_CACHE` | `false` | Disables the entity cache. |

## Duplicates

| Variable | Default | Description |
|---|---|---|
| `DYNATRACE_DUPLICATE_REJECT` | `""` | Specifies resource types for which duplicate detection rejects the import. |
| `DYNATRACE_DUPLICATE_HIJACK` | `""` | Specifies resource types for which duplicate detection hijacks the existing resource. |

## Tags

| Variable | Default | Description |
|---|---|---|
| `DYNATRACE_TAGS_ERR_ZERO_MATCHED` | `false` | Returns an error when a tag query matches zero entities. |

## Resources

| Variable | Default | Description |
|---|---|---|
| `DT_TERRAFORM_IMPORT` | `false` | Enables Terraform import mode. |
| `DYNATRACE_FORCE_NEW_ON_HEADERS` | `false` | Forces resource recreation when HTTP headers change. |
| `DYNATRACE_GOLDEN_STATE_ENABLED` | `false` | Enables golden state tracking for resources. |

## Management Zones

| Variable | Default | Min | Max | Description |
|---|---|---|---|---|
| `DT_MGMZ_RETRIES` | `50` | `50` | `600` | Number of retries when waiting for management zone readiness. Values outside [Min, Max] are clamped to the nearest boundary. |
| `DT_MGMZ_SUCCESSES` | `5` | `5` | `100` | Number of consecutive successes required for management zone readiness. Values outside [Min, Max] are clamped to the nearest boundary. |

## Dashboards

| Variable | Default | Description |
|---|---|---|
| `DYNATRACE_DASHBOARD_TESTS` | `""` | Configures dashboard test behavior. |

## Custom Device

| Variable | Default | Min | Max | Description |
|---|---|---|---|---|
| `DT_CUSTOM_DEVICE_APPLY_TIMEOUT` | `100` | `100` | `500` | Timeout in seconds for custom device apply operations. Values outside [Min, Max] are clamped to the nearest boundary. |

## DQL

| Variable | Default | Min | Max | Description |
|---|---|---|---|---|
| `DYNATRACE_DQL_POLL_SLEEP_DURATION` | `5000` | `0` | `60000` | Sleep duration in milliseconds between DQL poll attempts. Values outside [Min, Max] are clamped to the nearest boundary. |

## Host Monitoring

| Variable | Default | Description |
|---|---|---|
| `DYNATRACE_HOST_MONITORING_OFFLINE` | `false` | Allows managing host monitoring for offline hosts. |
| `DYNATRACE_HOST_MONITORING_WARNINGS` | `false` | Enables warnings for host monitoring issues. |
| `DYNATRACE_HOST_MONITORING_STRICT_UPDATE_RETRIES` | `""` | Sets the number of retries for strict host monitoring updates. |

## Workflows

| Variable | Default | Description |
|---|---|---|
| `DYNATRACE_WORKFLOW_TASKS_USE_TYPE_LIST` | `false` | Uses a type list for workflow task filtering. |

## HCL / Terraform Generation

| Variable | Default | Description |
|---|---|---|
| `DYNATRACE_HEREDOC` | `""` | Configures heredoc usage in generated HCL. |

## Export

| Variable | Default | Description |
|---|---|---|
| `DYNATRACE_TARGET_FOLDER` | `""` | Sets the target folder for export output. |
| `DYNATRACE_CLEAN_TARGET_FOLDER` | `false` | Cleans the target folder before export. |
| `DYNATRACE_PROVIDER_SOURCE` | `""` | Sets the provider source for generated Terraform files. |
| `DYNATRACE_PROVIDER_VERSION` | `""` | Sets the provider version for generated Terraform files. |
| `DYNATRACE_CUSTOM_PROVIDER_LOCATION` | `""` | Sets the custom provider binary location. |
| `DYNATRACE_NO_REFRESH_ON_IMPORT` | `false` | Skips refreshing resources after export. |
| `DYNATRACE_QUICK_INIT` | `false` | Enables quick initialization during export. |
| `DYNATRACE_ULTRA_PARALLEL` | `false` | Enables ultra-parallel export mode. |
| `DYNATRACE_PARALLEL` | `""` | Sets the number of parallel export workers. |
| `DYNATRACE_SHORTER_NAMES` | `false` | Shortens resource names that would exceed 240 characters in generated HCL. |
| `DYNATRACE_ENABLE_EXPORT_DASHBOARD` | `false` | Enables export of dashboard resources. |
| `DYNATRACE_ATOMIC_DEPENDENCIES` | `false` | Exports dependencies atomically. |
| `DYNATRACE_MIGRATION_CACHE_FOLDER` | `""` | Sets the migration cache folder path. |
| `DYNATRACE_FORMAT_HCL_FILES` | `false` | Formats exported HCL files. |
| `DYNATRACE_HCL_NO_FORMAT` | `false` | Disables HCL formatting during export. |
| `DYNATRACE_NAME_REPLACE_DASH` | `false` | Replaces dashes with underscores in resource names. |
| `DYNATRACE_BUILD_ADDRESS_FILES` | `false` | Generates address files during export. |
| `DYNATRACE_EXPORT_IGNORE_RESOURCES` | `""` | Comma-separated list of resource types to exclude from export. |
| `DYNATRACE_IGNORE_CHANGES_REQUIRES_ATTENTION` | `false` | Adds a `lifecycle { ignore_changes = [...] }` block to exported resources containing sensitive fields so Terraform won't overwrite values set manually after export. |
| `DYNATRACE_IMPORT_STATE_PATH` | `""` | Sets the path for importing Terraform state. |
| `DYNATRACE_PREV_STATE_ON` | `false` | Keeps resource identifiers stable between runs so downstream Terraform references don't break. |
| `DYNATRACE_PREV_STATE_PATH_THIS` | `""` | Sets the path to the previous state for the current environment. |
| `DYNATRACE_PREV_STATE_PATH_LINKED` | `""` | Sets the path to the previous state for linked environments. |

## Testing

| Variable | Default | Description |
|---|---|---|
| `DYNATRACE_ENV_URL` | `""` | Sets the Dynatrace environment URL for acceptance tests. |
| `DYNATRACE_API_TOKEN` | `""` | Sets the Dynatrace API token for acceptance tests. |
| `TF_ACC` | `""` | Enables Terraform acceptance tests when set to a non-empty value. |

## Migration

| Variable | Default | Description |
|---|---|---|
| `MIGRATION` | `false` | Enables migration mode. |

## Synthetic Monitors

| Variable | Default | Min | Max | Description |
|---|---|---|---|---|
| `DYNATRACE_CREATE_CONFIRM_SYNTHETIC_MONITORS_V2` | `8` | `1` | `50` | Number of confirmation retries when creating synthetic monitors v2. Falls back to the default if the value is outside [Min, Max]. |

## Web Application

| Variable | Default | Min | Max | Description |
|---|---|---|---|---|
| `DYNATRACE_CREATE_CONFIRM_WEB_APPLICATION` | `280` | `20` | `500` | Number of confirmation retries when creating web applications. Falls back to the default if the value is outside [Min, Max]. |

## Custom Tags

| Variable | Default | Min | Max | Description |
|---|---|---|---|---|
| `DYNATRACE_MAX_CONCURRENT_CUSTOM_TAG_LIST_REQUESTS` | `4` | `1` | `20` | Maximum number of concurrent requests for listing custom tags. Falls back to the default if the value is outside [Min, Max]. |

## Source Environment

These variables support multiple key aliases checked in priority order — the first non-empty value found is used.

| Keys (checked in order) | Description |
|---|---|
| `DYNATRACE_SOURCE_ENV_URL`, `DT_SOURCE_ENV_URL`, `DYNATRACE_SOURCE_ENVIRONMENT_URL`, `DT_SOURCE_ENVIRONMENT_URL` | URL of the source Dynatrace environment for migration. |
| `DYNATRACE_SOURCE_API_TOKEN`, `DT_SOURCE_API_TOKEN` | API token for the source Dynatrace environment. |
| `DT_SOURCE_CLIENT_ID`, `DYNATRACE_SOURCE_CLIENT_ID` | OAuth client ID for the source Dynatrace environment. |
| `DT_SOURCE_ACCOUNT_ID`, `DYNATRACE_SOURCE_ACCOUNT_ID` | Account ID for the source Dynatrace environment. |
| `DT_SOURCE_CLIENT_SECRET`, `DYNATRACE_SOURCE_CLIENT_SECRET` | OAuth client secret for the source Dynatrace environment. |
| `DYNATRACE_SOURCE_PLATFORM_TOKEN`, `DT_SOURCE_PLATFORM_TOKEN` | Platform token for the source Dynatrace environment. |
