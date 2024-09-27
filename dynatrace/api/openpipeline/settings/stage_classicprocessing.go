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
