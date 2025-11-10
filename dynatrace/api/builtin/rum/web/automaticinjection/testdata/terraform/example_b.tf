resource "dynatrace_web_app_auto_injection" "#name#" {
  application_id = "APPLICATION-1234567890000000"
  cache_control_headers {
    cache_control_headers = true
  }
  monitoring_code_source_section {
    code_source          = "OneAgent"
  }
  snippet_format {
    snippet_format    = "OneAgent JavaScript Tag"
  }
}
