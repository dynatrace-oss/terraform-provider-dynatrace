resource "dynatrace_session_replay_resource_capture" "#name#" {
  application_id                              = "environment"
  enable_resource_capturing                   = true
  resource_capture_url_exclusion_pattern_list = [ "www.google.com" ]
}