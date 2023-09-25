resource "dynatrace_host_monitoring" "#name#" {
  enabled        = true
  auto_injection = false
  host_id        = "HOST-1234567890000000"
}
