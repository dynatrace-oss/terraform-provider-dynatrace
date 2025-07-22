resource "dynatrace_web_app_manual_insertion" "#name#" {
  application_id = "APPLICATION-1234567890000000"
  code_snippet {
    code_snippet_type = "SYNCHRONOUSLY"
  }
  javascript_tag {
    cache_duration        = "1"
    crossorigin_anonymous = true
  }
  oneagent_javascript_tag {
  }
  oneagent_javascript_tag_sri {
  }
}
