resource "dynatrace_remote_environments" "env" {
  name          = "#name#"
  network_scope = "EXTERNAL"
  token         = "################"
  uri           = "https://example_#name#.live.dynatrace.com"
}