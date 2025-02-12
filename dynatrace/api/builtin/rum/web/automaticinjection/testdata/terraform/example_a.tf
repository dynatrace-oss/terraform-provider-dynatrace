resource "dynatrace_web_app_auto_injection" "#name#" {
  application_id = "APPLICATION-1234567890000000"
  cache_control_headers {
    cache_control_headers = true
  }
  monitoring_code_source_section {
    code_source          = "OneAgent"
    monitoring_code_path = "/testpath/"
  }
  snippet_format {
    code_snippet_type = "DEFERRED"
    snippet_format    = "Code Snippet"
  }
}
