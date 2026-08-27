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
)

func sampleEnvironmentEvent() environmentEvent {
	return environmentEvent{
		DtSettingsStage:  "dev",
		GithubRepository: "dynatrace-oss/terraform-provider-dynatrace",
		GithubRef:        "refs/heads/main",
		GithubEventName:  "push",
		GithubActor:      "some-actor",
		GithubRunID:      "123456",
		GithubJob:        "987654321",
	}
}

// ----- getEnvironmentEvent -----

func TestGetEnvironmentEvent_ReadsAllVariables(t *testing.T) {
	t.Setenv("DT_SETTINGS_STAGE", "dev")
	t.Setenv("GITHUB_REPOSITORY", "dynatrace-oss/terraform-provider-dynatrace")
	t.Setenv("GITHUB_REF", "refs/heads/main")
	t.Setenv("GITHUB_EVENT_NAME", "push")
	t.Setenv("GITHUB_ACTOR", "some-actor")
	t.Setenv("GITHUB_RUN_ID", "123456")
	t.Setenv("JOB_CHECK_RUN_ID", "987654321")

	assert.Equal(t, sampleEnvironmentEvent(), getEnvironmentEvent())
}

func TestGetEnvironmentEvent_UnsetVariablesYieldZeroValue(t *testing.T) {
	for _, name := range []string{
		"DT_SETTINGS_STAGE", "GITHUB_REPOSITORY", "GITHUB_REF",
		"GITHUB_EVENT_NAME", "GITHUB_ACTOR", "GITHUB_RUN_ID", "JOB_CHECK_RUN_ID",
	} {
		t.Setenv(name, "")
	}

	assert.Equal(t, environmentEvent{}, getEnvironmentEvent())
}

func TestGetEnvironmentEvent_IgnoresGithubJob(t *testing.T) {
	t.Setenv("GITHUB_JOB", "collect-test-results")
	t.Setenv("JOB_CHECK_RUN_ID", "987654321")

	assert.Equal(t, "987654321", getEnvironmentEvent().GithubJob)
}

// ----- buildResultsValue -----

func TestBuildResultsValue_Populated(t *testing.T) {
	payload, err := buildResultsValue([]result{
		{SchemaID: "builtin:audit-log", Version: "1", Status: "PASSED"},
		{SchemaID: "builtin:other", Version: "2", Status: "FAILED"},
	})

	require.NoError(t, err)
	assert.JSONEq(t, `[
		{"dt.settings.schema_id":"builtin:audit-log","dt.settings.schema_version":"1","status":"PASSED"},
		{"dt.settings.schema_id":"builtin:other","dt.settings.schema_version":"2","status":"FAILED"}
	]`, payload)
}

func TestBuildResultsValue_EmptySlice(t *testing.T) {
	payload, err := buildResultsValue([]result{})

	require.NoError(t, err)
	assert.Equal(t, "[]", payload)
}

func TestBuildResultsValue_NilSliceMarshalsToEmptyArray(t *testing.T) {
	payload, err := buildResultsValue(nil)

	require.NoError(t, err)
	assert.Equal(t, "[]", payload)
}

// ----- buildEvent -----

func TestBuildEvent_SetsConstantsAndPassesThroughEnvironment(t *testing.T) {
	ee := sampleEnvironmentEvent()
	results := []result{{SchemaID: "builtin:audit-log", Version: "1", Status: "PASSED"}}

	before := time.Now().UnixMilli()
	e, err := buildEvent(ee, results)
	after := time.Now().UnixMilli()

	require.NoError(t, err)
	assert.Equal(t, "casc.e2e.test", e.EventType)
	assert.Equal(t, "SDLC_EVENT", e.EventKind)
	assert.Equal(t, "1.0.0", e.EventVersion)
	assert.Equal(t, ee, e.environmentEvent)
	assert.GreaterOrEqual(t, e.Timestamp, before)
	assert.LessOrEqual(t, e.Timestamp, after)
	assert.JSONEq(t, `[{"dt.settings.schema_id":"builtin:audit-log","dt.settings.schema_version":"1","status":"PASSED"}]`, e.Results)
}

// ----- writeEvent -----

func TestWriteEvent_FlattensEmbeddedEnvironmentEvent(t *testing.T) {
	e, err := buildEvent(sampleEnvironmentEvent(), []result{
		{SchemaID: "builtin:audit-log", Version: "1", Status: "PASSED"},
	})
	require.NoError(t, err)
	path := filepath.Join(t.TempDir(), "test-report.json")

	require.NoError(t, writeEvent(path, e))

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	var payload map[string]any
	require.NoError(t, json.Unmarshal(data, &payload))

	// the seven environmentEvent fields must sit at the top level, not nested
	assert.Equal(t, "dev", payload["dt.settings.stage"])
	assert.Equal(t, "dynatrace-oss/terraform-provider-dynatrace", payload["github.repository"])
	assert.Equal(t, "refs/heads/main", payload["github.ref"])
	assert.Equal(t, "push", payload["github.event_name"])
	assert.Equal(t, "some-actor", payload["github.actor"])
	assert.Equal(t, "123456", payload["github.run_id"])
	assert.Equal(t, "987654321", payload["github.job"])

	assert.Equal(t, "casc.e2e.test", payload["event.type"])
	assert.Equal(t, "SDLC_EVENT", payload["event.kind"])
	assert.Equal(t, "1.0.0", payload["event.version"])
	assert.Contains(t, payload, "timestamp")

	// results is a string holding double-encoded JSON, and result.Version is
	// serialised as schema_version
	raw, ok := payload["results"].(string)
	require.True(t, ok, "results must be a JSON string, got %T", payload["results"])
	assert.JSONEq(t, `[{"dt.settings.schema_id":"builtin:audit-log","dt.settings.schema_version":"1","status":"PASSED"}]`, raw)
}

func TestWriteEvent_UnwritablePath(t *testing.T) {
	err := writeEvent(filepath.Join(t.TempDir(), "missing", "test-report.json"), event{})

	assert.Error(t, err)
}
