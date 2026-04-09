resource "dynatrace_default_launchpad" "launchpad" {
  group_launchpads {
    group_launchpad {
      is_enabled = false
      launchpad_id = "00000000-0000-0000-0000-000000000001"
      user_group_id = "00000000-0000-0000-0000-000000000001"
    }
    group_launchpad {
      is_enabled = true
      launchpad_id = "00000000-0000-0000-0000-000000000003"
      user_group_id = "00000000-0000-0000-0000-000000000003"
    }
  }
}
