---
layout: ""
page_title: "dynatrace_openpipeline_v2_system_events_pipelinegroups Resource - terraform-provider-dynatrace"
subcategory: "OpenPipeline V2"
description: |-
  The resource `dynatrace_openpipeline_v2_system_events_pipelinegroups` covers configuration of OpenPipeline for system events pipeline groups
---

# dynatrace_openpipeline_v2_system_events_pipelinegroups (Resource)

-> This resource requires the API token scopes **Read settings** (`settings.read`) and **Write settings** (`settings.write`)

-> This resource requires the OAuth scopes **Read settings** (`settings:objects:read`) and **Write settings** (`settings:objects:write`)

## Limitations
~> **Warning** If a resource is created using an API token or without setting `DYNATRACE_HTTP_OAUTH_PREFERENCE=true` (when both are used), the settings object's owner will remain empty.

An empty owner implies:
- The settings object becomes public, allowing other users with settings permissions to read and modify it.
- Changing the settings object's permissions will have no effect, meaning the `dynatrace_settings_permissions` resource can't alter its access.

When a settings object is created using platform credentials:
- The owner is set to the owner of the OAuth client or platform token.
- By default, the settings object is private; only the owner can read and modify it.
- Access modifiers can be managed using the `dynatrace_settings_permissions` resource.

We recommend using platform credentials to ensure a correct setup.
In case an API token is needed, we recommend setting `DYNATRACE_HTTP_OAUTH_PREFERENCE=true`.

## Dynatrace Documentation

- OpenPipeline - https://docs.dynatrace.com/docs/platform/openpipeline

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_openpipeline_v2_system_events_pipelinegroups` downloads all existing OpenPipeline definitions for system events pipeline groups

The full documentation of the export feature is available [here](https://dt-url.net/h203qmc).

## Resource Example Usage

```terraform
resource "dynatrace_openpipeline_v2_system_events_pipelines" "example1" {
  display_name = "#name#"
  custom_id = "#name#"
  processing {}
  davis {
    processors {
      processor {
        type        = "davis"
        id          = "processor_Create_warning_event_8226"
        description = "Create warning event"
        matcher     = "true"
        davis {
          properties {
            property {
              key = "event.type"
              value = "CUSTOM_ALERT"
            }
            property {
              key = "event.name"
              value = "Warning detected"
            }
            property {
              key = "event.description"
              value = "Warning: {dims:record.summary}"
            }
          }
        }
        enabled = true
      }
    }
  }
  metric_extraction {
    processors {
      processor {
        type        = "counterMetric"
        id          = "processor_Count_warning_events_6392"
        description = "Count warnings"
        matcher     = "true"
        counter_metric {
          metric_key = "warning.count"
          dimensions {
            dimension {
              source_field_name = "dt.cost.costcenter"
            }
            dimension {
              source_field_name = "dt.cost.product"
            }
            dimension {
              source_field_name = "dt.security_context"
            }
            dimension {
              source_field_name      = "record.category"
              destination_field_name = "warning_category"
            }
          }
        }
        enabled = true
      }
      processor {
        type        = "valueMetric"
        id          = "processor_Warning_timeout_1990"
        description = "Warning timeout"
        matcher     = "true"
        value_metric {
          metric_key    = "warning.timeout"
          field         = "recording.timeout_in_min"
          default_value = 60
          dimensions {
            dimension {
              source_field_name = "dt.cost.costcenter"
            }
            dimension {
              source_field_name = "dt.cost.product"
            }
            dimension {
              source_field_name = "dt.security_context"
            }
            dimension {
              source_field_name      = "record.category"
              destination_field_name = "warning_category"
            }
          }
        }
        enabled = true
      }
    }
  }
  security_context {}
  cost_allocation {}
  product_allocation {}
  storage {}
  data_extraction {}
}


resource "dynatrace_openpipeline_v2_system_events_pipelines" "example2" {
  display_name = "#name#-2"
  custom_id = "#name#-2"
  processing {}
  product_allocation {}
  data_extraction {}
  metric_extraction {}
  security_context {}
  storage {}
  cost_allocation {}
  davis {}
}

resource "dynatrace_openpipeline_v2_system_events_pipelinegroups" "example" {
  display_name = "#name#"
  included_pipelines {
    included_pipeline {
      is_target_pipeline_placeholder = false
      pipeline_id = dynatrace_openpipeline_v2_system_events_pipelines.example1.id
      enabled_stages = ["davis", "metricExtraction"]
    }
  }
  targeted_pipelines = [dynatrace_openpipeline_v2_system_events_pipelines.example2.id]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) Display name

### Optional

- `included_pipelines` (Block List, Max: 1) Included Pipelines (see [below for nested schema](#nestedblock--included_pipelines))
- `targeted_pipelines` (Set of String) Pipelines wrapped by this group

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--included_pipelines"></a>
### Nested Schema for `included_pipelines`

Required:

- `included_pipeline` (Block List, Min: 1) (see [below for nested schema](#nestedblock--included_pipelines--included_pipeline))

<a id="nestedblock--included_pipelines--included_pipeline"></a>
### Nested Schema for `included_pipelines.included_pipeline`

Required:

- `is_target_pipeline_placeholder` (Boolean) Placeholder for the wrapped Pipeline

Optional:

- `enabled_stages` (Set of String) Enabled Stages. Possible Values: `costAllocation`, `dataExtraction`, `davis`, `metricExtraction`, `processing`, `productAllocation`, `securityContext`, `smartscapeEdgeExtraction`, `smartscapeNodeExtraction`, `storage`.
- `pipeline_id` (String) Pipeline ID
