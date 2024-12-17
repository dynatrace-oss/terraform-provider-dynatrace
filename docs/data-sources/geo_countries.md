---
layout: ""
page_title: "dynatrace_geo_countries Data Source - terraform-provider-dynatrace"
subcategory: "Real User Monitoring"
description: |-
  The data source `dynatrace_geo_countries` covers queries for countries and their codes
---

# dynatrace_geo_countries (Data Source)

The `dynatrace_geo_countries` data source retrieves the list of countries and their codes.

Geographic regions API: GET countries - https://docs.dynatrace.com/docs/shortlink/api-v2-rum-geographic-regions-get-countries

## Example Usage

```terraform
data "dynatrace_geo_countries" "Example" {
}

output "Test" {
  value = data.dynatrace_geo_countries.Example
}

```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `countries` (List of Object) (see [below for nested schema](#nestedatt--countries))
- `id` (String) The ID of this resource.

<a id="nestedatt--countries"></a>
### Nested Schema for `countries`

Read-Only:

- `code` (String)
- `name` (String)