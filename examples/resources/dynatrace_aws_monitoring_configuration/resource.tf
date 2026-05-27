# Minimal example consumed by terraform-plugin-docs to generate
# docs/resources/aws_monitoring_configuration.md. For the full
# end-to-end flow (HAS connection -> IAM role + trust -> role-ARN patch
# -> monitoring config), see examples/aws-monitoring-end-to-end/.

resource "dynatrace_aws_monitoring_configuration" "this" {
  name          = "prod-aws-monitoring"
  connection_id = dynatrace_aws_connection.this.id
  account_id    = data.aws_caller_identity.current.account_id
  regions       = ["eu-central-1", "us-east-1"]

  # extension_version is Optional+Computed: omit it and the provider
  # resolves the highest semver installed on the tenant at create time
  # (mirrors `dtctl create aws`). Pin it only when you need a specific
  # version.
  # extension_version = "1.0.7"

  feature_sets = [
    "EC2_essential",
    "RDS_essential",
    "S3_essential",
  ]

  # Optional typed blocks. Each replaces what used to live behind a
  # generic JSON escape hatch; everything is schema-validated and
  # surfaces in plan diffs.

  # tag_enrichment = ["owner", "cost-center", "environment"]

  # tag_filter {
  #   key       = "environment"
  #   value     = "production"
  #   condition = "INCLUDE"
  # }

  # cloud_watch_logs {
  #   enabled = true
  #   regions = ["eu-central-1"]
  # }

  # dt_label_enrichment {
  #   label   = "dt.cost.product"
  #   tag_key = "product"   # XOR with `literal`
  # }

  # custom_namespace {
  #   namespace              = "MyApp/Metrics"
  #   auto_discovery_enabled = false
  #   metric {
  #     name         = "queue.depth"
  #     unit         = "Count"
  #     dimensions   = ["QueueName"]
  #     aggregations = ["Average", "Maximum"]
  #     type         = "CUSTOM"
  #   }
  # }
}
