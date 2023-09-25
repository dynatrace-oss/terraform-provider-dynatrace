resource "dynatrace_generic_setting" "ABC" {
  schema = "app:my.booking.analytics:connection"
  scope  = "environment"
  value = jsonencode({
    "client_id"     : "asdfhh",
    "client_secret" : "mysecret",
    "name"          : "ABC",
    "tenant_id"     : "asdf",
    "type"          : "client_secret",
    "user_id"       : "asdf"
  })
}
