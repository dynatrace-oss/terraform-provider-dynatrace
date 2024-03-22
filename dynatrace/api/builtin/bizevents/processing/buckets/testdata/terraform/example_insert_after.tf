resource "dynatrace_business_events_buckets" "first-instance" {
  enabled     = true
  bucket_name = "default_bizevents"
  matcher     = "matchesValue(event.type, \"com.easytrade.buy-assets\")"
  rule_name   = "#name#"
}

resource "dynatrace_business_events_buckets" "second-instance" {
  enabled      = true
  bucket_name  = "default_bizevents"
  matcher      = "matchesValue(event.type, \"com.easytrade.buy-assets\")"
  rule_name    = "#name#-second"
  insert_after = dynatrace_business_events_buckets.first-instance.id
}
