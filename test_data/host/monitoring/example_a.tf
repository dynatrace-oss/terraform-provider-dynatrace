resource "dynatrace_host_monitoring" "#name#" {
  enabled        = true
  auto_injection = false
  full_stack     = false
  host_id        = "HOST-1234567890000000"
}
