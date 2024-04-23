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
  }
}
