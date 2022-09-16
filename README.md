# Dynatrace Terraform Provider
## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.13.x+

## Using the provider

Please check the documentation within the [Terraform Registry](https://registry.terraform.io/providers/dynatrace-oss/dynatrace/latest/docs) for how to use and configure this Provider as well as for currrently supported resources and data sources.

## Exporting existing configuration from a Dynatrace Environment
In addition to acting as a Terraform Provider Plugin the executable `terraform-provider-dynatrace` (`terraform-provider-dynatrace.exe` on Windows) can also get directly invoked.
The utility then reaches out to the Dynatrace Environment specified by the command line arguments and fetches all currently supported configuration items. These results will then automatically get transformed into HCL (the configuration language to be used for `.tf` files) and places each configuration item into its own `.tf` file).
### Command Line Syntax
Invoking the export functionality requires
* The environment varibale `DYNATRACE_ENV_URL` as the URL of your Dynatrace Environment
* The environment variable `DYNATRACE_API_TOKEN` as the API Token with the following permissions:
    -  `Read configuration`
    -  `Capture request data`
    -  `Read SLO`
    -  `Read settings`
    -  `Read synthetic monitors, locations, and nodes`
* Optinonally the environment variable `DYNATRACE_TARGET_FOLDER`. If it's not set, the output folder `./configuration` is assumed
#### Windows
`terraform-provider-dynatrace.exe export *[<resourcename>[=<id>]]`
#### Linux
`./terraform-provider-dynatrace export *[<resourcename>[=<id>]]`
#### Usage Examples
* `./terraform-provider-dynatrace export` downloads all available configuration settings
* `./terraform-provider-dynatrace export dynatrace_dashboard` downloads all available dashboards
* `./terraform-provider-dynatrace export dynatrace_dashboard dynatrace_slo` downloads all available dashboards and all available SLOs
* `./terraform-provider-dynatrace export dynatrace_dashboard=4f5942d4-3450-40a8-818f-c5faeb3563d0` downloads only the dashboard with the id `4f5942d4-3450-40a8-818f-c5faeb3563d0`
* `./terraform-provider-dynatrace export dynatrace_dashboard=4f5942d4-3450-40a8-818f-c5faeb3563d0 dynatrace_dashboard=9c4b75f1-9a64-4b44-a8e4-149154fd5325` downloads only the dashboards with the ids `4f5942d4-3450-40a8-818f-c5faeb3563d0` and `9c4b75f1-9a64-4b44-a8e4-149154fd5325`
* `./terraform-provider-dynatrace export dynatrace_slo dynatrace_dashboard=4f5942d4-3450-40a8-818f-c5faeb3563d0 dynatrace_dashboard=9c4b75f1-9a64-4b44-a8e4-149154fd5325` downloads all available SLOs and only the dashboards with the ids `4f5942d4-3450-40a8-818f-c5faeb3563d0` and `9c4b75f1-9a64-4b44-a8e4-149154fd5
