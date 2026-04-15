resource "dynatrace_generic_relationships" "relationship" {
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
    # update source directly
    source {
      condition   = "$eq(terraformEdit)"
      source_type = "Metrics"
    }
    # update mapping rules
    source {
      source_type = "Entities"
      mapping_rules {
        mapping_rule {
          destination_property       = "dest1"
          destination_transformation = "Leave text as-is"
          source_property            = "source1"
          source_transformation      = "Leave text as-is"
        }
        mapping_rule {
          destination_property       = "destEdit"
          destination_transformation = "Leave text as-is"
          source_property            = "sourceEdit"
          source_transformation      = "Leave text as-is"
        }
      }
    }
  }
}
