---
layout: ""
page_title: "dynatrace_process_group_simple_detection Resource - terraform-provider-dynatrace"
subcategory: "Process Group Monitoring"
description: |-
  The resource `dynatrace_process_group_simple_detection` covers configuration for process group simple detection rules
---

# dynatrace_process_group_simple_detection (Resource)

-> This resource requires the API token scopes **Read settings** (`settings.read`) and **Write settings** (`settings.write`)

## Dynatrace Documentation

- Simple detection rules - https://www.dynatrace.com/support/help/platform-modules/infrastructure-monitoring/process-groups/configuration/pg-detection#simple

- Settings API - https://www.dynatrace.com/support/help/dynatrace-api/environment-api/settings (schemaId: `builtin:process-group.simple-detection-rule`)

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_process_group_simple_detection` downloads all existing process group simple detection rules

The full documentation of the export feature is available [here](https://dt-url.net/h203qmc).

## Resource Example Usage

```terraform
resource "dynatrace_process_group_simple_detection" "#name#" {
  enabled             = false
  group_identifier    = "GroupIdentifierExample"
  instance_identifier = "InstanceIdentifierExample"
  process_type        = "PROCESS_TYPE_GO"
  rule_type           = "prop"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)
- `group_identifier` (String) If Dynatrace detects this property at startup of a process, it will use its value to identify process groups more granular.
- `instance_identifier` (String) Use a variable to identify instances within a process group.

The type of variable is the same as selected in 'Property source'.
- `rule_type` (String) Possible Values: `Prop`, `Env`

### Optional

- `insert_after` (String) Because this resource allows for ordering you may specify the ID of the resource instance that comes before this instance regarding order. If not specified when creating the setting will be added to the end of the list. If not specified during update the order will remain untouched
- `process_type` (String) Note: Not all types can be detected at startup.

### Read-Only

- `id` (String) The ID of this resource.
 