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
