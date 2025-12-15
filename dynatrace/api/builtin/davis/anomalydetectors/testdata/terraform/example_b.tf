resource "dynatrace_davis_anomaly_detectors" "anomaly" {
  description = ""
  enabled     = true
  source      = "Davis Anomaly Detection"
  title       = "Detector with OAuth requirement #name#"
  analyzer {
    name = "dt.statistics.ui.anomaly_detection.SeasonalBaselineAnomalyDetectionAnalyzer"
    input {
      analyzer_input_field {
        key   = "query.expression"
        value =<<-EOT
          fetch dt.system.query_executions
                                 | filter status == "SUCCEEDED"
                                 | fields timestamp, scanned_bytes, user.email
                                 | fieldsAdd gb = toDouble(scanned_bytes) / (1024*1024*1024)
                                 | fieldsAdd cost = gb * 0.00106
                                 | summarize  cost_per_user = sum(cost), by:{user.email}
                                 | filter cost_per_user >= 0.01
                                 | makeTimeseries datapoints = sum(cost_per_user)
                                   , by: {user.email}
                                   , time: now() - 4m
                                   , interval: 4m
        EOT
      }
      analyzer_input_field {
        key   = "tolerance"
        value = "1"
      }
      analyzer_input_field {
        key   = "dealertingSamples"
        value = "5"
      }
      analyzer_input_field {
        key   = "slidingWindow"
        value = "5"
      }
      analyzer_input_field {
        key   = "alertOnMissingData"
        value = "false"
      }
      analyzer_input_field {
        key   = "alertCondition"
        value = "BELOW"
      }
      analyzer_input_field {
        key   = "violatingSamples"
        value = "3"
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
        key   = "event.description"
        value = "Alert detected"
      }
      property {
        key   = "event.name"
        value = "Potential log data anomaly"
      }
    }
  }
  execution_settings {
    query_offset = 3
  }
}
