---
layout: ""
page_title: "dynatrace_openpipeline_v2_events_security_pipelinegroups Resource - terraform-provider-dynatrace"
subcategory: "OpenPipeline V2"
description: |-
  The resource `dynatrace_openpipeline_v2_events_security_pipelinegroups` covers configuration of OpenPipeline for events security pipeline groups
---

# dynatrace_openpipeline_v2_events_security_pipelinegroups (Resource)

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

- `terraform-provider-dynatrace -export dynatrace_openpipeline_v2_events_security_pipelinegroups` downloads all existing OpenPipeline definitions for events security pipeline groups

The full documentation of the export feature is available [here](https://dt-url.net/h203qmc).

## Resource Example Usage

```terraform
resource "dynatrace_openpipeline_v2_events_security_pipelines" "example1" {
  display_name = "#name#"
  custom_id = "#name#"
  group_role = "compositionPipeline"
  routing = "notRoutable"
  processing {
    processors {
      processor {
        type        = "drop"
        id          = "processor_Drop_unnecessary_records_3802"
        description = "Drop unnecessary records"
        matcher     = "not matchesPhrase(record.name, \"Warning\")"
        enabled     = true
      }
      processor {
        type        = "fieldsAdd"
        id          = "processor_Add_warning_flag_5434"
        description = "Add warning flag"
        matcher     = "matchesPhrase(record.name, \"Warning\")"
        sample_data = "{\n  \"record.name\": \"Warning record\" \n}"
        fields_add {
          fields {
            field {
              name  = "is_warning"
              value = "true"
            }
          }
        }
        enabled = true
      }
      processor {
        type        = "fieldsRemove"
        id          = "processor_Remove_details_field_8539"
        description = "Remove details field"
        sample_data = "{\n  \"record.name\": \"Warning\",\n  \"record.details\": \"some record details\"\n}"
        matcher     = "isNotNull(record.details)"
        fields_remove {
          fields = ["record.details"]
        }
        enabled     = true
      }
      processor {
        type        = "fieldsRename"
        id          = "processor_Rename_name_to_title_8530"
        description = "Rename name to title"
        sample_data = "{\n  \"record.name\": \"Warning\"\n}"
        matcher     = "true"
        fields_rename {
          fields {
            field {
              from_name = "record.name"
              to_name   = "record.title"
            }
          }
        }
        enabled     = true
      }
      processor {
        type        = "dql"
        id          = "processor_Combine_title_and_summary_to_name_8808"
        description = "Combine title and summary to name"
        sample_data = "{\n  \"record.title\": \"Warning\",\n  \"record.summary\": \"Request failed\"\n}"
        matcher     = "true"
        dql {
          script = "fieldsAdd record.name = concat(record.title, \" - \", record.summary)"
        }
        enabled     = true
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


resource "dynatrace_openpipeline_v2_events_security_pipelines" "example2" {
  display_name = "#name#-2"
  custom_id = "#name#-2"
  group_role = "memberPipeline"
}

resource "dynatrace_openpipeline_v2_events_sdlc_pipelinegroups" "example" {
  display_name = "#name#"
  composition {
    pipeline_group_composition {
      is_pipeline_placeholder = true
    }
    pipeline_group_composition {
      is_pipeline_placeholder = false
      stages {
        type = "include"
        include = ["processing", "metricExtraction"]
      }
      pipeline_id = dynatrace_openpipeline_v2_events_security_pipelines.example1.id
    }
  }
  member_stages {
    include = ["processing", "metricExtraction"]
    type = "include"
  }
  member_pipelines = [dynatrace_openpipeline_v2_events_security_pipelines.example2.id]
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
