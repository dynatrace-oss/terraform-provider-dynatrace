resource "dynatrace_rum_provider_breakdown" "first-instance" {
  report_public_improvement = false
  resource_name             = "#name#"
  resource_type             = "ThirdParty"
  domain_name_pattern_list {
    domain_name_pattern {
      pattern = "Terraform3rdPartyExample.com"
    }
  }
}

resource "dynatrace_rum_provider_breakdown" "second-instance" {
  report_public_improvement = false
  resource_name             = "#name#-second"
  resource_type             = "ThirdParty"
  domain_name_pattern_list {
    domain_name_pattern {
      pattern = "Terraform3rdPartyExample.com"
    }
  }
  insert_after = dynatrace_rum_provider_breakdown.first-instance.id
}

