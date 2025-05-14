resource "dynatrace_service_detection_rules" "#name#" {
  enabled      = false
  scope        = "environment"
  rule {
    description                    = "Example description"
    additional_required_attributes = [ "attribute-1", "attribute-2" ]
    condition                      = "matchesValue(k8s.cluster.name,\"terraform\")"
    rule_name                      = "#name#"
    service_name_template          = "{k8s.workload.name}"
  }
}
