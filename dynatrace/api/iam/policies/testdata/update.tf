variable "ACCOUNT_ID" {
  description = "The ID of the Dynatrace account."
  type        = string
}

variable "POLICY_NAME" {
  description = "The name of the policy."
  type        = string
}

resource "dynatrace_iam_policy" "policy" {
  name            = var.POLICY_NAME
  account         = var.ACCOUNT_ID
  statement_query = "ALLOW settings:objects:read, settings:schemas:read WHERE settings:schemaId = \"other-string\";"
}
