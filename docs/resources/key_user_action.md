---
layout: ""
page_title: dynatrace_key_user_action Resource - terraform-provider-dynatrace"
subcategory: "Web Applications"
description: |-
  The resource `dynatrace_key_user_action` covers configuration of Key User Actions for web applications
---

# dynatrace_key_user_action (Resource)

-> This resource requires the API token scopes **Read configuration** (`ReadConfig`), **Write configuration** (`WriteConfig`) and **Read Entities** (`entities.read`)

## Dynatrace Documentation

- RUM setup and configuration for web applications - https://www.dynatrace.com/support/help/how-to-use-dynatrace/real-user-monitoring/setup-and-configuration/web-applications

- Web application configuration API for Key User Actions - https://www.dynatrace.com/support/help/dynatrace-api/configuration-api/rum/web-application-configuration-api#edit-key-user-actions-list

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_key_user_action` downloads all existing Key User Actions

The full documentation of the export feature is available [here](https://registry.terraform.io/providers/dynatrace-oss/dynatrace/latest/docs/guides/export-v2).

## Resource Example Usage

The following example showcases how to manage Key User Actions separately from Web Applications using a dedicated resource.
While it is still possible to embed `key_user_actions` into the resource `dynatrace_web_application`, doing so is discouraged.

```terraform
resource "dynatrace_key_user_action" "#name#" {
  application_id = dynatrace_web_application.terrafoobar.id
  domain         = "120.0.0.1"
  name           = "Loading of page /custom"
  type           = "Load"
}


resource "dynatrace_web_application" "terrafoobar" {
  # key_user_actions {
  #   action {
  #     domain = "120.0.0.1"
  #     name   = "Loading of page /custom"
  #     type   = "Load"
  #   }
  # }    
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
  meta_data_capture_settings {
    capture {
      name           = "VisitTag1"
      type           = "JAVA_SCRIPT_VARIABLE"
      capturing_name = "PTC.navigation.GLOBAL_USER"
      # public_metadata = false 
      unique_id = 1
      # use_last_value = false 
    }
    capture {
      name           = "PageIdentity"
      type           = "CSS_SELECTOR"
      capturing_name = "#infoPageIdentityObjectIdentifier"
      # public_metadata = false 
      unique_id = 2
      # use_last_value = false 
    }
    capture {
      name           = "GCLID - Google Click Identifier"
      type           = "QUERY_STRING"
      capturing_name = "gclid"
      # public_metadata = false 
      unique_id = 3
      # use_last_value = false 
    }
    capture {
      name           = "Session ID"
      type           = "COOKIE"
      capturing_name = "RES_SESSIONID"
      # public_metadata = false 
      unique_id = 4
      # use_last_value = false 
    }
    capture {
      name           = "Tracking ID"
      type           = "COOKIE"
      capturing_name = "RES_TRACKINGID"
      # public_metadata = false 
      unique_id = 5
      # use_last_value = false 
    }
  }
  monitoring_settings {
    add_cross_origin_anonymous_attribute = true
    cache_control_header_optimizations   = true
    # cookie_placement_domain = "" 
    # correlation_header_inclusion_regex = "" 
    # custom_configuration_properties = "" 
    # exclude_xhr_regex = "" 
    # fetch_requests = false 
    injection_mode = "JAVASCRIPT_TAG"
    # library_file_location = "" 
    # monitoring_data_path = "" 
    script_tag_cache_duration_in_hours = 1
    # secure_cookie_attribute = false 
    # server_request_path_id = "" 
    # xml_http_request = false 
    advanced_javascript_tag_settings {
      # instrument_unsupported_ajax_frameworks = false 
      max_action_name_length = 100
      max_errors_to_capture  = 10
      # special_characters_to_escape = "" 
      # sync_beacon_firefox = false 
      # sync_beacon_internet_explorer = false 
      additional_event_handlers {
        # blur = false 
        # change = false 
        # click = false 
        max_dom_nodes = 5000
        # mouseup = false 
        # to_string_method = false 
        # use_mouse_up_event_for_clicks = false 
      }
      global_event_capture_settings {
        # additional_event_captured_as_user_input = "" 
        click       = true
        doubleclick = true
        keydown     = true
        keyup       = true
        mousedown   = true
        mouseup     = true
        scroll      = true
      }
    }
    content_capture {
      javascript_errors                 = true
      visually_complete_and_speed_index = true
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
      visually_complete_settings {
        # exclude_url_regex = "" 
        # ignored_mutations_list = "" 
        inactivity_timeout = 1000
        mutation_timeout   = 50
        threshold          = 50
      }
    }
    javascript_framework_support {
      # active_x_object = false 
      angular = true
      # dojo = false 
      extjs = true
      # icefaces = false 
      jquery = true
      # moo_tools = false 
      prototype = true
    }
  }
  session_replay_config {
    enabled                                = false
    cost_control_percentage                = 100
    enable_css_resource_capturing          = true
    css_resource_capturing_exclusion_rules = []
  }
  user_action_and_session_properties {
    property {
      type         = "STRING"
      aggregation  = "LAST"
      display_name = "GCLID - Google Click Identifier"
      id           = 2
      # ignore_case = false 
      key                       = "google_gclid"
      metadata_id               = 3
      origin                    = "META_DATA"
      store_as_session_property = true
      # store_as_user_action_property = false 
    }
    property {
      type         = "STRING"
      aggregation  = "LAST"
      display_name = "Session ID"
      id           = 3
      # ignore_case = false 
      key                       = "certona_session_id"
      metadata_id               = 4
      origin                    = "META_DATA"
      store_as_session_property = true
      # store_as_user_action_property = false 
    }
    property {
      type         = "STRING"
      aggregation  = "LAST"
      display_name = "Tracking ID"
      id           = 4
      # ignore_case = false 
      key                       = "certona_tracking_id"
      metadata_id               = 5
      origin                    = "META_DATA"
      store_as_session_property = true
      # store_as_user_action_property = false 
    }
  }
  user_action_naming_settings {
    ignore_case                    = true
    query_parameter_cleanups       = ["cfid", "phpsessid", "__sid", "cftoken", "sid"]
    split_user_actions_by_domain   = true
    use_first_detected_load_action = true
    load_action_naming_rules {
      rule {
        template = "Loading of {pageTitle (default)}"
        # use_or_conditions = false 
      }
    }
    placeholders {
      placeholder {
        name            = "TrailingURL"
        input           = "PAGE_URL"
        processing_part = "ALL"
        # use_guessed_element_identifier = false 
        processing_steps {
          step {
            type = "SUBSTRING"
            # fallback_to_input = false 
            pattern_after_search_type  = "LAST"
            pattern_before             = "/Windchill/app/#ptc1"
            pattern_before_search_type = "FIRST"
          }
        }
      }
      placeholder {
        name            = "PageIdentity"
        input           = "METADATA"
        metadata_id     = 2
        processing_part = "ALL"
        # use_guessed_element_identifier = false 
      }
    }
    xhr_action_naming_rules {
      rule {
        template = "{pageTitle (default)}"
        # use_or_conditions = false 
      }
    }
  }
  user_tags {
    tag {
      id = 1
      # ignore_case = false 
      metadata_id = 1
    }
  }
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `application_id` (String) The ID of the WebApplication
- `name` (String) The name of the action
- `type` (String) The type of the action. Possible values are `Custom`, `Load` and `Xhr`

### Optional

- `domain` (String) The domain where the action is performed

### Read-Only

- `id` (String) The ID of this resource.
 