data "dynatrace_synthetic_location" "location" {
  name = "Location"
}

resource "dynatrace_http_monitor" "advanced" {
  name      = "#name#"
  frequency = 1
  locations = [data.dynatrace_synthetic_location.location.id]
  anomaly_detection {
    loading_time_thresholds {
    }
    outage_handling {
      global_outage = true
      global_outage_policy {
        consecutive_runs = 1
      }
    }
  }
  script {
    request {
      description     = "getOffice365ActiveUserCounts"
      method          = "GET"
      url             = "https://graph.microsoft.com/v1.0/reports/getOffice365ActiveUserCounts(period='D7')"
      configuration {
        accept_any_certificate = true
        follow_redirects       = true
        headers {
          header {
            name  = "name1"
            value = "value1"
          }
          # to update
          header {
            name  = "name2"
            value = "value2"
          }
        }
      }
    }
    custom_properties {
      custom_property {
        name = "hmRequestTimeoutInMs"
        value = "5000"
      }
      # to update
      custom_property {
        name = "hmConnectTimeoutInMs"
        value = "6000"
      }
    }
  }
}
