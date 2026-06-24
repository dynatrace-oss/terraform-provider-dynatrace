//go:build unit

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

package iam

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalPayload(t *testing.T) {
	t.Run("nil payload yields a nil body", func(t *testing.T) {
		body, err := marshalPayload(nil)
		require.NoError(t, err)
		assert.Nil(t, body)
	})

	t.Run("value payload is marshaled to JSON", func(t *testing.T) {
		payload := map[string]any{"name": "test", "count": 3}

		body, err := marshalPayload(payload)
		require.NoError(t, err)
		require.NotNil(t, body)

		data, err := io.ReadAll(body)
		require.NoError(t, err)

		expected, err := json.Marshal(payload)
		require.NoError(t, err)
		assert.JSONEq(t, string(expected), string(data))
	})

	t.Run("unmarshalable payload returns an error", func(t *testing.T) {
		body, err := marshalPayload(make(chan int))
		require.Error(t, err)
		assert.Nil(t, body)
	})
}
