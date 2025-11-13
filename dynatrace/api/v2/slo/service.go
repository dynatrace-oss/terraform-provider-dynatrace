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

package slo

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/shutdown"

	slo "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/slo/settings"
)

func Service(credentials *rest.Credentials) settings.CRUDService[*slo.SLO] {
	return &service{credentials: credentials}
}

type service struct {
	credentials *rest.Credentials
}

func (me *service) Get(ctx context.Context, id string, v *slo.SLO) error {
	err := me.get(ctx, id, v)
	if err != nil {
		if err.Error() == "Cannot access a disabled SLO." {
			return errors.New("inaccessible")
		}
		return err
	}
	return nil
}

func (me *service) get(ctx context.Context, id string, v *slo.SLO) error {
	var err error

	client := rest.APITokenClient(me.credentials)
	req := client.Get(ctx, fmt.Sprintf("/api/v2/slo/%s?evaluate=false", url.PathEscape(id)), 200)
	if err = req.Finish(v); err != nil {
		return err
	}
	return nil
}

func (me *service) SchemaID() string {
	return "v2:environment:slo"
}

type sloList struct {
	SLOs        []*sloListEntry `json:"slo"`
	PageSize    *int32          `json:"pageSize"`
	NextPageKey *string         `json:"nextPageKey,omitempty"`
	TotalCount  *int64          `json:"totalCount"`
}

type sloListEntry struct {
	slo.SLO
	ID string `json:"id"`
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	var err error

	client := rest.APITokenClient(me.credentials)
	req := client.Get(ctx, "/api/v2/slo?pageSize=4000&sort=name&timeFrame=CURRENT&pageIdx=1&demo=false&evaluate=false", 200)
	var slos sloList
	if err = req.Finish(&slos); err != nil {
		return nil, err
	}
	stubs := api.Stubs{}
	for _, slo := range slos.SLOs {
		stubs = append(stubs, &api.Stub{ID: slo.ID, Name: slo.SLO.Name, Value: &slo.SLO})
	}

	return stubs, nil
}

func (me *service) Validate(v *slo.SLO) error {
	return nil // no endpoint for that
}

func (me *service) Create(ctx context.Context, v *slo.SLO) (*api.Stub, error) {
	// mu.Lock()
	// defer mu.Unlock()

	var err error

	var id string

	client := rest.APITokenClient(me.credentials)
	req := client.Post(ctx, "/api/v2/slo", v, 201).OnResponse(func(resp *http.Response) {
		if resp == nil {
			return
		}
		if resp.Header == nil {
			return
		}
		location := resp.Header.Get("Location")
		if len(location) > 0 {
			parts := strings.Split(location, "/")
			if len(parts) > 0 {
				id = parts[len(parts)-1]
			}
		}
	})

	retry := true
	attempts := 0
	maxattempts := 100
	for retry {
		attempts = attempts + 1
		if err = req.Finish(); err != nil {
			if attempts == maxattempts {
				return &api.Stub{ID: id, Name: v.Name}, err
			}
			if !strings.Contains(err.Error(), "calc:") && !strings.Contains(err.Error(), "Metric selector is invalid") && !strings.Contains(err.Error(), "<title>HTTP Status 400") {
				return &api.Stub{ID: id, Name: v.Name}, err
			}
			time.Sleep(2 * time.Second)
			if shutdown.System.Stopped() {
				return &api.Stub{ID: id, Name: v.Name}, errors.New("execution interrupted")
			}
		} else {
			retry = false
		}
	}
	if v.Enabled {
		length := 0
		for length == 0 {
			var slos sloList
			if err = client.Get(ctx, fmt.Sprintf("/api/v2/slo?sloSelector=%s&pageSize=10000&sort=name&timeFrame=CURRENT&pageIdx=1&demo=false&evaluate=false", url.QueryEscape(fmt.Sprintf("id(\"%s\")", id))), 200).Finish(&slos); err != nil {
				return &api.Stub{ID: id, Name: v.Name}, err
			}
			length = len(slos.SLOs)
			if length == 0 {
				time.Sleep(time.Second * 2)
				if shutdown.System.Stopped() {
					return &api.Stub{ID: id, Name: v.Name}, errors.New("execution interrupted")
				}
			}
			for _, stub := range slos.SLOs {
				v.Timeframe = stub.SLO.Timeframe
			}
		}

		retry = true
		numRequiredSuccesses := 20
		for retry {
			req = client.Get(ctx, fmt.Sprintf("/api/v2/slo/%s?evaluate=false", url.PathEscape(id)), 200)
			if err = req.Finish(v); err != nil {
				if !strings.Contains(err.Error(), "not found") {
					return &api.Stub{ID: id, Name: v.Name}, err
				}
				time.Sleep(2 * time.Second)
				if shutdown.System.Stopped() {
					return &api.Stub{ID: id, Name: v.Name}, errors.New("execution interrupted")
				}
			} else {
				numRequiredSuccesses--
				if numRequiredSuccesses < 0 {
					retry = false
				}
			}
		}
	}

	return &api.Stub{ID: id, Name: v.Name}, nil
}

func (me *service) Update(ctx context.Context, id string, v *slo.SLO) error {
	return rest.APITokenClient(me.credentials).Put(ctx, fmt.Sprintf("/api/v2/slo/%s", url.PathEscape(id)), v, 200).Finish()
}

func (me *service) Delete(ctx context.Context, id string) error {
	return rest.APITokenClient(me.credentials).Delete(ctx, fmt.Sprintf("/api/v2/slo/%s", url.PathEscape(id)), 204).Finish()
}

func (me *service) New() *slo.SLO {
	return new(slo.SLO)
}
