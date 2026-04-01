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

package rest

import (
	"bytes"
	"context"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/oauth2"
)

const (
	maxRetries     = 10
	defaultBackoff = 5 * time.Second
)

type RetryTransport struct {
	Transport http.RoundTripper
}

func (t *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Buffer the request body so it can be replayed on retries.
	var bodyBytes []byte
	if req.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {
			return nil, err
		}
	}

	var resp *http.Response
	var err error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Restore the body for each attempt.
		if bodyBytes != nil {
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		resp, err = t.Transport.RoundTrip(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusTooManyRequests {
			return resp, nil
		}

		if attempt == maxRetries {
			break
		}

		wait := backoff(resp, attempt)
		resp.Body.Close()
		time.Sleep(wait)
	}

	return resp, nil
}

// backoff determines how long to wait before the next retry.
// It first checks the Retry-After header; if absent it falls back to exponential backoff.
func backoff(resp *http.Response, attempt int) time.Duration {
	if ra := resp.Header.Get("Retry-After"); ra != "" {
		if seconds, err := strconv.Atoi(ra); err == nil && seconds > 0 {
			return time.Duration(seconds) * time.Second
		}
		if t, err := http.ParseTime(ra); err == nil {
			if d := time.Until(t); d > 0 {
				return d
			}
		}
	}
	// exponential backoff: defaultBackoff * 2^attempt
	return time.Duration(float64(defaultBackoff) * math.Pow(2, float64(attempt)))
}

func NewContextWithOAuthRetryClient(ctx context.Context) context.Context {
	return context.WithValue(ctx, oauth2.HTTPClient, &http.Client{
		Transport: &RetryTransport{
			Transport: http.DefaultTransport,
		},
	})
}
