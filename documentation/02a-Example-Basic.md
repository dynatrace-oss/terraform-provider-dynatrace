# Hands-on: Terraform Dynatrace Basic Example

To get you started with managing configurations, this section will guide you through a simple example of creating a management zone with Terraform. You will learn how to create, update, and destroy a configuration.

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


To create a token that would work for all configuration, please include the following permissions in addition to the above.
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

## Build Configuration
This example will cover configuration of creating a management zone for a web application.

1. To get started, verify that you have completed the initial steps outlined in [Getting Started with Terraform and the Dynatrace Provider](LINK). You should have a working directory with the Dynatrace Terraform provider before proceeding to the next steps.

2. In the working directory, create a `main.tf` file with the following contents. Terraform configuration is made up of resource block(s) that represent configuration. For additional information on the management zone resource, please refer to the Terraform Registry [page](https://registry.terraform.io/providers/dynatrace-oss/dynatrace/latest/docs/resources/management_zone_v2).
```terraform
resource "dynatrace_management_zone_v2" "TerraformExample" {
  name = "Terraform Example"
  rules {
    rule {
      type    = "ME"
      enabled = true
      attribute_rule {
        entity_type = "WEB_APPLICATION"
        attribute_conditions {
          condition {
            case_sensitive = true
            key            = "WEB_APPLICATION_NAME"
            operator       = "EQUALS"
            string_value   = "easyTravel"
          }
        }
      }
    }
  }
}
```
-> Manually creating Terraform configuration can be difficult, the option to export existing configuration from a tenant is covered in the [export utility](LINK) guide.

1. Open a terminal window and set the tenant URL and API token environment variables. This will be the tenant we want to push configuration. On Windows, SaaS:
```
set DYNATRACE_ENV_URL=https://########.live.dynatrace.com
set DYNATRACE_API_TOKEN=dt0c01.########.########
```
Managed:
```
set DYNATRACE_ENV_URL=https://<dynatrace-host>/e/########
set DYNATRACE_API_TOKEN=dt0c01.########.########
```

1. In the terminal window, navigate to the working directory. The next step is to run `terraform plan` in order to get an execution plan which will provide a preview of the changes Terraform plans to make to the tenant.
```terraform
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # dynatrace_management_zone_v2.TerraformExample will be created
  + resource "dynatrace_management_zone_v2" "TerraformExample" {
      + id        = (known after apply)
      + legacy_id = (known after apply)
      + name      = "Terraform Example"

      + rules {
          + rule {
              + enabled = true
              + type    = "ME"

              + attribute_rule {
                  + entity_type = "WEB_APPLICATION"

                  + attribute_conditions {
                      + condition {
                          + case_sensitive = true
                          + key            = "WEB_APPLICATION_NAME"
                          + operator       = "EQUALS"
                          + string_value   = "easyTravel"
                        }
                    }
                }
            }
        }
    }

Plan: 1 to add, 0 to change, 0 to destroy.
```

5. Once the plan is verified, running `terraform apply` will execute the actions proposed in the plan. In this case, it will push the management zone configuration to the tenant.
```terraform
dynatrace_management_zone_v2.TerraformExample: Creating...
dynatrace_management_zone_v2.TerraformExample: Creation complete after 1s [id=*************)]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

Note the `terraform.tfstate` file that was created with the apply command. This is the Terraform state file that is automatically generated to keep track of the resources that Terraform is currently managing. On subsequent executions, Terraform will utilize this file and compare configuration to the code and make any adjustments.

This concludes this section, the following sections will cover changing and destroying configuration. 

## Change Configuration
In the previous section, you created a management zone for a web application. This section will continue building off of the example by learning how to modify a resource.

1. In the terminal window, verify you are in the working directory. Execute `terraform plan` which should return that no changes are needed.
```terraform
dynatrace_management_zone_v2.easyTravel: Refreshing state... [id=*************]

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
```

2. Next we will make a change to the configuration. Open `main.tf` then change the value for `string_value`, eg. "easyTravel" to "Terraform".
```terraform
resource "dynatrace_management_zone_v2" "TerraformExample" {
  name = "Terraform Example"
  rules {
    rule {
      type    = "ME"
      enabled = true
      attribute_rule {
        entity_type = "WEB_APPLICATION"
        attribute_conditions {
          condition {
            case_sensitive = true
            key            = "WEB_APPLICATION_NAME"
            operator       = "EQUALS"
            string_value   = "Terraform"
          }
        }
      }
    }
  }
}
```

3. Since we intend to apply the local changes we made, we can directly execute `terraform apply`. This will generate an execution plan displaying the changes and ask to apply.
```terraform
dynatrace_management_zone_v2.easyTravel: Modifying... [id=*************]
dynatrace_management_zone_v2.easyTravel: Modifications complete after 0s [id=*************]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```

The management zone configuration in the Dynatrace environment and the Terraform state file is now updated based off of the changes. 

This concludes this section of the exercise of changing configuration of a Terraform resource.

## Destroy Configuration
In the previous sections, you created and modified a management zone for a web application. This last section will cover how to delete configuration managed by Terraform.

1. In the terminal window, verify you are in the working directory. Execute `terraform plan` which should return that no changes are needed.
```terraform
dynatrace_management_zone_v2.easyTravel: Refreshing state... [id=*************]

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
```

2. Delete the management zone by executing `terraform destroy`. This will generate an execution plan and ask to destroy.

```terraform
dynatrace_management_zone_v2.easyTravel: Destroying... [id=*************]
dynatrace_management_zone_v2.easyTravel: Destruction complete after 0s

Destroy complete! Resources: 1 destroyed.
```

The management zone configuration in the Dynatrace environment has been destroyed and the Terraform state file is now empty.

## What's Next
-> The full list of provider resources is available on the [Terraform Registry](https://registry.terraform.io/providers/dynatrace-oss/dynatrace/latest/docs).

Hands-on: [Terraform Dynatrace Advanced Example](LINK)