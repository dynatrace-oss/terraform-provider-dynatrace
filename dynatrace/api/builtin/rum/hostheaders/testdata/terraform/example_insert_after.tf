resource "dynatrace_rum_host_headers" "first-instance" {
  header_name = "#name#"
}

resource "dynatrace_rum_host_headers" "second-instance" {
  header_name  = "#name#-second"
  insert_after = dynatrace_rum_host_headers.first-instance.id
}
