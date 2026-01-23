variable "GROUP_NAME" {
  description = "The name of the group."
  type        = string
}

resource "dynatrace_iam_group" "my-group" {
  name = var.GROUP_NAME
  description = "A group created for e2e testing."
  federated_attribute_values = ["some-value"]
}