resource "dynatrace_process_visibility" "#name#" {
  enabled       = true
  max_processes = 80
  scope         = "environment"
}