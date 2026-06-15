resource "dynatrace_gcp_principal" "principal" {}

output "principal_output" {
  value = dynatrace_gcp_principal.principal.principal
}
