---
name: release-notes
description: Turn a changelog commit list into user-facing release notes in the provider's release style.
---

## What I do

- Convert a `## Changelog` section into polished release notes that match recent provider releases such as `v1.95.0`, `v1.92.0`, and `v1.91.0`.
- Fetch the changelog directly from GitHub when given a release tag, release URL, or compare range.
- Start from the commit list only, then inspect commits and related PR metadata as needed to infer the user-facing changes.
- Group related commits into a small set of meaningful bullets instead of writing one bullet per commit.

## Prerequisites

- Either:
  - A markdown file with a `## Changelog` section listing commit hashes and commit messages, or
  - A GitHub release URL, compare URL, or enough information to identify the repo and tag or range.
- The `gh` CLI must be installed and authenticated.

## Accepted inputs

This skill can work from any of these starting points:

- A local markdown file containing a `## Changelog` section.
- A GitHub release URL such as `https://github.com/<owner>/<repo>/releases/tag/<tag>`.
- A compare URL such as `https://github.com/<owner>/<repo>/compare/<base>...<head>`.
- A repo plus a tag, for example `dynatrace-oss/terraform-provider-dynatrace` and `v1.95.0`.
- A repo plus a tag range, for example `v1.94.0...v1.95.0`.

## Example invocations

Use this skill with prompts like these:

### Release URL input

```text
/release-notes https://github.com/dynatrace-oss/terraform-provider-dynatrace/releases/tag/v1.95.0
```

Expected behavior:

- Use `gh` to read the release body.
- Extract only the `## Changelog` section.
- Inspect relevant commits and PR metadata.
- Return release notes in the repository's usual grouped format.

### Compare range input

```text
/release-notes dynatrace-oss/terraform-provider-dynatrace v1.94.0...v1.95.0
```

Expected behavior:

- Use the compare API to get the commits in the range.
- Treat that commit list as the changelog source.
- Group the changes into user-facing release-note sections.

### Local markdown file input

```text
/release-notes docs/releases/v1.95.0.md
```

Expected behavior:

- Read the local markdown file.
- Extract the existing `## Changelog` section.
- Insert generated release notes above `## Changelog`.
- Keep the original changelog below the generated notes.

## Style guide

Use the release-note style from recent releases in this repository.

- Prefer `gh` for all GitHub access. Do not use generic web fetching when `gh` can provide the same information.

- Prefer these sections, in this order:

```markdown
## New resources
## New data sources
## New fields / breaking changes
## Deprecated resources
## Bug fixes
## Misc
```

- Omit any section that would be empty.
- Start each bullet with the affected resource or data source in backticks when possible.
- Use one concise bullet per user-facing change or tightly related group of changes.
- Focus on what changed for Terraform users, not on internal implementation details.
- Be explicit about migration impact. If something is breaking, start the bullet with `**BREAKING**`.
- Do not list raw version bumps. Inspect the underlying commit diff and describe the actual user-visible change.
- Ignore purely internal `test:`, `ci:`, `refactor:`, `chore:`, and docs-only commits unless they clearly change user behavior, exported output, generated documentation, logging, or troubleshooting.
- If several commits fix the same resource or the same class of issue, merge them into one bullet.
- Keep wording crisp and factual.

## Classification rules

- `New resources`: brand new Terraform resources.
- `New data sources`: brand new Terraform data sources.
- `New fields / breaking changes`: new schema fields, new validation rules, changed defaults, fields that now force recreation, renamed or removed fields or blocks, behavior changes caused by schema updates, and anything marked breaking.
- `Deprecated resources`: explicit deprecations and replacements.
- `Bug fixes`: plan/apply fixes, panic fixes, export fixes, empty-set fixes, retry logic, incorrect payloads, computed or default handling, and similar corrections.
- `Misc`: notable user-facing documentation, export, logging, or description updates that do not fit elsewhere.

## Workflow

Follow these steps to generate release notes:

### 1. Read the changelog entries and collect the commit hashes

- If a local markdown file is provided, read the existing `## Changelog` section and use the hashes from there.
- If a release URL or tag is provided, fetch the release body from GitHub and extract the `## Changelog` section from it.
- If a compare URL or tag range is provided, fetch the commit list directly from the compare API and treat that as the changelog source.
- For each potentially relevant commit, inspect the commit details with `gh`.
- Inspect the changed files and, if available, the linked PR to determine the actual user-facing change.

Helpful commands:

```bash
gh release view <tag> --repo <owner>/<repo> --json body
gh api repos/<owner>/<repo>/releases/tags/<tag>
gh api repos/<owner>/<repo>/compare/<base>...<head>
gh api repos/dynatrace-oss/terraform-provider-dynatrace/commits/<hash>
gh api repos/dynatrace-oss/terraform-provider-dynatrace/commits/<hash>/pulls
gh pr view <number> --json title,body
```

- Prefer the release body when it already contains a curated `## Changelog` section.
- Fall back to the compare API when there is no release body or no changelog section to extract.

### 2. Extract only the `## Changelog` section when starting from a release body

- Parse the markdown returned by `gh release view` or the release API.
- Find the heading whose exact text is `## Changelog`.
- Treat everything after that heading as changelog content until the next heading of the same or higher level, or the end of the release body.
- Ignore all earlier sections such as `## New resources`, `## Bug fixes`, or prose above the changelog.
- If the release body does not contain `## Changelog`, use the compare API instead of guessing.
- From the extracted section, keep only commit-list entries and ignore unrelated text.

### 3. Group commits into user-facing changes

- Combine commits that belong to the same feature, fix, or migration.
- Drop noise such as test-only follow-up commits unless they reveal the intended behavior of a user-facing change.
- Prefer a smaller number of high-signal bullets over exhaustive commit-by-commit coverage.

### 4. Write the release notes at the top of the file

- Insert the chosen sections above the existing `## Changelog` section.
- Write bullets in repository style, for example:
  - `` `dynatrace_example_resource`: Added new optional field `foo`. ``
  - `` `dynatrace_example_resource`: Fixed non-empty plan after apply when `bar` is omitted. ``
  - `` `dynatrace_example_resource`: Deprecated in favor of `dynatrace_other_resource`. ``
  - `` `dynatrace_example_resource`: New resource for managing ... ``
- When the change affects multiple resources in the same way, list them together in the same bullet.
- When a change comes from a generated resource update, translate the diff into the resulting fields, validations, defaults, or behavior changes instead of mentioning code generation.

### 5. Remove empty sections and clean up wording

- Delete any empty heading.
- Deduplicate overlapping bullets.
- Make sure each bullet is understandable without reading the changelog.

### 6. Save the updated markdown file

- If working in a local markdown file, keep the original `## Changelog` section below the generated notes.
- If working directly from GitHub input only, return the generated release notes as markdown unless the user also provided a target file to update.
