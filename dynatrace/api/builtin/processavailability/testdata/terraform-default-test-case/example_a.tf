resource "dynatrace_process_availability" "#name#" {
  enabled = true
  name    = "#name#"
  rules {
    rule {
      property  = "executable"
      condition = "$contains(svc)"
    }
  }
  metadata {
    item {
      key   = "foo"
      value = "bar"
    }
  }
}
