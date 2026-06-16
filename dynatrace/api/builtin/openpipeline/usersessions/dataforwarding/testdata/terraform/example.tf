resource "dynatrace_openpipeline_v2_usersessions_dataforwarding" "example" {
  forwarding_name = "#name#"
  enabled = false
  matcher = "true"
  cloud_vendor_type = "aws"
  aws_connection {
    arn           = "arn:aws:iam::aws:role/#name#"
    connection_id = dynatrace_aws_connection.test-aws-connection.id
  }
  data_forwarding_type = "processed"
  pipelines = [dynatrace_openpipeline_v2_usersessions_pipelines.pipeline.id]
  bulk_pattern = "<YYYYMMDD>/<HH>/<HHmmss.SSSS>_<bulk-id>.json.gz"
}

resource "dynatrace_openpipeline_v2_usersessions_pipelines" "pipeline" {
  display_name = "Minimal pipeline"
  custom_id = "pipeline_Minimal_pipeline_1234_tf_#name#"
}

resource "dynatrace_aws_connection" "test-aws-connection" {
  name = "#name#"
  role_based_auth {
    consumers = ["SVC:com.dynatrace.openpipeline"]
  }
}
