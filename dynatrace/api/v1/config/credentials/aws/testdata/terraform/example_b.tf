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
}
