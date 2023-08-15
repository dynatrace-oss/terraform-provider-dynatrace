# Dynatrace Integration Configuration: AWS CloudWatch Logs Terraform Example

AWS log forwarding allows you to stream logs from Amazon CloudWatch into Dynatrace logs via an ActiveGate, additional information can be found under [CloudWatch Logs](https://www.dynatrace.com/support/help/setup-and-configuration/setup-on-cloud-platforms/amazon-web-services/amazon-web-services-integrations/aws-log-forwarder).

The configuration below is an example of how to configure the CloudWatch log forwarder with Terraform.

```
locals {
  stack_name = "dynatrace-aws-logs"
  # tags are optional
  tags = {
    source = "terraform"
  }

  # To discover potential names for log groups the command
  #   aws logs describe-log-groups --output text --query "logGroups[].[logGroupName]"
  # will be helpful
  log_groups = [
    "/aws/lambda/my-lambda",
    "/aws/apigateway/my-api"
  ]

  # https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html
  # empty string matches everything
  filter_pattern = ""
}

resource "local_file" "dynatrace_aws_log_forwarder" {
  filename = "${path.module}/dynatrace-aws-log-forwarder.zip"
  source   = "https://github.com/dynatrace-oss/dynatrace-aws-log-forwarder/releases/latest/download/dynatrace-aws-log-forwarder.zip"
}

resource "archive_file" "dynatrace_aws_log_forwarder" {
  type        = "zip"
  # You may want to choose a different folder to extract to here
  output_path = "${path.module}"
  source_file = file(local_file.dynatrace_aws_log_forwarder.filename)
}

# ----------------------------------------------------------------------------------------------------------------------
# Lambda Function code needs to be available via S3 Bucket
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_s3_bucket" "dynatrace_aws_log_forwarder" {
  bucket = "dynatrace_aws_log_forwarder"
}

resource "aws_s3_bucket_object" "dynatrace_aws_log_forwarder" {
  bucket = aws_s3_bucket.dynatrace_aws_log_forwarder.id
  key    = "dynatrace-aws-log-forwarder.zip"
  source = "${archive_file.dynatrace_aws_log_forwarder.output_path}/dynatrace-aws-log-forwarder-lambda.zip"
}

# ----------------------------------------------------------------------------------------------------------------------
# Create Stack
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_cloudformation_stack" "dynatrace_aws_log_forwarder" {
  name = local.stack_name
  capabilities = ["CAPABILITY_IAM"]
  tags = local.tags
  parameters = {
    # The URL to Your Dynatrace SaaS logs ingest target.
    # If you choose new ActiveGate deployment (UseExistingActiveGate = false), provide Tenant URL (https://<your_environment_ID>.live.dynatrace.com).
    # If you choose to use existing ActiveGate (UseExistingActiveGate = true), provide ActiveGate endpoint:
    #  - for Public ActiveGate: https://<your_environment_ID>.live.dynatrace.com
    #  - for Environment ActiveGate: https://<active_gate_address>:9999/e/<environment_id> (e.g. https://22.111.98.222:9999/e/abc12345)
    DynatraceEnvironmentUrl = var.target_url
    # Set to "" if you choose to new an existing ActiveGate (UseExistingActiveGate = false). It will not be used.
    # If you choose to use an existing ActiveGate (UseExistingActiveGate = true)
    #  - for Public ActiveGate: https://<TenantId>.live.dynatrace.com
    #  - for Environment ActiveGate: https://<active_gate_address>:9999/e/<TenantId> (e.g. https://22.111.98.222:9999/e/abc12345)
    TenantId = "tenant-id"
    # Dynatrace API token. Integration requires API v1 Log import Token permission.
    DynatraceApiKey = var.target_api_token
    # Enables checking SSL certificate of the target Active Gate
    # By default (if this option is not provided) certificates aren't validated
    VerifySSLTargetActiveGate = false
    # If you choose new ActiveGate deployment, set to 'false'
    # In such case, a new EC2 with ActiveGate will be added to log forwarder deployment (enclosed in VPC with log forwarder)
    # If you choose to use existing ActiveGate (either Public AG or Environment AG), set to 'true'.
    UseExistingActiveGate = true
    # Only needed when UseExistingActiveGate = false
    # PaaS token generated in Integration/Platform as a Service.
    # Used for ActiveGate installation.
    DynatracePaasToken = var.target_paas_token
    # Defines the max log length after which a log will be truncated
    # For values over 8192 there's also a change in Dynatrace settings needed
    # For that you need to contact Dynatrace One.
    MaxLogContentLength = 8192
  }

  ###############################################################
  # This is VERY messy - definitely room for improvement
  # When executing `dynatrace-aws-logs.sh deploy` it executes two things in a row:
  #   * `aws cloudformation deploy` creates a stack using the template file `dynatrace-aws-log-forwarder-template.yaml`
  #       - The Lambda Function defined within `dynatrace-aws-log-forwarder-template.yaml` contains a source code only a skeleton
  #   * `aws lambda update-function-code` immediately afterwards updates the source code contained within `dynatrace-aws-log-forwarder-lambda.zip`
  # The AWS Terraform Provider doesn't offer a counterpart for `aws lambda update-function-code`
  #   - The resource `aws_lambda_function` would CREATE a Lambda Function for us
  #   - But in that case we'd have to replicate the settings of it from `dynatrace-aws-log-forwarder-template.yaml`
  # The current solution replaces within the template the skeleton code 
  #     ZipFile: |
  #        def handler(event, context):
  #          raise Exception("Dynatrace Logs Lambda has not been uploaded")
  # with
  #     S3Bucket: ${aws_s3_bucket.dynatrace_aws_log_forwarder.id}
  #     S3Key: ${aws_s3_bucket_object.dynatrace_aws_log_forwarder.key}
  # in other words, we're referring to the zip archive that has been previously made available via S3 Bucket
  # That way there is no counterpart for `aws lambda update-function-code` necessary
  template_body = replace(file("dynatrace-aws-log-forwarder-template.yaml"), "        ZipFile: |\n          def handler(event, context):\n            raise Exception(\"Dynatrace Logs Lambda has not been uploaded\")", "        S3Bucket: ${aws_s3_bucket.dynatrace_aws_log_forwarder.id}\n        S3Key: ${aws_s3_bucket_object.dynatrace_aws_log_forwarder.key}")
}

resource "aws_cloudwatch_log_subscription_filter" "example" {
  for_each = local.log_groups

  name            = local.stack_name
  role_arn        = aws_cloudformation_stack.dynatrace_aws_log_forwarder.outputs["CloudWatchLogsRoleArn"]
  log_group_name  = each.key
  filter_pattern  = local.filter_pattern
  destination_arn = aws_cloudformation_stack.dynatrace_aws_log_forwarder.outputs["FirehoseArn"]
}

# aws logs put-subscription-filter --log-group-name "$LOG_GROUP" --filter-name "$SUBSCRIPTION_FILTER_NAME" --filter-pattern "$FILTER_PATTERN"  --role-arn "$ROLE_ARN"
```