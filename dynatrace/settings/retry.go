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

package settings

import (
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// ClassifyConstraintRetryError returns a retryable or non-retryable error on certain conflicts (constraint violations), e.g., data-source is not up to date
// If the provided violationMessage is inside a constraint violation, a retryable error is returned
func ClassifyConstraintRetryError(err error, violationMessage string) *retry.RetryError {
	if err == nil {
		return nil
	}
	if restError, ok := errors.AsType[rest.Error](err); ok && restError.Code == http.StatusBadRequest {
		containsConflict := slices.ContainsFunc(restError.ConstraintViolations, func(cv rest.ConstraintViolation) bool {
			return strings.Contains(cv.Message, violationMessage)
		})
		if containsConflict {
			return retry.RetryableError(err)
		}
	}

	return retry.NonRetryableError(err)
}
