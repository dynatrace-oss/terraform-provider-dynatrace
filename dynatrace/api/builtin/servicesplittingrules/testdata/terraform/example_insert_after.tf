resource "dynatrace_service_splitting" "first-instance" {
  enabled = false
  scope   = "environment"
  rule {
    description = "Example description"
    condition   = "matchesValue(k8s.cluster.name, \"terraform\")"
    rule_name   = "#name#-first"
    service_splitting_attributes {
      service_splitting_attribute {
        key = "Attribute-1"
      }
    }
  }
}

resource "dynatrace_service_splitting" "second-instance" {
  enabled = false
  insert_after = "${dynatrace_service_splitting.first-instance.id}"
  scope   = "environment"
  rule {
    description = "Example description"
    condition   = "matchesValue(k8s.cluster.name, \"terraform\")"
    rule_name   = "#name#-second"
    service_splitting_attributes {
      service_splitting_attribute {
        key = "Attribute-1"
      }
      service_splitting_attribute {
        key = "Attribute-2"
      }
      service_splitting_attribute {
        key = "Attribute-3"
      }
    }
  }
}
