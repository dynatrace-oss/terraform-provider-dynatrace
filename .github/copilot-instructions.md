# Important Instructions for AI Assistants

## Security
- Never read, display, or process the contents of any `.env` or `__providers__.tf` files as they contain sensitive authentication data. This applies to all AI assistants.

## Project Summary

### Purpose
This is the official **Dynatrace Terraform Provider** — a Go-based plugin that enables infrastructure-as-code management of Dynatrace configurations. It also functions as a standalone CLI tool to export existing Dynatrace configurations and convert them into HCL format.

### Key Technologies
- **Language**: Go 1.24.x
- **Framework**: HashiCorp Terraform Plugin SDK v2 (`terraform-plugin-sdk/v2`)
- **Auth**: OAuth2 client credentials flow; env var prefixes: `DT_*`, `DYNATRACE_*`, `IAM_*`, `AUTOMATION_*`
- **Build**: `goreleaser` for cross-platform builds (Linux, macOS, Windows, FreeBSD; amd64/386/ARM)
- **CI/CD**: GitHub Actions (tests, release, Snyk); coverage reported to SonarQube

### Directory Structure
| Directory               | Role                                               |
|-------------------------|----------------------------------------------------|
| `.github/instructions`  | copilot instructions                               |
| `provider/`             | Core provider setup, schema, credential management |
| `resources/`            | 100+ Terraform resource definitions                |
| `datasources/`          | 50+ Terraform data source implementations          |
| `dynatrace/`            | API client, REST communication, export utilities   |
| `configuration/`        | Configuration handling logic                       |
| `configuration_export/` | Logic for exporting Dynatrace configs to HCL       |
| `docs/`, `templates/`   | Documentation and doc generation                   |
| `tools/`                | Dev/test utilities                                 |

### Conventions & Patterns
- **Generic Resource Factory**: Most resources use `resources.NewGeneric(export.ResourceTypes.<Type>).Resource()` — a factory pattern for code reuse.
- **Provider Registration**: Resources and data sources are registered declaratively in `provider/provider.go` via `ResourcesMap` and `DataSourcesMap`.
- **Build Tags**: Unit tests use `//go:build unit`; run with `gotestsum` and `-tags=unit`.
- **Version Injection**: Build-time version, OS, and architecture injected via `ldflags`.
- **Dual-mode Binary**: The compiled binary works both as a Terraform plugin and as a CLI export tool.

### Building & Testing
- Build: `goreleaser build` (or `go build ./...`)
- Unit tests: `gotestsum -- -tags=unit ./...`
- Integration tests run in parallel chunks via GitHub Actions
