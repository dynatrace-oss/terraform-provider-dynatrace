resource "dynatrace_remote_environments" "env" {
  name          = "#name#"
  network_scope = "EXTERNAL"
  token         = "################"
  uri           = "https://example-#name#.live.dynatrace.com"
}