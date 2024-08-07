---
layout: ""
page_title: "dynatrace_platform_bucket Resource - terraform-provider-dynatrace"
subcategory: "Platform"
description: |-
  The resource `dynatrace_platform_bucket` covers configuration of Grail Buckets
---

# dynatrace_platform_bucket (Resource)


## Dynatrace Documentation

- Grail Buckets - https://www.dynatrace.com/support/help/platform/grail/data-model#custom-grail-buckets

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_platform_bucket` downloads all existing bucket definitions

The full documentation of the export feature is available [here](https://dt-url.net/h203qmc).

## Prerequisites

Using this resource requires an OAuth client to be configured within your account settings.
The scopes of the OAuth Client need to include `View bucket metadata (storage:bucket-definitions:read)`, `Write buckets (storage:bucket-definitions:write)` and `Delete buckets (storage:bucket-definitions:delete)`.

Finally the provider configuration requires the credentials for that OAuth Client.
The configuration section of your provider needs to look like this.
```terraform
provider "dynatrace" {
  dt_env_url   = "https://########.live.dynatrace.com/"  
  dt_api_token = "######.########################.################################################################"  

  # Usually not required. Terraform will deduct it if `dt_env_url` has been specified
  # automation_env_url = "https://########.apps.dynatrace.com/" 
  automation_client_id = "######.########"
  automation_client_secret = "######.########.################################################################"  
}
```
-> In order to handle credentials in a secure manner we recommend to use the environment variables `DYNATRACE_AUTOMATION_CLIENT_ID` and `DYNATRACE_AUTOMATION_CLIENT_SECRET` as an alternative.

## Resource Example Usage

```terraform
resource "dynatrace_platform_bucket" "#name#" {
  name         = "#name#"
  display_name = "Custom logs bucket playground"
  retention    = 67
  table        = "logs"
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name / id of the bucket definition
- `retention` (Number) The retention of stored data in days
- `table` (String) The table the bucket definition applies to. Possible values are `logs`, `spans`,	`events` and `bizevents`. Changing this attribute will result in deleting and re-creating the bucket definition

### Optional

- `display_name` (String) The name of the bucket definition when visualized within the UI

### Read-Only

- `id` (String) The ID of this resource.
- `status` (String) The status of the bucket definition. Usually has the value `active` unless an update or delete is currently happening
