---
layout: ""
page_title: "Environment Migration"
description: |-
  The environment migration guide provides information on migrating configuration from one environment into another
---

## Environment Migration

-> This guide assumes that the OneAgent communication is reconfigured via [oneagentctl](https://www.dynatrace.com/support/help/setup-and-configuration/dynatrace-oneagent/oneagent-configuration-via-command-line-interface) for migration. Moving the OneAgent without oneagentctl may generate new entity IDs which could lead to missing configuration.

This guide covers information required to migrate configuration from one environment into another with the [export utility](https://registry.terraform.io/providers/dynatrace-oss/dynatrace/latest/docs/guides/export-v2) of the Terraform provider. The functionality supports all types of migrations including Managed to SaaS, SaaS to SaaS, etc. Please reach out to your Dynatrace account team for any assistance with migration planning.

There are two main approaches for migration (more information below):
* Bulk - Migrate ALL configuration from one environment to a fresh target environment
* Iterative - Migrate configuration by resource group - safer approach for a target environment with existing configuration

A recommended approach is to pre-apply configuration prior to the OneAgent migration and then re-apply after to deploy any configuration that previously returned an error due to an entity dependency. It can be expected that most configuration will be successfully applied during the pre-apply phase. **Please keep exported working directories of configuration that are applied in case any configuration needs to be re-applied or destroyed.**

An Excel spreadsheet is available to assist with the migration process, specifically helpful for an iterative approach. Note that the spreadsheet has a custom function to track a static time/date for when the checkbox for column A "Complete" is selected. Leaving the document untrusted will simply make column C "Date Completed" unusable, everything else in the document will function and view properly.

For the iterative approach, it is important to push configuration in order of dependency. Please reference column D "Migration: Iterative Order" and push all configuration in order of the number present for each resource. Eg. Push all configuration with iterative order "1", then all with "2", etc. 

Link: [Terraform Migration Helper](https://github.com/dynatrace-oss/terraform-provider-dynatrace/blob/main/documentation/Terraform%20-%20Migration%20Helper.xlsm)

### Command Line Syntax
There are two options with invoking the export utility and subsequent Terraform execution. Option 1 is particularly helpful with the iterative method where you may be running multiple exports and Terraform executions within the same command line shell.

Option 1:
* Set the environment variables `DYNATRACE_SOURCE_ENV_URL` and `DYNATRACE_SOURCE_API_TOKEN` as the URL and API Token of your ***source*** Dynatrace environment.
* Set the environment variables `DYNATRACE_ENV_URL` and `DYNATRACE_API_TOKEN` as the URL and API Token of your ***target*** Dynatrace environment.

Option 2:
* Set the environment variables `DYNATRACE_ENV_URL` and `DYNATRACE_API_TOKEN` as the URL and API Token - same variables are used for your ***source and target*** Dynatrace environment. Change the values between the invocation of the export utility and Terraform execution.

Optionally the environment variable `DYNATRACE_TARGET_FOLDER` allows you to set a target directory of the export utility. If it's not set, the output folder `./configuration` is assumed.

### Bulk Migration
In order for the best results of a bulk migration, it is highly recommended that the target environment does not have any existing configuration.

For bulk migration, utilize the `-migrate` flag. This will create all of the necessary files with resource dependencies and hardcoded entity IDs.

Windows: `terraform-provider-dynatrace.exe -export -migrate`

Linux: `terraform-provider-dynatrace -export -migrate`

The following resources are excluded by default, please use the iterative method to export/apply after running the bulk migration.
* Excluded by default due to longer execution: `dynatrace_json_dashboard`, `dynatrace_custom_tags`, `dynatrace_custom_device`
* Application Security: `dynatrace_appsec_notification`, `dynatrace_vulnerability_alerting`, `dynatrace_vulnerability_settings`, `dynatrace_vulnerability_third_party`, `dynatrace_vulnerability_code`, `dynatrace_attack_alerting`, `dynatrace_attack_settings`, `dynatrace_attack_rules`, `dynatrace_attack_allowlist`
* AutomationEngine: `dynatrace_automation_workflow`, `dynatrace_automation_business_calendar`, `dynatrace_automation_scheduling_rule`
* Account Management: `dynatrace_iam_group`(SaaS), `dynatrace_iam_permission`(SaaS), `dynatrace_iam_policy_bindings`(SaaS), `dynatrace_iam_policy`(SaaS), `dynatrace_iam_user`(SaaS), `dynatrace_mgmz_permission`(Managed), `dynatrace_policy_bindings`(Managed), `dynatrace_policy`(Managed), `dynatrace_user_group`(Managed), `dynatrace_user`(Managed)

### Iterative Migration
The iterative approach is useful in scenarios where you would like to migrate configuration by resource group or to an environment with existing configuration. 

Most resources do not allow configuration with the same name, but there is a small number of resources where duplicates could be created. Please set an optional environment variable below depending on the preferred behavior.
* `DYNATRACE_DUPLICATE_REJECT=ALL` - Duplicates will not overwrite existing configuration
* `DYNATRACE_DUPLICATE_HIJACK=ALL` - Duplicates will overwrite existing configuration

A common approach with the configuration migration is to do an apply prior to the OneAgent migration and then a re-apply after to deploy any configuration that depend on entities existing. An underlying dependency which has entity verification could have cascading effects where a data source could be driven to `null` and the plan/apply prevents execution. 

Example: You may have a request attribute which has an entity filter, a calculated service metric that depends on this request attribute, and a dashboard that is using the calculated service metric. 

In order to have a successful execution, the environment variable `MIGRATION=true` can be set which will replace `null` data source references with a name prefixed with `TFMIGRATIONID-`. This feature is currently available for calculated service metric and SLO data sources. When enabling this feature, a log file will be produced which will contain all resources that are incomplete due to missing underlying dependencies. Reapplying the configuration once the underlying resource dependency exists will automatically repair any resources.

For an iterative migration, utilize the `-migrate -datasources` flags. This will create all of the necessary files with dependencies via data sources and hardcoded entity IDs.

Windows: `terraform-provider-dynatrace.exe -export -migrate -datasources <resourcename>`

Linux: `terraform-provider-dynatrace -export -migrate -datasources <resourcename>`

### Additional Information
* It is expected that there are resources that may fail during a migration, two types of errors in particular are noted below. If you experience any errors beyond what is specified below, please create a [GitHub issue](https://github.com/dynatrace-oss/terraform-provider-dynatrace/issues).
  1. "Grandfathered" configuration - certain API endpoints do not accept configuration that is no longer supported. eg. Calculated service metrics require a management zone or at least one condition marked with [Service property]. 
  2. Missing entities - certain API endpoints validate whether entity IDs in the configuration exists. If you encounter this issue, simply re-apply once the entities exist in the target environment. 