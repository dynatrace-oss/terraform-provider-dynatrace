variable "GROUP_NAME" {
  description = "The name of the group."
  type        = string
}

resource "dynatrace_iam_group" "my-group" {
  name = var.GROUP_NAME
  description = "A group created for e2e testing."
}

variable "SERVICE_USER_NAME" {
  description = "The name of the service user."
  type        = string
}

resource "dynatrace_iam_service_user" "test_service_user" {
  name        = var.SERVICE_USER_NAME
  description = "a service user for testing purposes"
  groups      = [  dynatrace_iam_group.my-group.id ]
}