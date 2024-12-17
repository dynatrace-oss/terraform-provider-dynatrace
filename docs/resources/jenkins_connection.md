---
layout: ""
page_title: dynatrace_jenkins_connection Resource - terraform-provider-dynatrace"
subcategory: "Connections"
description: |-
  The resource `dynatrace_jenkins_connection` covers configuration for Jenkins connections
---

# dynatrace_jenkins_connection (Resource)

-> This resource requires the API token scopes **Read settings** (`settings.read`) and **Write settings** (`settings.write`)

## Dynatrace Documentation

- Jenkins - https://docs.dynatrace.com/docs/analyze-explore-automate/workflows/actions/jenkins

- Settings API - https://www.dynatrace.com/support/help/dynatrace-api/environment-api/settings (schemaId: `app:dynatrace.jenkins.connector:connection`)

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_jenkins_connection` downloads all existing Jenkins connections

The full documentation of the export feature is available [here](https://dt-url.net/h203qmc).

## Resource Example Usage

```terraform
resource "dynatrace_jenkins_connection" "#name#"{
  name    = "#name#"
  url     = "https://www.#name#.com"
  username    = "#name#"
  password   = "#######"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the Jenkins connection
- `password` (String, Sensitive) The password of the user or API token obtained from the Jenkins UI (Dashboard > User > Configure > API Token)
- `url` (String) Base URL of your Jenkins instance (e.g. https://[YOUR_JENKINS_DOMAIN]/)
- `username` (String) The name of your Jenkins user (e.g. jenkins)

### Read-Only

- `id` (String) The ID of this resource.