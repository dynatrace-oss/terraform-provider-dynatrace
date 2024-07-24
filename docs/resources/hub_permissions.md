---
layout: ""
page_title: dynatrace_hub_permissions Resource - terraform-provider-dynatrace"
subcategory: "Environment Settings"
description: |-
  The resource `dynatrace_hub_permissions` covers configuration for hub permissions
---

# dynatrace_hub_permissions (Resource)

-> This resource requires the API token scopes **Read settings** (`settings.read`) and **Write settings** (`settings.write`)

This resource allows configuring email recipients for Dynatrace Hub app installation requests.

## Dynatrace Documentation

- Dynatrace Hub - https://docs.dynatrace.com/docs/manage/hub

- Settings API - https://www.dynatrace.com/support/help/dynatrace-api/environment-api/settings (schemaId: `app:dynatrace.hub:manage.permissions`)

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_hub_permissions` downloads all existing hub permissions

The full documentation of the export feature is available [here](https://dt-url.net/h203qmc).

## Resource Example Usage

```terraform
resource "dynatrace_hub_permissions" "#name#" {
  email = "terraform@dynatrace.com"
  description = "This is an example description"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `description` (String) Description
- `email` (String) Email

### Read-Only

- `id` (String) The ID of this resource.
 