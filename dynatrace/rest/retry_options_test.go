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

package rest_test

import (
	"net/http"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/stretchr/testify/assert"
)

func TestRetryIfTemporaryUnavailable(t *testing.T) {
	for _, status := range []int{http.StatusTooManyRequests, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout} {
		t.Run(http.StatusText(status)+" is retried", func(t *testing.T) {
			assert.True(t, rest.RetryIfTemporaryUnavailable(&http.Response{StatusCode: status}))
		})
	}

	for _, status := range []int{http.StatusOK, http.StatusBadRequest, http.StatusForbidden, http.StatusNotFound, http.StatusInternalServerError, http.StatusNotImplemented} {
		t.Run(http.StatusText(status)+" is not retried", func(t *testing.T) {
			assert.False(t, rest.RetryIfTemporaryUnavailable(&http.Response{StatusCode: status}))
		})
	}
}
