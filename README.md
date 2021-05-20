# Dynatrace Terraform Provider
[![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Terraform Website: [https://www.terraform.io](https://www.terraform.io)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)
## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12+

## Using the provider

If you want to run Terraform with the dynatrace provider plugin on your system, complete the following steps:

1. [Download](https://github.com/dynatrace-oss/terraform-provider-dynatrace/releases/latest) the dynatrace provider plugin for Terraform.

1. Unzip the release archive to extract the plugin binary (`terraform-provider-dynatrace_vX.Y.Z`).

For Terraform version 0.12.x

1. Move the binary into the Terraform [plugins directory] for the platform.
    - Linux/Unix/macOS: `~/.terraform.d/plugins`
    - Windows: `%APPDATA%\terraform.d\plugins`

1. Add the plug-in provider to the Terraform configuration file.

    ```hcl
    terraform {
        required_providers {
            dynatrace = {
                version = "1.1.0"
            }
        }
    }
    ```
    

For Terraform version 0.13.x

1. Move the binary into the Terraform [plugins directory] for the platform.
    - Linux: `~/.terraform.d/plugins/dynatrace.com/com/dynatrace/1.0.4/linux_amd64/`
    - macOS: `~/.terraform.d/plugins/dynatrace.com/com/dynatrace/1.0.4/darwin_amd64/`
    - Windows: `%APPDATA%\terraform.d\plugins\dynatrace.com\com\dynatrace\1.0.4\windows_amd64\`

1. Add the plug-in provider to the Terraform configuration file.

    ```hcl
    terraform {
        required_version = "~> 0.13.0"
        required_providers {
            dynatrace = {
                version = "1.1.0"
                source = "dynatrace.com/com/dynatrace"
            }
        }
    }
    ```

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
