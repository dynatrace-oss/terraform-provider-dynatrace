variable "aws_region" {
  type        = string
  description = "AWS region the provider operates in (also used as deployment region for the extension)."
  default     = "eu-central-1"
}

variable "connection_name" {
  type        = string
  description = "Display name of the Dynatrace AWS connection (HAS)."
  default     = "dac-tf-poc-connection"
}

variable "role_name" {
  type        = string
  description = "Name of the IAM role Dynatrace will assume."
  default     = "DynatraceMonitoringRole"
}

variable "monitoring_name" {
  type        = string
  description = "Display name of the Dynatrace AWS monitoring configuration."
  default     = "dac-tf-poc-monitoring"
}

variable "extension_version" {
  type        = string
  description = "Version of com.dynatrace.extension.da-aws to target. Leave null to let the provider resolve the latest installed version on the tenant (mirrors `dtctl create aws`)."
  default     = null
}

variable "regions" {
  type        = list(string)
  description = "AWS regions to monitor."
  default     = ["eu-central-1"]
}

variable "feature_sets" {
  type        = set(string)
  description = <<-EOT
    Extension feature sets to enable. Each AWS namespace shipped with
    com.dynatrace.extension.da-aws typically exposes two feature sets:
      *_essential     — production-grade default metric set
      *_autodiscovery — entity discovery metrics (lightweight)

    The list below mirrors the static catalogue shipped with the extension
    v1.0.x. Use 'dtctl get aws monitoring-feature-sets' against the live
    tenant to enumerate the values supported by your installed version.

    Full catalogue (uncomment what you need):
      AWS_MWAA_autodiscovery, AmazonMQ_autodiscovery, AmazonMQ_essential,
      Amazon_MWAA_autodiscovery, ApiGateway_autodiscovery, ApiGateway_essential,
      AppRunner_autodiscovery, AppRunner_essential, AppStream_autodiscovery,
      AppStream_essential, AppSync_autodiscovery, AppSync_essential,
      ApplicationELB_autodiscovery, ApplicationELB_essential,
      Athena_autodiscovery, Athena_essential, AutoScaling_autodiscovery,
      AutoScaling_essential, Backup_autodiscovery, Backup_essential,
      Bedrock_AgentAlias_autodiscovery, Bedrock_AgentAlias_essential,
      Bedrock_Guardrail_essential, Bedrock_Guardrails_autodiscovery,
      Cassandra_autodiscovery, Cassandra_essential,
      CertificateManager_autodiscovery, CertificateManager_essential,
      CloudFront_autodiscovery, CloudFront_essential, CloudHSM_autodiscovery,
      CloudHSM_essential, CloudTrail_autodiscovery, CloudTrail_essential,
      CodeBuild_autodiscovery, CodeBuild_essential, Cognito_autodiscovery,
      Cognito_essential, Connect_autodiscovery, Connect_essential,
      ContainerInsights_autodiscovery, ContainerInsights_essential,
      DAX_autodiscovery, DAX_essential, DMS_autodiscovery, DMS_essential,
      DataSync_autodiscovery, DataSync_essential, DirectConnect_autodiscovery,
      DirectConnect_essential, DocDB_autodiscovery, DocDB_essential,
      DynamoDB_autodiscovery, DynamoDB_essential, EBS_autodiscovery,
      EBS_essential, EC2_autodiscovery, EC2_essential, ECR_autodiscovery,
      ECR_essential, ECS_ContainerInsights_autodiscovery,
      ECS_ContainerInsights_essential, ECS_ManagedScaling_autodiscovery,
      ECS_autodiscovery, ECS_essential, EFS_autodiscovery, EFS_essential,
      EKS_autodiscovery, EKS_essential, ELB_autodiscovery, ELB_essential,
      EMR_EC2_autodiscovery, EMR_EC2_essential, EMR_Serverless_autodiscovery,
      EMR_Serverless_essential, ElastiCache_autodiscovery,
      ElastiCache_essential, ElasticBeanstalk_autodiscovery,
      ElasticBeanstalk_essential, Events_autodiscovery, Events_essential,
      FSx_autodiscovery, FSx_essential, Firehose_autodiscovery,
      Firehose_essential, GatewayELB_autodiscovery,
      GlobalAccelerator_autodiscovery, GlobalAccelerator_essential,
      Glue_autodiscovery, Glue_essential, IPAM_essential, KMS_essential,
      Kafka_Connect_autodiscovery, Kafka_Connect_essential, Kafka_autodiscovery,
      Kafka_essential, KinesisAnalytics_ApacheFlink_autodiscovery,
      KinesisAnalytics_ApacheFlink_essential, KinesisDataStreams_autodiscovery,
      KinesisDataStreams_essential, Lambda_autodiscovery, Lambda_essential,
      Logs_autodiscovery, Logs_essential, MWAA_autodiscovery,
      NATGateway_essential, NatGateway_autodiscovery, Neptune_autodiscovery,
      Neptune_essential, NetworkELB_autodiscovery, NetworkELB_essential,
      NetworkFirewall_autodiscovery, NetworkFirewall_essential,
      OpenSearch_Domain_autodiscovery, OpenSearch_Domain_essential,
      OpenSearch_Serverless_autodiscovery, OpenSearch_Serverless_essential,
      PrivateCA_autodiscovery, PrivateCA_essential,
      PrivateLinkEndpoints_autodiscovery, PrivateLinkEndpoints_essential,
      PrivateLinkServices_autodiscovery, PrivateLinkServices_essential,
      RDS_autodiscovery, RDS_essential, Redshift_autodiscovery,
      Redshift_essential, Route53_autodiscovery, Route53_essential,
      S3_autodiscovery, S3_essential, SNS_Topic_autodiscovery, SNS_essential,
      SQS_autodiscovery, SQS_essential, SageMakerEndpoint_autodiscovery,
      SageMakerEndpoint_essential, SageMakerInferenceComponent_autodiscovery,
      SageMakerPipeline_autodiscovery, SageMaker_autodiscovery,
      SageMaker_essential, StepFunctions_autodiscovery,
      StepFunctions_essential, StorageGateway_autodiscovery,
      StorageGateway_essential, TransitGateway_autodiscovery,
      TransitGateway_essential, VPN_SiteToSiteVPNConnection_autodiscovery,
      VPN_SiteToSiteVPNConnection_essential, WAFV2_essential, WAFv2_autodiscovery
  EOT
  default = [
    "EC2_essential",
    "RDS_essential",
    "S3_essential",
  ]
}

# ---------------------------------------------------------------------------
# Other monitoring-configuration attributes already implemented in the
# provider (see provider/dynatrace/api/extensions/dac/awsmonitoring/settings).
# Wire them through additional `variable` blocks if you need to override the
# defaults; the resource block in main.tf would then consume them.
# ---------------------------------------------------------------------------

# variable "monitoring_enabled" {
#   type        = bool
#   description = "Whether the monitoring configuration is active. Defaults to true."
#   default     = true
# }

# variable "deployment_region" {
#   type        = string
#   description = <<-EOT
#     AWS region the extension workload runs in. Defaults to the first entry
#     in var.regions when omitted (the provider computes it).
#     Set this explicitly only when the workload should live in a different
#     region than the first monitored region.
#   EOT
#   default     = null
# }

# variable "scope" {
#   type        = string
#   description = <<-EOT
#     Settings 2.0 scope for the monitoring configuration. The only supported
#     value today is "integration-aws" (default). Changing it forces
#     recreation.
#   EOT
#   default     = "integration-aws"
# }

# All of the following used to be served by a single `value_overrides_json`
# escape hatch. They are now first-class typed attributes on the resource —
# see main.tf for end-to-end examples of every block.
#
# Uncomment whichever variables you want to drive from outside, then wire
# them into the `dynatrace_aws_monitoring_configuration "this"` block in
# main.tf (e.g. `tag_enrichment = var.tag_enrichment`, or
# `dynamic "tag_filter" { for_each = var.tag_filters; content { … } }`).

# variable "activation_context" {
#   type        = string
#   description = "Extension activation context. Defaults to DATA_ACQUISITION."
#   default     = "DATA_ACQUISITION"
# }

# variable "deployment_scope" {
#   type        = string
#   description = "Deployment scope. SINGLE_ACCOUNT (default) or MULTI_ACCOUNT."
#   default     = "SINGLE_ACCOUNT"
# }

# variable "deployment_mode" {
#   type        = string
#   description = "Deployment mode. AUTOMATED (default) or MANUAL."
#   default     = "AUTOMATED"
# }

# variable "configuration_mode" {
#   type        = string
#   description = "Configuration mode. QUICK_START (default) or ADVANCED."
#   default     = "QUICK_START"
# }

# variable "smartscape_enabled" {
#   type        = bool
#   description = "Whether Smartscape topology mapping is enabled. Defaults to true."
#   default     = true
# }

# variable "tag_enrichment" {
#   type        = set(string)
#   description = "AWS tag keys whose values are promoted as Dynatrace tags on monitored entities."
#   default     = []
#   # Example: ["owner", "cost-center", "environment"]
# }

# variable "tag_filters" {
#   type = list(object({
#     key       = string
#     value     = string
#     condition = string # INCLUDE | EXCLUDE
#   }))
#   description = "Filter monitored AWS resources by tag. Use with a dynamic \"tag_filter\" block in main.tf."
#   default     = []
#   # Example:
#   # default = [
#   #   { key = "environment", value = "production", condition = "INCLUDE" },
#   #   { key = "owner",       value = "temp-team",  condition = "EXCLUDE" },
#   # ]
# }

# variable "cloud_watch_logs" {
#   type = object({
#     enabled = bool
#     regions = list(string)
#   })
#   description = "CloudWatch Logs ingestion. Pass null (default) to skip log ingestion entirely."
#   default     = null
#   # Example:
#   # default = {
#   #   enabled = true
#   #   regions = ["eu-central-1"]
#   # }
# }

# variable "dt_label_enrichments" {
#   type = list(object({
#     label   = string
#     literal = optional(string)
#     tag_key = optional(string)
#   }))
#   description = "Dynatrace labels (dt.*) added to every monitored entity. Each entry must set exactly one of literal or tag_key."
#   default     = []
#   # Example:
#   # default = [
#   #   { label = "dt.security_context", literal = "my-app" },
#   #   { label = "dt.cost.product",     tag_key = "product" },
#   # ]
# }

# variable "custom_namespaces" {
#   type = list(object({
#     namespace              = string
#     auto_discovery_enabled = optional(bool, false)
#     metrics = list(object({
#       name         = string
#       unit         = string
#       dimensions   = optional(list(string), [])
#       aggregations = list(string)
#       type         = string # CUSTOM_AWS | CUSTOM
#     }))
#   }))
#   description = "Additional CloudWatch namespaces to ingest (standard AWS/* via CUSTOM_AWS, or your own app metrics via CUSTOM)."
#   default     = []
#   # Example:
#   # default = [
#   #   {
#   #     namespace = "AWS/GroundStation"
#   #     metrics = [
#   #       {
#   #         name         = "AzimuthAngle"
#   #         unit         = "Count"
#   #         dimensions   = ["SatelliteId"]
#   #         aggregations = ["Sum", "SampleCount"]
#   #         type         = "CUSTOM_AWS"
#   #       },
#   #     ]
#   #   },
#   #   {
#   #     namespace = "MyApp/Metrics"
#   #     metrics = [
#   #       {
#   #         name         = "queue.depth"
#   #         unit         = "Count"
#   #         dimensions   = ["QueueName"]
#   #         aggregations = ["Average", "Maximum"]
#   #         type         = "CUSTOM"
#   #       },
#   #     ]
#   #   },
#   # ]
# }

variable "dynatrace_aws_principal_account_id" {
  type        = string
  description = "AWS account id of the Dynatrace tenant region that will assume the IAM role. Leave null to auto-detect from dt_env_url (314146291599 for production *.dynatrace.com tenants, 476114158034 for sprint/dev/labs)."
  default     = null
}

variable "dt_env_url" {
  type        = string
  description = "Dynatrace environment URL (e.g. https://abc12345.live.dynatrace.com or https://abc12345.apps.dynatrace.com). Required — no default."
}

variable "dynatrace_platform_token" {
  type        = string
  description = "Dynatrace platform token (dt0s16.*) used to authenticate to the tenant. Mark sensitive; pass via TF_VAR_dynatrace_platform_token or the DT_PLATFORM_TOKEN env var on the provider."
  sensitive   = true
  default     = null
}
