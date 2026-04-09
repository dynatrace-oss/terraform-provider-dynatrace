resource "dynatrace_aws_credentials" "cred" {
  label          = "#name#"
  partition_type = "AWS_DEFAULT"
  tagged_only    = false
  authentication_data {
    account_id = "246186168471"
    iam_role   = "Dynatrace_monitoring_role_demo1"
  }
  tags_to_monitor {
    name = "Environment"
    value = "Production"
  }
  tags_to_monitor {
    name = "Environment"
    value = "Hardening"
  }
}
