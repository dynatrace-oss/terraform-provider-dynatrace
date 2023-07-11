# Migration Guide

-> This guide assumes that the OneAgent communication is reconfigured via [oneagentctl](https://www.dynatrace.com/support/help/setup-and-configuration/dynatrace-oneagent/oneagent-configuration-via-command-line-interface) for migration. Moving the OneAgent without oneagentctl may generate new entity IDs which could lead to missing configuration.

Terraform can be used to migrate configuration between Dynatrace environments. The functionality supports all types of migrations including Managed to SaaS, SaaS to SaaS, etc. Please reach out to your Dynatrace account team for any assistance with migration planning.

There are two main approaches for migration (more information below):
* Bulk - Migrate ALL configuration from one environment to a fresh target environment
* Iterative - Migrate configuration by resource group - safer approach for a target environment with existing configuration

An Excel spreadsheet is available to assist with the migration process, specifically helpful for an iterative approach. Note that the spreadsheet has a custom function to track a static time/date for when the checkbox for column A "Complete" is selected. Leaving the document untrusted will simply make column C "Date Completed" unusable, everything else in the document will function and view properly.

For the iterative approach, it is important to push configuration in order of dependency. Please reference column D "Migration: Iterative Order" and push all configuration in order of the number present for each resource. Eg. Push all configuration with iterative order "1", then all with "2", etc. 

Download: [Terraform Migration Helper](https://github.com/dynatrace-oss/terraform-provider-dynatrace/blob/main/documentation/Terraform%20-%20Migration%20Helper.xlsm)

## Prerequisites
* Terraform CLI with the Dynatrace provider installed (see [Getting Started with Terraform and the Dynatrace Provider](LINK)) and available under PATH.
* A Dynatrace ***source*** environment and access to create environment tokens.
* A Dynatrace ***target*** environment and access to create environment tokens.
* Source and target Dynatrace tokens with at least the following permissions:
 * **Read configuration** (`ReadConfig`)
     * Required for reading API v1 configuration.
 * **Write configuration** (`WriteConfig`)
     * Required for creating API v1 configuration.
 * **Read settings** (`settings.read`)
     * Required for reading Settings 2.0 configuration.
 * **Write settings** (`settings.write`)
     * Required for modifying Settings 2.0 configuration.
 * **Create and read synthetic monitors, locations, and nodes** (`ExternalSyntheticIntegration`)
     * Required for reading and creating synthetic configuration.
 * **Capture request data** (`CaptureRequestData`)
     * Required for configuring request attributes.
 * **Read credential vault entries** (`credentialVault.read`)
     * Required for reading credentials.
 * **Write credential vault entries** (`credentialVault.write`)
     * Required for writing credentials.
 * **Read network zones** (`networkZones.read`)
     * Required for reading network zones.
 * **Write network zones** (`networkZones.write`)
     * Required for writing network zones.

To learn how to create tokens, see [Create an API token](https://www.dynatrace.com/support/help/manage/access-control/access-tokens#create-api-token).

## Guide

### Set Environment Variables
There are two options with invoking the export utility and subsequent Terraform execution. Option 1 is particularly helpful with the iterative method where you may be running multiple exports and Terraform executions within the same command line shell.

Option 1:
* Set the environment variables `DYNATRACE_SOURCE_ENV_URL` and `DYNATRACE_SOURCE_API_TOKEN` as the URL and API Token of your ***source*** Dynatrace environment.
* Set the environment variables `DYNATRACE_ENV_URL` and `DYNATRACE_API_TOKEN` as the URL and API Token of your ***target*** Dynatrace environment.

Example (Managed to SaaS on Windows):
```
set DYNATRACE_SOURCE_ENV_URL=https://<dynatrace-host>/e/########
set DYNATRACE_SOURCE_API_TOKEN=dt0c01.########.########
set DYNATRACE_ENV_URL=https://########.live.dynatrace.com
set DYNATRACE_API_TOKEN=dt0c01.########.########
```

Option 2:
* Set the environment variables `DYNATRACE_ENV_URL` and `DYNATRACE_API_TOKEN` as the URL and API Token - same variables are used for your ***source and target*** Dynatrace environment. Change the values between the invocation of the export utility and Terraform execution. 

Optionally, the environment variable `DYNATRACE_TARGET_FOLDER` can be set to specify an output directory. If the variable is not set, the default `./configuration` will be used.

### Bulk Migration
In order for the best results of a bulk migration, it is highly recommended that the target environment does not have any existing configuration.

For bulk migration, utilize the `-migrate` flag. This will create all of the necessary files with resource dependencies and hardcoded entity IDs.

Windows: `terraform-provider-dynatrace.exe -export -migrate`

Linux: `terraform-provider-dynatrace -export -migrate`

-> The `dynatrace_json_dashboard`, `dynatrace_aws_service`, `dynatrace_azure_service`, and account management resources are excluded by default, please use the iterative method to export/apply after running the bulk migration.

### Iterative Migration
The iterative approach is useful in scenarios where you would like to migrate configuration by resource group or to an environment with existing configuration. 

Most resources do not allow configuration with the same name, but there is a small number of resources where duplicates could be created. Please set an optional environment variable below depending on the preferred behavior.
* `DYNATRACE_DUPLICATE_REJECT=ALL` - Duplicates will not overwrite existing configuration
* `DYNATRACE_DUPLICATE_HIJACK=ALL` - Duplicates will overwrite existing configuration

For an iterative migration, utilize the `-migrate -datasources` flags. This will create all of the necessary files with dependencies via data sources and hardcoded entity IDs.

Windows: `terraform-provider-dynatrace.exe -export -migrate -datasources <resourcename>`

Linux: `terraform-provider-dynatrace -export -migrate -datasources <resourcename>`

### Additional Information
* It is expected that there are resources that may fail during a migration, two types of errors in particular are noted below. If you experience any errors beyond what is specified below, please create a [GitHub issue](https://github.com/dynatrace-oss/terraform-provider-dynatrace/issues).
  1. "Grandfathered" configuration - certain API endpoints do not accept configuration that is no longer supported. eg. Calculated service metrics require a management zone or at least one condition marked with [Service property]. 
  2. Missing entities - certain API endpoints validate whether entity IDs in the configuration exists. If you encounter this issue, simply re-apply once the entities exist in the target environment. 