//go:build unit

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

package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/retry"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/stretchr/testify/assert"
)

// TestDurationUntilDeadlineOrDefault_NoDeadline tests that DurationUntilDeadlineOrDefault returns the default timeout when the context has no deadline.
func TestDurationUntilDeadlineOrDefault_NoDeadline(t *testing.T) {
	defaultTimeout := 2 * time.Minute
	ctx := context.Background()

	retryTimeout := retry.DurationUntilDeadlineOrDefault(ctx, defaultTimeout)
	assert.Equal(t, defaultTimeout, retryTimeout)
}

// TestDurationUntilDeadlineOrDefault_WithDeadline_Plenty tests that DurationUntilDeadlineOrDefault returns the remaining duration when the context has a deadline set well in the future.
func TestDurationUntilDeadlineOrDefault_WithDeadline_Plenty(t *testing.T) {
	// caller provides 5 minutes; retryTimeout should be ~5 minutes
	parent := context.Background()
	deadline := time.Now().Add(5 * time.Minute)
	ctxWithDL, cancelParent := context.WithDeadline(parent, deadline)
	defer cancelParent()

	retryTimeout := retry.DurationUntilDeadlineOrDefault(ctxWithDL, 2*time.Minute)

	// expect ~5 minutes
	expected := 5 * time.Minute
	assert.InDelta(t, expected.Milliseconds(), retryTimeout.Milliseconds(), 200, "expected approximately %v, got %v", expected, retryTimeout)

}

// TestDurationUntilDeadlineOrDefault_ExpiredDeadline tests that wDurationUntilDeadlineOrDefault returns zero when the context has an expired deadline.
func TestDurationUntilDeadlineOrDefault_ExpiredDeadline(t *testing.T) {
	ctxWithExpiredDL, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()

	retryTimeout := retry.DurationUntilDeadlineOrDefault(ctxWithExpiredDL, 2*time.Minute)
	assert.Zero(t, retryTimeout, "expected zero retry timeout for expired deadline")
}

// TestClassifyRetryError tests that the ClassifyRetryError function correctly identifies retryable errors.
func TestClassifyRetryError(t *testing.T) {
	tests := []struct {
		err         error
		isRetryable bool
	}{
		{err: rest.Error{Code: 400, Message: "bad"}, isRetryable: true},
		{err: rest.Error{Code: 404, Message: "notfound"}, isRetryable: true},
		{err: rest.Error{Code: 401, Message: "unauth"}, isRetryable: false},
		{err: rest.Error{Code: 500, Message: "server"}, isRetryable: true},
		{err: errors.New("network failure"), isRetryable: false},
	}

	for _, tc := range tests {
		t.Run(tc.err.Error(), func(t *testing.T) {
			retryError := retry.ClassifyRetryError(tc.err)
			assert.Equal(t, tc.isRetryable, retryError.Retryable)
		})
	}
}
