resource "dynatrace_credentials_v2" "token" {
  name              = "#name#"
  token             = "token"
  owner_access_only = false
  scopes            = ["SYNTHETIC"]
}

resource "dynatrace_credentials_v2" "username_password" {
  name              = "#name#"
  username          = "username"
  password          = "password"
  owner_access_only = false
  scopes            = ["SYNTHETIC"]
}

data "dynatrace_credentials_v2" "by_name_and_type" {
  name       = "#name#"
  type       = "USERNAME_PASSWORD"
  depends_on = [dynatrace_credentials_v2.token, dynatrace_credentials_v2.username_password]
}

output "matched_ids" {
  value = jsonencode(sort(data.dynatrace_credentials_v2.by_name_and_type.items[*].id))
}

output "expected_ids" {
  value = jsonencode([dynatrace_credentials_v2.username_password.id])
}
