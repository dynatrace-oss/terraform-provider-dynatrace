resource "dynatrace_maintenance_windows" "recurring_window" {
  name        = "Example recurring #name#"
  description = "Example description"
  filter      = "matchesValue(result.state, \"FAIL\")"
  auto_delete = true
  enabled     = true

  schedule {
    duration = 60
    timezone = "UTC"

    trigger {
      type = "time"

      recurring {
        time           = "12:00:00"
        earliest_start = "2100-01-01"
        until          = "2100-01-01"
        # rule = dynatrace_automation_scheduling_rule.scheduling_rule.id
      }
    }
  }
}

resource "dynatrace_maintenance_windows" "once_window" {
  name        = "Example once #name#"
  description = "Example description"
  filter      = "matchesValue(result.state, \"FAIL\")"
  auto_delete = true
  enabled     = true

  schedule {
    duration = 60
    timezone = "UTC"

    trigger {
      type = "once"

      once {
        date = "2100-01-01T00:00:00"
      }
    }
  }
}
