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
