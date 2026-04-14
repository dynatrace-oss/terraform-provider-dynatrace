resource "dynatrace_dashboards_allowlist" "list" {
  allowlist {
    urlpattern {
      rule     = "equals"
      template = "https://www.dynatrace.com/"
    }
    # update => re-create due to set-hash change
    urlpattern {
      rule     = "equals"
      template = "https://www.docs.dynatrace.com/"
    }
  }
}
