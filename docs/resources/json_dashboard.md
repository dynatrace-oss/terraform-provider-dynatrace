---
layout: ""
page_title: "dynatrace_json_dashboard Resource - terraform-provider-dynatrace"
subcategory: "Dashboards"
description: |-
  The resource `dynatrace_json_dashboard` covers configuration for dashboards in JSON format
---

# dynatrace_json_dashboard (Resource)

-> This resource is excluded by default in the export utility since there could be a large amount of dashboards.

-> This resource requires the API token scopes **Read configuration** (`ReadConfig`) and **Write configuration** (`WriteConfig`)

## Dynatrace Documentation

- Dashboards and reports - https://www.dynatrace.com/support/help/how-to-use-dynatrace/dashboards-and-charts

- Dashboards API - https://www.dynatrace.com/support/help/dynatrace-api/configuration-api/dashboards-api

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_json_dashboard` downloads all existing dashboard configuration

The full documentation of the export feature is available [here](https://dt-url.net/h203qmc).

## Resource Example Usage

```terraform
resource "dynatrace_json_dashboard" "name" {
  contents = jsonencode(
    {
      "dashboardMetadata": {
        "name": "Performance overview",
        "shared": true,
        "owner": "Dynatrace",
        "tags": [
          "performance"
        ],
        "preset": true,
        "hasConsistentColors": false
      },
      "tiles": [
        {
          "name": "Performance",
          "tileType": "HEADER",
          "configured": true,
          "bounds": {
            "top": 0,
            "left": 38,
            "width": 1026,
            "height": 38
          },
          "tileFilter": {},
          "isAutoRefreshDisabled": true
        },
        {
          "name": "Failure rate by service",
          "tileType": "DATA_EXPLORER",
          "configured": true,
          "bounds": {
            "top": 342,
            "left": 38,
            "width": 342,
            "height": 304
          },
          "tileFilter": {},
          "isAutoRefreshDisabled": true,
          "customName": "Successful calls by service instance",
          "queries": [
            {
              "id": "A",
              "metric": "builtin:service.errors.total.rate",
              "spaceAggregation": "AUTO",
              "timeAggregation": "DEFAULT",
              "splitBy": [
                "dt.entity.service"
              ],
              "sortBy": "DESC",
              "sortByDimension": "",
              "filterBy": {
                "nestedFilters": [],
                "criteria": []
              },
              "limit": 20,
              "rate": "NONE",
              "enabled": true
            }
          ],
          "visualConfig": {
            "type": "TOP_LIST",
            "global": {
              "hideLegend": false
            },
            "rules": [
              {
                "matcher": "A:",
                "properties": {
                  "color": "DEFAULT"
                },
                "seriesOverrides": []
              }
            ],
            "axes": {
              "xAxis": {
                "visible": true
              },
              "yAxes": []
            },
            "heatmapSettings": {
              "yAxis": "VALUE"
            },
            "singleValueSettings": {
              "showSparkLine": true
            },
            "thresholds": [
              {
                "axisTarget": "LEFT",
                "rules": [
                  {
                    "color": "#7dc540"
                  },
                  {
                    "color": "#f5d30f"
                  },
                  {
                    "color": "#dc172a"
                  }
                ],
                "visible": true
              }
            ],
            "tableSettings": {
              "hiddenColumns": []
            },
            "graphChartSettings": {
              "connectNulls": false
            },
            "honeycombSettings": {
              "showHive": true,
              "showLegend": true,
              "showLabels": false
            }
          },
          "queriesSettings": {
            "resolution": ""
          },
          "metricExpressions": [
            "resolution=Inf&(builtin:service.errors.total.rate:splitBy(\"dt.entity.service\"):sort(value(auto,descending)):limit(20)):limit(100):names"
          ]
        },
        {
          "name": "Total calls",
          "tileType": "DATA_EXPLORER",
          "configured": true,
          "bounds": {
            "top": 38,
            "left": 38,
            "width": 342,
            "height": 304
          },
          "tileFilter": {},
          "isAutoRefreshDisabled": true,
          "customName": "Total calls",
          "queries": [
            {
              "id": "A",
              "metric": "builtin:service.errors.total.successCount",
              "spaceAggregation": "SUM",
              "timeAggregation": "DEFAULT",
              "splitBy": [
                "dt.entity.service"
              ],
              "sortBy": "DESC",
              "sortByDimension": "",
              "filterBy": {
                "nestedFilters": [],
                "criteria": []
              },
              "limit": 20,
              "rate": "NONE",
              "enabled": true
            },
            {
              "id": "B",
              "metric": "builtin:service.errors.fourxx.successCount",
              "spaceAggregation": "SUM",
              "timeAggregation": "DEFAULT",
              "splitBy": [
                "dt.entity.service"
              ],
              "sortBy": "DESC",
              "sortByDimension": "",
              "filterBy": {
                "nestedFilters": [],
                "criteria": []
              },
              "limit": 20,
              "rate": "NONE",
              "enabled": true
            },
            {
              "id": "C",
              "metric": "builtin:service.errors.fivexx.successCount",
              "spaceAggregation": "SUM",
              "timeAggregation": "DEFAULT",
              "splitBy": [
                "dt.entity.service"
              ],
              "sortBy": "DESC",
              "sortByDimension": "",
              "filterBy": {
                "nestedFilters": [],
                "criteria": []
              },
              "limit": 20,
              "rate": "NONE",
              "enabled": true
            }
          ],
          "visualConfig": {
            "type": "STACKED_AREA",
            "global": {
              "hideLegend": false
            },
            "rules": [
              {
                "matcher": "A:",
                "properties": {
                  "color": "DEFAULT"
                },
                "seriesOverrides": []
              },
              {
                "matcher": "B:",
                "properties": {
                  "color": "DEFAULT"
                },
                "seriesOverrides": []
              },
              {
                "matcher": "C:",
                "properties": {
                  "color": "DEFAULT"
                },
                "seriesOverrides": []
              }
            ],
            "axes": {
              "xAxis": {
                "displayName": "",
                "visible": true
              },
              "yAxes": [
                {
                  "displayName": "",
                  "visible": true,
                  "min": "AUTO",
                  "max": "AUTO",
                  "position": "LEFT",
                  "queryIds": [
                    "A",
                    "B",
                    "C"
                  ],
                  "defaultAxis": true
                }
              ]
            },
            "heatmapSettings": {
              "yAxis": "VALUE"
            },
            "singleValueSettings": {
              "showSparkLine": true
            },
            "thresholds": [
              {
                "axisTarget": "LEFT",
                "rules": [
                  {
                    "color": "#7dc540"
                  },
                  {
                    "color": "#f5d30f"
                  },
                  {
                    "color": "#dc172a"
                  }
                ],
                "visible": true
              }
            ],
            "tableSettings": {
              "hiddenColumns": []
            },
            "graphChartSettings": {
              "connectNulls": false
            },
            "honeycombSettings": {
              "showHive": true,
              "showLegend": true,
              "showLabels": false
            }
          },
          "queriesSettings": {
            "resolution": ""
          },
          "metricExpressions": [
            "resolution=null&(builtin:service.errors.total.successCount:splitBy(\"dt.entity.service\"):sum:sort(value(sum,descending)):limit(20)):limit(100):names,(builtin:service.errors.fourxx.successCount:splitBy(\"dt.entity.service\"):sum:sort(value(sum,descending)):limit(20)):limit(100):names,(builtin:service.errors.fivexx.successCount:splitBy(\"dt.entity.service\"):sum:sort(value(sum,descending)):limit(20)):limit(100):names"
          ]
        },
        {
          "name": "Total errors",
          "tileType": "DATA_EXPLORER",
          "configured": true,
          "bounds": {
            "top": 38,
            "left": 380,
            "width": 342,
            "height": 304
          },
          "tileFilter": {},
          "isAutoRefreshDisabled": true,
          "customName": "Total errors",
          "queries": [
            {
              "id": "A",
              "metric": "builtin:service.errors.total.count",
              "spaceAggregation": "SUM",
              "timeAggregation": "DEFAULT",
              "splitBy": [],
              "sortBy": "DESC",
              "sortByDimension": "",
              "filterBy": {
                "nestedFilters": [],
                "criteria": []
              },
              "limit": 20,
              "rate": "NONE",
              "enabled": true
            }
          ],
          "visualConfig": {
            "type": "SINGLE_VALUE",
            "global": {
              "hideLegend": false
            },
            "rules": [
              {
                "matcher": "A:",
                "properties": {
                  "color": "DEFAULT"
                },
                "seriesOverrides": []
              }
            ],
            "axes": {
              "xAxis": {
                "visible": true
              },
              "yAxes": []
            },
            "heatmapSettings": {
              "yAxis": "VALUE"
            },
            "singleValueSettings": {
              "showTrend": false,
              "showSparkLine": true,
              "linkTileColorToThreshold": false
            },
            "thresholds": [
              {
                "axisTarget": "LEFT",
                "rules": [
                  {
                    "color": "#7dc540"
                  },
                  {
                    "color": "#f5d30f"
                  },
                  {
                    "color": "#dc172a"
                  }
                ],
                "visible": true
              }
            ],
            "tableSettings": {
              "hiddenColumns": []
            },
            "graphChartSettings": {
              "connectNulls": false
            },
            "honeycombSettings": {
              "showHive": true,
              "showLegend": true,
              "showLabels": false
            }
          },
          "queriesSettings": {
            "resolution": ""
          },
          "metricExpressions": [
            "resolution=Inf&(builtin:service.errors.total.count:splitBy():sum:sort(value(sum,descending)):limit(20)):limit(100):names",
            "resolution=null&(builtin:service.errors.total.count:splitBy():sum:sort(value(sum,descending)):limit(20))"
          ]
        },
        {
          "name": "Client side errors",
          "tileType": "DATA_EXPLORER",
          "configured": true,
          "bounds": {
            "top": 38,
            "left": 722,
            "width": 342,
            "height": 304
          },
          "tileFilter": {},
          "isAutoRefreshDisabled": true,
          "customName": "Client & server errors by operation",
          "queries": [
            {
              "id": "A",
              "metric": "builtin:service.errors.client.count",
              "spaceAggregation": "SUM",
              "timeAggregation": "DEFAULT",
              "splitBy": [],
              "sortBy": "DESC",
              "sortByDimension": "",
              "filterBy": {
                "nestedFilters": [],
                "criteria": []
              },
              "limit": 20,
              "rate": "NONE",
              "enabled": true
            }
          ],
          "visualConfig": {
            "type": "SINGLE_VALUE",
            "global": {
              "hideLegend": false
            },
            "rules": [
              {
                "matcher": "A:",
                "properties": {
                  "color": "DEFAULT"
                },
                "seriesOverrides": []
              }
            ],
            "axes": {
              "xAxis": {
                "visible": true
              },
              "yAxes": []
            },
            "heatmapSettings": {
              "yAxis": "VALUE"
            },
            "singleValueSettings": {
              "showSparkLine": true
            },
            "thresholds": [
              {
                "axisTarget": "LEFT",
                "rules": [
                  {
                    "color": "#7dc540"
                  },
                  {
                    "color": "#f5d30f"
                  },
                  {
                    "color": "#dc172a"
                  }
                ],
                "visible": true
              }
            ],
            "tableSettings": {
              "hiddenColumns": []
            },
            "graphChartSettings": {
              "connectNulls": false
            },
            "honeycombSettings": {
              "showHive": true,
              "showLegend": true,
              "showLabels": false
            }
          },
          "queriesSettings": {
            "resolution": ""
          },
          "metricExpressions": [
            "resolution=Inf&(builtin:service.errors.client.count:splitBy():sum:sort(value(sum,descending)):limit(20)):limit(100):names",
            "resolution=null&(builtin:service.errors.client.count:splitBy():sum:sort(value(sum,descending)):limit(20))"
          ]
        }
      ]
    }
  )
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `contents` (String) Contains the JSON Code of the Dashboard

### Optional

- `link_id` (String) ID of the dashboard, used with the json_dashboard_base resource and variables to create circular dependencies between dashboards for hyperlinks. See the documentation for `dynatrace_json_dashboard_base` for a concrete example.

### Read-Only

- `id` (String) The ID of this resource.
 