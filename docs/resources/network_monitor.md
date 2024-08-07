---
layout: ""
page_title: dynatrace_network_monitor Resource - terraform-provider-dynatrace"
subcategory: "Network Availability Monitors"
description: |-
  The resource `dynatrace_network_monitor` covers configuration for network availability monitors
---

# dynatrace_network_monitor (Resource)

-> This resource requires the API token scopes **Read settings** (`settings.read`) and **Write settings** (`settings.write`)

## Dynatrace Documentation

- Network availability monitors - https://docs.dynatrace.com/docs/platform-modules/digital-experience/synthetic-monitoring/general-information/network-availability-monitors

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_network_monitor` downloads all existing network monitor configuration

The full documentation of the export feature is available [here](https://dt-url.net/h203qmc).

## Resource Example Usage

```terraform
resource "dynatrace_network_monitor" "DNS_Test" {
  name          = "DNS Test"
  description   = "This is an example DNS test"
  type          = "MULTI_PROTOCOL"
  enabled       = false
  frequency_min = 15
  locations     = [ "SYNTHETIC_LOCATION-39F97465BE7BF644" ]
  outage_handling {
    global_consecutive_outage_count_threshold = 1
    global_outages                            = true
  }
  performance_thresholds {
    enabled = true
    thresholds {
      threshold {
        aggregation        = "AVG"
        dealerting_samples = 5
        samples            = 5
        step_index         = 0
        threshold          = 100
        violating_samples  = 3
      }
    }
  }
  steps {
    step {
      name         = "DNS Test"
      request_type = "DNS"
      target_list  = [ "google.com", "yahoo.com" ]
      properties = {
        "DNS_RECORD_TYPES"  = "A"
        "EXECUTION_TIMEOUT" = "PT2S"
      }
      constraints {
        constraint {
          type = "SUCCESS_RATE_PERCENT"
          properties = {
            "value"    = "90"
            "operator" = ">="
          }
        }
      }
      request_configurations {
        request_configuration {
          constraints {
            constraint {
              type = "DNS_STATUS_CODE"
              properties = {
                "operator"   = "="
                "statusCode" = "0"
              }
            }
          }
        }
      }
    }
  }
  tags {
    tag {
      context = "CONTEXTLESS"
      key     = "Key1"
      source  = "USER"
      value   = "Value1"
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `locations` (Set of String) The locations to which the monitor is assigned
- `name` (String) Name of the monitor
- `steps` (Block List, Min: 1) The steps of the monitor (see [below for nested schema](#nestedblock--steps))
- `type` (String) Type of the monitor, possible values: `MULTI_PROTOCOL`

### Optional

- `description` (String) Description of the monitor
- `enabled` (Boolean) If true, the monitor is enabled
- `frequency_min` (Number) Frequency of the monitor, in minutes
- `outage_handling` (Block List, Max: 1) Outage handling configuration (see [below for nested schema](#nestedblock--outage_handling))
- `performance_thresholds` (Block List, Max: 1) Performance thresholds configuration (see [below for nested schema](#nestedblock--performance_thresholds))
- `tags` (Block List) A set of tags assigned to the monitor.

You can specify only the value of the tag here and the CONTEXTLESS context and source 'USER' will be added automatically. But preferred option is usage of SyntheticTagWithSourceDto model. (see [below for nested schema](#nestedblock--tags))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--steps"></a>
### Nested Schema for `steps`

Required:

- `step` (Block List, Min: 1) The step of a network availability monitor (see [below for nested schema](#nestedblock--steps--step))

<a id="nestedblock--steps--step"></a>
### Nested Schema for `steps.step`

Required:

- `constraints` (Block List, Min: 1) The list of constraints which apply to all requests in the step (see [below for nested schema](#nestedblock--steps--step--constraints))
- `name` (String) Step name
- `properties` (Map of String) Key/value pairs of properties which apply to all requests in the step
- `request_type` (String) Request type, possible values `ICMP`, `TCP`, `DNS`
- `target_list` (Set of String) Target list

Optional:

- `request_configurations` (Block List) Request configurations (see [below for nested schema](#nestedblock--steps--step--request_configurations))
- `target_filter` (String) Target filter

<a id="nestedblock--steps--step--constraints"></a>
### Nested Schema for `steps.step.constraints`

Required:

- `constraint` (Block Set, Min: 1) The network availability monitor constraint (see [below for nested schema](#nestedblock--steps--step--constraints--constraint))

<a id="nestedblock--steps--step--constraints--constraint"></a>
### Nested Schema for `steps.step.constraints.constraint`

Required:

- `properties` (Map of String) Key/value pairs of constraint properties
- `type` (String) Constraint type



<a id="nestedblock--steps--step--request_configurations"></a>
### Nested Schema for `steps.step.request_configurations`

Required:

- `request_configuration` (Block Set, Min: 1) The configuration of a network availability monitor request (see [below for nested schema](#nestedblock--steps--step--request_configurations--request_configuration))

<a id="nestedblock--steps--step--request_configurations--request_configuration"></a>
### Nested Schema for `steps.step.request_configurations.request_configuration`

Optional:

- `constraints` (Block List) Request constraints (see [below for nested schema](#nestedblock--steps--step--request_configurations--request_configuration--constraints))

<a id="nestedblock--steps--step--request_configurations--request_configuration--constraints"></a>
### Nested Schema for `steps.step.request_configurations.request_configuration.constraints`

Required:

- `constraint` (Block Set, Min: 1) The network availability monitor constraint (see [below for nested schema](#nestedblock--steps--step--request_configurations--request_configuration--constraints--constraint))

<a id="nestedblock--steps--step--request_configurations--request_configuration--constraints--constraint"></a>
### Nested Schema for `steps.step.request_configurations.request_configuration.constraints.constraint`

Required:

- `properties` (Map of String) Key/value pairs of constraint properties
- `type` (String) Constraint type







<a id="nestedblock--outage_handling"></a>
### Nested Schema for `outage_handling`

Optional:

- `global_consecutive_outage_count_threshold` (Number) Number of consecutive failures for all locations
- `global_outages` (Boolean) Generate a problem and send an alert when the monitor is unavailable at all configured locations
- `local_consecutive_outage_count_threshold` (Number) Number of consecutive failures
- `local_location_outage_count_threshold` (Number) Number of failing locations
- `local_outages` (Boolean) Generate a problem and send an alert when the monitor is unavailable for one or more consecutive runs at any location


<a id="nestedblock--performance_thresholds"></a>
### Nested Schema for `performance_thresholds`

Optional:

- `enabled` (Boolean) Performance threshold is enabled (true) or disabled (false)
- `thresholds` (Block List, Max: 1) The list of performance threshold rules (see [below for nested schema](#nestedblock--performance_thresholds--thresholds))

<a id="nestedblock--performance_thresholds--thresholds"></a>
### Nested Schema for `performance_thresholds.thresholds`

Optional:

- `threshold` (Block Set) The list of performance threshold rules (see [below for nested schema](#nestedblock--performance_thresholds--thresholds--threshold))

<a id="nestedblock--performance_thresholds--thresholds--threshold"></a>
### Nested Schema for `performance_thresholds.thresholds.threshold`

Optional:

- `aggregation` (String) Aggregation type, possible values: `AVG`, `MAX`, `MIN`
- `dealerting_samples` (Number) Number of most recent non-violating request executions that closes the problem
- `samples` (Number) Number of request executions in analyzed sliding window (sliding window size)
- `step_index` (Number) Specify the step's index to which a threshold applies
- `threshold` (Number) Notify if monitor request takes longer than X milliseconds to execute
- `violating_samples` (Number) Number of violating request executions in analyzed sliding window




<a id="nestedblock--tags"></a>
### Nested Schema for `tags`

Optional:

- `tag` (Block Set) Tag with source of a Dynatrace entity. (see [below for nested schema](#nestedblock--tags--tag))

<a id="nestedblock--tags--tag"></a>
### Nested Schema for `tags.tag`

Required:

- `key` (String) The key of the tag

Optional:

- `context` (String) The origin of the tag, such as AWS or Cloud Foundry.

Custom tags use the CONTEXTLESS value
- `source` (String) The source of the tag, possible values: `AUTO`, `RULE_BASED` or `USER`
- `value` (String) The value of the tag
 