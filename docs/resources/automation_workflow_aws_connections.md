---
layout: ""
page_title: dynatrace_automation_workflow_aws_connections Resource - terraform-provider-dynatrace"
subcategory: "Automation"
description: |-
  The resource `dynatrace_automation_workflow_aws_connections` covers configuration for AWS connections for Workflows app
---

# dynatrace_automation_workflow_aws_connections (Resource)

-> This resource requires the API token scopes **Read settings** (`settings.read`) and **Write settings** (`settings.write`)

## Dynatrace Documentation

- AWS for Workflows - https://docs.dynatrace.com/docs/platform-modules/automations/workflows/actions/aws

- Settings API - https://www.dynatrace.com/support/help/dynatrace-api/environment-api/settings (schemaId: `builtin:hyperscaler-authentication.aws.connection`)

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_automation_workflow_aws_connections` downloads existing AWS connections for Workflows configuration

The full documentation of the export feature is available [here](https://dt-url.net/h203qmc).

## Resource Example Usage

```terraform
resource "dynatrace_automation_workflow_aws_connections" "#name#" {
  name = "#name#"
  type = "webIdentity"
  web_identity {
    role_arn    = "arn:aws:iam::helloworld:role.helloworld"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name
- `type` (String) Possible Values: `WebIdentity`

### Optional

- `web_identity` (Block List, Max: 1) no documentation available (see [below for nested schema](#nestedblock--web_identity))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--web_identity"></a>
### Nested Schema for `web_identity`

Required:

- `role_arn` (String, Sensitive) The ARN of the AWS role that should be assumed

Optional:

- `policy_arns` (List of String, Sensitive) An optional list of policies that can be used to restrict the AWS role
 