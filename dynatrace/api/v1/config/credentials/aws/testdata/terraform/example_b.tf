resource "dynatrace_aws_credentials" "#name#" {
  label          = "#name#"
  partition_type = "AWS_CN"
  tagged_only    = true

  authentication_data {
    account_id = "246186168471"
    iam_role   = "Dynatrace_monitoring_role_demo1"
  }


  tags_to_monitor {
    name  = "string"
    value = "string2"
  }
  supporting_services_to_monitor {
    name = "KinesisFirehose"
    monitored_metrics {
      name       = "IncomingBytes"
      dimensions = ["DeliveryStreamName"]
      statistic  = "SUM"
    }
    monitored_metrics {
      name       = "IncomingRecords"
      dimensions = ["DeliveryStreamName"]
      statistic  = "SUM"
    }
  }
  supporting_services_to_monitor {
    name = "connect"
    monitored_metrics {
      name       = "CallsBreachingConcurrencyQuota"
      dimensions = ["InstanceId"]
      statistic  = "SUM"
    }
  }
  supporting_services_to_monitor {
    name = "logs"
    monitored_metrics {
      name       = "IncomingLogEvents"
      dimensions = ["LogGroupName"]
      statistic  = "SUM"
    }
    monitored_metrics {
      name       = "IncomingBytes"
      dimensions = ["LogGroupName"]
      statistic  = "SUM"
    }
    monitored_metrics {
      name       = "IncomingBytes"
      dimensions = ["Region"]
      statistic  = "SUM"
    }
    monitored_metrics {
      name       = "IncomingLogEvents"
      dimensions = ["Region"]
      statistic  = "SUM"
    }
  }
  supporting_services_to_monitor {
    name = "polly"
    monitored_metrics {
      name       = "4XXCount"
      dimensions = ["Region", "Operation"]
      statistic  = "AVERAGE"
    }
    monitored_metrics {
      name       = "4XXCount"
      dimensions = ["Region", "Operation"]
      statistic  = "SUM"
    }
    monitored_metrics {
      name       = "2XXCount"
      dimensions = ["Region", "Operation"]
      statistic  = "AVERAGE"
    }
    monitored_metrics {
      name       = "2XXCount"
      dimensions = ["Region", "Operation"]
      statistic  = "SUM"
    }
    monitored_metrics {
      name       = "5XXCount"
      dimensions = ["Region", "Operation"]
      statistic  = "AVERAGE"
    }
    monitored_metrics {
      name       = "5XXCount"
      dimensions = ["Region", "Operation"]
      statistic  = "SUM"
    }
    monitored_metrics {
      name       = "RequestCharacters"
      dimensions = ["Region", "Operation"]
      statistic  = "AVG_MIN_MAX"
    }
    monitored_metrics {
      name       = "RequestCharacters"
      dimensions = ["Region", "Operation"]
      statistic  = "SUM"
    }
    monitored_metrics {
      name       = "ResponseLatency"
      dimensions = ["Region", "Operation"]
      statistic  = "AVG_MIN_MAX"
    }
    monitored_metrics {
      name       = "ResponseLatency"
      dimensions = ["Region", "Operation"]
      statistic  = "SAMPLE_COUNT"
    }
  }
}
