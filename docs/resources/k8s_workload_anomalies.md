---
layout: ""
page_title: dynatrace_k8s_workload_anomalies Resource - terraform-provider-dynatrace"
description: |-
  The resource `dynatrace_k8s_workload_anomalies` covers configuration for Kubernetes workload anomalies
---

# dynatrace_k8s_workload_anomalies (Resource)

## Dynatrace Documentation

- Alert on common Kubernetes/OpenShift issues - https://www.dynatrace.com/support/help/platform-modules/infrastructure-monitoring/container-platform-monitoring/kubernetes-monitoring/alert-on-kubernetes-issues

- Settings API - https://www.dynatrace.com/support/help/dynatrace-api/environment-api/settings (schemaId: `builtin:anomaly-detection.kubernetes.workload`)

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_k8s_workload_anomalies` downloads all existing Kubernetes workload anomaly configuration

The full documentation of the export feature is available [here](https://registry.terraform.io/providers/dynatrace-oss/dynatrace/latest/docs/guides/export-v2).

## Resource Example Usage

```terraform
resource "dynatrace_k8s_workload_anomalies" "#name#" {
  scope = "environment"
  container_restarts {
    enabled = true
    configuration {
      observation_period_in_minutes = 6
      sample_period_in_minutes      = 4
      threshold                     = 2
    }
  }
  deployment_stuck {
    enabled = true
    configuration {
      observation_period_in_minutes = 5
      sample_period_in_minutes      = 4
    }
  }
  not_all_pods_ready {
    enabled = true
    configuration {
      observation_period_in_minutes = 6
      sample_period_in_minutes      = 4
    }
  }
  pending_pods {
    enabled = true
    configuration {
      observation_period_in_minutes = 16
      sample_period_in_minutes      = 11
      threshold                     = 2
    }
  }
  pod_stuck_in_terminating {
    enabled = true
    configuration {
      observation_period_in_minutes = 6
      sample_period_in_minutes      = 4
    }
  }
  workload_without_ready_pods {
    enabled = true
    configuration {
      observation_period_in_minutes = 6
      sample_period_in_minutes      = 4
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `container_restarts` (Block List, Min: 1, Max: 1) no documentation available (see [below for nested schema](#nestedblock--container_restarts))
- `deployment_stuck` (Block List, Min: 1, Max: 1) no documentation available (see [below for nested schema](#nestedblock--deployment_stuck))
- `not_all_pods_ready` (Block List, Min: 1, Max: 1) no documentation available (see [below for nested schema](#nestedblock--not_all_pods_ready))
- `pending_pods` (Block List, Min: 1, Max: 1) no documentation available (see [below for nested schema](#nestedblock--pending_pods))
- `pod_stuck_in_terminating` (Block List, Min: 1, Max: 1) no documentation available (see [below for nested schema](#nestedblock--pod_stuck_in_terminating))
- `workload_without_ready_pods` (Block List, Min: 1, Max: 1) no documentation available (see [below for nested schema](#nestedblock--workload_without_ready_pods))

### Optional

- `scope` (String) The scope of this setting (CLOUD_APPLICATION_NAMESPACE, KUBERNETES_CLUSTER). Omit this property if you want to cover the whole environment.

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--container_restarts"></a>
### Nested Schema for `container_restarts`

Required:

- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)

Optional:

- `configuration` (Block List, Max: 1) Alert if (see [below for nested schema](#nestedblock--container_restarts--configuration))

<a id="nestedblock--container_restarts--configuration"></a>
### Nested Schema for `container_restarts.configuration`

Required:

- `observation_period_in_minutes` (Number) within the last
- `sample_period_in_minutes` (Number) per minute, for any
- `threshold` (Number) there is at least



<a id="nestedblock--deployment_stuck"></a>
### Nested Schema for `deployment_stuck`

Required:

- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)

Optional:

- `configuration` (Block List, Max: 1) Alert if (see [below for nested schema](#nestedblock--deployment_stuck--configuration))

<a id="nestedblock--deployment_stuck--configuration"></a>
### Nested Schema for `deployment_stuck.configuration`

Required:

- `observation_period_in_minutes` (Number) within the last
- `sample_period_in_minutes` (Number) workload stops progressing for at least



<a id="nestedblock--not_all_pods_ready"></a>
### Nested Schema for `not_all_pods_ready`

Required:

- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)

Optional:

- `configuration` (Block List, Max: 1) Alert if (see [below for nested schema](#nestedblock--not_all_pods_ready--configuration))

<a id="nestedblock--not_all_pods_ready--configuration"></a>
### Nested Schema for `not_all_pods_ready.configuration`

Required:

- `observation_period_in_minutes` (Number) within the last
- `sample_period_in_minutes` (Number) some workload pods are not ready for at least



<a id="nestedblock--pending_pods"></a>
### Nested Schema for `pending_pods`

Required:

- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)

Optional:

- `configuration` (Block List, Max: 1) Alert if (see [below for nested schema](#nestedblock--pending_pods--configuration))

<a id="nestedblock--pending_pods--configuration"></a>
### Nested Schema for `pending_pods.configuration`

Required:

- `observation_period_in_minutes` (Number) within the last
- `sample_period_in_minutes` (Number) stuck in pending state for at least
- `threshold` (Number) there is at least



<a id="nestedblock--pod_stuck_in_terminating"></a>
### Nested Schema for `pod_stuck_in_terminating`

Required:

- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)

Optional:

- `configuration` (Block List, Max: 1) Alert if (see [below for nested schema](#nestedblock--pod_stuck_in_terminating--configuration))

<a id="nestedblock--pod_stuck_in_terminating--configuration"></a>
### Nested Schema for `pod_stuck_in_terminating.configuration`

Required:

- `observation_period_in_minutes` (Number) within the last
- `sample_period_in_minutes` (Number) pod termination stops progressing for at least



<a id="nestedblock--workload_without_ready_pods"></a>
### Nested Schema for `workload_without_ready_pods`

Required:

- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)

Optional:

- `configuration` (Block List, Max: 1) Alert if (see [below for nested schema](#nestedblock--workload_without_ready_pods--configuration))

<a id="nestedblock--workload_without_ready_pods--configuration"></a>
### Nested Schema for `workload_without_ready_pods.configuration`

Required:

- `observation_period_in_minutes` (Number) within the last
- `sample_period_in_minutes` (Number) workload has no ready pods for at least
 