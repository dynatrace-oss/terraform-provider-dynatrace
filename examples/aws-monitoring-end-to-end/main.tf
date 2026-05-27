terraform {
  required_version = ">= 1.10"

  required_providers {
    dynatrace = {
      source  = "dynatrace-oss/dynatrace"
      version = ">= 1.0.0"
    }
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

provider "dynatrace" {
  dt_env_url     = var.dt_env_url
  platform_token = var.dynatrace_platform_token
  dt_api_token   = var.dynatrace_platform_token
}

data "aws_caller_identity" "current" {}

locals {
  # Mirrors the CFN cProductionEnvironment condition in
  # da-aws-nested-monitoring-role.yaml: production tenants live under
  # *.dynatrace.com, everything else (sprint / dev / labs) is non-production.
  is_production_tenant = can(regex("\\.dynatrace\\.com(/|$)", var.dt_env_url))

  dynatrace_aws_principal_account_id_effective = coalesce(
    var.dynatrace_aws_principal_account_id,
    local.is_production_tenant ? "314146291599" : "476114158034",
  )
}

# 1. Create the Dynatrace AWS connection first. Its objectId is what the IAM
#    trust policy must accept as sts:ExternalId. role_arn is intentionally
#    omitted here and patched in later by dynatrace_aws_connection_role_arn,
#    mirroring Isaac's CFN flow (Step 1: HAS_CONFIGURATION with empty arn).
resource "dynatrace_aws_connection" "this" {
  name = var.connection_name
  role_based_auth {
    consumers = ["SVC:com.dynatrace.da"]
  }
}

# 2. Create the IAM role in the target AWS account. Trust policy allows the
#    Dynatrace principal to assume the role, gated on the connection's
#    objectId as external id.
data "aws_iam_policy_document" "trust" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::${local.dynatrace_aws_principal_account_id_effective}:root"]
    }
    condition {
      test     = "StringEquals"
      variable = "sts:ExternalId"
      values   = [dynatrace_aws_connection.this.id]
    }
  }
}

resource "aws_iam_policy" "monitoring_topology" {
  name        = "${var.role_name}-topology"
  description = "Dynatrace DAC AWS monitoring — primary TopologyInventory statement (split out of the full policy because AWS caps managed policies at 6144 chars compact)."
  policy = templatefile("${path.module}/iam_policy_topology.json.tftpl", {
    aws_account_id = data.aws_caller_identity.current.account_id
  })
}

resource "aws_iam_policy" "monitoring_extras" {
  name        = "${var.role_name}-extras"
  description = "Dynatrace DAC AWS monitoring — CloudWatch + Cassandra + secondary TopologyInventory statements."
  policy = templatefile("${path.module}/iam_policy_extras.json.tftpl", {
    aws_account_id = data.aws_caller_identity.current.account_id
  })
}

resource "aws_iam_role" "monitoring" {
  name               = var.role_name
  description        = "Assumed by Dynatrace to read CloudWatch metrics for the DAC AWS integration."
  assume_role_policy = data.aws_iam_policy_document.trust.json
}

resource "aws_iam_role_policy_attachment" "monitoring_topology" {
  role       = aws_iam_role.monitoring.name
  policy_arn = aws_iam_policy.monitoring_topology.arn
}

resource "aws_iam_role_policy_attachment" "monitoring_extras" {
  role       = aws_iam_role.monitoring.name
  policy_arn = aws_iam_policy.monitoring_extras.arn
}

# 3. Patch the role ARN into the connection so DAC can assume it.
resource "dynatrace_aws_connection_role_arn" "this" {
  aws_connection_id = dynatrace_aws_connection.this.id
  role_arn          = aws_iam_role.monitoring.arn
}

# 4. Create the monitoring configuration. Depends on the role ARN being
#    patched into the connection.
resource "dynatrace_aws_monitoring_configuration" "this" {
  name    = var.monitoring_name
  enabled = true
  # extension_version is intentionally optional: when omitted, the provider
  # resolves the highest semver version installed on the tenant (same as
  # `dtctl create aws`). Pin it here only when you need a specific version.
  extension_version = var.extension_version

  connection_id = dynatrace_aws_connection.this.id
  account_id    = data.aws_caller_identity.current.account_id

  regions      = var.regions
  feature_sets = var.feature_sets

  # ---------------------------------------------------------------------------
  # Optional first-class typed attributes. Uncomment what you need.
  # All of these used to live behind a `value_overrides_json` escape hatch —
  # they are now individually typed, validated, and tracked in plan diffs.
  # ---------------------------------------------------------------------------

  # smartscape_enabled  = true            # default true
  # activation_context  = "DATA_ACQUISITION"
  # deployment_scope    = "SINGLE_ACCOUNT"
  # deployment_mode     = "AUTOMATED"
  # configuration_mode  = "QUICK_START"

  # Promote selected AWS tag keys onto Dynatrace entities.
  # tag_enrichment = ["owner", "cost-center", "environment"]

  # Filter monitored resources by AWS tags.
  # tag_filter {
  #   key       = "environment"
  #   value     = "production"
  #   condition = "INCLUDE"
  # }
  # tag_filter {
  #   key       = "owner"
  #   value     = "temp-team"
  #   condition = "EXCLUDE"
  # }

  # Add static or tag-derived Dynatrace labels (dt.*) to every monitored entity.
  # dt_label_enrichment {
  #   label   = "dt.security_context"
  #   literal = "my-app"
  # }
  # dt_label_enrichment {
  #   label   = "dt.cost.product"
  #   tag_key = "product"
  # }

  # Ingest CloudWatch logs in addition to metrics.
  # cloud_watch_logs {
  #   enabled = true
  #   regions = ["eu-central-1"]
  # }

  # Custom CloudWatch namespace from a standard AWS service (CUSTOM_AWS).
  # custom_namespace {
  #   namespace              = "AWS/GroundStation"
  #   auto_discovery_enabled = false
  #   metric {
  #     name         = "AzimuthAngle"
  #     unit         = "Count"
  #     dimensions   = ["SatelliteId"]
  #     aggregations = ["Sum", "SampleCount"]
  #     type         = "CUSTOM_AWS"
  #   }
  # }

  # Custom non-AWS namespace published by your own application (CUSTOM).
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

  depends_on = [dynatrace_aws_connection_role_arn.this]
}
