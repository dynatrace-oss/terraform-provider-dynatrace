variable "BOUNDARY_NAME" {
  description = "The name of the policy boundary."
  type        = string
}

resource "dynatrace_iam_policy_boundary" "boundary" {
  name  = var.BOUNDARY_NAME
  query = "environment:management-zone startsWith \"[Bar]\";"
}