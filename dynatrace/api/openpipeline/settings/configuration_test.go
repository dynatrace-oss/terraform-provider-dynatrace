/**
* @license
* Copyright 2025 Dynatrace LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

//go:build unit

package openpipeline

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfiguration_MarshallUnmarshallJSON(t *testing.T) {
	tests := []struct {
		name      string
		givenJSON string
	}{
		{
			name: "basic data",
			givenJSON: `{"id": "events",
  "editable": true,
  "version": "1725444503541-d6dc888848df0116",
  "customBasePath": "/platform/ingest/custom/events",
  "endpoints": [],
  "pipelines": [],
  "routing": {
    "catchAllPipeline": {
      "editable": false,
      "pipelineId": "default"
    },
    "entries": []
  }
}`,
		},
		{
			name: "routing",
			givenJSON: `{"id": "events",
  "editable": true,
  "version": "1725444503541-d6dc888848df0116",
  "customBasePath": "/platform/ingest/custom/events",
  "endpoints": [],
  "pipelines": [],
  "routing": {
    "catchAllPipeline": {
      "editable": false,
      "pipelineId": "default"
    },
    "entries": [
      {
        "enabled": true,
        "editable": true,
        "builtin": false,
        "pipelineId": "default",
        "matcher": "not true ",
        "note": "Some dynamic route"
      },
      {
        "enabled": true,
        "editable": true,
        "builtin": false,
        "pipelineId": "different for above",
        "matcher": "not true ",
        "note": "Some dynamic route"
      }
    ]
  }
}`,
		},
		{
			name: "endpoints (w/o processors)",
			givenJSON: `{"id": "events",
  "editable": true,
  "version": "1725444503541-d6dc888848df0116",
  "customBasePath": "/platform/ingest/custom/events",
  "endpoints": [
    {
      "enabled": true,
      "builtin": false,
      "editable": true,
      "basePath": "/platform/ingest/custom/events",
      "segment": "test",
      "displayName": "My Events Source",
      "defaultBucket": "default_events",
      "processors": [],
      "routing": {
        "type": "static",
        "pipelineId": "default"
      }
    }
  ],
  "pipelines": [],
  "routing": {
    "catchAllPipeline": {
      "editable": false,
      "pipelineId": "default"
    },
    "entries": []
  }
}`,
		},
		{
			name: "endpoints with DQL processor",
			givenJSON: `{"id": "events",
  "editable": true,
  "version": "1725444503541-d6dc888848df0116",
  "customBasePath": "/platform/ingest/custom/events",
  "endpoints": [
    {
      "enabled": true,
      "builtin": false,
      "editable": true,
      "basePath": "/platform/ingest/custom/events",
      "segment": "test",
      "displayName": "My Events Source",
      "defaultBucket": "default_events",
      "processors": [
        {
          "type": "dql",
          "enabled": true,
          "editable": true,
          "id": "processor_Fetch_logs_5354",
          "description": "Fetch logs",
          "matcher": "true",
          "sampleData": "fetch logs",
          "dqlScript": "fields not true "
        }
      ],
      "routing": {
        "type": "static",
        "pipelineId": "default"
      }
    }
  ],
  "pipelines": [],
  "routing": {
    "catchAllPipeline": {
      "editable": false,
      "pipelineId": "default"
    },
    "entries": []
  }
}`,
		},
		{
			name: "endpoints with fieldsAdd processor",
			givenJSON: `{"id": "events",
  "editable": true,
  "version": "1725444503541-d6dc888848df0116",
  "customBasePath": "/platform/ingest/custom/events",
  "endpoints": [
    {
      "enabled": true,
      "builtin": false,
      "editable": true,
      "basePath": "/platform/ingest/custom/events",
      "segment": "test",
      "displayName": "My Events Source",
      "defaultBucket": "default_events",
      "processors": [
        {
          "type": "fieldsAdd",
          "enabled": true,
          "editable": true,
          "id": "processor_Some_add_field_processor_1382",
          "description": "Some add field processor",
          "matcher": "true",
          "fields": [
            {
              "name": "hi",
              "value": "bob"
            }
          ]
        }
      ],
      "routing": {
        "type": "static",
        "pipelineId": "default"
      }
    }
  ],
  "pipelines": [],
  "routing": {
    "catchAllPipeline": {
      "editable": false,
      "pipelineId": "default"
    },
    "entries": []
  }
}`,
		},
		{
			name: "endpoints with fieldsRemove processor",
			givenJSON: `{"id": "events",
  "editable": true,
  "version": "1725444503541-d6dc888848df0116",
  "customBasePath": "/platform/ingest/custom/events",
  "endpoints": [
    {
      "enabled": true,
      "builtin": false,
      "editable": true,
      "basePath": "/platform/ingest/custom/events",
      "segment": "test",
      "displayName": "My Events Source",
      "defaultBucket": "default_events",
      "processors": [
        {
          "type": "fieldsRemove",
          "enabled": true,
          "editable": true,
          "id": "processor_Some_remove_field_processor_7521",
          "description": "Some remove field processor",
          "matcher": "true",
          "fields": [
            "hi"
          ]
        }
      ],
      "routing": {
        "type": "static",
        "pipelineId": "default"
      }
    }
  ],
  "pipelines": [],
  "routing": {
    "catchAllPipeline": {
      "editable": false,
      "pipelineId": "default"
    },
    "entries": []
  }
}`,
		},
		{
			name: "endpoints with fieldsRename processor",
			givenJSON: `{"id": "events",
  "editable": true,
  "version": "1725444503541-d6dc888848df0116",
  "customBasePath": "/platform/ingest/custom/events",
  "endpoints": [
    {
      "enabled": true,
      "builtin": false,
      "editable": true,
      "basePath": "/platform/ingest/custom/events",
      "segment": "test",
      "displayName": "My Events Source",
      "defaultBucket": "default_events",
      "processors": [
        {
          "type": "fieldsRename",
          "enabled": true,
          "editable": true,
          "id": "processor_Some_rename_processor_2467",
          "description": "Some rename processor",
          "matcher": "true",
          "fields": [
            {
              "fromName": "sally",
              "toName": "bob"
            }
          ]
        }
      ],
      "routing": {
        "type": "static",
        "pipelineId": "default"
      }
    }
  ],
  "pipelines": [],
  "routing": {
    "catchAllPipeline": {
      "editable": false,
      "pipelineId": "default"
    },
    "entries": []
  }
}`,
		},
		{
			name: "pipeline plain",
			givenJSON: `{"id": "events",
  "editable": true,
  "version": "1725444503541-d6dc888848df0116",
  "customBasePath": "/platform/ingest/custom/events",
  "endpoints": [],
  "pipelines": [
    {
      "type": "default",
      "enabled": true,
      "editable": true,
      "id": "pipeline_My_custom_pipeline_6592",
      "builtin": false,
      "displayName": "My custom pipeline",
      "storage": {
        "editable": true,
        "catchAllBucketName": "default_events",
        "processors": []
      },
      "securityContext": {
        "editable": true,
        "processors": []
      },
      "metricExtraction": {
        "editable": true,
        "processors": []
      },
      "dataExtraction": {
        "editable": true,
        "processors": []
      },
      "processing": {
        "editable": true,
        "processors": []
      }
    }
  ],
  "routing": {
    "catchAllPipeline": {
      "editable": false,
      "pipelineId": "default"
    },
    "entries": []
  }
}`,
		},
		{
			name: "pipeline with processors",
			givenJSON: `{"id": "events",
  "editable": true,
  "version": "1725444503541-d6dc888848df0116",
  "customBasePath": "/platform/ingest/custom/events",
  "endpoints": [],
  "pipelines": [
    {
      "type": "default",
      "enabled": true,
      "editable": true,
      "id": "pipeline_My_custom_pipeline_6592",
      "builtin": false,
      "displayName": "My custom pipeline",
      "storage": {
        "editable": true,
        "catchAllBucketName": "default_events",
        "processors": []
      },
      "securityContext": {
        "editable": true,
        "processors": []
      },
      "metricExtraction": {
        "editable": true,
        "processors": []
      },
      "dataExtraction": {
        "editable": true,
        "processors": []
      },
      "processing": {
        "editable": true,
        "processors": [
          {
            "type": "dql",
            "enabled": true,
            "editable": true,
            "id": "processor_My_DQL_Processor_8501",
            "description": "My DQL Processor",
            "matcher": "true",
            "dqlScript": "fieldsAdd true "
          },
          {
            "type": "fieldsAdd",
            "enabled": true,
            "editable": true,
            "id": "processor_My_Fields_Add_1538",
            "description": "My Fields Add",
            "matcher": "true",
            "fields": [
              {
                "name": "Bob",
                "value": "hello"
              }
            ]
          },
          {
            "type": "fieldsRename",
            "enabled": true,
            "editable": true,
            "id": "processor_My_Rename_Field_7609",
            "description": "My Rename Field",
            "matcher": "true",
            "fields": [
              {
                "fromName": "Mary",
                "toName": "Bob"
              }
            ]
          },
          {
            "type": "fieldsRemove",
            "enabled": true,
            "editable": true,
            "id": "processor_My_Fields_Remove_6776",
            "description": "My Fields Remove",
            "matcher": "true",
            "fields": [
              "Bob"
            ]
          }
        ]
      }
    }
  ],
  "routing": {
    "catchAllPipeline": {
      "editable": false,
      "pipelineId": "default"
    },
    "entries": []
  }
}`,
		},
		{
			name: "pipeline with dataExtraction",
			givenJSON: `{"id": "events",
  "editable": true,
  "version": "1725444503541-d6dc888848df0116",
  "customBasePath": "/platform/ingest/custom/events",
  "endpoints": [],
  "pipelines": [
    {
      "type": "default",
      "enabled": true,
      "editable": true,
      "id": "pipeline_My_custom_pipeline_6592",
      "builtin": false,
      "displayName": "My custom pipeline",
      "storage": {
        "editable": true,
        "catchAllBucketName": "default_events",
        "processors": []
      },
      "securityContext": {
        "editable": true,
        "processors": []
      },
      "metricExtraction": {
        "editable": true,
        "processors": []
      },
      "dataExtraction": {
        "editable": true,
        "processors": [
          {
            "type": "davis",
            "enabled": true,
            "editable": true,
            "id": "processor_My_Davis_Event_4903",
            "description": "My Davis Event",
            "matcher": "true",
            "properties": [
              {
                "key": "event.type",
                "value": "CUSTOM_ALERT"
              },
              {
                "key": "event.name",
                "value": "event"
              },
              {
                "key": "event.description",
                "value": "Some description"
              }
            ]
          }
        ]
      },
      "processing": {
        "editable": true,
        "processors": []
      }
    }
  ],
  "routing": {
    "catchAllPipeline": {
      "editable": false,
      "pipelineId": "default"
    },
    "entries": []
  }
}`,
		},
		{
			name: "pipeline with metricExtraction",
			givenJSON: `{"id": "events",
  "editable": true,
  "version": "1725444503541-d6dc888848df0116",
  "customBasePath": "/platform/ingest/custom/events",
  "endpoints": [],
  "pipelines": [
    {
      "type": "default",
      "enabled": true,
      "editable": true,
      "id": "pipeline_My_custom_pipeline_6592",
      "builtin": false,
      "displayName": "My custom pipeline",
      "storage": {
        "editable": true,
        "catchAllBucketName": "default_events",
        "processors": []
      },
      "securityContext": {
        "editable": true,
        "processors": []
      },
      "metricExtraction": {
        "editable": true,
        "processors": [
          {
            "type": "valueMetric",
            "enabled": true,
            "editable": true,
            "id": "processor_My_Value_Metric_9760",
            "description": "My Value Metric",
            "matcher": "true",
            "metricKey": "events.bob",
            "dimensions": [
              "availability"
            ],
            "field": "Bob"
          },
          {
            "type": "counterMetric",
            "enabled": true,
            "editable": true,
            "id": "processor_My_Counter_Metric_8495",
            "description": "My Counter Metric",
            "matcher": "true",
            "metricKey": "events.counter",
            "dimensions": [
              "availability",
              "custom"
            ]
          }
        ]
      },
      "dataExtraction": {
        "editable": true,
        "processors": []
      },
      "processing": {
        "editable": true,
        "processors": []
      }
    }
  ],
  "routing": {
    "catchAllPipeline": {
      "editable": false,
      "pipelineId": "default"
    },
    "entries": []
  }
}`,
		},
		{
			name: "pipeline with securityContext",
			givenJSON: `{"id": "events",
  "editable": true,
  "version": "1725444503541-d6dc888848df0116",
  "customBasePath": "/platform/ingest/custom/events",
  "endpoints": [],
  "pipelines": [
    {
      "type": "default",
      "enabled": true,
      "editable": true,
      "id": "pipeline_My_custom_pipeline_6592",
      "builtin": false,
      "displayName": "My custom pipeline",
      "storage": {
        "editable": true,
        "catchAllBucketName": "default_events",
        "processors": []
      },
      "securityContext": {
        "editable": true,
        "processors": [
          {
            "type": "securityContext",
            "enabled": true,
            "editable": true,
            "id": "processor_My_Security_Context_8193",
            "description": "My Security Context",
            "matcher": "true",
            "value": {
              "type": "field",
              "field": "Bob"
            }
          }
        ]
      },
      "metricExtraction": {
        "editable": true,
        "processors": []
      },
      "dataExtraction": {
        "editable": true,
        "processors": []
      },
      "processing": {
        "editable": true,
        "processors": []
      }
    }
  ],
  "routing": {
    "catchAllPipeline": {
      "editable": false,
      "pipelineId": "default"
    },
    "entries": []
  }
}`,
		},
		{
			name: "pipeline with storage",
			givenJSON: `{"id": "events",
  "editable": true,
  "version": "1725444503541-d6dc888848df0116",
  "customBasePath": "/platform/ingest/custom/events",
  "endpoints": [],
  "pipelines": [
    {
      "type": "default",
      "enabled": true,
      "editable": true,
      "id": "pipeline_My_custom_pipeline_6592",
      "builtin": false,
      "displayName": "My custom pipeline",
      "storage": {
        "editable": true,
        "catchAllBucketName": "default_events",
        "processors": [
          {
            "type": "bucketAssignment",
            "enabled": true,
            "editable": true,
            "id": "processor_My_Bucket_Assignment_5491",
            "description": "My Bucket Assignment",
            "matcher": "true",
            "bucketName": "default_events"
          },
          {
            "type": "noStorage",
            "enabled": true,
            "editable": true,
            "id": "processor_My_No_Storage_Assignment_9916",
            "description": "My No Storage Assignment",
            "matcher": "true"
          }
        ]
      },
      "securityContext": {
        "editable": true,
        "processors": []
      },
      "metricExtraction": {
        "editable": true,
        "processors": []
      },
      "dataExtraction": {
        "editable": true,
        "processors": []
      },
      "processing": {
        "editable": true,
        "processors": []
      }
    }
  ],
  "routing": {
    "catchAllPipeline": {
      "editable": false,
      "pipelineId": "default"
    },
    "entries": []
  }
}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var d Configuration

			err := json.Unmarshal([]byte(tc.givenJSON), &d)
			require.NoError(t, err)

			actual, err := json.Marshal(d)
			require.NoError(t, err)

			require.JSONEq(t, tc.givenJSON, string(actual))
		})
	}
}

// TestConfiguration_UnmarshalJSON allow us to expect unmarshalling of JSON into Configuration struct.
// Need to expand this into an actual test.
func TestConfiguration_UnmarshalJSON(t *testing.T) {
	j := `{
  "id": "events",
  "editable": true,
  "version": "1725444503541-d6dc888848df0116",
  "customBasePath": "/platform/ingest/custom/events",
  "endpoints": [
    {
      "enabled": true,
      "builtin": false,
      "editable": true,
      "basePath": "/platform/ingest/custom/events",
      "segment": "test",
      "displayName": "My Events Source",
      "defaultBucket": "default_events",
      "processors": [
        {
          "type": "dql",
          "enabled": true,
          "editable": true,
          "id": "processor_Fetch_logs_5354",
          "description": "Fetch logs",
          "matcher": "true",
          "sampleData": "fetch logs",
          "dqlScript": "fields not true "
        },
        {
          "type": "fieldsAdd",
          "enabled": true,
          "editable": true,
          "id": "processor_Some_add_field_processor_1382",
          "description": "Some add field processor",
          "matcher": "true",
          "fields": [
            {
              "name": "hi",
              "value": "bob"
            }
          ]
        },
        {
          "type": "fieldsRemove",
          "enabled": true,
          "editable": true,
          "id": "processor_Some_remove_field_processor_7521",
          "description": "Some remove field processor",
          "matcher": "true",
          "fields": [
            "hi"
          ]
        },
        {
          "type": "fieldsRename",
          "enabled": true,
          "editable": true,
          "id": "processor_Some_rename_processor_2467",
          "description": "Some rename processor",
          "matcher": "true",
          "fields": [
            {
              "fromName": "sally",
              "toName": "bob"
            }
          ]
        }
      ],
      "routing": {
        "type": "static",
        "pipelineId": "default"
      }
    },
    {
      "enabled": true,
      "builtin": true,
      "editable": false,
      "basePath": "/platform/ingest/v1/events",
      "segment": "",
      "displayName": "Default API",
      "processors": [],
      "routing": {
        "type": "dynamic"
      }
    }
  ],
  "pipelines": [
    {
      "type": "default",
      "enabled": true,
      "editable": true,
      "id": "pipeline_My_custom_pipeline_6592",
      "builtin": false,
      "displayName": "My custom pipeline",
      "storage": {
        "editable": true,
        "catchAllBucketName": "default_events",
        "processors": [
          {
            "type": "bucketAssignment",
            "enabled": true,
            "editable": true,
            "id": "processor_My_Bucket_Assignment_5491",
            "description": "My Bucket Assignment",
            "matcher": "true",
            "bucketName": "default_events"
          },
          {
            "type": "noStorage",
            "enabled": true,
            "editable": true,
            "id": "processor_My_No_Storage_Assignment_9916",
            "description": "My No Storage Assignment",
            "matcher": "true"
          }
        ]
      },
      "securityContext": {
        "editable": true,
        "processors": [
          {
            "type": "securityContext",
            "enabled": true,
            "editable": true,
            "id": "processor_My_Security_Context_8193",
            "description": "My Security Context",
            "matcher": "true",
            "value": {
              "type": "field",
              "field": "Bob"
            }
          }
        ]
      },
      "metricExtraction": {
        "editable": true,
        "processors": [
          {
            "type": "valueMetric",
            "enabled": true,
            "editable": true,
            "id": "processor_My_Value_Metric_9760",
            "description": "My Value Metric",
            "matcher": "true",
            "metricKey": "events.bob",
            "dimensions": [
              "availability"
            ],
            "field": "Bob"
          },
          {
            "type": "counterMetric",
            "enabled": true,
            "editable": true,
            "id": "processor_My_Counter_Metric_8495",
            "description": "My Counter Metric",
            "matcher": "true",
            "metricKey": "events.counter",
            "dimensions": [
              "availability",
              "custom"
            ]
          }
        ]
      },
      "dataExtraction": {
        "editable": true,
        "processors": [
          {
            "type": "davis",
            "enabled": true,
            "editable": true,
            "id": "processor_My_Davis_Event_4903",
            "description": "My Davis Event",
            "matcher": "true",
            "properties": [
              {
                "key": "event.type",
                "value": "CUSTOM_ALERT"
              },
              {
                "key": "event.name",
                "value": "event"
              },
              {
                "key": "event.description",
                "value": "Some description"
              }
            ]
          }
        ]
      },
      "processing": {
        "editable": true,
        "processors": [
          {
            "type": "dql",
            "enabled": true,
            "editable": true,
            "id": "processor_My_DQL_Processor_8501",
            "description": "My DQL Processor",
            "matcher": "true",
            "dqlScript": "fieldsAdd true "
          },
          {
            "type": "fieldsAdd",
            "enabled": true,
            "editable": true,
            "id": "processor_My_Fields_Add_1538",
            "description": "My Fields Add",
            "matcher": "true",
            "fields": [
              {
                "name": "Bob",
                "value": "hello"
              }
            ]
          },
          {
            "type": "fieldsRename",
            "enabled": true,
            "editable": true,
            "id": "processor_My_Rename_Field_7609",
            "description": "My Rename Field",
            "matcher": "true",
            "fields": [
              {
                "fromName": "Mary",
                "toName": "Bob"
              }
            ]
          },
          {
            "type": "fieldsRemove",
            "enabled": true,
            "editable": true,
            "id": "processor_My_Fields_Remove_6776",
            "description": "My Fields Remove",
            "matcher": "true",
            "fields": [
              "Bob"
            ]
          }
        ]
      }
    },
    {
      "type": "default",
      "enabled": true,
      "editable": false,
      "id": "default",
      "builtin": true,
      "displayName": "events",
      "storage": {
        "editable": false,
        "catchAllBucketName": "default_events",
        "processors": []
      },
      "securityContext": {
        "editable": false,
        "processors": []
      },
      "metricExtraction": {
        "editable": false,
        "processors": []
      },
      "dataExtraction": {
        "editable": false,
        "processors": []
      },
      "processing": {
        "editable": false,
        "processors": []
      }
    }
  ],
  "routing": {
    "editable": true,
    "catchAllPipeline": {
      "editable": false,
      "pipelineId": "default"
    },
    "entries": [
      {
        "enabled": true,
        "editable": true,
        "builtin": false,
        "pipelineId": "default",
        "matcher": "not true ",
        "note": "Some dynamic route"
      }
    ]
  }
}
`

	d := Configuration{}
	err := json.Unmarshal([]byte(j), &d)
	require.NoError(t, err)

	b, err := json.Marshal(d)

	require.NoError(t, err)
	s := string(b)

	require.JSONEq(t, j, s)
}
