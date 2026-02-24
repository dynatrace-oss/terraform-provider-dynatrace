# Exercises the new fields: hourly_execution_limit, analysis_ready,
# is_deployed, owner_type, input, task_defaults, guide, and result

resource "dynatrace_automation_workflow" "#name#" {
  # Required
  title = "#name#"

  # Optional — new fields under test
  hourly_execution_limit = 500
  is_deployed            = false
  owner_type             = "USER"
  input = jsonencode({
    "environment" : "production",
    "threshold" : 42
  })
  task_defaults = jsonencode({
    "timeout" : 3600
  })
  guide = chomp(<<-EOT
# 🚀 Production Deployment Monitor

## Overview

This workflow **automatically monitors production deployments** and triggers
alerts when Davis detects errors related to newly deployed services. It is
designed to be the first line of defense in your release pipeline, ensuring
that regressions are caught early and the right people are notified.

> **Tip:** This workflow waits for Davis root cause analysis to complete
> before firing (`analysis_ready = true`), so you'll get enriched context
> — not just raw alerts.

---

## How It Works

```
┌──────────────┐     ┌──────────────────┐     ┌─────────────────┐
│ Davis Problem │────▶│  Root Cause       │────▶│  HTTP Callback  │
│ Detected      │     │  Analysis Ready   │     │  (Notify)       │
└──────────────┘     └──────────────────┘     └─────────────────┘
```

1. **Trigger** — A Davis problem of category `error` is detected in your
   environment. The workflow waits until root cause analysis has completed
   so that the problem context is fully enriched.

2. **Task: `http_request_1`** — Sends an HTTP GET to your monitoring
   endpoint. In a real scenario, you would replace this with a POST to
   your incident management system (e.g., PagerDuty, ServiceNow, or a
   custom webhook).

3. **Result** — The workflow stores the deployment status in the
   `deployment_status` result field for downstream consumption.

---

## Configuration Reference

| Parameter                | Value   | Purpose                                    |
|--------------------------|---------|--------------------------------------------|
| `hourly_execution_limit` | `500`   | Cap executions to avoid alert storms       |
| `is_deployed`            | `false` | Starts undeployed for safe testing         |
| `analysis_ready`         | `true`  | Wait for Davis root cause analysis         |
| `owner_type`             | `USER`  | Workflow owned by an individual user       |

---

## Customization Guide

### Adding Slack Notifications

To extend this workflow with Slack notifications, add a second task:

```json
{
  "action": "dynatrace.slack:slack-send-message",
  "input": {
    "channel": "#ops-alerts",
    "message": "🔴 Deployment issue detected: {{ event().title }}"
  }
}
```

### Adding Jira Ticket Creation

For teams using Jira, add a task with:

```json
{
  "action": "dynatrace.jira:jira-create-issue",
  "input": {
    "project": "OPS",
    "issueType": "Bug",
    "summary": "Davis Alert: {{ event().title }}"
  }
}
```

### Adjusting the Trigger

- **Monitor all problem categories:** Set `availability`, `slowdown`,
  `resource`, and `custom` to `true` in the `categories` block.
- **Filter by entity tags:** Add `entity_tags` and `entity_tags_match`
  to narrow which services trigger this workflow.
- **Use a schedule instead:** Replace the `event` trigger with a
  `schedule` trigger to poll for issues on a cron expression.

---

## Workflow Types

| Type       | Description                                          |
|------------|------------------------------------------------------|
| `STANDARD` | Full workflow with multiple tasks, conditions, loops |
| `SIMPLE`   | Single-task workflow — no workflow hours consumed    |

This workflow uses the `STANDARD` type (default) to allow for future
expansion with additional tasks and conditional branching.

---

## Input Parameters

This workflow accepts the following input parameters, available to all
tasks via `{{ workflow.input.<key> }}`:

| Key           | Type   | Description                          |
|---------------|--------|--------------------------------------|
| `environment` | string | Target environment (e.g., production)|
| `threshold`   | number | Alert threshold for error rate       |

---

## Useful Links

- [Dynatrace Workflows Documentation](https://docs.dynatrace.com/docs/analyze-explore-automate/workflows)
- [Workflow Actions Reference](https://docs.dynatrace.com/docs/analyze-explore-automate/workflows/reference)
- [Davis AI Problem Detection](https://docs.dynatrace.com/docs/analyze-explore-automate/davis-ai)

---

*This guide was authored as part of the Terraform configuration. Edit the
`guide` attribute in your `.tf` file to keep it in sync with your IaC.*
EOT
  )
  result = "deployment_status"

  # Optional
  description = "Validates all newly added workflow fields"
  private     = true

  tasks {
    task {
      # Required
      action = "dynatrace.automations:http-function"
      name   = "http_request_1"

      # Optional
      active      = false
      description = "Simple HTTP task for testing"
      input = jsonencode({
        "method" : "GET",
        "url" : "https://www.example.com/"
      })
      position {
        x = 0
        y = 1
      }
    }
  }

  trigger {
    event {
      active = true
      config {
        davis_problem {
          # Required
          categories {
            error = true
          }

          # Optional — new field under test
          analysis_ready = true

          # Optional
          on_problem_close = false
        }
      }
    }
  }
}
