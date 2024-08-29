resource "dynatrace_openpipeline" "this" {
  custom_base_path = "/platform/ingest/custom/logs"
  kind             = "logs"
  editable         = true
  endpoints {
    endpoint {
      base_path    = "-"
      builtin      = true
      display_name = "OneAgent"
      editable     = true
      enabled      = true
      segment      = ""
      routing {
        type = "dynamic"
      }
    }
    endpoint {
      base_path    = "-"
      builtin      = true
      display_name = "Extensions"
      editable     = false
      enabled      = true
      segment      = ""
      routing {
        type = "dynamic"
      }
    }
    endpoint {
      base_path    = "/api/v2/otlp/v1/logs"
      builtin      = true
      display_name = "OpenTelemetry"
      editable     = false
      enabled      = true
      segment      = ""
      routing {
        type = "dynamic"
      }
    }
    endpoint {
      base_path    = "/api/v2/logs/ingest"
      builtin      = true
      display_name = "Classic Environment API"
      editable     = false
      enabled      = true
      segment      = ""
      routing {
        type = "dynamic"
      }
    }
  }
  pipelines {
    pipeline {
      default_pipeline {
        builtin      = false
        display_name = "My Pipeline"
        editable     = true
        enabled      = true
        id           = "pipeline_My_Pipeline_9555"
        data_extraction {
          editable = true
          processors {
            processor {
              bizevent_extraction_processor {
                id          = "processor_My_Business_event_Data_extraction_2926"
                description = "My Business event Data extraction"
                enabled     = true
                editable    = true
                matcher     = "true"
                event_provider {
                  type  = "field"
                  field = "some-event-field"
                }
                event_type {
                  type  = "field"
                  field = "some-field"
                }
                field_extraction {
                  semantic = "INCLUDE_ALL"
                  fields = []
                }
              }
            }
          }
        }
        metric_extraction {
          editable = true
          processors {
            processor {
              value_metric_extraction_processor {
                id          = "processor_My_Value_Metric_extraction_6134"
                description = "My Value Metric extraction"
                enabled     = true
                editable    = true
                matcher     = "true"
                metric_key  = "log.some-key"
                dimensions = ["availability"]
                field       = "some-field"
              }
            }
          }
        }
        security_context {
          editable = true
          processors {
            processor {
              security_context_processor {
                id          = "processor_My_Permission_Processor_1999"
                description = "My Permission Processor"
                enabled     = true
                editable    = true
                matcher     = "true"
                value {
                  type  = "field"
                  field = "d"
                }
              }
            }
          }
        }
        storage {
          editable              = true
          catch_all_bucket_name = "default_logs"
          processors {
            processor {
              no_storage_processor {
                id          = "processor_My_No_storage_Assignment_Processor_6580"
                description = "My No storage Assignment Processor"
                enabled     = true
                editable    = true
                matcher     = "true"
              }
            }
          }
        }
      }
    }
    pipeline {
      classic_pipeline {
        builtin         = true
        enabled         = true
        editable        = false
        id              = "default"
        settings_schema = "builtin:logmonitoring.log-dpp-rules"
        processing {
          processors {
            processor {
              sqlx_processor {
                id          = "dome-id"
                matcher     = "matcher"
                description = "sqlx"
                enabled     = true
                sqlx_script = "some-script"
              }
            }
          }
        }
      }
    }
  }
  routing {
    editable = true
    catch_all_pipeline {
      editable    = true
      pipeline_id = "default"
    }
    entries {
      entry {
        builtin     = true
        editable    = true
        enabled     = true
        matcher     = "matcher"
        note        = "note"
        pipeline_id = "p1"
      }
    }
  }
}