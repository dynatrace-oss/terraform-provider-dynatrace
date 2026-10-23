---
layout: ""
page_title: "dynatrace_openpipeline_v2_smartscape_events_pipelinegroups Resource - terraform-provider-dynatrace"
subcategory: "OpenPipeline V2"
description: |-
  The resource `dynatrace_openpipeline_v2_smartscape_events_pipelinegroups` covers configuration of OpenPipeline for smartscape events pipeline groups
---

# dynatrace_openpipeline_v2_smartscape_events_pipelinegroups (Resource)

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

- `terraform-provider-dynatrace -export dynatrace_openpipeline_v2_smartscape_events_pipelinegroups` downloads all existing OpenPipeline definitions for smartscape events pipeline groups

The full documentation of the export feature is available [here](https://dt-url.net/h203qmc).

## Resource Example Usage

```terraform
resource "dynatrace_openpipeline_v2_smartscape_events_pipelines" "example1" {
  display_name = "#name#"
  custom_id    = "#name#"
  group_role   = "compositionPipeline"
  routing      = "notRoutable"
  processing {
    processors {
      processor {
        type        = "drop"
        id          = "processor_Drop_unnecessary_records_3802"
        description = "Drop unnecessary records"
        matcher     = "not matchesPhrase(name, \"Warning\")"
        enabled     = true
      }
      processor {
        type        = "fieldsAdd"
        id          = "processor_Add_warning_flag_5434"
        description = "Add warning flag"
        matcher     = "matchesPhrase(name, \"Warning\")"
        sample_data = "{\n  \"name\": \"Warning record\" \n}"
        fields_add {
          fields {
            field {
              name       = "custom.is_warning"
              value      = "true"
              value_type = "boolean"
            }
          }
        }
        enabled = true
      }
      processor {
        type        = "fieldsRemove"
        id          = "processor_Remove_details_field_8539"
        description = "Remove details field"
        sample_data = "{\n  \"name\": \"Warning\",\n  \"custom.details\": \"some record details\"\n}"
        matcher     = "isNotNull(custom.details)"
        fields_remove {
          fields = ["custom.details"]
        }
        enabled = true
      }
      processor {
        type        = "fieldsRename"
        id          = "processor_Rename_title_to_summary_8530"
        description = "Rename title to summary"
        sample_data = "{\n  \"custom.title\": \"Warning\"\n}"
        matcher     = "true"
        fields_rename {
          fields {
            field {
              from_name = "custom.title"
              to_name   = "custom.summary"
            }
          }
        }
        enabled = true
      }
      processor {
        type        = "dql"
        id          = "processor_Combine_summary_and_source_8808"
        description = "Combine summary and source"
        sample_data = "{\n  \"custom.summary\": \"Warning\",\n  \"custom.source\": \"Request failed\"\n}"
        matcher     = "true"
        dql {
          script = "fieldsAdd custom.combined = concat(custom.summary, \" - \", custom.source)"
        }
        enabled = true
      }
    }
  }
}


resource "dynatrace_openpipeline_v2_smartscape_events_pipelines" "example2" {
  display_name = "#name#-2"
  custom_id    = "#name#-2"
  group_role   = "memberPipeline"
}


resource "dynatrace_openpipeline_v2_smartscape_events_pipelinegroups" "example" {
  display_name = "#name#"
  composition {
    pipeline_group_composition {
      is_pipeline_placeholder = true
    }
    pipeline_group_composition {
      is_pipeline_placeholder = false
      stages {
        type    = "include"
        include = ["processing"]
      }
      pipeline_id = dynatrace_openpipeline_v2_smartscape_events_pipelines.example1.id
    }
  }
  member_stages {
    include = ["processing"]
    type    = "include"
  }
  member_pipelines = [dynatrace_openpipeline_v2_smartscape_events_pipelines.example2.id]
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

- `exclude` (Set of String) exclude stages. Possible values: `costAllocation`, `dataExtraction`, `davis`, `metricExtraction`, `processing`, `productAllocation`, `security`, `securityContext`, `smartscapeEdgeExtraction`, `smartscapeNodeExtraction`, `storage`
- `include` (Set of String) include stages. Possible values: `costAllocation`, `dataExtraction`, `davis`, `metricExtraction`, `processing`, `productAllocation`, `security`, `securityContext`, `smartscapeEdgeExtraction`, `smartscapeNodeExtraction`, `storage`


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

- `exclude` (Set of String) exclude stages. Possible values: `costAllocation`, `dataExtraction`, `davis`, `metricExtraction`, `processing`, `productAllocation`, `security`, `securityContext`, `smartscapeEdgeExtraction`, `smartscapeNodeExtraction`, `storage`
- `include` (Set of String) include stages. Possible values: `costAllocation`, `dataExtraction`, `davis`, `metricExtraction`, `processing`, `productAllocation`, `security`, `securityContext`, `smartscapeEdgeExtraction`, `smartscapeNodeExtraction`, `storage`
