# Example with monitoring_data_path set by API
resource "dynatrace_web_application" "#name#" {
  cost_control_user_session_percentage = 100
  load_action_key_performance_metric   = "VISUALLY_COMPLETE"

  name                              = "#name#"
  type                              = "MANUALLY_INJECTED"
  xhr_action_key_performance_metric = "VISUALLY_COMPLETE"
  real_user_monitoring_enabled      = true

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
    injection_mode = "JAVASCRIPT_TAG"

    advanced_javascript_tag_settings {
      max_action_name_length = 100
      max_errors_to_capture  = 10
      additional_event_handlers {
        max_dom_nodes = 5000
      }

      global_event_capture_settings {
        change      = true
        click       = true
        doubleclick = true
        keydown     = true
        keyup       = true
        mousedown   = true
        mouseup     = true
        scroll      = true
        touch_end   = true
        touch_start = true
      }
    }

    xml_http_request = true
    fetch_requests   = true
    content_capture {
      resource_timing_settings {
        w3c_resource_timings  = true
        instrumentation_delay = 0
      }

      timeout_settings {
        temporary_action_limit         = 3
        temporary_action_total_timeout = 300000
      }

      javascript_errors                 = true
      visually_complete_and_speed_index = true
    }
    secure_cookie_attribute = true
  }

  waterfall_settings {
    resource_browser_caching_threshold            = 50
    resources_threshold                           = 100000
    slow_cnd_resources_threshold                  = 200
    slow_first_party_resources_threshold          = 200
    slow_third_party_resources_threshold          = 200
    speed_index_visually_complete_ratio_threshold = 50
    uncompressed_resources_threshold              = 860
  }

  xhr_action_apdex_settings {
    frustrating_fallback_threshold = 12000
    frustrating_threshold          = 12000
    tolerated_fallback_threshold   = 3000
    tolerated_threshold            = 3000
  }
  user_action_naming_settings {
  }
}
