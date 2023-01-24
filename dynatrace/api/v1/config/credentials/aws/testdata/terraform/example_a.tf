resource "dynatrace_aws_credentials" "#name#" {
  label          = "#name#"
  partition_type = "AWS_CN"
  tagged_only    = true

  authentication_data {
    access_key = "########02"
    secret_key = "########12"
  }

  tags_to_monitor {
    name  = "string"
    value = "string"
  }
}
