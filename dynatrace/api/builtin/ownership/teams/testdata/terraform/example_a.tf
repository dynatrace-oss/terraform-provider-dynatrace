# ID vu9U3hXa3q0AAAABABdidWlsdGluOm93bmVyc2hpcC50ZWFtcwAGdGVuYW50AAZ0ZW5hbnQAJDYzMDE3YzMzLTdlYzUtMzc1Zi1iODdkLTcyNzM0MmRkMTlkZb7vVN4V2t6t
resource "dynatrace_ownership_teams" "#name#" {
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
      email            = "kodai.ishikawa@dynatrace.com"
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
}
