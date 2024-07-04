---
layout: ""
page_title: "Provider: Dynatrace"
description: |-
  The Dynatrace Terraform Provider provides resources to interact with the Dynatrace REST API.
---

# Dynatrace Terraform Provider

The Dynatrace Terraform Provider is used to interact with the resources supported by the Dynatrace REST API. The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to learn about the available resources and data sources. For additional information, refer to the official Dynatrace documentation on Terraform available [here](https://docs.dynatrace.com/docs/manage/configuration-as-code/terraform).

## Example

# Terraform 0.13+ uses the Terraform Registry:
```
terraform {
    required_providers {
        dynatrace = {
            version = "~> 1.0"
            source = "dynatrace-oss/dynatrace"
        }
    }
} 
```
## Configure the Dynatrace provider
The recommended approach is to configure the provider via environment variables.

Define `DYNATRACE_ENV_URL` for the Dynatrace environment URL.
* SaaS Environments: `https://########.live.dynatrace.com`
* Managed Enviroments: `https://<dynatrace-host>/e/#####################`

Define `DYNATRACE_API_TOKEN` as an access token with the following permissions.
 * **Read configuration** (`ReadConfig`)
 * **Write configuration** (`WriteConfig`)
 * **Read settings** (`settings.read`)
 * **Write settings** (`settings.write`)
 * **Create and read synthetic monitors, locations, and nodes** (`ExternalSyntheticIntegration`)
 * **Capture request data** (`CaptureRequestData`)
 * **Read credential vault entries** (`credentialVault.read`)
 * **Write credential vault entries** (`credentialVault.write`)
 * **Read network zones** (`networkZones.read`)
 * **Write network zones** (`networkZones.write`)
 * **Read security problems** (`securityProblems.read`)
 * **Write security problems** (`securityProblems.write`)
 * **Read attacks** (`attacks.read`)
 * **Write attacks** (`attacks.write`)

Configure an OAuth client with all of the permissions below to be compatible with all OAuth based Terraform resources, or provide a subset of permissions based off of required use cases - refer to the resource specific pages for additional information. 

Define `DT_CLIENT_ID`, `DT_CLIENT_SECRET`, `DT_ACCOUNT_ID` based off of the created OAuth client.
 * **View settings objects for schema** (`settings:objects:read`)
 * **Create settings objects for schema** (`settings:objects:write`)
 * **View workflows** (`automation:workflows:read`)
 * **Create and edit workflows** (`automation:workflows:write`)
 * **View calendars** (`automation:calendars:read`)
 * **Create and edit calendars** (`automation:calendars:write`)
 * **View rules** (`automation:rules:read`)
 * **Create and edit rules** (`automation:rules:write`)
 * **Admin permission to manage all workflows and executions** (`automation:workflows:admin`)
 * **View documents** (`document:documents:read`)
 * **Create and edit documents** (`document:documents:write`)
 * **Delete documents** (`document:documents:delete`)
 * **Read direct-shares** (`document:direct-shares:read`)
 * **Write direct-shares** (`document:direct-shares:write`)
 * **Delete direct-shares** (`document:direct-shares:delete`)
 * **View bizevents** (`storage:bizevents:read`)
 * **View bucket metadata** (`storage:bucket-definitions:read`)
 * **Write buckets** (`storage:bucket-definitions:write`)
 * **View users and groups** (`account-idm-read`)
 * **Manage users and groups** (`account-idm-write`)
 * **View and manage policies** (`iam-policies-management`)
 * **View environments** (`account-env-read`)

## Exporting existing configuration from a Dynatrace environment
In addition to the out-of-the-box functionality of Terraform, the provider has the ability to be executed as a standalone executable to export an existing configuration from a Dynatrace environment. Refer to the [Export Utility](https://docs.dynatrace.com/docs/manage/configuration-as-code/terraform/guides/export-utility) page for more information.
