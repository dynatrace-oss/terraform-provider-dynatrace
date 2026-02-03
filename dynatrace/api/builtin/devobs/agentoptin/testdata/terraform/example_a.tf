variable "PROCESS_GROUP_ID" {
  type = string
}

resource "dynatrace_devobs_agent_optin" "optin" {
  scope   = var.PROCESS_GROUP_ID
  enabled = false
}
