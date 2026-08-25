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

// metadata is the ambient GitHub run context that populates the github.* event
// fields. File I/O is passed via flags; this context is read from the environment.
type metadata struct {
	organization string
	repository   string
	ref          string
	branch       string
	sha          string
	eventName    string
	actor        string
	workflow     string
	runID        string
	runNumber    string
	runAttempt   string
	runURL       string
	job          string
}

func metadataFromEnvironment() metadata {
	repository := os.Getenv("GITHUB_REPOSITORY")
	return metadata{
		organization: os.Getenv("GITHUB_REPOSITORY_OWNER"),
		repository:   repositoryName(repository),
		ref:          os.Getenv("GITHUB_REF"),
		branch:       os.Getenv("GITHUB_REF_NAME"),
		sha:          os.Getenv("GITHUB_SHA"),
		eventName:    os.Getenv("GITHUB_EVENT_NAME"),
		actor:        os.Getenv("GITHUB_ACTOR"),
		workflow:     os.Getenv("GITHUB_WORKFLOW"),
		runID:        os.Getenv("GITHUB_RUN_ID"),
		runNumber:    os.Getenv("GITHUB_RUN_NUMBER"),
		runAttempt:   os.Getenv("GITHUB_RUN_ATTEMPT"),
		runURL:       os.Getenv("GITHUB_SERVER_URL") + "/" + repository + "/actions/runs/" + os.Getenv("GITHUB_RUN_ID"),
		job:          os.Getenv("JOB_CHECK_RUN_ID"),
	}
}

// repositoryName returns the repository portion of an "owner/repo" slug,
// mirroring the ${GITHUB_REPOSITORY#*/} expansion the action previously used.
func repositoryName(repository string) string {
	for i := 0; i < len(repository); i++ {
		if repository[i] == '/' {
			return repository[i+1:]
		}
	}
	return repository
}

type Status string

const (
	StatusPass Status = "PASS"
	StatusFail Status = "FAIL"
)

type event struct {
	EventType     string  `json:"event.type"`
	EventKind     string  `json:"event.kind"`
	EventVersion  string  `json:"event.version"`
	Component     string  `json:"component"`
	Timestamp     int64   `json:"timestamp"`
	Stage         string  `json:"dt.settings.stage"`
	SchemaId      string  `json:"dt.settings.schema_id"`
	SchemaVersion string  `json:"dt.settings.schema_version"`
	Status        Status  `json:"status"`
	Error         *string `json:"error,omitempty"`
	TestName      string  `json:"test.name"`
	Org           string  `json:"github.organization"`
	Repository    string  `json:"github.repository"`
	Ref           string  `json:"github.ref"`
	Branch        string  `json:"github.branch"`
	SHA           string  `json:"github.sha"`
	EventName     string  `json:"github.event_name"`
	Actor         string  `json:"github.actor"`
	Workflow      string  `json:"github.workflow"`
	RunID         string  `json:"github.run_id"`
	RunNumber     string  `json:"github.run_number"`
	RunAttempt    string  `json:"github.run_attempt"`
	RunURL        string  `json:"github.run_url"`
	Job           string  `json:"github.job"`
}

func buildEvents(testNames []string, meta metadata) []event {
	events := make([]event, 0, len(testNames))
	for _, testName := range testNames {
		events = append(events, event{
			EventType:     "casc.e2e.test",
			EventKind:     "SDLC_EVENT",
			EventVersion:  "1.0.0",
			Component:     "terraform-provider-dynatrace",
			Timestamp:     time.Now().UnixMilli(),
			SchemaId:      "casc.e2e.test", // For demonstration purposes, using a fixed schema ID
			SchemaVersion: "1.0.0",         // For demonstration purposes, using a fixed schema version
			Stage:         "PROD",          // For demonstration purposes, using a fixed stage
			Status:        StatusPass,      // For demonstration purposes, all tests are marked as PASS
			TestName:      testName,
			Error:         nil, // For demonstration purposes, no error is provided
			Org:           meta.organization,
			Repository:    meta.repository,
			Ref:           meta.ref,
			Branch:        meta.branch,
			SHA:           meta.sha,
			EventName:     meta.eventName,
			Actor:         meta.actor,
			Workflow:      meta.workflow,
			RunID:         meta.runID,
			RunNumber:     meta.runNumber,
			RunAttempt:    meta.runAttempt,
			RunURL:        meta.runURL,
			Job:           meta.job,
		})
	}
	return events
}

func writeEvents(path string, events []event) error {
	payload, err := json.Marshal(events)
	if err != nil {
		return err
	}
	return os.WriteFile(path, payload, 0o644)
}
