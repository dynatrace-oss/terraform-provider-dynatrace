resource "dynatrace_app_monitoring" "monitor" {
  default_log_level = "off"
  default_trace_level = "off"
  app_monitoring {
    app_monitoring {
      app_id           = "app:dynatrace.github.connector:connection"
      custom_log_level = "info"
    }
    app_monitoring {
      app_id           = "app:dynatrace.site.reliability.guardian:guardians"
      custom_log_level = "debug"
    }
  }
}
