resource "dynatrace_credentials_v2" "synthetic" {
  name              = "#name#"
  token             = "token"
  owner_access_only = false
  scopes            = ["SYNTHETIC"]
}

resource "dynatrace_credentials_v2" "app_engine" {
  name              = "#name#"
  token             = "token"
  owner_access_only = false
  scopes            = ["APP_ENGINE"]
}

data "dynatrace_credentials_v2" "by_name_and_scope" {
  name       = "#name#"
  scope      = "SYNTHETIC"
  depends_on = [dynatrace_credentials_v2.synthetic, dynatrace_credentials_v2.app_engine]
}

output "matched_ids" {
  value = jsonencode(sort(data.dynatrace_credentials_v2.by_name_and_scope.items[*].id))
}

output "expected_ids" {
  value = jsonencode([dynatrace_credentials_v2.synthetic.id])
}
