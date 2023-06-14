resource "dynatrace_vmware" "#name#" {
  enabled   = false
  ipaddress = "vcenter01"
  label     = "#name#"
  password  = "################"
  username  = "terraform"
}