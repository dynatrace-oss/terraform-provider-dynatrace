resource "dynatrace_web_app_manual_insertion" "#name#" {
  application_id = "APPLICATION-1234567890000000"
  code_snippet {
    code_snippet_type = "DEFERRED"
  }
  javascript_tag {
    cache_duration             = "12"
    crossorigin_anonymous      = false
    script_execution_attribute = "async"
  }
  oneagent_javascript_tag {
    script_execution_attribute = "defer"
  }
  oneagent_javascript_tag_sri {
    script_execution_attribute = "defer"
  }
}