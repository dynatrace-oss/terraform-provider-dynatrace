resource "dynatrace_rum_ip_determination" "first-instance" {
  header_name = "X-Header-#name#-3"
}

resource "dynatrace_rum_ip_determination" "second-instance" {
  header_name  = "X-Header-#name#-4"
  insert_after = dynatrace_rum_ip_determination.first-instance.id
}
