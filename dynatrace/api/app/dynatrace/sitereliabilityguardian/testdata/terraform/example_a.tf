resource "dynatrace_site_reliability_guardian" "#name#" {
  name = "Test"
  tags = [ "stage:staging" ]
  event_kind = "BIZ_EVENT"
  objectives {
    objective {
      name                = "Error rate"
      comparison_operator = "LESS_THAN_OR_EQUAL"
      dql_query           =<<-EOT
        fetch logs
        | fieldsAdd errors = toLong(loglevel == "ERROR")
        | summarize errorRate = sum(errors)/count() * 100
      EOT
      objective_type      = "DQL"
      target              = 8
      warning             = 6
    }
    objective {
      name                = "Count bizevents"
      comparison_operator = "GREATER_THAN_OR_EQUAL"
      dql_query           =<<-EOT
        fetch bizevents
        | summarize count()
      EOT
      objective_type      = "DQL"
      target              = 50
      warning             = 55
    }
  }
}