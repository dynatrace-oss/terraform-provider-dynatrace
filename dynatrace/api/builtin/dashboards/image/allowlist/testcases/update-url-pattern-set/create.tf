resource "dynatrace_dashboards_allowlist" "list" {
  allowlist {
    urlpattern {
      rule     = "equals"
      template = "https://www.dynatrace.com/"
    }
    urlpattern {
      rule     = "startsWith"
      template = "https://www.docs.dynatrace.com/"
    }
  }
}
