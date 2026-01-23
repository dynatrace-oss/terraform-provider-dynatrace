variable "GROUP_NAME" {
  description = "The name of the group."
  type        = string
}

variable "ACCOUNT_ID" {
  description = "The ID of the Dynatrace account."
  type        = string
}

resource "dynatrace_iam_group" "my-group" {
  name = var.GROUP_NAME
  description = "A group created for e2e testing with specific permissions."
  permissions {
    permission {
      name  = "account-viewer"
      type  = "account"
      scope = var.ACCOUNT_ID
    }
  }
}