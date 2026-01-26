variable "ACCOUNT_ID" {
  description = "The ID of the Dynatrace account."
  type        = string
}

variable "GROUP_NAME" {
  description = "The name of the group."
  type        = string
}

resource "dynatrace_iam_group" "my-group" {
  name = var.GROUP_NAME
  description = "A group created for e2e testing."
  federated_attribute_values = ["some-value"]

  lifecycle {
    ignore_changes = [permissions]
  }
}

resource "dynatrace_iam_permission" "perm_a" {
  name            = "account-viewer"
  group           = dynatrace_iam_group.my-group.id
  account = var.ACCOUNT_ID
}