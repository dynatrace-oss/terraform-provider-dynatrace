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

## IAM

| Variable | Default | Description |
|---|---|---|
| `DYNATRACE_DISABLE_IAM_RATE_LIMITER` | `false` | Disables the IAM rate limiter. |

## Settings 2.0

| Variable | Default | Description |
|---|---|---|
| `DT_NO_REPAIR_INPUT` | `false` | Disables automatic repair of invalid settings input. |
| `DYNATRACE_DISABLE_ORDERING_SUPPORT` | `false` | Disables ordering support for settings resources. |

## Cache

| Variable | Default | Description |
|---|---|---|
| `CACHE_OFFLINE_MODE` | `false` | Enables offline mode using cached data. |
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

## Documents

| Variable | Default | Description |
|---|---|---|
| `DT_DOCUMENTS_IGNORE_UNEXPECTED_EOF` | `false` | Ignores unexpected EOF errors when managing documents. |

## Entity

| Variable | Default | Description |
|---|---|---|
| `DYNATRACE_DISABLE_ENTITY_CACHE` | `false` | Disables the entity cache. |

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

## Host Monitoring

| Variable | Default | Description |
|---|---|---|
| `DYNATRACE_HOST_MONITORING_OFFLINE` | `false` | Allows managing host monitoring for offline hosts. |
| `DYNATRACE_HOST_MONITORING_WARNINGS` | `false` | Enables warnings for host monitoring issues. |

## Workflows

| Variable | Default | Description |
|---|---|---|
| `DYNATRACE_WORKFLOW_TASKS_USE_TYPE_LIST` | `false` | Uses a type list for workflow task filtering. |

## Export

| Variable | Default | Description |
|---|---|---|
| `DYNATRACE_CLEAN_TARGET_FOLDER` | `false` | Cleans the target folder before export. |
| `DYNATRACE_NO_REFRESH_ON_IMPORT` | `false` | Skips refreshing resources after export. |
| `DYNATRACE_QUICK_INIT` | `false` | Enables quick initialization during export. |
| `DYNATRACE_ULTRA_PARALLEL` | `false` | Enables ultra-parallel export mode. |
| `DYNATRACE_SHORTER_NAMES` | `false` | Shortens resource names that would exceed 240 characters in generated HCL. |
| `DYNATRACE_ENABLE_EXPORT_DASHBOARD` | `false` | Enables export of dashboard resources. |
| `DYNATRACE_ATOMIC_DEPENDENCIES` | `false` | Exports dependencies atomically. |
| `DYNATRACE_FORMAT_HCL_FILES` | `false` | Formats exported HCL files. |
| `DYNATRACE_HCL_NO_FORMAT` | `false` | Disables HCL formatting during export. |
| `DYNATRACE_NAME_REPLACE_DASH` | `false` | Replaces dashes with underscores in resource names. |
| `DYNATRACE_BUILD_ADDRESS_FILES` | `false` | Generates address files during export. |
| `DYNATRACE_IGNORE_CHANGES_REQUIRES_ATTENTION` | `false` | Adds a `lifecycle { ignore_changes = [...] }` block to exported resources containing sensitive fields so Terraform won't overwrite values set manually after export. |
| `DYNATRACE_PREV_STATE_ON` | `false` | Keeps resource identifiers stable between runs so downstream Terraform references don't break. |

## Migration

| Variable | Default | Description |
|---|---|---|
| `MIGRATION` | `false` | Enables migration mode. |
