resource "dynatrace_ownership_teams" "teams" {
  name        = "#name#"
  description = "Created by Terraform"
  identifier  = "Terraform_#name#"
  additional_information {
    additional_information {
      key   = "HashiCorp"
      url   = "https://www.terraform.io/"
      value = "Terraform"
    }
  }
  contact_details {
    contact_detail {
      email            = "terraform@dynatrace.com"
      integration_type = "EMAIL"
    }
  }
  links {
    link {
      link_type = "URL"
      url       = "https://www.google.com"
    }
  }
  responsibilities {
    development      = true
    infrastructure   = false
    line_of_business = false
    operations       = true
    security         = false
  }
  supplementary_identifiers {
    supplementary_identifier {
      supplementary_identifier = "identifier1"
    }
    # updated
    supplementary_identifier {
      supplementary_identifier = "identifierEdit"
    }
  }
}
