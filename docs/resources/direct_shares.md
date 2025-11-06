---
layout: ""
page_title: "Resource - terraform-provider-dynatrace"
subcategory: "Documents"
description: |-
  The resource `dynatrace_direct_shares` covers configuration of direct shares for Documents (dashboards and notebooks) in Dynatrace.
---

# dynatrace_direct_shares (Resource)

-> **Dynatrace SaaS only**

-> To utilize this resource, please define the environment variables `DT_CLIENT_ID`, `DT_CLIENT_SECRET`, `DT_ACCOUNT_ID` with an OAuth client including the following permissions: **Read direct-shares** (`document:direct-shares:read`), **Write direct-shares** (`document:direct-shares:write`), and **Delete direct-shares** (`document:direct-shares:delete`).

-> This resource is currently not covered by the export utility.

## Dynatrace Documentation

- Dynatrace Documents - https://########.apps.dynatrace.com/platform/swagger-ui/index.html?urls.primaryName=Document%20Service

## Resource Example Usage

```terraform
resource "dynatrace_direct_shares" "this" {
  document_id = dynatrace_document.this.id
  access      = "read-write"

  recipients {
    recipient {
      id   = "441664f0-23c9-40ef-b344-18c02c23d787"
      type = "user"
    }

    recipient {
      id   = "441664f0-23c9-40ef-b344-18c02c23d788"
      type = "group"
    }
  }
}

resource "dynatrace_document" "this" {
  type = "dashboard"
  name = "#name#"
  content = jsonencode(
    {
      "version" : 13,
      "variables" : [],
      "tiles" : {
        "0" : {
          "type" : "markdown",
          "title" : "",
          "content" : "![Image of a Dashboard](https://dt-cdn.net/wp-content/uploads/2022/09/pic1____Dashboard-Preset___PNG.png)"
        },
        "1" : {
          "type" : "data",
          "title" : "",
          "query" : "timeseries avg(dt.host.cpu.user)",
          "queryConfig" : {
            "additionalFilters" : {},
            "version" : "4.3.1",
            "datatype" : "metrics",
            "metricKey" : "dt.host.cpu.user",
            "aggregation" : "avg",
            "by" : []
          },
          "subType" : "dql-builder-metrics",
          "visualization" : "lineChart",
          "visualizationSettings" : {
            "thresholds" : [],
            "chartSettings" : {
              "gapPolicy" : "connect",
              "circleChartSettings" : {
                "groupingThresholdType" : "relative",
                "groupingThresholdValue" : 0,
                "valueType" : "relative"
              },
              "categoryOverrides" : {},
              "fieldMapping" : {
                "timestamp" : "timeframe",
                "leftAxisValues" : [
                  "avg(dt.host.cpu.user)"
                ],
                "leftAxisDimensions" : [],
                "fields" : [],
                "values" : []
              }
            },
            "singleValue" : {
              "showLabel" : true,
              "label" : "",
              "prefixIcon" : "",
              "autoscale" : true,
              "alignment" : "center",
              "colorThresholdTarget" : "value"
            },
            "table" : {
              "rowDensity" : "condensed",
              "enableSparklines" : false,
              "hiddenColumns" : [],
              "lineWrapIds" : [],
              "columnWidths" : {}
            }
          }
        },
        "2" : {
          "type" : "data",
          "title" : "",
          "query" : "timeseries avg(dt.host.memory.used)",
          "queryConfig" : {
            "additionalFilters" : {},
            "version" : "4.3.1",
            "datatype" : "metrics",
            "metricKey" : "dt.host.memory.used",
            "aggregation" : "avg",
            "by" : []
          },
          "subType" : "dql-builder-metrics",
          "visualization" : "lineChart",
          "visualizationSettings" : {
            "thresholds" : [],
            "chartSettings" : {
              "gapPolicy" : "connect",
              "circleChartSettings" : {
                "groupingThresholdType" : "relative",
                "groupingThresholdValue" : 0,
                "valueType" : "relative"
              },
              "categoryOverrides" : {},
              "fieldMapping" : {
                "timestamp" : "timeframe",
                "leftAxisValues" : [
                  "avg(dt.host.memory.used)"
                ],
                "leftAxisDimensions" : [],
                "fields" : [],
                "values" : []
              },
              "categoricalBarChartSettings" : {}
            },
            "singleValue" : {
              "showLabel" : true,
              "label" : "",
              "prefixIcon" : "",
              "autoscale" : true,
              "alignment" : "center",
              "colorThresholdTarget" : "value"
            },
            "table" : {
              "rowDensity" : "condensed",
              "enableSparklines" : false,
              "hiddenColumns" : [],
              "lineWrapIds" : [],
              "columnWidths" : {}
            }
          }
        }
      },
      "layouts" : {
        "0" : {
          "x" : 0,
          "y" : 0,
          "w" : 24,
          "h" : 14
        },
        "1" : {
          "x" : 0,
          "y" : 14,
          "w" : 9,
          "h" : 6
        },
        "2" : {
          "x" : 15,
          "y" : 14,
          "w" : 9,
          "h" : 6
        }
      }
    }
  )
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `document_id` (String) Document ID
- `recipients` (Block List, Min: 1, Max: 1) Recipients of the direct share (see [below for nested schema](#nestedblock--recipients))

### Optional

- `access` (String) Access grants. Possible values are `read` and `read-write`

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--recipients"></a>
### Nested Schema for `recipients`

Optional:

- `recipient` (Block Set, Max: 1000) Recipient of the direct share (see [below for nested schema](#nestedblock--recipients--recipient))

<a id="nestedblock--recipients--recipient"></a>
### Nested Schema for `recipients.recipient`

Required:

- `id` (String) Identifier of the recipient

Optional:

- `type` (String) Type of the recipient. Possible values are `group' and `user'
