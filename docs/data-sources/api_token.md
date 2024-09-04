---
layout: ""
page_title: "dynatrace_api_token Data Source - terraform-provider-dynatrace"
subcategory: "Access Tokens"
description: |-
  The data source `dynatrace_api_token` covers queries for an access token
---

# dynatrace_api_token (Data Source)

The API token data source allows a single access token to be retrieved by its name, note the token value is not included in the response.

If multiple tokens match the given name, the first result will be retrieved. To retrieve multiple tokens of the same name, please utilize the `dynatrace_api_tokens` data source.

## Example Usage

```terraform
data "dynatrace_api_token" "example" {
  name = "Terraform"
}

output "example" {
  value = data.dynatrace_api_token.example
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String)

### Read-Only

- `creation_date` (String) Token creation date in ISO 8601 format (yyyy-MM-dd'T'HH:mm:ss.SSS'Z')
- `enabled` (Boolean) The token is enabled (true) or disabled (false), default disabled (false).
- `expiration_date` (String) The expiration date of the token.
- `id` (String) The ID of this resource.
- `owner` (String) The owner of the token
- `personal_access_token` (Boolean) The token is a personal access token (true) or an API token (false).
- `scopes` (Set of String) A list of the scopes to be assigned to the token.