resource "dynatrace_update_windows" "window" {
  count = 2
  name       = "#name#-${count.index}"
  enabled    = true
  recurrence = "ONCE"
  once_recurrence {
    recurrence_range {
      end   = "2023-02-15T04:00:00Z"
      start = "2023-02-15T02:00:00Z"
    }
  }
}

resource "dynatrace_oneagent_updates" "update" {
  scope          = "environment"
  target_version = "latest"
  update_mode    = "AUTOMATIC_DURING_MW"
  maintenance_windows {
    dynamic "maintenance_window" {
      for_each = toset(dynatrace_update_windows.window.*.id)
      content {
        maintenance_window = maintenance_window.value
      }
    }
  }
}
