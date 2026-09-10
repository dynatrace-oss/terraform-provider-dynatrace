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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/tools/internal/githubevent"
)

type event struct {
	githubevent.Event

	Release      bool   `json:"github.release"`
	Timestamp    int64  `json:"timestamp"`
	EventType    string `json:"event.type"`
	EventKind    string `json:"event.kind"`
	EventVersion string `json:"event.version"`
	Results      string `json:"results"`
}

type result struct {
	SchemaID     string `json:"dt.settings.schema_id"`
	Version      string `json:"dt.settings.schema_version"`
	ResourceName string `json:"tf.resource_name"`
}

func buildEvent(ghe githubevent.Event, release bool, results []result) (event, error) {
	resultsPayload, err := buildResultsValue(results)
	if err != nil {
		return event{}, err
	}

	return event{
		Timestamp:    time.Now().UnixMilli(),
		EventType:    "terraform.schema_info",
		EventKind:    "SDLC_EVENT",
		EventVersion: "1.0.0",
		Event:        ghe,
		Release:      release,
		Results:      resultsPayload,
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
