variable "ACCOUNT_ID" {
  description = "The ID of the Dynatrace account."
  type        = string
}

resource "dynatrace_iam_group" "my-group" {
  name        = "#name#"
  description = "A group created for e2e testing with specific permissions."
  permissions {
    permission {
      name  = "account-viewer"
      type  = "account"
      scope = var.ACCOUNT_ID
    }
  }
}
