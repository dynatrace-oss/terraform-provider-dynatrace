resource "dynatrace_limit_outbound_connections" "#name#" {
  allowed_outbound_connections {
    enforced  = true
    host_list = [ "www.dynatrace.com", "www.google.com" ]
  }
}