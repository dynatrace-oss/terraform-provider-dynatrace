---
layout: ""
page_title: "dynatrace_openpipeline_v2_spans_pipelines Resource - terraform-provider-dynatrace"
subcategory: "OpenPipeline V2"
description: |-
  The resource `dynatrace_openpipeline_v2_spans_pipelines` covers configuration of OpenPipeline for spans pipelines
---

# dynatrace_openpipeline_v2_spans_pipelines (Resource)

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

- `terraform-provider-dynatrace -export dynatrace_openpipeline_v2_spans_pipelines` downloads all existing OpenPipeline definitions for spans pipelines

The full documentation of the export feature is available [here](https://dt-url.net/h203qmc).

## Resource Example Usage

```terraform
resource "dynatrace_openpipeline_v2_spans_pipelines" "max-pipeline" {
  display_name = "Warning pipeline"
  custom_id = "pipeline_Warning_pipeline_2773_tf_#name#"
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
            type = "multiValueConstant"
            multi_value_constant = ["ACME1", "ACME2"]
          }
        }
        enabled = true
      }
    }
  }
  cost_allocation {

  }
  product_allocation {

  }
  storage {
    processors {
      processor {
        type        = "bucketAssignment"
        id          = "processor_Add_to_default_bucket_5010"
        description = "Add to default bucket"
        matcher     = "true"
        bucket_assignment {
          bucket_name = "default_events"
        }
        enabled = true
      }
    }
  }
  data_extraction {}
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cost_allocation` (Block List, Min: 1, Max: 1) Cost allocation stage (see [below for nested schema](#nestedblock--cost_allocation))
- `custom_id` (String) Custom pipeline id
- `data_extraction` (Block List, Min: 1, Max: 1) Data extraction stage (see [below for nested schema](#nestedblock--data_extraction))
- `davis` (Block List, Min: 1, Max: 1) Davis event extraction stage (see [below for nested schema](#nestedblock--davis))
- `display_name` (String) Display name
- `metric_extraction` (Block List, Min: 1, Max: 1) Metrics extraction stage (see [below for nested schema](#nestedblock--metric_extraction))
- `processing` (Block List, Min: 1, Max: 1) Processing stage (see [below for nested schema](#nestedblock--processing))
- `product_allocation` (Block List, Min: 1, Max: 1) Product allocation stage (see [below for nested schema](#nestedblock--product_allocation))
- `security_context` (Block List, Min: 1, Max: 1) Security context stage (see [below for nested schema](#nestedblock--security_context))
- `storage` (Block List, Min: 1, Max: 1) Storage stage (see [below for nested schema](#nestedblock--storage))

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

- `description` (String) no documentation available
- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)
- `id` (String) Processor identifier
- `type` (String) Possible Values: `AzureLogForwarding`, `Bizevent`, `BucketAssignment`, `CostAllocation`, `CounterMetric`, `Davis`, `Dql`, `Drop`, `FieldsAdd`, `FieldsRemove`, `FieldsRename`, `HistogramMetric`, `NoStorage`, `ProductAllocation`, `SamplingAwareCounterMetric`, `SamplingAwareValueMetric`, `SecurityContext`, `SecurityEvent`, `Technology`, `ValueMetric`

Optional:

- `azure_log_forwarding` (Block List, Max: 1) Azure log forwarding processor attributes (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--azure_log_forwarding))
- `bizevent` (Block List, Max: 1) Bizevent extraction processor attributes (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--bizevent))
- `bucket_assignment` (Block List, Max: 1) Bucket assignment processor attributes (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--bucket_assignment))
- `cost_allocation` (Block List, Max: 1) Cost allocation processor attributes (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--cost_allocation))
- `counter_metric` (Block List, Max: 1) Counter metric processor attributes (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--counter_metric))
- `davis` (Block List, Max: 1) Davis event extraction processor attributes (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--davis))
- `dql` (Block List, Max: 1) DQL processor attributes (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--dql))
- `fields_add` (Block List, Max: 1) Fields add processor attributes (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--fields_add))
- `fields_remove` (Block List, Max: 1) Fields remove processor attributes (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--fields_remove))
- `fields_rename` (Block List, Max: 1) Fields rename processor attributes (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--fields_rename))
- `histogram_metric` (Block List, Max: 1) Histogram metric processor attributes (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--histogram_metric))
- `matcher` (String) [See our documentation](https://dt-url.net/bp234rv)
- `product_allocation` (Block List, Max: 1) Product allocation processor attributes (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--product_allocation))
- `sample_data` (String) Sample data
- `sampling_aware_counter_metric` (Block List, Max: 1) Sampling-aware counter metric processor attributes (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--sampling_aware_counter_metric))
- `sampling_aware_value_metric` (Block List, Max: 1) Sampling aware value metric processor attributes (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--sampling_aware_value_metric))
- `security_context` (Block List, Max: 1) Security context processor attributes (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--security_context))
- `security_event` (Block List, Max: 1) Security event extraction processor attributes (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--security_event))
- `technology` (Block List, Max: 1) Technology processor attributes (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--technology))
- `value_metric` (Block List, Max: 1) Value metric processor attributes (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--value_metric))

<a id="nestedblock--cost_allocation--processors--processor--azure_log_forwarding"></a>
### Nested Schema for `cost_allocation.processors.processor.azure_log_forwarding`

Required:

- `field_extraction` (Block List, Min: 1, Max: 1) Field Extraction (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--azure_log_forwarding--field_extraction))
- `forwarder_config_id` (String) no documentation available

<a id="nestedblock--cost_allocation--processors--processor--azure_log_forwarding--field_extraction"></a>
### Nested Schema for `cost_allocation.processors.processor.azure_log_forwarding.field_extraction`

Required:

- `type` (String) Possible Values: `Exclude`, `Include`, `IncludeAll`

Optional:

- `exclude` (Set of String) Fields
- `include` (Block List, Max: 1) Fields (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--azure_log_forwarding--field_extraction--include))

<a id="nestedblock--cost_allocation--processors--processor--azure_log_forwarding--field_extraction--include"></a>
### Nested Schema for `cost_allocation.processors.processor.azure_log_forwarding.field_extraction.include`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--azure_log_forwarding--field_extraction--include--dimension))

<a id="nestedblock--cost_allocation--processors--processor--azure_log_forwarding--field_extraction--include--dimension"></a>
### Nested Schema for `cost_allocation.processors.processor.azure_log_forwarding.field_extraction.include.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name





<a id="nestedblock--cost_allocation--processors--processor--bizevent"></a>
### Nested Schema for `cost_allocation.processors.processor.bizevent`

Required:

- `event_provider` (Block List, Min: 1, Max: 1) Event provider (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--bizevent--event_provider))
- `field_extraction` (Block List, Min: 1, Max: 1) Field extraction (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--bizevent--field_extraction))

Optional:

- `event_type` (Block List, Max: 1) Event type (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--bizevent--event_type))

<a id="nestedblock--cost_allocation--processors--processor--bizevent--event_provider"></a>
### Nested Schema for `cost_allocation.processors.processor.bizevent.event_provider`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--bizevent--event_provider--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--cost_allocation--processors--processor--bizevent--event_provider--field"></a>
### Nested Schema for `cost_allocation.processors.processor.bizevent.event_provider.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value



<a id="nestedblock--cost_allocation--processors--processor--bizevent--field_extraction"></a>
### Nested Schema for `cost_allocation.processors.processor.bizevent.field_extraction`

Required:

- `type` (String) Possible Values: `Exclude`, `Include`, `IncludeAll`

Optional:

- `exclude` (Set of String) Fields
- `include` (Block List, Max: 1) Fields (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--bizevent--field_extraction--include))

<a id="nestedblock--cost_allocation--processors--processor--bizevent--field_extraction--include"></a>
### Nested Schema for `cost_allocation.processors.processor.bizevent.field_extraction.include`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--bizevent--field_extraction--include--dimension))

<a id="nestedblock--cost_allocation--processors--processor--bizevent--field_extraction--include--dimension"></a>
### Nested Schema for `cost_allocation.processors.processor.bizevent.field_extraction.include.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--cost_allocation--processors--processor--bizevent--event_type"></a>
### Nested Schema for `cost_allocation.processors.processor.bizevent.event_type`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--bizevent--event_type--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--cost_allocation--processors--processor--bizevent--event_type--field"></a>
### Nested Schema for `cost_allocation.processors.processor.bizevent.event_type.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value




<a id="nestedblock--cost_allocation--processors--processor--bucket_assignment"></a>
### Nested Schema for `cost_allocation.processors.processor.bucket_assignment`

Required:

- `bucket_name` (String) Bucket name


<a id="nestedblock--cost_allocation--processors--processor--cost_allocation"></a>
### Nested Schema for `cost_allocation.processors.processor.cost_allocation`

Required:

- `value` (Block List, Min: 1, Max: 1) The strategy to set the cost allocation field (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--cost_allocation--value))

<a id="nestedblock--cost_allocation--processors--processor--cost_allocation--value"></a>
### Nested Schema for `cost_allocation.processors.processor.cost_allocation.value`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--cost_allocation--value--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--cost_allocation--processors--processor--cost_allocation--value--field"></a>
### Nested Schema for `cost_allocation.processors.processor.cost_allocation.value.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value




<a id="nestedblock--cost_allocation--processors--processor--counter_metric"></a>
### Nested Schema for `cost_allocation.processors.processor.counter_metric`

Required:

- `metric_key` (String) Metric key

Optional:

- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--counter_metric--dimensions))

<a id="nestedblock--cost_allocation--processors--processor--counter_metric--dimensions"></a>
### Nested Schema for `cost_allocation.processors.processor.counter_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--counter_metric--dimensions--dimension))

<a id="nestedblock--cost_allocation--processors--processor--counter_metric--dimensions--dimension"></a>
### Nested Schema for `cost_allocation.processors.processor.counter_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--cost_allocation--processors--processor--davis"></a>
### Nested Schema for `cost_allocation.processors.processor.davis`

Required:

- `properties` (Block List, Min: 1, Max: 1) no documentation available (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--davis--properties))

<a id="nestedblock--cost_allocation--processors--processor--davis--properties"></a>
### Nested Schema for `cost_allocation.processors.processor.davis.properties`

Required:

- `property` (Block List, Min: 1) (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--davis--properties--property))

<a id="nestedblock--cost_allocation--processors--processor--davis--properties--property"></a>
### Nested Schema for `cost_allocation.processors.processor.davis.properties.property`

Required:

- `key` (String) no documentation available
- `value` (String) no documentation available




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




<a id="nestedblock--cost_allocation--processors--processor--histogram_metric"></a>
### Nested Schema for `cost_allocation.processors.processor.histogram_metric`

Required:

- `field` (String) Field with metric value
- `metric_key` (String) Metric key

Optional:

- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--histogram_metric--dimensions))

<a id="nestedblock--cost_allocation--processors--processor--histogram_metric--dimensions"></a>
### Nested Schema for `cost_allocation.processors.processor.histogram_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--histogram_metric--dimensions--dimension))

<a id="nestedblock--cost_allocation--processors--processor--histogram_metric--dimensions--dimension"></a>
### Nested Schema for `cost_allocation.processors.processor.histogram_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--cost_allocation--processors--processor--product_allocation"></a>
### Nested Schema for `cost_allocation.processors.processor.product_allocation`

Required:

- `value` (Block List, Min: 1, Max: 1) The strategy to set product allocation field (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--product_allocation--value))

<a id="nestedblock--cost_allocation--processors--processor--product_allocation--value"></a>
### Nested Schema for `cost_allocation.processors.processor.product_allocation.value`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--product_allocation--value--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--cost_allocation--processors--processor--product_allocation--value--field"></a>
### Nested Schema for `cost_allocation.processors.processor.product_allocation.value.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value




<a id="nestedblock--cost_allocation--processors--processor--sampling_aware_counter_metric"></a>
### Nested Schema for `cost_allocation.processors.processor.sampling_aware_counter_metric`

Required:

- `metric_key` (String) Metric key

Optional:

- `aggregation` (String) Possible Values: `Disabled`, `Enabled`
- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--sampling_aware_counter_metric--dimensions))
- `sampling` (String) Possible Values: `Disabled`, `Enabled`

<a id="nestedblock--cost_allocation--processors--processor--sampling_aware_counter_metric--dimensions"></a>
### Nested Schema for `cost_allocation.processors.processor.sampling_aware_counter_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--sampling_aware_counter_metric--dimensions--dimension))

<a id="nestedblock--cost_allocation--processors--processor--sampling_aware_counter_metric--dimensions--dimension"></a>
### Nested Schema for `cost_allocation.processors.processor.sampling_aware_counter_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--cost_allocation--processors--processor--sampling_aware_value_metric"></a>
### Nested Schema for `cost_allocation.processors.processor.sampling_aware_value_metric`

Required:

- `measurement` (String) Possible Values: `Duration`, `Field`
- `metric_key` (String) Metric key

Optional:

- `aggregation` (String) Possible Values: `Disabled`, `Enabled`
- `default_value` (String) Default value with metric value
- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--sampling_aware_value_metric--dimensions))
- `field` (String) Field with metric value
- `sampling` (String) Possible Values: `Disabled`, `Enabled`

<a id="nestedblock--cost_allocation--processors--processor--sampling_aware_value_metric--dimensions"></a>
### Nested Schema for `cost_allocation.processors.processor.sampling_aware_value_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--sampling_aware_value_metric--dimensions--dimension))

<a id="nestedblock--cost_allocation--processors--processor--sampling_aware_value_metric--dimensions--dimension"></a>
### Nested Schema for `cost_allocation.processors.processor.sampling_aware_value_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--cost_allocation--processors--processor--security_context"></a>
### Nested Schema for `cost_allocation.processors.processor.security_context`

Required:

- `value` (Block List, Min: 1, Max: 1) Security context value assignment (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--security_context--value))

<a id="nestedblock--cost_allocation--processors--processor--security_context--value"></a>
### Nested Schema for `cost_allocation.processors.processor.security_context.value`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

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




<a id="nestedblock--cost_allocation--processors--processor--security_event"></a>
### Nested Schema for `cost_allocation.processors.processor.security_event`

Required:

- `field_extraction` (Block List, Min: 1, Max: 1) Field Extraction (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--security_event--field_extraction))

<a id="nestedblock--cost_allocation--processors--processor--security_event--field_extraction"></a>
### Nested Schema for `cost_allocation.processors.processor.security_event.field_extraction`

Required:

- `type` (String) Possible Values: `Exclude`, `Include`, `IncludeAll`

Optional:

- `exclude` (Set of String) Fields
- `include` (Block List, Max: 1) Fields (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--security_event--field_extraction--include))

<a id="nestedblock--cost_allocation--processors--processor--security_event--field_extraction--include"></a>
### Nested Schema for `cost_allocation.processors.processor.security_event.field_extraction.include`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--security_event--field_extraction--include--dimension))

<a id="nestedblock--cost_allocation--processors--processor--security_event--field_extraction--include--dimension"></a>
### Nested Schema for `cost_allocation.processors.processor.security_event.field_extraction.include.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name





<a id="nestedblock--cost_allocation--processors--processor--technology"></a>
### Nested Schema for `cost_allocation.processors.processor.technology`

Required:

- `technology_id` (String) Technology ID

Optional:

- `custom_matcher` (String) Custom matching condition which should be used instead of technology matcher.


<a id="nestedblock--cost_allocation--processors--processor--value_metric"></a>
### Nested Schema for `cost_allocation.processors.processor.value_metric`

Required:

- `field` (String) Field with metric value
- `metric_key` (String) Metric key

Optional:

- `default_value` (String) Default value with metric value
- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--value_metric--dimensions))

<a id="nestedblock--cost_allocation--processors--processor--value_metric--dimensions"></a>
### Nested Schema for `cost_allocation.processors.processor.value_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--cost_allocation--processors--processor--value_metric--dimensions--dimension))

<a id="nestedblock--cost_allocation--processors--processor--value_metric--dimensions--dimension"></a>
### Nested Schema for `cost_allocation.processors.processor.value_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name







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

- `description` (String) no documentation available
- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)
- `id` (String) Processor identifier
- `type` (String) Possible Values: `AzureLogForwarding`, `Bizevent`, `BucketAssignment`, `CostAllocation`, `CounterMetric`, `Davis`, `Dql`, `Drop`, `FieldsAdd`, `FieldsRemove`, `FieldsRename`, `HistogramMetric`, `NoStorage`, `ProductAllocation`, `SamplingAwareCounterMetric`, `SamplingAwareValueMetric`, `SecurityContext`, `SecurityEvent`, `Technology`, `ValueMetric`

Optional:

- `azure_log_forwarding` (Block List, Max: 1) Azure log forwarding processor attributes (see [below for nested schema](#nestedblock--data_extraction--processors--processor--azure_log_forwarding))
- `bizevent` (Block List, Max: 1) Bizevent extraction processor attributes (see [below for nested schema](#nestedblock--data_extraction--processors--processor--bizevent))
- `bucket_assignment` (Block List, Max: 1) Bucket assignment processor attributes (see [below for nested schema](#nestedblock--data_extraction--processors--processor--bucket_assignment))
- `cost_allocation` (Block List, Max: 1) Cost allocation processor attributes (see [below for nested schema](#nestedblock--data_extraction--processors--processor--cost_allocation))
- `counter_metric` (Block List, Max: 1) Counter metric processor attributes (see [below for nested schema](#nestedblock--data_extraction--processors--processor--counter_metric))
- `davis` (Block List, Max: 1) Davis event extraction processor attributes (see [below for nested schema](#nestedblock--data_extraction--processors--processor--davis))
- `dql` (Block List, Max: 1) DQL processor attributes (see [below for nested schema](#nestedblock--data_extraction--processors--processor--dql))
- `fields_add` (Block List, Max: 1) Fields add processor attributes (see [below for nested schema](#nestedblock--data_extraction--processors--processor--fields_add))
- `fields_remove` (Block List, Max: 1) Fields remove processor attributes (see [below for nested schema](#nestedblock--data_extraction--processors--processor--fields_remove))
- `fields_rename` (Block List, Max: 1) Fields rename processor attributes (see [below for nested schema](#nestedblock--data_extraction--processors--processor--fields_rename))
- `histogram_metric` (Block List, Max: 1) Histogram metric processor attributes (see [below for nested schema](#nestedblock--data_extraction--processors--processor--histogram_metric))
- `matcher` (String) [See our documentation](https://dt-url.net/bp234rv)
- `product_allocation` (Block List, Max: 1) Product allocation processor attributes (see [below for nested schema](#nestedblock--data_extraction--processors--processor--product_allocation))
- `sample_data` (String) Sample data
- `sampling_aware_counter_metric` (Block List, Max: 1) Sampling-aware counter metric processor attributes (see [below for nested schema](#nestedblock--data_extraction--processors--processor--sampling_aware_counter_metric))
- `sampling_aware_value_metric` (Block List, Max: 1) Sampling aware value metric processor attributes (see [below for nested schema](#nestedblock--data_extraction--processors--processor--sampling_aware_value_metric))
- `security_context` (Block List, Max: 1) Security context processor attributes (see [below for nested schema](#nestedblock--data_extraction--processors--processor--security_context))
- `security_event` (Block List, Max: 1) Security event extraction processor attributes (see [below for nested schema](#nestedblock--data_extraction--processors--processor--security_event))
- `technology` (Block List, Max: 1) Technology processor attributes (see [below for nested schema](#nestedblock--data_extraction--processors--processor--technology))
- `value_metric` (Block List, Max: 1) Value metric processor attributes (see [below for nested schema](#nestedblock--data_extraction--processors--processor--value_metric))

<a id="nestedblock--data_extraction--processors--processor--azure_log_forwarding"></a>
### Nested Schema for `data_extraction.processors.processor.azure_log_forwarding`

Required:

- `field_extraction` (Block List, Min: 1, Max: 1) Field Extraction (see [below for nested schema](#nestedblock--data_extraction--processors--processor--azure_log_forwarding--field_extraction))
- `forwarder_config_id` (String) no documentation available

<a id="nestedblock--data_extraction--processors--processor--azure_log_forwarding--field_extraction"></a>
### Nested Schema for `data_extraction.processors.processor.azure_log_forwarding.field_extraction`

Required:

- `type` (String) Possible Values: `Exclude`, `Include`, `IncludeAll`

Optional:

- `exclude` (Set of String) Fields
- `include` (Block List, Max: 1) Fields (see [below for nested schema](#nestedblock--data_extraction--processors--processor--azure_log_forwarding--field_extraction--include))

<a id="nestedblock--data_extraction--processors--processor--azure_log_forwarding--field_extraction--include"></a>
### Nested Schema for `data_extraction.processors.processor.azure_log_forwarding.field_extraction.include`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--data_extraction--processors--processor--azure_log_forwarding--field_extraction--include--dimension))

<a id="nestedblock--data_extraction--processors--processor--azure_log_forwarding--field_extraction--include--dimension"></a>
### Nested Schema for `data_extraction.processors.processor.azure_log_forwarding.field_extraction.include.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name





<a id="nestedblock--data_extraction--processors--processor--bizevent"></a>
### Nested Schema for `data_extraction.processors.processor.bizevent`

Required:

- `event_provider` (Block List, Min: 1, Max: 1) Event provider (see [below for nested schema](#nestedblock--data_extraction--processors--processor--bizevent--event_provider))
- `field_extraction` (Block List, Min: 1, Max: 1) Field extraction (see [below for nested schema](#nestedblock--data_extraction--processors--processor--bizevent--field_extraction))

Optional:

- `event_type` (Block List, Max: 1) Event type (see [below for nested schema](#nestedblock--data_extraction--processors--processor--bizevent--event_type))

<a id="nestedblock--data_extraction--processors--processor--bizevent--event_provider"></a>
### Nested Schema for `data_extraction.processors.processor.bizevent.event_provider`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--data_extraction--processors--processor--bizevent--event_provider--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--data_extraction--processors--processor--bizevent--event_provider--field"></a>
### Nested Schema for `data_extraction.processors.processor.bizevent.event_provider.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value



<a id="nestedblock--data_extraction--processors--processor--bizevent--field_extraction"></a>
### Nested Schema for `data_extraction.processors.processor.bizevent.field_extraction`

Required:

- `type` (String) Possible Values: `Exclude`, `Include`, `IncludeAll`

Optional:

- `exclude` (Set of String) Fields
- `include` (Block List, Max: 1) Fields (see [below for nested schema](#nestedblock--data_extraction--processors--processor--bizevent--field_extraction--include))

<a id="nestedblock--data_extraction--processors--processor--bizevent--field_extraction--include"></a>
### Nested Schema for `data_extraction.processors.processor.bizevent.field_extraction.include`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--data_extraction--processors--processor--bizevent--field_extraction--include--dimension))

<a id="nestedblock--data_extraction--processors--processor--bizevent--field_extraction--include--dimension"></a>
### Nested Schema for `data_extraction.processors.processor.bizevent.field_extraction.include.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--data_extraction--processors--processor--bizevent--event_type"></a>
### Nested Schema for `data_extraction.processors.processor.bizevent.event_type`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--data_extraction--processors--processor--bizevent--event_type--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--data_extraction--processors--processor--bizevent--event_type--field"></a>
### Nested Schema for `data_extraction.processors.processor.bizevent.event_type.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value




<a id="nestedblock--data_extraction--processors--processor--bucket_assignment"></a>
### Nested Schema for `data_extraction.processors.processor.bucket_assignment`

Required:

- `bucket_name` (String) Bucket name


<a id="nestedblock--data_extraction--processors--processor--cost_allocation"></a>
### Nested Schema for `data_extraction.processors.processor.cost_allocation`

Required:

- `value` (Block List, Min: 1, Max: 1) The strategy to set the cost allocation field (see [below for nested schema](#nestedblock--data_extraction--processors--processor--cost_allocation--value))

<a id="nestedblock--data_extraction--processors--processor--cost_allocation--value"></a>
### Nested Schema for `data_extraction.processors.processor.cost_allocation.value`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--data_extraction--processors--processor--cost_allocation--value--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--data_extraction--processors--processor--cost_allocation--value--field"></a>
### Nested Schema for `data_extraction.processors.processor.cost_allocation.value.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value




<a id="nestedblock--data_extraction--processors--processor--counter_metric"></a>
### Nested Schema for `data_extraction.processors.processor.counter_metric`

Required:

- `metric_key` (String) Metric key

Optional:

- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--data_extraction--processors--processor--counter_metric--dimensions))

<a id="nestedblock--data_extraction--processors--processor--counter_metric--dimensions"></a>
### Nested Schema for `data_extraction.processors.processor.counter_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--data_extraction--processors--processor--counter_metric--dimensions--dimension))

<a id="nestedblock--data_extraction--processors--processor--counter_metric--dimensions--dimension"></a>
### Nested Schema for `data_extraction.processors.processor.counter_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--data_extraction--processors--processor--davis"></a>
### Nested Schema for `data_extraction.processors.processor.davis`

Required:

- `properties` (Block List, Min: 1, Max: 1) no documentation available (see [below for nested schema](#nestedblock--data_extraction--processors--processor--davis--properties))

<a id="nestedblock--data_extraction--processors--processor--davis--properties"></a>
### Nested Schema for `data_extraction.processors.processor.davis.properties`

Required:

- `property` (Block List, Min: 1) (see [below for nested schema](#nestedblock--data_extraction--processors--processor--davis--properties--property))

<a id="nestedblock--data_extraction--processors--processor--davis--properties--property"></a>
### Nested Schema for `data_extraction.processors.processor.davis.properties.property`

Required:

- `key` (String) no documentation available
- `value` (String) no documentation available




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




<a id="nestedblock--data_extraction--processors--processor--histogram_metric"></a>
### Nested Schema for `data_extraction.processors.processor.histogram_metric`

Required:

- `field` (String) Field with metric value
- `metric_key` (String) Metric key

Optional:

- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--data_extraction--processors--processor--histogram_metric--dimensions))

<a id="nestedblock--data_extraction--processors--processor--histogram_metric--dimensions"></a>
### Nested Schema for `data_extraction.processors.processor.histogram_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--data_extraction--processors--processor--histogram_metric--dimensions--dimension))

<a id="nestedblock--data_extraction--processors--processor--histogram_metric--dimensions--dimension"></a>
### Nested Schema for `data_extraction.processors.processor.histogram_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--data_extraction--processors--processor--product_allocation"></a>
### Nested Schema for `data_extraction.processors.processor.product_allocation`

Required:

- `value` (Block List, Min: 1, Max: 1) The strategy to set product allocation field (see [below for nested schema](#nestedblock--data_extraction--processors--processor--product_allocation--value))

<a id="nestedblock--data_extraction--processors--processor--product_allocation--value"></a>
### Nested Schema for `data_extraction.processors.processor.product_allocation.value`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--data_extraction--processors--processor--product_allocation--value--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--data_extraction--processors--processor--product_allocation--value--field"></a>
### Nested Schema for `data_extraction.processors.processor.product_allocation.value.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value




<a id="nestedblock--data_extraction--processors--processor--sampling_aware_counter_metric"></a>
### Nested Schema for `data_extraction.processors.processor.sampling_aware_counter_metric`

Required:

- `metric_key` (String) Metric key

Optional:

- `aggregation` (String) Possible Values: `Disabled`, `Enabled`
- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--data_extraction--processors--processor--sampling_aware_counter_metric--dimensions))
- `sampling` (String) Possible Values: `Disabled`, `Enabled`

<a id="nestedblock--data_extraction--processors--processor--sampling_aware_counter_metric--dimensions"></a>
### Nested Schema for `data_extraction.processors.processor.sampling_aware_counter_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--data_extraction--processors--processor--sampling_aware_counter_metric--dimensions--dimension))

<a id="nestedblock--data_extraction--processors--processor--sampling_aware_counter_metric--dimensions--dimension"></a>
### Nested Schema for `data_extraction.processors.processor.sampling_aware_counter_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--data_extraction--processors--processor--sampling_aware_value_metric"></a>
### Nested Schema for `data_extraction.processors.processor.sampling_aware_value_metric`

Required:

- `measurement` (String) Possible Values: `Duration`, `Field`
- `metric_key` (String) Metric key

Optional:

- `aggregation` (String) Possible Values: `Disabled`, `Enabled`
- `default_value` (String) Default value with metric value
- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--data_extraction--processors--processor--sampling_aware_value_metric--dimensions))
- `field` (String) Field with metric value
- `sampling` (String) Possible Values: `Disabled`, `Enabled`

<a id="nestedblock--data_extraction--processors--processor--sampling_aware_value_metric--dimensions"></a>
### Nested Schema for `data_extraction.processors.processor.sampling_aware_value_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--data_extraction--processors--processor--sampling_aware_value_metric--dimensions--dimension))

<a id="nestedblock--data_extraction--processors--processor--sampling_aware_value_metric--dimensions--dimension"></a>
### Nested Schema for `data_extraction.processors.processor.sampling_aware_value_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--data_extraction--processors--processor--security_context"></a>
### Nested Schema for `data_extraction.processors.processor.security_context`

Required:

- `value` (Block List, Min: 1, Max: 1) Security context value assignment (see [below for nested schema](#nestedblock--data_extraction--processors--processor--security_context--value))

<a id="nestedblock--data_extraction--processors--processor--security_context--value"></a>
### Nested Schema for `data_extraction.processors.processor.security_context.value`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

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




<a id="nestedblock--data_extraction--processors--processor--security_event"></a>
### Nested Schema for `data_extraction.processors.processor.security_event`

Required:

- `field_extraction` (Block List, Min: 1, Max: 1) Field Extraction (see [below for nested schema](#nestedblock--data_extraction--processors--processor--security_event--field_extraction))

<a id="nestedblock--data_extraction--processors--processor--security_event--field_extraction"></a>
### Nested Schema for `data_extraction.processors.processor.security_event.field_extraction`

Required:

- `type` (String) Possible Values: `Exclude`, `Include`, `IncludeAll`

Optional:

- `exclude` (Set of String) Fields
- `include` (Block List, Max: 1) Fields (see [below for nested schema](#nestedblock--data_extraction--processors--processor--security_event--field_extraction--include))

<a id="nestedblock--data_extraction--processors--processor--security_event--field_extraction--include"></a>
### Nested Schema for `data_extraction.processors.processor.security_event.field_extraction.include`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--data_extraction--processors--processor--security_event--field_extraction--include--dimension))

<a id="nestedblock--data_extraction--processors--processor--security_event--field_extraction--include--dimension"></a>
### Nested Schema for `data_extraction.processors.processor.security_event.field_extraction.include.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name





<a id="nestedblock--data_extraction--processors--processor--technology"></a>
### Nested Schema for `data_extraction.processors.processor.technology`

Required:

- `technology_id` (String) Technology ID

Optional:

- `custom_matcher` (String) Custom matching condition which should be used instead of technology matcher.


<a id="nestedblock--data_extraction--processors--processor--value_metric"></a>
### Nested Schema for `data_extraction.processors.processor.value_metric`

Required:

- `field` (String) Field with metric value
- `metric_key` (String) Metric key

Optional:

- `default_value` (String) Default value with metric value
- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--data_extraction--processors--processor--value_metric--dimensions))

<a id="nestedblock--data_extraction--processors--processor--value_metric--dimensions"></a>
### Nested Schema for `data_extraction.processors.processor.value_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--data_extraction--processors--processor--value_metric--dimensions--dimension))

<a id="nestedblock--data_extraction--processors--processor--value_metric--dimensions--dimension"></a>
### Nested Schema for `data_extraction.processors.processor.value_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name







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

- `description` (String) no documentation available
- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)
- `id` (String) Processor identifier
- `type` (String) Possible Values: `AzureLogForwarding`, `Bizevent`, `BucketAssignment`, `CostAllocation`, `CounterMetric`, `Davis`, `Dql`, `Drop`, `FieldsAdd`, `FieldsRemove`, `FieldsRename`, `HistogramMetric`, `NoStorage`, `ProductAllocation`, `SamplingAwareCounterMetric`, `SamplingAwareValueMetric`, `SecurityContext`, `SecurityEvent`, `Technology`, `ValueMetric`

Optional:

- `azure_log_forwarding` (Block List, Max: 1) Azure log forwarding processor attributes (see [below for nested schema](#nestedblock--davis--processors--processor--azure_log_forwarding))
- `bizevent` (Block List, Max: 1) Bizevent extraction processor attributes (see [below for nested schema](#nestedblock--davis--processors--processor--bizevent))
- `bucket_assignment` (Block List, Max: 1) Bucket assignment processor attributes (see [below for nested schema](#nestedblock--davis--processors--processor--bucket_assignment))
- `cost_allocation` (Block List, Max: 1) Cost allocation processor attributes (see [below for nested schema](#nestedblock--davis--processors--processor--cost_allocation))
- `counter_metric` (Block List, Max: 1) Counter metric processor attributes (see [below for nested schema](#nestedblock--davis--processors--processor--counter_metric))
- `davis` (Block List, Max: 1) Davis event extraction processor attributes (see [below for nested schema](#nestedblock--davis--processors--processor--davis))
- `dql` (Block List, Max: 1) DQL processor attributes (see [below for nested schema](#nestedblock--davis--processors--processor--dql))
- `fields_add` (Block List, Max: 1) Fields add processor attributes (see [below for nested schema](#nestedblock--davis--processors--processor--fields_add))
- `fields_remove` (Block List, Max: 1) Fields remove processor attributes (see [below for nested schema](#nestedblock--davis--processors--processor--fields_remove))
- `fields_rename` (Block List, Max: 1) Fields rename processor attributes (see [below for nested schema](#nestedblock--davis--processors--processor--fields_rename))
- `histogram_metric` (Block List, Max: 1) Histogram metric processor attributes (see [below for nested schema](#nestedblock--davis--processors--processor--histogram_metric))
- `matcher` (String) [See our documentation](https://dt-url.net/bp234rv)
- `product_allocation` (Block List, Max: 1) Product allocation processor attributes (see [below for nested schema](#nestedblock--davis--processors--processor--product_allocation))
- `sample_data` (String) Sample data
- `sampling_aware_counter_metric` (Block List, Max: 1) Sampling-aware counter metric processor attributes (see [below for nested schema](#nestedblock--davis--processors--processor--sampling_aware_counter_metric))
- `sampling_aware_value_metric` (Block List, Max: 1) Sampling aware value metric processor attributes (see [below for nested schema](#nestedblock--davis--processors--processor--sampling_aware_value_metric))
- `security_context` (Block List, Max: 1) Security context processor attributes (see [below for nested schema](#nestedblock--davis--processors--processor--security_context))
- `security_event` (Block List, Max: 1) Security event extraction processor attributes (see [below for nested schema](#nestedblock--davis--processors--processor--security_event))
- `technology` (Block List, Max: 1) Technology processor attributes (see [below for nested schema](#nestedblock--davis--processors--processor--technology))
- `value_metric` (Block List, Max: 1) Value metric processor attributes (see [below for nested schema](#nestedblock--davis--processors--processor--value_metric))

<a id="nestedblock--davis--processors--processor--azure_log_forwarding"></a>
### Nested Schema for `davis.processors.processor.azure_log_forwarding`

Required:

- `field_extraction` (Block List, Min: 1, Max: 1) Field Extraction (see [below for nested schema](#nestedblock--davis--processors--processor--azure_log_forwarding--field_extraction))
- `forwarder_config_id` (String) no documentation available

<a id="nestedblock--davis--processors--processor--azure_log_forwarding--field_extraction"></a>
### Nested Schema for `davis.processors.processor.azure_log_forwarding.field_extraction`

Required:

- `type` (String) Possible Values: `Exclude`, `Include`, `IncludeAll`

Optional:

- `exclude` (Set of String) Fields
- `include` (Block List, Max: 1) Fields (see [below for nested schema](#nestedblock--davis--processors--processor--azure_log_forwarding--field_extraction--include))

<a id="nestedblock--davis--processors--processor--azure_log_forwarding--field_extraction--include"></a>
### Nested Schema for `davis.processors.processor.azure_log_forwarding.field_extraction.include`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--davis--processors--processor--azure_log_forwarding--field_extraction--include--dimension))

<a id="nestedblock--davis--processors--processor--azure_log_forwarding--field_extraction--include--dimension"></a>
### Nested Schema for `davis.processors.processor.azure_log_forwarding.field_extraction.include.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name





<a id="nestedblock--davis--processors--processor--bizevent"></a>
### Nested Schema for `davis.processors.processor.bizevent`

Required:

- `event_provider` (Block List, Min: 1, Max: 1) Event provider (see [below for nested schema](#nestedblock--davis--processors--processor--bizevent--event_provider))
- `field_extraction` (Block List, Min: 1, Max: 1) Field extraction (see [below for nested schema](#nestedblock--davis--processors--processor--bizevent--field_extraction))

Optional:

- `event_type` (Block List, Max: 1) Event type (see [below for nested schema](#nestedblock--davis--processors--processor--bizevent--event_type))

<a id="nestedblock--davis--processors--processor--bizevent--event_provider"></a>
### Nested Schema for `davis.processors.processor.bizevent.event_provider`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--davis--processors--processor--bizevent--event_provider--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--davis--processors--processor--bizevent--event_provider--field"></a>
### Nested Schema for `davis.processors.processor.bizevent.event_provider.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value



<a id="nestedblock--davis--processors--processor--bizevent--field_extraction"></a>
### Nested Schema for `davis.processors.processor.bizevent.field_extraction`

Required:

- `type` (String) Possible Values: `Exclude`, `Include`, `IncludeAll`

Optional:

- `exclude` (Set of String) Fields
- `include` (Block List, Max: 1) Fields (see [below for nested schema](#nestedblock--davis--processors--processor--bizevent--field_extraction--include))

<a id="nestedblock--davis--processors--processor--bizevent--field_extraction--include"></a>
### Nested Schema for `davis.processors.processor.bizevent.field_extraction.include`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--davis--processors--processor--bizevent--field_extraction--include--dimension))

<a id="nestedblock--davis--processors--processor--bizevent--field_extraction--include--dimension"></a>
### Nested Schema for `davis.processors.processor.bizevent.field_extraction.include.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--davis--processors--processor--bizevent--event_type"></a>
### Nested Schema for `davis.processors.processor.bizevent.event_type`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--davis--processors--processor--bizevent--event_type--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--davis--processors--processor--bizevent--event_type--field"></a>
### Nested Schema for `davis.processors.processor.bizevent.event_type.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value




<a id="nestedblock--davis--processors--processor--bucket_assignment"></a>
### Nested Schema for `davis.processors.processor.bucket_assignment`

Required:

- `bucket_name` (String) Bucket name


<a id="nestedblock--davis--processors--processor--cost_allocation"></a>
### Nested Schema for `davis.processors.processor.cost_allocation`

Required:

- `value` (Block List, Min: 1, Max: 1) The strategy to set the cost allocation field (see [below for nested schema](#nestedblock--davis--processors--processor--cost_allocation--value))

<a id="nestedblock--davis--processors--processor--cost_allocation--value"></a>
### Nested Schema for `davis.processors.processor.cost_allocation.value`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--davis--processors--processor--cost_allocation--value--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--davis--processors--processor--cost_allocation--value--field"></a>
### Nested Schema for `davis.processors.processor.cost_allocation.value.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value




<a id="nestedblock--davis--processors--processor--counter_metric"></a>
### Nested Schema for `davis.processors.processor.counter_metric`

Required:

- `metric_key` (String) Metric key

Optional:

- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--davis--processors--processor--counter_metric--dimensions))

<a id="nestedblock--davis--processors--processor--counter_metric--dimensions"></a>
### Nested Schema for `davis.processors.processor.counter_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--davis--processors--processor--counter_metric--dimensions--dimension))

<a id="nestedblock--davis--processors--processor--counter_metric--dimensions--dimension"></a>
### Nested Schema for `davis.processors.processor.counter_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--davis--processors--processor--davis"></a>
### Nested Schema for `davis.processors.processor.davis`

Required:

- `properties` (Block List, Min: 1, Max: 1) no documentation available (see [below for nested schema](#nestedblock--davis--processors--processor--davis--properties))

<a id="nestedblock--davis--processors--processor--davis--properties"></a>
### Nested Schema for `davis.processors.processor.davis.properties`

Required:

- `property` (Block List, Min: 1) (see [below for nested schema](#nestedblock--davis--processors--processor--davis--properties--property))

<a id="nestedblock--davis--processors--processor--davis--properties--property"></a>
### Nested Schema for `davis.processors.processor.davis.properties.property`

Required:

- `key` (String) no documentation available
- `value` (String) no documentation available




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




<a id="nestedblock--davis--processors--processor--histogram_metric"></a>
### Nested Schema for `davis.processors.processor.histogram_metric`

Required:

- `field` (String) Field with metric value
- `metric_key` (String) Metric key

Optional:

- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--davis--processors--processor--histogram_metric--dimensions))

<a id="nestedblock--davis--processors--processor--histogram_metric--dimensions"></a>
### Nested Schema for `davis.processors.processor.histogram_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--davis--processors--processor--histogram_metric--dimensions--dimension))

<a id="nestedblock--davis--processors--processor--histogram_metric--dimensions--dimension"></a>
### Nested Schema for `davis.processors.processor.histogram_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--davis--processors--processor--product_allocation"></a>
### Nested Schema for `davis.processors.processor.product_allocation`

Required:

- `value` (Block List, Min: 1, Max: 1) The strategy to set product allocation field (see [below for nested schema](#nestedblock--davis--processors--processor--product_allocation--value))

<a id="nestedblock--davis--processors--processor--product_allocation--value"></a>
### Nested Schema for `davis.processors.processor.product_allocation.value`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--davis--processors--processor--product_allocation--value--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--davis--processors--processor--product_allocation--value--field"></a>
### Nested Schema for `davis.processors.processor.product_allocation.value.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value




<a id="nestedblock--davis--processors--processor--sampling_aware_counter_metric"></a>
### Nested Schema for `davis.processors.processor.sampling_aware_counter_metric`

Required:

- `metric_key` (String) Metric key

Optional:

- `aggregation` (String) Possible Values: `Disabled`, `Enabled`
- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--davis--processors--processor--sampling_aware_counter_metric--dimensions))
- `sampling` (String) Possible Values: `Disabled`, `Enabled`

<a id="nestedblock--davis--processors--processor--sampling_aware_counter_metric--dimensions"></a>
### Nested Schema for `davis.processors.processor.sampling_aware_counter_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--davis--processors--processor--sampling_aware_counter_metric--dimensions--dimension))

<a id="nestedblock--davis--processors--processor--sampling_aware_counter_metric--dimensions--dimension"></a>
### Nested Schema for `davis.processors.processor.sampling_aware_counter_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--davis--processors--processor--sampling_aware_value_metric"></a>
### Nested Schema for `davis.processors.processor.sampling_aware_value_metric`

Required:

- `measurement` (String) Possible Values: `Duration`, `Field`
- `metric_key` (String) Metric key

Optional:

- `aggregation` (String) Possible Values: `Disabled`, `Enabled`
- `default_value` (String) Default value with metric value
- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--davis--processors--processor--sampling_aware_value_metric--dimensions))
- `field` (String) Field with metric value
- `sampling` (String) Possible Values: `Disabled`, `Enabled`

<a id="nestedblock--davis--processors--processor--sampling_aware_value_metric--dimensions"></a>
### Nested Schema for `davis.processors.processor.sampling_aware_value_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--davis--processors--processor--sampling_aware_value_metric--dimensions--dimension))

<a id="nestedblock--davis--processors--processor--sampling_aware_value_metric--dimensions--dimension"></a>
### Nested Schema for `davis.processors.processor.sampling_aware_value_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--davis--processors--processor--security_context"></a>
### Nested Schema for `davis.processors.processor.security_context`

Required:

- `value` (Block List, Min: 1, Max: 1) Security context value assignment (see [below for nested schema](#nestedblock--davis--processors--processor--security_context--value))

<a id="nestedblock--davis--processors--processor--security_context--value"></a>
### Nested Schema for `davis.processors.processor.security_context.value`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

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




<a id="nestedblock--davis--processors--processor--security_event"></a>
### Nested Schema for `davis.processors.processor.security_event`

Required:

- `field_extraction` (Block List, Min: 1, Max: 1) Field Extraction (see [below for nested schema](#nestedblock--davis--processors--processor--security_event--field_extraction))

<a id="nestedblock--davis--processors--processor--security_event--field_extraction"></a>
### Nested Schema for `davis.processors.processor.security_event.field_extraction`

Required:

- `type` (String) Possible Values: `Exclude`, `Include`, `IncludeAll`

Optional:

- `exclude` (Set of String) Fields
- `include` (Block List, Max: 1) Fields (see [below for nested schema](#nestedblock--davis--processors--processor--security_event--field_extraction--include))

<a id="nestedblock--davis--processors--processor--security_event--field_extraction--include"></a>
### Nested Schema for `davis.processors.processor.security_event.field_extraction.include`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--davis--processors--processor--security_event--field_extraction--include--dimension))

<a id="nestedblock--davis--processors--processor--security_event--field_extraction--include--dimension"></a>
### Nested Schema for `davis.processors.processor.security_event.field_extraction.include.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name





<a id="nestedblock--davis--processors--processor--technology"></a>
### Nested Schema for `davis.processors.processor.technology`

Required:

- `technology_id` (String) Technology ID

Optional:

- `custom_matcher` (String) Custom matching condition which should be used instead of technology matcher.


<a id="nestedblock--davis--processors--processor--value_metric"></a>
### Nested Schema for `davis.processors.processor.value_metric`

Required:

- `field` (String) Field with metric value
- `metric_key` (String) Metric key

Optional:

- `default_value` (String) Default value with metric value
- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--davis--processors--processor--value_metric--dimensions))

<a id="nestedblock--davis--processors--processor--value_metric--dimensions"></a>
### Nested Schema for `davis.processors.processor.value_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--davis--processors--processor--value_metric--dimensions--dimension))

<a id="nestedblock--davis--processors--processor--value_metric--dimensions--dimension"></a>
### Nested Schema for `davis.processors.processor.value_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name







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

- `description` (String) no documentation available
- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)
- `id` (String) Processor identifier
- `type` (String) Possible Values: `AzureLogForwarding`, `Bizevent`, `BucketAssignment`, `CostAllocation`, `CounterMetric`, `Davis`, `Dql`, `Drop`, `FieldsAdd`, `FieldsRemove`, `FieldsRename`, `HistogramMetric`, `NoStorage`, `ProductAllocation`, `SamplingAwareCounterMetric`, `SamplingAwareValueMetric`, `SecurityContext`, `SecurityEvent`, `Technology`, `ValueMetric`

Optional:

- `azure_log_forwarding` (Block List, Max: 1) Azure log forwarding processor attributes (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--azure_log_forwarding))
- `bizevent` (Block List, Max: 1) Bizevent extraction processor attributes (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--bizevent))
- `bucket_assignment` (Block List, Max: 1) Bucket assignment processor attributes (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--bucket_assignment))
- `cost_allocation` (Block List, Max: 1) Cost allocation processor attributes (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--cost_allocation))
- `counter_metric` (Block List, Max: 1) Counter metric processor attributes (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--counter_metric))
- `davis` (Block List, Max: 1) Davis event extraction processor attributes (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--davis))
- `dql` (Block List, Max: 1) DQL processor attributes (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--dql))
- `fields_add` (Block List, Max: 1) Fields add processor attributes (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--fields_add))
- `fields_remove` (Block List, Max: 1) Fields remove processor attributes (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--fields_remove))
- `fields_rename` (Block List, Max: 1) Fields rename processor attributes (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--fields_rename))
- `histogram_metric` (Block List, Max: 1) Histogram metric processor attributes (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--histogram_metric))
- `matcher` (String) [See our documentation](https://dt-url.net/bp234rv)
- `product_allocation` (Block List, Max: 1) Product allocation processor attributes (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--product_allocation))
- `sample_data` (String) Sample data
- `sampling_aware_counter_metric` (Block List, Max: 1) Sampling-aware counter metric processor attributes (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--sampling_aware_counter_metric))
- `sampling_aware_value_metric` (Block List, Max: 1) Sampling aware value metric processor attributes (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--sampling_aware_value_metric))
- `security_context` (Block List, Max: 1) Security context processor attributes (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--security_context))
- `security_event` (Block List, Max: 1) Security event extraction processor attributes (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--security_event))
- `technology` (Block List, Max: 1) Technology processor attributes (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--technology))
- `value_metric` (Block List, Max: 1) Value metric processor attributes (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--value_metric))

<a id="nestedblock--metric_extraction--processors--processor--azure_log_forwarding"></a>
### Nested Schema for `metric_extraction.processors.processor.azure_log_forwarding`

Required:

- `field_extraction` (Block List, Min: 1, Max: 1) Field Extraction (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--azure_log_forwarding--field_extraction))
- `forwarder_config_id` (String) no documentation available

<a id="nestedblock--metric_extraction--processors--processor--azure_log_forwarding--field_extraction"></a>
### Nested Schema for `metric_extraction.processors.processor.azure_log_forwarding.field_extraction`

Required:

- `type` (String) Possible Values: `Exclude`, `Include`, `IncludeAll`

Optional:

- `exclude` (Set of String) Fields
- `include` (Block List, Max: 1) Fields (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--azure_log_forwarding--field_extraction--include))

<a id="nestedblock--metric_extraction--processors--processor--azure_log_forwarding--field_extraction--include"></a>
### Nested Schema for `metric_extraction.processors.processor.azure_log_forwarding.field_extraction.include`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--azure_log_forwarding--field_extraction--include--dimension))

<a id="nestedblock--metric_extraction--processors--processor--azure_log_forwarding--field_extraction--include--dimension"></a>
### Nested Schema for `metric_extraction.processors.processor.azure_log_forwarding.field_extraction.include.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name





<a id="nestedblock--metric_extraction--processors--processor--bizevent"></a>
### Nested Schema for `metric_extraction.processors.processor.bizevent`

Required:

- `event_provider` (Block List, Min: 1, Max: 1) Event provider (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--bizevent--event_provider))
- `field_extraction` (Block List, Min: 1, Max: 1) Field extraction (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--bizevent--field_extraction))

Optional:

- `event_type` (Block List, Max: 1) Event type (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--bizevent--event_type))

<a id="nestedblock--metric_extraction--processors--processor--bizevent--event_provider"></a>
### Nested Schema for `metric_extraction.processors.processor.bizevent.event_provider`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--bizevent--event_provider--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--metric_extraction--processors--processor--bizevent--event_provider--field"></a>
### Nested Schema for `metric_extraction.processors.processor.bizevent.event_provider.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value



<a id="nestedblock--metric_extraction--processors--processor--bizevent--field_extraction"></a>
### Nested Schema for `metric_extraction.processors.processor.bizevent.field_extraction`

Required:

- `type` (String) Possible Values: `Exclude`, `Include`, `IncludeAll`

Optional:

- `exclude` (Set of String) Fields
- `include` (Block List, Max: 1) Fields (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--bizevent--field_extraction--include))

<a id="nestedblock--metric_extraction--processors--processor--bizevent--field_extraction--include"></a>
### Nested Schema for `metric_extraction.processors.processor.bizevent.field_extraction.include`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--bizevent--field_extraction--include--dimension))

<a id="nestedblock--metric_extraction--processors--processor--bizevent--field_extraction--include--dimension"></a>
### Nested Schema for `metric_extraction.processors.processor.bizevent.field_extraction.include.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--metric_extraction--processors--processor--bizevent--event_type"></a>
### Nested Schema for `metric_extraction.processors.processor.bizevent.event_type`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--bizevent--event_type--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--metric_extraction--processors--processor--bizevent--event_type--field"></a>
### Nested Schema for `metric_extraction.processors.processor.bizevent.event_type.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value




<a id="nestedblock--metric_extraction--processors--processor--bucket_assignment"></a>
### Nested Schema for `metric_extraction.processors.processor.bucket_assignment`

Required:

- `bucket_name` (String) Bucket name


<a id="nestedblock--metric_extraction--processors--processor--cost_allocation"></a>
### Nested Schema for `metric_extraction.processors.processor.cost_allocation`

Required:

- `value` (Block List, Min: 1, Max: 1) The strategy to set the cost allocation field (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--cost_allocation--value))

<a id="nestedblock--metric_extraction--processors--processor--cost_allocation--value"></a>
### Nested Schema for `metric_extraction.processors.processor.cost_allocation.value`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--cost_allocation--value--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--metric_extraction--processors--processor--cost_allocation--value--field"></a>
### Nested Schema for `metric_extraction.processors.processor.cost_allocation.value.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value




<a id="nestedblock--metric_extraction--processors--processor--counter_metric"></a>
### Nested Schema for `metric_extraction.processors.processor.counter_metric`

Required:

- `metric_key` (String) Metric key

Optional:

- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--counter_metric--dimensions))

<a id="nestedblock--metric_extraction--processors--processor--counter_metric--dimensions"></a>
### Nested Schema for `metric_extraction.processors.processor.counter_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--counter_metric--dimensions--dimension))

<a id="nestedblock--metric_extraction--processors--processor--counter_metric--dimensions--dimension"></a>
### Nested Schema for `metric_extraction.processors.processor.counter_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--metric_extraction--processors--processor--davis"></a>
### Nested Schema for `metric_extraction.processors.processor.davis`

Required:

- `properties` (Block List, Min: 1, Max: 1) no documentation available (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--davis--properties))

<a id="nestedblock--metric_extraction--processors--processor--davis--properties"></a>
### Nested Schema for `metric_extraction.processors.processor.davis.properties`

Required:

- `property` (Block List, Min: 1) (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--davis--properties--property))

<a id="nestedblock--metric_extraction--processors--processor--davis--properties--property"></a>
### Nested Schema for `metric_extraction.processors.processor.davis.properties.property`

Required:

- `key` (String) no documentation available
- `value` (String) no documentation available




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




<a id="nestedblock--metric_extraction--processors--processor--histogram_metric"></a>
### Nested Schema for `metric_extraction.processors.processor.histogram_metric`

Required:

- `field` (String) Field with metric value
- `metric_key` (String) Metric key

Optional:

- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--histogram_metric--dimensions))

<a id="nestedblock--metric_extraction--processors--processor--histogram_metric--dimensions"></a>
### Nested Schema for `metric_extraction.processors.processor.histogram_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--histogram_metric--dimensions--dimension))

<a id="nestedblock--metric_extraction--processors--processor--histogram_metric--dimensions--dimension"></a>
### Nested Schema for `metric_extraction.processors.processor.histogram_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--metric_extraction--processors--processor--product_allocation"></a>
### Nested Schema for `metric_extraction.processors.processor.product_allocation`

Required:

- `value` (Block List, Min: 1, Max: 1) The strategy to set product allocation field (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--product_allocation--value))

<a id="nestedblock--metric_extraction--processors--processor--product_allocation--value"></a>
### Nested Schema for `metric_extraction.processors.processor.product_allocation.value`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--product_allocation--value--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--metric_extraction--processors--processor--product_allocation--value--field"></a>
### Nested Schema for `metric_extraction.processors.processor.product_allocation.value.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value




<a id="nestedblock--metric_extraction--processors--processor--sampling_aware_counter_metric"></a>
### Nested Schema for `metric_extraction.processors.processor.sampling_aware_counter_metric`

Required:

- `metric_key` (String) Metric key

Optional:

- `aggregation` (String) Possible Values: `Disabled`, `Enabled`
- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--sampling_aware_counter_metric--dimensions))
- `sampling` (String) Possible Values: `Disabled`, `Enabled`

<a id="nestedblock--metric_extraction--processors--processor--sampling_aware_counter_metric--dimensions"></a>
### Nested Schema for `metric_extraction.processors.processor.sampling_aware_counter_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--sampling_aware_counter_metric--dimensions--dimension))

<a id="nestedblock--metric_extraction--processors--processor--sampling_aware_counter_metric--dimensions--dimension"></a>
### Nested Schema for `metric_extraction.processors.processor.sampling_aware_counter_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--metric_extraction--processors--processor--sampling_aware_value_metric"></a>
### Nested Schema for `metric_extraction.processors.processor.sampling_aware_value_metric`

Required:

- `measurement` (String) Possible Values: `Duration`, `Field`
- `metric_key` (String) Metric key

Optional:

- `aggregation` (String) Possible Values: `Disabled`, `Enabled`
- `default_value` (String) Default value with metric value
- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--sampling_aware_value_metric--dimensions))
- `field` (String) Field with metric value
- `sampling` (String) Possible Values: `Disabled`, `Enabled`

<a id="nestedblock--metric_extraction--processors--processor--sampling_aware_value_metric--dimensions"></a>
### Nested Schema for `metric_extraction.processors.processor.sampling_aware_value_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--sampling_aware_value_metric--dimensions--dimension))

<a id="nestedblock--metric_extraction--processors--processor--sampling_aware_value_metric--dimensions--dimension"></a>
### Nested Schema for `metric_extraction.processors.processor.sampling_aware_value_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--metric_extraction--processors--processor--security_context"></a>
### Nested Schema for `metric_extraction.processors.processor.security_context`

Required:

- `value` (Block List, Min: 1, Max: 1) Security context value assignment (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--security_context--value))

<a id="nestedblock--metric_extraction--processors--processor--security_context--value"></a>
### Nested Schema for `metric_extraction.processors.processor.security_context.value`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

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




<a id="nestedblock--metric_extraction--processors--processor--security_event"></a>
### Nested Schema for `metric_extraction.processors.processor.security_event`

Required:

- `field_extraction` (Block List, Min: 1, Max: 1) Field Extraction (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--security_event--field_extraction))

<a id="nestedblock--metric_extraction--processors--processor--security_event--field_extraction"></a>
### Nested Schema for `metric_extraction.processors.processor.security_event.field_extraction`

Required:

- `type` (String) Possible Values: `Exclude`, `Include`, `IncludeAll`

Optional:

- `exclude` (Set of String) Fields
- `include` (Block List, Max: 1) Fields (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--security_event--field_extraction--include))

<a id="nestedblock--metric_extraction--processors--processor--security_event--field_extraction--include"></a>
### Nested Schema for `metric_extraction.processors.processor.security_event.field_extraction.include`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--security_event--field_extraction--include--dimension))

<a id="nestedblock--metric_extraction--processors--processor--security_event--field_extraction--include--dimension"></a>
### Nested Schema for `metric_extraction.processors.processor.security_event.field_extraction.include.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name





<a id="nestedblock--metric_extraction--processors--processor--technology"></a>
### Nested Schema for `metric_extraction.processors.processor.technology`

Required:

- `technology_id` (String) Technology ID

Optional:

- `custom_matcher` (String) Custom matching condition which should be used instead of technology matcher.


<a id="nestedblock--metric_extraction--processors--processor--value_metric"></a>
### Nested Schema for `metric_extraction.processors.processor.value_metric`

Required:

- `field` (String) Field with metric value
- `metric_key` (String) Metric key

Optional:

- `default_value` (String) Default value with metric value
- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--value_metric--dimensions))

<a id="nestedblock--metric_extraction--processors--processor--value_metric--dimensions"></a>
### Nested Schema for `metric_extraction.processors.processor.value_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--metric_extraction--processors--processor--value_metric--dimensions--dimension))

<a id="nestedblock--metric_extraction--processors--processor--value_metric--dimensions--dimension"></a>
### Nested Schema for `metric_extraction.processors.processor.value_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name







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

- `description` (String) no documentation available
- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)
- `id` (String) Processor identifier
- `type` (String) Possible Values: `AzureLogForwarding`, `Bizevent`, `BucketAssignment`, `CostAllocation`, `CounterMetric`, `Davis`, `Dql`, `Drop`, `FieldsAdd`, `FieldsRemove`, `FieldsRename`, `HistogramMetric`, `NoStorage`, `ProductAllocation`, `SamplingAwareCounterMetric`, `SamplingAwareValueMetric`, `SecurityContext`, `SecurityEvent`, `Technology`, `ValueMetric`

Optional:

- `azure_log_forwarding` (Block List, Max: 1) Azure log forwarding processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--azure_log_forwarding))
- `bizevent` (Block List, Max: 1) Bizevent extraction processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--bizevent))
- `bucket_assignment` (Block List, Max: 1) Bucket assignment processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--bucket_assignment))
- `cost_allocation` (Block List, Max: 1) Cost allocation processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--cost_allocation))
- `counter_metric` (Block List, Max: 1) Counter metric processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--counter_metric))
- `davis` (Block List, Max: 1) Davis event extraction processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--davis))
- `dql` (Block List, Max: 1) DQL processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--dql))
- `fields_add` (Block List, Max: 1) Fields add processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--fields_add))
- `fields_remove` (Block List, Max: 1) Fields remove processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--fields_remove))
- `fields_rename` (Block List, Max: 1) Fields rename processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--fields_rename))
- `histogram_metric` (Block List, Max: 1) Histogram metric processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--histogram_metric))
- `matcher` (String) [See our documentation](https://dt-url.net/bp234rv)
- `product_allocation` (Block List, Max: 1) Product allocation processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--product_allocation))
- `sample_data` (String) Sample data
- `sampling_aware_counter_metric` (Block List, Max: 1) Sampling-aware counter metric processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--sampling_aware_counter_metric))
- `sampling_aware_value_metric` (Block List, Max: 1) Sampling aware value metric processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--sampling_aware_value_metric))
- `security_context` (Block List, Max: 1) Security context processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--security_context))
- `security_event` (Block List, Max: 1) Security event extraction processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--security_event))
- `technology` (Block List, Max: 1) Technology processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--technology))
- `value_metric` (Block List, Max: 1) Value metric processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--value_metric))

<a id="nestedblock--processing--processors--processor--azure_log_forwarding"></a>
### Nested Schema for `processing.processors.processor.azure_log_forwarding`

Required:

- `field_extraction` (Block List, Min: 1, Max: 1) Field Extraction (see [below for nested schema](#nestedblock--processing--processors--processor--azure_log_forwarding--field_extraction))
- `forwarder_config_id` (String) no documentation available

<a id="nestedblock--processing--processors--processor--azure_log_forwarding--field_extraction"></a>
### Nested Schema for `processing.processors.processor.azure_log_forwarding.field_extraction`

Required:

- `type` (String) Possible Values: `Exclude`, `Include`, `IncludeAll`

Optional:

- `exclude` (Set of String) Fields
- `include` (Block List, Max: 1) Fields (see [below for nested schema](#nestedblock--processing--processors--processor--azure_log_forwarding--field_extraction--include))

<a id="nestedblock--processing--processors--processor--azure_log_forwarding--field_extraction--include"></a>
### Nested Schema for `processing.processors.processor.azure_log_forwarding.field_extraction.include`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--processing--processors--processor--azure_log_forwarding--field_extraction--include--dimension))

<a id="nestedblock--processing--processors--processor--azure_log_forwarding--field_extraction--include--dimension"></a>
### Nested Schema for `processing.processors.processor.azure_log_forwarding.field_extraction.include.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name





<a id="nestedblock--processing--processors--processor--bizevent"></a>
### Nested Schema for `processing.processors.processor.bizevent`

Required:

- `event_provider` (Block List, Min: 1, Max: 1) Event provider (see [below for nested schema](#nestedblock--processing--processors--processor--bizevent--event_provider))
- `field_extraction` (Block List, Min: 1, Max: 1) Field extraction (see [below for nested schema](#nestedblock--processing--processors--processor--bizevent--field_extraction))

Optional:

- `event_type` (Block List, Max: 1) Event type (see [below for nested schema](#nestedblock--processing--processors--processor--bizevent--event_type))

<a id="nestedblock--processing--processors--processor--bizevent--event_provider"></a>
### Nested Schema for `processing.processors.processor.bizevent.event_provider`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--processing--processors--processor--bizevent--event_provider--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--processing--processors--processor--bizevent--event_provider--field"></a>
### Nested Schema for `processing.processors.processor.bizevent.event_provider.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value



<a id="nestedblock--processing--processors--processor--bizevent--field_extraction"></a>
### Nested Schema for `processing.processors.processor.bizevent.field_extraction`

Required:

- `type` (String) Possible Values: `Exclude`, `Include`, `IncludeAll`

Optional:

- `exclude` (Set of String) Fields
- `include` (Block List, Max: 1) Fields (see [below for nested schema](#nestedblock--processing--processors--processor--bizevent--field_extraction--include))

<a id="nestedblock--processing--processors--processor--bizevent--field_extraction--include"></a>
### Nested Schema for `processing.processors.processor.bizevent.field_extraction.include`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--processing--processors--processor--bizevent--field_extraction--include--dimension))

<a id="nestedblock--processing--processors--processor--bizevent--field_extraction--include--dimension"></a>
### Nested Schema for `processing.processors.processor.bizevent.field_extraction.include.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--processing--processors--processor--bizevent--event_type"></a>
### Nested Schema for `processing.processors.processor.bizevent.event_type`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--processing--processors--processor--bizevent--event_type--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--processing--processors--processor--bizevent--event_type--field"></a>
### Nested Schema for `processing.processors.processor.bizevent.event_type.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value




<a id="nestedblock--processing--processors--processor--bucket_assignment"></a>
### Nested Schema for `processing.processors.processor.bucket_assignment`

Required:

- `bucket_name` (String) Bucket name


<a id="nestedblock--processing--processors--processor--cost_allocation"></a>
### Nested Schema for `processing.processors.processor.cost_allocation`

Required:

- `value` (Block List, Min: 1, Max: 1) The strategy to set the cost allocation field (see [below for nested schema](#nestedblock--processing--processors--processor--cost_allocation--value))

<a id="nestedblock--processing--processors--processor--cost_allocation--value"></a>
### Nested Schema for `processing.processors.processor.cost_allocation.value`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--processing--processors--processor--cost_allocation--value--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--processing--processors--processor--cost_allocation--value--field"></a>
### Nested Schema for `processing.processors.processor.cost_allocation.value.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value




<a id="nestedblock--processing--processors--processor--counter_metric"></a>
### Nested Schema for `processing.processors.processor.counter_metric`

Required:

- `metric_key` (String) Metric key

Optional:

- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--processing--processors--processor--counter_metric--dimensions))

<a id="nestedblock--processing--processors--processor--counter_metric--dimensions"></a>
### Nested Schema for `processing.processors.processor.counter_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--processing--processors--processor--counter_metric--dimensions--dimension))

<a id="nestedblock--processing--processors--processor--counter_metric--dimensions--dimension"></a>
### Nested Schema for `processing.processors.processor.counter_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--processing--processors--processor--davis"></a>
### Nested Schema for `processing.processors.processor.davis`

Required:

- `properties` (Block List, Min: 1, Max: 1) no documentation available (see [below for nested schema](#nestedblock--processing--processors--processor--davis--properties))

<a id="nestedblock--processing--processors--processor--davis--properties"></a>
### Nested Schema for `processing.processors.processor.davis.properties`

Required:

- `property` (Block List, Min: 1) (see [below for nested schema](#nestedblock--processing--processors--processor--davis--properties--property))

<a id="nestedblock--processing--processors--processor--davis--properties--property"></a>
### Nested Schema for `processing.processors.processor.davis.properties.property`

Required:

- `key` (String) no documentation available
- `value` (String) no documentation available




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




<a id="nestedblock--processing--processors--processor--histogram_metric"></a>
### Nested Schema for `processing.processors.processor.histogram_metric`

Required:

- `field` (String) Field with metric value
- `metric_key` (String) Metric key

Optional:

- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--processing--processors--processor--histogram_metric--dimensions))

<a id="nestedblock--processing--processors--processor--histogram_metric--dimensions"></a>
### Nested Schema for `processing.processors.processor.histogram_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--processing--processors--processor--histogram_metric--dimensions--dimension))

<a id="nestedblock--processing--processors--processor--histogram_metric--dimensions--dimension"></a>
### Nested Schema for `processing.processors.processor.histogram_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--processing--processors--processor--product_allocation"></a>
### Nested Schema for `processing.processors.processor.product_allocation`

Required:

- `value` (Block List, Min: 1, Max: 1) The strategy to set product allocation field (see [below for nested schema](#nestedblock--processing--processors--processor--product_allocation--value))

<a id="nestedblock--processing--processors--processor--product_allocation--value"></a>
### Nested Schema for `processing.processors.processor.product_allocation.value`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--processing--processors--processor--product_allocation--value--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--processing--processors--processor--product_allocation--value--field"></a>
### Nested Schema for `processing.processors.processor.product_allocation.value.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value




<a id="nestedblock--processing--processors--processor--sampling_aware_counter_metric"></a>
### Nested Schema for `processing.processors.processor.sampling_aware_counter_metric`

Required:

- `metric_key` (String) Metric key

Optional:

- `aggregation` (String) Possible Values: `Disabled`, `Enabled`
- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--processing--processors--processor--sampling_aware_counter_metric--dimensions))
- `sampling` (String) Possible Values: `Disabled`, `Enabled`

<a id="nestedblock--processing--processors--processor--sampling_aware_counter_metric--dimensions"></a>
### Nested Schema for `processing.processors.processor.sampling_aware_counter_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--processing--processors--processor--sampling_aware_counter_metric--dimensions--dimension))

<a id="nestedblock--processing--processors--processor--sampling_aware_counter_metric--dimensions--dimension"></a>
### Nested Schema for `processing.processors.processor.sampling_aware_counter_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--processing--processors--processor--sampling_aware_value_metric"></a>
### Nested Schema for `processing.processors.processor.sampling_aware_value_metric`

Required:

- `measurement` (String) Possible Values: `Duration`, `Field`
- `metric_key` (String) Metric key

Optional:

- `aggregation` (String) Possible Values: `Disabled`, `Enabled`
- `default_value` (String) Default value with metric value
- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--processing--processors--processor--sampling_aware_value_metric--dimensions))
- `field` (String) Field with metric value
- `sampling` (String) Possible Values: `Disabled`, `Enabled`

<a id="nestedblock--processing--processors--processor--sampling_aware_value_metric--dimensions"></a>
### Nested Schema for `processing.processors.processor.sampling_aware_value_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--processing--processors--processor--sampling_aware_value_metric--dimensions--dimension))

<a id="nestedblock--processing--processors--processor--sampling_aware_value_metric--dimensions--dimension"></a>
### Nested Schema for `processing.processors.processor.sampling_aware_value_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--processing--processors--processor--security_context"></a>
### Nested Schema for `processing.processors.processor.security_context`

Required:

- `value` (Block List, Min: 1, Max: 1) Security context value assignment (see [below for nested schema](#nestedblock--processing--processors--processor--security_context--value))

<a id="nestedblock--processing--processors--processor--security_context--value"></a>
### Nested Schema for `processing.processors.processor.security_context.value`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

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




<a id="nestedblock--processing--processors--processor--security_event"></a>
### Nested Schema for `processing.processors.processor.security_event`

Required:

- `field_extraction` (Block List, Min: 1, Max: 1) Field Extraction (see [below for nested schema](#nestedblock--processing--processors--processor--security_event--field_extraction))

<a id="nestedblock--processing--processors--processor--security_event--field_extraction"></a>
### Nested Schema for `processing.processors.processor.security_event.field_extraction`

Required:

- `type` (String) Possible Values: `Exclude`, `Include`, `IncludeAll`

Optional:

- `exclude` (Set of String) Fields
- `include` (Block List, Max: 1) Fields (see [below for nested schema](#nestedblock--processing--processors--processor--security_event--field_extraction--include))

<a id="nestedblock--processing--processors--processor--security_event--field_extraction--include"></a>
### Nested Schema for `processing.processors.processor.security_event.field_extraction.include`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--processing--processors--processor--security_event--field_extraction--include--dimension))

<a id="nestedblock--processing--processors--processor--security_event--field_extraction--include--dimension"></a>
### Nested Schema for `processing.processors.processor.security_event.field_extraction.include.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name





<a id="nestedblock--processing--processors--processor--technology"></a>
### Nested Schema for `processing.processors.processor.technology`

Required:

- `technology_id` (String) Technology ID

Optional:

- `custom_matcher` (String) Custom matching condition which should be used instead of technology matcher.


<a id="nestedblock--processing--processors--processor--value_metric"></a>
### Nested Schema for `processing.processors.processor.value_metric`

Required:

- `field` (String) Field with metric value
- `metric_key` (String) Metric key

Optional:

- `default_value` (String) Default value with metric value
- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--processing--processors--processor--value_metric--dimensions))

<a id="nestedblock--processing--processors--processor--value_metric--dimensions"></a>
### Nested Schema for `processing.processors.processor.value_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--processing--processors--processor--value_metric--dimensions--dimension))

<a id="nestedblock--processing--processors--processor--value_metric--dimensions--dimension"></a>
### Nested Schema for `processing.processors.processor.value_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name







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

- `description` (String) no documentation available
- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)
- `id` (String) Processor identifier
- `type` (String) Possible Values: `AzureLogForwarding`, `Bizevent`, `BucketAssignment`, `CostAllocation`, `CounterMetric`, `Davis`, `Dql`, `Drop`, `FieldsAdd`, `FieldsRemove`, `FieldsRename`, `HistogramMetric`, `NoStorage`, `ProductAllocation`, `SamplingAwareCounterMetric`, `SamplingAwareValueMetric`, `SecurityContext`, `SecurityEvent`, `Technology`, `ValueMetric`

Optional:

- `azure_log_forwarding` (Block List, Max: 1) Azure log forwarding processor attributes (see [below for nested schema](#nestedblock--product_allocation--processors--processor--azure_log_forwarding))
- `bizevent` (Block List, Max: 1) Bizevent extraction processor attributes (see [below for nested schema](#nestedblock--product_allocation--processors--processor--bizevent))
- `bucket_assignment` (Block List, Max: 1) Bucket assignment processor attributes (see [below for nested schema](#nestedblock--product_allocation--processors--processor--bucket_assignment))
- `cost_allocation` (Block List, Max: 1) Cost allocation processor attributes (see [below for nested schema](#nestedblock--product_allocation--processors--processor--cost_allocation))
- `counter_metric` (Block List, Max: 1) Counter metric processor attributes (see [below for nested schema](#nestedblock--product_allocation--processors--processor--counter_metric))
- `davis` (Block List, Max: 1) Davis event extraction processor attributes (see [below for nested schema](#nestedblock--product_allocation--processors--processor--davis))
- `dql` (Block List, Max: 1) DQL processor attributes (see [below for nested schema](#nestedblock--product_allocation--processors--processor--dql))
- `fields_add` (Block List, Max: 1) Fields add processor attributes (see [below for nested schema](#nestedblock--product_allocation--processors--processor--fields_add))
- `fields_remove` (Block List, Max: 1) Fields remove processor attributes (see [below for nested schema](#nestedblock--product_allocation--processors--processor--fields_remove))
- `fields_rename` (Block List, Max: 1) Fields rename processor attributes (see [below for nested schema](#nestedblock--product_allocation--processors--processor--fields_rename))
- `histogram_metric` (Block List, Max: 1) Histogram metric processor attributes (see [below for nested schema](#nestedblock--product_allocation--processors--processor--histogram_metric))
- `matcher` (String) [See our documentation](https://dt-url.net/bp234rv)
- `product_allocation` (Block List, Max: 1) Product allocation processor attributes (see [below for nested schema](#nestedblock--product_allocation--processors--processor--product_allocation))
- `sample_data` (String) Sample data
- `sampling_aware_counter_metric` (Block List, Max: 1) Sampling-aware counter metric processor attributes (see [below for nested schema](#nestedblock--product_allocation--processors--processor--sampling_aware_counter_metric))
- `sampling_aware_value_metric` (Block List, Max: 1) Sampling aware value metric processor attributes (see [below for nested schema](#nestedblock--product_allocation--processors--processor--sampling_aware_value_metric))
- `security_context` (Block List, Max: 1) Security context processor attributes (see [below for nested schema](#nestedblock--product_allocation--processors--processor--security_context))
- `security_event` (Block List, Max: 1) Security event extraction processor attributes (see [below for nested schema](#nestedblock--product_allocation--processors--processor--security_event))
- `technology` (Block List, Max: 1) Technology processor attributes (see [below for nested schema](#nestedblock--product_allocation--processors--processor--technology))
- `value_metric` (Block List, Max: 1) Value metric processor attributes (see [below for nested schema](#nestedblock--product_allocation--processors--processor--value_metric))

<a id="nestedblock--product_allocation--processors--processor--azure_log_forwarding"></a>
### Nested Schema for `product_allocation.processors.processor.azure_log_forwarding`

Required:

- `field_extraction` (Block List, Min: 1, Max: 1) Field Extraction (see [below for nested schema](#nestedblock--product_allocation--processors--processor--azure_log_forwarding--field_extraction))
- `forwarder_config_id` (String) no documentation available

<a id="nestedblock--product_allocation--processors--processor--azure_log_forwarding--field_extraction"></a>
### Nested Schema for `product_allocation.processors.processor.azure_log_forwarding.field_extraction`

Required:

- `type` (String) Possible Values: `Exclude`, `Include`, `IncludeAll`

Optional:

- `exclude` (Set of String) Fields
- `include` (Block List, Max: 1) Fields (see [below for nested schema](#nestedblock--product_allocation--processors--processor--azure_log_forwarding--field_extraction--include))

<a id="nestedblock--product_allocation--processors--processor--azure_log_forwarding--field_extraction--include"></a>
### Nested Schema for `product_allocation.processors.processor.azure_log_forwarding.field_extraction.include`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--product_allocation--processors--processor--azure_log_forwarding--field_extraction--include--dimension))

<a id="nestedblock--product_allocation--processors--processor--azure_log_forwarding--field_extraction--include--dimension"></a>
### Nested Schema for `product_allocation.processors.processor.azure_log_forwarding.field_extraction.include.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name





<a id="nestedblock--product_allocation--processors--processor--bizevent"></a>
### Nested Schema for `product_allocation.processors.processor.bizevent`

Required:

- `event_provider` (Block List, Min: 1, Max: 1) Event provider (see [below for nested schema](#nestedblock--product_allocation--processors--processor--bizevent--event_provider))
- `field_extraction` (Block List, Min: 1, Max: 1) Field extraction (see [below for nested schema](#nestedblock--product_allocation--processors--processor--bizevent--field_extraction))

Optional:

- `event_type` (Block List, Max: 1) Event type (see [below for nested schema](#nestedblock--product_allocation--processors--processor--bizevent--event_type))

<a id="nestedblock--product_allocation--processors--processor--bizevent--event_provider"></a>
### Nested Schema for `product_allocation.processors.processor.bizevent.event_provider`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--product_allocation--processors--processor--bizevent--event_provider--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--product_allocation--processors--processor--bizevent--event_provider--field"></a>
### Nested Schema for `product_allocation.processors.processor.bizevent.event_provider.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value



<a id="nestedblock--product_allocation--processors--processor--bizevent--field_extraction"></a>
### Nested Schema for `product_allocation.processors.processor.bizevent.field_extraction`

Required:

- `type` (String) Possible Values: `Exclude`, `Include`, `IncludeAll`

Optional:

- `exclude` (Set of String) Fields
- `include` (Block List, Max: 1) Fields (see [below for nested schema](#nestedblock--product_allocation--processors--processor--bizevent--field_extraction--include))

<a id="nestedblock--product_allocation--processors--processor--bizevent--field_extraction--include"></a>
### Nested Schema for `product_allocation.processors.processor.bizevent.field_extraction.include`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--product_allocation--processors--processor--bizevent--field_extraction--include--dimension))

<a id="nestedblock--product_allocation--processors--processor--bizevent--field_extraction--include--dimension"></a>
### Nested Schema for `product_allocation.processors.processor.bizevent.field_extraction.include.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--product_allocation--processors--processor--bizevent--event_type"></a>
### Nested Schema for `product_allocation.processors.processor.bizevent.event_type`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--product_allocation--processors--processor--bizevent--event_type--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--product_allocation--processors--processor--bizevent--event_type--field"></a>
### Nested Schema for `product_allocation.processors.processor.bizevent.event_type.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value




<a id="nestedblock--product_allocation--processors--processor--bucket_assignment"></a>
### Nested Schema for `product_allocation.processors.processor.bucket_assignment`

Required:

- `bucket_name` (String) Bucket name


<a id="nestedblock--product_allocation--processors--processor--cost_allocation"></a>
### Nested Schema for `product_allocation.processors.processor.cost_allocation`

Required:

- `value` (Block List, Min: 1, Max: 1) The strategy to set the cost allocation field (see [below for nested schema](#nestedblock--product_allocation--processors--processor--cost_allocation--value))

<a id="nestedblock--product_allocation--processors--processor--cost_allocation--value"></a>
### Nested Schema for `product_allocation.processors.processor.cost_allocation.value`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--product_allocation--processors--processor--cost_allocation--value--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--product_allocation--processors--processor--cost_allocation--value--field"></a>
### Nested Schema for `product_allocation.processors.processor.cost_allocation.value.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value




<a id="nestedblock--product_allocation--processors--processor--counter_metric"></a>
### Nested Schema for `product_allocation.processors.processor.counter_metric`

Required:

- `metric_key` (String) Metric key

Optional:

- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--product_allocation--processors--processor--counter_metric--dimensions))

<a id="nestedblock--product_allocation--processors--processor--counter_metric--dimensions"></a>
### Nested Schema for `product_allocation.processors.processor.counter_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--product_allocation--processors--processor--counter_metric--dimensions--dimension))

<a id="nestedblock--product_allocation--processors--processor--counter_metric--dimensions--dimension"></a>
### Nested Schema for `product_allocation.processors.processor.counter_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--product_allocation--processors--processor--davis"></a>
### Nested Schema for `product_allocation.processors.processor.davis`

Required:

- `properties` (Block List, Min: 1, Max: 1) no documentation available (see [below for nested schema](#nestedblock--product_allocation--processors--processor--davis--properties))

<a id="nestedblock--product_allocation--processors--processor--davis--properties"></a>
### Nested Schema for `product_allocation.processors.processor.davis.properties`

Required:

- `property` (Block List, Min: 1) (see [below for nested schema](#nestedblock--product_allocation--processors--processor--davis--properties--property))

<a id="nestedblock--product_allocation--processors--processor--davis--properties--property"></a>
### Nested Schema for `product_allocation.processors.processor.davis.properties.property`

Required:

- `key` (String) no documentation available
- `value` (String) no documentation available




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




<a id="nestedblock--product_allocation--processors--processor--histogram_metric"></a>
### Nested Schema for `product_allocation.processors.processor.histogram_metric`

Required:

- `field` (String) Field with metric value
- `metric_key` (String) Metric key

Optional:

- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--product_allocation--processors--processor--histogram_metric--dimensions))

<a id="nestedblock--product_allocation--processors--processor--histogram_metric--dimensions"></a>
### Nested Schema for `product_allocation.processors.processor.histogram_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--product_allocation--processors--processor--histogram_metric--dimensions--dimension))

<a id="nestedblock--product_allocation--processors--processor--histogram_metric--dimensions--dimension"></a>
### Nested Schema for `product_allocation.processors.processor.histogram_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--product_allocation--processors--processor--product_allocation"></a>
### Nested Schema for `product_allocation.processors.processor.product_allocation`

Required:

- `value` (Block List, Min: 1, Max: 1) The strategy to set product allocation field (see [below for nested schema](#nestedblock--product_allocation--processors--processor--product_allocation--value))

<a id="nestedblock--product_allocation--processors--processor--product_allocation--value"></a>
### Nested Schema for `product_allocation.processors.processor.product_allocation.value`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--product_allocation--processors--processor--product_allocation--value--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--product_allocation--processors--processor--product_allocation--value--field"></a>
### Nested Schema for `product_allocation.processors.processor.product_allocation.value.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value




<a id="nestedblock--product_allocation--processors--processor--sampling_aware_counter_metric"></a>
### Nested Schema for `product_allocation.processors.processor.sampling_aware_counter_metric`

Required:

- `metric_key` (String) Metric key

Optional:

- `aggregation` (String) Possible Values: `Disabled`, `Enabled`
- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--product_allocation--processors--processor--sampling_aware_counter_metric--dimensions))
- `sampling` (String) Possible Values: `Disabled`, `Enabled`

<a id="nestedblock--product_allocation--processors--processor--sampling_aware_counter_metric--dimensions"></a>
### Nested Schema for `product_allocation.processors.processor.sampling_aware_counter_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--product_allocation--processors--processor--sampling_aware_counter_metric--dimensions--dimension))

<a id="nestedblock--product_allocation--processors--processor--sampling_aware_counter_metric--dimensions--dimension"></a>
### Nested Schema for `product_allocation.processors.processor.sampling_aware_counter_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--product_allocation--processors--processor--sampling_aware_value_metric"></a>
### Nested Schema for `product_allocation.processors.processor.sampling_aware_value_metric`

Required:

- `measurement` (String) Possible Values: `Duration`, `Field`
- `metric_key` (String) Metric key

Optional:

- `aggregation` (String) Possible Values: `Disabled`, `Enabled`
- `default_value` (String) Default value with metric value
- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--product_allocation--processors--processor--sampling_aware_value_metric--dimensions))
- `field` (String) Field with metric value
- `sampling` (String) Possible Values: `Disabled`, `Enabled`

<a id="nestedblock--product_allocation--processors--processor--sampling_aware_value_metric--dimensions"></a>
### Nested Schema for `product_allocation.processors.processor.sampling_aware_value_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--product_allocation--processors--processor--sampling_aware_value_metric--dimensions--dimension))

<a id="nestedblock--product_allocation--processors--processor--sampling_aware_value_metric--dimensions--dimension"></a>
### Nested Schema for `product_allocation.processors.processor.sampling_aware_value_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--product_allocation--processors--processor--security_context"></a>
### Nested Schema for `product_allocation.processors.processor.security_context`

Required:

- `value` (Block List, Min: 1, Max: 1) Security context value assignment (see [below for nested schema](#nestedblock--product_allocation--processors--processor--security_context--value))

<a id="nestedblock--product_allocation--processors--processor--security_context--value"></a>
### Nested Schema for `product_allocation.processors.processor.security_context.value`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

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




<a id="nestedblock--product_allocation--processors--processor--security_event"></a>
### Nested Schema for `product_allocation.processors.processor.security_event`

Required:

- `field_extraction` (Block List, Min: 1, Max: 1) Field Extraction (see [below for nested schema](#nestedblock--product_allocation--processors--processor--security_event--field_extraction))

<a id="nestedblock--product_allocation--processors--processor--security_event--field_extraction"></a>
### Nested Schema for `product_allocation.processors.processor.security_event.field_extraction`

Required:

- `type` (String) Possible Values: `Exclude`, `Include`, `IncludeAll`

Optional:

- `exclude` (Set of String) Fields
- `include` (Block List, Max: 1) Fields (see [below for nested schema](#nestedblock--product_allocation--processors--processor--security_event--field_extraction--include))

<a id="nestedblock--product_allocation--processors--processor--security_event--field_extraction--include"></a>
### Nested Schema for `product_allocation.processors.processor.security_event.field_extraction.include`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--product_allocation--processors--processor--security_event--field_extraction--include--dimension))

<a id="nestedblock--product_allocation--processors--processor--security_event--field_extraction--include--dimension"></a>
### Nested Schema for `product_allocation.processors.processor.security_event.field_extraction.include.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name





<a id="nestedblock--product_allocation--processors--processor--technology"></a>
### Nested Schema for `product_allocation.processors.processor.technology`

Required:

- `technology_id` (String) Technology ID

Optional:

- `custom_matcher` (String) Custom matching condition which should be used instead of technology matcher.


<a id="nestedblock--product_allocation--processors--processor--value_metric"></a>
### Nested Schema for `product_allocation.processors.processor.value_metric`

Required:

- `field` (String) Field with metric value
- `metric_key` (String) Metric key

Optional:

- `default_value` (String) Default value with metric value
- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--product_allocation--processors--processor--value_metric--dimensions))

<a id="nestedblock--product_allocation--processors--processor--value_metric--dimensions"></a>
### Nested Schema for `product_allocation.processors.processor.value_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--product_allocation--processors--processor--value_metric--dimensions--dimension))

<a id="nestedblock--product_allocation--processors--processor--value_metric--dimensions--dimension"></a>
### Nested Schema for `product_allocation.processors.processor.value_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name







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

- `description` (String) no documentation available
- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)
- `id` (String) Processor identifier
- `type` (String) Possible Values: `AzureLogForwarding`, `Bizevent`, `BucketAssignment`, `CostAllocation`, `CounterMetric`, `Davis`, `Dql`, `Drop`, `FieldsAdd`, `FieldsRemove`, `FieldsRename`, `HistogramMetric`, `NoStorage`, `ProductAllocation`, `SamplingAwareCounterMetric`, `SamplingAwareValueMetric`, `SecurityContext`, `SecurityEvent`, `Technology`, `ValueMetric`

Optional:

- `azure_log_forwarding` (Block List, Max: 1) Azure log forwarding processor attributes (see [below for nested schema](#nestedblock--security_context--processors--processor--azure_log_forwarding))
- `bizevent` (Block List, Max: 1) Bizevent extraction processor attributes (see [below for nested schema](#nestedblock--security_context--processors--processor--bizevent))
- `bucket_assignment` (Block List, Max: 1) Bucket assignment processor attributes (see [below for nested schema](#nestedblock--security_context--processors--processor--bucket_assignment))
- `cost_allocation` (Block List, Max: 1) Cost allocation processor attributes (see [below for nested schema](#nestedblock--security_context--processors--processor--cost_allocation))
- `counter_metric` (Block List, Max: 1) Counter metric processor attributes (see [below for nested schema](#nestedblock--security_context--processors--processor--counter_metric))
- `davis` (Block List, Max: 1) Davis event extraction processor attributes (see [below for nested schema](#nestedblock--security_context--processors--processor--davis))
- `dql` (Block List, Max: 1) DQL processor attributes (see [below for nested schema](#nestedblock--security_context--processors--processor--dql))
- `fields_add` (Block List, Max: 1) Fields add processor attributes (see [below for nested schema](#nestedblock--security_context--processors--processor--fields_add))
- `fields_remove` (Block List, Max: 1) Fields remove processor attributes (see [below for nested schema](#nestedblock--security_context--processors--processor--fields_remove))
- `fields_rename` (Block List, Max: 1) Fields rename processor attributes (see [below for nested schema](#nestedblock--security_context--processors--processor--fields_rename))
- `histogram_metric` (Block List, Max: 1) Histogram metric processor attributes (see [below for nested schema](#nestedblock--security_context--processors--processor--histogram_metric))
- `matcher` (String) [See our documentation](https://dt-url.net/bp234rv)
- `product_allocation` (Block List, Max: 1) Product allocation processor attributes (see [below for nested schema](#nestedblock--security_context--processors--processor--product_allocation))
- `sample_data` (String) Sample data
- `sampling_aware_counter_metric` (Block List, Max: 1) Sampling-aware counter metric processor attributes (see [below for nested schema](#nestedblock--security_context--processors--processor--sampling_aware_counter_metric))
- `sampling_aware_value_metric` (Block List, Max: 1) Sampling aware value metric processor attributes (see [below for nested schema](#nestedblock--security_context--processors--processor--sampling_aware_value_metric))
- `security_context` (Block List, Max: 1) Security context processor attributes (see [below for nested schema](#nestedblock--security_context--processors--processor--security_context))
- `security_event` (Block List, Max: 1) Security event extraction processor attributes (see [below for nested schema](#nestedblock--security_context--processors--processor--security_event))
- `technology` (Block List, Max: 1) Technology processor attributes (see [below for nested schema](#nestedblock--security_context--processors--processor--technology))
- `value_metric` (Block List, Max: 1) Value metric processor attributes (see [below for nested schema](#nestedblock--security_context--processors--processor--value_metric))

<a id="nestedblock--security_context--processors--processor--azure_log_forwarding"></a>
### Nested Schema for `security_context.processors.processor.azure_log_forwarding`

Required:

- `field_extraction` (Block List, Min: 1, Max: 1) Field Extraction (see [below for nested schema](#nestedblock--security_context--processors--processor--azure_log_forwarding--field_extraction))
- `forwarder_config_id` (String) no documentation available

<a id="nestedblock--security_context--processors--processor--azure_log_forwarding--field_extraction"></a>
### Nested Schema for `security_context.processors.processor.azure_log_forwarding.field_extraction`

Required:

- `type` (String) Possible Values: `Exclude`, `Include`, `IncludeAll`

Optional:

- `exclude` (Set of String) Fields
- `include` (Block List, Max: 1) Fields (see [below for nested schema](#nestedblock--security_context--processors--processor--azure_log_forwarding--field_extraction--include))

<a id="nestedblock--security_context--processors--processor--azure_log_forwarding--field_extraction--include"></a>
### Nested Schema for `security_context.processors.processor.azure_log_forwarding.field_extraction.include`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--security_context--processors--processor--azure_log_forwarding--field_extraction--include--dimension))

<a id="nestedblock--security_context--processors--processor--azure_log_forwarding--field_extraction--include--dimension"></a>
### Nested Schema for `security_context.processors.processor.azure_log_forwarding.field_extraction.include.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name





<a id="nestedblock--security_context--processors--processor--bizevent"></a>
### Nested Schema for `security_context.processors.processor.bizevent`

Required:

- `event_provider` (Block List, Min: 1, Max: 1) Event provider (see [below for nested schema](#nestedblock--security_context--processors--processor--bizevent--event_provider))
- `field_extraction` (Block List, Min: 1, Max: 1) Field extraction (see [below for nested schema](#nestedblock--security_context--processors--processor--bizevent--field_extraction))

Optional:

- `event_type` (Block List, Max: 1) Event type (see [below for nested schema](#nestedblock--security_context--processors--processor--bizevent--event_type))

<a id="nestedblock--security_context--processors--processor--bizevent--event_provider"></a>
### Nested Schema for `security_context.processors.processor.bizevent.event_provider`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--security_context--processors--processor--bizevent--event_provider--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--security_context--processors--processor--bizevent--event_provider--field"></a>
### Nested Schema for `security_context.processors.processor.bizevent.event_provider.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value



<a id="nestedblock--security_context--processors--processor--bizevent--field_extraction"></a>
### Nested Schema for `security_context.processors.processor.bizevent.field_extraction`

Required:

- `type` (String) Possible Values: `Exclude`, `Include`, `IncludeAll`

Optional:

- `exclude` (Set of String) Fields
- `include` (Block List, Max: 1) Fields (see [below for nested schema](#nestedblock--security_context--processors--processor--bizevent--field_extraction--include))

<a id="nestedblock--security_context--processors--processor--bizevent--field_extraction--include"></a>
### Nested Schema for `security_context.processors.processor.bizevent.field_extraction.include`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--security_context--processors--processor--bizevent--field_extraction--include--dimension))

<a id="nestedblock--security_context--processors--processor--bizevent--field_extraction--include--dimension"></a>
### Nested Schema for `security_context.processors.processor.bizevent.field_extraction.include.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--security_context--processors--processor--bizevent--event_type"></a>
### Nested Schema for `security_context.processors.processor.bizevent.event_type`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--security_context--processors--processor--bizevent--event_type--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--security_context--processors--processor--bizevent--event_type--field"></a>
### Nested Schema for `security_context.processors.processor.bizevent.event_type.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value




<a id="nestedblock--security_context--processors--processor--bucket_assignment"></a>
### Nested Schema for `security_context.processors.processor.bucket_assignment`

Required:

- `bucket_name` (String) Bucket name


<a id="nestedblock--security_context--processors--processor--cost_allocation"></a>
### Nested Schema for `security_context.processors.processor.cost_allocation`

Required:

- `value` (Block List, Min: 1, Max: 1) The strategy to set the cost allocation field (see [below for nested schema](#nestedblock--security_context--processors--processor--cost_allocation--value))

<a id="nestedblock--security_context--processors--processor--cost_allocation--value"></a>
### Nested Schema for `security_context.processors.processor.cost_allocation.value`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--security_context--processors--processor--cost_allocation--value--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--security_context--processors--processor--cost_allocation--value--field"></a>
### Nested Schema for `security_context.processors.processor.cost_allocation.value.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value




<a id="nestedblock--security_context--processors--processor--counter_metric"></a>
### Nested Schema for `security_context.processors.processor.counter_metric`

Required:

- `metric_key` (String) Metric key

Optional:

- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--security_context--processors--processor--counter_metric--dimensions))

<a id="nestedblock--security_context--processors--processor--counter_metric--dimensions"></a>
### Nested Schema for `security_context.processors.processor.counter_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--security_context--processors--processor--counter_metric--dimensions--dimension))

<a id="nestedblock--security_context--processors--processor--counter_metric--dimensions--dimension"></a>
### Nested Schema for `security_context.processors.processor.counter_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--security_context--processors--processor--davis"></a>
### Nested Schema for `security_context.processors.processor.davis`

Required:

- `properties` (Block List, Min: 1, Max: 1) no documentation available (see [below for nested schema](#nestedblock--security_context--processors--processor--davis--properties))

<a id="nestedblock--security_context--processors--processor--davis--properties"></a>
### Nested Schema for `security_context.processors.processor.davis.properties`

Required:

- `property` (Block List, Min: 1) (see [below for nested schema](#nestedblock--security_context--processors--processor--davis--properties--property))

<a id="nestedblock--security_context--processors--processor--davis--properties--property"></a>
### Nested Schema for `security_context.processors.processor.davis.properties.property`

Required:

- `key` (String) no documentation available
- `value` (String) no documentation available




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




<a id="nestedblock--security_context--processors--processor--histogram_metric"></a>
### Nested Schema for `security_context.processors.processor.histogram_metric`

Required:

- `field` (String) Field with metric value
- `metric_key` (String) Metric key

Optional:

- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--security_context--processors--processor--histogram_metric--dimensions))

<a id="nestedblock--security_context--processors--processor--histogram_metric--dimensions"></a>
### Nested Schema for `security_context.processors.processor.histogram_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--security_context--processors--processor--histogram_metric--dimensions--dimension))

<a id="nestedblock--security_context--processors--processor--histogram_metric--dimensions--dimension"></a>
### Nested Schema for `security_context.processors.processor.histogram_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--security_context--processors--processor--product_allocation"></a>
### Nested Schema for `security_context.processors.processor.product_allocation`

Required:

- `value` (Block List, Min: 1, Max: 1) The strategy to set product allocation field (see [below for nested schema](#nestedblock--security_context--processors--processor--product_allocation--value))

<a id="nestedblock--security_context--processors--processor--product_allocation--value"></a>
### Nested Schema for `security_context.processors.processor.product_allocation.value`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--security_context--processors--processor--product_allocation--value--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--security_context--processors--processor--product_allocation--value--field"></a>
### Nested Schema for `security_context.processors.processor.product_allocation.value.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value




<a id="nestedblock--security_context--processors--processor--sampling_aware_counter_metric"></a>
### Nested Schema for `security_context.processors.processor.sampling_aware_counter_metric`

Required:

- `metric_key` (String) Metric key

Optional:

- `aggregation` (String) Possible Values: `Disabled`, `Enabled`
- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--security_context--processors--processor--sampling_aware_counter_metric--dimensions))
- `sampling` (String) Possible Values: `Disabled`, `Enabled`

<a id="nestedblock--security_context--processors--processor--sampling_aware_counter_metric--dimensions"></a>
### Nested Schema for `security_context.processors.processor.sampling_aware_counter_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--security_context--processors--processor--sampling_aware_counter_metric--dimensions--dimension))

<a id="nestedblock--security_context--processors--processor--sampling_aware_counter_metric--dimensions--dimension"></a>
### Nested Schema for `security_context.processors.processor.sampling_aware_counter_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--security_context--processors--processor--sampling_aware_value_metric"></a>
### Nested Schema for `security_context.processors.processor.sampling_aware_value_metric`

Required:

- `measurement` (String) Possible Values: `Duration`, `Field`
- `metric_key` (String) Metric key

Optional:

- `aggregation` (String) Possible Values: `Disabled`, `Enabled`
- `default_value` (String) Default value with metric value
- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--security_context--processors--processor--sampling_aware_value_metric--dimensions))
- `field` (String) Field with metric value
- `sampling` (String) Possible Values: `Disabled`, `Enabled`

<a id="nestedblock--security_context--processors--processor--sampling_aware_value_metric--dimensions"></a>
### Nested Schema for `security_context.processors.processor.sampling_aware_value_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--security_context--processors--processor--sampling_aware_value_metric--dimensions--dimension))

<a id="nestedblock--security_context--processors--processor--sampling_aware_value_metric--dimensions--dimension"></a>
### Nested Schema for `security_context.processors.processor.sampling_aware_value_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--security_context--processors--processor--security_context"></a>
### Nested Schema for `security_context.processors.processor.security_context`

Required:

- `value` (Block List, Min: 1, Max: 1) Security context value assignment (see [below for nested schema](#nestedblock--security_context--processors--processor--security_context--value))

<a id="nestedblock--security_context--processors--processor--security_context--value"></a>
### Nested Schema for `security_context.processors.processor.security_context.value`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

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




<a id="nestedblock--security_context--processors--processor--security_event"></a>
### Nested Schema for `security_context.processors.processor.security_event`

Required:

- `field_extraction` (Block List, Min: 1, Max: 1) Field Extraction (see [below for nested schema](#nestedblock--security_context--processors--processor--security_event--field_extraction))

<a id="nestedblock--security_context--processors--processor--security_event--field_extraction"></a>
### Nested Schema for `security_context.processors.processor.security_event.field_extraction`

Required:

- `type` (String) Possible Values: `Exclude`, `Include`, `IncludeAll`

Optional:

- `exclude` (Set of String) Fields
- `include` (Block List, Max: 1) Fields (see [below for nested schema](#nestedblock--security_context--processors--processor--security_event--field_extraction--include))

<a id="nestedblock--security_context--processors--processor--security_event--field_extraction--include"></a>
### Nested Schema for `security_context.processors.processor.security_event.field_extraction.include`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--security_context--processors--processor--security_event--field_extraction--include--dimension))

<a id="nestedblock--security_context--processors--processor--security_event--field_extraction--include--dimension"></a>
### Nested Schema for `security_context.processors.processor.security_event.field_extraction.include.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name





<a id="nestedblock--security_context--processors--processor--technology"></a>
### Nested Schema for `security_context.processors.processor.technology`

Required:

- `technology_id` (String) Technology ID

Optional:

- `custom_matcher` (String) Custom matching condition which should be used instead of technology matcher.


<a id="nestedblock--security_context--processors--processor--value_metric"></a>
### Nested Schema for `security_context.processors.processor.value_metric`

Required:

- `field` (String) Field with metric value
- `metric_key` (String) Metric key

Optional:

- `default_value` (String) Default value with metric value
- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--security_context--processors--processor--value_metric--dimensions))

<a id="nestedblock--security_context--processors--processor--value_metric--dimensions"></a>
### Nested Schema for `security_context.processors.processor.value_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--security_context--processors--processor--value_metric--dimensions--dimension))

<a id="nestedblock--security_context--processors--processor--value_metric--dimensions--dimension"></a>
### Nested Schema for `security_context.processors.processor.value_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name







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

- `description` (String) no documentation available
- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)
- `id` (String) Processor identifier
- `type` (String) Possible Values: `AzureLogForwarding`, `Bizevent`, `BucketAssignment`, `CostAllocation`, `CounterMetric`, `Davis`, `Dql`, `Drop`, `FieldsAdd`, `FieldsRemove`, `FieldsRename`, `HistogramMetric`, `NoStorage`, `ProductAllocation`, `SamplingAwareCounterMetric`, `SamplingAwareValueMetric`, `SecurityContext`, `SecurityEvent`, `Technology`, `ValueMetric`

Optional:

- `azure_log_forwarding` (Block List, Max: 1) Azure log forwarding processor attributes (see [below for nested schema](#nestedblock--storage--processors--processor--azure_log_forwarding))
- `bizevent` (Block List, Max: 1) Bizevent extraction processor attributes (see [below for nested schema](#nestedblock--storage--processors--processor--bizevent))
- `bucket_assignment` (Block List, Max: 1) Bucket assignment processor attributes (see [below for nested schema](#nestedblock--storage--processors--processor--bucket_assignment))
- `cost_allocation` (Block List, Max: 1) Cost allocation processor attributes (see [below for nested schema](#nestedblock--storage--processors--processor--cost_allocation))
- `counter_metric` (Block List, Max: 1) Counter metric processor attributes (see [below for nested schema](#nestedblock--storage--processors--processor--counter_metric))
- `davis` (Block List, Max: 1) Davis event extraction processor attributes (see [below for nested schema](#nestedblock--storage--processors--processor--davis))
- `dql` (Block List, Max: 1) DQL processor attributes (see [below for nested schema](#nestedblock--storage--processors--processor--dql))
- `fields_add` (Block List, Max: 1) Fields add processor attributes (see [below for nested schema](#nestedblock--storage--processors--processor--fields_add))
- `fields_remove` (Block List, Max: 1) Fields remove processor attributes (see [below for nested schema](#nestedblock--storage--processors--processor--fields_remove))
- `fields_rename` (Block List, Max: 1) Fields rename processor attributes (see [below for nested schema](#nestedblock--storage--processors--processor--fields_rename))
- `histogram_metric` (Block List, Max: 1) Histogram metric processor attributes (see [below for nested schema](#nestedblock--storage--processors--processor--histogram_metric))
- `matcher` (String) [See our documentation](https://dt-url.net/bp234rv)
- `product_allocation` (Block List, Max: 1) Product allocation processor attributes (see [below for nested schema](#nestedblock--storage--processors--processor--product_allocation))
- `sample_data` (String) Sample data
- `sampling_aware_counter_metric` (Block List, Max: 1) Sampling-aware counter metric processor attributes (see [below for nested schema](#nestedblock--storage--processors--processor--sampling_aware_counter_metric))
- `sampling_aware_value_metric` (Block List, Max: 1) Sampling aware value metric processor attributes (see [below for nested schema](#nestedblock--storage--processors--processor--sampling_aware_value_metric))
- `security_context` (Block List, Max: 1) Security context processor attributes (see [below for nested schema](#nestedblock--storage--processors--processor--security_context))
- `security_event` (Block List, Max: 1) Security event extraction processor attributes (see [below for nested schema](#nestedblock--storage--processors--processor--security_event))
- `technology` (Block List, Max: 1) Technology processor attributes (see [below for nested schema](#nestedblock--storage--processors--processor--technology))
- `value_metric` (Block List, Max: 1) Value metric processor attributes (see [below for nested schema](#nestedblock--storage--processors--processor--value_metric))

<a id="nestedblock--storage--processors--processor--azure_log_forwarding"></a>
### Nested Schema for `storage.processors.processor.azure_log_forwarding`

Required:

- `field_extraction` (Block List, Min: 1, Max: 1) Field Extraction (see [below for nested schema](#nestedblock--storage--processors--processor--azure_log_forwarding--field_extraction))
- `forwarder_config_id` (String) no documentation available

<a id="nestedblock--storage--processors--processor--azure_log_forwarding--field_extraction"></a>
### Nested Schema for `storage.processors.processor.azure_log_forwarding.field_extraction`

Required:

- `type` (String) Possible Values: `Exclude`, `Include`, `IncludeAll`

Optional:

- `exclude` (Set of String) Fields
- `include` (Block List, Max: 1) Fields (see [below for nested schema](#nestedblock--storage--processors--processor--azure_log_forwarding--field_extraction--include))

<a id="nestedblock--storage--processors--processor--azure_log_forwarding--field_extraction--include"></a>
### Nested Schema for `storage.processors.processor.azure_log_forwarding.field_extraction.include`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--storage--processors--processor--azure_log_forwarding--field_extraction--include--dimension))

<a id="nestedblock--storage--processors--processor--azure_log_forwarding--field_extraction--include--dimension"></a>
### Nested Schema for `storage.processors.processor.azure_log_forwarding.field_extraction.include.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name





<a id="nestedblock--storage--processors--processor--bizevent"></a>
### Nested Schema for `storage.processors.processor.bizevent`

Required:

- `event_provider` (Block List, Min: 1, Max: 1) Event provider (see [below for nested schema](#nestedblock--storage--processors--processor--bizevent--event_provider))
- `field_extraction` (Block List, Min: 1, Max: 1) Field extraction (see [below for nested schema](#nestedblock--storage--processors--processor--bizevent--field_extraction))

Optional:

- `event_type` (Block List, Max: 1) Event type (see [below for nested schema](#nestedblock--storage--processors--processor--bizevent--event_type))

<a id="nestedblock--storage--processors--processor--bizevent--event_provider"></a>
### Nested Schema for `storage.processors.processor.bizevent.event_provider`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--storage--processors--processor--bizevent--event_provider--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--storage--processors--processor--bizevent--event_provider--field"></a>
### Nested Schema for `storage.processors.processor.bizevent.event_provider.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value



<a id="nestedblock--storage--processors--processor--bizevent--field_extraction"></a>
### Nested Schema for `storage.processors.processor.bizevent.field_extraction`

Required:

- `type` (String) Possible Values: `Exclude`, `Include`, `IncludeAll`

Optional:

- `exclude` (Set of String) Fields
- `include` (Block List, Max: 1) Fields (see [below for nested schema](#nestedblock--storage--processors--processor--bizevent--field_extraction--include))

<a id="nestedblock--storage--processors--processor--bizevent--field_extraction--include"></a>
### Nested Schema for `storage.processors.processor.bizevent.field_extraction.include`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--storage--processors--processor--bizevent--field_extraction--include--dimension))

<a id="nestedblock--storage--processors--processor--bizevent--field_extraction--include--dimension"></a>
### Nested Schema for `storage.processors.processor.bizevent.field_extraction.include.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--storage--processors--processor--bizevent--event_type"></a>
### Nested Schema for `storage.processors.processor.bizevent.event_type`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--storage--processors--processor--bizevent--event_type--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--storage--processors--processor--bizevent--event_type--field"></a>
### Nested Schema for `storage.processors.processor.bizevent.event_type.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value




<a id="nestedblock--storage--processors--processor--bucket_assignment"></a>
### Nested Schema for `storage.processors.processor.bucket_assignment`

Required:

- `bucket_name` (String) Bucket name


<a id="nestedblock--storage--processors--processor--cost_allocation"></a>
### Nested Schema for `storage.processors.processor.cost_allocation`

Required:

- `value` (Block List, Min: 1, Max: 1) The strategy to set the cost allocation field (see [below for nested schema](#nestedblock--storage--processors--processor--cost_allocation--value))

<a id="nestedblock--storage--processors--processor--cost_allocation--value"></a>
### Nested Schema for `storage.processors.processor.cost_allocation.value`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--storage--processors--processor--cost_allocation--value--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--storage--processors--processor--cost_allocation--value--field"></a>
### Nested Schema for `storage.processors.processor.cost_allocation.value.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value




<a id="nestedblock--storage--processors--processor--counter_metric"></a>
### Nested Schema for `storage.processors.processor.counter_metric`

Required:

- `metric_key` (String) Metric key

Optional:

- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--storage--processors--processor--counter_metric--dimensions))

<a id="nestedblock--storage--processors--processor--counter_metric--dimensions"></a>
### Nested Schema for `storage.processors.processor.counter_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--storage--processors--processor--counter_metric--dimensions--dimension))

<a id="nestedblock--storage--processors--processor--counter_metric--dimensions--dimension"></a>
### Nested Schema for `storage.processors.processor.counter_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--storage--processors--processor--davis"></a>
### Nested Schema for `storage.processors.processor.davis`

Required:

- `properties` (Block List, Min: 1, Max: 1) no documentation available (see [below for nested schema](#nestedblock--storage--processors--processor--davis--properties))

<a id="nestedblock--storage--processors--processor--davis--properties"></a>
### Nested Schema for `storage.processors.processor.davis.properties`

Required:

- `property` (Block List, Min: 1) (see [below for nested schema](#nestedblock--storage--processors--processor--davis--properties--property))

<a id="nestedblock--storage--processors--processor--davis--properties--property"></a>
### Nested Schema for `storage.processors.processor.davis.properties.property`

Required:

- `key` (String) no documentation available
- `value` (String) no documentation available




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




<a id="nestedblock--storage--processors--processor--histogram_metric"></a>
### Nested Schema for `storage.processors.processor.histogram_metric`

Required:

- `field` (String) Field with metric value
- `metric_key` (String) Metric key

Optional:

- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--storage--processors--processor--histogram_metric--dimensions))

<a id="nestedblock--storage--processors--processor--histogram_metric--dimensions"></a>
### Nested Schema for `storage.processors.processor.histogram_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--storage--processors--processor--histogram_metric--dimensions--dimension))

<a id="nestedblock--storage--processors--processor--histogram_metric--dimensions--dimension"></a>
### Nested Schema for `storage.processors.processor.histogram_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--storage--processors--processor--product_allocation"></a>
### Nested Schema for `storage.processors.processor.product_allocation`

Required:

- `value` (Block List, Min: 1, Max: 1) The strategy to set product allocation field (see [below for nested schema](#nestedblock--storage--processors--processor--product_allocation--value))

<a id="nestedblock--storage--processors--processor--product_allocation--value"></a>
### Nested Schema for `storage.processors.processor.product_allocation.value`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--storage--processors--processor--product_allocation--value--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--storage--processors--processor--product_allocation--value--field"></a>
### Nested Schema for `storage.processors.processor.product_allocation.value.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value




<a id="nestedblock--storage--processors--processor--sampling_aware_counter_metric"></a>
### Nested Schema for `storage.processors.processor.sampling_aware_counter_metric`

Required:

- `metric_key` (String) Metric key

Optional:

- `aggregation` (String) Possible Values: `Disabled`, `Enabled`
- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--storage--processors--processor--sampling_aware_counter_metric--dimensions))
- `sampling` (String) Possible Values: `Disabled`, `Enabled`

<a id="nestedblock--storage--processors--processor--sampling_aware_counter_metric--dimensions"></a>
### Nested Schema for `storage.processors.processor.sampling_aware_counter_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--storage--processors--processor--sampling_aware_counter_metric--dimensions--dimension))

<a id="nestedblock--storage--processors--processor--sampling_aware_counter_metric--dimensions--dimension"></a>
### Nested Schema for `storage.processors.processor.sampling_aware_counter_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--storage--processors--processor--sampling_aware_value_metric"></a>
### Nested Schema for `storage.processors.processor.sampling_aware_value_metric`

Required:

- `measurement` (String) Possible Values: `Duration`, `Field`
- `metric_key` (String) Metric key

Optional:

- `aggregation` (String) Possible Values: `Disabled`, `Enabled`
- `default_value` (String) Default value with metric value
- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--storage--processors--processor--sampling_aware_value_metric--dimensions))
- `field` (String) Field with metric value
- `sampling` (String) Possible Values: `Disabled`, `Enabled`

<a id="nestedblock--storage--processors--processor--sampling_aware_value_metric--dimensions"></a>
### Nested Schema for `storage.processors.processor.sampling_aware_value_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--storage--processors--processor--sampling_aware_value_metric--dimensions--dimension))

<a id="nestedblock--storage--processors--processor--sampling_aware_value_metric--dimensions--dimension"></a>
### Nested Schema for `storage.processors.processor.sampling_aware_value_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name




<a id="nestedblock--storage--processors--processor--security_context"></a>
### Nested Schema for `storage.processors.processor.security_context`

Required:

- `value` (Block List, Min: 1, Max: 1) Security context value assignment (see [below for nested schema](#nestedblock--storage--processors--processor--security_context--value))

<a id="nestedblock--storage--processors--processor--security_context--value"></a>
### Nested Schema for `storage.processors.processor.security_context.value`

Required:

- `type` (String) Possible Values: `Constant`, `Field`, `MultiValueConstant`

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




<a id="nestedblock--storage--processors--processor--security_event"></a>
### Nested Schema for `storage.processors.processor.security_event`

Required:

- `field_extraction` (Block List, Min: 1, Max: 1) Field Extraction (see [below for nested schema](#nestedblock--storage--processors--processor--security_event--field_extraction))

<a id="nestedblock--storage--processors--processor--security_event--field_extraction"></a>
### Nested Schema for `storage.processors.processor.security_event.field_extraction`

Required:

- `type` (String) Possible Values: `Exclude`, `Include`, `IncludeAll`

Optional:

- `exclude` (Set of String) Fields
- `include` (Block List, Max: 1) Fields (see [below for nested schema](#nestedblock--storage--processors--processor--security_event--field_extraction--include))

<a id="nestedblock--storage--processors--processor--security_event--field_extraction--include"></a>
### Nested Schema for `storage.processors.processor.security_event.field_extraction.include`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--storage--processors--processor--security_event--field_extraction--include--dimension))

<a id="nestedblock--storage--processors--processor--security_event--field_extraction--include--dimension"></a>
### Nested Schema for `storage.processors.processor.security_event.field_extraction.include.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name





<a id="nestedblock--storage--processors--processor--technology"></a>
### Nested Schema for `storage.processors.processor.technology`

Required:

- `technology_id` (String) Technology ID

Optional:

- `custom_matcher` (String) Custom matching condition which should be used instead of technology matcher.


<a id="nestedblock--storage--processors--processor--value_metric"></a>
### Nested Schema for `storage.processors.processor.value_metric`

Required:

- `field` (String) Field with metric value
- `metric_key` (String) Metric key

Optional:

- `default_value` (String) Default value with metric value
- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--storage--processors--processor--value_metric--dimensions))

<a id="nestedblock--storage--processors--processor--value_metric--dimensions"></a>
### Nested Schema for `storage.processors.processor.value_metric.dimensions`

Required:

- `dimension` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--storage--processors--processor--value_metric--dimensions--dimension))

<a id="nestedblock--storage--processors--processor--value_metric--dimensions--dimension"></a>
### Nested Schema for `storage.processors.processor.value_metric.dimensions.dimension`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name
