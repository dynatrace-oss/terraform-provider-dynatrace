resource "dynatrace_activegate_updates" "example" {
  scope          = "environment"
  target_version = "latest"
  update_mode    = "AUTOMATIC_DURING_UW"
  update_windows {
    // one window removed
    update_window {
      update_window = dynatrace_update_windows.example[1].id
    }
  }
}

resource "dynatrace_update_windows" "example" {
  count      = 2
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
