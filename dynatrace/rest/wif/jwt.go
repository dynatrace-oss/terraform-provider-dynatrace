/**
* @license
* Copyright 2026 Dynatrace LLC
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

package wif

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// jwtSegmentCount is the number of dot-separated segments of a JWT: header, payload and signature.
const jwtSegmentCount = 3

// expiryOf reads the expiry claim out of a JWT without verifying the token.
//
// No signature check is attempted, and none is missing: the provider is not the relying party for
// these tokens - the Dynatrace platform is - so there is nothing here to verify against. Only the
// payload segment is decoded; the signature segment is never decoded, inspected or logged. For the
// same reason none of the errors below embed any part of the token.
func expiryOf(token string) (time.Time, error) {
	segments := strings.Split(token, ".")
	if len(segments) != jwtSegmentCount {
		return time.Time{}, errors.New("the ID token is not a JWT: expected three dot-separated segments")
	}

	// TrimRight covers encoders that pad the segment. GitHub does not, but the cost is one call.
	payload, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(segments[1], "="))
	if err != nil {
		return time.Time{}, errors.New("the payload segment of the ID token is not valid base64url")
	}

	var claims struct {
		Expiry int64 `json:"exp"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return time.Time{}, errors.New("the payload segment of the ID token is not valid JSON")
	}
	if claims.Expiry == 0 {
		return time.Time{}, errors.New("the ID token has no `exp` claim, so the provider cannot tell when to obtain a replacement")
	}

	return time.Unix(claims.Expiry, 0), nil
}
