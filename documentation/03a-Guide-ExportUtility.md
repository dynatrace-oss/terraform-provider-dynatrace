# Export Utility

In addition to the out-of-the-box functionality of Terraform, the provider has the ability to be executed as a standalone executable to export existing configuration from a Dynatrace environment. This functionality provides an alternative to manually creating Terraform configuration and provides an easy way to create templates based off of existing configuration.

## Prerequisites
* Terraform CLI with the Dynatrace provider installed (see [Getting Started with Terraform and the Dynatrace Provider](LINK)) and available under PATH.
* A Dynatrace environment and access to create environment tokens.
* A Dynatrace token with at least the following permissions:
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
1. Open a terminal window and set the tenant URL and API token environment variables. This will be the tenant we want to retrieve configuration. On Windows, SaaS:
```
set DYNATRACE_ENV_URL=https://########.live.dynatrace.com
set DYNATRACE_API_TOKEN=dt0c01.########.########
```
Managed:
```
set DYNATRACE_ENV_URL=https://<dynatrace-host>/e/########
set DYNATRACE_API_TOKEN=dt0c01.########.########
```
Optionally, the environment variable `DYNATRACE_TARGET_FOLDER` can be set to specify an output directory. If the variable is not set, the default `./configuration` will be used.

2. In the terminal window, locate and navigate to the Terraform Dynatrace Provider executable retrieved from [Getting Started with Terraform and the Dynatrace Provider](LINK). In most cases, the executable should be available in the working directory under `/.terraform/providers/registry.terraform.io/dynatrace-oss/dynatrace/{provider_version}/{os_version}/terraform-provider-dynatrace_x.y.z`

3. The executable can be directly invoked with the following options. Please refer to the Usage Examples section for examples.

Windows: `terraform-provider-dynatrace.exe -export [-ref] [-migrate] [-id] [-flat] [-exclude] [<resourcename>[=<id>]]`

Linux: `./terraform-provider-dynatrace -export [-ref] [-migrate] [-import-state-v2] [-id] [-flat] [-exclude] [<resourcename>[=<id>]]`

### Options
* `-ref` Enable resources with data sources and dependencies
* `-migrate` Enable resources with dependencies, no data sources. More information available in the Migration guide (LINK).
* `-import-state-v2` Automatically initializes the terraform modules and imports downloaded resources into the state
* `-id` Enable commented id output in resource files
* `-flat` All downloaded resources end up directly within the target folder - no module structure will be created
* `-exclude` Exclude specified resource(s) from export

**NOTE:** Dashboards are currently excluded from the default export due to the amount of results this may return. Please specify the resource `dynatrace_json_dashboard` directly in the command line arguments to retrieve dashboards.

### Usage Examples
* `./terraform-provider-dynatrace -export` downloads all available configuration settings without data sources / dependency references
* `./terraform-provider-dynatrace -export -ref -id` downloads all available configuration settings with data sources / dependency references and adds commented ids in resource output
* `./terraform-provider-dynatrace -export -ref dynatrace_json_dashboard dynatrace_web_application` downloads all available dashboards, web applications and resource dependencies with references
* `./terraform-provider-dynatrace -export -ref dynatrace_alerting=4f5942d4-3450-40a8-818f-c5faeb3563d0 dynatrace_alerting=9c4b75f1-9a64-4b44-a8e4-149154fd5325` downloads the alerting profiles with the ids `4f5942d4-3450-40a8-818f-c5faeb3563d0` and `9c4b75f1-9a64-4b44-a8e4-149154fd5325`, includes all resource dependencies with references
* `./terraform-provider-dynatrace -export -ref dynatrace_calculated_service_metric dynatrace_alerting=4f5942d4-3450-40a8-818f-c5faeb3563d0` downloads all available calculated service metrics and also the alerting profile with the id `4f5942d4-3450-40a8-818f-c5faeb3563d0`, includes all resource dependencies with references
* `./terraform-provider-dynatrace -export -import-state-v2` downloads all available configuration settings and imports resources into the state
* `./terraform-provider-dynatrace -export -import-state-v2 dynatrace_web_application` downloads all web applications and imports resources into the state
* `./terraform-provider-dynatrace -export -ref -exclude dynatrace_calculated_service_metric dynatrace_alerting` download all available configuration settings except `dynatrace_calculated_service_metric` and `dynatrace_alerting`, includes all resource dependencies with references

### Additional Information
* There may be instances where the exported configuration is deprecated and/or is unable to be used without modification. In these instances, the files will be moved into `.flawed` of the output folder and the reasoning will be provided as a comment at the top of the the resource file. 
    -  E.g. A dashboard with no tiles can be created and can be retrieved via the export, but the subsequent `terraform apply` would fail without any tiles. 
* There are instances where the returned configuration does not contain all of the required information to run an `terraform apply` due to sensitive data or  instances where the files require additional attention. The files that apply to this scenario will be automatically moved to `.requires_attention`, the reasoning will be provided as a comment at the top of the the resource file. 
    -  E.g. `dynatrace_credentials` confidential strings are not available via the API.