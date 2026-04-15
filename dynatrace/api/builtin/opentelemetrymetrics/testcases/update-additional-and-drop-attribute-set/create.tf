resource "dynatrace_opentelemetry_metrics" "metrics" {
  additional_attributes_to_dimension_enabled = true
  meter_name_to_dimension_enabled            = true
  scope                                      = "environment"
  enable_mint_v_2_ingest = true
  additional_attributes {
    additional_attribute {
      enabled       = true
      attribute_key = "terraform.test.add1"
    }
    # to update
    additional_attribute {
      enabled       = true
      attribute_key = "terraform.test.add2"
    }
  }
  to_drop_attributes {
    to_drop_attribute {
      enabled       = true
      attribute_key = "terraform.test.drop1"
    }
    # to update
    to_drop_attribute {
      enabled       = true
      attribute_key = "terraform.test.drop2"
    }
  }
}
