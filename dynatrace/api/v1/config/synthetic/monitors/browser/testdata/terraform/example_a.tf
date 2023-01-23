resource "dynatrace_browser_monitor" "#name#" {
  name                   = "#name#"
  frequency              = 15
  locations              = ["GEOLOCATION-B4B9167CAAA88F6A", "GEOLOCATION-03E96F97A389F96A"]
  manually_assigned_apps = ["APPLICATION-EA7C4B59F27D43EB"]
  anomaly_detection {
    loading_time_thresholds {
      enabled = true
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
          name  = "kjh"
          value = "kjh"
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
        description = "Loading of \"https://www.orf.at\""
        navigate {
          url = "https://www.orf.at"
          authentication {
            type  = "http_authentication"
            creds = "CREDENTIALS_VAULT-26F62024BC3ABBCF"
          }
          wait {
            wait_for = "page_complete"
          }
        }
      }
      event {
        description = "jhjhjh"
        navigate {
          url = "https://www.orf.at"
          authentication {
            type  = "http_authentication"
            creds = "CREDENTIALS_VAULT-26F62024BC3ABBCF"
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
}
