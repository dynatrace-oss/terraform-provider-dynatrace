data "dynatrace_synthetic_location" "location" {
  name = "Location"
}

resource "dynatrace_credentials" "credentials_vault" {
  name        = "#name#"
  description = "my credentials vault"
  scopes      = ["SYNTHETIC"]
  username = "username"
  password = "password"
}

resource "dynatrace_http_monitor" "monitor" {
  name = "#name#"
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
      description = "portal.office.com"
      method = "GET"
      url = "https://portal.office.com"
      configuration {
        accept_any_certificate = true
      }
      authentication {
        type = "KERBEROS"
        credentials = dynatrace_credentials.credentials_vault.id
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
