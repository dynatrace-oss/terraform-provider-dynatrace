
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

resource "dynatrace_calculated_web_metric" "dimensions" {
  name           = "#name#"
  enabled        = true
  app_identifier = dynatrace_web_application.application.id
  metric_key     = "calc:apps.web.#name#"
  metric_definition {
    metric = "UserActionDuration"
  }
}
