resource "dynatrace_gcp_principal" "principal" {}

output "records" {
  value = dynatrace_gcp_principal.principal.principal
}
