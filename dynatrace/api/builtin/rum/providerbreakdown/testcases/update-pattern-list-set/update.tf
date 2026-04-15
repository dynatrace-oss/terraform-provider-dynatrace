resource "dynatrace_rum_provider_breakdown" "breakdown" {
  report_public_improvement = false
  resource_name             = "#name#"
  resource_type             = "ThirdParty"
  domain_name_pattern_list {
    domain_name_pattern {
      pattern = "dynatrace.com/1"
    }
    # updated
    domain_name_pattern {
      pattern = "dynatrace.com/edit"
    }
  }
}

