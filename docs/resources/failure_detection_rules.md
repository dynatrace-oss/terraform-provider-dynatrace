---
layout: ""
page_title: dynatrace_failure_detection_rules Resource - terraform-provider-dynatrace"
subcategory: "Failure Detection"
description: |-
  The resource `dynatrace_failure_detection_rules` covers configuration for failure detection rules
---

# dynatrace_failure_detection_rules (Resource)

-> This resource requires the API token scopes **Read settings** (`settings.read`) and **Write settings** (`settings.write`)

## Dynatrace Documentation

- Configure service failure detection - https://www.dynatrace.com/support/help/platform-modules/applications-and-microservices/services/service-monitoring-settings/configure-service-failure-detection

- Settings API - https://www.dynatrace.com/support/help/dynatrace-api/environment-api/settings (schemaId: `builtin:failure-detection.environment.rules`)

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_failure_detection_rules` downloads all existing failure detection rules

The full documentation of the export feature is available [here](https://dt-url.net/h203qmc).

## Resource Example Usage

```terraform
resource "dynatrace_failure_detection_rules" "#name#" {
  name         = ""
  enabled      = true
  parameter_id = "${dynatrace_failure_detection_parameters.parameter.id}"
  conditions {
    condition {
      attribute = "SERVICE_NAME"
      predicate {
        case_sensitive = true
        predicate_type = "STRING_EQUALS"
        text_values    = [ "Terraform" ]
      }
    }
  }
}

resource "dynatrace_failure_detection_parameters" "parameter" {
  name        = "#name#"
  description = "Created by Terraform"
  broken_links {
    http_404_not_found_failures = false
  }
  exception_rules {
    ignore_all_exceptions         = false
    ignore_span_failure_detection = true
    custom_error_rules {
      custom_error_rule {
        request_attribute = "195b205c-5c01-4563-b29b-e33caf24ec7d"
        condition {
          compare_operation_type = "STRING_EXISTS"
        }
      }
    }
    custom_handled_exceptions {
      custom_handled_exception {
        class_pattern   = "ClassPattern"
        message_pattern = "ExceptionPattern"
      }
    }
    ignored_exceptions {
      custom_handled_exception {
        class_pattern   = "ClassPattern"
        message_pattern = "ExceptionPattern"
      }
    }
    success_forcing_exceptions {
      custom_handled_exception {
        class_pattern   = "ClassPattern"
        message_pattern = "ExceptionPattern"
      }
    }
  }
  http_response_codes {
    client_side_errors                        = "400-599"
    fail_on_missing_response_code_client_side = false
    fail_on_missing_response_code_server_side = true
    server_side_errors                        = "500-599"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `conditions` (Block List, Min: 1, Max: 1) Conditions (see [below for nested schema](#nestedblock--conditions))
- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)
- `name` (String) Rule name
- `parameter_id` (String) Failure detection parameters

### Optional

- `description` (String) Rule description
- `insert_after` (String) Because this resource allows for ordering you may specify the ID of the resource instance that comes before this instance regarding order. If not specified when creating the setting will be added to the end of the list. If not specified during update the order will remain untouched

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--conditions"></a>
### Nested Schema for `conditions`

Required:

- `condition` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--conditions--condition))

<a id="nestedblock--conditions--condition"></a>
### Nested Schema for `conditions.condition`

Required:

- `attribute` (String) Possible Values: `PG_NAME`, `PG_TAG`, `SERVICE_MANAGEMENT_ZONE`, `SERVICE_NAME`, `SERVICE_TAG`, `SERVICE_TYPE`
- `predicate` (Block List, Min: 1, Max: 1) Condition to check the attribute against (see [below for nested schema](#nestedblock--conditions--condition--predicate))

<a id="nestedblock--conditions--condition--predicate"></a>
### Nested Schema for `conditions.condition.predicate`

Required:

- `predicate_type` (String) Predicate type

Optional:

- `case_sensitive` (Boolean) Case sensitive
- `management_zones` (Set of String) Management zones
- `service_type` (Set of String) Service types
- `tag_keys` (Set of String) Tag keys
- `tags` (Set of String) Tags (exact match)
- `text_values` (Set of String) Names
 