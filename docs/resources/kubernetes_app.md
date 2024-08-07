---
layout: ""
page_title: dynatrace_kubernetes_app Resource - terraform-provider-dynatrace"
subcategory: "Cloud Platforms"
description: |-
  The resource `dynatrace_kubernetes_app` covers configuration to enable the new Kubernetes app
---

# dynatrace_kubernetes_app (Resource)

-> This resource requires the API token scopes **Read settings** (`settings.read`) and **Write settings** (`settings.write`)

## Dynatrace Documentation

- Kubernetes app - https://docs.dynatrace.com/docs/platform-modules/infrastructure-monitoring/container-platform-monitoring/kubernetes-app

- Settings API - https://www.dynatrace.com/support/help/dynatrace-api/environment-api/settings (schemaId: `builtin:app-transition.kubernetes`)

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_kubernetes_app` downloads all existing Kubernetes app settings

The full documentation of the export feature is available [here](https://dt-url.net/h203qmc).

## Resource Example Usage

```terraform
resource "dynatrace_kubernetes_app" "#name#" {
  scope = "environment"
  kubernetes_app_options {
    enable_kubernetes_app = true
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `kubernetes_app_options` (Block List, Min: 1, Max: 1) no documentation available (see [below for nested schema](#nestedblock--kubernetes_app_options))

### Optional

- `scope` (String) The scope of this setting (KUBERNETES_CLUSTER). Omit this property if you want to cover the whole environment.

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--kubernetes_app_options"></a>
### Nested Schema for `kubernetes_app_options`

Required:

- `enable_kubernetes_app` (Boolean) New Kubernetes experience
 