resource "dynatrace_credentials_v2" "credential" {
  name              = "#name#"
  token             = "token"
  owner_access_only = false
  scopes            = ["SYNTHETIC"]
}

data "dynatrace_credentials_v2" "by_name" {
  name = upper(dynatrace_credentials_v2.credential.name)
}

output "matched_ids" {
  value = jsonencode(sort(data.dynatrace_credentials_v2.by_name.items[*].id))
}

output "expected_ids" {
  value = jsonencode([dynatrace_credentials_v2.credential.id])
}
