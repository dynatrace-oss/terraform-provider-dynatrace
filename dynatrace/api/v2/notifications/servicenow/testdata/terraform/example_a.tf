resource "dynatrace_service_now_notification" "#name#" { # replace #name# with the name you would like your resource be known within your Terraform Module
  active    = false
  name      = "#name#" # replace #name# with the name you would like your entry to be displayed within the Dynatrace Web UI
  profile   = dynatrace_alerting.Default.id
  instance  = "#name#"
  username  = "service-now-username"
  password  = "service-now-password"
  message   = "service-now-message"
  incidents = true
  events    = true
}

resource "dynatrace_alerting" "Default" {
  name = "#name#"
}
