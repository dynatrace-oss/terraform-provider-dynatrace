---
layout: ""
page_title: "Resource - terraform-provider-dynatrace"
subcategory: "Documents"
description: |-
  The resource `dynatrace_document` covers configuration of Documents (dashboards and notebooks) in Dynatrace.
---

# dynatrace_document (Resource)

-> **Dynatrace SaaS only**

-> To utilize this resource, please define the environment variables `DT_CLIENT_ID`, `DT_CLIENT_SECRET`, `DT_ACCOUNT_ID` with an OAuth client including the following permissions: **Create and edit documents** (`document:documents:write`), **View documents** (`document:documents:read`), **Delete documents** (`document:documents:delete`), and  **Delete documents from trash** (`document:trash.documents:delete`).

-> This resource is excluded by default in the export utility, please explicitly specify the resource to retrieve existing configuration.

## Dynatrace Documentation

- Dynatrace Documents - https://########.apps.dynatrace.com/platform/swagger-ui/index.html?urls.primaryName=Document%20Service

## Resource Example Usage

```terraform
resource "dynatrace_document" "this" {
  type = "dashboard"
  name = "Example Dashboard"
  custom_id = "#name#"
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

# data "dynatrace_documents" "all-dashboard-and-notebooks" {}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `content` (String) Document content as JSON
- `name` (String) The name/name of the document
- `type` (String) Type of the document. Possible Values are `dashboard`, `launchpad` and `notebook`

### Optional

- `custom_id` (String) If provided, this will be the id of the document. If not provided, a system-generated id is used.
- `private` (Boolean) Specifies whether the document is private or readable by everybody

### Read-Only

- `id` (String) The ID of this resource.
- `owner` (String) The ID of the owner of this document
- `version` (Number) The version of the document
