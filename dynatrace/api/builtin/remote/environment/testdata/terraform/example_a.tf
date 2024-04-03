resource "dynatrace_remote_environments" "#name#" {
  name          = "TerraformExample"
  network_scope = "EXTERNAL"
  token         = "################"
  uri           = "https://terraformexample.live.dynatrace.com"
}