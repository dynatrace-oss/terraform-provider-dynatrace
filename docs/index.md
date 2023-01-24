---
layout: ""
page_title: "Provider: Dynatrace"
description: |-
  The Dynatrace Terraform Provider provides resources to interact with the Dynatrace REST API.
---

# Dynatrace Terraform Provider

The Dynatrace Terraform Provider is used to interact with the resources supported by the Dynatrace REST API. The provider needs to be configured with the proper credentials before it can be used.

Use the links to the left to learn about the available resources and data sources.

## Example

# Terraform 0.13+ uses the Terraform Registry:
```
terraform {
    required_providers {
        dynatrace = {
            version = "1.16.0"
            source = "dynatrace-oss/dynatrace"
        }
    }
} 
```
## Configure the Dynatrace provider
```
provider "dynatrace" {
    dt_env_url    = "https://########.live.dynatrace.com"
    dt_api_token  = "################"
}
```
where `dt_env_url` represents the URL of your Dynatrace Environment and `dt_api_token` needs to be an API Token with the following permissions:
* `Read configuration`
* `Capture request data`
* `Read SLO`
* `Read settings`
* `Read synthetic monitors, locations, and nodes`

## Exporting existing configuration from a Dynatrace Environment
In addition to acting as a Terraform Provider Plugin, the executable `terraform-provider-dynatrace` (`terraform-provider-dynatrace.exe` on Windows) can be directly invoked.
The utility reaches out to the Dynatrace Environment specified by the command line arguments and fetches all currently supported configuration items. These results will then automatically get transformed into HCL (the configuration language to be used for `.tf` files) and places each configuration item into its own `.tf` file.

With release 1.15.0, a new and improved version of the export utility is now available with various improvements (outlined below).
* Data source references automatically generated in resources
* Data source blocks created for dependencies in data_source.tf file(s)
* Creation of all resources required based off of dependencies
* Resources organized in modules with preconfigured main.tf based on output
* Pre-populated providers.tf file in the parent folder and all modules
* Option to include ID as a commented line in generated resources
* Option to include/exclude specific resource IDs from export
* Ability to convert from HCL to JSON formatted files
 
Please provide any feedback of the export utility via [GitHub Issues](https://github.com/dynatrace-oss/terraform-provider-dynatrace/issues). Details on the execution of the export utility is provided below.

### Command Line Syntax
Invoking the export functionality requires
* The environment variable `DYNATRACE_ENV_URL` as the URL of your Dynatrace environment
* The environment variable `DYNATRACE_API_TOKEN` as the API Token of your Dynatrace environment
* Optionally the environment variable `DYNATRACE_TARGET_FOLDER`. If it's not set, the output folder `./configuration` is assumed

## Export Utility
Windows: `terraform-provider-dynatrace.exe -export [-v] [-ref] [-id] [-migrate] [-exclude] [<resourcename>[=<id>]]`

Linux: `./terraform-provider-dynatrace -export [-v] [-ref] [-id] [-migrate] [-exclude] [<resourcename>[=<id>]]`
### Options
* `-v` Enable verbose logging
* `-ref` Enable resources with data sources and dependencies
* `-id` Enable commented id output in resource files
* `-migrate` Enable output specific to environment migration
    -  Removes node IDs from private synthetic locations
* `-exclude` Exclude specified resource(s) from export

**NOTE:** Dashboards (because there could be thousands of them) are currently excluded from the export unless the resource is directly specified in the command line arguments.

### Usage Examples
* `./terraform-provider-dynatrace -export` downloads all available configuration settings without data sources and dependency references (export similar to previous version)
* `./terraform-provider-dynatrace -export -ref -id` downloads all available configuration settings with data sources / dependency references and adds commented ids in resource output
* `./terraform-provider-dynatrace -export -ref dynatrace_dashboard dynatrace_web_application` downloads all available dashboards, web applications and resource dependencies with references
* `./terraform-provider-dynatrace -export -ref dynatrace_alerting=4f5942d4-3450-40a8-818f-c5faeb3563d0 dynatrace_alerting=9c4b75f1-9a64-4b44-a8e4-149154fd5325` downloads the alerting profiles with the ids `4f5942d4-3450-40a8-818f-c5faeb3563d0` and `9c4b75f1-9a64-4b44-a8e4-149154fd5325`, includes all resource dependencies with references
* `./terraform-provider-dynatrace -export -ref dynatrace_calculated_service_metric dynatrace_alerting=4f5942d4-3450-40a8-818f-c5faeb3563d0` downloads all available calculated service metrics and also the alerting profile with the id `4f5942d4-3450-40a8-818f-c5faeb3563d0`, includes all resource dependencies with references
* `./terraform-provider-dynatrace -export -ref -exclude dynatrace_calculated_service_metric dynatrace_alerting` download all available configuration settings except `dynatrace_calculated_service_metric` and `dynatrace_alerting`, includes all resource dependencies with references

### Additional Information
* There may be instances where the exported configuration is deprecated and/or is unable to be used for a create/update. In these instances, the files will be moved into `.flawed` of the output folder and the explanation will be available as a commented output in the resource file. 
    -  E.g. A dashboard with no tiles can be created and can be retrieved via the export, but the subsequent `terraform apply` would fail without any tiles. 
* There are instances where the returned configuration does not contain all of the required information to run an `terraform apply` due to sensitive data or  instances where the files require additional attention. The files that apply to this scenario will be automatically moved to `.requires_attention`, the explanation will be available as a commented output in the resource file.
    -  E.g. `dynatrace_credentials` confidential strings are not available via the API.

### Known Issues
* Due to the complexity of dashboards, there may be cases where the `terraform apply` fails after an export. Dashboard validation will be improved in a future release.

## Export Utility (Legacy method)
Windows: `terraform-provider-dynatrace.exe exportv1 *[<resourcename>[=<id>]]`

Linux: `./terraform-provider-dynatrace exportv1 *[<resourcename>[=<id>]]`

### Usage Examples
* `./terraform-provider-dynatrace exportv1` downloads all available configuration settings
* `./terraform-provider-dynatrace exportv1 dynatrace_dashboard` downloads all available dashboards
* `./terraform-provider-dynatrace exportv1 dynatrace_dashboard dynatrace_slo` downloads all available dashboards and all available SLOs
* `./terraform-provider-dynatrace exportv1 dynatrace_dashboard=4f5942d4-3450-40a8-818f-c5faeb3563d0` downloads only the dashboard with the id `4f5942d4-3450-40a8-818f-c5faeb3563d0`
* `./terraform-provider-dynatrace exportv1 dynatrace_dashboard=4f5942d4-3450-40a8-818f-c5faeb3563d0 dynatrace_dashboard=9c4b75f1-9a64-4b44-a8e4-149154fd5325` downloads only the dashboards with the ids `4f5942d4-3450-40a8-818f-c5faeb3563d0` and `9c4b75f1-9a64-4b44-a8e4-149154fd5325`
* `./terraform-provider-dynatrace exportv1 dynatrace_slo dynatrace_dashboard=4f5942d4-3450-40a8-818f-c5faeb3563d0 dynatrace_dashboard=9c4b75f1-9a64-4b44-a8e4-149154fd5325` downloads all available SLOs and only the dashboards with the ids `4f5942d4-3450-40a8-818f-c5faeb3563d0` and `9c4b75f1-9a64-4b44-a8e4-149154fd5

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `dt_api_token` (String, Sensitive)
- `dt_cluster_api_token` (String, Sensitive)
- `dt_cluster_url` (String, Sensitive)
- `dt_env_url` (String)

- `iam_account_id` (String, Sensitive)
- `iam_client_id` (String, Sensitive)
- `iam_client_secret` (String, Sensitive)

