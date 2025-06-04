resource "dynatrace_web_app_ip_address_exclusion" "#name#" {
  application_id               = "APPLICATION-1234567890000000"
  ip_address_exclusion_include = false
  ip_exclusion_list {
    ip_exclusion {
      ip = "192.168.1.5"
    }
    ip_exclusion {
      ip    = "10.0.0.1"
      ip_to = "10.0.0.5"
    }
  }
}
