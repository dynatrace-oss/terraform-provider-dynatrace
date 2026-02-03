resource "dynatrace_synthetic_location" "location" {
  name                                  = "#name#"
  city                                  = "San Francisco de Asis"
  country_code                          = "VE"
  region_code                           = "04"
  deployment_type                       = "STANDARD"
  latitude                              = 10.0756
  location_node_outage_delay_in_minutes = 3
  longitude                             = -67.5442
}

resource "dynatrace_credentials" "credentials_vault" {
  name        = "#name#"
  description = "my credentials vault"
  scopes      = ["SYNTHETIC"]
  username = "username"
  password = "password"
}

resource "time_sleep" "wait_5_seconds" {
  depends_on = [dynatrace_synthetic_location.location]
  create_duration = "5s"
}

resource "dynatrace_http_monitor" "monitor" {
  depends_on = [time_sleep.wait_5_seconds]
  name = "#name#"
  frequency = 1
  locations = [dynatrace_synthetic_location.location.id]
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
