resource "dynatrace_alerting" "#name#" {
  name = "#name#"
  filters {
    filter {
      custom {
        metadata {
          items {
            filter {
              key   = "POC"
              value = "GRAIL"
            }
          }
        }
      }
    }
  }
}