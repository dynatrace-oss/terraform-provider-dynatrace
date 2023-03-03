resource "dynatrace_generic_relationships" "#name#" {
  enabled          = true
  created_by       = "Terraform"
  from_role        = "terraformrole"
  from_type        = "os:service"
  to_role          = "terraformrole"
  to_type          = "terraformdestination"
  type_of_relation = "PART_OF"
  sources {
    source {
      condition   = "$eq(terraform)"
      source_type = "Metrics"
    }
  }
}