resource "dynatrace_rum_provider_breakdown" "#name#" {
  report_public_improvement = false
  resource_name             = "#name#"
  resource_type             = "ThirdParty"
  domain_name_pattern_list {
    domain_name_pattern {
      pattern = "Terraform3rdPartyExample.com"
    }
  }
}

