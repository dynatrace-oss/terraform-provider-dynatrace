resource "dynatrace_session_replay_web_privacy" "#name#" {
  application_id             = "APPLICATION-1234567890000000"
  enable_opt_in_mode         = false
  url_exclusion_pattern_list = [ "www.google.com" ]
  masking_presets {
    playback_masking_preset  = "MASK_ALL"
    recording_masking_preset = "ALLOW_LIST"
    recording_masking_allow_list_rules {
      allow_list_rule {
        css_expression = "selector.example"
        target         = "ELEMENT"
      }
    }
  }
}