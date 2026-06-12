---
layout: ""
page_title: "dynatrace_gcp_monitoring_configuration Resource - terraform-provider-dynatrace"
subcategory: "Extensions 2.0"
description: |-
  The resource `dynatrace_gcp_monitoring_configuration` manages an Extensions 2.0 monitoring configuration for `com.dynatrace.extension.da-gcp`.
---

# dynatrace_gcp_monitoring_configuration (Resource)

-> This resource is part of the Dynatrace Application Cloud (DAC) typed-resource set. It wraps an Extensions 2.0 monitoring configuration object (extension `com.dynatrace.extension.da-gcp`) with a strongly-typed schema, mirroring the wire shape produced by `dtctl create gcp monitoring`.

-> This resource requires the platform token scopes **Read settings** (`settings:objects:read`) and **Write settings** (`settings:objects:write`), plus **Read extensions** (`extensions.read`) and **Write extension monitoring configurations** (`extensionConfigurations.write`).

## Requirements

This resource depends on a [`dynatrace_gcp_connection`](./gcp_connection.md) (Hyperscaler Authentication Service, `serviceAccountImpersonation` mode). The customer GCP service account referenced via `credential.service_account` must:

1. Exist in the target GCP project.
2. Hold the project-level read role required by the integration (typically `roles/viewer`, optionally augmented with `roles/monitoring.viewer`).
3. Grant `roles/iam.serviceAccountTokenCreator` to the Dynatrace-managed principal exposed by the [`dynatrace_gcp_principal`](./gcp_principal.md) resource — this is what enables Dynatrace to impersonate the customer SA.

The GCP APIs the extension calls (`compute.googleapis.com`, `monitoring.googleapis.com`, `cloudresourcemanager.googleapis.com`, `iamcredentials.googleapis.com`, plus any service-specific API for the enabled feature sets) must be enabled on the project.

## Dynatrace Documentation

- Google Cloud Platform integration — https://docs.dynatrace.com/docs/ingest-from/google-cloud-platform
- Extensions 2.0 monitoring configurations API — https://docs.dynatrace.com/docs/dynatrace-api/environment-api/extensions-v2

## Resource Example Usage

```terraform
resource "dynatrace_gcp_principal" "this" {}

resource "google_service_account" "monitoring" {
  account_id   = "dynatrace-dac"
  display_name = "Dynatrace monitoring account"
  project      = var.gcp_project_id
}

resource "google_project_iam_member" "viewer" {
  project = var.gcp_project_id
  role    = "roles/viewer"
  member  = "serviceAccount:${google_service_account.monitoring.email}"
}

resource "google_service_account_iam_member" "dynatrace_token_creator" {
  service_account_id = google_service_account.monitoring.name
  role               = "roles/iam.serviceAccountTokenCreator"
  member             = "serviceAccount:${dynatrace_gcp_principal.this.principal}"
}

resource "dynatrace_gcp_connection" "this" {
  name = "dac-tf-poc-connection"
  type = "serviceAccountImpersonation"
  service_account_impersonation {
    service_account_id = google_service_account.monitoring.email
    consumers          = ["SVC:com.dynatrace.da"]
  }

  depends_on = [google_service_account_iam_member.dynatrace_token_creator]
}

resource "dynatrace_gcp_monitoring_configuration" "this" {
  depends_on = [
    google_project_iam_member.viewer,
    google_service_account_iam_member.dynatrace_token_creator,
  ]

  name    = "dac-tf-poc-monitoring"
  enabled = true

  credential {
    connection_id   = dynatrace_gcp_connection.this.id
    service_account = google_service_account.monitoring.email
  }

  regions        = ["us-central1"]
  project_filter = [var.gcp_project_id]

  feature_sets = [
    "compute_engine_essential",
    "sql_essential",
    "kubernetes_engine_essential",
  ]

  # Optional: filter monitored resources by GCP resource-manager tag.
  # Resource-manager tags are addressed by numeric IDs, not by display name:
  #   tagKeys/<TAG_KEY_ID>  /  tagValues/<TAG_VALUE_ID>
  tag_filter {
    key       = "tagKeys/123456789012"
    value     = "tagValues/123456789013"
    condition = "INCLUDE"
  }

  # Optional: filter by classic GCP label (plain key/value).
  label_filter {
    key       = "team"
    value     = "platform"
    condition = "INCLUDE"
  }

  # Optional: promote GCP tag / label keys onto Dynatrace entities.
  # tag_enrichment uses the same `tagKeys/<ID>` form as tag_filter;
  # label_enrichment uses plain label keys.
  tag_enrichment   = ["tagKeys/123456789012"]
  label_enrichment = ["team"]
}
```

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_gcp_monitoring_configuration` downloads all existing GCP monitoring configurations.

## Notes

- `extension_version` is **sticky in state**. When omitted on Create, the provider queries the installed extension list and pins the highest version available on the tenant. Updates do not touch the pinned version unless the user changes it explicitly.
- The GCP integration distinguishes between two GCP concepts that look similar but are distinct: **tags** (resource-manager tags addressed via `tagKeys/…`, modelled by `tag_filter` / `tag_enrichment`) and **labels** (classic per-resource key/value labels, modelled by `label_filter` / `label_enrichment`). Both arrays are written to separate top-level fields on the wire.
- The wire API echoes `namespaces[]` and `services[]` arrays assembled server-side from the configured feature sets. The provider deliberately ignores those echoes during read to avoid eternal drift — manual per-namespace pinning is not supported.
- `observability_scopes_enabled` is omitted from the wire payload when false to keep the JSON shape identical to `dtctl create gcp monitoring`.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `credential` (Block List, Min: 1) HAS connection + GCP service-account binding. At least one is required. dtctl always writes exactly one, but the API accepts a list. (see [below for nested schema](#nestedblock--credential))
- `name` (String) Human-readable name of the monitoring configuration (written to `description` in the extension config payload).

### Optional

- `enabled` (Boolean) Whether the monitoring configuration is active. Defaults to true.
- `extension_version` (String) Version of `com.dynatrace.extension.da-gcp` that this configuration targets. Optional — when omitted at create time, the provider picks the highest semver version installed on the tenant (same behavior as `dtctl create gcp monitoring`). The resolved value is persisted to state. On subsequent refreshes the provider reads back whatever version Dynatrace currently reports for this configuration; if the extension was auto-updated (or bumped manually) the new version surfaces as drift in `terraform plan`, but no Terraform-driven update silently re-resolves it. To pin a version, set it explicitly here.
- `feature_sets` (Set of String) GCP feature sets to enable (e.g. `compute_engine_essential`, `kubernetes_engine_essential`, `sql_essential`). When empty, the extension defaults are used.
- `folder_filter` (Set of String) GCP folder IDs — fans out to all projects under each folder. Maps to `folderFiltering` on the wire.
- `label_enrichment` (Set of String) GCP label keys whose values are copied as Dynatrace tags on monitored entities.
- `label_filter` (Block List) Filter monitored resources by GCP label (classic key/value labels per resource). Distinct from `tag_filter`. Repeat the block to define multiple filters. (see [below for nested schema](#nestedblock--label_filter))
- `observability_scopes_enabled` (Boolean) Whether observability scopes are enabled. Defaults to false.
- `project_filter` (Set of String) GCP project IDs the customer service account reads. Empty set means "all projects the SA can impersonate into". Maps to `projectFiltering` on the wire.
- `regions` (Set of String) GCP regions (locations) to monitor, e.g. `us-central1`. Empty set = all locations the extension knows about. Maps to `locationFiltering` on the wire.
- `resource_autodiscovery` (Block List) Per-resource-type autodiscovery override. Repeat the block once per `resource_type` you want to override. (see [below for nested schema](#nestedblock--resource_autodiscovery))
- `scope` (String) Settings 2.0 scope. Defaults to `integration-gcp`. Changing it forces recreation.
- `tag_enrichment` (Set of String) GCP tag keys whose values are copied as Dynatrace tags on monitored entities.
- `tag_filter` (Block List) Filter monitored resources by GCP resource-manager tag (`tagKeys/…`). Repeat the block to define multiple filters. (see [below for nested schema](#nestedblock--tag_filter))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--credential"></a>
### Nested Schema for `credential`

Required:

- `connection_id` (String) ObjectId of the `dynatrace_gcp_connection` resource (or the manually-created connection).
- `service_account` (String) The customer GCP service account email Dynatrace impersonates (e.g. `dynatrace-integration@<project>.iam.gserviceaccount.com`).

Optional:

- `description` (String) Free-form description for this credential. Defaults to the top-level `name`.
- `enabled` (Boolean) Per-credential enable flag. Defaults to true. Distinct from the top-level `enabled`.


<a id="nestedblock--label_filter"></a>
### Nested Schema for `label_filter`

Required:

- `condition` (String) `INCLUDE` to only monitor matching resources, `EXCLUDE` to skip them.
- `key` (String) GCP tag or label key.
- `value` (String) GCP tag or label value to match.


<a id="nestedblock--resource_autodiscovery"></a>
### Nested Schema for `resource_autodiscovery`

Required:

- `auto_discovery_enabled` (Boolean) Whether autodiscovery is enabled for this resource type.
- `resource_type` (String) GCP monitored resource type in the form `<service>.googleapis.com/<Kind>`, e.g. `compute.googleapis.com/Instance`.

Optional:

- `exclude_metric_type` (Set of String) Metric types to exclude from autodiscovery.


<a id="nestedblock--tag_filter"></a>
### Nested Schema for `tag_filter`

Required:

- `condition` (String) `INCLUDE` to only monitor matching resources, `EXCLUDE` to skip them.
- `key` (String) GCP tag or label key.
- `value` (String) GCP tag or label value to match.
