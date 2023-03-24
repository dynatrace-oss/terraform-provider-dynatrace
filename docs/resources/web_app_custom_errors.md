---
layout: ""
page_title: "dynatrace_web_app_custom_errors Resource - terraform-provider-dynatrace"
description: |-
  The resource `dynatrace_web_app_custom_errors` covers configuration for web application custom errors
---

# dynatrace_web_app_custom_errors (Resource)

-> **Settings 2.0** Certain field(s) of this resource has overlap with `dynatrace_web_app_custom_errors`, therefore it is excluded from the default export. To retrieve this resource via export, directly specify it as a command line argument. 

## Dynatrace Documentation

- Configure custom errors - https://www.dynatrace.com/support/help/platform-modules/digital-experience/web-applications/additional-configuration/configure-errors#configure-custom-errors

- Settings API - https://www.dynatrace.com/support/help/dynatrace-api/environment-api/settings (schemaId: `builtin:rum.web.custom-errors`)

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_web_app_custom_errors` downloads all existing custom error configuration

The full documentation of the export feature is available [here](https://registry.terraform.io/providers/dynatrace-oss/dynatrace/latest/docs/guides/export-v2).

## Resource Example Usage

```terraform
resource "dynatrace_web_app_custom_errors" "#name#" {
  ignore_custom_errors_in_apdex_calculation = true
  scope                                     = "APPLICATION-1234567890000000"
  error_rules {
    error_rule {
      key_matcher   = "EQUALS"
      key_pattern   = "hashicorp"
      value_matcher = "ALL"
      capture_settings {
        capture         = true
        consider_for_ai = true
        impact_apdex    = false
      }
    }
    error_rule {
      key_matcher   = "CONTAINS"
      key_pattern   = "TF"
      value_matcher = "ENDS_WITH"
      value_pattern = "EX"
      capture_settings {
        capture = false
      }
    }
    error_rule {
      key_matcher   = "BEGINS_WITH"
      key_pattern   = "terraform"
      value_matcher = "CONTAINS"
      value_pattern = "example"
      capture_settings {
        capture         = true
        consider_for_ai = true
        impact_apdex    = true
      }
    }
    error_rule {
      key_matcher   = "ALL"
      value_matcher = "ALL"
      capture_settings {
        capture         = true
        consider_for_ai = false
        impact_apdex    = true
      }
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `ignore_custom_errors_in_apdex_calculation` (Boolean) (Field has overlap with `dynatrace_application_error_rules`) This setting overrides Apdex settings for individual rules listed below
- `scope` (String) The scope of this setting (APPLICATION)

### Optional

- `error_rules` (Block List, Max: 1) (Field has overlap with `dynatrace_application_error_rules`) (see [below for nested schema](#nestedblock--error_rules))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--error_rules"></a>
### Nested Schema for `error_rules`

Required:

- `error_rule` (Block List, Min: 1) (see [below for nested schema](#nestedblock--error_rules--error_rule))

<a id="nestedblock--error_rules--error_rule"></a>
### Nested Schema for `error_rules.error_rule`

Required:

- `capture_settings` (Block List, Min: 1, Max: 1) Capture settings (see [below for nested schema](#nestedblock--error_rules--error_rule--capture_settings))
- `key_matcher` (String) Possible Values: `ALL`, `BEGINS_WITH`, `CONTAINS`, `ENDS_WITH`, `EQUALS`
- `value_matcher` (String) Possible Values: `ALL`, `BEGINS_WITH`, `CONTAINS`, `ENDS_WITH`, `EQUALS`

Optional:

- `key_pattern` (String) A case-insensitive key pattern
- `value_pattern` (String) A case-insensitive value pattern

<a id="nestedblock--error_rules--error_rule--capture_settings"></a>
### Nested Schema for `error_rules.error_rule.capture_settings`

Required:

- `capture` (Boolean) Capture this error

Optional:

- `consider_for_ai` (Boolean) [View more details](https://dt-url.net/hd580p2k)
- `impact_apdex` (Boolean) Include error in Apdex calculations
 