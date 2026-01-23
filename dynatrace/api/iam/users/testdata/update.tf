variable "USER_EMAIL" {
  description = "The email of the Dynatrace IAM user."
  type        = string
}

// Every user always belongs to at least "Default group with all users"
data "dynatrace_iam_group" "all_users_group" {
  name = "Default group with all users"
}

resource "dynatrace_iam_user" "my_user" {
  email  = var.USER_EMAIL
  groups = [data.dynatrace_iam_group.all_users_group.id]
}
