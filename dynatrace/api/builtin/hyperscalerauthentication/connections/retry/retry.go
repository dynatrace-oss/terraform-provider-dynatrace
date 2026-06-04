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

package retry

import (
	"context"
	"errors"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// RetryDeadlineBuffer is reserved before the context deadline so the final retry attempt can
// complete against a still-live context. Without it, the retry budget would coincide with the
// context deadline: a request in-flight at expiry would be cancelled and surface a generic
// "context deadline exceeded" error instead of the real, actionable API error (and, for retries
// driven by eventual consistency, would mask the underlying constraint violation).
const RetryDeadlineBuffer = 5 * time.Second

// DurationUntilDeadlineOrDefault computes the retry budget: the duration until the context deadline
// minus RetryDeadlineBuffer, or the provided defaultTimeout if no deadline is set. Reserving the
// buffer ensures the retry loop gives up slightly before the context is cancelled, so the last
// attempt finishes and returns the real error rather than a context-cancellation error.
func DurationUntilDeadlineOrDefault(ctx context.Context, defaultTimeout time.Duration) time.Duration {
	dl, hasDeadline := ctx.Deadline()
	if !hasDeadline {
		// no deadline: use default. There is no external context cancellation to race against, so
		// no buffer is needed here.
		return defaultTimeout
	}

	remaining := time.Until(dl) - RetryDeadlineBuffer
	if remaining <= 0 {
		// already expired, or too little time left to retry and still reserve the buffer
		return 0
	}

	return remaining
}

// ClassifyRetryError encapsulates which errors should be retried.
// - 400 and 404 are considered retryable due to eventual consistency.
// - other 4xx are non-retryable, and non-HTTP (network) errors are non-retryable.
// - 5xx are retryable.
func ClassifyRetryError(err error) *retry.RetryError {
	if err == nil {
		return nil
	}

	var restError rest.Error
	if errors.As(err, &restError) {
		code := restError.Code
		// Retry on specific client errors that can be transient (eventual consistency).
		if code == 400 || code == 404 {
			return retry.RetryableError(err)
		}
		// Treat other 4xx as non-retryable client errors.
		if code >= 400 && code < 500 {
			return retry.NonRetryableError(err)
		}
		// 5xx and others -> retryable
		return retry.RetryableError(err)
	}
	// Non-HTTP errors (network, timeouts, context) -> non-retryable
	return retry.NonRetryableError(err)
}
