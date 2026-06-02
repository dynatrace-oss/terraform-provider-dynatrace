data "dynatrace_gcp_principal" "principal" {
}

output "records" {
  value = data.dynatrace_gcp_principal.principal.principal
}
