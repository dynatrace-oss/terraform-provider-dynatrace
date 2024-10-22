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
