resource "dynatrace_process_availability" "example" {
  enabled = true
  name    = "${randomize}"
}
