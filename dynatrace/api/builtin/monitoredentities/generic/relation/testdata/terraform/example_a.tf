resource "dynatrace_generic_types" "#name#" {
  name         = "terraform:type"
  enabled      = true
  created_by   = "Terraform"
  display_name = "TerraformTest"
  rules {
    rule {
      icon_pattern          = "{TerraformIcon}"
      id_pattern            = "{TerraformPlaceholder}"
      instance_name_pattern = "{TerraformInstance}"
      attributes {
        attribute {
          key     = "TerraformAttribute"
          pattern = "{TerraformExtraction}"
        }
      }
      required_dimensions {
        required_dimension {
          key = "TerraformDimension"
        }
      }
      sources {
        source {
          condition   = "$eq(TerraformCondition)"
          source_type = "Events"
        }
      }
    }
  }
}
