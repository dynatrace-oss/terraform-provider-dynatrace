---
layout: ""
page_title: dynatrace_web_application Resource - terraform-provider-dynatrace"
subcategory: "Web Applications"
description: |-
  The resource `dynatrace_web_application` covers configuration for web applications
---

# dynatrace_web_application (Resource)

-> This resource requires the API token scopes **Read configuration** (`ReadConfig`) and **Write configuration** (`WriteConfig`)

## Dynatrace Documentation

- RUM setup and configuration for web applications - https://www.dynatrace.com/support/help/how-to-use-dynatrace/real-user-monitoring/setup-and-configuration/web-applications

- Web application configuration API - https://www.dynatrace.com/support/help/dynatrace-api/configuration-api/rum/web-application-configuration-api

## Environment Variables (Optional)

There may be a delay for this resource to be fully available as a dependency for a subsequent resource. E.g. Utilizing this resource and application detection rules together.
 
A default polling mechanism exists to validate the creation but may require adjustment due to load. The following environment variable can be used to fine tune this setting.

- `DYNATRACE_CREATE_CONFIRM_WEB_APPLICATION` (Default: 60, Max: 300) configures the number of successful consecutive retries expected.

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_web_application` downloads all existing web application configuration

The full documentation of the export feature is available [here](https://dt-url.net/h203qmc).

## Resource Example Usage

```terraform
resource "dynatrace_web_application" "#name#" {
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

- `cost_control_user_session_percentage` (Number) (Field has overlap with `dynatrace_web_app_enablement`) Analize *X*% of user sessions
- `custom_action_apdex_settings` (Block List, Min: 1, Max: 1) Defines the Custom Action Apdex settings of an application (see [below for nested schema](#nestedblock--custom_action_apdex_settings))
- `load_action_apdex_settings` (Block List, Min: 1, Max: 1) Defines the Load Action Apdex settings of an application (see [below for nested schema](#nestedblock--load_action_apdex_settings))
- `load_action_key_performance_metric` (String) The key performance metric of load actions. Possible values are `ACTION_DURATION`, `CUMULATIVE_LAYOUT_SHIFT`, `DOM_INTERACTIVE`, `FIRST_INPUT_DELAY`, `LARGEST_CONTENTFUL_PAINT`, `LOAD_EVENT_END`, `LOAD_EVENT_START`, `RESPONSE_END`, `RESPONSE_START`, `SPEED_INDEX` and `VISUALLY_COMPLETE`
- `monitoring_settings` (Block List, Min: 1, Max: 1) Real user monitoring settings (see [below for nested schema](#nestedblock--monitoring_settings))
- `name` (String) The name of the web application, displayed in the UI
- `type` (String) The type of the web application. Possible values are `AUTO_INJECTED`, `BROWSER_EXTENSION_INJECTED` and `MANUALLY_INJECTED`
- `waterfall_settings` (Block List, Min: 1, Max: 1) These settings influence the monitoring data you receive for 3rd party, CDN, and 1st party resources (see [below for nested schema](#nestedblock--waterfall_settings))
- `xhr_action_apdex_settings` (Block List, Min: 1, Max: 1) Defines the XHR Action Apdex settings of an application (see [below for nested schema](#nestedblock--xhr_action_apdex_settings))
- `xhr_action_key_performance_metric` (String) The key performance metric of XHR actions. Possible values are `ACTION_DURATION`, `RESPONSE_END`, `RESPONSE_START` and `VISUALLY_COMPLETE`.

### Optional

- `conversion_goals` (Block List, Max: 1) A list of conversion goals of the application (see [below for nested schema](#nestedblock--conversion_goals))
- `key_user_actions` (Block List, Deprecated) User Action names to be flagged as Key User Actions (see [below for nested schema](#nestedblock--key_user_actions))
- `meta_data_capture_settings` (Block List, Max: 1) Java script agent meta data capture settings (see [below for nested schema](#nestedblock--meta_data_capture_settings))
- `real_user_monitoring_enabled` (Boolean) (Field has overlap with `dynatrace_web_app_enablement`) Real user monitoring enabled/disabled
- `session_replay_config` (Block List, Max: 1) Settings regarding Session Replay (see [below for nested schema](#nestedblock--session_replay_config))
- `url_injection_pattern` (String) URL injection pattern for manual web application
- `user_action_and_session_properties` (Block List, Max: 1) User action and session properties settings (see [below for nested schema](#nestedblock--user_action_and_session_properties))
- `user_action_naming_settings` (Block List, Max: 1) The settings of user action naming (see [below for nested schema](#nestedblock--user_action_naming_settings))
- `user_tags` (Block List, Max: 1) User tags settings (see [below for nested schema](#nestedblock--user_tags))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--custom_action_apdex_settings"></a>
### Nested Schema for `custom_action_apdex_settings`

Required:

- `frustrating_fallback_threshold` (Number) Fallback threshold of an XHR action, defining a tolerable user experience, when the configured KPM is not available. Values between 0 and 240000 are allowed.
- `frustrating_threshold` (Number) Maximal value of apdex, which is considered as tolerable user experience. Values between 0 and 240000 are allowed.
- `tolerated_fallback_threshold` (Number) Fallback threshold of an XHR action, defining a satisfied user experience, when the configured KPM is not available. Values between 0 and 60000 are allowed.
- `tolerated_threshold` (Number) Maximal value of apdex, which is considered as satisfied user experience. Values between 0 and 60000 are allowed.

Optional:

- `threshold` (Number, Deprecated) no documentation available


<a id="nestedblock--load_action_apdex_settings"></a>
### Nested Schema for `load_action_apdex_settings`

Required:

- `frustrating_fallback_threshold` (Number) Fallback threshold of an XHR action, defining a tolerable user experience, when the configured KPM is not available. Values between 0 and 240000 are allowed.
- `frustrating_threshold` (Number) Maximal value of apdex, which is considered as tolerable user experience. Values between 0 and 240000 are allowed.
- `tolerated_fallback_threshold` (Number) Fallback threshold of an XHR action, defining a satisfied user experience, when the configured KPM is not available. Values between 0 and 60000 are allowed.
- `tolerated_threshold` (Number) Maximal value of apdex, which is considered as satisfied user experience. Values between 0 and 60000 are allowed.

Optional:

- `threshold` (Number, Deprecated) no documentation available


<a id="nestedblock--monitoring_settings"></a>
### Nested Schema for `monitoring_settings`

Required:

- `advanced_javascript_tag_settings` (Block List, Min: 1, Max: 1) Advanced JavaScript tag settings (see [below for nested schema](#nestedblock--monitoring_settings--advanced_javascript_tag_settings))
- `content_capture` (Block List, Min: 1, Max: 1) Settings for content capture (see [below for nested schema](#nestedblock--monitoring_settings--content_capture))
- `injection_mode` (String) Possible valures are `CODE_SNIPPET`, `CODE_SNIPPET_ASYNC`, `INLINE_CODE`, `JAVASCRIPT_TAG`, `JAVASCRIPT_TAG_COMPLETE`, `JAVASCRIPT_TAG_SRI`

Optional:

- `add_cross_origin_anonymous_attribute` (Boolean) Add the cross origin = anonymous attribute to capture JavaScript error messages and W3C resource timings
- `angular_package_name` (String) The name of the angular package
- `browser_restriction_settings` (Block List, Max: 1) Settings for restricting certain browser type, version, platform and, comparator. It also restricts the mode (see [below for nested schema](#nestedblock--monitoring_settings--browser_restriction_settings))
- `cache_control_header_optimizations` (Boolean) Optimize the value of cache control headers for use with Dynatrace real user monitoring enabled/disabled
- `cookie_placement_domain` (String) Domain for cookie placement. Maximum 150 characters.
- `correlation_header_inclusion_regex` (String) To enable RUM for XHR calls to AWS Lambda, define a regular expression matching these calls, Dynatrace can then automatically add a custom header (`x-dtc`) to each such request to the respective endpoints in AWS.

Important: These endpoints must accept the `x-dtc` header, or the requests will fail
- `custom_configuration_properties` (String) The location to send monitoring data from the JavaScript tag.

 Specify either a relative or an absolute URL. If you use an absolute URL, data will be sent using CORS. 

 **Required** for auto-injected applications, optional for agentless applications. Maximum 512 characters.
- `exclude_xhr_regex` (String) You can exclude some actions from becoming XHR actions.

Put a regular expression, matching all the required URLs, here.

If noting specified the feature is disabled
- `fetch_requests` (Boolean) `fetch()` request capture enabled/disabled
- `ignore_ip_address_restriction_settings` (Boolean) Manage IP address exclusion settings with `dynatrace_web_app_ip_address_exclusion` resource
- `instrumented_web_server` (Boolean) Instrumented web or app server.
- `ip_address_restriction_settings` (Block List, Max: 1) Settings for restricting certain ip addresses and for introducing subnet mask. It also restricts the mode (see [below for nested schema](#nestedblock--monitoring_settings--ip_address_restriction_settings))
- `javascript_framework_support` (Block List, Max: 1) Support of various JavaScript frameworks (see [below for nested schema](#nestedblock--monitoring_settings--javascript_framework_support))
- `javascript_injection_rules` (Block List, Max: 1) Java script injection rules (see [below for nested schema](#nestedblock--monitoring_settings--javascript_injection_rules))
- `library_file_from_cdn` (Boolean) Get the JavaScript library file from the CDN. Not supported by agentless applications and assumed to be false for auto-injected applications if omitted.
- `library_file_location` (String) The location of your application’s custom JavaScript library file. 

 If nothing specified the root directory of your web server is used. 

 **Required** for auto-injected applications, not supported by agentless applications. Maximum 512 characters.
- `monitoring_data_path` (String) The location to send monitoring data from the JavaScript tag.

 Specify either a relative or an absolute URL. If you use an absolute URL, data will be sent using CORS. 

 **Required** for auto-injected applications, optional for agentless applications. Maximum 512 characters.
- `same_site_cookie_attribute` (String) Same site cookie attribute
- `script_tag_cache_duration_in_hours` (Number) Time duration for the cache settings
- `secure_cookie_attribute` (Boolean) Secure attribute usage for Dynatrace cookies enabled/disabled
- `server_request_path_id` (String) Path to identify the server’s request ID. Maximum 150 characters.
- `use_cors` (Boolean) Send beacon data via CORS.
- `xml_http_request` (Boolean) `XmlHttpRequest` support enabled/disabled

<a id="nestedblock--monitoring_settings--advanced_javascript_tag_settings"></a>
### Nested Schema for `monitoring_settings.advanced_javascript_tag_settings`

Required:

- `additional_event_handlers` (Block List, Min: 1, Max: 1) Additional event handlers and wrappers (see [below for nested schema](#nestedblock--monitoring_settings--advanced_javascript_tag_settings--additional_event_handlers))
- `max_action_name_length` (Number) Maximum character length for action names. Valid values range from 5 to 10000.
- `max_errors_to_capture` (Number) Maximum number of errors to be captured per page. Valid values range from 0 to 50.

Optional:

- `event_wrapper_settings` (Block List, Max: 1) In addition to the event handlers, events called using `addEventListener` or `attachEvent` can be captured. Be careful with this option! Event wrappers can conflict with the JavaScript code on a web page (see [below for nested schema](#nestedblock--monitoring_settings--advanced_javascript_tag_settings--event_wrapper_settings))
- `global_event_capture_settings` (Block List, Max: 1) Global event capture settings (see [below for nested schema](#nestedblock--monitoring_settings--advanced_javascript_tag_settings--global_event_capture_settings))
- `instrument_unsupported_ajax_frameworks` (Boolean) Instrumentation of unsupported Ajax frameworks enabled/disabled
- `proxy_wrapper_enabled` (Boolean) Proxy wrapper enabled/disabled
- `special_characters_to_escape` (String) Additional special characters that are to be escaped using non-alphanumeric characters in HTML escape format. Maximum length 30 character. Allowed characters are `^`, `\`, `<` and `>`.
- `sync_beacon_firefox` (Boolean) Send the beacon signal as a synchronous XMLHttpRequest using Firefox enabled/disabled
- `sync_beacon_internet_explorer` (Boolean) Send the beacon signal as a synchronous XMLHttpRequest using Internet Explorer enabled/disabled
- `user_action_name_attribute` (String) User action name attribute

<a id="nestedblock--monitoring_settings--advanced_javascript_tag_settings--additional_event_handlers"></a>
### Nested Schema for `monitoring_settings.advanced_javascript_tag_settings.additional_event_handlers`

Required:

- `max_dom_nodes` (Number) Max. number of DOM nodes to instrument. Valid values range from 0 to 100000.

Optional:

- `blur` (Boolean) Blur event handler enabled/disabled
- `change` (Boolean) Change event handler enabled/disabled
- `click` (Boolean) Click event handler enabled/disabled
- `mouseup` (Boolean) Mouseup event handler enabled/disabled
- `to_string_method` (Boolean) toString method enabled/disabled
- `use_mouse_up_event_for_clicks` (Boolean) Use mouseup event for clicks enabled/disabled


<a id="nestedblock--monitoring_settings--advanced_javascript_tag_settings--event_wrapper_settings"></a>
### Nested Schema for `monitoring_settings.advanced_javascript_tag_settings.event_wrapper_settings`

Optional:

- `blur` (Boolean) Blur enabled/disabled
- `change` (Boolean) Change enabled/disabled
- `click` (Boolean) Click enabled/disabled
- `mouseup` (Boolean) MouseUp enabled/disabled
- `touch_end` (Boolean) TouchEnd enabled/disabled
- `touch_start` (Boolean) TouchStart enabled/disabled


<a id="nestedblock--monitoring_settings--advanced_javascript_tag_settings--global_event_capture_settings"></a>
### Nested Schema for `monitoring_settings.advanced_javascript_tag_settings.global_event_capture_settings`

Optional:

- `additional_event_captured_as_user_input` (String) Additional events to be captured globally as user input. 

For example `DragStart` or `DragEnd`. Maximum 100 characters.
- `change` (Boolean) Change enabled/disabled
- `click` (Boolean) Click enabled/disabled
- `doubleclick` (Boolean) DoubleClick enabled/disabled
- `keydown` (Boolean) KeyDown enabled/disabled
- `keyup` (Boolean) KeyUp enabled/disabled
- `mousedown` (Boolean) MouseDown enabled/disabled
- `mouseup` (Boolean) MouseUp enabled/disabled
- `scroll` (Boolean) Scroll enabled/disabled
- `touch_end` (Boolean) TouchEnd enabled/disabled
- `touch_start` (Boolean) TouchStart enabled/disabled



<a id="nestedblock--monitoring_settings--content_capture"></a>
### Nested Schema for `monitoring_settings.content_capture`

Required:

- `resource_timing_settings` (Block List, Min: 1, Max: 1) Settings for resource timings capture (see [below for nested schema](#nestedblock--monitoring_settings--content_capture--resource_timing_settings))
- `timeout_settings` (Block List, Min: 1, Max: 1) Settings for timed action capture (see [below for nested schema](#nestedblock--monitoring_settings--content_capture--timeout_settings))

Optional:

- `javascript_errors` (Boolean) JavaScript errors monitoring enabled/disabled
- `visually_complete_and_speed_index` (Boolean) Visually complete and Speed index support enabled/disabled
- `visually_complete_settings` (Block List, Max: 1) Settings for VisuallyComplete (see [below for nested schema](#nestedblock--monitoring_settings--content_capture--visually_complete_settings))

<a id="nestedblock--monitoring_settings--content_capture--resource_timing_settings"></a>
### Nested Schema for `monitoring_settings.content_capture.resource_timing_settings`

Required:

- `instrumentation_delay` (Number) Instrumentation delay for monitoring resource and image resource impact in browsers that don't offer W3C resource timings. 

Valid values range from 0 to 9999.

Only effective if `nonW3cResourceTimings` is enabled

Optional:

- `non_w3c_resource_timings` (Boolean) Timing for JavaScript files and images on non-W3C supported browsers enabled/disabled
- `resource_timing_capture_type` (String) Defines how detailed resource timings are captured.

Only effective if **w3cResourceTimings** or **nonW3cResourceTimings** is enabled. Possible values are `CAPTURE_ALL_SUMMARIES`, `CAPTURE_FULL_DETAILS` and `CAPTURE_LIMITED_SUMMARIES`
- `resource_timings_domain_limit` (Number) Limits the number of domains for which W3C resource timings are captured.

Only effective if **resourceTimingCaptureType** is `CAPTURE_LIMITED_SUMMARIES`. Valid values range from 0 to 50.
- `w3c_resource_timings` (Boolean) W3C resource timings for third party/CDN enabled/disabled


<a id="nestedblock--monitoring_settings--content_capture--timeout_settings"></a>
### Nested Schema for `monitoring_settings.content_capture.timeout_settings`

Required:

- `temporary_action_limit` (Number) Defines how deep temporary actions may cascade. 0 disables temporary actions completely. Recommended value if enabled is 3
- `temporary_action_total_timeout` (Number) The total timeout of all cascaded timeouts that should still be able to create a temporary action

Optional:

- `timed_action_support` (Boolean) Timed action support enabled/disabled. 

Enable to detect actions that trigger sending of XHRs via `setTimout` methods


<a id="nestedblock--monitoring_settings--content_capture--visually_complete_settings"></a>
### Nested Schema for `monitoring_settings.content_capture.visually_complete_settings`

Optional:

- `exclude_url_regex` (String) A RegularExpression used to exclude images and iframes from being detected by the VC module
- `ignored_mutations_list` (String) Query selector for mutation nodes to ignore in VC and SI calculation
- `inactivity_timeout` (Number) The time in ms the VC module waits for no mutations happening on the page after the load action. Defaults to 1000. Valid values range from 0 to 30000.
- `mutation_timeout` (Number) Determines the time in ms VC waits after an action closes to start calculation. Defaults to 50. Valid values range from 0 to 5000.
- `threshold` (Number) Minimum visible area in pixels of elements to be counted towards VC and SI. Defaults to 50. Valid values range from 0 to 10000.



<a id="nestedblock--monitoring_settings--browser_restriction_settings"></a>
### Nested Schema for `monitoring_settings.browser_restriction_settings`

Required:

- `mode` (String) The mode of the list of browser restrictions. Possible values area `EXCLUDE` and `INCLUDE`.

Optional:

- `restrictions` (Block List, Max: 1) A list of browser restrictions (see [below for nested schema](#nestedblock--monitoring_settings--browser_restriction_settings--restrictions))

<a id="nestedblock--monitoring_settings--browser_restriction_settings--restrictions"></a>
### Nested Schema for `monitoring_settings.browser_restriction_settings.restrictions`

Required:

- `restriction` (Block List, Min: 1) Browser exclusion rules for the browsers that are to be excluded (see [below for nested schema](#nestedblock--monitoring_settings--browser_restriction_settings--restrictions--restriction))

<a id="nestedblock--monitoring_settings--browser_restriction_settings--restrictions--restriction"></a>
### Nested Schema for `monitoring_settings.browser_restriction_settings.restrictions.restriction`

Required:

- `browser_type` (String) The type of the browser that is used. Possible values are `ANDROID_WEBKIT`, `BOTS_SPIDERS`, `CHROME`, `EDGE`, `FIREFOX`, `INTERNET_EXPLORER, `OPERA` and `SAFARI`

Optional:

- `browser_version` (String) The version of the browser that is used
- `comparator` (String) No documentation available. Possible values are `EQUALS`, `GREATER_THAN_OR_EQUAL` and `LOWER_THAN_OR_EQUAL`.
- `platform` (String) The platform on which the browser is being used. Possible values are `ALL`, `DESKTOP` and `MOBILE`




<a id="nestedblock--monitoring_settings--ip_address_restriction_settings"></a>
### Nested Schema for `monitoring_settings.ip_address_restriction_settings`

Required:

- `mode` (String) The mode of the list of ip address restrictions. Possible values area `EXCLUDE` and `INCLUDE`.

Optional:

- `restrictions` (Block List, Max: 1) The IP addresses or the IP address ranges to be mapped to the location (see [below for nested schema](#nestedblock--monitoring_settings--ip_address_restriction_settings--restrictions))

<a id="nestedblock--monitoring_settings--ip_address_restriction_settings--restrictions"></a>
### Nested Schema for `monitoring_settings.ip_address_restriction_settings.restrictions`

Required:

- `range` (Block List, Min: 1) The IP address or the IP address range to be mapped to the location (see [below for nested schema](#nestedblock--monitoring_settings--ip_address_restriction_settings--restrictions--range))

<a id="nestedblock--monitoring_settings--ip_address_restriction_settings--restrictions--range"></a>
### Nested Schema for `monitoring_settings.ip_address_restriction_settings.restrictions.range`

Required:

- `address` (String) The IP address to be mapped. 

For an IP address range, this is the **from** address.

Optional:

- `address_to` (String) The **to** address of the IP address range.
- `subnet_mask` (Number) The subnet mask of the IP address range. Valid values range from 0 to 128.




<a id="nestedblock--monitoring_settings--javascript_framework_support"></a>
### Nested Schema for `monitoring_settings.javascript_framework_support`

Optional:

- `active_x_object` (Boolean) ActiveXObject support enabled/disabled
- `angular` (Boolean) AngularJS and Angular support enabled/disabled
- `dojo` (Boolean) Dojo support enabled/disabled
- `extjs` (Boolean) ExtJS, Sencha Touch support enabled/disabled
- `icefaces` (Boolean) ICEfaces support enabled/disabled
- `jquery` (Boolean) jQuery, Backbone.js support enabled/disabled
- `moo_tools` (Boolean) MooTools support enabled/disabled
- `prototype` (Boolean) Prototype support enabled/disabled


<a id="nestedblock--monitoring_settings--javascript_injection_rules"></a>
### Nested Schema for `monitoring_settings.javascript_injection_rules`

Required:

- `rule` (Block List, Min: 1) Java script injection rule (see [below for nested schema](#nestedblock--monitoring_settings--javascript_injection_rules--rule))

<a id="nestedblock--monitoring_settings--javascript_injection_rules--rule"></a>
### Nested Schema for `monitoring_settings.javascript_injection_rules.rule`

Required:

- `rule` (String) The url rule of the java script injection. Possible values are `AFTER_SPECIFIC_HTML`, `AUTOMATIC_INJECTION`, `BEFORE_SPECIFIC_HTML` and `DO_NOT_INJECT`.
- `url_operator` (String) The url operator of the java script injection. Possible values are `ALL_PAGES`, `CONTAINS`, `ENDS_WITH`, `EQUALS` and `STARTS_WITH`.

Optional:

- `enabled` (Boolean) `fetch()` request capture enabled/disabled
- `html_pattern` (String) The HTML pattern of the java script injection
- `target` (String) The target against which the rule of the java script injection should be matched. Possible values are `PAGE_QUERY` and `URL`.
- `url_pattern` (String) The url pattern of the java script injection




<a id="nestedblock--waterfall_settings"></a>
### Nested Schema for `waterfall_settings`

Required:

- `resource_browser_caching_threshold` (Number) Warn about resources with a lower browser cache rate above *X*%. Values between 1 and 100 are allowed.
- `resources_threshold` (Number) Warn about resources larger than *X* bytes. Values between 0 and 99999000 are allowed.
- `slow_cnd_resources_threshold` (Number) Warn about slow CDN resources with a response time above *X* ms. Values between 0 and 99999000 are allowed.
- `slow_first_party_resources_threshold` (Number) Warn about slow 1st party resources with a response time above *X* ms. Values between 0 and 99999000 are allowed.
- `slow_third_party_resources_threshold` (Number) Warn about slow 3rd party resources with a response time above *X* ms. Values between 0 and 99999000 are allowed.
- `speed_index_visually_complete_ratio_threshold` (Number) Warn if Speed index exceeds *X* % of Visually complete. Values between 1 and 99 are allowed.
- `uncompressed_resources_threshold` (Number) Warn about uncompressed resources larger than *X* bytes. Values between 0 and 99999 are allowed.


<a id="nestedblock--xhr_action_apdex_settings"></a>
### Nested Schema for `xhr_action_apdex_settings`

Required:

- `frustrating_fallback_threshold` (Number) Fallback threshold of an XHR action, defining a tolerable user experience, when the configured KPM is not available. Values between 0 and 240000 are allowed.
- `frustrating_threshold` (Number) Maximal value of apdex, which is considered as tolerable user experience. Values between 0 and 240000 are allowed.
- `tolerated_fallback_threshold` (Number) Fallback threshold of an XHR action, defining a satisfied user experience, when the configured KPM is not available. Values between 0 and 60000 are allowed.
- `tolerated_threshold` (Number) Maximal value of apdex, which is considered as satisfied user experience. Values between 0 and 60000 are allowed.

Optional:

- `threshold` (Number, Deprecated) no documentation available


<a id="nestedblock--conversion_goals"></a>
### Nested Schema for `conversion_goals`

Required:

- `goal` (Block List, Min: 1) A conversion goal of the application (see [below for nested schema](#nestedblock--conversion_goals--goal))

<a id="nestedblock--conversion_goals--goal"></a>
### Nested Schema for `conversion_goals.goal`

Required:

- `name` (String) The name of the conversion goal. Valid length within 1 and 50 characters.

Optional:

- `destination` (Block List, Max: 1) Configuration of a destination-based conversion goal (see [below for nested schema](#nestedblock--conversion_goals--goal--destination))
- `id` (String) The ID of conversion goal. 

 Omit it while creating a new conversion goal
- `type` (String) The type of the conversion goal. Possible values are `Destination`, `UserAction`, `VisitDuration` and `VisitNumActions`
- `user_action` (Block List, Max: 1) Configuration of a destination-based conversion goal (see [below for nested schema](#nestedblock--conversion_goals--goal--user_action))
- `visit_duration` (Block List, Max: 1) Configuration of a destination-based conversion goal (see [below for nested schema](#nestedblock--conversion_goals--goal--visit_duration))
- `visit_num_action` (Block List, Max: 1) Configuration of a destination-based conversion goal (see [below for nested schema](#nestedblock--conversion_goals--goal--visit_num_action))

<a id="nestedblock--conversion_goals--goal--destination"></a>
### Nested Schema for `conversion_goals.goal.destination`

Required:

- `url_or_path` (String) The path to be reached to hit the conversion goal

Optional:

- `case_sensitive` (Boolean) The match is case-sensitive (`true`) or (`false`)
- `match_type` (String) The operator of the match. Possible values are `Begins`, `Contains` and `Ends`.


<a id="nestedblock--conversion_goals--goal--user_action"></a>
### Nested Schema for `conversion_goals.goal.user_action`

Optional:

- `action_type` (String) Type of the action to which the rule applies. Possible values are `Custom`, `Load` and `Xhr`.
- `case_sensitive` (Boolean) The match is case-sensitive (`true`) or (`false`)
- `match_entity` (String) The type of the entity to which the rule applies. Possible values are `ActionName`, `CssSelector`, `JavaScriptVariable`, `MetaTag`, `PagePath`, `PageTitle`, `PageUrl`, `UrlAnchor` and `XhrUrl`.
- `match_type` (String) The operator of the match. Possible values are `Begins`, `Contains` and `Ends`.
- `value` (String) The value to be matched to hit the conversion goal


<a id="nestedblock--conversion_goals--goal--visit_duration"></a>
### Nested Schema for `conversion_goals.goal.visit_duration`

Required:

- `duration` (Number) The duration of session to hit the conversion goal, in milliseconds


<a id="nestedblock--conversion_goals--goal--visit_num_action"></a>
### Nested Schema for `conversion_goals.goal.visit_num_action`

Optional:

- `num_user_actions` (Number) The number of user actions to hit the conversion goal




<a id="nestedblock--key_user_actions"></a>
### Nested Schema for `key_user_actions`

Required:

- `action` (Block Set, Min: 1) Configuration of the key user action (see [below for nested schema](#nestedblock--key_user_actions--action))

<a id="nestedblock--key_user_actions--action"></a>
### Nested Schema for `key_user_actions.action`

Required:

- `name` (String) The name of the action
- `type` (String) The type of the action. Possible values are `Custom`, `Load` and `Xhr`.

Optional:

- `domain` (String) The domain where the action is performed.



<a id="nestedblock--meta_data_capture_settings"></a>
### Nested Schema for `meta_data_capture_settings`

Optional:

- `capture` (Block List) Java script agent meta data capture settings (see [below for nested schema](#nestedblock--meta_data_capture_settings--capture))

<a id="nestedblock--meta_data_capture_settings--capture"></a>
### Nested Schema for `meta_data_capture_settings.capture`

Required:

- `capturing_name` (String) The name of the meta data to capture
- `name` (String) Name for displaying the captured values in Dynatrace
- `type` (String) The type of the meta data to capture. Possible values are `COOKIE`, `CSS_SELECTOR`, `JAVA_SCRIPT_FUNCTION`, `JAVA_SCRIPT_VARIABLE`, `META_TAG` and `QUERY_STRING`.

Optional:

- `public_metadata` (Boolean) `true` if this metadata should be captured regardless of the privacy settings, `false` otherwise
- `unique_id` (Number) The unique ID of the meta data to capture
- `use_last_value` (Boolean) `true` if the last captured value should be used for this metadata. By default the first value will be used.



<a id="nestedblock--session_replay_config"></a>
### Nested Schema for `session_replay_config`

Required:

- `cost_control_percentage` (Number) (Field has overlap with `dynatrace_web_app_enablement`) Session replay sampling rating in percent

Optional:

- `css_resource_capturing_exclusion_rules` (List of String) (Field has overlap with `dynatrace_session_replay_resource_capture`) A list of URLs to be excluded from CSS resource capturing
- `enable_css_resource_capturing` (Boolean) (Field has overlap with `dynatrace_session_replay_resource_capture`) Capture (`true`) or don't capture (`false`) CSS resources from the session
- `enabled` (Boolean) (Field has overlap with `dynatrace_web_app_enablement`) SessionReplay Enabled/Disabled


<a id="nestedblock--user_action_and_session_properties"></a>
### Nested Schema for `user_action_and_session_properties`

Optional:

- `property` (Block List) User action and session properties settings (see [below for nested schema](#nestedblock--user_action_and_session_properties--property))

<a id="nestedblock--user_action_and_session_properties--property"></a>
### Nested Schema for `user_action_and_session_properties.property`

Required:

- `id` (Number) Unique id among all userTags and properties of this application
- `key` (String) Key of the property
- `origin` (String) The origin of the property. Possible values are `JAVASCRIPT_API`, `META_DATA` and `SERVER_SIDE_REQUEST_ATTRIBUTE`.
- `type` (String) The data type of the property. Possible values are `DATE`, `DOUBLE`, `LONG`, `LONG_STRING` and `STRING`.

Optional:

- `aggregation` (String) The aggregation type of the property. 

  It defines how multiple values of the property are aggregated. Possible values are `AVERAGE`, `FIRST`, `LAST`, `MAXIMUM`, `MINIMUM` and `SUM`.
- `cleanup_rule` (String) The cleanup rule of the property. 

Defines how to extract the data you need from a string value. Specify the [regular expression](https://dt-url.net/k9e0iaq) for the data you need there
- `display_name` (String) The display name of the property
- `ignore_case` (Boolean) If `true`, the value of this property will always be stored in lower case. Defaults to `false`.
- `long_string_length` (Number) If the `type` is `LONG_STRING`, the max length for this property. Must be a multiple of `100`. Defaults to `200`. Maximum is `1000`.
- `metadata_id` (Number) If the origin is `META_DATA`, metaData id of the property
- `server_side_request_attribute` (String) The ID of the request attribute. 

Only applicable when the **origin** is set to `SERVER_SIDE_REQUEST_ATTRIBUTE`
- `store_as_session_property` (Boolean) If `true`, the property is stored as a session property
- `store_as_user_action_property` (Boolean) If `true`, the property is stored as a user action property



<a id="nestedblock--user_action_naming_settings"></a>
### Nested Schema for `user_action_naming_settings`

Optional:

- `custom_action_naming_rules` (Block List, Max: 1) User action naming rules for custom actions (see [below for nested schema](#nestedblock--user_action_naming_settings--custom_action_naming_rules))
- `ignore_case` (Boolean) Case insensitive naming
- `load_action_naming_rules` (Block List, Max: 1) User action naming rules for loading actions (see [below for nested schema](#nestedblock--user_action_naming_settings--load_action_naming_rules))
- `placeholders` (Block List, Max: 1) User action placeholders (see [below for nested schema](#nestedblock--user_action_naming_settings--placeholders))
- `query_parameter_cleanups` (Set of String) User action naming rules for custom actions. If not specified Dynatrace assumes `__sid`, `cfid`, `cftoken`, `phpsessid` and `sid`.
- `split_user_actions_by_domain` (Boolean) Deactivate this setting if different domains should not result in separate user actions
- `use_first_detected_load_action` (Boolean) First load action found under an XHR action should be used when true. Else the deepest one under the xhr action is used
- `xhr_action_naming_rules` (Block List, Max: 1) User action naming rules for XHR actions (see [below for nested schema](#nestedblock--user_action_naming_settings--xhr_action_naming_rules))

<a id="nestedblock--user_action_naming_settings--custom_action_naming_rules"></a>
### Nested Schema for `user_action_naming_settings.custom_action_naming_rules`

Required:

- `rule` (Block List, Min: 1) The settings of naming rule (see [below for nested schema](#nestedblock--user_action_naming_settings--custom_action_naming_rules--rule))

<a id="nestedblock--user_action_naming_settings--custom_action_naming_rules--rule"></a>
### Nested Schema for `user_action_naming_settings.custom_action_naming_rules.rule`

Required:

- `template` (String) Naming pattern. Use Curly brackets `{}` to select placeholders

Optional:

- `conditions` (Block List, Max: 1) Defines the conditions when the naming rule should apply (see [below for nested schema](#nestedblock--user_action_naming_settings--custom_action_naming_rules--rule--conditions))
- `use_or_conditions` (Boolean) If set to `true` the conditions will be connected by logical OR instead of logical AND

<a id="nestedblock--user_action_naming_settings--custom_action_naming_rules--rule--conditions"></a>
### Nested Schema for `user_action_naming_settings.custom_action_naming_rules.rule.conditions`

Required:

- `condition` (Block List, Min: 1) Defines the conditions when the naming rule should apply (see [below for nested schema](#nestedblock--user_action_naming_settings--custom_action_naming_rules--rule--conditions--condition))

<a id="nestedblock--user_action_naming_settings--custom_action_naming_rules--rule--conditions--condition"></a>
### Nested Schema for `user_action_naming_settings.custom_action_naming_rules.rule.conditions.condition`

Required:

- `operand1` (String) Must be a defined placeholder wrapped in curly braces
- `operator` (String) The operator of the condition. Possible values are `CONTAINS`, `ENDS_WITH`, `EQUALS`, `IS_EMPTY`, `IS_NOT_EMPTY`, `MATCHES_REGULAR_EXPRESSION`, `NOT_CONTAINS`, `NOT_ENDS_WITH`, `NOT_EQUALS`, `NOT_MATCHES_REGULAR_EXPRESSION`, `NOT_STARTS_WITH` and `STARTS_WITH`.

Optional:

- `operand2` (String) Must be null if operator is `IS_EMPTY`, a regex if operator is `MATCHES_REGULAR_ERPRESSION`. In all other cases the value can be a freetext or a placeholder wrapped in curly braces





<a id="nestedblock--user_action_naming_settings--load_action_naming_rules"></a>
### Nested Schema for `user_action_naming_settings.load_action_naming_rules`

Required:

- `rule` (Block List, Min: 1) The settings of naming rule (see [below for nested schema](#nestedblock--user_action_naming_settings--load_action_naming_rules--rule))

<a id="nestedblock--user_action_naming_settings--load_action_naming_rules--rule"></a>
### Nested Schema for `user_action_naming_settings.load_action_naming_rules.rule`

Required:

- `template` (String) Naming pattern. Use Curly brackets `{}` to select placeholders

Optional:

- `conditions` (Block List, Max: 1) Defines the conditions when the naming rule should apply (see [below for nested schema](#nestedblock--user_action_naming_settings--load_action_naming_rules--rule--conditions))
- `use_or_conditions` (Boolean) If set to `true` the conditions will be connected by logical OR instead of logical AND

<a id="nestedblock--user_action_naming_settings--load_action_naming_rules--rule--conditions"></a>
### Nested Schema for `user_action_naming_settings.load_action_naming_rules.rule.conditions`

Required:

- `condition` (Block List, Min: 1) Defines the conditions when the naming rule should apply (see [below for nested schema](#nestedblock--user_action_naming_settings--load_action_naming_rules--rule--conditions--condition))

<a id="nestedblock--user_action_naming_settings--load_action_naming_rules--rule--conditions--condition"></a>
### Nested Schema for `user_action_naming_settings.load_action_naming_rules.rule.conditions.condition`

Required:

- `operand1` (String) Must be a defined placeholder wrapped in curly braces
- `operator` (String) The operator of the condition. Possible values are `CONTAINS`, `ENDS_WITH`, `EQUALS`, `IS_EMPTY`, `IS_NOT_EMPTY`, `MATCHES_REGULAR_EXPRESSION`, `NOT_CONTAINS`, `NOT_ENDS_WITH`, `NOT_EQUALS`, `NOT_MATCHES_REGULAR_EXPRESSION`, `NOT_STARTS_WITH` and `STARTS_WITH`.

Optional:

- `operand2` (String) Must be null if operator is `IS_EMPTY`, a regex if operator is `MATCHES_REGULAR_ERPRESSION`. In all other cases the value can be a freetext or a placeholder wrapped in curly braces





<a id="nestedblock--user_action_naming_settings--placeholders"></a>
### Nested Schema for `user_action_naming_settings.placeholders`

Required:

- `placeholder` (Block List, Min: 1) User action placeholders (see [below for nested schema](#nestedblock--user_action_naming_settings--placeholders--placeholder))

<a id="nestedblock--user_action_naming_settings--placeholders--placeholder"></a>
### Nested Schema for `user_action_naming_settings.placeholders.placeholder`

Required:

- `input` (String) The input for the place holder. Possible values are `ELEMENT_IDENTIFIER`, `INPUT_TYPE`, `METADATA`, `PAGE_TITLE`, `PAGE_URL`, `SOURCE_URL`, `TOP_XHR_URL` and `XHR_URL`
- `name` (String) Placeholder name. Valid length needs to be between 1 and 50 characters
- `processing_part` (String) The part to process. Possible values are `ALL`, `ANCHOR` and `PATH`

Optional:

- `metadata_id` (Number) The ID of the metadata
- `processing_steps` (Block List, Max: 1) The processing step settings (see [below for nested schema](#nestedblock--user_action_naming_settings--placeholders--placeholder--processing_steps))
- `use_guessed_element_identifier` (Boolean) Use the element identifier that was selected by Dynatrace

<a id="nestedblock--user_action_naming_settings--placeholders--placeholder--processing_steps"></a>
### Nested Schema for `user_action_naming_settings.placeholders.placeholder.processing_steps`

Required:

- `step` (Block List, Min: 1) The processing step (see [below for nested schema](#nestedblock--user_action_naming_settings--placeholders--placeholder--processing_steps--step))

<a id="nestedblock--user_action_naming_settings--placeholders--placeholder--processing_steps--step"></a>
### Nested Schema for `user_action_naming_settings.placeholders.placeholder.processing_steps.step`

Required:

- `type` (String) An action to be taken by the processing: 

* `SUBSTRING`: Extracts the string between `patternBefore` and `patternAfter`. 
* `REPLACEMENT`: Replaces the string between `patternBefore` and `patternAfter` with the specified `replacement`.
* `REPLACE_WITH_PATTERN`: Replaces the **patternToReplace** with the specified **replacement**. 
* `EXTRACT_BY_REGULAR_EXPRESSION`: Extracts the part of the string that matches the **regularExpression**. 
* `REPLACE_WITH_REGULAR_EXPRESSION`: Replaces all occurrences that match **regularExpression** with the specified **replacement**. 
* `REPLACE_IDS`: Replaces all IDs and UUIDs with the specified **replacement**. Possible values are `EXTRACT_BY_REGULAR_EXPRESSION`, `REPLACEMENT`, `REPLACE_IDS`, `REPLACE_WITH_PATTERN`, `REPLACE_WITH_REGULAR_EXPRESSION` and `SUBSTRING`.

Optional:

- `fallback_to_input` (Boolean) If set to `true`: Returns the input if `patternBefore` or `patternAfter` cannot be found and the `type` is `SUBSTRING`. Returns the input if `regularExpression` doesn't match and `type` is `EXTRACT_BY_REGULAR_EXPRESSION`. 

 Otherwise `null` is returned.
- `pattern_after` (String) The pattern after the required value. It will be removed.
- `pattern_after_search_type` (String) The required occurrence of **patternAfter**. Possible values are `FIRST` and `LAST`.
- `pattern_before` (String) The pattern before the required value. It will be removed.
- `pattern_before_search_type` (String) The required occurrence of **patternBefore**. Possible values are `FIRST` and `LAST`.
- `pattern_to_replace` (String) The pattern to be replaced. 

 Only applicable if the `type` is `REPLACE_WITH_PATTERN`.
- `regular_expression` (String) A regular expression for the string to be extracted or replaced. Only applicable if the `type` is `EXTRACT_BY_REGULAR_EXPRESSION` or `REPLACE_WITH_REGULAR_EXPRESSION`.
- `replacement` (String) Replacement for the original value





<a id="nestedblock--user_action_naming_settings--xhr_action_naming_rules"></a>
### Nested Schema for `user_action_naming_settings.xhr_action_naming_rules`

Required:

- `rule` (Block List, Min: 1) The settings of naming rule (see [below for nested schema](#nestedblock--user_action_naming_settings--xhr_action_naming_rules--rule))

<a id="nestedblock--user_action_naming_settings--xhr_action_naming_rules--rule"></a>
### Nested Schema for `user_action_naming_settings.xhr_action_naming_rules.rule`

Required:

- `template` (String) Naming pattern. Use Curly brackets `{}` to select placeholders

Optional:

- `conditions` (Block List, Max: 1) Defines the conditions when the naming rule should apply (see [below for nested schema](#nestedblock--user_action_naming_settings--xhr_action_naming_rules--rule--conditions))
- `use_or_conditions` (Boolean) If set to `true` the conditions will be connected by logical OR instead of logical AND

<a id="nestedblock--user_action_naming_settings--xhr_action_naming_rules--rule--conditions"></a>
### Nested Schema for `user_action_naming_settings.xhr_action_naming_rules.rule.conditions`

Required:

- `condition` (Block List, Min: 1) Defines the conditions when the naming rule should apply (see [below for nested schema](#nestedblock--user_action_naming_settings--xhr_action_naming_rules--rule--conditions--condition))

<a id="nestedblock--user_action_naming_settings--xhr_action_naming_rules--rule--conditions--condition"></a>
### Nested Schema for `user_action_naming_settings.xhr_action_naming_rules.rule.conditions.condition`

Required:

- `operand1` (String) Must be a defined placeholder wrapped in curly braces
- `operator` (String) The operator of the condition. Possible values are `CONTAINS`, `ENDS_WITH`, `EQUALS`, `IS_EMPTY`, `IS_NOT_EMPTY`, `MATCHES_REGULAR_EXPRESSION`, `NOT_CONTAINS`, `NOT_ENDS_WITH`, `NOT_EQUALS`, `NOT_MATCHES_REGULAR_EXPRESSION`, `NOT_STARTS_WITH` and `STARTS_WITH`.

Optional:

- `operand2` (String) Must be null if operator is `IS_EMPTY`, a regex if operator is `MATCHES_REGULAR_ERPRESSION`. In all other cases the value can be a freetext or a placeholder wrapped in curly braces






<a id="nestedblock--user_tags"></a>
### Nested Schema for `user_tags`

Required:

- `tag` (Block List, Min: 1) User tag settings (see [below for nested schema](#nestedblock--user_tags--tag))

<a id="nestedblock--user_tags--tag"></a>
### Nested Schema for `user_tags.tag`

Optional:

- `cleanup_rule` (String) Cleanup rule expression of the userTag
- `id` (Number) A unique ID among all userTags and properties of this application. Minimum value is 1. Do not set that attribute anymore - terraform will handle it. Kept for backwards compatibility
- `ignore_case` (Boolean) If `true`, the value of this tag will always be stored in lower case. Defaults to `false`.
- `metadata_id` (Number) If it's of type metaData, metaData id of the userTag
- `server_side_request_attribute` (String) The ID of the RrequestAttribute for the userTag

Read-Only:

- `unique_id` (Number) A unique ID among all userTags and properties of this application. Minimum value is 1.
 