resource "dynatrace_ip_address_masking" "#name#" {
  scope = "APPLICATION-1234567890000000"
  enabled = true
  type = "public"
}