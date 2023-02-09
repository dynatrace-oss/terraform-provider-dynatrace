resource "dynatrace_web_app_enablement" "#name#" {
  rum {
    enabled                  = true
    cost_and_traffic_control = 25
  }
  session_replay {
    enabled = false
    cost_and_traffic_control = 25
  }
}
