/**
* @license
* Copyright 2020 Dynatrace LLC
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

package workflows

type Status string

var Statuses = struct {
	Success Status
	Error   Status
	Any     Status
	Ok      Status // Success or Skipped
	NOk     Status // Error or Cancelled
}{
	"SUCCESS",
	"ERROR",
	"ANY",
	"OK",
	"NOK",
}

type State string

var States = struct {
	Running   State
	Success   State
	Error     State
	Waiting   State
	Idle      State
	Paused    State
	Delayed   State
	Cancelled State
}{
	"RUNNING",
	"SUCCESS",
	"ERROR",
	"WAITING",
	"IDLE",
	"PAUSED",
	"DELAYED",
	"CANCELLED",
}

type Else string

var Elses = struct {
	Success Else
	Error   Else
}{
	"SKIP",
	"STOP",
}

type ScheduleInputType string

var ScheduleInputTypes = struct {
	Static ScheduleInputType
}{
	"static",
}

type ScheduleTriggerType string

var ScheduleTriggerTypes = struct {
	Cron     ScheduleTriggerType
	Time     ScheduleTriggerType
	Interval ScheduleTriggerType
}{
	"cron",
	"time",
	"interval",
}

type EventTriggerType string

func (me EventTriggerType) IsKnown() bool {
	switch me {
	case EventTriggerTypes.DavisEvent:
		return true
	case EventTriggerTypes.DavisProblem:
		return true
	case EventTriggerTypes.Event:
		return true
	default:
		return false
	}
}

var EventTriggerTypes = struct {
	Event        EventTriggerType
	DavisProblem EventTriggerType
	DavisEvent   EventTriggerType
}{
	"event",
	"davis-problem",
	"davis-event",
}

type EntityTagsMatch string

var EntityTagsMatches = struct {
	All EntityTagsMatch
	Any EntityTagsMatch
}{
	"all",
	"any",
}

type EventType string

var EventTypes = struct {
	Events    EventType
	BizEvents EventType
}{
	"events",
	"bizevents",
}

type TriggerType string

var TriggerTypes = struct {
	Schedule TriggerType
	Event    TriggerType
	Manual   TriggerType
}{
	"Schedule",
	"Event",
	"Manual",
}
