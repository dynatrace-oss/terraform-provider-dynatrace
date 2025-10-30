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

package pipelines_test

import (
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/metrics/pipelines/settings"
	testing2 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenPipelineMetricsPipelines(t *testing.T) {
	api.TestAcc(t)
}

func TestOpenPipelineMetricsPipelinesUnmarshal(t *testing.T) {
	entries := new(pipelines.FieldExtractionEntries)
	validEntries := []*pipelines.FieldExtractionEntry{
		{
			DefaultValue:         testing2.ToPointer("value"),
			DestinationFieldName: nil,
			SourceFieldName:      "",
		},
		{
			DefaultValue:         nil,
			DestinationFieldName: testing2.ToPointer("value"),
			SourceFieldName:      "",
		},
		{
			DefaultValue:         nil,
			DestinationFieldName: nil,
			SourceFieldName:      "value",
		},
		{
			DefaultValue:         testing2.ToPointer("value1"),
			DestinationFieldName: testing2.ToPointer("value2"),
			SourceFieldName:      "value3",
		},
	}
	validWithEmpty := pipelines.FieldExtractionEntries{
		{
			DefaultValue:         nil,
			DestinationFieldName: nil,
			SourceFieldName:      "",
		},
	}
	validWithEmpty = append(validWithEmpty, validEntries...)
	err := entries.UnmarshalHCL(testing2.MockDecoder{Elements: map[string]any{"dimension": validWithEmpty}})
	require.NoError(t, err)
	assert.Len(t, *entries, len(validEntries))
	assert.ElementsMatch(t, *entries, validEntries)
}
