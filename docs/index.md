---
layout: ""
page_title: "Provider: Dynatrace"
description: |-
  The Dynatrace Terraform provider provides resources to interact with the Dynatrace API.
---

# Dynatrace Terraform Provider

The Dynatrace Terraform provider is used to interact with the resources supported by the Dynatrace API. The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to learn about the available resources and data sources. For additional information, refer to the official Dynatrace documentation on Terraform available [here](https://dt-url.net/3s63qyj).

The Dynatrace Terraform provider is officially supported by Dynatrace.

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
* SaaS environments: `https://########.live.dynatrace.com`
* Managed environments: `https://<dynatrace-host>/e/#####################`

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

The Dynatrace Terraform provider supports using OAuth credentials for authentication with endpoints that allow both API token and OAuth-based access.

> **Note:** Not all resources support OAuth. For example, the `dynatrace_json_dashboard` resource can only be configured using API tokens.

To enable OAuth-based authentication, set the environment variable:

```
DYNATRACE_HTTP_OAUTH_PREFERENCE=true
```

When this variable is set and OAuth credentials (e.g., `DT_CLIENT_ID` and `DT_CLIENT_SECRET`) are provided, the provider will prioritize using REST endpoints that support OAuth.

If `DYNATRACE_HTTP_OAUTH_PREFERENCE` is not set or is not `true`, the provider will default to using an API token for authentication.

---

### Authenticating with platform tokens

You can authenticate using a platform token by setting the environment variable:

```
DYNATRACE_PLATFORM_TOKEN=<<PLATFORM_TOKEN>>
```

Alternatively, you can use the `platform_token` attribute in the provider configuration.

If `DYNATRACE_PLATFORM_TOKEN` is not defined, the provider will use the configured OAuth credentials (`DT_CLIENT_ID` and `DT_CLIENT_SECRET`) to obtain a Bearer token.

Platform token authentication follows the same selection rules as OAuth credentials:
When the environment variable `DYNATRACE_HTTP_OAUTH_PREFERENCE` is set to `true`, the provider will favor platform or OAuth tokens over API tokens.

---

### Authenticating with Workload Identity Federation (GitHub Actions OIDC)

Inside a GitHub Actions job the provider can obtain a short lived OIDC token and use it for the Dynatrace platform APIs. No long lived secret has to be stored in the repository.

Set the vendor and the audience:

```
DYNATRACE_WIF_VENDOR=github
DYNATRACE_WIF_AUDIENCE=<<AUDIENCE>>
```

Alternatively, you can use the `wif_vendor` and `wif_audience` attributes in the provider configuration. `wif_audience` has no default: the value has to match the audience that your Dynatrace environment expects.

The job needs the `id-token: write` permission, otherwise GitHub does not issue an OIDC token to it:

```yaml
name: terraform

on:
  push:
    branches: [main]

permissions:
  contents: read
  id-token: write   # required - lets the job request a GitHub OIDC token

jobs:
  apply:
    runs-on: ubuntu-latest
    env:
      DYNATRACE_ENV_URL: https://########.apps.dynatrace.com
      DYNATRACE_WIF_VENDOR: github
      DYNATRACE_WIF_AUDIENCE: <<AUDIENCE>>
    steps:
      - uses: actions/checkout@v4
      - uses: hashicorp/setup-terraform@v3
      - run: terraform init
      - run: terraform apply -auto-approve
```

> **Note:** If the permission is missing, the provider fails with `Unable to get ACTIONS_ID_TOKEN_REQUEST_URL environment variable`.

The provider obtains a new token shortly before the current one expires, so a long running `terraform apply` does not fail halfway through when the token reaches the end of its short lifetime.

Workload Identity Federation applies to the Dynatrace **platform** APIs only. IAM (Account Management) resources still require OAuth credentials, and classic resources still require an API token. When it is configured, it takes precedence over `platform_token` and over OAuth credentials for platform requests, and the provider emits a warning naming the credentials that are being ignored.

#### Using an OIDC token obtained elsewhere

If you obtain the token by other means, you can pass it directly:

```
DYNATRACE_WIF_OIDC_TOKEN=<<OIDC_TOKEN>>
```

or use the `wif_oidc_token` attribute. It conflicts with `wif_vendor` - specify one or the other, never both.

> **Warning:** A token supplied this way is **never replaced**. It carries a fixed `exp` claim and can expire in the middle of a long running `terraform apply`, which shows up as `401` responses partway through. Prefer `wif_vendor` whenever the workload runs on a supported CI platform.


## Exporting existing configuration from a Dynatrace environment
In addition to the out-of-the-box functionality of Terraform, the provider has the ability to be executed as a standalone executable to export an existing configuration from a Dynatrace environment. Refer to the [Export Utility](https://dt-url.net/h203qmc) page for more information.

## How can I troubleshoot the Dynatrace Terraform provider?
Set the following environment variables to generate logs in the working directory.
```
export DYNATRACE_DEBUG=true
export DYNATRACE_LOG_HTTP=terraform-provider-dynatrace.http.log
export DYNATRACE_HTTP_RESPONSE=true
```
For assistance, please contact the Dynatrace Support team as described on the [support page](https://support.dynatrace.com/).

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `account_id` (String, Sensitive) The Dynatrace account ID (UUID). Required for IAM (Account Management) resources. Also serves as a fallback for `iam_account_id`.
- `automation_client_id` (String, Sensitive) The client ID of an OAuth client used for platform APIs. Falls back to `client_id` if not specified.
- `automation_client_secret` (String, Sensitive) The client secret of an OAuth client used for platform APIs. Falls back to `client_secret` if not specified.
- `automation_env_url` (String) The URL of the Dynatrace platform environment (`https://#####.apps.dynatrace.com`). Falls back to `dt_env_url` if not specified.
- `automation_token_url` (String) The token URL for obtaining access tokens via OAuth for the platform APIs. Default: `https://sso.dynatrace.com/sso/oauth2/token`.
- `client_id` (String, Sensitive) The client ID of an OAuth client used for  platform APIs. Also serves as a fallback for `iam_client_id` and `automation_client_id`.
- `client_secret` (String, Sensitive) The client secret of an OAuth client used for platform APIs. Also serves as a fallback for `iam_client_secret` and `automation_client_secret`.
- `dt_api_token` (String, Sensitive) The API token for classic Dynatrace APIs.
- `dt_cluster_api_token` (String, Sensitive) The API token for Dynatrace Managed cluster APIs.
- `dt_cluster_url` (String, Sensitive) The URL of the Dynatrace Managed cluster.
- `dt_env_url` (String) The URL of the Dynatrace environment (e.g. `https://#####.live.dynatrace.com` or `https://#####.apps.dynatrace.com`).
- `iam_account_id` (String, Sensitive) The Dynatrace account ID (UUID). Required for IAM (Account Management) resources. Falls back to `account_id` if not specified.
- `iam_client_id` (String, Sensitive) The client ID of an OAuth client used for the IAM (Account Management) API. Falls back to `client_id` if not specified.
- `iam_client_secret` (String, Sensitive) The client secret of an OAuth client used for the IAM (Account Management) API. Falls back to `client_secret` if not specified.
- `iam_endpoint_url` (String, Sensitive) The endpoint URL for the IAM (Account Management) API. Default: `https://api.dynatrace.com`.
- `iam_token_url` (String, Sensitive) The token URL for obtaining access tokens via OAuth for the IAM (Account Management) API. Default: `https://sso.dynatrace.com/sso/oauth2/token`.
- `platform_token` (String, Sensitive) The Dynatrace platform token used for platform APIs. When specified, it is used in preference to `client_id`, `client_secret`, `automation_client_id`, `automation_client_secret`, `automation_token_url`, and `automation_env_url` for platform requests. Platform tokens can't be used for IAM (Account Management) or classic resources.
- `wif_audience` (String) The audience (`aud` claim) requested for the Workload Identity Federation OIDC token. Required whenever `wif_vendor` is specified. There is no default: the value has to match the audience that the Dynatrace environment expects.
- `wif_oidc_token` (String, Sensitive) An already obtained OIDC token (JWT) that is sent to the platform APIs as it is. Conflicts with `wif_vendor`. This token is never replaced: it carries a fixed `exp` claim and can expire in the middle of a long running `terraform apply`. Prefer `wif_vendor` whenever the workload runs on a supported CI platform.
- `wif_vendor` (String) Obtains a short lived OIDC token from the given workload identity provider and uses it for platform APIs. The only supported value is `github`, which requires the job to run in GitHub Actions with the `id-token: write` permission. Requires `wif_audience` and conflicts with `wif_oidc_token`. When specified, it is used in preference to `platform_token` and to OAuth credentials for platform requests. Workload Identity Federation can't be used for IAM (Account Management) or classic resources.
