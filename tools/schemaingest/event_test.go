//go:build unit

/*
 * @license
 * Copyright 2026 Dynatrace LLC
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/tools/internal/githubevent"
)

func sampleGithubEvent() githubevent.Event {
	return githubevent.Event{
		Repository: "dynatrace-oss/terraform-provider-dynatrace",
		Ref:        "refs/heads/main",
		EventName:  "push",
		Actor:      "some-actor",
		RunID:      "123456",
		Job:        "987654321",
	}
}

// ----- buildResultsValue -----

func TestBuildResultsValue_EmptyYieldsEmptyArray(t *testing.T) {
	payload, err := buildResultsValue(nil)
	require.NoError(t, err)
	assert.Equal(t, "[]", payload)
}

func TestBuildResultsValue_StringifiesArray(t *testing.T) {
	results := []result{
		{SchemaID: "builtin:activegate-token", Version: "1.1", ResourceName: "dynatrace_activegate_token"},
	}
	payload, err := buildResultsValue(results)
	require.NoError(t, err)
	assert.JSONEq(t, `[{"dt.settings.schema_id":"builtin:activegate-token","dt.settings.schema_version":"1.1","tf.resource_name":"dynatrace_activegate_token"}]`, payload)
}

// ----- buildEvent -----

func TestBuildEvent_SetsFixedFields(t *testing.T) {
	before := time.Now().UnixMilli()
	e, err := buildEvent(sampleGithubEvent(), true, nil)
	after := time.Now().UnixMilli()
	require.NoError(t, err)

	assert.Equal(t, "terraform.schema_info", e.EventType)
	assert.Equal(t, "SDLC_EVENT", e.EventKind)
	assert.Equal(t, "1.0.0", e.EventVersion)
	assert.Equal(t, sampleGithubEvent(), e.Event)
	assert.True(t, e.Release)
	assert.Equal(t, "[]", e.Results)
	assert.GreaterOrEqual(t, e.Timestamp, before)
	assert.LessOrEqual(t, e.Timestamp, after)
}

// ----- writeEvent -----

func TestWriteEvent_WritesJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "event.json")
	e, err := buildEvent(sampleGithubEvent(), true, []result{
		{SchemaID: "builtin:activegate-token", Version: "1.1", ResourceName: "dynatrace_activegate_token"},
	})
	require.NoError(t, err)
	require.NoError(t, writeEvent(path, e))

	data, err := os.ReadFile(path)
	require.NoError(t, err)

	var got map[string]any
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Equal(t, "terraform.schema_info", got["event.type"])
	assert.Equal(t, "dynatrace-oss/terraform-provider-dynatrace", got["github.repository"])
	assert.Equal(t, true, got["github.release"])
	assert.Equal(t, `[{"dt.settings.schema_id":"builtin:activegate-token","dt.settings.schema_version":"1.1","tf.resource_name":"dynatrace_activegate_token"}]`, got["results"])
}
