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
- Pipeline groups - https://docs.dynatrace.com/docs/shortlink/openpipeline-pipeline-groups

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_openpipeline_v2_system_events_pipelinegroups` downloads all existing OpenPipeline definitions for system events pipeline groups

The full documentation of the export feature is available [here](https://dt-url.net/h203qmc).

## Resource Example Usage

```terraform
resource "dynatrace_openpipeline_v2_system_events_pipelines" "example1" {
  display_name = "#name#"
  custom_id = "#name#"
  group_role = "compositionPipeline"
  routing = "notRoutable"
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
}


resource "dynatrace_openpipeline_v2_system_events_pipelines" "example2" {
  display_name = "#name#-2"
  custom_id = "#name#-2"
  group_role = "memberPipeline"
}

resource "dynatrace_openpipeline_v2_system_events_pipelinegroups" "example" {
  display_name = "#name#"
  composition {
    pipeline_group_composition {
      is_pipeline_placeholder = true
    }
    pipeline_group_composition {
      is_pipeline_placeholder = false
      stages {
        type = "include"
        include = ["davis", "metricExtraction"]
      }
      pipeline_id = dynatrace_openpipeline_v2_system_events_pipelines.example1.id
    }
  }
  member_stages {
    include = ["davis", "metricExtraction"]
    type = "include"
  }
  member_pipelines = [dynatrace_openpipeline_v2_system_events_pipelines.example2.id]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) Display name
- `member_stages` (Block List, Min: 1, Max: 1) stage configuration of the member pipelines (see [below for nested schema](#nestedblock--member_stages))

### Optional

- `composition` (Block List, Max: 1) Composition (see [below for nested schema](#nestedblock--composition))
- `member_pipelines` (Set of String) Pipelines wrapped by this group

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--member_stages"></a>
### Nested Schema for `member_stages`

Required:

- `type` (String) Stage configuration type. Possible values: `exclude`, `include`, `includeAll`

Optional:

- `exclude` (Set of String) exclude stages. Possible values: `costAllocation`, `dataExtraction`, `davis`, `metricExtraction`, `processing`, `productAllocation`, `securityContext`, `smartscapeEdgeExtraction`, `smartscapeNodeExtraction`, `storage`
- `include` (Set of String) include stages. Possible values: `costAllocation`, `dataExtraction`, `davis`, `metricExtraction`, `processing`, `productAllocation`, `securityContext`, `smartscapeEdgeExtraction`, `smartscapeNodeExtraction`, `storage`


<a id="nestedblock--composition"></a>
### Nested Schema for `composition`

Required:

- `pipeline_group_composition` (Block List, Min: 1) (see [below for nested schema](#nestedblock--composition--pipeline_group_composition))

<a id="nestedblock--composition--pipeline_group_composition"></a>
### Nested Schema for `composition.pipeline_group_composition`

Required:

- `is_pipeline_placeholder` (Boolean) Placeholder for the wrapped pipeline

Optional:

- `pipeline_id` (String) Pipeline ID
- `stages` (Block List, Max: 1) stage configuration for this pipelines (see [below for nested schema](#nestedblock--composition--pipeline_group_composition--stages))

<a id="nestedblock--composition--pipeline_group_composition--stages"></a>
### Nested Schema for `composition.pipeline_group_composition.stages`

Required:

- `type` (String) Stage configuration type. Possible values: `exclude`, `include`, `includeAll`

Optional:

- `exclude` (Set of String) exclude stages. Possible values: `costAllocation`, `dataExtraction`, `davis`, `metricExtraction`, `processing`, `productAllocation`, `securityContext`, `smartscapeEdgeExtraction`, `smartscapeNodeExtraction`, `storage`
- `include` (Set of String) include stages. Possible values: `costAllocation`, `dataExtraction`, `davis`, `metricExtraction`, `processing`, `productAllocation`, `securityContext`, `smartscapeEdgeExtraction`, `smartscapeNodeExtraction`, `storage`
