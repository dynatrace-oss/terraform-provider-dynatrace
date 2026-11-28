---
layout: ""
page_title: "dynatrace_davis_anomaly_detectors Resource - terraform-provider-dynatrace"
subcategory: "Anomaly Detection"
description: |-
  The resource `dynatrace_davis_anomaly_detectors` covers configuration for Davis anomaly detectors
---

# dynatrace_davis_anomaly_detectors (Resource)

-> **Dynatrace SaaS only**

-> This resource requires an OAuth client or platform token configured with the permissions outlined in the [Davis Anomaly Detection](https://docs.dynatrace.com/docs/shortlink/davis-ai-anomaly-detection-app#expand--prerequisites--1) app documentation.
Please set the environment variables `DT_CLIENT_ID` and `DT_CLIENT_SECRET`, or alternatively `DT_PLATFORM_TOKEN`.

-> Depending on the anomaly detector configuration, additional **storage permissions** may be required for DQL-related access (e.g. `storage:bizevents:read`, `storage:logs:read`, `storage:entities:read`).

## Dynatrace Documentation

- Davis Anomaly Detection App - https://docs.dynatrace.com/docs/platform/davis-ai/anomaly-detection/anomaly-detection-app

- Settings API - https://www.dynatrace.com/support/help/dynatrace-api/environment-api/settings (schemaId: `builtin:davis.anomaly-detectors`)

- Static thresholds for anomaly detection - https://docs.dynatrace.com/docs/dynatrace-intelligence/anomaly-detection/static-thresholds

- Auto-adaptive thresholds for anomaly detection - https://docs.dynatrace.com/docs/dynatrace-intelligence/anomaly-detection/auto-adaptive-threshold

- Anomaly detection configuration (seasonal baseline and record analyzers) - https://docs.dynatrace.com/docs/dynatrace-intelligence/anomaly-detection/anomaly-detection-configuration

- Davis Analyzers API (available analyzers and their input fields) - https://########.apps.dynatrace.com/platform/davis/analyzers/v1/analyzers

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_davis_anomaly_detectors` downloads all existing Davis anomaly detector configuration

The full documentation of the export feature is available [here](https://dt-url.net/h203qmc).

## Resource Example Usage

### Static Threshold

```terraform
resource "dynatrace_davis_anomaly_detectors" "#name#" {
  description = "Sample Description"
  enabled     = false
  source      = "Davis Anomaly Detection"
  title       = "#name#"
  analyzer {
    name = "dt.statistics.ui.anomaly_detection.StaticThresholdAnomalyDetectionAnalyzer"
    input {
      analyzer_input_field {
        key   = "alertCondition"
        value = "ABOVE"
      }
      analyzer_input_field {
        key   = "alertOnMissingData"
        value = "false"
      }
      analyzer_input_field {
        key   = "violatingSamples"
        value = "3"
      }
      analyzer_input_field {
        key   = "slidingWindow"
        value = "5"
      }
      analyzer_input_field {
        key   = "dealertingSamples"
        value = "5"
      }
      analyzer_input_field {
        key   = "query"
        value = "fetch bizevents"
      }
      analyzer_input_field {
        key   = "threshold"
        value = "12345678"
      }
    }
  }
  event_template {
    properties {
      property {
        key   = "dt.source_entity"
        value = "{dims:dt.source_entity}"
      }
      property {
        key   = "event.type"
        value = "CUSTOM_ALERT"
      }
      property {
        key   = "event.name"
        value = "Event Name Example"
      }
      property {
        key   = "event.description"
        value = "Event Description Example"
      }
    }
  }
  execution_settings {
    actor = dynatrace_iam_service_user.user.id
  }
}

resource "dynatrace_iam_service_user" "user" {
  name = "#name#"
}
```

### Record Analyzer

```terraform
resource "dynatrace_davis_anomaly_detectors" "#name#" {
  title       = "#name#"
  description = "Raises one event per service with more than 100 ERROR logs in the last hour"
  enabled     = false
  source      = "Davis Anomaly Detection"
  analyzer {
    name = "dt.statistics.anomaly_detection.RecordAnomalyDetectionAnalyzer"
    input {
      analyzer_input_field {
        key   = "query.expression"
        value = <<-EOT
          fetch logs, from: -1h
          | filter loglevel == "ERROR"
          | summarize error_count = count(), by:{service.name}
          | filter error_count > 100
        EOT
      }
      analyzer_input_field {
        key   = "alertIdentityFields[0]"
        value = "service.name"
      }
    }
  }
  event_template {
    properties {
      property {
        key   = "event.type"
        value = "ERROR_EVENT"
      }
      property {
        key   = "event.name"
        value = "{error_count} errors on {service.name}"
      }
      property {
        key   = "event.description"
        value = "Service {service.name} logged {error_count} ERROR lines in the last hour"
      }
    }
  }
  execution_settings {
    actor = dynatrace_iam_service_user.user.id
    delay = 60
  }
}

resource "dynatrace_iam_service_user" "user" {
  name = "#name#"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `analyzer` (Block List, Min: 1, Max: 1) Analyzer input to initialize the analyzer (see [below for nested schema](#nestedblock--analyzer))
- `description` (String) The description of the anomaly detector
- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)
- `event_template` (Block List, Min: 1, Max: 1) Defines additional fields on the davis events triggered by the anomaly detector (see [below for nested schema](#nestedblock--event_template))
- `execution_settings` (Block List, Min: 1, Max: 1) Defines the configuration parameters that influence how and under what context a query or evaluation is executed. (see [below for nested schema](#nestedblock--execution_settings))
- `source` (String) The source which created the anomaly detector
- `title` (String) The title of the anomaly detector

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--analyzer"></a>
### Nested Schema for `analyzer`

Required:

- `name` (String) Fully qualified name of the analyzer

Optional:

- `input` (Block List, Max: 1) Input fields for the specified analyzer (see [below for nested schema](#nestedblock--analyzer--input))

<a id="nestedblock--analyzer--input"></a>
### Nested Schema for `analyzer.input`

Required:

- `analyzer_input_field` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--analyzer--input--analyzer_input_field))

<a id="nestedblock--analyzer--input--analyzer_input_field"></a>
### Nested Schema for `analyzer.input.analyzer_input_field`

Required:

- `key` (String) Analyzer input field key
- `value` (String) Analyzer input field value




<a id="nestedblock--event_template"></a>
### Nested Schema for `event_template`

Optional:

- `properties` (Block List, Max: 1) Set of additional key-value properties to be attached to the triggered event. (see [below for nested schema](#nestedblock--event_template--properties))

<a id="nestedblock--event_template--properties"></a>
### Nested Schema for `event_template.properties`

Required:

- `property` (Block List, Min: 1) (see [below for nested schema](#nestedblock--event_template--properties--property))

<a id="nestedblock--event_template--properties--property"></a>
### Nested Schema for `event_template.properties.property`

Required:

- `key` (String) Property key
- `value` (String) Property value. Supports substitution of placeholders placed in curly braces {}.




<a id="nestedblock--execution_settings"></a>
### Nested Schema for `execution_settings`

Optional:

- `actor` (String) UUID of a service user. Queries will be executed on behalf of the service user.
- `delay` (Number) Fixed delay between executions (in seconds)
- `query_offset` (Number) Minute offset of sliding evaluation window for metrics with latency
