resource "dynatrace_iam_group" "my-group" {
  name                       = "#name#"
  description                = "A group created for e2e testing."
  federated_attribute_values = ["some-value"]
}
