resource "dynatrace_credentials_v2" "secret_id" {
  name              = "#name#"
  token             = "my-secret-id"
  owner_access_only = true
  scopes            = ["SYNTHETIC"]
}

resource "dynatrace_credentials_v2" "hashicorp_approle" {
  name              = "#name#"
  owner_access_only = true
  external {
    vault_url            = "https://example.com:8020"
    vault_namespace      = "my-namespace"
    path_to_credentials  = "secret/my-credentials"
    roleid               = "my-role-id"
    secretid             = dynatrace_credentials_v2.secret_id.id
    username_secret_name = "my-username-secret"
    password_secret_name = "my-password-secret"
  }
  scopes = ["SYNTHETIC"]
}
