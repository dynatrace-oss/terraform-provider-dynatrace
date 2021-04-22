resource "dynatrace_aws_credentials" "#name#" {
  label          = "#name#"
  partition_type = "AWS_CN"
  tagged_only    = true

  authentication_data {
    access_key = "########0"
    secret_key = "########1"
  }

  tags_to_monitor {
    name  = "string"
    value = "string"
  }
}
