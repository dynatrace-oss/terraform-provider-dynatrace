# stays the same
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

# new window
resource "dynatrace_update_windows" "window-update" {
  name       = "#name#-new"
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
    # use one old and one new window
    maintenance_window {
      maintenance_window = dynatrace_update_windows.window[0].id
    }
    maintenance_window {
      maintenance_window = dynatrace_update_windows.window-update.id
    }
  }
}
