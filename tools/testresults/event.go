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
	"time"
)

type environmentEvent struct {
	DtSettingsStage  string `json:"dt.settings.stage"`
	GithubRepository string `json:"github.repository"`
	GithubRef        string `json:"github.ref"`
	GithubEventName  string `json:"github.event_name"`
	GithubActor      string `json:"github.actor"`
	GithubRunID      string `json:"github.run_id"`
	GithubJob        string `json:"github.job"`
}

func getEnvironmentEvent() environmentEvent {
	return environmentEvent{
		DtSettingsStage:  os.Getenv("DT_SETTINGS_STAGE"),
		GithubRepository: os.Getenv("GITHUB_REPOSITORY"),
		GithubRef:        os.Getenv("GITHUB_REF"),
		GithubEventName:  os.Getenv("GITHUB_EVENT_NAME"),
		GithubActor:      os.Getenv("GITHUB_ACTOR"),
		GithubRunID:      os.Getenv("GITHUB_RUN_ID"),
		GithubJob:        os.Getenv("JOB_CHECK_RUN_ID"),
	}
}

type event struct {
	environmentEvent

	Timestamp    int64  `json:"timestamp"`
	EventType    string `json:"event.type"`
	EventKind    string `json:"event.kind"`
	EventVersion string `json:"event.version"`
	Results      string `json:"results"`
}

type result struct {
	SchemaID string `json:"dt.settings.schema_id"`
	Version  string `json:"dt.settings.schema_version"`
	Status   string `json:"status"`
}

func buildEvent(ee environmentEvent, results []result) (event, error) {
	resultsPayload, err := buildResultsValue(results)
	if err != nil {
		return event{}, err
	}

	return event{
		Timestamp:        time.Now().UnixMilli(),
		EventType:        "casc.e2e.test",
		EventKind:        "SDLC_EVENT",
		EventVersion:     "1.0.0",
		environmentEvent: ee,
		Results:          resultsPayload,
	}, nil
}

func buildResultsValue(results []result) (string, error) {
	if len(results) == 0 {
		return "[]", nil
	}

	payload, err := json.Marshal(results)
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

func writeEvent(path string, e event) error {
	payload, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return os.WriteFile(path, payload, 0o644)
}
