variable "ACCOUNT_ID" {
  description = "The ID of the Dynatrace account."
  type        = string
}

resource "dynatrace_iam_group" "my_group" {
  name          = "#name#"
  lifecycle {
    ignore_changes = [permissions]
  }
}

resource "dynatrace_iam_policy" "policy" {
  name            = "#name#"
  account = var.ACCOUNT_ID
  statement_query = "ALLOW settings:objects:read;"
}

resource "dynatrace_iam_policy" "acc_policy" {
  name            = "#name#-2"
  account     = var.ACCOUNT_ID
  statement_query = "ALLOW settings:schemas:read;"
}

resource "dynatrace_iam_policy_bindings_v2" "acc_bindings" {
  group       = dynatrace_iam_group.my_group.id
  account     = var.ACCOUNT_ID
  policy {
    id = dynatrace_iam_policy.policy.id
  }
  policy {
    id = dynatrace_iam_policy.acc_policy.id
    parameters = {
      "prop-b" : "value-b"
    }
    metadata = {
      "prop-b" : "value-c"
    }
  }
}
