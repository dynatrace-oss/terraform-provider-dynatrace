resource "dynatrace_process_availability" "first-instance" {
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

resource "dynatrace_process_availability" "second-instance" {
  enabled = true
  name    = "#name#-second"
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
  insert_after = dynatrace_process_availability.first-instance.id
}
