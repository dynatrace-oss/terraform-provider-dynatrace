resource "dynatrace_credentials_v2" "first" {
  name              = "#name#"
  token             = "token"
  owner_access_only = false
  scopes            = ["SYNTHETIC"]
}

resource "dynatrace_credentials_v2" "second" {
  name              = "#name#"
  token             = "token"
  owner_access_only = false
  scopes            = ["SYNTHETIC"]
}

data "dynatrace_credentials_v2" "by_name" {
  name       = "#name#"
  depends_on = [dynatrace_credentials_v2.first, dynatrace_credentials_v2.second]
}

output "matched_ids" {
  value = jsonencode(sort(data.dynatrace_credentials_v2.by_name.items[*].id))
}

output "expected_ids" {
  value = jsonencode(sort([dynatrace_credentials_v2.first.id, dynatrace_credentials_v2.second.id]))
}
