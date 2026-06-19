---
layout: ""
page_title: "dynatrace_openpipeline_v2_events_security_dataforwarding Resource - terraform-provider-dynatrace"
subcategory: "OpenPipeline V2"
description: |-
  The resource `dynatrace_openpipeline_v2_events_security_dataforwarding` covers configuration of OpenPipeline for events security data forwarding
---

# dynatrace_openpipeline_v2_events_security_dataforwarding (Resource)

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
- Forwarding - https://docs.dynatrace.com/docs/platform/openpipeline/concepts/forwarding

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_openpipeline_v2_events_security_dataforwarding` downloads all existing OpenPipeline definitions for events security data forwarding

The full documentation of the export feature is available [here](https://dt-url.net/h203qmc).

## Resource Example Usage

```terraform
resource "dynatrace_openpipeline_v2_events_security_dataforwarding" "example" {
  forwarding_name   = "#name#"
  enabled           = false
  matcher           = "true"
  cloud_vendor_type = "gcp"
  gcp_connection {
    bucket_name   = "my-bucket"
    connection_id = dynatrace_gcp_connection.connection.id
  }
  data_forwarding_type = "processed"
  pipelines            = [dynatrace_openpipeline_v2_events_security_pipelines.pipeline.id]
  bulk_pattern         = "<YYYYMMDD>/<HH>/<HHmmss.SSSS>_<bulk-id>.json.gz"
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
        enabled = true
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
        enabled = true
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
        enabled = true
      }
    }
  }
}

resource "dynatrace_openpipeline_v2_events_security_pipelines" "pipeline" {
  display_name = "Minimal pipeline"
  custom_id    = "pipeline_Minimal_pipeline_1234_tf_#name#"
}

# Create GCP connection
variable "DT_GCP_TEST_IMPERSONABLE_SERVICE_ACCOUNT" {
  type        = string
  description = "The service account that should be used for the GCP connection setup end-to-end test."
  sensitive   = true
}

resource "dynatrace_gcp_connection" "connection" {
  name = "#name#"
  type = "serviceAccountImpersonation"
  service_account_impersonation {
    service_account_id = var.DT_GCP_TEST_IMPERSONABLE_SERVICE_ACCOUNT
    consumers = [
      "SVC:com.dynatrace.openpipeline"
    ]
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `bulk_pattern` (String) Segmentation and prefix of the data
- `cloud_vendor_type` (String) Cloud Vendor Type. Possible values: `aws`, `azure`, `gcp`
- `data_forwarding_type` (String) Pipeline Type. Possible values: `processed`, `raw`
- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)
- `forwarding_name` (String) Forwarding name
- `matcher` (String) Query which determines whether the record should be routed to the target pipeline of this rule.

### Optional

- `aws_connection` (Block List, Max: 1) AWS Connection (see [below for nested schema](#nestedblock--aws_connection))
- `azure_connection` (Block List, Max: 1) Azure Connection (see [below for nested schema](#nestedblock--azure_connection))
- `builtin_ingest_sources` (Set of String) List of built-in ingest sources
- `builtin_pipelines` (Set of String) Built-in pipelines
- `bulk_size` (Number) Bulk size for transmission
- `gcp_connection` (Block List, Max: 1) GCP Connection (see [below for nested schema](#nestedblock--gcp_connection))
- `ingest_sources` (Set of String) List of ingest sources
- `pipelines` (Set of String) Pipelines
- `processing` (Block List, Max: 1) Processing (see [below for nested schema](#nestedblock--processing))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--aws_connection"></a>
### Nested Schema for `aws_connection`

Required:

- `arn` (String) S3 Bucket ARN
- `connection_id` (String) AWS connection


<a id="nestedblock--azure_connection"></a>
### Nested Schema for `azure_connection`

Required:

- `connection_id` (String) Azure connection
- `container_url` (String) Container Connection URL


<a id="nestedblock--gcp_connection"></a>
### Nested Schema for `gcp_connection`

Required:

- `bucket_name` (String) GCS Bucket Name
- `connection_id` (String) GCP connection


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
- `type` (String) Processor type. Possible values: `azureLogForwarding`, `bizevent`, `bucketAssignment`, `costAllocation`, `counterMetric`, `davis`, `dql`, `drop`, `fieldsAdd`, `fieldsRemove`, `fieldsRename`, `geoLookup`, `histogramMetric`, `noStorage`, `productAllocation`, `samplingAwareCounterMetric`, `samplingAwareHistogramMetric`, `samplingAwareValueMetric`, `sdlcEvent`, `securityContext`, `securityEvent`, `smartscapeEdge`, `smartscapeNode`, `technology`, `valueMetric`

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
- `geo_lookup` (Block List, Max: 1) Geo lookup processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--geo_lookup))
- `histogram_metric` (Block List, Max: 1) Histogram metric processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--histogram_metric))
- `matcher` (String) [See our documentation](https://dt-url.net/bp234rv)
- `product_allocation` (Block List, Max: 1) Product allocation processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--product_allocation))
- `sample_data` (String) Sample data
- `sampling_aware_counter_metric` (Block List, Max: 1) Sampling-aware counter metric processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--sampling_aware_counter_metric))
- `sampling_aware_histogram_metric` (Block List, Max: 1) Sampling aware histogram metric processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--sampling_aware_histogram_metric))
- `sampling_aware_value_metric` (Block List, Max: 1) Sampling aware value metric processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--sampling_aware_value_metric))
- `sdlc_event` (Block List, Max: 1) SdlcEvent extraction processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--sdlc_event))
- `security_context` (Block List, Max: 1) Security context processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--security_context))
- `security_event` (Block List, Max: 1) Security event extraction processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--security_event))
- `smartscape_edge` (Block List, Max: 1) Smartscape edge extraction processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--smartscape_edge))
- `smartscape_node` (Block List, Max: 1) Smartscape node extraction processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--smartscape_node))
- `technology` (Block List, Max: 1) Technology processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--technology))
- `value_metric` (Block List, Max: 1) Value metric processor attributes (see [below for nested schema](#nestedblock--processing--processors--processor--value_metric))

<a id="nestedblock--processing--processors--processor--azure_log_forwarding"></a>
### Nested Schema for `processing.processors.processor.azure_log_forwarding`

Required:

- `field_extraction` (Block List, Min: 1, Max: 1) Field Extraction (see [below for nested schema](#nestedblock--processing--processors--processor--azure_log_forwarding--field_extraction))
- `forwarder_config_id` (String) No documentation available

<a id="nestedblock--processing--processors--processor--azure_log_forwarding--field_extraction"></a>
### Nested Schema for `processing.processors.processor.azure_log_forwarding.field_extraction`

Required:

- `type` (String) Fields Extraction type. Possible values: `exclude`, `include`, `includeAll`

Optional:

- `exclude` (Set of String) Fields
- `include` (Block List, Max: 1) Fields (see [below for nested schema](#nestedblock--processing--processors--processor--azure_log_forwarding--field_extraction--include))

<a id="nestedblock--processing--processors--processor--azure_log_forwarding--field_extraction--include"></a>
### Nested Schema for `processing.processors.processor.azure_log_forwarding.field_extraction.include`

Required:

- `field_extraction_entry` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--processing--processors--processor--azure_log_forwarding--field_extraction--include--field_extraction_entry))

<a id="nestedblock--processing--processors--processor--azure_log_forwarding--field_extraction--include--field_extraction_entry"></a>
### Nested Schema for `processing.processors.processor.azure_log_forwarding.field_extraction.include.field_extraction_entry`

Optional:

- `constant_field_name` (String) Destination field name
- `constant_value` (String) Constant value to be assigned to field
- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name
- `extraction_type` (String) Field value extraction type. Possible values: `constant`, `field`
- `source_field_name` (String) Source field name
- `strategy` (String) Strategy for field extraction. Possible values: `equals`, `startsWith`





<a id="nestedblock--processing--processors--processor--bizevent"></a>
### Nested Schema for `processing.processors.processor.bizevent`

Required:

- `event_provider` (Block List, Min: 1, Max: 1) Event provider (see [below for nested schema](#nestedblock--processing--processors--processor--bizevent--event_provider))
- `event_type` (Block List, Min: 1, Max: 1) Event type (see [below for nested schema](#nestedblock--processing--processors--processor--bizevent--event_type))
- `field_extraction` (Block List, Min: 1, Max: 1) Field extraction (see [below for nested schema](#nestedblock--processing--processors--processor--bizevent--field_extraction))

<a id="nestedblock--processing--processors--processor--bizevent--event_provider"></a>
### Nested Schema for `processing.processors.processor.bizevent.event_provider`

Required:

- `type` (String) Type of value assignment. Possible values: `constant`, `field`, `multiValueConstant`

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



<a id="nestedblock--processing--processors--processor--bizevent--event_type"></a>
### Nested Schema for `processing.processors.processor.bizevent.event_type`

Required:

- `type` (String) Type of value assignment. Possible values: `constant`, `field`, `multiValueConstant`

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



<a id="nestedblock--processing--processors--processor--bizevent--field_extraction"></a>
### Nested Schema for `processing.processors.processor.bizevent.field_extraction`

Required:

- `type` (String) Fields Extraction type. Possible values: `exclude`, `include`, `includeAll`

Optional:

- `exclude` (Set of String) Fields
- `include` (Block List, Max: 1) Fields (see [below for nested schema](#nestedblock--processing--processors--processor--bizevent--field_extraction--include))

<a id="nestedblock--processing--processors--processor--bizevent--field_extraction--include"></a>
### Nested Schema for `processing.processors.processor.bizevent.field_extraction.include`

Required:

- `field_extraction_entry` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--processing--processors--processor--bizevent--field_extraction--include--field_extraction_entry))

<a id="nestedblock--processing--processors--processor--bizevent--field_extraction--include--field_extraction_entry"></a>
### Nested Schema for `processing.processors.processor.bizevent.field_extraction.include.field_extraction_entry`

Optional:

- `constant_field_name` (String) Destination field name
- `constant_value` (String) Constant value to be assigned to field
- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name
- `extraction_type` (String) Field value extraction type. Possible values: `constant`, `field`
- `source_field_name` (String) Source field name
- `strategy` (String) Strategy for field extraction. Possible values: `equals`, `startsWith`





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

- `type` (String) Type of value assignment. Possible values: `constant`, `field`, `multiValueConstant`

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

- `field_extraction_entry` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--processing--processors--processor--counter_metric--dimensions--field_extraction_entry))

<a id="nestedblock--processing--processors--processor--counter_metric--dimensions--field_extraction_entry"></a>
### Nested Schema for `processing.processors.processor.counter_metric.dimensions.field_extraction_entry`

Optional:

- `constant_field_name` (String) Destination field name
- `constant_value` (String) Constant value to be assigned to field
- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name
- `extraction_type` (String) Field value extraction type. Possible values: `constant`, `field`
- `source_field_name` (String) Source field name
- `strategy` (String) Strategy for field extraction. Possible values: `equals`, `startsWith`




<a id="nestedblock--processing--processors--processor--davis"></a>
### Nested Schema for `processing.processors.processor.davis`

Required:

- `properties` (Block List, Min: 1, Max: 1) No documentation available (see [below for nested schema](#nestedblock--processing--processors--processor--davis--properties))

<a id="nestedblock--processing--processors--processor--davis--properties"></a>
### Nested Schema for `processing.processors.processor.davis.properties`

Required:

- `property` (Block List, Min: 1) (see [below for nested schema](#nestedblock--processing--processors--processor--davis--properties--property))

<a id="nestedblock--processing--processors--processor--davis--properties--property"></a>
### Nested Schema for `processing.processors.processor.davis.properties.property`

Required:

- `key` (String) No documentation available

Optional:

- `strategy` (String) Strategy for field extraction. Possible values: `equals`, `startsWith`
- `value` (String) No documentation available




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




<a id="nestedblock--processing--processors--processor--geo_lookup"></a>
### Nested Schema for `processing.processors.processor.geo_lookup`

Required:

- `ip_field_key` (String) The field key that contains the IP address to be resolved to a geo location.

Optional:

- `geo_field_prefix` (String) Optional prefix for all output geo fields. If specified, output fields will be prefixed as <prefix>.geo.<field>. If omitted, output fields will be geo.<field>.
- `output_fields` (Set of String) The geo fields to enrich the record with. If empty or not specified, the default fields (city name, country ISO code, country name, location) are used. Possible values: `cityName`, `continentIsoCode`, `continentName`, `countryIsoCode`, `countryName`, `location`, `postalCode`, `regionIsoCode`, `regionName`, `subdivisionIsoCodes`


<a id="nestedblock--processing--processors--processor--histogram_metric"></a>
### Nested Schema for `processing.processors.processor.histogram_metric`

Required:

- `field` (String) Field with metric value
- `metric_key` (String) Metric key

Optional:

- `default_value` (String) Default value with metric value
- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--processing--processors--processor--histogram_metric--dimensions))

<a id="nestedblock--processing--processors--processor--histogram_metric--dimensions"></a>
### Nested Schema for `processing.processors.processor.histogram_metric.dimensions`

Required:

- `field_extraction_entry` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--processing--processors--processor--histogram_metric--dimensions--field_extraction_entry))

<a id="nestedblock--processing--processors--processor--histogram_metric--dimensions--field_extraction_entry"></a>
### Nested Schema for `processing.processors.processor.histogram_metric.dimensions.field_extraction_entry`

Optional:

- `constant_field_name` (String) Destination field name
- `constant_value` (String) Constant value to be assigned to field
- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name
- `extraction_type` (String) Field value extraction type. Possible values: `constant`, `field`
- `source_field_name` (String) Source field name
- `strategy` (String) Strategy for field extraction. Possible values: `equals`, `startsWith`




<a id="nestedblock--processing--processors--processor--product_allocation"></a>
### Nested Schema for `processing.processors.processor.product_allocation`

Required:

- `value` (Block List, Min: 1, Max: 1) The strategy to set product allocation field (see [below for nested schema](#nestedblock--processing--processors--processor--product_allocation--value))

<a id="nestedblock--processing--processors--processor--product_allocation--value"></a>
### Nested Schema for `processing.processors.processor.product_allocation.value`

Required:

- `type` (String) Type of value assignment. Possible values: `constant`, `field`, `multiValueConstant`

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

- `aggregation` (String) Possible values: `disabled`, `enabled`
- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--processing--processors--processor--sampling_aware_counter_metric--dimensions))
- `sampling` (String) Possible values: `disabled`, `enabled`

<a id="nestedblock--processing--processors--processor--sampling_aware_counter_metric--dimensions"></a>
### Nested Schema for `processing.processors.processor.sampling_aware_counter_metric.dimensions`

Required:

- `field_extraction_entry` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--processing--processors--processor--sampling_aware_counter_metric--dimensions--field_extraction_entry))

<a id="nestedblock--processing--processors--processor--sampling_aware_counter_metric--dimensions--field_extraction_entry"></a>
### Nested Schema for `processing.processors.processor.sampling_aware_counter_metric.dimensions.field_extraction_entry`

Optional:

- `constant_field_name` (String) Destination field name
- `constant_value` (String) Constant value to be assigned to field
- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name
- `extraction_type` (String) Field value extraction type. Possible values: `constant`, `field`
- `source_field_name` (String) Source field name
- `strategy` (String) Strategy for field extraction. Possible values: `equals`, `startsWith`




<a id="nestedblock--processing--processors--processor--sampling_aware_histogram_metric"></a>
### Nested Schema for `processing.processors.processor.sampling_aware_histogram_metric`

Required:

- `measurement` (String) Possible values: `duration`, `field`
- `metric_key` (String) Metric key

Optional:

- `aggregation` (String) Possible values: `disabled`, `enabled`
- `default_value` (String) Default value with metric value
- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--processing--processors--processor--sampling_aware_histogram_metric--dimensions))
- `field` (String) Field with metric value
- `sampling` (String) Possible values: `disabled`, `enabled`

<a id="nestedblock--processing--processors--processor--sampling_aware_histogram_metric--dimensions"></a>
### Nested Schema for `processing.processors.processor.sampling_aware_histogram_metric.dimensions`

Required:

- `field_extraction_entry` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--processing--processors--processor--sampling_aware_histogram_metric--dimensions--field_extraction_entry))

<a id="nestedblock--processing--processors--processor--sampling_aware_histogram_metric--dimensions--field_extraction_entry"></a>
### Nested Schema for `processing.processors.processor.sampling_aware_histogram_metric.dimensions.field_extraction_entry`

Optional:

- `constant_field_name` (String) Destination field name
- `constant_value` (String) Constant value to be assigned to field
- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name
- `extraction_type` (String) Field value extraction type. Possible values: `constant`, `field`
- `source_field_name` (String) Source field name
- `strategy` (String) Strategy for field extraction. Possible values: `equals`, `startsWith`




<a id="nestedblock--processing--processors--processor--sampling_aware_value_metric"></a>
### Nested Schema for `processing.processors.processor.sampling_aware_value_metric`

Required:

- `measurement` (String) Possible values: `duration`, `field`
- `metric_key` (String) Metric key

Optional:

- `aggregation` (String) Possible values: `disabled`, `enabled`
- `default_value` (String) Default value with metric value
- `dimensions` (Block List, Max: 1) List of dimensions (see [below for nested schema](#nestedblock--processing--processors--processor--sampling_aware_value_metric--dimensions))
- `field` (String) Field with metric value
- `sampling` (String) Possible values: `disabled`, `enabled`

<a id="nestedblock--processing--processors--processor--sampling_aware_value_metric--dimensions"></a>
### Nested Schema for `processing.processors.processor.sampling_aware_value_metric.dimensions`

Required:

- `field_extraction_entry` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--processing--processors--processor--sampling_aware_value_metric--dimensions--field_extraction_entry))

<a id="nestedblock--processing--processors--processor--sampling_aware_value_metric--dimensions--field_extraction_entry"></a>
### Nested Schema for `processing.processors.processor.sampling_aware_value_metric.dimensions.field_extraction_entry`

Optional:

- `constant_field_name` (String) Destination field name
- `constant_value` (String) Constant value to be assigned to field
- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name
- `extraction_type` (String) Field value extraction type. Possible values: `constant`, `field`
- `source_field_name` (String) Source field name
- `strategy` (String) Strategy for field extraction. Possible values: `equals`, `startsWith`




<a id="nestedblock--processing--processors--processor--sdlc_event"></a>
### Nested Schema for `processing.processors.processor.sdlc_event`

Required:

- `event_category` (Block List, Min: 1, Max: 1) Event category (see [below for nested schema](#nestedblock--processing--processors--processor--sdlc_event--event_category))
- `event_provider` (Block List, Min: 1, Max: 1) Event provider (see [below for nested schema](#nestedblock--processing--processors--processor--sdlc_event--event_provider))
- `event_status` (Block List, Min: 1, Max: 1) Event status (see [below for nested schema](#nestedblock--processing--processors--processor--sdlc_event--event_status))
- `event_type` (Block List, Min: 1, Max: 1) Event type (see [below for nested schema](#nestedblock--processing--processors--processor--sdlc_event--event_type))
- `field_extraction` (Block List, Min: 1, Max: 1) Field extraction (see [below for nested schema](#nestedblock--processing--processors--processor--sdlc_event--field_extraction))

<a id="nestedblock--processing--processors--processor--sdlc_event--event_category"></a>
### Nested Schema for `processing.processors.processor.sdlc_event.event_category`

Required:

- `type` (String) Type of value assignment. Possible values: `constant`, `field`, `multiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--processing--processors--processor--sdlc_event--event_category--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--processing--processors--processor--sdlc_event--event_category--field"></a>
### Nested Schema for `processing.processors.processor.sdlc_event.event_category.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value



<a id="nestedblock--processing--processors--processor--sdlc_event--event_provider"></a>
### Nested Schema for `processing.processors.processor.sdlc_event.event_provider`

Required:

- `type` (String) Type of value assignment. Possible values: `constant`, `field`, `multiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--processing--processors--processor--sdlc_event--event_provider--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--processing--processors--processor--sdlc_event--event_provider--field"></a>
### Nested Schema for `processing.processors.processor.sdlc_event.event_provider.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value



<a id="nestedblock--processing--processors--processor--sdlc_event--event_status"></a>
### Nested Schema for `processing.processors.processor.sdlc_event.event_status`

Required:

- `type` (String) Type of value assignment. Possible values: `constant`, `field`, `multiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--processing--processors--processor--sdlc_event--event_status--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--processing--processors--processor--sdlc_event--event_status--field"></a>
### Nested Schema for `processing.processors.processor.sdlc_event.event_status.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value



<a id="nestedblock--processing--processors--processor--sdlc_event--event_type"></a>
### Nested Schema for `processing.processors.processor.sdlc_event.event_type`

Required:

- `type` (String) Type of value assignment. Possible values: `constant`, `field`, `multiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--processing--processors--processor--sdlc_event--event_type--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--processing--processors--processor--sdlc_event--event_type--field"></a>
### Nested Schema for `processing.processors.processor.sdlc_event.event_type.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value



<a id="nestedblock--processing--processors--processor--sdlc_event--field_extraction"></a>
### Nested Schema for `processing.processors.processor.sdlc_event.field_extraction`

Required:

- `type` (String) Fields Extraction type. Possible values: `exclude`, `include`, `includeAll`

Optional:

- `exclude` (Set of String) Fields
- `include` (Block List, Max: 1) Fields (see [below for nested schema](#nestedblock--processing--processors--processor--sdlc_event--field_extraction--include))

<a id="nestedblock--processing--processors--processor--sdlc_event--field_extraction--include"></a>
### Nested Schema for `processing.processors.processor.sdlc_event.field_extraction.include`

Required:

- `field_extraction_entry` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--processing--processors--processor--sdlc_event--field_extraction--include--field_extraction_entry))

<a id="nestedblock--processing--processors--processor--sdlc_event--field_extraction--include--field_extraction_entry"></a>
### Nested Schema for `processing.processors.processor.sdlc_event.field_extraction.include.field_extraction_entry`

Optional:

- `constant_field_name` (String) Destination field name
- `constant_value` (String) Constant value to be assigned to field
- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name
- `extraction_type` (String) Field value extraction type. Possible values: `constant`, `field`
- `source_field_name` (String) Source field name
- `strategy` (String) Strategy for field extraction. Possible values: `equals`, `startsWith`





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




<a id="nestedblock--processing--processors--processor--security_event"></a>
### Nested Schema for `processing.processors.processor.security_event`

Required:

- `field_extraction` (Block List, Min: 1, Max: 1) Field Extraction (see [below for nested schema](#nestedblock--processing--processors--processor--security_event--field_extraction))

<a id="nestedblock--processing--processors--processor--security_event--field_extraction"></a>
### Nested Schema for `processing.processors.processor.security_event.field_extraction`

Required:

- `type` (String) Fields Extraction type. Possible values: `exclude`, `include`, `includeAll`

Optional:

- `exclude` (Set of String) Fields
- `include` (Block List, Max: 1) Fields (see [below for nested schema](#nestedblock--processing--processors--processor--security_event--field_extraction--include))

<a id="nestedblock--processing--processors--processor--security_event--field_extraction--include"></a>
### Nested Schema for `processing.processors.processor.security_event.field_extraction.include`

Required:

- `field_extraction_entry` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--processing--processors--processor--security_event--field_extraction--include--field_extraction_entry))

<a id="nestedblock--processing--processors--processor--security_event--field_extraction--include--field_extraction_entry"></a>
### Nested Schema for `processing.processors.processor.security_event.field_extraction.include.field_extraction_entry`

Optional:

- `constant_field_name` (String) Destination field name
- `constant_value` (String) Constant value to be assigned to field
- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name
- `extraction_type` (String) Field value extraction type. Possible values: `constant`, `field`
- `source_field_name` (String) Source field name
- `strategy` (String) Strategy for field extraction. Possible values: `equals`, `startsWith`





<a id="nestedblock--processing--processors--processor--smartscape_edge"></a>
### Nested Schema for `processing.processors.processor.smartscape_edge`

Required:

- `edge_type` (String) Edge type
- `source_id_field_name` (String) Source ID field name
- `source_type` (String) Source type
- `target_id_field_name` (String) Target ID field name
- `target_type` (String) Target type


<a id="nestedblock--processing--processors--processor--smartscape_node"></a>
### Nested Schema for `processing.processors.processor.smartscape_node`

Required:

- `extract_node` (Boolean) Extract node
- `id_components` (Block List, Min: 1, Max: 1) ID components (see [below for nested schema](#nestedblock--processing--processors--processor--smartscape_node--id_components))
- `node_id_field_name` (String) Node ID field name
- `node_type` (String) Node type

Optional:

- `fields_to_extract` (Block List, Max: 1) Fields to extract (see [below for nested schema](#nestedblock--processing--processors--processor--smartscape_node--fields_to_extract))
- `node_name` (Block List, Max: 1) Node name (see [below for nested schema](#nestedblock--processing--processors--processor--smartscape_node--node_name))
- `static_edges_to_extract` (Block List, Max: 1) Static edges to extract (see [below for nested schema](#nestedblock--processing--processors--processor--smartscape_node--static_edges_to_extract))

<a id="nestedblock--processing--processors--processor--smartscape_node--id_components"></a>
### Nested Schema for `processing.processors.processor.smartscape_node.id_components`

Required:

- `id_component` (Block List, Min: 1) (see [below for nested schema](#nestedblock--processing--processors--processor--smartscape_node--id_components--id_component))

<a id="nestedblock--processing--processors--processor--smartscape_node--id_components--id_component"></a>
### Nested Schema for `processing.processors.processor.smartscape_node.id_components.id_component`

Required:

- `id_component` (String) ID component
- `referenced_field_name` (String) Referenced field name



<a id="nestedblock--processing--processors--processor--smartscape_node--fields_to_extract"></a>
### Nested Schema for `processing.processors.processor.smartscape_node.fields_to_extract`

Required:

- `smartscape_field_extraction_entry` (Block List, Min: 1) (see [below for nested schema](#nestedblock--processing--processors--processor--smartscape_node--fields_to_extract--smartscape_field_extraction_entry))

<a id="nestedblock--processing--processors--processor--smartscape_node--fields_to_extract--smartscape_field_extraction_entry"></a>
### Nested Schema for `processing.processors.processor.smartscape_node.fields_to_extract.smartscape_field_extraction_entry`

Required:

- `referenced_field_name` (String) Referenced field name

Optional:

- `field_name` (String) Field name
- `strategy` (String) Strategy for field extraction. Possible values: `equals`, `startsWith`



<a id="nestedblock--processing--processors--processor--smartscape_node--node_name"></a>
### Nested Schema for `processing.processors.processor.smartscape_node.node_name`

Required:

- `type` (String) Type of value assignment. Possible values: `constant`, `field`, `multiValueConstant`

Optional:

- `constant` (String) Constant value
- `field` (Block List, Max: 1) Value from field (see [below for nested schema](#nestedblock--processing--processors--processor--smartscape_node--node_name--field))
- `multi_value_constant` (List of String) Constant multi value

<a id="nestedblock--processing--processors--processor--smartscape_node--node_name--field"></a>
### Nested Schema for `processing.processors.processor.smartscape_node.node_name.field`

Required:

- `source_field_name` (String) Source field name

Optional:

- `default_value` (String) Default value



<a id="nestedblock--processing--processors--processor--smartscape_node--static_edges_to_extract"></a>
### Nested Schema for `processing.processors.processor.smartscape_node.static_edges_to_extract`

Required:

- `smartscape_static_edge_extraction_entry` (Block List, Min: 1) (see [below for nested schema](#nestedblock--processing--processors--processor--smartscape_node--static_edges_to_extract--smartscape_static_edge_extraction_entry))

<a id="nestedblock--processing--processors--processor--smartscape_node--static_edges_to_extract--smartscape_static_edge_extraction_entry"></a>
### Nested Schema for `processing.processors.processor.smartscape_node.static_edges_to_extract.smartscape_static_edge_extraction_entry`

Required:

- `edge_type` (String) Edge type
- `target_id_field_name` (String) Target ID field name
- `target_type` (String) Target type




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

- `field_extraction_entry` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--processing--processors--processor--value_metric--dimensions--field_extraction_entry))

<a id="nestedblock--processing--processors--processor--value_metric--dimensions--field_extraction_entry"></a>
### Nested Schema for `processing.processors.processor.value_metric.dimensions.field_extraction_entry`

Optional:

- `constant_field_name` (String) Destination field name
- `constant_value` (String) Constant value to be assigned to field
- `default_value` (String) Default value
- `destination_field_name` (String) Destination field name
- `extraction_type` (String) Field value extraction type. Possible values: `constant`, `field`
- `source_field_name` (String) Source field name
- `strategy` (String) Strategy for field extraction. Possible values: `equals`, `startsWith`
