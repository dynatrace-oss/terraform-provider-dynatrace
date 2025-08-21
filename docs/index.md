---
layout: ""
page_title: "Provider: Dynatrace"
description: |-
  The Dynatrace Terraform Provider provides resources to interact with the Dynatrace REST API.
---

# Dynatrace Terraform Provider

The Dynatrace Terraform Provider is used to interact with the resources supported by the Dynatrace REST API. The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to learn about the available resources and data sources. For additional information, refer to the official Dynatrace documentation on Terraform available [here](https://dt-url.net/3s63qyj).

The Dynatrace Terraform Provider is officially supported by Dynatrace.

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

To create an API token with full access to all resources, define `DYNATRACE_API_TOKEN` as an [access token](https://docs.dynatrace.com/docs/manage/identity-access-management/access-tokens-and-oauth-clients/access-tokens) with the following permissions. 
 * **Read settings** (`settings.read`)
 * **Write settings** (`settings.write`)
 * **Read configuration** (`ReadConfig`)
 * **Write configuration** (`WriteConfig`)
 * **Capture request data** (`CaptureRequestData`)
 * **Create and read synthetic monitors, locations, and nodes** (`ExternalSyntheticIntegration`)
 * **Create ActiveGate tokens** (`activeGateTokenManagement.create`)
 * **Read ActiveGate tokens** (`activeGateTokenManagement.read`)
 * **Write ActiveGate tokens** (`activeGateTokenManagement.write`)
 * **Read API tokens** (`apiTokens.read`)
 * **Write API tokens** (`apiTokens.write`)
 * **Read attacks** (`attacks.read`)
 * **Write attacks** (`attacks.write`)
 * **Read credential vault entries** (`credentialVault.read`)
 * **Write credential vault entries** (`credentialVault.write`)
 * **Read Entities** (`entities.read`)
 * **Write extensions** (`extensions.write`)
 * **Read extensions environment configuration** (`extensionEnvironment.read`)
 * **Write extensions environment configuration** (`extensionEnvironment.write`)
 * **Read network zones** (`networkZones.read`)
 * **Write network zones** (`networkZones.write`)
 * **Read security problems** (`securityProblems.read`)
 * **Write security problems** (`securityProblems.write`)
 * **Read SLO** (`slo.read`)
 * **Write SLO** (`slo.write`)

Configure an [OAuth client](https://dt-url.net/fj43qif) with all of the permissions below to be compatible with all OAuth based Terraform resources, or provide a subset of permissions based off of required use cases - refer to the resource specific pages for additional information. 

Define `DT_CLIENT_ID`, `DT_CLIENT_SECRET`, `DT_ACCOUNT_ID` based off of the created OAuth client.
 * **View settings objects for schema** (`settings:objects:read`)
 * **Create settings objects for schema** (`settings:objects:write`)
 * **Admin permission to manage all setting objects** (`settings:objects:admin`)
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
 * **Delete documents from trash** (`document:trash.documents:delete`)
 * **Read direct-shares** (`document:direct-shares:read`)
 * **Write direct-shares** (`document:direct-shares:write`)
 * **Delete direct-shares** (`document:direct-shares:delete`)
 * **Read OpenPipeline configurations** (`openpipeline:configurations:read`)
 * **Write OpenPipeline configurations** (`openpipeline:configurations:write`)
 * **View SLOs** (`slo:slos:read`) 
 * **Create and edit SLOs** (`slo:slos:write`)
 * **View SLO objective templates** (`slo:objective-templates:read`)
 * **View bizevents** (`storage:bizevents:read`)
 * **View bucket metadata** (`storage:bucket-definitions:read`)
 * **Write buckets** (`storage:bucket-definitions:write`)
 * **View Filter-Segments** (`storage:filter-segments:read`)
 * **Create and Update Filter-Segments** (`storage:filter-segments:write`)
 * **Share Filter-Segments** (`storage:filter-segments:share`)
 * **Delete Filter-Segments** (`storage:filter-segments:delete`)
 * **Maintain all Filter-Segments on the environment** (`storage:filter-segments:admin`)
 * **View users and groups** (`account-idm-read`)
 * **Manage users and groups** (`account-idm-write`)
 * **View and manage policies** (`iam-policies-management`)
 * **View environments** (`account-env-read`)


### Authenticating with OAuth Credentials

The Terraform Provider supports using OAuth credentials for authentication with endpoints that allow both API Token and OAuth-based access.

> **Note:** Not all resources currently support OAuth. For example, the `dynatrace_json_dashboard` resource can only be configured using API Tokens.

To enable OAuth-based authentication, set the environment variable:

```
DYNATRACE_HTTP_OAUTH_PREFERENCE=true
```

When this variable is set and OAuth credentials (e.g., `DT_CLIENT_ID` and `DT_CLIENT_SECRET`) are provided, the provider will prioritize using REST endpoints that support OAuth.

If `DYNATRACE_HTTP_OAUTH_PREFERENCE` is not set or is not `true`, the provider will default to using an API Token for authentication.

---

### Authenticating with Platform Tokens

You can authenticate using a Platform Token by setting the environment variable:

```
DYNATRACE_PLATFORM_TOKEN
```

Alternatively, you can use the `platform_token` attribute in the provider configuration.

If `DYNATRACE_PLATFORM_TOKEN` is not defined, the provider will use the configured OAuth credentials (`DT_CLIENT_ID` and `DT_CLIENT_SECRET`) to obtain a Bearer token.

Platform Token authentication follows the same selection rules as OAuth credentials:  
When the environment variable `DYNATRACE_HTTP_OAUTH_PREFERENCE` is set to `true`, the provider will favor Platform or OAuth tokens over API Tokens.


## Exporting existing configuration from a Dynatrace environment
In addition to the out-of-the-box functionality of Terraform, the provider has the ability to be executed as a standalone executable to export an existing configuration from a Dynatrace environment. Refer to the [Export Utility](https://dt-url.net/h203qmc) page for more information.

## How can I troubleshoot the Dynatrace Terraform provider?
Set the following environment variables to generate logs in the working directory.
```
export DYNATRACE_DEBUG=true
export DYNATRACE_LOG_HTTP=terraform-provider-dynatrace.http.log
export DYNATRACE_HTTP_RESPONSE=true
```
For any assistance, please create a [GitHub](https://github.com/dynatrace-oss/terraform-provider-dynatrace/issues) issue.