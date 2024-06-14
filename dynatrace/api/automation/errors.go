package automation

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
)

/*
	{
	  "error": {
	    "message": "Cannot create a new workflow. The workflow limit according to license of 3 has been reached.",
		"code": 400,
		"details": {
		  "errorCode": "licenseLimitReached"
		}
	  }
	}
*/

type ErrorEnvelope struct {
	Err *Error `json:"error"`
}

func (e ErrorEnvelope) Error() string {
	return e.Err.Error()
}

func (e *ErrorEnvelope) Unmarshal(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	if err := json.Unmarshal(data, e); err != nil {
		return false
	}
	return e.Err != nil
}

type Error struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Details json.RawMessage `json:"details"`
}

func (e Error) ToRESTError() error {
	return rest.Error{Code: e.Code, Message: e.Error()}
}

func (e Error) Error() string {
	result := strings.TrimSpace(e.Message)
	if len(e.Details) == 0 {
		return result
	}
	if len(result) == 0 {
		return fmt.Sprintf("Status Code: %d", e.Code)
	}
	return fmt.Sprintf("%s %s", result, string(e.Details))
}
