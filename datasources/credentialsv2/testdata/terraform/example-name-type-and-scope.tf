resource "dynatrace_credentials_v2" "token_synthetic" {
  name              = "#name#"
  token             = "token"
  owner_access_only = false
  scopes            = ["SYNTHETIC"]
}

resource "dynatrace_credentials_v2" "token_app_engine" {
  name              = "#name#"
  token             = "token"
  owner_access_only = false
  scopes            = ["APP_ENGINE"]
}

resource "dynatrace_credentials_v2" "username_password_synthetic" {
  name              = "#name#"
  username          = "username"
  password          = "password"
  owner_access_only = false
  scopes            = ["SYNTHETIC"]
}

data "dynatrace_credentials_v2" "by_name_type_and_scope" {
  name  = "#name#"
  type  = "TOKEN"
  scope = "SYNTHETIC"
  depends_on = [
    dynatrace_credentials_v2.token_synthetic,
    dynatrace_credentials_v2.token_app_engine,
    dynatrace_credentials_v2.username_password_synthetic,
  ]
}

output "matched_ids" {
  value = jsonencode(sort(data.dynatrace_credentials_v2.by_name_type_and_scope.items[*].id))
}

output "expected_ids" {
  value = jsonencode([dynatrace_credentials_v2.token_synthetic.id])
}
