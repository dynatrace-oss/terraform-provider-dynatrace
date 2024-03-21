resource "dynatrace_business_events_processing" "first-instance" {
  enabled   = true
  matcher   = "matchesValue(event.type, \"com.easytrade.buy-assets\")"
  rule_name = "#name#"
  script    = "FIELDS_ADD(trading_volume:price*amount)"
  rule_testing {
    sample_event = jsonencode({
      "id" : "OR-838475",
      "paymentType" : "paypal",
      "plannedDeliveryDate" : "01.01.2021",
      "total" : 234
    })
  }
  transformation_fields {
    transformation_field {
      name     = "amount"
      type     = "DOUBLE"
      array    = false
      optional = false
      readonly = true
    }
    transformation_field {
      name     = "price"
      type     = "DOUBLE"
      array    = false
      optional = false
      readonly = true
    }
  }
}

resource "dynatrace_business_events_processing" "second-instance" {
  enabled   = true
  matcher   = "matchesValue(event.type, \"com.easytrade.buy-assets\")"
  rule_name = "#name#-second"
  script    = "FIELDS_ADD(trading_volume:price*amount)"
  rule_testing {
    sample_event = jsonencode({
      "id" : "OR-838475",
      "paymentType" : "paypal",
      "plannedDeliveryDate" : "01.01.2021",
      "total" : 234
    })
  }
  transformation_fields {
    transformation_field {
      name     = "amount"
      type     = "DOUBLE"
      array    = false
      optional = false
      readonly = true
    }
    transformation_field {
      name     = "price"
      type     = "DOUBLE"
      array    = false
      optional = false
      readonly = true
    }
  }
  insert_after = dynatrace_business_events_processing.first-instance.id
}
