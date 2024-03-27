resource "dynatrace_rum_ip_determination" "first-instance" {
  header_name = "X-Header-Example3"
}

resource "dynatrace_rum_ip_determination" "second-instance" {
  header_name  = "X-Header-Example4"
  insert_after = dynatrace_rum_ip_determination.first-instance.id
}
