resource "dynatrace_automation_workflow_aws_connections" "#name#" {
  name = "#name#"
  type = "webIdentity"
  web_identity {
    role_arn    = "arn:aws:iam::123456789012:role/helloworld"
  }
}