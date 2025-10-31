//go:build unit

package rest_test

import (
	"net/http"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/stretchr/testify/assert"

	"github.com/dynatrace/dynatrace-configuration-as-code-core/api"
)

func TestDoWithAdminAccessRetry(t *testing.T) {
	type apiReturnValues struct {
		response api.Response
		error    error
	}

	t.Run("returns on first try with adminAccess", func(t *testing.T) {
		response := api.Response{StatusCode: http.StatusOK}
		call := func(_ bool) (api.Response, error) {
			return response, nil
		}
		resp, err, adminAccess := rest.DoWithAdminAccessRetry(call)
		assert.NoError(t, err)
		assert.True(t, adminAccess)
		assert.Equal(t, response, resp)
	})

	t.Run("returns on second try without adminAccess", func(t *testing.T) {
		responsesRetrySuccess := []apiReturnValues{
			{
				api.Response{},
				api.APIError{StatusCode: http.StatusForbidden},
			},
			{
				api.Response{StatusCode: 200},
				nil,
			},
		}
		calls := 0
		call := func(_ bool) (api.Response, error) {
			resp := responsesRetrySuccess[calls]
			calls++
			return resp.response, resp.error
		}
		resp, err, adminAccess := rest.DoWithAdminAccessRetry(call)
		assert.Equal(t, calls, 2)
		assert.NoError(t, err)
		assert.False(t, adminAccess)
		assert.Equal(t, responsesRetrySuccess[1].response, resp)
	})

	t.Run("errors if first response is not 403", func(t *testing.T) {
		respErr := api.APIError{StatusCode: http.StatusNotFound}
		call := func(_ bool) (api.Response, error) {
			return api.Response{}, respErr
		}
		_, err, adminAccess := rest.DoWithAdminAccessRetry(call)
		assert.Equal(t, err, respErr)
		assert.False(t, adminAccess)
	})

	t.Run("errors if second try without adminAccess errors", func(t *testing.T) {
		calls := 0
		responses := []apiReturnValues{
			{
				api.Response{},
				api.APIError{StatusCode: http.StatusForbidden},
			},
			{
				api.Response{},
				// body here just to differentiate between the other error and still have a 403 check,
				// just to be sure that sequential 403 will not lead to a loop
				api.APIError{StatusCode: http.StatusForbidden, Body: []byte("{}")},
			},
		}
		call := func(_ bool) (api.Response, error) {
			resp := responses[calls]
			calls++
			return resp.response, resp.error
		}
		_, err, adminAccess := rest.DoWithAdminAccessRetry(call)
		assert.Equal(t, calls, 2)
		assert.Equal(t, err, responses[1].error)
		assert.False(t, adminAccess)
	})
}
