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

resource "dynatrace_web_application" "application" {
  name                                 = "#name#"
  type                                 = "AUTO_INJECTED"
  cost_control_user_session_percentage = 100
  load_action_key_performance_metric   = "VISUALLY_COMPLETE"
  real_user_monitoring_enabled         = true
  xhr_action_key_performance_metric    = "VISUALLY_COMPLETE"
  custom_action_apdex_settings {
    frustrating_fallback_threshold = 12000
    frustrating_threshold          = 12000
    tolerated_fallback_threshold   = 3000
    tolerated_threshold            = 3000
  }
  load_action_apdex_settings {
    frustrating_fallback_threshold = 12000
    frustrating_threshold          = 12000
    tolerated_fallback_threshold   = 3000
    tolerated_threshold            = 3000
  }
  monitoring_settings {
    add_cross_origin_anonymous_attribute = true
    cache_control_header_optimizations   = true
    injection_mode = "JAVASCRIPT_TAG"
    script_tag_cache_duration_in_hours = 1
    advanced_javascript_tag_settings {
      max_action_name_length = 100
      max_errors_to_capture  = 10
      additional_event_handlers {
        max_dom_nodes = 5000
      }
    }
    content_capture {
      resource_timing_settings {
        instrumentation_delay    = 53
        non_w3c_resource_timings = true
        w3c_resource_timings     = true
      }
      timeout_settings {
        temporary_action_limit         = 3
        temporary_action_total_timeout = 100
        timed_action_support           = true
      }
    }
  }
  user_action_naming_settings {}
  waterfall_settings {
    resource_browser_caching_threshold            = 50
    resources_threshold                           = 100000
    slow_cnd_resources_threshold                  = 200000
    slow_first_party_resources_threshold          = 200000
    slow_third_party_resources_threshold          = 200000
    speed_index_visually_complete_ratio_threshold = 50
    uncompressed_resources_threshold              = 860
  }
  xhr_action_apdex_settings {
    frustrating_fallback_threshold = 12000
    frustrating_threshold          = 12000
    tolerated_fallback_threshold   = 3000
    tolerated_threshold            = 3000
  }
}

resource "dynatrace_credentials" "credentials_vault" {
  name        = "#name#"
  description = "my credentials vault"
  scopes      = ["SYNTHETIC"]
  username = "username"
  password = "password"
}

resource "time_sleep" "wait_5_seconds" {
  depends_on = [dynatrace_synthetic_location.location, dynatrace_web_application.application, dynatrace_credentials.credentials_vault]
  create_duration = "5s"
}

resource "dynatrace_browser_monitor" "monitor" {
  depends_on = [time_sleep.wait_5_seconds]
  name                   = "#name#"
  frequency              = 15
  locations              = [dynatrace_synthetic_location.location.id]
  manually_assigned_apps = [dynatrace_web_application.application.id]
  anomaly_detection {
    loading_time_thresholds {
      enabled = true
      thresholds {
        threshold {
          type  = "TOTAL"
          value_ms = 1000
          event_index = 0
        }
        # updated
        threshold {
          type  = "ACTION"
          value_ms = 3000
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
          name  = "nameEdit"
          value = "valueEdit"
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
    # updated
    tag {
      context = "CONTEXTLESS"
      key = "keyEdit"
      source = "USER"
    }
  }
}

resource "time_sleep" "wait_for_monitor" {
  depends_on = [dynatrace_browser_monitor.monitor]
  create_duration = "10s"
}
