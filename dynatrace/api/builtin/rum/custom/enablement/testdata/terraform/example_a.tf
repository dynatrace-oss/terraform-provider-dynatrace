resource "dynatrace_custom_app_enablement" "#name#" {
  scope = "environment"
  rum {
    enabled                  = false
    cost_and_traffic_control = 40
  }
}
