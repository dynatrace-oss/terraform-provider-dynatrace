# Using this fork before the changes land in `dynatrace-oss/terraform-provider-dynatrace`

This fork adds the `dynatrace_aws_monitoring_configuration` resource — the
typed Terraform equivalent of `dtctl create aws` — on top of the upstream
provider. Until the corresponding PR is merged into
`dynatrace-oss/terraform-provider-dynatrace`, you can consume this fork
directly. Two supported paths below.

---

## Path A — `dev_overrides` (fastest; recommended for evaluation)

Use this when you want to try the resource against your own tenant
without going through a registry.

### 1. Build the provider locally

```bash
git clone https://github.com/pawelsiwek/terraform-provider-dynatrace.git
cd terraform-provider-dynatrace
git checkout feat/dac-aws-monitoring
go install .                       # binary lands in $(go env GOPATH)/bin
```

### 2. Point Terraform at the local binary

`~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "dynatrace-oss/dynatrace" = "/Users/<you>/go/bin"
  }
  direct {}
}
```

(`go env GOPATH` gives you the right path if it is not `~/go`.)

### 3. Use the upstream `source` in your HCL

Because of `dev_overrides`, Terraform resolves the upstream address to
your local fork binary. No `required_providers` change needed.

```hcl
terraform {
  required_providers {
    dynatrace = { source = "dynatrace-oss/dynatrace" }
    aws       = { source = "hashicorp/aws" }
  }
}
```

You will see `Warning: Provider development overrides are in effect` on
every command — expected.

### 4. End-to-end example

The full sample (HAS connection → IAM role with trust gated on the
connection's `objectId` as `sts:ExternalId` → role-ARN patch → monitoring
configuration) lives in the PoC repo at
[`aws-poc-sample/`](https://github.com/pawelsiwek/dac-terraform-provider-poc/tree/main/aws-poc-sample).
Copy the four files (`main.tf`, `variables.tf`,
`iam_policy_topology.json.tftpl`, `iam_policy_extras.json.tftpl`) into a
working directory.

Minimum env:

```bash
export DYNATRACE_ENV_URL="https://<tenant>.apps.dynatrace.com"
export DYNATRACE_API_TOKEN="dt0s16.XXXXXXXX...."  # platform token, scopes:
                                                  # extensions.read/write,
                                                  # settings.read/write,
                                                  # hub.read
export TF_VAR_dynatrace_platform_token="$DYNATRACE_API_TOKEN"
# plus standard AWS_PROFILE / AWS_REGION as usual

terraform init -upgrade
terraform apply -auto-approve
terraform plan                       # MUST be "No changes." — round-trip proof
```

---

## Path B — Consume as a separate registry provider

For shared modules / CI where `dev_overrides` is impractical. Published
under a distinct namespace so it never clashes with the upstream:

```hcl
terraform {
  required_providers {
    dynatrace-dac = {
      source  = "pawelsiwek/dynatrace-dac"
      version = "~> 0.1"
    }
    aws = { source = "hashicorp/aws" }
  }
}

provider "dynatrace-dac" {
  environment_url = var.dt_env_url
  platform_token  = var.dt_platform_token
}

resource "dynatrace-dac_aws_monitoring_configuration" "this" {
  # ... same schema as dynatrace_aws_monitoring_configuration
}
```

Only the **provider local name** differs (`dynatrace-dac`); resource
schemas are identical. When the upstream merge lands, migration is a
`source` swap + `terraform state replace-provider`.

> Status: registry publication of `pawelsiwek/dynatrace-dac` is planned
> but not yet live. Until it ships, use Path A.

---

## Resource quick reference

```hcl
resource "dynatrace_aws_monitoring_configuration" "this" {
  name          = "prod-aws-monitoring"
  connection_id = dynatrace_aws_connection.this.id
  account_id    = data.aws_caller_identity.current.account_id
  regions       = ["eu-central-1", "us-east-1"]

  # All optional; sensible defaults applied if omitted.
  # extension_version    = "1.0.7"           # else: latest installed
  # deployment_region    = "eu-central-1"    # else: regions[0]
  # deployment_scope     = "SINGLE_ACCOUNT"
  # deployment_mode      = "AUTOMATED"
  # configuration_mode   = "QUICK_START"
  # activation_context   = "DATA_ACQUISITION"
  # smartscape_enabled   = true
  # enabled              = true
  # scope                = "integration-aws" # ForceNew
  # feature_sets         = ["…enum values…"]
  # tag_enrichment       = ["Environment", "Owner"]

  # tag_filter { key = "Env" value = "prod" condition = "INCLUDE" }
  # cloud_watch_logs { enabled = true  regions = ["eu-central-1"] }
  # dt_label_enrichment { label = "team" tag_key = "Team" }   # XOR: literal | tag_key
  # custom_namespace { … nested metric { … } … }
}
```

Full attribute table: [`spec/task.md` §12](https://github.com/pawelsiwek/dac-terraform-provider-poc/blob/main/spec/task.md#12-capability-coverage-dynatrace_aws_monitoring_configuration-vs-dtctl-create-aws).

---

## Migration back to upstream (when the PR merges)

1. Bump `required_providers.dynatrace.version` to whichever upstream
   release contains `dynatrace_aws_monitoring_configuration`.
2. Drop the `dev_overrides` block from `~/.terraformrc` (Path A) **or**
   change `source` from `pawelsiwek/dynatrace-dac` to
   `dynatrace-oss/dynatrace` and run
   `terraform state replace-provider pawelsiwek/dynatrace-dac dynatrace-oss/dynatrace`
   (Path B).
3. `terraform init -upgrade && terraform plan` — must be "No changes."

If it is **not** no-op, file an issue in this fork before applying — the
schemas were intended to be byte-equivalent.

---

## Known gotchas (from operator runbook)

- **`400 Account ID must be unique`** on create: the DAC API enforces one
  monitoring config per AWS account ID per tenant (regions do not
  segment). Cleanup snippet:
  [`spec/task.md` §14.4](https://github.com/pawelsiwek/dac-terraform-provider-poc/blob/main/spec/task.md#144-clean-leftover-monitoring-config-for-the-same-aws-account).
- **`Cannot exceed quota for PolicySize: 6144`**: the full CFN-equivalent
  IAM policy (~8.1 KB) does not fit in a single `aws_iam_policy`. The
  sample splits it into `iam_policy_topology.json.tftpl` +
  `iam_policy_extras.json.tftpl`, both attached via
  `aws_iam_role_policy_attachment`.
- **`No API Token has been specified`** on `dynatrace_aws_connection`:
  set `DYNATRACE_API_TOKEN` in the same shell as `terraform apply` —
  `TF_VAR_*` alone does not propagate.
- **Plan drift on `cloud_watch_logs`**: fixed in the fork
  ([`settings.go::UnmarshalJSON`](provider/dynatrace/api/extensions/dac/awsmonitoring/settings/settings.go))
  by only populating the block when `enabled || len(regions) > 0`.
- **CRLF**: do not edit provider sources on Windows without `core.autocrlf=input`;
  CI's `goimports -l` will flag them.

---

## Where the new resource lives in this fork

- Resource code: `provider/dynatrace/api/extensions/dac/awsmonitoring/`
- Registration: `provider/provider.go` (next to `dynatrace_aws_connection`)
- Tests: `provider/dynatrace/api/extensions/dac/awsmonitoring/settings/settings_test.go`
