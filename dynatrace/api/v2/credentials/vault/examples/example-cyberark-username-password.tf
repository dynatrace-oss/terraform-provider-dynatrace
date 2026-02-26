variable "credentials_username" {
  sensitive = true
}

variable "credentials_password" {
  sensitive = true
}

resource "dynatrace_credentials" "username_password_credentials" {
  name              = "#name#"
  username          = var.credentials_username
  password          = var.credentials_password
  owner_access_only = true
  scopes            = ["SYNTHETIC"]
}

resource "dynatrace_credentials" "cyberark_username_password" {
  name              = "#name#"
  owner_access_only = true
  external {
    vault_url                 = "https://example.com"
    application_id            = "my-application-id"
    safe_name                 = "my-safe-name"
    folder_name               = "my-folder-name"
    account_name              = "my-account-name"
    username_password_for_cpm = dynatrace_credentials.username_password_credentials.id
  }
  scopes = ["SYNTHETIC"]
}
