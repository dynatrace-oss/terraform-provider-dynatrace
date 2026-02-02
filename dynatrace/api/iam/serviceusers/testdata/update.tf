variable "SERVICE_USER_NAME" {
  description = "The name of the service user."
  type        = string
}

resource "dynatrace_iam_service_user" "test_service_user" {
  name        = var.SERVICE_USER_NAME
  description = "an updated description"
  groups      = []
}