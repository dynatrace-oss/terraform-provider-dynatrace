# Repository instructions for GitHub Copilot

These instructions apply to **every** Copilot interaction in this repository
(chat, code completion). They encode the hard contracts the
`dynatrace_aws_monitoring_configuration` resource must keep.

---

## Project stack (do not change without an issue)

- Go `1.26.3`
- `github.com/hashicorp/terraform-plugin-sdk/v2 v2.40.1` — **SDK v2, not Plugin Framework**
- `github.com/Masterminds/semver/v3` for version resolution
- Generic resource wiring through `resources.NewGeneric(...)` (see existing entries in `provider/provider.go`)

## Hard rules for `provider/dynatrace/api/extensions/dac/awsmonitoring/`

These are non-negotiable. A PR that violates any of them must be revised.

1. **No JSON escape hatch.** Every wire field in
   `builtin:com.dynatrace.extension.da-aws` must be a first-class typed
   attribute in the SDK v2 schema. Forbidden constructs:
   - `value_overrides_json`, `extra_json`, or any free-form blob
   - `schema.TypeMap` with `Elem: TypeString` used as a generic bag
   - `json.RawMessage`, `map[string]any`, `map[string]interface{}` as
     exported struct fields
   - `interface{}` / `any` in marshalled payloads
2. **One sub-type per file.** Mirror the upstream convention used by
   `dynatrace/api/builtin/alerting/profile/settings/event_filter.go`. Each
   sub-type owns its struct, `Schema()`, `MarshalHCL`, `UnmarshalHCL`.
   Top-level `Settings` additionally owns `MarshalJSON` / `UnmarshalJSON`.
3. **Enums are validated.** Use `validation.StringInSlice` (or an inline
   `ValidateFunc`) listing the values **taken from the schema delta**,
   never invented.
4. **Optional/Required decisions are sourced.** When introducing or
   modifying an attribute, add a one-line `//` comment with the JSON
   Pointer in the extension schema, e.g.
   `// schema: /properties/aws/properties/foo (required, enum)`.
5. **XOR on repeatable blocks** is enforced in `UnmarshalHCL`, not
   declaratively. `ExactlyOneOf` with a wildcard path is invalid in SDK v2
   for `MaxItems != 1` blocks (see `dt_label_enrichment.go`).
6. **Lists vs sets**:
   - `TypeList` when order matters to the API (e.g. `regions`)
   - `TypeSet` when it does not (e.g. `feature_sets`, `tag_enrichment`)
   Decision must be visible in a comment if non-obvious.
7. **Round-trip tests are mandatory.** For every new/changed attribute:
   - HCL round-trip in `settings_test.go` (`MarshalHCL` then `UnmarshalHCL`)
   - JSON round-trip with a wire-shape literal pinned in the test source
8. **No TODO/WIP/FIXME** in committed code. No `panic("TODO")`. No `.bak`,
   `.orig`, or commented-out blocks.
9. **`extension_version` is `Optional + Computed`.** Never make it
   `Required`. Resolution rules live in `service.go::resolveLatestExtensionVersion`.
10. **`deployment_region` is `Optional + Computed`**, defaults to first
    `regions` entry. **`scope` is `Optional + ForceNew`**, defaults to
    `"integration-aws"`.
11. **Drift guard for echo fields.** Server echoes
    `cloudWatchLogsConfiguration: {enabled:false, regions:[]}` even when
    the user never wrote the block. `UnmarshalJSON` only populates
    `me.CloudWatchLogs` when `enabled || len(regions) > 0`. Apply the same
    pattern to any future echo-only-default field.

## Line endings

Files must be LF-only. CRLF breaks `goimports -l` in CI. Verify with
`file <path>` — must report `ASCII text`, not `ASCII text, with CRLF line
terminators`.

## Out of scope

- Plugin Framework migration
- Azure / GCP monitoring configuration resources (separate milestones)
- Per-region log/event StackSet ingest

If a request would require touching any of the above, stop and ask the
maintainer instead of guessing.
