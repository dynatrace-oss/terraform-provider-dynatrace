resource "dynatrace_credentials_v2" "token_credentials" {
  name              = "#name#"
  token             = "token"
  owner_access_only = false
  scopes            = ["EXTENSION_AUTHENTICATION"]
}
