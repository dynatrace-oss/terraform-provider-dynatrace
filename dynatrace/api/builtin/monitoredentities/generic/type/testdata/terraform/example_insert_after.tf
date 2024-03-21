resource "dynatrace_generic_types" "first-instance" {
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

resource "dynatrace_generic_types" "second-instance" {
  name         = "terraform:type:second"
  enabled      = true
  created_by   = "Terraform"
  display_name = "TerraformTest-second"
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
  insert_after = dynatrace_generic_types.first-instance.id
}
