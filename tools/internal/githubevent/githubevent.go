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

// Package githubevent holds the shared GitHub CI context attached to SDLC events.
package githubevent

import "os"

// Event carries the GitHub context read from the pipeline env variables.
type Event struct {
	Repository string `json:"github.repository"`
	Ref        string `json:"github.ref"`
	EventName  string `json:"github.event_name"`
	Actor      string `json:"github.actor"`
	RunID      string `json:"github.run_id"`
	Job        string `json:"github.job"`
}

// Get reads the GitHub context from the process environment.
func Get() Event {
	return Event{
		Repository: os.Getenv("GITHUB_REPOSITORY"),
		Ref:        os.Getenv("GITHUB_REF"),
		EventName:  os.Getenv("GITHUB_EVENT_NAME"),
		Actor:      os.Getenv("GITHUB_ACTOR"),
		RunID:      os.Getenv("GITHUB_RUN_ID"),
		Job:        os.Getenv("JOB_CHECK_RUN_ID"),
	}
}
