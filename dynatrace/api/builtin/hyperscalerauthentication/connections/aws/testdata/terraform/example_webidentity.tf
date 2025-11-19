resource "dynatrace_aws_connection" "test-aws-connection" {
  name = "#name#"
  web_identity {
    consumers = ["APP:dynatrace.aws.connector"]
  }
}

resource "aws_iam_openid_connect_provider" "dynatrace-oidc-provider" {
  url = "https://token.dynatrace.com"

  client_id_list = [
    "<TENANT_URL>/app-id/dynatrace.aws.connector",
  ]
}

resource "aws_iam_role" "example_role" {
  name = "#name#"
  assume_role_policy = jsonencode(
    {
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Principal": {
                    "Federated": "${aws_iam_openid_connect_provider.dynatrace-oidc-provider.arn}"
                },
                "Action": "sts:AssumeRoleWithWebIdentity",
                "Condition": {
                    "StringEquals": {
                        "${aws_iam_openid_connect_provider.dynatrace-oidc-provider.url}:sub": "dt:connection-id/${dynatrace_aws_connection.test-aws-connection.id}",
                        "${aws_iam_openid_connect_provider.dynatrace-oidc-provider.url}:aud": "<TENANT_URL>/app-id/dynatrace.aws.connector"
                    }
                }
            }
        ]
    })
}

resource "dynatrace_aws_connection_role_arn" "test-aws-connection-arn" {
  aws_connection_id = dynatrace_aws_connection.test-aws-connection.id
  role_arn          = aws_iam_role.example_role.arn

  timeouts {
    create = "15s"
  }
}
