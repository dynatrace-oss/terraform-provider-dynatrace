---
page_title: "dynatrace_json_dashboard_base Resource - terraform-provider-dynatrace"
subcategory: "Dashboards"
description: |-
  The resource `dynatrace_json_dashboard_base` is of use when `dynatrace_json_dashboard` resources within the same module refer to each other. You're able to avoid cyclic dependencies here.
---

# dynatrace_json_dashboard_base (Resource)

The resource `dynatrace_json_dashboard_base` doesn't contain any attributes itself. It acts as some sort of anchor resource that defines the eventual `ID` of a Dashboard - without having to
refer to the resource `dynatrace_json_dashboard` explicitly.

## Resource Example Usage

This example shows how to apply two Dashboards that refer to each other using Markdown Tiles.

The blunt approach would be to create two `dynatrace_json_dashboard` resource blocks.

```terraform
resource "dynatrace_json_dashboard" "dashboard-a" {
  contents = jsonencode({
      "dashboardMetadata": {
        "name": "dashboard-a",
        "owner": "me@home.com"
      },
      "tiles": [
        {
          "bounds": {
            "height": 152,
            "left": 0,
            "top": 0,
            "width": 304
          },
          "configured": true,
          "markdown": "## This is a reference to [Dashboard B](https://#########.live.dynatrace.com/#dashboard;gtf=-2h;gf=all;id=${dynatrace_json_dashboard.dashboard-b.id})",
          "name": "Markdown",
          "tileType": "MARKDOWN"
        }
      ]
    })
}

resource "dynatrace_json_dashboard" "dashboard-b" {
  contents = jsonencode({
      "dashboardMetadata": {
        "name": "dashboard-b",
        "owner": "me@home.com"
      },
      "tiles": [
        {
          "bounds": {
            "height": 152,
            "left": 0,
            "top": 0,
            "width": 304
          },
          "configured": true,
          "markdown": "## This is a reference to [Dashboard A](https://#########.live.dynatrace.com/#dashboard;gtf=-2h;gf=all;id=${dynatrace_json_dashboard.dashboard-a.id})",
          "name": "Markdown",
          "tileType": "MARKDOWN"
        }
      ]
    }) 
}
```

Terraform will detect a cyclic dependency here and will refuse to apply it.

The resource `dynatrace_json_dashboard_base` allows us to resolve that cycle.
Purpose of the resource block `dynatrace_json_dashboard_base` is solely to allow other resource to refer to its `id` attribute. It essentially represents the
Dashboard without specifying what tiles it consists of.

```terraform
resource "dynatrace_json_dashboard_base" "dashboard-a" {
}
```
That's why it is important that any resource block `dynatrace_json_dashboard_base` is accompanied by  a resource block `dynatrace_json_dashboard` that contributes
the actual contents of that Dashboard.

The attribute `link_id` tells Terraform that the `ID` of that Dashboard needs to align with the `ID` defined by the
`dynatrace_json_dashboard_base` resource block.

```terraform
resource "dynatrace_json_dashboard" "dashboard-a" {
  contents = jsonencode({
  ...
  })
  link_id  = "${dynatrace_json_dashboard_base.dashboard-a.id}"
}
```
If another Dashboards needs to reference `dashboard-a`, it can do so without referencing another `dynatrace_json_dashboard` resource. It can refer to the `dynatrace_json_dashboard_base`
resource.
```terraform
resource "dynatrace_json_dashboard" "dashboard-b" {
  contents = jsonencode({
      ...
      "tiles": [
          ...
          "markdown": "## This is a reference to [Dashboard B](https://#########.live.dynatrace.com/#dashboard;gtf=-2h;gf=all;id=${dynatrace_json_dashboard_base.dashboard-a.id})",
          ....
      ]
    })
  ...
}
```
Here the full example with two dashboards cross referencing each other - utilizing the resource `dynatrace_json_dashboard_base`.
```terraform

resource "dynatrace_json_dashboard_base" "dashboard-a" {
}

resource "dynatrace_json_dashboard" "dashboard-a" {
  contents = jsonencode({
      "dashboardMetadata": {
        "name": "dashboard-a",
        "owner": "me@home.com"
      },
      "tiles": [
        {
          "bounds": {
            "height": 152,
            "left": 0,
            "top": 0,
            "width": 304
          },
          "configured": true,
          "markdown": "## This is a reference to [Dashboard B](https://#########.live.dynatrace.com/#dashboard;gtf=-2h;gf=all;id=${dynatrace_json_dashboard_base.dashboard-b.id})",
          "name": "Markdown",
          "tileType": "MARKDOWN"
        }
      ]
    })
  link_id  = "${dynatrace_json_dashboard_base.dashboard-a.id}"
}

resource "dynatrace_json_dashboard_base" "dashboard-b" {
}

resource "dynatrace_json_dashboard" "dashboard-b" {
  contents = jsonencode({
      "dashboardMetadata": {
        "name": "dashboard-b",
        "owner": "me@home.com"
      },
      "tiles": [
        {
          "bounds": {
            "height": 152,
            "left": 0,
            "top": 0,
            "width": 304
          },
          "configured": true,
          "markdown": "## This is a reference to [Dashboard A](https://#########.live.dynatrace.com/#dashboard;gtf=-2h;gf=all;id=${dynatrace_json_dashboard_base.dashboard-a.id})",
          "name": "Markdown",
          "tileType": "MARKDOWN"
        }
      ]
    }) 
    link_id  = "${dynatrace_json_dashboard_base.dashboard-b.id}"
}
```