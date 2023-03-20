resource "dynatrace_remote_environments" "TerraformExample" {
  name          = "TerraformExample"
  network_scope = "EXTERNAL"
  token         = "################"
  uri           = "https://terraformexample.live.dynatrace.com"
}