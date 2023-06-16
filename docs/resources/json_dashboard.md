---
layout: ""
page_title: dynatrace_json_dashboard Resource - terraform-provider-dynatrace"
description: |-
  The resource `dynatrace_json_dashboard` covers configuration for dashboards in JSON format
---

# dynatrace_json_dashboard (Resource)

-> This resource is excluded by default in the export utility since there could be a large amount of dashboards.

## Dynatrace Documentation

- Dashboards and reports - https://www.dynatrace.com/support/help/how-to-use-dynatrace/dashboards-and-charts

- Dashboards API - https://www.dynatrace.com/support/help/dynatrace-api/configuration-api/dashboards-api

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_json_dashboard` downloads all existing dashboard configuration

The full documentation of the export feature is available [here](https://registry.terraform.io/providers/dynatrace-oss/dynatrace/latest/docs/guides/export-v2).

## Resource Example Usage

```terraform
resource "dynatrace_json_dashboard" "#name#" {
  contents = jsonencode(
    {
      "dashboardMetadata" : {
        "name" : "#name#",
        "shared" : false,
        "preset" : false,
        "owner" : "Dynatrace",
        "tags" : [
          "Kubernetes"
        ],
        "dynamicFilters" : {
          "filters" : [
            "KUBERNETES_CLUSTER"
          ]
        },
        "hasConsistentColors" : false
      },
      "tiles" : [
        {
          "name" : "Markdown",
          "tileType" : "MARKDOWN",
          "configured" : true,
          "bounds" : {
            "top" : 0,
            "left" : 0,
            "width" : 684,
            "height" : 38
          },
          "tileFilter" : {},
          "markdown" : "## Cluster resource overview"
        },
        {
          "name" : "",
          "tileType" : "HOSTS",
          "configured" : true,
          "bounds" : {
            "top" : 38,
            "left" : 342,
            "width" : 342,
            "height" : 304
          },
          "tileFilter" : {
            "timeframe" : "-5m"
          },
          "filterConfig" : {
            "type" : "HOST",
            "customName" : "Full-Stack Kubernetes nodes",
            "defaultName" : "Full-Stack Kubernetes nodes",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "TIMESERIES",
              "series" : [],
              "resultMetadata" : {}
            },
            "filtersPerEntityType" : {
              "HOST" : {
                "HOST_SOFTWARE_TECH" : [
                  "KUBERNETES"
                ]
              }
            }
          },
          "chartVisible" : true
        },
        {
          "name" : "Custom chart",
          "tileType" : "CUSTOM_CHARTING",
          "configured" : true,
          "bounds" : {
            "top" : 418,
            "left" : 190,
            "width" : 190,
            "height" : 152
          },
          "tileFilter" : {
            "timeframe" : "-5m"
          },
          "filterConfig" : {
            "type" : "MIXED",
            "customName" : "CPU available",
            "defaultName" : "Custom chart",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "SINGLE_VALUE",
              "series" : [
                {
                  "metric" : "builtin:cloud.kubernetes.cluster.cpuAvailable",
                  "aggregation" : "SUM_DIMENSIONS",
                  "type" : "LINE",
                  "entityType" : "KUBERNETES_CLUSTER",
                  "dimensions" : [],
                  "sortAscending" : false,
                  "sortColumn" : true,
                  "aggregationRate" : "TOTAL"
                }
              ],
              "resultMetadata" : {}
            },
            "filtersPerEntityType" : {}
          }
        },
        {
          "name" : "",
          "tileType" : "CUSTOM_CHARTING",
          "configured" : true,
          "bounds" : {
            "top" : 38,
            "left" : 684,
            "width" : 304,
            "height" : 304
          },
          "tileFilter" : {
            "timeframe" : "-5m"
          },
          "filterConfig" : {
            "type" : "MIXED",
            "customName" : "Pods",
            "defaultName" : "Custom chart",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "PIE",
              "series" : [
                {
                  "metric" : "builtin:cloud.kubernetes.workload.pods",
                  "aggregation" : "SUM_DIMENSIONS",
                  "type" : "LINE",
                  "entityType" : "CLOUD_APPLICATION",
                  "dimensions" : [
                    {
                      "id" : "1",
                      "name" : "Pod phase",
                      "values" : [],
                      "entityDimension" : false
                    }
                  ],
                  "sortAscending" : false,
                  "sortColumn" : true,
                  "aggregationRate" : "TOTAL"
                }
              ],
              "resultMetadata" : {
                "null¦Pod phase»Pending»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION" : {
                  "lastModified" : 1597234457744,
                  "customColor" : "#f5d30f"
                },
                "null¦Pod phase»Succeeded»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION" : {
                  "lastModified" : 1597237249882,
                  "customColor" : "#008cdb"
                },
                "null¦Pod phase»Running»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION" : {
                  "lastModified" : 1597234642722,
                  "customColor" : "#64bd64"
                },
                "null¦Pod phase»Failed»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION" : {
                  "lastModified" : 1597234118116,
                  "customColor" : "#ff0000"
                }
              }
            },
            "filtersPerEntityType" : {}
          }
        },
        {
          "name" : "Custom chart",
          "tileType" : "CUSTOM_CHARTING",
          "configured" : true,
          "bounds" : {
            "top" : 418,
            "left" : 608,
            "width" : 190,
            "height" : 152
          },
          "tileFilter" : {
            "timeframe" : "-5m"
          },
          "filterConfig" : {
            "type" : "MIXED",
            "customName" : "Memory available",
            "defaultName" : "Custom chart",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "SINGLE_VALUE",
              "series" : [
                {
                  "metric" : "builtin:cloud.kubernetes.cluster.memoryAvailable",
                  "aggregation" : "SUM_DIMENSIONS",
                  "type" : "LINE",
                  "entityType" : "KUBERNETES_CLUSTER",
                  "dimensions" : [],
                  "sortAscending" : false,
                  "sortColumn" : true,
                  "aggregationRate" : "TOTAL"
                }
              ],
              "resultMetadata" : {}
            },
            "filtersPerEntityType" : {}
          }
        },
        {
          "name" : "Markdown",
          "tileType" : "MARKDOWN",
          "configured" : true,
          "bounds" : {
            "top" : 380,
            "left" : 0,
            "width" : 1634,
            "height" : 38
          },
          "tileFilter" : {},
          "markdown" : "## Node resource usage"
        },
        {
          "name" : "Custom chart",
          "tileType" : "CUSTOM_CHARTING",
          "configured" : true,
          "bounds" : {
            "top" : 38,
            "left" : 0,
            "width" : 342,
            "height" : 304
          },
          "tileFilter" : {
            "timeframe" : "-5m"
          },
          "filterConfig" : {
            "type" : "MIXED",
            "customName" : "Cluster nodes",
            "defaultName" : "Custom chart",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "PIE",
              "series" : [
                {
                  "metric" : "builtin:cloud.kubernetes.cluster.nodes",
                  "aggregation" : "AVG",
                  "type" : "LINE",
                  "entityType" : "KUBERNETES_CLUSTER",
                  "dimensions" : [
                    {
                      "id" : "0",
                      "name" : "dt.entity.kubernetes_cluster",
                      "values" : [],
                      "entityDimension" : true
                    }
                  ],
                  "sortAscending" : false,
                  "sortColumn" : true,
                  "aggregationRate" : "TOTAL"
                }
              ],
              "resultMetadata" : {}
            },
            "filtersPerEntityType" : {}
          }
        },
        {
          "name" : "",
          "tileType" : "CUSTOM_CHARTING",
          "configured" : true,
          "bounds" : {
            "top" : 418,
            "left" : 1026,
            "width" : 190,
            "height" : 152
          },
          "tileFilter" : {
            "timeframe" : "-5m"
          },
          "filterConfig" : {
            "type" : "MIXED",
            "customName" : "Disk available",
            "defaultName" : "Custom chart",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "SINGLE_VALUE",
              "series" : [
                {
                  "metric" : "builtin:host.disk.avail",
                  "aggregation" : "SUM_DIMENSIONS",
                  "type" : "LINE",
                  "entityType" : "HOST",
                  "dimensions" : [],
                  "sortAscending" : false,
                  "sortColumn" : true,
                  "aggregationRate" : "TOTAL"
                }
              ],
              "resultMetadata" : {}
            },
            "filtersPerEntityType" : {
              "HOST" : {
                "HOST_SOFTWARE_TECH" : [
                  "KUBERNETES"
                ]
              }
            }
          }
        },
        {
          "name" : "Custom chart",
          "tileType" : "CUSTOM_CHARTING",
          "configured" : true,
          "bounds" : {
            "top" : 570,
            "left" : 0,
            "width" : 418,
            "height" : 304
          },
          "tileFilter" : {},
          "filterConfig" : {
            "type" : "MIXED",
            "customName" : "CPU usage % ",
            "defaultName" : "Custom chart",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "TIMESERIES",
              "series" : [
                {
                  "metric" : "builtin:host.cpu.usage",
                  "aggregation" : "AVG",
                  "type" : "LINE",
                  "entityType" : "HOST",
                  "dimensions" : [
                    {
                      "id" : "0",
                      "name" : "dt.entity.host",
                      "values" : [],
                      "entityDimension" : true
                    }
                  ],
                  "sortAscending" : false,
                  "sortColumn" : true,
                  "aggregationRate" : "TOTAL"
                }
              ],
              "resultMetadata" : {}
            },
            "filtersPerEntityType" : {
              "HOST" : {
                "HOST_SOFTWARE_TECH" : [
                  "KUBERNETES"
                ]
              }
            }
          }
        },
        {
          "name" : "Custom chart",
          "tileType" : "CUSTOM_CHARTING",
          "configured" : true,
          "bounds" : {
            "top" : 570,
            "left" : 418,
            "width" : 418,
            "height" : 304
          },
          "tileFilter" : {},
          "filterConfig" : {
            "type" : "MIXED",
            "customName" : "Memory usage % ",
            "defaultName" : "Custom chart",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "TIMESERIES",
              "series" : [
                {
                  "metric" : "builtin:host.mem.usage",
                  "aggregation" : "AVG",
                  "type" : "LINE",
                  "entityType" : "HOST",
                  "dimensions" : [
                    {
                      "id" : "0",
                      "name" : "dt.entity.host",
                      "values" : [],
                      "entityDimension" : true
                    }
                  ],
                  "sortAscending" : false,
                  "sortColumn" : true,
                  "aggregationRate" : "TOTAL"
                }
              ],
              "resultMetadata" : {}
            },
            "filtersPerEntityType" : {
              "HOST" : {
                "HOST_SOFTWARE_TECH" : [
                  "KUBERNETES"
                ]
              }
            }
          }
        },
        {
          "name" : "Custom chart",
          "tileType" : "CUSTOM_CHARTING",
          "configured" : true,
          "bounds" : {
            "top" : 570,
            "left" : 836,
            "width" : 418,
            "height" : 304
          },
          "tileFilter" : {},
          "filterConfig" : {
            "type" : "MIXED",
            "customName" : "Disk usage % ",
            "defaultName" : "Custom chart",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "TIMESERIES",
              "series" : [
                {
                  "metric" : "builtin:host.disk.usedPct",
                  "aggregation" : "AVG",
                  "type" : "LINE",
                  "entityType" : "HOST",
                  "dimensions" : [
                    {
                      "id" : "0",
                      "name" : "dt.entity.host",
                      "values" : [],
                      "entityDimension" : true
                    },
                    {
                      "id" : "1",
                      "name" : "dt.entity.disk",
                      "values" : [],
                      "entityDimension" : true
                    }
                  ],
                  "sortAscending" : false,
                  "sortColumn" : true,
                  "aggregationRate" : "TOTAL"
                }
              ],
              "resultMetadata" : {}
            },
            "filtersPerEntityType" : {
              "HOST" : {
                "HOST_SOFTWARE_TECH" : [
                  "KUBERNETES"
                ]
              }
            }
          }
        },
        {
          "name" : "Custom chart",
          "tileType" : "CUSTOM_CHARTING",
          "configured" : true,
          "bounds" : {
            "top" : 418,
            "left" : 0,
            "width" : 190,
            "height" : 152
          },
          "tileFilter" : {
            "timeframe" : "-5m"
          },
          "filterConfig" : {
            "type" : "MIXED",
            "customName" : "Total CPU requests",
            "defaultName" : "Custom chart",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "SINGLE_VALUE",
              "series" : [
                {
                  "metric" : "builtin:cloud.kubernetes.cluster.cpuRequested",
                  "aggregation" : "SUM_DIMENSIONS",
                  "type" : "LINE",
                  "entityType" : "KUBERNETES_CLUSTER",
                  "dimensions" : [],
                  "sortAscending" : false,
                  "sortColumn" : true,
                  "aggregationRate" : "TOTAL"
                }
              ],
              "resultMetadata" : {}
            },
            "filtersPerEntityType" : {}
          }
        },
        {
          "name" : "Custom chart",
          "tileType" : "CUSTOM_CHARTING",
          "configured" : true,
          "bounds" : {
            "top" : 418,
            "left" : 418,
            "width" : 190,
            "height" : 152
          },
          "tileFilter" : {
            "timeframe" : "-5m"
          },
          "filterConfig" : {
            "type" : "MIXED",
            "customName" : "Total memory requests",
            "defaultName" : "Custom chart",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "SINGLE_VALUE",
              "series" : [
                {
                  "metric" : "builtin:cloud.kubernetes.cluster.memoryRequested",
                  "aggregation" : "SUM_DIMENSIONS",
                  "type" : "LINE",
                  "entityType" : "KUBERNETES_CLUSTER",
                  "dimensions" : [],
                  "sortAscending" : false,
                  "sortColumn" : true,
                  "aggregationRate" : "TOTAL"
                }
              ],
              "resultMetadata" : {}
            },
            "filtersPerEntityType" : {}
          }
        },
        {
          "name" : "",
          "tileType" : "CUSTOM_CHARTING",
          "configured" : true,
          "bounds" : {
            "top" : 418,
            "left" : 836,
            "width" : 190,
            "height" : 152
          },
          "tileFilter" : {
            "timeframe" : "-5m"
          },
          "filterConfig" : {
            "type" : "MIXED",
            "customName" : "Total disk used",
            "defaultName" : "Custom chart",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "SINGLE_VALUE",
              "series" : [
                {
                  "metric" : "builtin:host.disk.used",
                  "aggregation" : "SUM_DIMENSIONS",
                  "type" : "LINE",
                  "entityType" : "HOST",
                  "dimensions" : [],
                  "sortAscending" : false,
                  "sortColumn" : true,
                  "aggregationRate" : "TOTAL"
                }
              ],
              "resultMetadata" : {}
            },
            "filtersPerEntityType" : {
              "HOST" : {
                "HOST_SOFTWARE_TECH" : [
                  "KUBERNETES"
                ]
              }
            }
          }
        },
        {
          "name" : "Custom chart",
          "tileType" : "CUSTOM_CHARTING",
          "configured" : true,
          "bounds" : {
            "top" : 570,
            "left" : 1254,
            "width" : 380,
            "height" : 304
          },
          "tileFilter" : {},
          "filterConfig" : {
            "type" : "MIXED",
            "customName" : "Traffic in/out",
            "defaultName" : "Custom chart",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "TIMESERIES",
              "series" : [
                {
                  "metric" : "builtin:host.net.nic.trafficIn",
                  "aggregation" : "SUM_DIMENSIONS",
                  "type" : "LINE",
                  "entityType" : "HOST",
                  "dimensions" : [
                    {
                      "id" : "0",
                      "name" : "dt.entity.host",
                      "values" : [],
                      "entityDimension" : true
                    }
                  ],
                  "sortAscending" : false,
                  "sortColumn" : false,
                  "aggregationRate" : "TOTAL"
                },
                {
                  "metric" : "builtin:host.net.nic.trafficOut",
                  "aggregation" : "SUM_DIMENSIONS",
                  "type" : "LINE",
                  "entityType" : "HOST",
                  "dimensions" : [
                    {
                      "id" : "0",
                      "name" : "dt.entity.host",
                      "values" : [],
                      "entityDimension" : true
                    }
                  ],
                  "sortAscending" : false,
                  "sortColumn" : true,
                  "aggregationRate" : "TOTAL"
                }
              ],
              "resultMetadata" : {}
            },
            "filtersPerEntityType" : {
              "HOST" : {
                "HOST_SOFTWARE_TECH" : [
                  "KUBERNETES"
                ]
              }
            }
          }
        },
        {
          "name" : "",
          "tileType" : "CUSTOM_CHARTING",
          "configured" : true,
          "bounds" : {
            "top" : 418,
            "left" : 1444,
            "width" : 190,
            "height" : 152
          },
          "tileFilter" : {
            "timeframe" : "-5m"
          },
          "filterConfig" : {
            "type" : "MIXED",
            "customName" : "Traffic out",
            "defaultName" : "Custom chart",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "SINGLE_VALUE",
              "series" : [
                {
                  "metric" : "builtin:host.net.nic.trafficOut",
                  "aggregation" : "SUM_DIMENSIONS",
                  "type" : "LINE",
                  "entityType" : "HOST",
                  "dimensions" : [],
                  "sortAscending" : false,
                  "sortColumn" : true,
                  "aggregationRate" : "TOTAL"
                }
              ],
              "resultMetadata" : {}
            },
            "filtersPerEntityType" : {
              "HOST" : {
                "HOST_SOFTWARE_TECH" : [
                  "KUBERNETES"
                ]
              }
            }
          }
        },
        {
          "name" : "",
          "tileType" : "CUSTOM_CHARTING",
          "configured" : true,
          "bounds" : {
            "top" : 418,
            "left" : 1254,
            "width" : 190,
            "height" : 152
          },
          "tileFilter" : {
            "timeframe" : "-5m"
          },
          "filterConfig" : {
            "type" : "MIXED",
            "customName" : "Traffic in",
            "defaultName" : "Custom chart",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "SINGLE_VALUE",
              "series" : [
                {
                  "metric" : "builtin:host.net.nic.trafficIn",
                  "aggregation" : "SUM_DIMENSIONS",
                  "type" : "LINE",
                  "entityType" : "HOST",
                  "dimensions" : [],
                  "sortAscending" : false,
                  "sortColumn" : true,
                  "aggregationRate" : "TOTAL"
                }
              ],
              "resultMetadata" : {}
            },
            "filtersPerEntityType" : {
              "HOST" : {
                "HOST_SOFTWARE_TECH" : [
                  "KUBERNETES"
                ]
              }
            }
          }
        },
        {
          "name" : "Markdown",
          "tileType" : "MARKDOWN",
          "configured" : true,
          "bounds" : {
            "top" : 0,
            "left" : 684,
            "width" : 950,
            "height" : 38
          },
          "tileFilter" : {},
          "markdown" : "## [Workloads overview](#dashboard;id=6b38732e-d26b-45c7-b107-ed85e87ff288)"
        },
        {
          "name" : "",
          "tileType" : "CUSTOM_CHARTING",
          "configured" : true,
          "bounds" : {
            "top" : 38,
            "left" : 1330,
            "width" : 304,
            "height" : 304
          },
          "tileFilter" : {},
          "filterConfig" : {
            "type" : "MIXED",
            "customName" : "Workloads",
            "defaultName" : "Custom chart",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "PIE",
              "series" : [
                {
                  "metric" : "builtin:cloud.kubernetes.namespace.workloads",
                  "aggregation" : "SUM_DIMENSIONS",
                  "type" : "LINE",
                  "entityType" : "CLOUD_APPLICATION_NAMESPACE",
                  "dimensions" : [
                    {
                      "id" : "1",
                      "name" : "Deployment type",
                      "values" : [],
                      "entityDimension" : false
                    }
                  ],
                  "sortAscending" : false,
                  "sortColumn" : true,
                  "aggregationRate" : "TOTAL"
                }
              ],
              "resultMetadata" : {
                "null¦Pod phase»Pending»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION" : {
                  "lastModified" : 1597234457744,
                  "customColor" : "#f5d30f"
                },
                "null¦Pod phase»Succeeded»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION" : {
                  "lastModified" : 1597237249882,
                  "customColor" : "#008cdb"
                },
                "null¦Deployment type»DaemonSet»falsebuiltin:cloud.kubernetes.namespace.workloads|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION_NAMESPACE" : {
                  "lastModified" : 1597858600132,
                  "customColor" : "#ffa86c"
                },
                "null¦Pod phase»Running»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION" : {
                  "lastModified" : 1597234642722,
                  "customColor" : "#64bd64"
                },
                "null¦Pod phase»Failed»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION" : {
                  "lastModified" : 1597234118116,
                  "customColor" : "#ff0000"
                }
              }
            },
            "filtersPerEntityType" : {}
          }
        },
        {
          "name" : "",
          "tileType" : "CUSTOM_CHARTING",
          "configured" : true,
          "bounds" : {
            "top" : 38,
            "left" : 988,
            "width" : 342,
            "height" : 304
          },
          "tileFilter" : {
            "timeframe" : "-5m"
          },
          "filterConfig" : {
            "type" : "MIXED",
            "customName" : "Running pods",
            "defaultName" : "Custom chart",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "TOP_LIST",
              "series" : [
                {
                  "metric" : "builtin:cloud.kubernetes.namespace.runningPods",
                  "aggregation" : "SUM_DIMENSIONS",
                  "type" : "LINE",
                  "entityType" : "CLOUD_APPLICATION_NAMESPACE",
                  "dimensions" : [
                    {
                      "id" : "0",
                      "name" : "dt.entity.cloud_application_namespace",
                      "values" : [],
                      "entityDimension" : true
                    }
                  ],
                  "sortAscending" : false,
                  "sortColumn" : true,
                  "aggregationRate" : "TOTAL"
                }
              ],
              "resultMetadata" : {
                "null¦Pod phase»Pending»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION" : {
                  "lastModified" : 1597234457744,
                  "customColor" : "#f5d30f"
                },
                "null¦Pod phase»Succeeded»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION" : {
                  "lastModified" : 1597237249882,
                  "customColor" : "#008cdb"
                },
                "null¦Deployment type»DaemonSet»falsebuiltin:cloud.kubernetes.namespace.workloads|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION_NAMESPACE" : {
                  "lastModified" : 1597858600132,
                  "customColor" : "#ffa86c"
                },
                "null¦Pod phase»Running»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION" : {
                  "lastModified" : 1597234642722,
                  "customColor" : "#64bd64"
                },
                "null¦Pod phase»Failed»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION" : {
                  "lastModified" : 1597234118116,
                  "customColor" : "#ff0000"
                }
              }
            },
            "filtersPerEntityType" : {}
          }
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

### Read-Only

- `id` (String) The ID of this resource.
 