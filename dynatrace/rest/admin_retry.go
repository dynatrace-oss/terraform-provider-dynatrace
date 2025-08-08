package rest

import (
	"errors"
	"net/http"

	"github.com/dynatrace/dynatrace-configuration-as-code-core/api"
)

type AdminAccessRequestFn = func(adminAccess bool) (api.Response, error)

// DoWithAdminAccessRetry calls a given function with adminAccess and retries it without adminAccess in case there is a 403.
//
// returns:
//   - response: The response of the API call with adminAccess enabled if the permission is given
//     or the response of the API call without adminAccess if the permission is not given.
//   - err: Any occurring error, not related to the permission error of the adminAccess enabled call.
//   - adminAccess: The used adminAccess for the returned response.
func DoWithAdminAccessRetry(requestFn AdminAccessRequestFn) (response api.Response, err error, adminAccess bool) {
	var apiErr api.APIError
	resp, err := requestFn(true)
	if err != nil {
		if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusForbidden {
			resp, err = requestFn(false)
			return resp, err, false
		}
		return api.Response{}, err, false
	}
	return resp, nil, true
}
