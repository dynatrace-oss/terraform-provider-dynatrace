resource "dynatrace_document" "this" {
  type = "dashboard"
  name = "Example Dashboard"
  custom_id = "00000000-0000-0000-0000-000000000000"
  content = jsonencode(
    {
      "version" : 13,
      "variables" : [],
      "tiles" : {}
    }
  )
}
