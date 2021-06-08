# Dynatrace Terraform Provider
## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.13.x+

## Using the provider

If you want to run Terraform with the dynatrace provider plugin on your system, add the plug-in provider to the Terraform configuration file.

    ```
    terraform {
        required_version = "~> 0.13.0"
        required_providers {
            dynatrace = {
                version = "1.2.2"
                source = "dynatrace-oss/dynatrace"
            }
        }
    }   
    ```

In order to configure the provider, add a code section like this into your Terraform configuration file

    ```
       provider "dynatrace" {
           dt_env_url   = "https://#######.live.dynatrace.com"
           dt_api_token = "##########################################"
       }    
    ```
where `dt_env_url` represents the URL of your Dynatrace Environment and `dt_api_token` needs to be an API Token with the permissions `Read configuration` and `Capture request data`.

## Currently supported configuration settings
* Dashboards
* Management Zones
* Custom Services
* Request Attributes
* Alerting Profiles
* Problem Notifiations
* Auto Tag Configuration
* Kubernetes Credentials
* AWS Credentials
* Azure Credentials
* Maintenance Windows

## Exporting existing configuration from a Dynatrace Environment
In addition to acting as a Terraform Provider Plugin the executable `terraform-provider-dynatrace` (`terraform-provider-dynatrace.exe` on Windows) can also get directly invoked.
The utility then reaches out to the Dynatrace Environment specified by the command line arguments and fetches all currently supported configuration items. These results will then automatically get transformed into HCL (the configuration language to be used for `.tf` files) and places each configuration item into its own `.tf` file).
### Command Line Syntax
Invoking the download functionality requires
* The URL of your Dynatrace Environment
* An API Token with the permissions `Read configuration` and `Capture request data`
* If the optional argument for the output folder is not specified `./configuration` is assumed
#### Windows
`terraform-provider-dynatrace.exe download https://<environment-id>.live.dynatrace.com <api-token> [<output-folder>]`
#### Linux
`./terraform-provider-dynatrace download https://<environment-id>.live.dynatrace.com <api-token> [<output-folder>]`
