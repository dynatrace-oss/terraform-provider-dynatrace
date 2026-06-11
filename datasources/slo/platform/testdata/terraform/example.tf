resource "dynatrace_platform_slo" "slo" {
  name        = "#name#"
  description = "Sample custom SLO"
  tags        = [ "ExampleKey:ExampleValue" ]
  criteria {
    criteria_detail {
      target         = 96
      timeframe_from = "now-30d"
      timeframe_to   = "now"
      warning        = 99
    }
  }
  custom_sli {
    indicator =<<-EOT
      timeseries { total=sum(dt.service.request.count) ,failures=sum(dt.service.request.failure_count) }, by: { dt.entity.service }
      | fieldsAdd tags=entityAttr(dt.entity.service, "tags")
      | filter in(tags, "criticality:Gold")
      | fieldsAdd entityName = entityName(dt.entity.service)
      | fieldsAdd sli=(((total[]-failures[])/total[])*(100))
      | fieldsRemove total, failures, tags
    EOT
  }
}

data "dynatrace_platform_slo" "slo" {
  name = dynatrace_platform_slo.slo.name
}


output "id" {
  value = data.dynatrace_platform_slo.slo.id
}

output "sli" {
  value = data.dynatrace_platform_slo.slo.custom_sli[0].indicator
}
