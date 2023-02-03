resource "dynatrace_muted_requests" "#name#" {
  muted_request_names = [ "/healthcheck", "/heartbeat" ]
  service_id          = "SERVICE-1234567890000000"
}
