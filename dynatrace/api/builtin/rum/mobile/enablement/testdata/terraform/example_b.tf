resource "dynatrace_mobile_app_enablement" "#name#" {
  rum {
    enabled                  = true
    cost_and_traffic_control = 25
  }
  session_replay {
    on_crash = false
    cost_and_traffic_control = 50
    full_session_replay = true
  }
}
