resource "dynatrace_host_monitoring_advanced" "#name#" {
  host_id        = "HOST-1234567890000000"
  process_agent_injection     = true
  code_module_injection = true
}
