//go:build unit

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
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/stretchr/testify/assert"

	api2 "github.com/dynatrace/dynatrace-configuration-as-code-core/api"
)

type testCase struct {
	Name string
	Err  error
}

func TestError_IsNotFoundError_True(t *testing.T) {
	testCases := []testCase{
		{
			Name: "Not found error",
			Err:  rest.Error{Code: http.StatusNotFound},
		},
		{
			Name: "Nested not found error",
			Err:  fmt.Errorf("not found error: %w", rest.Error{Code: http.StatusNotFound}),
		},
		{
			Name: "Core not found error",
			Err:  api2.APIError{StatusCode: http.StatusNotFound},
		},
		{
			Name: "Core nested not found error",
			Err:  fmt.Errorf("not found error: %w", api2.APIError{StatusCode: http.StatusNotFound}),
		},
	}

	for _, tcs := range testCases {
		t.Run(tcs.Name, func(t *testing.T) {
			assert.True(t, rest.IsNotFoundError(tcs.Err))
		})
	}
}

func TestError_IsNotFoundError_False(t *testing.T) {
	testCases := []testCase{
		{
			Name: "Nil error",
			Err:  nil,
		},
		{
			Name: "Different error",
			Err:  rest.Error{Code: http.StatusBadRequest},
		},
		{
			Name: "Nested different error",
			Err:  fmt.Errorf("not found error: %w", errors.New("not found error")),
		},
		{
			Name: "Different core error",
			Err:  api2.APIError{StatusCode: http.StatusBadRequest},
		},
		{
			Name: "Core nested different error",
			Err:  fmt.Errorf("not found error: %w", api2.APIError{StatusCode: http.StatusBadRequest}),
		},
	}

	for _, tcs := range testCases {
		t.Run(tcs.Name, func(t *testing.T) {
			assert.False(t, rest.IsNotFoundError(tcs.Err))
		})
	}
}
