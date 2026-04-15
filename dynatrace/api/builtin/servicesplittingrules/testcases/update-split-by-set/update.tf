resource "dynatrace_service_splitting" "split" {
  enabled = false
  scope   = "environment"
  rule {
    description = "Example description"
    condition   = "matchesValue(k8s.cluster.name, \"terraform\")"
    rule_name   = "#name#"
    service_splitting_attributes {
      service_splitting_attribute {
        key = "Attribute-1"
      }
      # updated
      service_splitting_attribute {
        key = "Attribute-edit"
      }
    }
  }
}
