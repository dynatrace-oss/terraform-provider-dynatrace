resource "dynatrace_generic_types" "types" {
  name         = "terraform:type"
  enabled      = true
  created_by   = "Terraform"
  display_name = "TerraformTest"
  rules {
    rule {
      icon_pattern          = "{TerraformIcon}"
      id_pattern            = "{TerraformPlaceholder}"
      instance_name_pattern = "{TerraformInstance}"
      # update sources set
      sources {
        source {
          condition   = "$eq(TerraformCondition1)"
          source_type = "Events"
        }
        source {
          condition   = "$eq(TerraformCondition2)"
          source_type = "Events"
        }
      }
      # update attributes set
      attributes {
        attribute {
          key     = "TerraformAttribute1"
          pattern = "{TerraformExtraction1}"
        }
        attribute {
          key     = "TerraformAttribute2"
          pattern = "{TerraformExtraction2}"
        }
      }
      # update required dimensions set
      required_dimensions {
        required_dimension {
          key = "TerraformDimension1"
        }
        required_dimension {
          key = "TerraformDimension2"
        }
      }
    }
  }
}
