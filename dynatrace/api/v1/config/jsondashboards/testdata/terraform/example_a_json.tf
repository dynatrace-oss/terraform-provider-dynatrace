resource "dynatrace_json_dashboard" "#name#" {
  contents = jsonencode(
    {
      "dashboardMetadata" : {
        "name" : "Azure Cognitive Services",
        "shared" : true,
        "owner" : "Dynatrace",
        "tags" : [
          "Azure"
        ],
        "preset" : true,
        "hasConsistentColors" : false
      },
      "tiles" : [
        {
          "name" : "Custom chart",
          "tileType" : "CUSTOM_CHARTING",
          "configured" : true,
          "bounds" : {
            "top" : 342,
            "left" : 38,
            "width" : 342,
            "height" : 304
          },
          "tileFilter" : {},
          "filterConfig" : {
            "type" : "MIXED",
            "customName" : "Successful calls by service instance",
            "defaultName" : "Custom chart",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "TOP_LIST",
              "series" : [
                {
                  "metric" : "ext:cloud.azure.microsoft_cognitiveservices.accounts.successfulcalls",
                  "aggregation" : "SUM",
                  "type" : "LINE",
                  "entityType" : "IOT",
                  "dimensions" : [
                    {
                      "id" : "0",
                      "name" : "dt.entity.custom_device",
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
          "name" : "Custom chart",
          "tileType" : "CUSTOM_CHARTING",
          "configured" : true,
          "bounds" : {
            "top" : 342,
            "left" : 722,
            "width" : 342,
            "height" : 304
          },
          "tileFilter" : {},
          "filterConfig" : {
            "type" : "MIXED",
            "customName" : "Latency by service instance",
            "defaultName" : "Custom chart",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "TOP_LIST",
              "series" : [
                {
                  "metric" : "ext:cloud.azure.microsoft_cognitiveservices.accounts.latency",
                  "aggregation" : "AVG",
                  "type" : "LINE",
                  "entityType" : "IOT",
                  "dimensions" : [
                    {
                      "id" : "0",
                      "name" : "dt.entity.custom_device",
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
          "name" : "Custom chart",
          "tileType" : "CUSTOM_CHARTING",
          "configured" : true,
          "bounds" : {
            "top" : 38,
            "left" : 380,
            "width" : 342,
            "height" : 304
          },
          "tileFilter" : {},
          "filterConfig" : {
            "type" : "MIXED",
            "customName" : "Total errors",
            "defaultName" : "Custom chart",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "TIMESERIES",
              "series" : [
                {
                  "metric" : "ext:cloud.azure.microsoft_cognitiveservices.accounts.totalerrors",
                  "aggregation" : "SUM",
                  "type" : "BAR",
                  "entityType" : "IOT",
                  "dimensions" : [
                    {
                      "id" : "0",
                      "name" : "dt.entity.custom_device",
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
          "name" : "Custom chart",
          "tileType" : "CUSTOM_CHARTING",
          "configured" : true,
          "bounds" : {
            "top" : 38,
            "left" : 38,
            "width" : 342,
            "height" : 304
          },
          "tileFilter" : {},
          "filterConfig" : {
            "type" : "MIXED",
            "customName" : "Total calls",
            "defaultName" : "Custom chart",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "TIMESERIES",
              "series" : [
                {
                  "metric" : "ext:cloud.azure.microsoft_cognitiveservices.accounts.blockedcalls",
                  "aggregation" : "SUM",
                  "type" : "BAR",
                  "entityType" : "IOT",
                  "dimensions" : [],
                  "sortAscending" : false,
                  "sortColumn" : false,
                  "aggregationRate" : "TOTAL"
                },
                {
                  "metric" : "ext:cloud.azure.microsoft_cognitiveservices.accounts.successfulcalls",
                  "aggregation" : "SUM",
                  "type" : "BAR",
                  "entityType" : "IOT",
                  "dimensions" : [],
                  "sortAscending" : false,
                  "sortColumn" : true,
                  "aggregationRate" : "TOTAL"
                }
              ],
              "resultMetadata" : {
                "nullext:cloud.azure.microsoft_cognitiveservices.accounts.blockedcalls|SUM|TOTAL|BAR|IOT" : {
                  "lastModified" : 1595493591071,
                  "customColor" : "#FF0000"
                },
                "nullext:cloud.azure.microsoft_cognitiveservices.accounts.totalerrors|SUM|TOTAL|LINE|IOT" : {
                  "lastModified" : 1595343403084,
                  "customColor" : "#4fd5e0"
                },
                "nullext:cloud.azure.microsoft_cognitiveservices.accounts.successfulcalls|SUM|TOTAL|BAR|IOT" : {
                  "lastModified" : 1595493603433,
                  "customColor" : "#4fd5e0"
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
            "top" : 38,
            "left" : 722,
            "width" : 342,
            "height" : 304
          },
          "tileFilter" : {},
          "filterConfig" : {
            "type" : "MIXED",
            "customName" : "Client & server errors by operation",
            "defaultName" : "Custom chart",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "TOP_LIST",
              "series" : [
                {
                  "metric" : "ext:cloud.azure.microsoft_cognitiveservices.accounts.clienterrors",
                  "aggregation" : "SUM",
                  "type" : "BAR",
                  "entityType" : "IOT",
                  "dimensions" : [
                    {
                      "id" : "2",
                      "name" : "Operation name",
                      "values" : [],
                      "entityDimension" : false
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
          "name" : "Custom chart",
          "tileType" : "CUSTOM_CHARTING",
          "configured" : true,
          "bounds" : {
            "top" : 646,
            "left" : 38,
            "width" : 342,
            "height" : 304
          },
          "tileFilter" : {},
          "filterConfig" : {
            "type" : "MIXED",
            "customName" : "Successful calls by operation",
            "defaultName" : "Custom chart",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "TOP_LIST",
              "series" : [
                {
                  "metric" : "ext:cloud.azure.microsoft_cognitiveservices.accounts.successfulcalls",
                  "aggregation" : "SUM",
                  "type" : "LINE",
                  "entityType" : "IOT",
                  "dimensions" : [
                    {
                      "id" : "2",
                      "name" : "Operation name",
                      "values" : [],
                      "entityDimension" : false
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
          "name" : "Custom chart",
          "tileType" : "CUSTOM_CHARTING",
          "configured" : true,
          "bounds" : {
            "top" : 646,
            "left" : 380,
            "width" : 342,
            "height" : 304
          },
          "tileFilter" : {},
          "filterConfig" : {
            "type" : "MIXED",
            "customName" : "Blocked calls by operation",
            "defaultName" : "Custom chart",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "TOP_LIST",
              "series" : [
                {
                  "metric" : "ext:cloud.azure.microsoft_cognitiveservices.accounts.blockedcalls",
                  "aggregation" : "SUM",
                  "type" : "LINE",
                  "entityType" : "IOT",
                  "dimensions" : [
                    {
                      "id" : "2",
                      "name" : "Operation name",
                      "values" : [],
                      "entityDimension" : false
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
          "name" : "Custom chart",
          "tileType" : "CUSTOM_CHARTING",
          "configured" : true,
          "bounds" : {
            "top" : 646,
            "left" : 722,
            "width" : 342,
            "height" : 304
          },
          "tileFilter" : {},
          "filterConfig" : {
            "type" : "MIXED",
            "customName" : "Latency by operation",
            "defaultName" : "Custom chart",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "TOP_LIST",
              "series" : [
                {
                  "metric" : "ext:cloud.azure.microsoft_cognitiveservices.accounts.latency",
                  "aggregation" : "AVG",
                  "type" : "LINE",
                  "entityType" : "IOT",
                  "dimensions" : [
                    {
                      "id" : "2",
                      "name" : "Operation name",
                      "values" : [],
                      "entityDimension" : false
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
          "name" : "Custom chart",
          "tileType" : "CUSTOM_CHARTING",
          "configured" : true,
          "bounds" : {
            "top" : 38,
            "left" : 1064,
            "width" : 342,
            "height" : 304
          },
          "tileFilter" : {},
          "filterConfig" : {
            "type" : "MIXED",
            "customName" : "Data IN/OUT",
            "defaultName" : "Custom chart",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "TIMESERIES",
              "series" : [
                {
                  "metric" : "ext:cloud.azure.microsoft_cognitiveservices.accounts.datain",
                  "aggregation" : "SUM",
                  "type" : "BAR",
                  "entityType" : "IOT",
                  "dimensions" : [],
                  "sortAscending" : false,
                  "sortColumn" : false,
                  "aggregationRate" : "TOTAL"
                },
                {
                  "metric" : "ext:cloud.azure.microsoft_cognitiveservices.accounts.dataout",
                  "aggregation" : "SUM",
                  "type" : "BAR",
                  "entityType" : "IOT",
                  "dimensions" : [],
                  "sortAscending" : false,
                  "sortColumn" : true,
                  "aggregationRate" : "TOTAL"
                }
              ],
              "resultMetadata" : {
                "nullext:cloud.azure.microsoft_cognitiveservices.accounts.dataout|SUM|TOTAL|BAR|IOT" : {
                  "lastModified" : 1595495569439,
                  "customColor" : "#008cdb"
                },
                "nullext:cloud.azure.microsoft_cognitiveservices.accounts.datain|SUM|TOTAL|BAR|IOT" : {
                  "lastModified" : 1595495567476,
                  "customColor" : "#aeebf0"
                },
                "nullext:cloud.azure.microsoft_cognitiveservices.accounts.datain|SUM|TOTAL|BAR|IOT|IOT" : {
                  "lastModified" : 1595494661543,
                  "customColor" : "#008cdb"
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
            "top" : 342,
            "left" : 1064,
            "width" : 342,
            "height" : 304
          },
          "tileFilter" : {},
          "filterConfig" : {
            "type" : "MIXED",
            "customName" : "Data IN by operation",
            "defaultName" : "Custom chart",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "TOP_LIST",
              "series" : [
                {
                  "metric" : "ext:cloud.azure.microsoft_cognitiveservices.accounts.datain",
                  "aggregation" : "SUM",
                  "type" : "LINE",
                  "entityType" : "IOT",
                  "dimensions" : [
                    {
                      "id" : "2",
                      "name" : "Operation name",
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
                "nullext:cloud.azure.microsoft_cognitiveservices.accounts.dataout|SUM|TOTAL|LINE|IOT" : {
                  "lastModified" : 1595344576968,
                  "customColor" : "#ef651f"
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
            "top" : 646,
            "left" : 1064,
            "width" : 342,
            "height" : 304
          },
          "tileFilter" : {},
          "filterConfig" : {
            "type" : "MIXED",
            "customName" : "Data OUT by operation",
            "defaultName" : "Custom chart",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "TOP_LIST",
              "series" : [
                {
                  "metric" : "ext:cloud.azure.microsoft_cognitiveservices.accounts.dataout",
                  "aggregation" : "SUM",
                  "type" : "LINE",
                  "entityType" : "IOT",
                  "dimensions" : [
                    {
                      "id" : "2",
                      "name" : "Operation name",
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
                "nullÂ¦Operation nameÂ»ReadÂ»falseext:cloud.azure.microsoft_cognitiveservices.accounts.dataout|SUM|TOTAL|LINE|IOT" : {
                  "lastModified" : 1595344576968,
                  "customColor" : "#ef651f"
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
            "top" : 342,
            "left" : 380,
            "width" : 342,
            "height" : 304
          },
          "tileFilter" : {},
          "filterConfig" : {
            "type" : "MIXED",
            "customName" : "Blocked calls by service instance",
            "defaultName" : "Custom chart",
            "chartConfig" : {
              "legendShown" : true,
              "type" : "TOP_LIST",
              "series" : [
                {
                  "metric" : "ext:cloud.azure.microsoft_cognitiveservices.accounts.blockedcalls",
                  "aggregation" : "SUM",
                  "type" : "LINE",
                  "entityType" : "IOT",
                  "dimensions" : [
                    {
                      "id" : "0",
                      "name" : "dt.entity.custom_device",
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
          "name" : "Performance ",
          "tileType" : "HEADER",
          "configured" : true,
          "bounds" : {
            "top" : 0,
            "left" : 38,
            "width" : 1026,
            "height" : 38
          },
          "tileFilter" : {}
        },
        {
          "name" : "Data volume",
          "tileType" : "HEADER",
          "configured" : true,
          "bounds" : {
            "top" : 0,
            "left" : 1064,
            "width" : 342,
            "height" : 38
          },
          "tileFilter" : {}
        }
      ]
  })
}
