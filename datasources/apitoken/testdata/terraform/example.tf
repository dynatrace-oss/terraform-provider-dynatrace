locals {
  token_name = "#name#"
}

data "dynatrace_api_token" "example" {
  name = local.token_name
}

data "dynatrace_api_tokens" "example" {
}

output "auth_token" {
  value = data.dynatrace_api_token.example.creation_date
}

output "filtered_api_token" {
  value = [for token in data.dynatrace_api_tokens.example.api_tokens : token if token.name == local.token_name][0].name
}
