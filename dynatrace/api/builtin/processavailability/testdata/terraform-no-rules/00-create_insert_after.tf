resource "dynatrace_process_availability" "first-instance" {
  enabled = true
  name    = "${randomize}"
}

resource "dynatrace_process_availability" "second-instance" {
  enabled      = true
  name         = "${randomize}-2"
  insert_after = dynatrace_process_availability.first-instance.id
}
