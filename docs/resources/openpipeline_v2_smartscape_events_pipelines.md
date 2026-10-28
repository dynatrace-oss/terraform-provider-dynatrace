---
layout: ""
page_title: "dynatrace_openpipeline_v2_smartscape_events_pipelines Resource - terraform-provider-dynatrace"
subcategory: "OpenPipeline V2"
description: |-
  The resource `dynatrace_openpipeline_v2_smartscape_events_pipelines` covers configuration of OpenPipeline for smartscape events pipelines
---

# dynatrace_openpipeline_v2_smartscape_events_pipelines (Resource)

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

- `terraform-provider-dynatrace -export dynatrace_openpipeline_v2_smartscape_events_pipelines` downloads all existing OpenPipeline definitions for smartscape events pipelines

The full documentation of the export feature is available [here](https://dt-url.net/h203qmc).

## Resource Example Usage

```terraform
resource "dynatrace_openpipeline_v2_smartscape_events_pipelines" "max-pipeline" {
  display_name = "Warning pipeline"
  custom_id    = "pipeline_Warning_pipeline_2773_tf_#name#"
  metadata_list {
    metadata {
      entry_key   = "environment"
      entry_value = "production"
    }
  }
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
  security_context {
    processors {
      processor {
        type        = "securityContext"
        id          = "processor_Use_dt.security_context_if_set_1080"
        description = "Use dt.security_context if set"
        matcher     = "isNotNull(dt.security_context)"
        security_context {
          value {
            type = "field"
            field {
              source_field_name = "dt.security_context"
            }
          }
        }
        enabled = true
      }
      processor {
        type        = "securityContext"
        id          = "processor_Assign_warnings_to_ACME_teams_if_no_context_set_5465"
        description = "Assign warnings to ACME teams if no context set"
        matcher     = "isNull(dt.security_context)"
        security_context {
          value {
            type                 = "multiValueConstant"
            multi_value_constant = ["ACME1", "ACME2"]
          }
        }
        enabled = true
      }
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `custom_id` (String) Custom pipeline id
- `display_name` (String) Display name

### Optional

- `cost_allocation` (Block List, Max: 1) Cost allocation stage (see [below for nested schema](#nestedblock--cost_allocation))
- `data_extraction` (Block List, Max: 1) Data extraction stage (see [below for nested schema](#nestedblock--data_extraction))
- `davis` (Block List, Max: 1) Davis event extraction stage (see [below for nested schema](#nestedblock--davis))
- `group_role` (String) Group role. Possible values: `basePipeline`, `compositionPipeline`, `memberPipeline`
- `metadata_list` (Block List, Max: 1) Pipeline metadata list (see [below for nested schema](#nestedblock--metadata_list))
- `metric_extraction` (Block List, Max: 1) Metrics extraction stage (see [below for nested schema](#nestedblock--metric_extraction))
- `processing` (Block List, Max: 1) Processing stage (see [below for nested schema](#nestedblock--processing))
- `product_allocation` (Block List, Max: 1) Product allocation stage (see [below for nested schema](#nestedblock--product_allocation))
- `routing` (String) Routing. Possible values: `notRoutable`, `routable`
- `security` (Block List, Max: 1) Security stage (see [below for nested schema](#nestedblock--security))
- `security_context` (Block List, Max: 1) Security context stage (see [below for nested schema](#nestedblock--security_context))
- `smartscape_edge_extraction` (Block List, Max: 1) Smartscape edge extraction stage (see [below for nested schema](#nestedblock--smartscape_edge_extraction))
- `smartscape_node_extraction` (Block List, Max: 1) Smartscape node extraction stage (see [below for nested schema](#nestedblock--smartscape_node_extraction))
- `storage` (Block List, Max: 1) Storage stage (see [below for nested schema](#nestedblock--storage))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--cost_allocation"></a>
### Nested Schema for `cost_allocation`

Optional:

- `processors` (Block List, Max: 1) Processors of stage (see [below for nested schema](#nestedblock--cost_allocation--processors))

<a id="nestedblock--cost_allocation--processors"></a>
### Nested Schema for `cost_allocation.processors`

Required:

- `processor` (Block List, Min: 1) (see [below for nested schema](#nestedblock--cost_allocation--processors--processor))

<a id="nestedblock--cost_allocation--processors--processor"></a>
### Nested Schema for `cost_allocation.processors.processor`

Required:

- `description` (String) No documentation available
- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)
- `id` (String) Processor identifier
- `type` (String) Processor type. Possible values: `dql`, `drop`, `fieldsAdd`, `fieldsRemove`, `fieldsRename`, `securityContext`

Optional:

- `dql` (Block List, Max: 1) DQL processor attributes (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--dql))
- `fields_add` (Block List, Max: 1) Fields add processor attributes (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--fields_add))
- `fields_remove` (Block List, Max: 1) Fields remove processor attributes (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--fields_remove))
- `fields_rename` (Block List, Max: 1) Fields rename processor attributes (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--fields_rename))
- `matcher` (String) [See our documentation](https://dt-url.net/bp234rv)
- `sample_data` (String) Sample data
- `security_context` (Block List, Max: 1) Security context processor attributes (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--security_context))

<a id="nestedblock--cost_allocation--processors--processor--dql"></a>
### Nested Schema for `cost_allocation.processors.processor.dql`

Required:

- `script` (String) DQL script


<a id="nestedblock--cost_allocation--processors--processor--fields_add"></a>
### Nested Schema for `cost_allocation.processors.processor.fields_add`

Required:

- `fields` (Block List, Min: 1, Max: 1) Fields to Add (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--fields_add--fields))

<a id="nestedblock--cost_allocation--processors--processor--fields_add--fields"></a>
### Nested Schema for `cost_allocation.processors.processor.fields_add.fields`

Required:

- `field` (Block List, Min: 1) (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--fields_add--fields--field))

<a id="nestedblock--cost_allocation--processors--processor--fields_add--fields--field"></a>
### Nested Schema for `cost_allocation.processors.processor.fields_add.fields.field`

Required:

- `name` (String) Fields's name
- `value` (String) Field's value

Optional:

- `value_type` (String) Field value type. Possible values: `boolean`, `double`, `long`, `string`




<a id="nestedblock--cost_allocation--processors--processor--fields_remove"></a>
### Nested Schema for `cost_allocation.processors.processor.fields_remove`

Required:

- `fields` (Set of String) Fields to remove


<a id="nestedblock--cost_allocation--processors--processor--fields_rename"></a>
### Nested Schema for `cost_allocation.processors.processor.fields_rename`

Required:

- `fields` (Block List, Min: 1, Max: 1) Fields to rename (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--fields_rename--fields))

<a id="nestedblock--cost_allocation--processors--processor--fields_rename--fields"></a>
### Nested Schema for `cost_allocation.processors.processor.fields_rename.fields`

Required:

- `field` (Block List, Min: 1) (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--fields_rename--fields--field))

<a id="nestedblock--cost_allocation--processors--processor--fields_rename--fields--field"></a>
### Nested Schema for `cost_allocation.processors.processor.fields_rename.fields.field`

Required:

- `from_name` (String) Fields's name
- `to_name` (String) New field's name




<a id="nestedblock--cost_allocation--processors--processor--security_context"></a>
### Nested Schema for `cost_allocation.processors.processor.security_context`

Required:

- `value` (Block List, Min: 1, Max: 1) Security context value assignment (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--security_context--value))

<a id="nestedblock--cost_allocation--processors--processor--security_context--value"></a>
### Nested Schema for `cost_allocation.processors.processor.security_context.value`

Required:

- `type` (String) Type of value assignment. Possible values: `constant`, `field`, `multiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--security_context--value--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--cost_allocation--processors--processor--security_context--value--field"></a>
### Nested Schema for `cost_allocation.processors.processor.security_context.value.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value







<a id="nestedblock--data_extraction"></a>
### Nested Schema for `data_extraction`

Optional:

- `processors` (Block List, Max: 1) Processors of stage (see [below for nested schema](#nestedblock--data_extraction--processors))

<a id="nestedblock--data_extraction--processors"></a>
### Nested Schema for `data_extraction.processors`

Required:

- `processor` (Block List, Min: 1) (see [below for nested schema](#nestedblock--data_extraction--processors--processor))

<a id="nestedblock--data_extraction--processors--processor"></a>
### Nested Schema for `data_extraction.processors.processor`

Required:

- `description` (String) No documentation available
- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)
- `id` (String) Processor identifier
- `type` (String) Processor type. Possible values: `dql`, `drop`, `fieldsAdd`, `fieldsRemove`, `fieldsRename`, `securityContext`

Optional:

- `dql` (Block List, Max: 1) DQL processor attributes (see [below for nested schema](#nestedblock--data_extraction--processors--processor--dql))
- `fields_add` (Block List, Max: 1) Fields add processor attributes (see [below for nested schema](#nestedblock--data_extraction--processors--processor--fields_add))
- `fields_remove` (Block List, Max: 1) Fields remove processor attributes (see [below for nested schema](#nestedblock--data_extraction--processors--processor--fields_remove))
- `fields_rename` (Block List, Max: 1) Fields rename processor attributes (see [below for nested schema](#nestedblock--data_extraction--processors--processor--fields_rename))
- `matcher` (String) [See our documentation](https://dt-url.net/bp234rv)
- `sample_data` (String) Sample data
- `security_context` (Block List, Max: 1) Security context processor attributes (see [below for nested schema](#nestedblock--data_extraction--processors--processor--security_context))

<a id="nestedblock--data_extraction--processors--processor--dql"></a>
### Nested Schema for `data_extraction.processors.processor.dql`

Required:

- `script` (String) DQL script


<a id="nestedblock--data_extraction--processors--processor--fields_add"></a>
### Nested Schema for `data_extraction.processors.processor.fields_add`

Required:

- `fields` (Block List, Min: 1, Max: 1) Fields to Add (see [below for nested schema](#nestedblock--data_extraction--processors--processor--fields_add--fields))

<a id="nestedblock--data_extraction--processors--processor--fields_add--fields"></a>
### Nested Schema for `data_extraction.processors.processor.fields_add.fields`

Required:

- `field` (Block List, Min: 1) (see [below for nested schema](#nestedblock--data_extraction--processors--processor--fields_add--fields--field))

<a id="nestedblock--data_extraction--processors--processor--fields_add--fields--field"></a>
### Nested Schema for `data_extraction.processors.processor.fields_add.fields.field`

Required:

- `name` (String) Fields's name
- `value` (String) Field's value

Optional:

- `value_type` (String) Field value type. Possible values: `boolean`, `double`, `long`, `string`




<a id="nestedblock--data_extraction--processors--processor--fields_remove"></a>
### Nested Schema for `data_extraction.processors.processor.fields_remove`

Required:

- `fields` (Set of String) Fields to remove


<a id="nestedblock--data_extraction--processors--processor--fields_rename"></a>
### Nested Schema for `data_extraction.processors.processor.fields_rename`

Required:

- `fields` (Block List, Min: 1, Max: 1) Fields to rename (see [below for nested schema](#nestedblock--data_extraction--processors--processor--fields_rename--fields))

<a id="nestedblock--data_extraction--processors--processor--fields_rename--fields"></a>
### Nested Schema for `data_extraction.processors.processor.fields_rename.fields`

Required:

- `field` (Block List, Min: 1) (see [below for nested schema](#nestedblock--data_extraction--processors--processor--fields_rename--fields--field))

<a id="nestedblock--data_extraction--processors--processor--fields_rename--fields--field"></a>
### Nested Schema for `data_extraction.processors.processor.fields_rename.fields.field`

Required:

- `from_name` (String) Fields's name
- `to_name` (String) New field's name




<a id="nestedblock--data_extraction--processors--processor--security_context"></a>
### Nested Schema for `data_extraction.processors.processor.security_context`

Required:

- `value` (Block List, Min: 1, Max: 1) Security context value assignment (see [below for nested schema](#nestedblock--data_extraction--processors--processor--security_context--value))

<a id="nestedblock--data_extraction--processors--processor--security_context--value"></a>
### Nested Schema for `data_extraction.processors.processor.security_context.value`

Required:

- `type` (String) Type of value assignment. Possible values: `constant`, `field`, `multiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--data_extraction--processors--processor--security_context--value--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--data_extraction--processors--processor--security_context--value--field"></a>
### Nested Schema for `data_extraction.processors.processor.security_context.value.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value







<a id="nestedblock--davis"></a>
### Nested Schema for `davis`

Optional:

- `processors` (Block List, Max: 1) Processors of stage (see [below for nested schema](#nestedblock--davis--processors))

<a id="nestedblock--davis--processors"></a>
### Nested Schema for `davis.processors`

Required:

- `processor` (Block List, Min: 1) (see [below for nested schema](#nestedblock--davis--processors--processor))

<a id="nestedblock--davis--processors--processor"></a>
### Nested Schema for `davis.processors.processor`

Required:

- `description` (String) No documentation available
- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)
- `id` (String) Processor identifier
- `type` (String) Processor type. Possible values: `dql`, `drop`, `fieldsAdd`, `fieldsRemove`, `fieldsRename`, `securityContext`

Optional:

- `dql` (Block List, Max: 1) DQL processor attributes (see [below for nested schema](#nestedblock--davis--processors--processor--dql))
- `fields_add` (Block List, Max: 1) Fields add processor attributes (see [below for nested schema](#nestedblock--davis--processors--processor--fields_add))
- `fields_remove` (Block List, Max: 1) Fields remove processor attributes (see [below for nested schema](#nestedblock--davis--processors--processor--fields_remove))
- `fields_rename` (Block List, Max: 1) Fields rename processor attributes (see [below for nested schema](#nestedblock--davis--processors--processor--fields_rename))
- `matcher` (String) [See our documentation](https://dt-url.net/bp234rv)
- `sample_data` (String) Sample data
- `security_context` (Block List, Max: 1) Security context processor attributes (see [below for nested schema](#nestedblock--davis--processors--processor--security_context))

<a id="nestedblock--davis--processors--processor--dql"></a>
### Nested Schema for `davis.processors.processor.dql`

Required:

- `script` (String) DQL script


<a id="nestedblock--davis--processors--processor--fields_add"></a>
### Nested Schema for `davis.processors.processor.fields_add`

Required:

- `fields` (Block List, Min: 1, Max: 1) Fields to Add (see [below for nested schema](#nestedblock--davis--processors--processor--fields_add--fields))

<a id="nestedblock--davis--processors--processor--fields_add--fields"></a>
### Nested Schema for `davis.processors.processor.fields_add.fields`

Required:

- `field` (Block List, Min: 1) (see [below for nested schema](#nestedblock--davis--processors--processor--fields_add--fields--field))

<a id="nestedblock--davis--processors--processor--fields_add--fields--field"></a>
### Nested Schema for `davis.processors.processor.fields_add.fields.field`

Required:

- `name` (String) Fields's name
- `value` (String) Field's value

Optional:

- `value_type` (String) Field value type. Possible values: `boolean`, `double`, `long`, `string`




<a id="nestedblock--davis--processors--processor--fields_remove"></a>
### Nested Schema for `davis.processors.processor.fields_remove`

Required:

- `fields` (Set of String) Fields to remove


<a id="nestedblock--davis--processors--processor--fields_rename"></a>
### Nested Schema for `davis.processors.processor.fields_rename`

Required:

- `fields` (Block List, Min: 1, Max: 1) Fields to rename (see [below for nested schema](#nestedblock--davis--processors--processor--fields_rename--fields))

<a id="nestedblock--davis--processors--processor--fields_rename--fields"></a>
### Nested Schema for `davis.processors.processor.fields_rename.fields`

Required:

- `field` (Block List, Min: 1) (see [below for nested schema](#nestedblock--davis--processors--processor--fields_rename--fields--field))

<a id="nestedblock--davis--processors--processor--fields_rename--fields--field"></a>
### Nested Schema for `davis.processors.processor.fields_rename.fields.field`

Required:

- `from_name` (String) Fields's name
- `to_name` (String) New field's name




<a id="nestedblock--davis--processors--processor--security_context"></a>
### Nested Schema for `davis.processors.processor.security_context`

Required:

- `value` (Block List, Min: 1, Max: 1) Security context value assignment (see [below for nested schema](#nestedblock--davis--processors--processor--security_context--value))

<a id="nestedblock--davis--processors--processor--security_context--value"></a>
### Nested Schema for `davis.processors.processor.security_context.value`

Required:

- `type` (String) Type of value assignment. Possible values: `constant`, `field`, `multiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--davis--processors--processor--security_context--value--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--davis--processors--processor--security_context--value--field"></a>
### Nested Schema for `davis.processors.processor.security_context.value.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value







<a id="nestedblock--metadata_list"></a>
### Nested Schema for `metadata_list`

Required:

- `metadata` (Block List, Min: 1) (see [below for nested schema](#nestedblock--metadata_list--metadata))

<a id="nestedblock--metadata_list--metadata"></a>
### Nested Schema for `metadata_list.metadata`

Required:

- `entry_key` (String) Metadata entry key

Optional:

- `entry_value` (String) Metadata entry value



<a id="nestedblock--metric_extraction"></a>
### Nested Schema for `metric_extraction`

Optional:

- `processors` (Block List, Max: 1) Processors of stage (see [below for nested schema](#nestedblock--metric_extraction--processors))

<a id="nestedblock--metric_extraction--processors"></a>
### Nested Schema for `metric_extraction.processors`

Required:

- `processor` (Block List, Min: 1) (see [below for nested schema](#nestedblock--metric_extraction--processors--processor))

<a id="nestedblock--metric_extraction--processors--processor"></a>
### Nested Schema for `metric_extraction.processors.processor`

Required:

- `description` (String) No documentation available
- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)
- `id` (String) Processor identifier
- `type` (String) Processor type. Possible values: `dql`, `drop`, `fieldsAdd`, `fieldsRemove`, `fieldsRename`, `securityContext`

Optional:

- `dql` (Block List, Max: 1) DQL processor attributes (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--dql))
- `fields_add` (Block List, Max: 1) Fields add processor attributes (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--fields_add))
- `fields_remove` (Block List, Max: 1) Fields remove processor attributes (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--fields_remove))
- `fields_rename` (Block List, Max: 1) Fields rename processor attributes (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--fields_rename))
- `matcher` (String) [See our documentation](https://dt-url.net/bp234rv)
- `sample_data` (String) Sample data
- `security_context` (Block List, Max: 1) Security context processor attributes (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--security_context))

<a id="nestedblock--metric_extraction--processors--processor--dql"></a>
### Nested Schema for `metric_extraction.processors.processor.dql`

Required:

- `script` (String) DQL script


<a id="nestedblock--metric_extraction--processors--processor--fields_add"></a>
### Nested Schema for `metric_extraction.processors.processor.fields_add`

Required:

- `fields` (Block List, Min: 1, Max: 1) Fields to Add (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--fields_add--fields))

<a id="nestedblock--metric_extraction--processors--processor--fields_add--fields"></a>
### Nested Schema for `metric_extraction.processors.processor.fields_add.fields`

Required:

- `field` (Block List, Min: 1) (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--fields_add--fields--field))

<a id="nestedblock--metric_extraction--processors--processor--fields_add--fields--field"></a>
### Nested Schema for `metric_extraction.processors.processor.fields_add.fields.field`

Required:

- `name` (String) Fields's name
- `value` (String) Field's value

Optional:

- `value_type` (String) Field value type. Possible values: `boolean`, `double`, `long`, `string`




<a id="nestedblock--metric_extraction--processors--processor--fields_remove"></a>
### Nested Schema for `metric_extraction.processors.processor.fields_remove`

Required:

- `fields` (Set of String) Fields to remove


<a id="nestedblock--metric_extraction--processors--processor--fields_rename"></a>
### Nested Schema for `metric_extraction.processors.processor.fields_rename`

Required:

- `fields` (Block List, Min: 1, Max: 1) Fields to rename (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--fields_rename--fields))

<a id="nestedblock--metric_extraction--processors--processor--fields_rename--fields"></a>
### Nested Schema for `metric_extraction.processors.processor.fields_rename.fields`

Required:

- `field` (Block List, Min: 1) (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--fields_rename--fields--field))

<a id="nestedblock--metric_extraction--processors--processor--fields_rename--fields--field"></a>
### Nested Schema for `metric_extraction.processors.processor.fields_rename.fields.field`

Required:

- `from_name` (String) Fields's name
- `to_name` (String) New field's name




<a id="nestedblock--metric_extraction--processors--processor--security_context"></a>
### Nested Schema for `metric_extraction.processors.processor.security_context`

Required:

- `value` (Block List, Min: 1, Max: 1) Security context value assignment (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--security_context--value))

<a id="nestedblock--metric_extraction--processors--processor--security_context--value"></a>
### Nested Schema for `metric_extraction.processors.processor.security_context.value`

Required:

- `type` (String) Type of value assignment. Possible values: `constant`, `field`, `multiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--security_context--value--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--metric_extraction--processors--processor--security_context--value--field"></a>
### Nested Schema for `metric_extraction.processors.processor.security_context.value.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value







<a id="nestedblock--processing"></a>
### Nested Schema for `processing`

Optional:

- `processors` (Block List, Max: 1) Processors of stage (see [below for nested schema](#nestedblock--processing--processors))

<a id="nestedblock--processing--processors"></a>
### Nested Schema for `processing.processors`

Required:

- `processor` (Block List, Min: 1) (see [below for nested schema](#nestedblock--processing--processors--processor))

<a id="nestedblock--processing--processors--processor"></a>
### Nested Schema for `processing.processors.processor`

Required:

- `description` (String) No documentation available
- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)
- `id` (String) Processor identifier
- `type` (String) Processor type. Possible values: `dql`, `drop`, `fieldsAdd`, `fieldsRemove`, `fieldsRename`, `securityContext`

Optional:

- `dql` (Block List, Max: 1) DQL processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--dql))
- `fields_add` (Block List, Max: 1) Fields add processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--fields_add))
- `fields_remove` (Block List, Max: 1) Fields remove processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--fields_remove))
- `fields_rename` (Block List, Max: 1) Fields rename processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--fields_rename))
- `matcher` (String) [See our documentation](https://dt-url.net/bp234rv)
- `sample_data` (String) Sample data
- `security_context` (Block List, Max: 1) Security context processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--security_context))

<a id="nestedblock--processing--processors--processor--dql"></a>
### Nested Schema for `processing.processors.processor.dql`

Required:

- `script` (String) DQL script


<a id="nestedblock--processing--processors--processor--fields_add"></a>
### Nested Schema for `processing.processors.processor.fields_add`

Required:

- `fields` (Block List, Min: 1, Max: 1) Fields to Add (see [below for nested schema](#nestedblock--processing--processors--processor--fields_add--fields))

<a id="nestedblock--processing--processors--processor--fields_add--fields"></a>
### Nested Schema for `processing.processors.processor.fields_add.fields`

Required:

- `field` (Block List, Min: 1) (see [below for nested schema](#nestedblock--processing--processors--processor--fields_add--fields--field))

<a id="nestedblock--processing--processors--processor--fields_add--fields--field"></a>
### Nested Schema for `processing.processors.processor.fields_add.fields.field`

Required:

- `name` (String) Fields's name
- `value` (String) Field's value

Optional:

- `value_type` (String) Field value type. Possible values: `boolean`, `double`, `long`, `string`




<a id="nestedblock--processing--processors--processor--fields_remove"></a>
### Nested Schema for `processing.processors.processor.fields_remove`

Required:

- `fields` (Set of String) Fields to remove


<a id="nestedblock--processing--processors--processor--fields_rename"></a>
### Nested Schema for `processing.processors.processor.fields_rename`

Required:

- `fields` (Block List, Min: 1, Max: 1) Fields to rename (see [below for nested schema](#nestedblock--processing--processors--processor--fields_rename--fields))

<a id="nestedblock--processing--processors--processor--fields_rename--fields"></a>
### Nested Schema for `processing.processors.processor.fields_rename.fields`

Required:

- `field` (Block List, Min: 1) (see [below for nested schema](#nestedblock--processing--processors--processor--fields_rename--fields--field))

<a id="nestedblock--processing--processors--processor--fields_rename--fields--field"></a>
### Nested Schema for `processing.processors.processor.fields_rename.fields.field`

Required:

- `from_name` (String) Fields's name
- `to_name` (String) New field's name




<a id="nestedblock--processing--processors--processor--security_context"></a>
### Nested Schema for `processing.processors.processor.security_context`

Required:

- `value` (Block List, Min: 1, Max: 1) Security context value assignment (see [below for nested schema](#nestedblock--processing--processors--processor--security_context--value))

<a id="nestedblock--processing--processors--processor--security_context--value"></a>
### Nested Schema for `processing.processors.processor.security_context.value`

Required:

- `type` (String) Type of value assignment. Possible values: `constant`, `field`, `multiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--processing--processors--processor--security_context--value--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--processing--processors--processor--security_context--value--field"></a>
### Nested Schema for `processing.processors.processor.security_context.value.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value







<a id="nestedblock--product_allocation"></a>
### Nested Schema for `product_allocation`

Optional:

- `processors` (Block List, Max: 1) Processors of stage (see [below for nested schema](#nestedblock--product_allocation--processors))

<a id="nestedblock--product_allocation--processors"></a>
### Nested Schema for `product_allocation.processors`

Required:

- `processor` (Block List, Min: 1) (see [below for nested schema](#nestedblock--product_allocation--processors--processor))

<a id="nestedblock--product_allocation--processors--processor"></a>
### Nested Schema for `product_allocation.processors.processor`

Required:

- `description` (String) No documentation available
- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)
- `id` (String) Processor identifier
- `type` (String) Processor type. Possible values: `dql`, `drop`, `fieldsAdd`, `fieldsRemove`, `fieldsRename`, `securityContext`

Optional:

- `dql` (Block List, Max: 1) DQL processor attributes (see [below for nested schema](#nestedblock--product_allocation--processors--processor--dql))
- `fields_add` (Block List, Max: 1) Fields add processor attributes (see [below for nested schema](#nestedblock--product_allocation--processors--processor--fields_add))
- `fields_remove` (Block List, Max: 1) Fields remove processor attributes (see [below for nested schema](#nestedblock--product_allocation--processors--processor--fields_remove))
- `fields_rename` (Block List, Max: 1) Fields rename processor attributes (see [below for nested schema](#nestedblock--product_allocation--processors--processor--fields_rename))
- `matcher` (String) [See our documentation](https://dt-url.net/bp234rv)
- `sample_data` (String) Sample data
- `security_context` (Block List, Max: 1) Security context processor attributes (see [below for nested schema](#nestedblock--product_allocation--processors--processor--security_context))

<a id="nestedblock--product_allocation--processors--processor--dql"></a>
### Nested Schema for `product_allocation.processors.processor.dql`

Required:

- `script` (String) DQL script


<a id="nestedblock--product_allocation--processors--processor--fields_add"></a>
### Nested Schema for `product_allocation.processors.processor.fields_add`

Required:

- `fields` (Block List, Min: 1, Max: 1) Fields to Add (see [below for nested schema](#nestedblock--product_allocation--processors--processor--fields_add--fields))

<a id="nestedblock--product_allocation--processors--processor--fields_add--fields"></a>
### Nested Schema for `product_allocation.processors.processor.fields_add.fields`

Required:

- `field` (Block List, Min: 1) (see [below for nested schema](#nestedblock--product_allocation--processors--processor--fields_add--fields--field))

<a id="nestedblock--product_allocation--processors--processor--fields_add--fields--field"></a>
### Nested Schema for `product_allocation.processors.processor.fields_add.fields.field`

Required:

- `name` (String) Fields's name
- `value` (String) Field's value

Optional:

- `value_type` (String) Field value type. Possible values: `boolean`, `double`, `long`, `string`




<a id="nestedblock--product_allocation--processors--processor--fields_remove"></a>
### Nested Schema for `product_allocation.processors.processor.fields_remove`

Required:

- `fields` (Set of String) Fields to remove


<a id="nestedblock--product_allocation--processors--processor--fields_rename"></a>
### Nested Schema for `product_allocation.processors.processor.fields_rename`

Required:

- `fields` (Block List, Min: 1, Max: 1) Fields to rename (see [below for nested schema](#nestedblock--product_allocation--processors--processor--fields_rename--fields))

<a id="nestedblock--product_allocation--processors--processor--fields_rename--fields"></a>
### Nested Schema for `product_allocation.processors.processor.fields_rename.fields`

Required:

- `field` (Block List, Min: 1) (see [below for nested schema](#nestedblock--product_allocation--processors--processor--fields_rename--fields--field))

<a id="nestedblock--product_allocation--processors--processor--fields_rename--fields--field"></a>
### Nested Schema for `product_allocation.processors.processor.fields_rename.fields.field`

Required:

- `from_name` (String) Fields's name
- `to_name` (String) New field's name




<a id="nestedblock--product_allocation--processors--processor--security_context"></a>
### Nested Schema for `product_allocation.processors.processor.security_context`

Required:

- `value` (Block List, Min: 1, Max: 1) Security context value assignment (see [below for nested schema](#nestedblock--product_allocation--processors--processor--security_context--value))

<a id="nestedblock--product_allocation--processors--processor--security_context--value"></a>
### Nested Schema for `product_allocation.processors.processor.security_context.value`

Required:

- `type` (String) Type of value assignment. Possible values: `constant`, `field`, `multiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--product_allocation--processors--processor--security_context--value--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--product_allocation--processors--processor--security_context--value--field"></a>
### Nested Schema for `product_allocation.processors.processor.security_context.value.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value







<a id="nestedblock--security"></a>
### Nested Schema for `security`

Optional:

- `processors` (Block List, Max: 1) Processors of stage (see [below for nested schema](#nestedblock--security--processors))

<a id="nestedblock--security--processors"></a>
### Nested Schema for `security.processors`

Required:

- `processor` (Block List, Min: 1) (see [below for nested schema](#nestedblock--security--processors--processor))

<a id="nestedblock--security--processors--processor"></a>
### Nested Schema for `security.processors.processor`

Required:

- `description` (String) No documentation available
- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)
- `id` (String) Processor identifier
- `type` (String) Processor type. Possible values: `dql`, `drop`, `fieldsAdd`, `fieldsRemove`, `fieldsRename`, `securityContext`

Optional:

- `dql` (Block List, Max: 1) DQL processor attributes (see [below for nested schema](#nestedblock--security--processors--processor--dql))
- `fields_add` (Block List, Max: 1) Fields add processor attributes (see [below for nested schema](#nestedblock--security--processors--processor--fields_add))
- `fields_remove` (Block List, Max: 1) Fields remove processor attributes (see [below for nested schema](#nestedblock--security--processors--processor--fields_remove))
- `fields_rename` (Block List, Max: 1) Fields rename processor attributes (see [below for nested schema](#nestedblock--security--processors--processor--fields_rename))
- `matcher` (String) [See our documentation](https://dt-url.net/bp234rv)
- `sample_data` (String) Sample data
- `security_context` (Block List, Max: 1) Security context processor attributes (see [below for nested schema](#nestedblock--security--processors--processor--security_context))

<a id="nestedblock--security--processors--processor--dql"></a>
### Nested Schema for `security.processors.processor.dql`

Required:

- `script` (String) DQL script


<a id="nestedblock--security--processors--processor--fields_add"></a>
### Nested Schema for `security.processors.processor.fields_add`

Required:

- `fields` (Block List, Min: 1, Max: 1) Fields to Add (see [below for nested schema](#nestedblock--security--processors--processor--fields_add--fields))

<a id="nestedblock--security--processors--processor--fields_add--fields"></a>
### Nested Schema for `security.processors.processor.fields_add.fields`

Required:

- `field` (Block List, Min: 1) (see [below for nested schema](#nestedblock--security--processors--processor--fields_add--fields--field))

<a id="nestedblock--security--processors--processor--fields_add--fields--field"></a>
### Nested Schema for `security.processors.processor.fields_add.fields.field`

Required:

- `name` (String) Fields's name
- `value` (String) Field's value

Optional:

- `value_type` (String) Field value type. Possible values: `boolean`, `double`, `long`, `string`




<a id="nestedblock--security--processors--processor--fields_remove"></a>
### Nested Schema for `security.processors.processor.fields_remove`

Required:

- `fields` (Set of String) Fields to remove


<a id="nestedblock--security--processors--processor--fields_rename"></a>
### Nested Schema for `security.processors.processor.fields_rename`

Required:

- `fields` (Block List, Min: 1, Max: 1) Fields to rename (see [below for nested schema](#nestedblock--security--processors--processor--fields_rename--fields))

<a id="nestedblock--security--processors--processor--fields_rename--fields"></a>
### Nested Schema for `security.processors.processor.fields_rename.fields`

Required:

- `field` (Block List, Min: 1) (see [below for nested schema](#nestedblock--security--processors--processor--fields_rename--fields--field))

<a id="nestedblock--security--processors--processor--fields_rename--fields--field"></a>
### Nested Schema for `security.processors.processor.fields_rename.fields.field`

Required:

- `from_name` (String) Fields's name
- `to_name` (String) New field's name




<a id="nestedblock--security--processors--processor--security_context"></a>
### Nested Schema for `security.processors.processor.security_context`

Required:

- `value` (Block List, Min: 1, Max: 1) Security context value assignment (see [below for nested schema](#nestedblock--security--processors--processor--security_context--value))

<a id="nestedblock--security--processors--processor--security_context--value"></a>
### Nested Schema for `security.processors.processor.security_context.value`

Required:

- `type` (String) Type of value assignment. Possible values: `constant`, `field`, `multiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--security--processors--processor--security_context--value--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--security--processors--processor--security_context--value--field"></a>
### Nested Schema for `security.processors.processor.security_context.value.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value







<a id="nestedblock--security_context"></a>
### Nested Schema for `security_context`

Optional:

- `processors` (Block List, Max: 1) Processors of stage (see [below for nested schema](#nestedblock--security_context--processors))

<a id="nestedblock--security_context--processors"></a>
### Nested Schema for `security_context.processors`

Required:

- `processor` (Block List, Min: 1) (see [below for nested schema](#nestedblock--security_context--processors--processor))

<a id="nestedblock--security_context--processors--processor"></a>
### Nested Schema for `security_context.processors.processor`

Required:

- `description` (String) No documentation available
- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)
- `id` (String) Processor identifier
- `type` (String) Processor type. Possible values: `dql`, `drop`, `fieldsAdd`, `fieldsRemove`, `fieldsRename`, `securityContext`

Optional:

- `dql` (Block List, Max: 1) DQL processor attributes (see [below for nested schema](#nestedblock--security_context--processors--processor--dql))
- `fields_add` (Block List, Max: 1) Fields add processor attributes (see [below for nested schema](#nestedblock--security_context--processors--processor--fields_add))
- `fields_remove` (Block List, Max: 1) Fields remove processor attributes (see [below for nested schema](#nestedblock--security_context--processors--processor--fields_remove))
- `fields_rename` (Block List, Max: 1) Fields rename processor attributes (see [below for nested schema](#nestedblock--security_context--processors--processor--fields_rename))
- `matcher` (String) [See our documentation](https://dt-url.net/bp234rv)
- `sample_data` (String) Sample data
- `security_context` (Block List, Max: 1) Security context processor attributes (see [below for nested schema](#nestedblock--security_context--processors--processor--security_context))

<a id="nestedblock--security_context--processors--processor--dql"></a>
### Nested Schema for `security_context.processors.processor.dql`

Required:

- `script` (String) DQL script


<a id="nestedblock--security_context--processors--processor--fields_add"></a>
### Nested Schema for `security_context.processors.processor.fields_add`

Required:

- `fields` (Block List, Min: 1, Max: 1) Fields to Add (see [below for nested schema](#nestedblock--security_context--processors--processor--fields_add--fields))

<a id="nestedblock--security_context--processors--processor--fields_add--fields"></a>
### Nested Schema for `security_context.processors.processor.fields_add.fields`

Required:

- `field` (Block List, Min: 1) (see [below for nested schema](#nestedblock--security_context--processors--processor--fields_add--fields--field))

<a id="nestedblock--security_context--processors--processor--fields_add--fields--field"></a>
### Nested Schema for `security_context.processors.processor.fields_add.fields.field`

Required:

- `name` (String) Fields's name
- `value` (String) Field's value

Optional:

- `value_type` (String) Field value type. Possible values: `boolean`, `double`, `long`, `string`




<a id="nestedblock--security_context--processors--processor--fields_remove"></a>
### Nested Schema for `security_context.processors.processor.fields_remove`

Required:

- `fields` (Set of String) Fields to remove


<a id="nestedblock--security_context--processors--processor--fields_rename"></a>
### Nested Schema for `security_context.processors.processor.fields_rename`

Required:

- `fields` (Block List, Min: 1, Max: 1) Fields to rename (see [below for nested schema](#nestedblock--security_context--processors--processor--fields_rename--fields))

<a id="nestedblock--security_context--processors--processor--fields_rename--fields"></a>
### Nested Schema for `security_context.processors.processor.fields_rename.fields`

Required:

- `field` (Block List, Min: 1) (see [below for nested schema](#nestedblock--security_context--processors--processor--fields_rename--fields--field))

<a id="nestedblock--security_context--processors--processor--fields_rename--fields--field"></a>
### Nested Schema for `security_context.processors.processor.fields_rename.fields.field`

Required:

- `from_name` (String) Fields's name
- `to_name` (String) New field's name




<a id="nestedblock--security_context--processors--processor--security_context"></a>
### Nested Schema for `security_context.processors.processor.security_context`

Required:

- `value` (Block List, Min: 1, Max: 1) Security context value assignment (see [below for nested schema](#nestedblock--security_context--processors--processor--security_context--value))

<a id="nestedblock--security_context--processors--processor--security_context--value"></a>
### Nested Schema for `security_context.processors.processor.security_context.value`

Required:

- `type` (String) Type of value assignment. Possible values: `constant`, `field`, `multiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--security_context--processors--processor--security_context--value--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--security_context--processors--processor--security_context--value--field"></a>
### Nested Schema for `security_context.processors.processor.security_context.value.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value







<a id="nestedblock--smartscape_edge_extraction"></a>
### Nested Schema for `smartscape_edge_extraction`

Optional:

- `processors` (Block List, Max: 1) Processors of stage (see [below for nested schema](#nestedblock--smartscape_edge_extraction--processors))

<a id="nestedblock--smartscape_edge_extraction--processors"></a>
### Nested Schema for `smartscape_edge_extraction.processors`

Required:

- `processor` (Block List, Min: 1) (see [below for nested schema](#nestedblock--smartscape_edge_extraction--processors--processor))

<a id="nestedblock--smartscape_edge_extraction--processors--processor"></a>
### Nested Schema for `smartscape_edge_extraction.processors.processor`

Required:

- `description` (String) No documentation available
- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)
- `id` (String) Processor identifier
- `type` (String) Processor type. Possible values: `dql`, `drop`, `fieldsAdd`, `fieldsRemove`, `fieldsRename`, `securityContext`

Optional:

- `dql` (Block List, Max: 1) DQL processor attributes (see [below for nested schema](#nestedblock--smartscape_edge_extraction--processors--processor--dql))
- `fields_add` (Block List, Max: 1) Fields add processor attributes (see [below for nested schema](#nestedblock--smartscape_edge_extraction--processors--processor--fields_add))
- `fields_remove` (Block List, Max: 1) Fields remove processor attributes (see [below for nested schema](#nestedblock--smartscape_edge_extraction--processors--processor--fields_remove))
- `fields_rename` (Block List, Max: 1) Fields rename processor attributes (see [below for nested schema](#nestedblock--smartscape_edge_extraction--processors--processor--fields_rename))
- `matcher` (String) [See our documentation](https://dt-url.net/bp234rv)
- `sample_data` (String) Sample data
- `security_context` (Block List, Max: 1) Security context processor attributes (see [below for nested schema](#nestedblock--smartscape_edge_extraction--processors--processor--security_context))

<a id="nestedblock--smartscape_edge_extraction--processors--processor--dql"></a>
### Nested Schema for `smartscape_edge_extraction.processors.processor.dql`

Required:

- `script` (String) DQL script


<a id="nestedblock--smartscape_edge_extraction--processors--processor--fields_add"></a>
### Nested Schema for `smartscape_edge_extraction.processors.processor.fields_add`

Required:

- `fields` (Block List, Min: 1, Max: 1) Fields to Add (see [below for nested schema](#nestedblock--smartscape_edge_extraction--processors--processor--fields_add--fields))

<a id="nestedblock--smartscape_edge_extraction--processors--processor--fields_add--fields"></a>
### Nested Schema for `smartscape_edge_extraction.processors.processor.fields_add.fields`

Required:

- `field` (Block List, Min: 1) (see [below for nested schema](#nestedblock--smartscape_edge_extraction--processors--processor--fields_add--fields--field))

<a id="nestedblock--smartscape_edge_extraction--processors--processor--fields_add--fields--field"></a>
### Nested Schema for `smartscape_edge_extraction.processors.processor.fields_add.fields.field`

Required:

- `name` (String) Fields's name
- `value` (String) Field's value

Optional:

- `value_type` (String) Field value type. Possible values: `boolean`, `double`, `long`, `string`




<a id="nestedblock--smartscape_edge_extraction--processors--processor--fields_remove"></a>
### Nested Schema for `smartscape_edge_extraction.processors.processor.fields_remove`

Required:

- `fields` (Set of String) Fields to remove


<a id="nestedblock--smartscape_edge_extraction--processors--processor--fields_rename"></a>
### Nested Schema for `smartscape_edge_extraction.processors.processor.fields_rename`

Required:

- `fields` (Block List, Min: 1, Max: 1) Fields to rename (see [below for nested schema](#nestedblock--smartscape_edge_extraction--processors--processor--fields_rename--fields))

<a id="nestedblock--smartscape_edge_extraction--processors--processor--fields_rename--fields"></a>
### Nested Schema for `smartscape_edge_extraction.processors.processor.fields_rename.fields`

Required:

- `field` (Block List, Min: 1) (see [below for nested schema](#nestedblock--smartscape_edge_extraction--processors--processor--fields_rename--fields--field))

<a id="nestedblock--smartscape_edge_extraction--processors--processor--fields_rename--fields--field"></a>
### Nested Schema for `smartscape_edge_extraction.processors.processor.fields_rename.fields.field`

Required:

- `from_name` (String) Fields's name
- `to_name` (String) New field's name




<a id="nestedblock--smartscape_edge_extraction--processors--processor--security_context"></a>
### Nested Schema for `smartscape_edge_extraction.processors.processor.security_context`

Required:

- `value` (Block List, Min: 1, Max: 1) Security context value assignment (see [below for nested schema](#nestedblock--smartscape_edge_extraction--processors--processor--security_context--value))

<a id="nestedblock--smartscape_edge_extraction--processors--processor--security_context--value"></a>
### Nested Schema for `smartscape_edge_extraction.processors.processor.security_context.value`

Required:

- `type` (String) Type of value assignment. Possible values: `constant`, `field`, `multiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--smartscape_edge_extraction--processors--processor--security_context--value--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--smartscape_edge_extraction--processors--processor--security_context--value--field"></a>
### Nested Schema for `smartscape_edge_extraction.processors.processor.security_context.value.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value







<a id="nestedblock--smartscape_node_extraction"></a>
### Nested Schema for `smartscape_node_extraction`

Optional:

- `processors` (Block List, Max: 1) Processors of stage (see [below for nested schema](#nestedblock--smartscape_node_extraction--processors))

<a id="nestedblock--smartscape_node_extraction--processors"></a>
### Nested Schema for `smartscape_node_extraction.processors`

Required:

- `processor` (Block List, Min: 1) (see [below for nested schema](#nestedblock--smartscape_node_extraction--processors--processor))

<a id="nestedblock--smartscape_node_extraction--processors--processor"></a>
### Nested Schema for `smartscape_node_extraction.processors.processor`

Required:

- `description` (String) No documentation available
- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)
- `id` (String) Processor identifier
- `type` (String) Processor type. Possible values: `dql`, `drop`, `fieldsAdd`, `fieldsRemove`, `fieldsRename`, `securityContext`

Optional:

- `dql` (Block List, Max: 1) DQL processor attributes (see [below for nested schema](#nestedblock--smartscape_node_extraction--processors--processor--dql))
- `fields_add` (Block List, Max: 1) Fields add processor attributes (see [below for nested schema](#nestedblock--smartscape_node_extraction--processors--processor--fields_add))
- `fields_remove` (Block List, Max: 1) Fields remove processor attributes (see [below for nested schema](#nestedblock--smartscape_node_extraction--processors--processor--fields_remove))
- `fields_rename` (Block List, Max: 1) Fields rename processor attributes (see [below for nested schema](#nestedblock--smartscape_node_extraction--processors--processor--fields_rename))
- `matcher` (String) [See our documentation](https://dt-url.net/bp234rv)
- `sample_data` (String) Sample data
- `security_context` (Block List, Max: 1) Security context processor attributes (see [below for nested schema](#nestedblock--smartscape_node_extraction--processors--processor--security_context))

<a id="nestedblock--smartscape_node_extraction--processors--processor--dql"></a>
### Nested Schema for `smartscape_node_extraction.processors.processor.dql`

Required:

- `script` (String) DQL script


<a id="nestedblock--smartscape_node_extraction--processors--processor--fields_add"></a>
### Nested Schema for `smartscape_node_extraction.processors.processor.fields_add`

Required:

- `fields` (Block List, Min: 1, Max: 1) Fields to Add (see [below for nested schema](#nestedblock--smartscape_node_extraction--processors--processor--fields_add--fields))

<a id="nestedblock--smartscape_node_extraction--processors--processor--fields_add--fields"></a>
### Nested Schema for `smartscape_node_extraction.processors.processor.fields_add.fields`

Required:

- `field` (Block List, Min: 1) (see [below for nested schema](#nestedblock--smartscape_node_extraction--processors--processor--fields_add--fields--field))

<a id="nestedblock--smartscape_node_extraction--processors--processor--fields_add--fields--field"></a>
### Nested Schema for `smartscape_node_extraction.processors.processor.fields_add.fields.field`

Required:

- `name` (String) Fields's name
- `value` (String) Field's value

Optional:

- `value_type` (String) Field value type. Possible values: `boolean`, `double`, `long`, `string`




<a id="nestedblock--smartscape_node_extraction--processors--processor--fields_remove"></a>
### Nested Schema for `smartscape_node_extraction.processors.processor.fields_remove`

Required:

- `fields` (Set of String) Fields to remove


<a id="nestedblock--smartscape_node_extraction--processors--processor--fields_rename"></a>
### Nested Schema for `smartscape_node_extraction.processors.processor.fields_rename`

Required:

- `fields` (Block List, Min: 1, Max: 1) Fields to rename (see [below for nested schema](#nestedblock--smartscape_node_extraction--processors--processor--fields_rename--fields))

<a id="nestedblock--smartscape_node_extraction--processors--processor--fields_rename--fields"></a>
### Nested Schema for `smartscape_node_extraction.processors.processor.fields_rename.fields`

Required:

- `field` (Block List, Min: 1) (see [below for nested schema](#nestedblock--smartscape_node_extraction--processors--processor--fields_rename--fields--field))

<a id="nestedblock--smartscape_node_extraction--processors--processor--fields_rename--fields--field"></a>
### Nested Schema for `smartscape_node_extraction.processors.processor.fields_rename.fields.field`

Required:

- `from_name` (String) Fields's name
- `to_name` (String) New field's name




<a id="nestedblock--smartscape_node_extraction--processors--processor--security_context"></a>
### Nested Schema for `smartscape_node_extraction.processors.processor.security_context`

Required:

- `value` (Block List, Min: 1, Max: 1) Security context value assignment (see [below for nested schema](#nestedblock--smartscape_node_extraction--processors--processor--security_context--value))

<a id="nestedblock--smartscape_node_extraction--processors--processor--security_context--value"></a>
### Nested Schema for `smartscape_node_extraction.processors.processor.security_context.value`

Required:

- `type` (String) Type of value assignment. Possible values: `constant`, `field`, `multiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--smartscape_node_extraction--processors--processor--security_context--value--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--smartscape_node_extraction--processors--processor--security_context--value--field"></a>
### Nested Schema for `smartscape_node_extraction.processors.processor.security_context.value.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value







<a id="nestedblock--storage"></a>
### Nested Schema for `storage`

Optional:

- `processors` (Block List, Max: 1) Processors of stage (see [below for nested schema](#nestedblock--storage--processors))

<a id="nestedblock--storage--processors"></a>
### Nested Schema for `storage.processors`

Required:

- `processor` (Block List, Min: 1) (see [below for nested schema](#nestedblock--storage--processors--processor))

<a id="nestedblock--storage--processors--processor"></a>
### Nested Schema for `storage.processors.processor`

Required:

- `description` (String) No documentation available
- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)
- `id` (String) Processor identifier
- `type` (String) Processor type. Possible values: `dql`, `drop`, `fieldsAdd`, `fieldsRemove`, `fieldsRename`, `securityContext`

Optional:

- `dql` (Block List, Max: 1) DQL processor attributes (see [below for nested schema](#nestedblock--storage--processors--processor--dql))
- `fields_add` (Block List, Max: 1) Fields add processor attributes (see [below for nested schema](#nestedblock--storage--processors--processor--fields_add))
- `fields_remove` (Block List, Max: 1) Fields remove processor attributes (see [below for nested schema](#nestedblock--storage--processors--processor--fields_remove))
- `fields_rename` (Block List, Max: 1) Fields rename processor attributes (see [below for nested schema](#nestedblock--storage--processors--processor--fields_rename))
- `matcher` (String) [See our documentation](https://dt-url.net/bp234rv)
- `sample_data` (String) Sample data
- `security_context` (Block List, Max: 1) Security context processor attributes (see [below for nested schema](#nestedblock--storage--processors--processor--security_context))

<a id="nestedblock--storage--processors--processor--dql"></a>
### Nested Schema for `storage.processors.processor.dql`

Required:

- `script` (String) DQL script


<a id="nestedblock--storage--processors--processor--fields_add"></a>
### Nested Schema for `storage.processors.processor.fields_add`

Required:

- `fields` (Block List, Min: 1, Max: 1) Fields to Add (see [below for nested schema](#nestedblock--storage--processors--processor--fields_add--fields))

<a id="nestedblock--storage--processors--processor--fields_add--fields"></a>
### Nested Schema for `storage.processors.processor.fields_add.fields`

Required:

- `field` (Block List, Min: 1) (see [below for nested schema](#nestedblock--storage--processors--processor--fields_add--fields--field))

<a id="nestedblock--storage--processors--processor--fields_add--fields--field"></a>
### Nested Schema for `storage.processors.processor.fields_add.fields.field`

Required:

- `name` (String) Fields's name
- `value` (String) Field's value

Optional:

- `value_type` (String) Field value type. Possible values: `boolean`, `double`, `long`, `string`




<a id="nestedblock--storage--processors--processor--fields_remove"></a>
### Nested Schema for `storage.processors.processor.fields_remove`

Required:

- `fields` (Set of String) Fields to remove


<a id="nestedblock--storage--processors--processor--fields_rename"></a>
### Nested Schema for `storage.processors.processor.fields_rename`

Required:

- `fields` (Block List, Min: 1, Max: 1) Fields to rename (see [below for nested schema](#nestedblock--storage--processors--processor--fields_rename--fields))

<a id="nestedblock--storage--processors--processor--fields_rename--fields"></a>
### Nested Schema for `storage.processors.processor.fields_rename.fields`

Required:

- `field` (Block List, Min: 1) (see [below for nested schema](#nestedblock--storage--processors--processor--fields_rename--fields--field))

<a id="nestedblock--storage--processors--processor--fields_rename--fields--field"></a>
### Nested Schema for `storage.processors.processor.fields_rename.fields.field`

Required:

- `from_name` (String) Fields's name
- `to_name` (String) New field's name




<a id="nestedblock--storage--processors--processor--security_context"></a>
### Nested Schema for `storage.processors.processor.security_context`

Required:

- `value` (Block List, Min: 1, Max: 1) Security context value assignment (see [below for nested schema](#nestedblock--storage--processors--processor--security_context--value))

<a id="nestedblock--storage--processors--processor--security_context--value"></a>
### Nested Schema for `storage.processors.processor.security_context.value`

Required:

- `type` (String) Type of value assignment. Possible values: `constant`, `field`, `multiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--storage--processors--processor--security_context--value--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--storage--processors--processor--security_context--value--field"></a>
### Nested Schema for `storage.processors.processor.security_context.value.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
