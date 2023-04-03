resource "dynatrace_network_traffic" "#name#" {
  host_id = "HOST-1234567890000000"
  exclude_ip {
    ip_address_form {
      ip_address = "192.168.0.1"
    }
  }
  exclude_nic {
    nic_form {
      interface = "terraform_test"
      os        = "OS_TYPE_WINDOWS"
    }
  }
}