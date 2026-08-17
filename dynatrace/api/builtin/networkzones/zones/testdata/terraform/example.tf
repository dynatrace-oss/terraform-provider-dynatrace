resource "dynatrace_network_zone_v2" "example" {
  identifier        = "#name#"
  alternative_zones = []
  description       = "my description"
  fallback_mode     = "ANY_ACTIVE_GATE"
}
