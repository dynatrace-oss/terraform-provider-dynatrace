data "dynatrace_application" "web_application" {
  name = "Web Application"
}

resource "dynatrace_key_user_action" "acc" {
  application_id = data.dynatrace_application.web_application.id
  domain         = "120.0.0.1"
  name           = "Loading of page /custom"
  type           = "Load"
}
