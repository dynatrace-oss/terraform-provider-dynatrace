package settings_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/pipeline/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/processors"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var minimalPipeline = settings.Pipeline{
	Kind:        "events",
	CustomID:    "customId",
	DisplayName: "displayName",
}

var sampleData = "sampleData"
var matcher = "not true"
var pipelineWithProcessor = settings.Pipeline{
	Kind:        "events",
	CustomID:    "customId",
	DisplayName: "displayName",
	Processing: &settings.Stage{
		Processors: []*processors.Processor{
			{
				Enabled:     true,
				Id:          "proc-2",
				Type:        processors.DqlProcessorType,
				Description: "my-proc-2",
				SampleData:  &sampleData,
				Matcher:     &matcher,
				Dql:         &processors.DqlAttributes{Script: "fieldsAdd true"},
			},
		},
	},
}

func TestPipeline_MarshalHCL(t *testing.T) {
	cases := []struct {
		name     string
		input    settings.Pipeline
		expected hcl.Properties
	}{
		{
			name:  "minimum-set",
			input: minimalPipeline,
			expected: hcl.Properties{
				"kind":         "events",
				"custom_id":    "customId",
				"display_name": "displayName",
			},
		},
		{
			name:  "all-set",
			input: pipelineWithProcessor,
			expected: hcl.Properties{
				"kind":         "events",
				"custom_id":    "customId",
				"display_name": "displayName",
				"processing": []interface{}{
					hcl.Properties{
						"processor": []interface{}{
							hcl.Properties{
								"enabled":     true,
								"id":          "proc-2",
								"type":        processors.DqlProcessorType,
								"description": "my-proc-2",
								"sample_data": sampleData,
								"matcher":     matcher,
								"dql": []interface{}{
									hcl.Properties{
										"script": "fieldsAdd true",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var actual = hcl.Properties{}
			err := tc.input.MarshalHCL(actual)
			assert.Equal(t, tc.expected, actual)
			assert.NoError(t, err)
		})
	}
}

func TestPipeline_UnmarshalHCL(t *testing.T) {
	s := new(settings.Pipeline).Schema()

	cases := []struct {
		name     string
		input    map[string]interface{}
		expected settings.Pipeline
	}{
		{
			name: "minimal fields set",
			input: map[string]interface{}{
				"kind":         "events",
				"custom_id":    "customId",
				"display_name": "displayName",
			},
			expected: minimalPipeline,
		},
		{
			name: "all fields set",
			input: map[string]interface{}{
				"kind":         "events",
				"custom_id":    "customId",
				"display_name": "displayName",
				"processing": []interface{}{
					map[string]interface{}{
						"processor": []interface{}{
							map[string]interface{}{
								"enabled":     true,
								"id":          "proc-2",
								"type":        processors.DqlProcessorType,
								"description": "my-proc-2",
								"sample_data": sampleData,
								"matcher":     matcher,
								"dql": []interface{}{
									map[string]interface{}{
										"script": "fieldsAdd true",
									},
								},
							},
						},
					},
				},
			},
			expected: pipelineWithProcessor,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, s, c.input)
			assert.NotNil(t, d)

			var actual settings.Pipeline
			decoder := hcl.DecoderFrom(d)
			err := actual.UnmarshalHCL(decoder)
			assert.Equal(t, c.expected, actual)
			assert.NoError(t, err)
		})
	}
}

func TestIngestSource_MarshalJSON(t *testing.T) {
	cases := []struct {
		name     string
		input    settings.Pipeline
		expected []byte
	}{
		{
			name:  "minimal-fields-set",
			input: minimalPipeline,
			expected: []byte( // The API demands that "processing" is present, even if it is null, so that seems correct.
				`{
					"customId": "customId",
					"displayName": "displayName",
					"processing": {},
					"securityContext": {},
					"costAllocation": {},
					"productAllocation": {},
					"storage": {},
					"metricExtraction": {},
					"davis": {},
					"dataExtraction": {}
				}`),
		},
		{
			name:  "pipeline-with-processor",
			input: pipelineWithProcessor,
			expected: []byte(
				`{
					"customId": "customId",
					"displayName": "displayName",
					"processing": {
						"processors": [
							{
								"id": "proc-2",
								"type": "dql",
								"matcher": "not true",
								"description": "my-proc-2",
								"sampleData": "sampleData",
								"enabled": true,
								"dql": {
									"script": "fieldsAdd true"
								}
							}
						]
					},
					"securityContext": {},
					"costAllocation": {},
					"productAllocation": {},
					"storage": {},
					"metricExtraction": {},
					"davis": {},
					"dataExtraction": {}
				}`),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := tc.input.MarshalJSON()
			require.NoError(t, err)

			var actualJSON map[string]interface{}
			err = json.Unmarshal(actual, &actualJSON)
			require.NoError(t, err)

			var expectedJSON map[string]interface{}
			err = json.Unmarshal(tc.expected, &expectedJSON)
			require.NoError(t, err)

			assert.Equal(t, expectedJSON, actualJSON)
		})
	}
}

func TestIngestSource_Validate(t *testing.T) {
	cases := []struct {
		name    string
		input   map[string]interface{}
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid",
			input: map[string]interface{}{
				"kind":         "events",
				"display_name": "displayName",
				"custom_id":    "customId",
			},
			wantErr: false,
		},
		{
			name: "valid with 1 processor set",
			input: map[string]interface{}{
				"kind":         "events",
				"custom_id":    "customId",
				"display_name": "displayName",
				"processing": []interface{}{
					map[string]interface{}{
						"processor": []interface{}{
							map[string]interface{}{
								"enabled":     true,
								"id":          "proc-2",
								"type":        processors.DqlProcessorType,
								"description": "my-proc-2",
								"sample_data": sampleData,
								"matcher":     matcher,
								"dql": []interface{}{
									map[string]interface{}{
										"script": "fieldsAdd true",
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "kind missing",
			input: map[string]interface{}{
				"display_name": "displayName",
				"custom_id":    "customId",
			},
			wantErr: true,
			errMsg:  "Missing required argument",
		},
		{
			name: "kind invalid",
			input: map[string]interface{}{
				"kind":         "invalid",
				"display_name": "displayName",
				"custom_id":    "customId",
			},
			wantErr: true,
			errMsg:  fmt.Sprintf("expected kind to be one of %q, got invalid", settings.AllowedKinds),
		},
		{
			name: "display_name missing",
			input: map[string]interface{}{
				"kind":      "events",
				"custom_id": "customId",
			},
			wantErr: true,
			errMsg:  "Missing required argument",
		},
		{
			name: "display_name present but blank",
			input: map[string]interface{}{
				"kind":         "events",
				"display_name": "",
				"custom_id":    "customId",
			},
			wantErr: true,
			errMsg:  "expected length of display_name to be in the range (1 - 500), got ",
		},
		{
			name: "display_name too long",
			input: map[string]interface{}{
				"kind":         "events",
				"display_name": strings.Repeat("a", settings.DisplayNameMaxLength+1),
				"custom_id":    "customId",
			},
			wantErr: true,
			errMsg:  "expected length of display_name to be in the range (1 - 500), got " + strings.Repeat("a", settings.DisplayNameMaxLength+1),
		},
		{
			name: "custom_id missing",
			input: map[string]interface{}{
				"kind":         "events",
				"display_name": "displayName",
			},
			wantErr: true,
			errMsg:  "Missing required argument",
		},
		{
			name: "custom_id present but blank",
			input: map[string]interface{}{
				"kind":         "events",
				"display_name": "displayName",
				"custom_id":    "",
			},
			wantErr: true,
			errMsg:  "expected length of custom_id to be in the range (4 - 100), got ",
		},
		{
			name: "custom_id too short",
			input: map[string]interface{}{
				"kind":         "events",
				"display_name": "displayName",
				"custom_id":    strings.Repeat("a", settings.CustomIDMinLength-1),
			},
			wantErr: true,
			errMsg:  "expected length of custom_id to be in the range (4 - 100), got " + strings.Repeat("a", settings.CustomIDMinLength-1),
		},
		{
			name: "custom_id too long",
			input: map[string]interface{}{
				"kind":         "events",
				"display_name": "displayName",
				"custom_id":    strings.Repeat("a", settings.CustomIDMaxLength+1),
			},
			wantErr: true,
			errMsg:  "expected length of custom_id to be in the range (4 - 100), got " + strings.Repeat("a", settings.CustomIDMaxLength+1),
		},
		{
			name: "custom_id starts with dt. ",
			input: map[string]interface{}{
				"kind":         "events",
				"display_name": "displayName",
				"custom_id":    "dt.customId",
			},
			wantErr: true,
			errMsg:  "custom_id must not start with 'dt.' or 'dynatrace.'",
		},
		{
			name: "custom_id starts with dynatrace. ",
			input: map[string]interface{}{
				"kind":         "events",
				"display_name": "displayName",
				"custom_id":    "dynatrace.customId",
			},
			wantErr: true,
			errMsg:  "custom_id must not start with 'dt.' or 'dynatrace.'",
		},
	}

	r := &schema.Resource{Schema: new(settings.Pipeline).Schema()}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			c := terraform.NewResourceConfigRaw(tt.input)

			diags := r.Validate(c)
			assert.Equal(t, tt.wantErr, diags.HasError())

			if diags.HasError() {
				assert.Equal(t, tt.errMsg, diags[0].Summary)
			}
		})
	}
}
