resource "dynatrace_http_monitor" "#name#" {
  name = "#name#"
  frequency = 1
  locations = ["GEOLOCATION-03E96F97A389F96A","GEOLOCATION-9999453BE4BDB3CD","GEOLOCATION-2FD31C834DE4D601","GEOLOCATION-924D253001531722","GEOLOCATION-7F39AED31559436D","GEOLOCATION-DDAA176627F5667A"]
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
      description = "portal.office.com"
      method = "GET"
      url = "https://portal.office.com"
      configuration {
        accept_any_certificate = true
      }
      authentication {
        type = "KERBEROS"
        credentials = "CREDENTIALS_VAULT-4DFB50E5F50A21A4"
        realm_name = "ABCDE"
        kdc_ip = "10.0.0.1"
      }
    }
  }
  tags {
    tag {
      context = "CONTEXTLESS"
      key = "Office365"
      source = "USER"
    }
  }
}
