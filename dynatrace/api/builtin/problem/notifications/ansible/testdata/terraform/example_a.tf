resource "dynatrace_ansible_tower_notification" "#name#" { # replace #name# with the name you would like your resource be known within your Terraform Module
  active           = false
  name             = "#name#" # replace #name# with the name you would like your entry to be displayed within the Dynatrace Web UI
  profile          = dynatrace_alerting.Default.id
  insecure         = true
  job_template_url = "https://localhost/#/templates/job_template/999"
  username         = "foo"
  password         = "bar"
  custom_message   = "some-custom-message"
}

resource "dynatrace_alerting" "Default" {
  name = "#name#"
}
