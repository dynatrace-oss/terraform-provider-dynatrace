# Dynatrace Terraform Provider
## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.13.x+

## Using the provider

Please check the documentation within the [Terraform Registry](https://registry.terraform.io/providers/dynatrace-oss/dynatrace/latest/docs) for how to use and configure this Provider as well as for currrently supported resources and data sources.

Additional information can be found under [Dynatrace Configuration as Code via Terraform](https://www.dynatrace.com/support/help/manage/configuration-as-code/terraform).

## Exporting existing configuration from a Dynatrace Environment

In addition to acting as a Terraform Provider Plugin the executable `terraform-provider-dynatrace` (`terraform-provider-dynatrace.exe` on Windows) can also get directly invoked.
The utility then reaches out to the Dynatrace Environment specified by the command line arguments and fetches all currently supported configuration items. These results will then automatically get transformed into HCL (the configuration language to be used for `.tf` files) and places each configuration item into its own `.tf` file).

Please check out the documentation within the [Terraform Registry](https://registry.terraform.io/providers/dynatrace-oss/dynatrace/latest/docs#exporting-existing-configuration-from-a-dynatrace-environment) for detailed information about how to use that functionality.
