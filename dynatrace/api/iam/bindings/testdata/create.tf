variable "ACCOUNT_ID" {
  description = "The ID of the Dynatrace account."
  type        = string
}

variable "GROUP_NAME" {
  description = "The name of the group."
  type        = string
}

resource "dynatrace_iam_group" "group" {
  name = var.GROUP_NAME
}

variable "POLICY_NAME_1" {
  description = "The name of the first policy."
  type        = string
}

resource "dynatrace_iam_policy" "policy-1" {
  name            = var.POLICY_NAME_1
  account         = var.ACCOUNT_ID
  statement_query = "ALLOW settings:objects:read, settings:schemas:read WHERE settings:schemaId = \"some-schema-id\";"
}

resource "dynatrace_iam_policy_bindings" "bindings" {
  group       = dynatrace_iam_group.group.id
  account     = var.ACCOUNT_ID
  policies    = [dynatrace_iam_policy.policy-1.id]
}


