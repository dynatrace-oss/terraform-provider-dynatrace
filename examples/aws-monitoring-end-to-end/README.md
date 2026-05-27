# AWS monitoring — end-to-end example

Creates a full Dynatrace AWS monitoring setup with a single
`terraform apply`:

1. **`dynatrace_aws_connection`** — the HAS object. Its `id` becomes the
   `sts:ExternalId` the IAM trust policy expects (mirrors Step 1 of the
   official CloudFormation flow: HAS_CONFIGURATION with empty role ARN).
2. **`aws_iam_role` + two `aws_iam_policy` attachments** — the IAM role
   that Dynatrace assumes to read CloudWatch metrics. The full
   CFN-equivalent permission set (~8.1 KB compact) exceeds AWS's 6144 B
   managed-policy quota, so it is split across
   [`iam_policy_topology.json.tftpl`](iam_policy_topology.json.tftpl)
   (primary TopologyInventory) and
   [`iam_policy_extras.json.tftpl`](iam_policy_extras.json.tftpl)
   (CloudWatch + Cassandra + secondary TopologyInventory).
3. **`dynatrace_aws_connection_role_arn`** — patches the role ARN back
   into the HAS object (Step 2 of the CFN flow).
4. **`dynatrace_aws_monitoring_configuration`** — this fork's new
   resource. Posts the typed monitoring configuration to
   `/api/v2/extensions/com.dynatrace.extension.da-aws/monitoringConfigurations`.

The order matters: the HAS connection must exist before the IAM trust
policy, the IAM role must exist before the role-ARN patch, and the role
ARN must be patched in before the monitoring configuration is created.
`depends_on` makes step 4 wait on step 3 explicitly.

## Prerequisites

- Terraform `>= 1.10`.
- AWS credentials with permission to create IAM roles in the target
  account (e.g. `AWS_PROFILE` env var pointing to an admin profile).
- A Dynatrace tenant URL and a **platform token** (`dt0s16.*`) with
  scopes covering `extensions.read/write`, `settings.read/write`,
  `hub.read`. Create one under *Access tokens → Platform tokens*.
- Either the fork installed via `dev_overrides` (see
  [`../../USING_THE_FORK.md`](../../USING_THE_FORK.md)) or the upstream
  release once this resource lands there.

## Running

```bash
export TF_VAR_dt_env_url="https://<tenant>.apps.dynatrace.com"
export TF_VAR_dynatrace_platform_token="dt0s16.XXXXXXXX...."
export AWS_PROFILE=<sandbox-profile>

terraform init
terraform plan
terraform apply
```

Expected first-apply diff: **8 to add** — 2 IAM policies, 1 IAM role,
2 attachments, 1 connection, 1 role-ARN patch, 1 monitoring configuration.

After apply, a second `terraform plan` **must report no changes** — that
is the round-trip proof that `MarshalJSON` and `UnmarshalJSON` agree
with the live API echo.

## Production vs non-production Dynatrace AWS account

The IAM trust policy needs to allow the right Dynatrace-owned AWS account
to assume the role. The example auto-detects from the tenant URL:

| `dt_env_url` matches | Dynatrace AWS account |
|---|---|
| `*.dynatrace.com` (production: `live` or `apps`) | `314146291599` |
| anything else (sprint / dev / labs) | `476114158034` |

These IDs come from the public CloudFormation templates at
`https://dynatrace-data-acquisition.s3.amazonaws.com/aws/deployment/cfn/latest/`.
Override manually with `TF_VAR_dynatrace_aws_principal_account_id` if the
auto-detection ever misses a new region split.

## Files

- [`main.tf`](main.tf) — providers, locals, all four resources, dependency
  wiring.
- [`variables.tf`](variables.tf) — every input the example accepts. Every
  optional typed attribute on `dynatrace_aws_monitoring_configuration` is
  documented as a commented-out `variable` block — uncomment to expose it
  to your wrapper.
- [`outputs.tf`](outputs.tf) — useful IDs / ARNs for downstream pipelines.
- [`iam_policy_topology.json.tftpl`](iam_policy_topology.json.tftpl) /
  [`iam_policy_extras.json.tftpl`](iam_policy_extras.json.tftpl) — IAM
  managed-policy bodies, split because of the 6144 B per-policy quota.

## Cleaning up

```bash
terraform destroy
```

If you previously created a monitoring configuration for the same AWS
account ID through the Dynatrace UI or `dtctl`, the API will reject the
new one with `400 Account ID must be unique`. Wipe the orphan first — see
the runbook in the spec referenced by [`../../USING_THE_FORK.md`](../../USING_THE_FORK.md).
