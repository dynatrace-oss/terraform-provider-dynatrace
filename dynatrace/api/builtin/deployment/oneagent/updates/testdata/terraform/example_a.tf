resource "dynatrace_oneagent_updates" "#name#" {
  scope          = "environment"
  target_version = "latest"
  update_mode    = "AUTOMATIC"
}