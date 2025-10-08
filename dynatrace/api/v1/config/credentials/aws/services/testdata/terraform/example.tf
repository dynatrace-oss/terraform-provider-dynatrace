resource "dynatrace_aws_credentials" "Example" {
  label                                 = "#name#"
  partition_type                        = "AWS_DEFAULT"
  tagged_only                           = false
  authentication_data {
    account_id = "123456789"
    iam_role   = "aws-monitoring-role"
  }
}

resource "dynatrace_aws_service" "ElastiCache" {
  name           = "ElastiCache"
  credentials_id = dynatrace_aws_credentials.Example.id
  metric {
    name       = "NetworkBandwidthOutAllowanceExceeded"
    dimensions = [ "CacheClusterId" ]
    statistic  = "SUM"
  }
  metric {
    name       = "CPUUtilization"
    dimensions = [ "CacheClusterId" ]
    statistic  = "AVG_MIN_MAX"
  }
}