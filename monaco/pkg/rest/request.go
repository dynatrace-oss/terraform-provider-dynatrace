/*
 * @license
 * Copyright 2023 Dynatrace LLC
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
	"math/rand"

	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const MinWaitDuration = 1 * time.Second
const MaxWaitDuration = 1 * time.Minute

func Get(client *http.Client, url string) (response Response, err error) {
	req, err := request(http.MethodGet, url)

	if err != nil {
		return response, err
	}

	return executeRequest(client, req)
}

// the name delete() would collide with the built-in function
func DeleteConfig(client *http.Client, url string, id string) error {
	fullPath := url + "/" + id
	req, err := request(http.MethodDelete, fullPath)

	if err != nil {
		return err
	}

	resp, err := executeRequest(client, req)

	if err != nil {
		return err
	}

	if resp.StatusCode == 404 {
		log.Printf("[DEBUG] No config with id '%s' found to delete (HTTP 404 response)", id)
		return nil
	}

	if !resp.Success() {
		return fmt.Errorf("failed call to DELETE %s (HTTP %d)!\n Response was:\n %s", fullPath, resp.StatusCode, string(resp.Body))
	}

	return nil
}

func Post(client *http.Client, url string, data []byte) (Response, error) {
	req, err := requestWithBody(http.MethodPost, url, bytes.NewBuffer(data))

	if err != nil {
		return Response{}, err
	}

	return executeRequest(client, req)
}

func PostMultiPartFile(client *http.Client, url string, data *bytes.Buffer, contentType string) (Response, error) {
	req, err := requestWithBody(http.MethodPost, url, data)

	if err != nil {
		return Response{}, err
	}

	req.Header.Set("Content-type", contentType)

	return executeRequest(client, req)
}

func Put(client *http.Client, url string, data []byte) (Response, error) {
	req, err := requestWithBody(http.MethodPut, url, bytes.NewBuffer(data))

	if err != nil {
		return Response{}, err
	}

	return executeRequest(client, req)
}

func request(method string, url string) (*http.Request, error) {
	return requestWithBody(method, url, nil)
}

func requestWithBody(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-type", "application/json")
	return req, nil
}

func executeRequest(client *http.Client, request *http.Request) (Response, error) {

	request.Header.Set("User-Agent", "Dynatrace Terraform Provider")

	response, err := executeWithRateLimiter(func() (Response, error) {
		resp, err := client.Do(request)
		if err != nil {
			log.Printf("[DEBUG] HTTP Request failed with Error: " + err.Error())
			return Response{}, err
		}
		defer func() {
			err = resp.Body.Close()
		}()
		body, err := io.ReadAll(resp.Body)

		return Response{
			StatusCode: resp.StatusCode,
			Body:       body,
			Headers:    resp.Header,
			Pagination: getPaginationValues(body),
		}, err
	})

	if err != nil {
		return Response{}, err
	}
	return response, nil
}

func executeWithRateLimiter(callback func() (Response, error)) (Response, error) {
	response, err := callback()
	if err != nil {
		return Response{}, err
	}

	maxIterations := 5
	curIteration := 0

	for response.StatusCode == http.StatusTooManyRequests && curIteration < maxIterations {

		sleepDuration, err := response.SleepDuration()

		if err != nil {
			// The API response didn't contain any rate limiting details. Need to generate wait time autonomously
			sleepDuration = GenerateSleepDuration(curIteration)
		}

		// That's why we need plausible min/max wait time defaults:
		sleepDuration = ApplyMinMaxDefaults(sleepDuration)

		log.Printf("[DEBUG] Rate limit reached (iteration: %d/%d). Sleeping for %s", curIteration+1, maxIterations, sleepDuration)

		time.Sleep(sleepDuration)

		// Checking again:
		curIteration++

		response, err = callback()
		if err != nil {
			return Response{}, err
		}
	}

	return response, nil
}

func microsecondsToUnixTime(timeInMicroseconds int64) time.Time {
	return time.Unix(timeInMicroseconds/1000000, (timeInMicroseconds%1000000)*1000)
}

func GenerateSleepDuration(multiplier int) (duration time.Duration) {
	if multiplier < 1 {
		multiplier = 1
	}
	addedWaitMillis := rand.Int63n(MinWaitDuration.Milliseconds())
	return MinWaitDuration + time.Duration(addedWaitMillis*int64(multiplier))
}

func ApplyMinMaxDefaults(sleepDuration time.Duration) time.Duration {
	if sleepDuration < MinWaitDuration {
		sleepDuration = MinWaitDuration
	}
	if sleepDuration > MaxWaitDuration {
		sleepDuration = MaxWaitDuration
	}
	return sleepDuration
}
