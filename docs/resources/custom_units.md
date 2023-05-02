---
layout: ""
page_title: dynatrace_custom_units Resource - terraform-provider-dynatrace"
description: |-
  The resource `dynatrace_custom_units` covers configuration for custom units
---

# dynatrace_custom_units (Resource)

## Dynatrace Documentation

- Settings API - https://www.dynatrace.com/support/help/dynatrace-api/environment-api/settings (schemaId: `builtin:custom-unit`)

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_custom_units` downloads all existing custom unit configuration

The full documentation of the export feature is available [here](https://registry.terraform.io/providers/dynatrace-oss/dynatrace/latest/docs/guides/export-v2).

## Resource Example Usage

```terraform
resource "dynatrace_custom_units" "#name#" {
  name        = "#name#"
  description = "Created by Terraform"
  plural_name = "TerraformUnits"
  symbol      = "T/u"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `description` (String) Unit description should provide additional information about the new unit
- `name` (String) Unit name has to be unique and is used as identifier.
- `plural_name` (String) Unit plural name represent the plural form of the unit name.
- `symbol` (String) Unit symbol has to be unique.

### Read-Only

- `id` (String) The ID of this resource.
 