resource "dynatrace_notification" "#name#" {
  ansible_tower {
    name                   = "#name#"
    accept_any_certificate = true
    active                 = false
    alerting_profile       = "f75e68ef-aca7-3a07-9c21-94eb00ecfc56"
    custom_message         = "some-custom-message"
    job_template_id        = 444
    job_template_url       = "https://localhost/#/templates/job_template/444"
    password               = "#######"
    username               = "foo"
  }
}