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

package workflows

import (
	"context"
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	buckets "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/platform/buckets/settings"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/api/clients"
	bucket "github.com/dynatrace/dynatrace-configuration-as-code-core/api/clients/buckets"
	"golang.org/x/oauth2/clientcredentials"
)

func Service(credentials *settings.Credentials) settings.CRUDService[*buckets.Bucket] {
	return &service{credentials}
}

type service struct {
	credentials *settings.Credentials
}

func (me *service) client() *bucket.Client {
	factory := clients.Factory().
		WithUserAgent("Dynatrace Terraform Provider").
		WithEnvironmentURL(me.credentials.Automation.EnvironmentURL).
		WithOAuthCredentials(clientcredentials.Config{
			ClientID:     me.credentials.Automation.ClientID,
			ClientSecret: me.credentials.Automation.ClientSecret,
			TokenURL:     me.credentials.Automation.TokenURL,
		})

	bucketClient, _ := factory.BucketClient()
	return bucketClient
}

func (me *service) Get(id string, v *buckets.Bucket) (err error) {
	var result bucket.Response
	if result, err = me.client().Get(context.TODO(), id); err != nil {
		return err
	}
	return json.Unmarshal(result.Data, &v)
}

func (me *service) SchemaID() string {
	return "platform:buckets"
}

func (me *service) List() (api.Stubs, error) {
	result, err := me.client().List(context.TODO())
	if err != nil {
		return nil, err
	}
	var stubs api.Stubs
	for _, r := range result.Objects {
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

func (me *service) Create(v *buckets.Bucket) (stub *api.Stub, err error) {
	var data []byte
	if data, err = json.Marshal(v); err != nil {
		return nil, err
	}
	client := me.client()
	if _, err = client.Create(context.TODO(), v.Name, data); err != nil {
		return nil, err
	}

	var responseBucket buckets.Bucket
	for responseBucket.Status == nil || string(*responseBucket.Status) == string(buckets.Statuses.Creating) {
		response, err := client.Get(context.TODO(), v.Name)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(response.Data, &responseBucket)
	}

	return &api.Stub{Name: v.Name, ID: v.Name}, nil
}

func (me *service) Update(id string, v *buckets.Bucket) (err error) {
	var data []byte
	if data, err = json.Marshal(v); err != nil {
		return err
	}
	_, err = me.client().Update(context.TODO(), id, data)
	return err
}

func (me *service) Delete(id string) error {
	_, err := me.client().Delete(context.TODO(), id)
	return err
}

func (me *service) New() *buckets.Bucket {
	return new(buckets.Bucket)
}

func (me *service) Name() string {
	return me.SchemaID()
}
