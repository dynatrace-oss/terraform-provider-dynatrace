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

package openpipeline

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ClassicProcessingStage struct {
	Editable   *bool                              `json:"editable,omitempty"`
	Processors []*ClassicProcessingStageProcessor `json:"processors"`
}

type ClassicProcessingStageProcessor struct {
	sqlxProcessor *SqlxProcessor
}

func (ep ClassicProcessingStageProcessor) MarshalJSON() ([]byte, error) {
	if ep.sqlxProcessor != nil {
		return json.Marshal(ep.sqlxProcessor)
	}

	return nil, errors.New("missing ClassicProcessingStageProcessor value")
}

func (ep *ClassicProcessingStageProcessor) UnmarshalJSON(b []byte) error {
	ttype, err := ExtractType(b)
	if err != nil {
		return err
	}

	switch ttype {
	case SecurityContextProcessorType:
		sqlxProcessor := SqlxProcessor{}
		if err := json.Unmarshal(b, &sqlxProcessor); err != nil {
			return err
		}
		ep.sqlxProcessor = &sqlxProcessor

	default:
		return fmt.Errorf("unknown ClassicProcessingStageProcessor type: %s", ttype)
	}

	return nil
}
