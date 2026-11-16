resource "dynatrace_credentials_v2" "client_secret" {
  name              = "#name#"
  token             = "my-client-secret"
  owner_access_only = true
  scopes            = ["SYNTHETIC"]
}

resource "dynatrace_credentials_v2" "azure_client_secret" {
  name              = "#name#"
  owner_access_only = true
  external {
    vault_url            = "https://example.vault.azure.net"
    tenantid             = "sometenant"
    clientid             = dynatrace_credentials_v2.client_secret.id
    client_secret        = dynatrace_credentials_v2.client_secret.id
    username_secret_name = "my-username-secret"
    password_secret_name = "my-password-secret"
  }
  scopes = ["SYNTHETIC"]
}
