resource "dynatrace_service_detection_rules" "first-instance" {
  enabled      = false
  scope        = "environment"
  rule {
    condition                      = "matchesValue(k8s.cluster.name,\"terraform\")"
    rule_name                      = "#name#-first"
    service_name_template          = "{k8s.workload.name}"
  }
}
  
resource "dynatrace_service_detection_rules" "second-instance" {
  enabled      = false
  insert_after = "${dynatrace_service_detection_rules.first-instance.id}"
  scope        = "environment"
  rule {
    description                    = "Example description"
    additional_required_attributes = [ "attribute-1", "attribute-2" ]
    condition                      = "matchesValue(k8s.cluster.name,\"terraform\")"
    rule_name                      = "#name#-second"
    service_name_template          = "{k8s.workload.name}"
  }
}