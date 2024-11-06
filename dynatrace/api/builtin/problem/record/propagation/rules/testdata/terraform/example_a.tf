resource "dynatrace_problem_record_propagation_rules" "#name#" {
    enabled = false
    source_attribute_key = "terraformSource"
    target_attribute_key = "terraformTarget"
}
