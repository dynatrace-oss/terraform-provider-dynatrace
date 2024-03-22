resource "dynatrace_process_availability" "first-instance" {
  enabled = true
  name    = "${randomize}"
  rules {
    rule {
      property  = "executable"
      condition = "$contains(svc)"
    }
  }
}

resource "dynatrace_process_availability" "second-instance" {
  enabled = true
  name    = "${randomize}-second"
  rules {
    rule {
      property  = "executable"
      condition = "$contains(svc)"
    }
  }
  insert_after = dynatrace_process_availability.first-instance.id
}
