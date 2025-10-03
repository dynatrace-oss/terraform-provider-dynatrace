resource "dynatrace_document" "this" {
  type = "dashboard"
  name = "Example Dashboard"
  content = jsonencode(
    {
      "version" : 13,
      "variables" : [],
      "tiles" : {}
    }
  )
}
