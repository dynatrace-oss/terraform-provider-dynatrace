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

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_davis_anomaly_detectors` downloads all existing Davis anomaly detector configuration

The full documentation of the export feature is available [here](https://dt-url.net/h203qmc).

## Analyzers

The `analyzer` block chooses how the detector is evaluated. `analyzer.name` is the fully qualified analyzer name, and `analyzer.input` holds its settings as `analyzer_input_field` key/value pairs. Values are always strings (`"5"`, `"false"`). The provider does not check the keys; Dynatrace validates them on create and update.

| Analyzer name | Evaluates | Input keys |
|---|---|---|
| `dt.statistics.ui.anomaly_detection.StaticThresholdAnomalyDetectionAnalyzer` | a timeseries against a fixed value | `query`, `threshold`, `alertCondition` (`ABOVE`, `BELOW`), `alertOnMissingData`, `violatingSamples`, `slidingWindow`, `dealertingSamples` |
| `dt.statistics.ui.anomaly_detection.AutoAdaptiveAnomalyDetectionAnalyzer` | a timeseries against a threshold Davis learns | `query`, `numberOfSignalFluctuations`, `alertCondition` (`ABOVE`, `BELOW`, `OUTSIDE`), `alertOnMissingData`, `violatingSamples`, `slidingWindow`, `dealertingSamples` |
| `dt.statistics.ui.anomaly_detection.SeasonalBaselineAnomalyDetectionAnalyzer` | a timeseries against a seasonal confidence band | `query`, `tolerance`, `alertCondition` (`ABOVE`, `BELOW`, `OUTSIDE`), `alertOnMissingData`, `violatingSamples`, `slidingWindow`, `dealertingSamples` |
| `dt.statistics.anomaly_detection.RecordAnomalyDetectionAnalyzer` | any DQL result, one row per violation | `query.expression`, `alertIdentityFields[N]`, `query.filterSegments[0].id` |

-> The three timeseries analyzers live under `dt.statistics.ui.anomaly_detection`. The record analyzer does not have `.ui` in its name.

-> The timeseries analyzers accept the DQL in `query` or in `query.expression`.

### Record analyzer

With `RecordAnomalyDetectionAnalyzer` the DQL query is the alert condition, and it can be any query (`fetch logs`, `fetch events`, a `timeseries` flattened with `arrayAvg`), not just a timeseries. Put the condition in a `filter`: no row means no alert, and each row raises one Davis event.

- `query.expression` is the DQL query.
- `alertIdentityFields[0]`, `alertIdentityFields[1]`, ... name the result columns that identify one row. A row with the same values on the next evaluation keeps the same alert open, and a row that disappears closes it. Without them, every evaluation raises a new event for every row.
- `query.filterSegments[0].id` restricts the query to a filter segment.
- `execution_settings.delay` is the fixed delay between two evaluations, in seconds. With a long window such as `from: -1h`, the same rows come back on every evaluation, which is why `alertIdentityFields[N]` matters.
- `alertCondition`, `violatingSamples`, `slidingWindow`, `dealertingSamples` and `alertOnMissingData` are not used by this analyzer.
- In `event_template`, every result column is available as `{column}` in `event.name` and `event.description`. For a column that holds an entity ID, `{dims:column}` gives the entity's display name, for example `dt.source_entity = "{dims:dt.entity.host}"`.

## Resource Example Usage

### Static threshold

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

### Record analyzer

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
