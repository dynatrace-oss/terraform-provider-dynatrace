locals {
    mobile_application_name = "Application"
}

data "dynatrace_mobile_application" "setup_application" {
  name = local.mobile_application_name
}

import {
  to = dynatrace_mobile_application.application
  for_each = try(data.dynatrace_mobile_application.setup_application.id, null) == null ? [] : [data.dynatrace_mobile_application.setup_application.id]
  id = each.value
}

resource "dynatrace_mobile_application" "application" {
  name = local.mobile_application_name
  beacon_endpoint_type    = "INSTRUMENTED_WEB_SERVER"
  beacon_endpoint_url     = "https://dynatrace.com/dtmb"
  application_type         = "MOBILE_APPLICATION"
  user_session_percentage = 100
  apdex {
    frustrated          = 12000
    frustrated_on_error = true
    tolerable           = 3000
  }

  lifecycle {
    prevent_destroy = true
  }
}
