resource "dynatrace_generic_setting" "ABC" {
  schema = "app:dynatrace.site.reliability.guardian:guardians"
  scope  = "environment"
  value = jsonencode({
    "name": "#name#",
    "tags": [ "stage:staging" ],
    "eventKind": "BIZ_EVENT",
    "objectives": [
      {
        "name": "Error rate",
        "comparisonOperator": "LESS_THAN_OR_EQUAL",
        "dqlQuery": <<-EOT
        fetch logs
        | fieldsAdd errors = toLong(loglevel == "ERROR")
        | summarize errorRate = sum(errors)/count() * 100
      EOT
        "objectiveType": "DQL",
        "target": 8,
        "warning": 6
      }
    ]
  })
}
