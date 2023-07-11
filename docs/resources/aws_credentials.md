---
layout: ""
page_title: dynatrace_aws_credentials Resource - terraform-provider-dynatrace"
subcategory: "Credentials"
description: |-
  The resource `dynatrace_aws_credentials` covers configuration for AWS credentials
---

# dynatrace_aws_credentials (Resource)

-> This resource requires the API token scopes **Read configuration** (`ReadConfig`) and **Write configuration** (`WriteConfig`)

## Dynatrace Documentation

- Set up Dynatrace on Amazon Web Services - https://www.dynatrace.com/support/help/setup-and-configuration/setup-on-cloud-platforms/amazon-web-services

- AWS credentials API - https://www.dynatrace.com/support/help/dynatrace-api/configuration-api/aws-credentials-api

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_aws_credentials` downloads all existing AWS credentials configuration

The full documentation of the export feature is available [here](https://registry.terraform.io/providers/dynatrace-oss/dynatrace/latest/docs/guides/export-v2).

## Resource Example Usage

```terraform
resource "dynatrace_aws_credentials" "#name#" {
  label          = "#name#"
  partition_type = "AWS_DEFAULT"
  tagged_only    = false
  authentication_data {
    account_id = "246186168471"
    iam_role   = "Dynatrace_monitoring_role_demo1"
  }
  supporting_services_to_monitor {
    name = "ECS"
    monitored_metrics {
      name       = "CPUReservation"
      dimensions = ["ClusterName"]
      statistic  = "AVG_MIN_MAX"
    }
    monitored_metrics {
      name       = "CPUUtilization"
      dimensions = ["ClusterName", "ServiceName"]
      statistic  = "AVG_MIN_MAX"
    }
    monitored_metrics {
      name       = "MemoryReservation"
      dimensions = ["ClusterName"]
      statistic  = "AVG_MIN_MAX"
    }
    monitored_metrics {
      name       = "MemoryUtilization"
      dimensions = ["ClusterName", "ServiceName"]
      statistic  = "AVG_MIN_MAX"
    }
  }
  supporting_services_to_monitor {
    name = "ecscontainerinsights"
    monitored_metrics {
      name       = "CpuUtilized"
      dimensions = ["ClusterName"]
      statistic  = "AVG_MIN_MAX"
    }
    monitored_metrics {
      name       = "MemoryUtilized"
      dimensions = ["ClusterName"]
      statistic  = "AVG_MIN_MAX"
    }
    monitored_metrics {
      name       = "NetworkRxBytes"
      dimensions = ["ClusterName"]
      statistic  = "AVG_MIN_MAX"
    }
    monitored_metrics {
      name       = "RunningTaskCount"
      dimensions = ["ClusterName", "ServiceName"]
      statistic  = "AVG_MIN_MAX"
    }
    monitored_metrics {
      name       = "MemoryUtilized"
      dimensions = ["ClusterName", "ServiceName"]
      statistic  = "AVG_MIN_MAX"
    }
    monitored_metrics {
      name       = "StorageReadBytes"
      dimensions = ["ClusterName"]
      statistic  = "AVG_MIN_MAX"
    }
    monitored_metrics {
      name       = "StorageReadBytes"
      dimensions = ["ClusterName", "ServiceName"]
      statistic  = "AVG_MIN_MAX"
    }
    monitored_metrics {
      name       = "CpuUtilized"
      dimensions = ["ClusterName", "ServiceName"]
      statistic  = "AVG_MIN_MAX"
    }
    monitored_metrics {
      name       = "StorageWriteBytes"
      dimensions = ["ClusterName"]
      statistic  = "AVG_MIN_MAX"
    }
    monitored_metrics {
      name       = "StorageWriteBytes"
      dimensions = ["ClusterName", "ServiceName"]
      statistic  = "AVG_MIN_MAX"
    }
    monitored_metrics {
      name       = "TaskCount"
      dimensions = ["ClusterName"]
      statistic  = "AVG_MIN_MAX"
    }
    monitored_metrics {
      name       = "NetworkTxBytes"
      dimensions = ["ClusterName", "ServiceName"]
      statistic  = "AVG_MIN_MAX"
    }
    monitored_metrics {
      name       = "NetworkTxBytes"
      dimensions = ["ClusterName"]
      statistic  = "AVG_MIN_MAX"
    }
    monitored_metrics {
      name       = "NetworkRxBytes"
      dimensions = ["ClusterName", "ServiceName"]
      statistic  = "AVG_MIN_MAX"
    }
    monitored_metrics {
      name       = "instance_memory_utilization"
      dimensions = ["ClusterName"]
      statistic  = "AVG_MIN_MAX"
    }
    monitored_metrics {
      name       = "instance_number_of_running_tasks"
      dimensions = ["ClusterName"]
      statistic  = "AVG_MIN_MAX"
    }
    monitored_metrics {
      name       = "instance_network_total_bytes"
      dimensions = ["ClusterName"]
      statistic  = "AVG_MIN_MAX"
    }
    monitored_metrics {
      name       = "instance_cpu_utilization"
      dimensions = ["ClusterName"]
      statistic  = "AVG_MIN_MAX"
    }
    monitored_metrics {
      name       = "instance_filesystem_utilization"
      dimensions = ["ClusterName"]
      statistic  = "AVG_MIN_MAX"
    }
    monitored_metrics {
      name       = "CpuUtilized"
      dimensions = ["ClusterName", "TaskDefinitionFamily"]
      statistic  = "AVG_MIN_MAX"
    }
    monitored_metrics {
      name       = "MemoryUtilized"
      dimensions = ["ClusterName", "TaskDefinitionFamily"]
      statistic  = "AVG_MIN_MAX"
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `authentication_data` (Block List, Min: 1, Max: 1) credentials for the AWS authentication (see [below for nested schema](#nestedblock--authentication_data))
- `partition_type` (String) The type of the AWS partition
- `tagged_only` (Boolean) Monitor only resources which have specified AWS tags (`true`) or all resources (`false`)

### Optional

- `label` (String) The name of the credentials
- `supporting_services_managed_in_dynatrace` (Boolean) If enabled (`true`) the attribute `supporting_services` will not get synchronized with Dynatrace. You will be able to manage them via WebUI without interference by Terraform.
- `supporting_services_to_monitor` (Block List) supporting services to be monitored (see [below for nested schema](#nestedblock--supporting_services_to_monitor))
- `tags_to_monitor` (Block List, Max: 10) AWS tags to be monitored. You can specify up to 10 tags. Only applicable when the **tagged_only** parameter is set to `true` (see [below for nested schema](#nestedblock--tags_to_monitor))
- `unknowns` (String) Any attributes that aren't yet supported by this provider

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--authentication_data"></a>
### Nested Schema for `authentication_data`

Optional:

- `access_key` (String) the access key
- `account_id` (String) the ID of the Amazon account
- `iam_role` (String) the IAM role to be used by Dynatrace to get monitoring data
- `secret_key` (String, Sensitive) the secret access key
- `unknowns` (String) Any attributes that aren't yet supported by this provider

Read-Only:

- `external_id` (String) (Read only) the external ID token for setting an IAM role. You can obtain it with the `GET /aws/iamExternalId` request


<a id="nestedblock--supporting_services_to_monitor"></a>
### Nested Schema for `supporting_services_to_monitor`

Optional:

- `monitored_metrics` (Block Set) a list of metrics to be monitored for this service (see [below for nested schema](#nestedblock--supporting_services_to_monitor--monitored_metrics))
- `name` (String) the name of the supporting service
- `unknowns` (String) Any attributes that aren't yet supported by this provider

<a id="nestedblock--supporting_services_to_monitor--monitored_metrics"></a>
### Nested Schema for `supporting_services_to_monitor.monitored_metrics`

Optional:

- `dimensions` (List of String) a list of metric's dimensions names
- `name` (String) the name of the metric of the supporting service
- `statistic` (String) the statistic (aggregation) to be used for the metric. AVG_MIN_MAX value is 3 statistics at once: AVERAGE, MINIMUM and MAXIMUM
- `unknowns` (String) Any attributes that aren't yet supported by this provider



<a id="nestedblock--tags_to_monitor"></a>
### Nested Schema for `tags_to_monitor`

Optional:

- `name` (String) the key of the AWS tag.
- `unknowns` (String) Any attributes that aren't yet supported by this provider
- `value` (String) the value of the AWS tag
 