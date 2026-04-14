# stays the same
resource "dynatrace_update_windows" "window" {
  count = 3
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
    maintenance_window {
      maintenance_window = dynatrace_update_windows.window[0].id
    }
    # updated
    maintenance_window {
      maintenance_window = dynatrace_update_windows.window[2].id
    }
  }
}
