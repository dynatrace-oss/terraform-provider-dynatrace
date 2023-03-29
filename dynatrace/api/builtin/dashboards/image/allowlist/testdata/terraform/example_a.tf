resource "dynatrace_dashboards_allowlist" "#name#" {
  allowlist {
    urlpattern {
      rule     = "equals"
      template = "https://www.terraform.io/"
    }
    urlpattern {
      rule     = "startsWith"
      template = "https://www.google.com/"
    }
  }
}