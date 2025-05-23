---
layout: ""
page_title: dynatrace_processgroup_naming_order Resource - terraform-provider-dynatrace"
subcategory: "Process Group Monitoring"
description: |-
  The resource `dynatrace_processgroup_naming_order` covers defining the order of rules defined for process group naming
---

# dynatrace_processgroup_naming_order (Resource)

-> This resource requires the API token scopes **Read configuration** (`ReadConfig`) and **Write configuration** (`WriteConfig`)

## Dynatrace Documentation

- Process group naming - https://www.dynatrace.com/support/help/how-to-use-dynatrace/process-groups/configuration/pg-naming

- Conditional naming API - https://www.dynatrace.com/support/help/dynatrace-api/configuration-api/conditional-naming

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_processgroup_naming_order` downloads an ordered list of process group naming rule IDs

The full documentation of the export feature is available [here](https://dt-url.net/h203qmc).

## Resource Example Usage

The example below contains configuration for two process group naming rules.
The resource `dynatrace_processgroup_naming_order` defines what order these naming rules should get evaluated in.

In case the Dynatrace environment contains additional process group naming rules that are not managed by Terraform, these
naming rules will end up ranked with lower priority than the ones listed within the `dynatrace_processgroup_naming_order` resource.

-> Be aware of the fact that the resource `dynatrace_processgroup_naming_order` should be treated as a singleton. Maintaining more than one instances of `dynatrace_processgroup_naming_order` within the same Terraform module is not recommended. It will result in non-empty plans.

```terraform
resource "dynatrace_processgroup_naming_order" "process_group_naming_order" {
  naming_rule_ids = [
    dynatrace_processgroup_naming.first.id,
    dynatrace_processgroup_naming.second.id,
  ]  
}

resource "dynatrace_processgroup_naming" "first" {
  name    = "first-one"
  enabled = true
  ...
}

resource "dynatrace_processgroup_naming" "second" {
  name    = "second-one"
  enabled = true
  ...
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `naming_rule_ids` (List of String) The IDs of the naming rules to define the order for

### Read-Only

- `id` (String) The ID of this resource.
 