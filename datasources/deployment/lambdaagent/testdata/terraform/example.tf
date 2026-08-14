data "dynatrace_lambda_agent_version" "example" {
}

output "latest" {
  value = data.dynatrace_lambda_agent_version.example.collector
}
