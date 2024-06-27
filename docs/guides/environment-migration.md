---
layout: ""
page_title: "Environment Migration"
description: |-
  Export existing Dynatrace configurations using Dynatrace Configuration as Code via Terraform.
---

# Environment Migration

This guide covers both bulk and iterative methods of migrating configurations between Dynatrace environments with Dynatrace Configuration as Code via Terraform.

-> Ensure OneAgent communication is reconfigured via `oneagentctl` for migration. Direct OneAgent migration might generate new entity IDs, potentially causing configuration loss.
**Improper migration can cause significant environment issues**.

## Prerequisites

* [Terraform CLI with the Dynatrace provider installed](https://docs.dynatrace.com/docs/manage/configuration-as-code/terraform/terraform-cli) and available under `PATH`.
* [Access token](https://docs.dynatrace.com/docs/manage/access-control/access-tokens) with at least the following permissions:

  * **Read settings** (`settings.read`)
  * **Write settings** (`settings.write`)
  * **Read configuration** (`ReadConfig`)
  * **Write configuration** (`WriteConfig`)
  * **Create and read synthetic monitors, locations, and nodes** (`ExternalSyntheticIntegration`)
  * **Capture request data** (`CaptureRequestData`)
  * **Read credential vault entries** (`credentialVault.read`)
  * **Write credential vault entries** (`credentialVault.write`)
  * **Read network zones** (`networkZones.read`)
  * **Write network zones** (`networkZones.write`)

-> Certain resources require an OAuth client for authentication (eg. automation, document, account management APIs), please refer to the resource specific pages for additional information.

## Set environment variables

There are two ways to invoke the export utility and subsequent Terraform execution.

### Source and target environment variables set separately
Setting the source and target environment variables separately is beneficial with the iterative method, where you might run multiple exports and Terraform executions within the same command line shell.

1. Set environment variables `DYNATRACE_SOURCE_ENV_URL` and `DYNATRACE_SOURCE_API_TOKEN` as the URL and API token of your source Dynatrace environment.
2. Set environment variables `DYNATRACE_ENV_URL` and `DYNATRACE_API_TOKEN` as the URL and API token of your target Dynatrace environment.

For example (Managed to SaaS on Windows):

```
set DYNATRACE_SOURCE_ENV_URL=https://<dynatrace-host>/e/########
set DYNATRACE_SOURCE_API_TOKEN=dt0c01.########.########
set DYNATRACE_ENV_URL=https://########.live.dynatrace.com
set DYNATRACE_API_TOKEN=dt0c01.########.########
```

### Same variables for both environments
With this method, you use the same variables for both environments.

Set the environment variables `DYNATRACE_ENV_URL` and `DYNATRACE_API_TOKEN` as the URL and API Token---same variables are used for your source and target Dynatrace environment. Change the values between the invocation of the export utility and Terraform execution.

Optionally, the environment variable `DYNATRACE_TARGET_FOLDER` can be set to specify an output directory. If the variable is not set, the default `./configuration` will be used.

## Migration methods

There are two main approaches for migration:

* [Bulk](#bulk-migration)---Transfer all configurations from a source to a fresh target environment.
* [Iterative](#iterative-migration)---Transfer configurations by resource group. Ideal for target environments with existing configurations.

  For iterative migration, use the [Terraform Migration Helper](https://dt-url.net/m5i37ar) spreadsheet. It aids in tracking migration order and completion.

### Bulk migration

Ensure the target environment doesn't have any existing configurations.

Use the `-migrate` flag. This will create all of the necessary files with resource dependencies and hardcoded entity IDs.

**Windows**: `terraform-provider-dynatrace.exe -export -migrate`

**Linux**: `terraform-provider-dynatrace -export -migrate`

-> View the resources that are excluded from the default export by running the provider with the `-list-exclusions` flag. To export any of the excluded resources, please refer to the [Usage Examples](https://docs.dynatrace.com/docs/manage/configuration-as-code/terraform/guides/export-utility#usage-examples).

### Iterative migration

The iterative approach is useful for migrating specific resource groups or to environments with existing configurations.

Set environment variables to handle duplicates:

* `DYNATRACE_DUPLICATE_REJECT=ALL`: Prevents overwriting existing configurations.
* `DYNATRACE_DUPLICATE_HIJACK=ALL`: Allows overwriting existing configurations.

Use the `-migrate -datasources` flag. This will create all of the necessary files with dependencies via data sources and hardcoded entity IDs.

**Windows**: `terraform-provider-dynatrace.exe -export -migrate -datasources <resourcename>`

**Linux**: `terraform-provider-dynatrace -export -migrate -datasources <resourcename>`

### Troubleshooting

Some resources might fail during migration. Common issues include:

* Some API endpoints might reject outdated configurations. For instance, calculated service metrics now require a management zone or a condition marked with service property.
* Some API endpoints validate whether entity IDs in the configuration exist. If missing in the target environment, re-apply once the entities are present.

If you experience any other errors, please create a [GitHub issue](https://dt-url.net/4bg37q8).