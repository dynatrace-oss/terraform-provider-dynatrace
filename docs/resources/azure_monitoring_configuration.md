---
layout: ""
page_title: "dynatrace_azure_monitoring_configuration Resource - terraform-provider-dynatrace"
subcategory: "Extensions 2.0"
description: |-
  The resource `dynatrace_azure_monitoring_configuration` manages an Extensions 2.0 monitoring configuration for `com.dynatrace.extension.da-azure`.
---

# dynatrace_azure_monitoring_configuration (Resource)

-> This resource is part of the Dynatrace Application Cloud (DAC) typed-resource set. It wraps an Extensions 2.0 monitoring configuration object (extension `com.dynatrace.extension.da-azure`) with a strongly-typed schema, mirroring the wire shape produced by `dtctl create azure`.

-> This resource requires the platform token scopes **Read settings** (`settings:objects:read`) and **Write settings** (`settings:objects:write`), plus **Read extensions** (`extensions.read`) and **Write extension monitoring configurations** (`extensionConfigurations.write`).

## Requirements

This resource depends on a `dynatrace_azure_connection` (HAS) and, when using workload-identity federation, the `dynatrace_azure_connection_authentication` patch resource. The Azure Service Principal referenced via `service_principal_id` must have the required RBAC role assignments (typically `Monitoring Reader` + `Reader` at subscription scope) before the monitoring configuration is created.

A complete end-to-end Terraform example wiring the Entra ID app registration, the federated identity credential, the HAS connection and the monitoring configuration together is available under [`azure-poc-sample/`](https://github.com/dynatrace-oss/terraform-provider-dynatrace/tree/main/azure-poc-sample) in this fork.

## Dynatrace Documentation

- Microsoft Azure integration — https://docs.dynatrace.com/docs/ingest-from/microsoft-azure
- Extensions 2.0 monitoring configurations API — https://docs.dynatrace.com/docs/dynatrace-api/environment-api/extensions-v2

## Resource Example Usage — federated identity

```terraform
resource "dynatrace_azure_connection" "this" {
  name = "dac-tf-poc-connection"
  type = "federatedIdentityCredential"
  federated_identity_credential {
    consumers = ["SVC:com.dynatrace.da"]
  }
}

resource "azuread_application_federated_identity_credential" "dynatrace" {
  application_id = azuread_application.dynatrace.id
  display_name   = "dt-fd"
  audiences      = ["${var.dt_env_url}/svc-id/com.dynatrace.da"]
  issuer         = "https://dev.token.dynatracelabs.com"
  subject        = "dt:connection-id/${dynatrace_azure_connection.this.id}"
}

resource "dynatrace_azure_connection_authentication" "this" {
  azure_connection_id = dynatrace_azure_connection.this.id
  directory_id        = data.azuread_client_config.current.tenant_id
  application_id      = azuread_application.dynatrace.client_id
}

resource "dynatrace_azure_monitoring_configuration" "this" {
  depends_on = [dynatrace_azure_connection_authentication.this]

  name    = "dac-tf-poc-monitoring"
  enabled = true

  credential {
    connection_id        = dynatrace_azure_connection.this.id
    service_principal_id = azuread_application.dynatrace.client_id
    type                 = "FEDERATED"
  }

  regions = ["eastus"]

  feature_sets = [
    "microsoft_compute.virtualmachines_essential",
    "microsoft_storage.storageaccounts_essential",
  ]

  # Optional: filter monitored resources by Azure tag
  tag_filter {
    key       = "environment"
    value     = "production"
    condition = "INCLUDE"
  }

  # Optional: promote Azure tag keys onto Dynatrace entities
  tag_enrichment = ["owner", "cost-center"]

  # Optional: enrich monitored entities with static or tag-derived dt.* labels
  dt_label_enrichment {
    label   = "dt.security_context"
    literal = "my-app"
  }
}
```

## Resource Example Usage — client secret

```terraform
resource "dynatrace_azure_connection" "this" {
  name = "dac-tf-poc-connection"
  type = "clientSecret"
  client_secret {
    client_secret  = azuread_application_password.dynatrace.value
    application_id = azuread_application.dynatrace.client_id
    directory_id   = data.azuread_client_config.current.tenant_id
    consumers      = ["SVC:com.dynatrace.da"]
  }
}

resource "dynatrace_azure_monitoring_configuration" "this" {
  name    = "dac-tf-poc-monitoring"
  enabled = true

  credential {
    connection_id        = dynatrace_azure_connection.this.id
    service_principal_id = azuread_application.dynatrace.client_id
    type                 = "SECRET"
  }

  regions      = ["eastus"]
  feature_sets = ["microsoft_compute.virtualmachines_essential"]
}
```

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_azure_monitoring_configuration` downloads all existing Azure monitoring configurations.

## Notes

- `extension_version` is **sticky in state**. When omitted on Create, the provider queries the installed extension list and pins the highest version available on the tenant. Updates do not touch the pinned version unless the user changes it explicitly.
- The wire API echoes `namespaces[]` and `eventHubsConfiguration[]` arrays assembled server-side from the configured feature sets and tag-based Event Hub autodiscovery (`dt-log-ingest-activated` tag). The provider deliberately ignores those echoes during read to avoid eternal drift — manual per-namespace pinning is not supported.
- `credential.type` defaults to `FEDERATED` when absent in the API response, matching dtctl's `ResolveCredential` behaviour.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `credential` (Block List, Min: 1) HAS connection + Service Principal binding. At least one is required. dtctl always writes exactly one, but the API accepts a list. (see [below for nested schema](#nestedblock--credential))
- `name` (String) Human-readable name of the monitoring configuration (written to `description` in the extension config payload).

### Optional

- `activation_context` (String) Extension activation context. Defaults to `DATA_ACQUISITION`.
- `configuration_mode` (String) Configuration mode. Defaults to `ADVANCED`.
- `deployment_mode` (String) Deployment mode. Defaults to `AUTOMATED`.
- `deployment_scope` (String) Deployment scope. Defaults to `SUBSCRIPTION`. The other observed value is `MANAGEMENT_GROUP`.
- `dt_label_enrichment` (Block List) Dynatrace labels (`dt.*`) applied to every monitored entity. Each block sets exactly one of `literal` or `tag_key`. (see [below for nested schema](#nestedblock--dt_label_enrichment))
- `enabled` (Boolean) Whether the monitoring configuration is active. Defaults to true.
- `extension_version` (String) Version of `com.dynatrace.extension.da-azure` that this configuration targets. Optional — when omitted at create time, the provider picks the highest semver version installed on the tenant (same behavior as `dtctl create azure monitoring`). The resolved value is persisted to state. On subsequent refreshes the provider reads back whatever version Dynatrace currently reports for this configuration; if the extension was auto-updated (or bumped manually) the new version surfaces as drift in `terraform plan`, but no Terraform-driven update silently re-resolves it. To pin a version, set it explicitly here.
- `feature_sets` (Set of String) Azure feature sets to enable (e.g. `microsoft_compute.virtualmachines_essential`). When empty, the extension defaults are used.
- `regions` (Set of String) Azure regions (locations) to monitor, e.g. `eastus`. Empty set = all locations the extension knows about. Maps to `locationFiltering` on the wire.
- `scope` (String) Settings 2.0 scope. Defaults to `integration-azure`. Changing it forces recreation.
- `subscription_filter` (Set of String) Subscription GUIDs to include or exclude (per `subscription_filtering_mode`). Empty set means "all subscriptions reachable by the Service Principal".
- `subscription_filtering_mode` (String) How to interpret `subscription_filter`. Defaults to `INCLUDE`.
- `tag_enrichment` (Set of String) Azure tag keys whose values are copied as Dynatrace tags on monitored entities.
- `tag_filter` (Block List) Filter monitored resources by Azure tag. Repeat the block to define multiple filters. (see [below for nested schema](#nestedblock--tag_filter))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--credential"></a>
### Nested Schema for `credential`

Required:

- `connection_id` (String) ObjectId of the `dynatrace_azure_connection` resource.
- `service_principal_id` (String) Azure application (client) id of the Service Principal that Dynatrace impersonates. Typically `azuread_application.<name>.client_id`.

Optional:

- `description` (String) Free-form description for this credential. Defaults to the top-level `name`.
- `enabled` (Boolean) Per-credential enable flag. Defaults to true. Distinct from the top-level `enabled`.
- `type` (String) `FEDERATED` (workload-identity federation) or `SECRET` (client secret). Defaulted from the connection's auth mode when omitted.


<a id="nestedblock--dt_label_enrichment"></a>
### Nested Schema for `dt_label_enrichment`

Required:

- `label` (String) Dynatrace label key, e.g. `dt.security_context` or `dt.cost.product`.

Optional:

- `literal` (String) Static value applied to every monitored entity. Mutually exclusive with `tag_key`.
- `tag_key` (String) Azure tag key whose value will be copied into the Dynatrace label. Mutually exclusive with `literal`.


<a id="nestedblock--tag_filter"></a>
### Nested Schema for `tag_filter`

Required:

- `condition` (String) `INCLUDE` to only monitor matching resources, `EXCLUDE` to skip them.
- `key` (String) Azure tag key.
- `value` (String) Azure tag value to match.
