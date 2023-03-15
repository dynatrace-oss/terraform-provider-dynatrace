resource "dynatrace_activegate_token" "#name#" {
  auth_token_enforcement_manually_enabled = false
  expiring_token_notifications_enabled    = true
}