resource "dynatrace_web_app_custom_prop_restrictions" "#name#" {
  application_id = "APPLICATION-1234567890000000"
  custom_event_properties_allow_list {
    custom_session_properties_allow {
      field_data_type = "STRING"
      field_name      = "ExampleEvent"
    }
    custom_session_properties_allow {
      field_data_type = "BOOLEAN"
      field_name      = "ExampleEvent2"
    }
  }
  custom_session_properties_allow_list {
    custom_session_properties_allow {
      field_data_type = "STRING"
      field_name      = "ExampleSession"
    }
    custom_session_properties_allow {
      field_data_type = "BOOLEAN"
      field_name      = "ExampleSession2"
    }
  }
}
