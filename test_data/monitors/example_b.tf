resource "dynatrace_browser_monitor" "#name#" {
  name      = "#name#"
  frequency = 15
  locations = ["GEOLOCATION-57F63BAD1C6A415C"]
  anomaly_detection {
    loading_time_thresholds {
      enabled = true
      thresholds {
        threshold {
          type        = "ACTION"
          event_index = 1
          value_ms    = 10000
        }
      }
    }
    outage_handling {
      global_outage  = true
      local_outage   = true
      retry_on_error = true
      local_outage_policy {
        affected_locations = 1
        consecutive_runs   = 5
      }
      global_outage_policy {
        consecutive_runs = 1
      }
    }
  }
  key_performance_metrics {
    load_action_kpm = "CUMULATIVE_LAYOUT_SHIFT"
    xhr_action_kpm  = "RESPONSE_END"
  }
  script {
    type = "availability"
    configuration {
      block          = ["https://www.google.at/"]
      bypass_csp     = true
      monitor_frames = true
      user_agent     = "reini-user-agent"
      bandwidth {
        network_type = "WiFi"
      }
      cookies {
        cookie {
          name   = "cookiea"
          domain = "google.com"
          path   = "/"
          value  = "b"
        }
      }
      device {
        height        = 1080
        mobile        = true
        scale_factor  = 2
        touch_enabled = true
        width         = 1920
      }
      headers {
        restrictions = ["https://www.google.at/"]
        header {
          name  = "X-foo"
          value = "bar"
        }
      }
      ignored_error_codes {
        matching_document_requests = "asdf*"
        status_codes               = "567"
      }
      javascript_setttings {
        custom_properties = "cux=1"
        timeout_settings {
          action_limit  = 4
          total_timeout = 102
        }
        visually_complete_options {
          excluded_elements    = ["#some", "#boo"]
          excluded_urls        = ["logout", "asdf"]
          image_size_threshold = 51
          inactivity_timeout   = 1001
          mutation_timeout     = 51
        }
      }
    }
    events {
      event {
        description = "Loading of \"https://www.google.at\""
        navigate {
          url = "https://www.google.at"
          authentication {
            type  = "webform"
            creds = "CREDENTIALS_VAULT-26F62024BC3ABBCF"
          }
          validate {
            validation {
              type  = "element_match"
              match = "ml,klj"
              regex = true
              target {
                window = "klkkj"
                locators {
                  locator {
                    type  = "css"
                    value = "kjjk"
                  }
                }
              }
            }
          }
          wait {
            wait_for = "page_complete"
          }
        }
      }
    }
  }
}
