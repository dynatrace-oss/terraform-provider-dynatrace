resource "dynatrace_log_metrics" "#name#" {
  enabled           = true
  dimensions        = ["dt.os.type", "dt.entity.process_group"]
  key               = "log.#name#"
  measure           = "ATTRIBUTE"
  measure_attribute = "dt.entity.host"
  query             = "matchesPhrase(content, \"terratest\")"
}
