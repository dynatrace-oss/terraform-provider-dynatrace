data "dynatrace_synthetic_location" "location" {
  name = "Location"
}

data "dynatrace_application" "web_application" {
  name = "Web Application"
}

resource "dynatrace_credentials" "credentials_vault" {
  name        = "#name#"
  description = "my credentials vault"
  scopes      = ["SYNTHETIC"]
  username = "username"
  password = "password"
}

resource "dynatrace_browser_monitor" "monitor" {
  name                   = "#name#"
  frequency              = 15
  locations              = [data.dynatrace_synthetic_location.location.id]
  manually_assigned_apps = [data.dynatrace_application.web_application.id]
  anomaly_detection {
    loading_time_thresholds {
      enabled = true
      thresholds {
        threshold {
          type  = "TOTAL"
          value_ms = 1000
          event_index = 0
        }
        # to update
        threshold {
          type  = "ACTION"
          value_ms = 2000
          event_index = 2
        }
      }
    }
    outage_handling {
      global_outage  = true
      retry_on_error = true
      global_outage_policy {
        consecutive_runs = 1
      }
    }
  }
  key_performance_metrics {
    load_action_kpm = "VISUALLY_COMPLETE"
    xhr_action_kpm  = "VISUALLY_COMPLETE"
  }
  script {
    type = "clickpath"
    configuration {
      bypass_csp = true
      user_agent = "Mozilla"
      bandwidth {
        network_type = "GPRS"
      }
      device {
        name        = "Apple iPhone 8"
        orientation = "landscape"
      }
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
      ignored_error_codes {
        status_codes = "400"
      }
      javascript_setttings {
        timeout_settings {
          action_limit  = 3
          total_timeout = 100
        }
        visually_complete_options {
          image_size_threshold = 0
          inactivity_timeout   = 0
          mutation_timeout     = 0
        }
      }
    }
    events {
      event {
        description = "Loading of \"https://example.com\""
        navigate {
          url = "https://example.com"
          authentication {
            type  = "http_authentication"
            creds = dynatrace_credentials.credentials_vault.id
          }
          wait {
            wait_for = "page_complete"
          }
        }
      }
      event {
        description = "jhjhjh"
        navigate {
          url = "https://example.com"
          authentication {
            type  = "http_authentication"
            creds = dynatrace_credentials.credentials_vault.id
          }
          validate {
            validation {
              type  = "text_match"
              match = "kkl"
              regex = true
              target {
                window = "k"
              }
            }
          }
          wait {
            timeout  = 60000
            wait_for = "validation"
            validation {
              type  = "element_match"
              match = "kjkj"
              target {
                locators {
                  locator {
                    type  = "css"
                    value = "jjj"
                  }
                }
              }
            }
          }
        }
      }
      event {
        description = "fvf"
        click {
          button = 0
          validate {
            validation {
              type = "text_match"
            }
          }
          wait {
            wait_for = "page_complete"
          }
        }
      }
      event {
        description = "jsfoo"
        javascript {
          code = <<-EOT
            let x = 3;
            for (var i = 0; i < x; x++) {
                console.log("asdf");
            }
          EOT
          wait {
            wait_for = "page_complete"
          }
        }
      }
    }
  }
  tags {
    tag {
      context = "CONTEXTLESS"
      key = "key1"
      source = "USER"
    }
    # to update
    tag {
      context = "CONTEXTLESS"
      key = "key2"
      source = "USER"
    }
  }
}
