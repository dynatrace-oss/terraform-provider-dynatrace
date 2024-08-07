---
layout: ""
page_title: dynatrace_discovery_feature_flags Resource - terraform-provider-dynatrace"
subcategory: "AppEngine"
description: |-
  The resource `dynatrace_discovery_feature_flags` covers configuration for Discovery and Coverage app feature flags
---

# dynatrace_discovery_feature_flags (Resource)

-> This resource requires the API token scopes **Read settings** (`settings.read`) and **Write settings** (`settings.write`)

## Dynatrace Documentation

- Discovery & Coverage - https://www.dynatrace.com/hub/detail/discovery-coverage/

- Settings API - https://www.dynatrace.com/support/help/dynatrace-api/environment-api/settings (schemaId: `app:dynatrace.discovery.coverage:feature-flags`)

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_discovery_feature_flags` downloads existing Discovery and Coverage app feature flags

The full documentation of the export feature is available [here](https://dt-url.net/h203qmc).

## Resource Example Usage

```terraform
resource "dynatrace_discovery_feature_flags" "#name#" {
  name            = "#name#"
  type            = "boolean"
  # boolean_value = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the feature
- `type` (String) Possible Values: `Boolean`, `Number`, `String`

### Optional

- `boolean_value` (Boolean) State of boolean feature
- `number_value` (Number) State of numeric feature
- `string_value` (String) State of textual feature

### Read-Only

- `id` (String) The ID of this resource.
 