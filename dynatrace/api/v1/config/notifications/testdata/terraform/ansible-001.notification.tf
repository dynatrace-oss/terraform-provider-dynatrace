resource "dynatrace_notification" "#name#" {
  ansible_tower {
    name                   = "#name#"
    accept_any_certificate = true
    active                 = false
    alerting_profile       = dynatrace_alerting_profile.Default.id
    custom_message         = "some-custom-message"
    job_template_id        = 444
    job_template_url       = "https://localhost/#/templates/job_template/444"
    password               = "#######"
    username               = "foo"
  }
}

resource "dynatrace_alerting_profile" "Default" {
  display_name = "#name#"
}
