/**
* @license
* Copyright 2020 Dynatrace LLC
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

package buckets

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	buckets "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/platform/buckets/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/shutdown"
	coreapi "github.com/dynatrace/dynatrace-configuration-as-code-core/api"
	bucket "github.com/dynatrace/dynatrace-configuration-as-code-core/clients/buckets"
)

func Service(credentials *rest.Credentials) settings.CRUDService[*buckets.Bucket] {
	return &service{credentials}
}

type service struct {
	credentials *rest.Credentials
}

func (me *service) client(ctx context.Context) (*bucket.Client, error) {
	platformClient, err := rest.CreatePlatformClient(ctx, me.credentials.OAuth.EnvironmentURL, me.credentials)
	if err != nil {
		return nil, err
	}
	return bucket.NewClient(platformClient), nil
}

var IGNORE_UNEXPECTED_EOF = (os.Getenv("DT_BUCKETS_IGNORE_UNEXPECTED_EOF") == "true")

func (me *service) Get(ctx context.Context, id string, v *buckets.Bucket) (err error) {
	err = me.get(ctx, id, v)
	if IGNORE_UNEXPECTED_EOF && err != nil {
		if strings.Contains(err.Error(), "unexpected EOF") {
			cfg := ctx.Value(settings.ContextKeyStateConfig)
			if stateBucket, ok := cfg.(*buckets.Bucket); ok {
				v.DisplayName = stateBucket.DisplayName
				v.Name = stateBucket.Name
				v.RetentionDays = stateBucket.RetentionDays
				v.Status = stateBucket.Status
				v.Table = stateBucket.Table
				v.Version = stateBucket.Version
				return nil
			}
		}
	}
	return err
}

func (me *service) get(ctx context.Context, id string, v *buckets.Bucket) (err error) {
	client, err := me.client(ctx)
	if err != nil {
		return err
	}
	var result bucket.Response
	if result, err = client.Get(ctx, id); err != nil {
		return err
	}
	if !result.IsSuccess() {
		return rest.Envelope(result.Data, me.credentials.OAuth.EnvironmentURL, "GET")
	}
	return json.Unmarshal(result.Data, &v)
}

func (me *service) SchemaID() string {
	return "platform:buckets"
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	client, err := me.client(ctx)
	if err != nil {
		return nil, err
	}
	result, err := client.List(ctx)
	if err != nil {
		return nil, err
	}
	var stubs api.Stubs
	for _, r := range result.All() {
		var bucket buckets.Bucket
		if err := json.Unmarshal(r, &bucket); err != nil {
			return nil, err
		}
		stubs = append(stubs, &api.Stub{ID: bucket.Name, Name: bucket.Name})
	}
	return stubs, nil
}

func (me *service) Validate(v *buckets.Bucket) error {
	return nil // no endpoint for that
}

const DefaultNumRequiredSuccesses = 10
const MinNumRequiredSuccesses = 10
const MaxNumRequiredSuccesses = 50

const DefaultMaxConfirmationRetries = 180
const MaxMaxConfirmationRetries = 360
const MinMaxConfirmationRetries = 180

func getEnv(key string, def int, min int, max int) int {
	value := os.Getenv(key)
	if len(value) == 0 {
		return def
	}
	iValue, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return def
	}
	if iValue > max {
		iValue = max
	}
	if iValue < min {
		iValue = min
	}
	return iValue
}

func (me *service) Create(ctx context.Context, v *buckets.Bucket) (stub *api.Stub, err error) {
	client, err := me.client(ctx)
	if err != nil {
		return nil, err
	}
	var data []byte
	if data, err = json.Marshal(v); err != nil {
		return nil, err
	}
	var response bucket.Response
	if response, err = client.Create(ctx, v.Name, data); err != nil {
		return nil, err
	}
	if !response.IsSuccess() {
		return nil, rest.Envelope(response.Data, me.credentials.OAuth.EnvironmentURL, "POST")
	}

	maxConfirmationRetries := getEnv("DT_BUCKETS_RETRIES", DefaultMaxConfirmationRetries, MinMaxConfirmationRetries, MaxMaxConfirmationRetries)
	numRequiredSuccesses := getEnv("DT_BUCKETS_NUM_SUCCESSES", DefaultNumRequiredSuccesses, MinNumRequiredSuccesses, MaxNumRequiredSuccesses)
	requiredSuccessesLeft := numRequiredSuccesses
	retries := 0
	var responseBucket buckets.Bucket
	for requiredSuccessesLeft > 0 || len(responseBucket.Status) == 0 || responseBucket.Status == buckets.Statuses.Creating {
		responseBucket = buckets.Bucket{}
		response, err := client.Get(ctx, v.Name)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(response.Data, &responseBucket)
		if responseBucket.Status == buckets.Statuses.Active {
			requiredSuccessesLeft--
		} else {
			requiredSuccessesLeft = numRequiredSuccesses
		}
		retries++
		if retries > maxConfirmationRetries {
			break
		}
		if shutdown.System.Stopped() {
			return &api.Stub{}, nil
		}
		time.Sleep(2 * time.Second)
	}

	return &api.Stub{Name: v.Name, ID: v.Name}, nil
}

func (me *service) Update(ctx context.Context, id string, v *buckets.Bucket) (err error) {
	client, err := me.client(ctx)
	if err != nil {
		return err
	}
	var oldBucket buckets.Bucket
	me.Get(ctx, id, &oldBucket)
	oldVersion := oldBucket.Version
	var data []byte
	if data, err = json.Marshal(v); err != nil {
		return err
	}
	var response bucket.Response
	response, err = client.Update(ctx, id, data)
	if err != nil {
		return err
	}
	if !response.IsSuccess() {
		return rest.Envelope(response.Data, me.credentials.OAuth.EnvironmentURL, "PUT")
	}
	maxConfirmationRetries := getEnv("DT_BUCKETS_RETRIES", DefaultMaxConfirmationRetries, MinMaxConfirmationRetries, MaxMaxConfirmationRetries)
	retries := 0
	for {
		var bucket buckets.Bucket
		err = me.Get(ctx, id, &bucket)
		if err == nil {
			if bucket.Version > oldVersion && len(bucket.Status) > 0 && bucket.Status != buckets.Statuses.Updating {
				break
			}
		}
		retries++
		if retries > maxConfirmationRetries {
			break
		}
		if shutdown.System.Stopped() {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return err
}

func (me *service) Delete(ctx context.Context, id string) error {
	client, err := me.client(ctx)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, id)
	if err != nil {
		return err
	}
	maxConfirmationRetries := getEnv("DT_BUCKETS_RETRIES", DefaultMaxConfirmationRetries, MinMaxConfirmationRetries, MaxMaxConfirmationRetries)
	retries := 0
	response, err := client.Get(ctx, id)
	for response.StatusCode != 404 {
		if err != nil {
			if apiErr, ok := err.(coreapi.APIError); ok {
				if apiErr.StatusCode == 404 {
					return nil
				}
			}
		}
		response, err = client.Get(ctx, id)
		retries++
		if retries > maxConfirmationRetries {
			break
		}
		if shutdown.System.Stopped() {
			return nil
		}
		if response.StatusCode == 404 {
			return nil
		}
		time.Sleep(5 * time.Second)
	}
	return err
}

func (me *service) New() *buckets.Bucket {
	return new(buckets.Bucket)
}
