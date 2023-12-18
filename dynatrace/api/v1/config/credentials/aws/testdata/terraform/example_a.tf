resource "dynatrace_aws_credentials" "#name#" {
  label          = "#name#"
  partition_type = "AWS_DEFAULT"
  tagged_only    = false
  authentication_data {
    account_id = "246186168471"
    iam_role   = "Dynatrace_monitoring_role_demo1"
  }
}
