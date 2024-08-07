resource "dynatrace_discovery_default_rules" "#name#" {
  rule {
    description       = "Through end-to-end tracing services on these hosts were found to be externally facing"
    category          = "Deep monitoring"
    environment_scope = false
    id                = "#name#"
    priority          = "CRITICAL"
    query             =<<-EOT
      fetch dt.entity.host, from:-15m
              | filter hasPublicTraffic == true OR ebpfHasPublicTraffic == true
              | fields host=entity.name, monitoringMode, host.id=id, compliant=(monitoringMode == "FULL_STACK")
              
    EOT
    title             = "#name#"
    actions {
      action {
        name = "setMonitoringMode"
        parameters {
          parameter {
            name  = "mode"
            value = "FULL_STACK"
          }
        }
      }
    }
  }
  settings {
    muted = false
  }
}