output "aws_role_arn" {
  description = "ARN of the IAM role Dynatrace assumes."
  value       = aws_iam_role.monitoring.arn
}

output "dynatrace_connection_id" {
  description = "Object id of the Dynatrace AWS connection. Used as sts:ExternalId."
  value       = dynatrace_aws_connection.this.id
}

output "dynatrace_monitoring_configuration_id" {
  description = "Object id of the Dynatrace AWS monitoring configuration."
  value       = dynatrace_aws_monitoring_configuration.this.id
}
