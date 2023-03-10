resource "dynatrace_service_http_failure" "#name#" {
  enabled    = true
  service_id = "SERVICE-D892CFE7DFAB0D08"
  broken_links {
    broken_link_domains         = [ "www.google.com" ]
    http_404_not_found_failures = true
  }
  http_response_codes {
    client_side_errors                        = "401-599"
    fail_on_missing_response_code_client_side = true
    fail_on_missing_response_code_server_side = true
    server_side_errors                        = "501-599"
  }
}