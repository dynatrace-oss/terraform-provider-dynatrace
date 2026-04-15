resource "dynatrace_process_availability" "availability" {
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
      key   = "key1"
      value = "value1"
    }
    # updated
    item {
      key   = "keyEdit"
      value = "valueEdit"
    }
  }
}
