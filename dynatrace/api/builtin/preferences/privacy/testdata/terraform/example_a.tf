resource "dynatrace_data_privacy" "#name#" {
  application_id = "environment"
  data_collection {
    opt_in_mode_enabled = true
  }
  do_not_track {
    comply_with_do_not_track = false
  }
  masking {
    ip_address_masking                = "public"
    ip_address_masking_enabled        = true
    personal_data_uri_masking_enabled = true
    user_action_masking_enabled       = true
  }
  user_tracking {
    persistent_cookie_enabled = true
  }
}