data "dynatrace_credentials_v2" "nothing" {
  name = "#name#"
}

output "matched_ids" {
  value = jsonencode(sort(data.dynatrace_credentials_v2.nothing.items[*].id))
}

output "expected_ids" {
  value = jsonencode([])
}
