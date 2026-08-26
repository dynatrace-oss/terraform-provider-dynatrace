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

package wif

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// jwtWithPayload assembles a JWT out of the given raw payload JSON. Its signature segment is not a
// real signature: nothing in this package verifies it, and supplying a valid one would suggest
// otherwise.
func jwtWithPayload(payload string) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"RS256","typ":"JWT"}`))
	return header + "." + base64.RawURLEncoding.EncodeToString([]byte(payload)) + ".not-a-real-signature"
}

func TestExpiryOfReadsExpiryClaim(t *testing.T) {
	expiry, err := expiryOf(jwtWithPayload(`{"aud":"dynatrace","exp":1767225600,"iat":1767225000}`))

	require.NoError(t, err)
	assert.Equal(t, time.Unix(1767225600, 0), expiry)
}

func TestExpiryOfRejectsTokenThatIsNotAJWT(t *testing.T) {
	_, err := expiryOf("not-a-jwt")

	assert.ErrorIs(t, err, errNotAJWT)
}

func TestExpiryOfRejectsPayloadThatIsNotBase64(t *testing.T) {
	_, err := expiryOf("header.not!valid!base64.signature")

	assert.ErrorIs(t, err, errPayloadNotBase64URL)
}

func TestExpiryOfRejectsPayloadThatIsNotJSON(t *testing.T) {
	_, err := expiryOf(jwtWithPayload("this is not JSON"))

	assert.ErrorIs(t, err, errPayloadNotJSON)
}

func TestExpiryOfRejectsPayloadWithoutExpiryClaim(t *testing.T) {
	_, err := expiryOf(jwtWithPayload(`{"aud":"dynatrace","iat":1767225000}`))

	assert.ErrorIs(t, err, errNoExpiryClaim)
}
