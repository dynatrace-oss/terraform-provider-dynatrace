resource "dynatrace_aws_connection" "test-aws-connection" {
  name = "Test connection Web identity"
  web_identity {
    consumers = ["APP:dynatrace.aws.connector"]
  }
}

resource "dynatrace_aws_connection_role_arn" "test-aws-connection-arn" {
  aws_connection_id = dynatrace_aws_connection.test-aws-connection.id
  role_arn          = "<role-arn>"
}