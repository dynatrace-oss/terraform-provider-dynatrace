---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "dynatrace_application_error_rules Resource - terraform-provider-dynatrace"
subcategory: ""
description: |-
  
---

# dynatrace_application_error_rules (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- **custom_errors** (Block List, Max: 1) An ordered list of HTTP errors.

 Rules are evaluated from top to bottom; the first matching rule applies (see [below for nested schema](#nestedblock--custom_errors))
- **http_errors** (Block List, Max: 1) An ordered list of HTTP errors.

 Rules are evaluated from top to bottom; the first matching rule applies (see [below for nested schema](#nestedblock--http_errors))
- **id** (String) The ID of this resource.
- **ignore_custom_errors_apdex** (Boolean) Exclude (`true`) or include (`false`) custom errors listed in **customErrorRules** in Apdex calculation
- **ignore_http_errors_apdex** (Boolean) Exclude (`true`) or include (`false`) HTTP errors listed in **httpErrorRules** in Apdex calculation
- **ignore_js_errors_apdex** (Boolean) Exclude (`true`) or include (`false`) JavaScript errors in Apdex calculation
- **web_application_id** (String) The EntityID of the the WebApplication

<a id="nestedblock--custom_errors"></a>
### Nested Schema for `custom_errors`

Required:

- **rule** (Block List, Min: 1) Configuration of the custom error in the web application (see [below for nested schema](#nestedblock--custom_errors--rule))

<a id="nestedblock--custom_errors--rule"></a>
### Nested Schema for `custom_errors.rule`

Optional:

- **capture** (Boolean) Capture (`true`) or ignore (`false`) the error
- **custom_alerting** (Boolean) Include (`true`) or exclude (`false`) the error in Davis AI [problem detection and analysis](https://dt-url.net/a963kd2)
- **impact_apdex** (Boolean) Include (`true`) or exclude (`false`) the error in Apdex calculation
- **key_matcher** (String) The matching operation for the **keyPattern**. Possible values are `BEGINS_WITH`, `CONTAINS`, `ENDS_WITH` and `EQUALS`
- **key_pattern** (String) The key of the error to look for
- **value_matcher** (String) The matching operation for the **valuePattern**. Possible values are `BEGINS_WITH`, `CONTAINS`, `ENDS_WITH` and `EQUALS`.
- **value_pattern** (String) The value of the error to look for



<a id="nestedblock--http_errors"></a>
### Nested Schema for `http_errors`

Required:

- **rule** (Block List, Min: 1) Configuration of the HTTP error in the web application (see [below for nested schema](#nestedblock--http_errors--rule))

<a id="nestedblock--http_errors--rule"></a>
### Nested Schema for `http_errors.rule`

Optional:

- **capture** (Boolean) Capture (`true`) or ignore (`false`) the error
- **consider_blocked_requests** (Boolean) If `true`, match by errors that have CSP Rule violations
- **consider_for_ai** (Boolean) Include (`true`) or exclude (`false`) the error in Davis AI [problem detection and analysis](https://dt-url.net/a963kd2)
- **consider_unknown_error_code** (Boolean) If `true`, match by errors that have unknown HTTP status code
- **error_codes** (String) The HTTP status code or status code range to match by. 

This field is required if **considerUnknownErrorCode** AND **considerBlockedRequests** are both set to `false`
- **filter** (String) The matching rule for the URL. Popssible values are `BEGINS_WITH`, `CONTAINS`, `ENDS_WITH` and `EQUALS`.
- **filter_by_url** (Boolean) If `true`, filter errors by URL
- **impact_apdex** (Boolean) Include (`true`) or exclude (`false`) the error in Apdex calculation
- **url** (String) The URL to look for

