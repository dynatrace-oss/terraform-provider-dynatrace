# Hands-on: Terraform Dynatrace Advanced Example

This section will guide you through an advanced example of creating a Dynatrace configuration template for a newly deployed application with Terraform.

-> This example refers to a generic Dynatrace configuration template and not a Terraform [templatefile](https://developer.hashicorp.com/terraform/language/functions/templatefile).

## Prerequisites
* Terraform CLI with the Dynatrace provider installed (see [Getting Started with Terraform and the Dynatrace Provider](LINK)) and available under PATH.
* A Dynatrace environment and access to create environment tokens.
* A Dynatrace token with at least the following permissions:
 * **Read settings** (`settings.read`)
     * Required for reading Settings 2.0 configuration.
 * **Write settings** (`settings.write`)
     * Required for modifying Settings 2.0 configuration.

To create a token that would work for all configuration, please include the following permissions in addition to the above.
 * **Read configuration** (`ReadConfig`)
     * Required for reading API v1 configuration.
 * **Write configuration** (`WriteConfig`)
     * Required for creating API v1 configuration.
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

## Example

In this example, we provide a JSON formatted input file which gets passed into the Terraform configuration to automatically create a management zone, alerting profile, and email notification for each defined application.

1. To get started, verify that you have completed the initial steps outlined in [Getting Started with Terraform and the Dynatrace Provider](LINK). You should have a working directory with the Dynatrace Terraform provider before proceeding to the next steps.

2. In the working directory, create a `main.tf` file with the following contents. For more information on each resource, please refer to the Terraform Registry [page](https://registry.terraform.io/providers/dynatrace-oss/dynatrace/latest/docs).

The `locals` block makes the contents of `data.json` available for lookup within the `resource` blocks.

```terraform
locals {
  app_data = jsondecode(file("data.json"))
}

resource "dynatrace_management_zone_v2" "mgmz_per_app" {
  for_each = local.app_data
  name = each.key
  rules {
    rule {
      type    = "ME"
      enabled = true
      attribute_rule {
        entity_type           = "HOST"
        host_to_pgpropagation = true
        attribute_conditions {
          condition {
            case_sensitive = true
            key            = "HOST_GROUP_NAME"
            operator       = "EQUALS"
            string_value   = each.value["host-group"]
          }
        }
      }
    }
  }
}

resource "dynatrace_alerting" "alerting_per_app" {
  for_each = dynatrace_management_zone_v2.mgmz_per_app
  name            = each.value.name
  management_zone = each.value.legacy_id
  rules {
    rule {
      delay_in_minutes = local.app_data[each.value.name]["delay-in-minutes"]
      include_mode     = "NONE"
      severity_level   = "MONITORING_UNAVAILABLE"
   }
  }
}

resource "dynatrace_email_notification" "email_notification_per_app" {
  for_each = dynatrace_alerting.alerting_per_app

  name                   = each.value.name
  subject                = "{State} Problem {ProblemID}: {ImpactedEntity}"
  to                     = local.app_data[each.value.name]["notify"]
  cc                     = ["terraform@dynatrace.com"]
  body                   = "{ProblemDetailsHTML}"
  profile                = each.value.id
  active                 = true
  notify_closed_problems = true
}
```

3. In the working directory, create a `data.json` file with the following contents. 
```json
{
    "App-A": {
        "host-group": "group-a",
        "delay-in-minutes": 20,
        "notify": ["app.a.owner@dynatrace.com"]
    },
    "App-B": {
        "host-group": "group-b",
        "delay-in-minutes": 30,
        "notify": ["app.b.owner@dynatrace.com"]
    }
}
```

In the JSON example above, we have two applications "App-A" and "App-B". The same name will also be used for the management zone, alerting profile, and email notification. Each application contains the following fields:
* **host-group:** The [host group](https://www.dynatrace.com/support/help/shortlink/host-groups#assign-a-host-to-a-host-group) that should be tied to the application, used for the management zone rule
* **delay-in-minutes:** The [alerting profile](https://www.dynatrace.com/support/help/observe-and-explore/notifications-and-alerting/alerting-profiles) configuration on the length in minutes for the monitoring unavailable event is open before an alert is sent out
* **notify:** The [email notification](https://www.dynatrace.com/support/help/observe-and-explore/notifications-and-alerting/problem-notifications/email-integration) recipients

4. Open a terminal window and set the tenant URL and API token environment variables. This will be the tenant we want to push configuration. On Windows, SaaS:
```
set DYNATRACE_ENV_URL=https://########.live.dynatrace.com
set DYNATRACE_API_TOKEN=dt0c01.########.########
```
Managed:
```
set DYNATRACE_ENV_URL=https://<dynatrace-host>/e/########
set DYNATRACE_API_TOKEN=dt0c01.########.########
```

5. In the terminal window navigate to the working directory and run `terraform apply -auto-approve`. This will generate an execution plan displaying the changes and ask to apply.
```terraform
dynatrace_management_zone_v2.mgmz_per_app["App-A"]: Creating...
dynatrace_management_zone_v2.mgmz_per_app["App-A"]: Creation complete after 1s [id=*************]
dynatrace_management_zone_v2.mgmz_per_app["App-B"]: Creating...
dynatrace_management_zone_v2.mgmz_per_app["App-B"]: Creation complete after 0s [id=*************]
dynatrace_alerting.alerting_per_app["App-B"]: Creating...
dynatrace_alerting.alerting_per_app["App-B"]: Creation complete after 1s [id=*************]
dynatrace_alerting.alerting_per_app["App-A"]: Creating...
dynatrace_alerting.alerting_per_app["App-A"]: Creation complete after 0s [id=*************]
dynatrace_email_notification.email_notification_per_app["App-B"]: Creating...
dynatrace_email_notification.email_notification_per_app["App-B"]: Creation complete after 1s [id=*************]
dynatrace_email_notification.email_notification_per_app["App-A"]: Creating...
dynatrace_email_notification.email_notification_per_app["App-A"]: Creation complete after 1s [id=*************]

Apply complete! Resources: 6 added, 0 changed, 0 destroyed.
```

The "App-A" and "App-B" management zone, alerting profile, and email notification configuration has now been created in the environment. To change or destroy configuration, please follow similar steps outlined in the [Terraform Dynatrace Basic Example](LINK).

## What's Next
[Export Utility](LINK)