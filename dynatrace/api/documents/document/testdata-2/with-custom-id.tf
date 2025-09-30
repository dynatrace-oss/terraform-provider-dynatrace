resource "dynatrace_document" "this" {
  type = "dashboard"
  name = "Example Dashboard"
  custom_id = "#name#"
  content = jsonencode(
    {
      "version" : 13,
      "variables" : [],
      "tiles" : {}
    }
  )
}
