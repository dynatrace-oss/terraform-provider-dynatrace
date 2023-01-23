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
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	slo "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/slo/settings"
)

func Service(credentials *settings.Credentials) settings.CRUDService[*slo.SLO] {
	return &service{credentials: credentials}
}

type service struct {
	credentials *settings.Credentials
}

func (me *service) Get(id string, v *slo.SLO) error {
	var err error

	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	req := client.Get(fmt.Sprintf("/api/v2/slo/%s", url.PathEscape(id)), 200)
	if err = req.Finish(v); err != nil {
		return err
	}

	numRequiredSuccesses := 10
	for numRequiredSuccesses > 0 {
		length := 0

		for length == 0 {
			req = client.Get(fmt.Sprintf("/api/v2/slo?sloSelector=id(\"%s\")&pageSize=10000&sort=name&timeFrame=CURRENT&pageIdx=1&demo=false&evaluate=false", url.QueryEscape(id)), 200)
			var slos sloList
			if err = req.Finish(&slos); err != nil {
				return err
			}
			length = len(slos.SLOs)
			if length == 0 {
				time.Sleep(time.Second * 2)
			}
			for _, stub := range slos.SLOs {
				v.Timeframe = stub.SLO.Timeframe
			}
		}
		numRequiredSuccesses--
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

func (me *service) List() (settings.Stubs, error) {
	var err error

	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	req := client.Get("/api/v2/slo?pageSize=10000&sort=name&timeFrame=CURRENT&pageIdx=1&demo=false&evaluate=false", 200)
	var slos sloList
	if err = req.Finish(&slos); err != nil {
		return nil, err
	}
	stubs := settings.Stubs{}
	for _, slo := range slos.SLOs {
		stubs = append(stubs, &settings.Stub{ID: slo.ID, Name: slo.SLO.Name, Value: &slo.SLO})
	}

	return stubs, nil
}

func (me *service) Validate(v *slo.SLO) error {
	return nil // no endpoint for that
}

func (me *service) Create(v *slo.SLO) (*settings.Stub, error) {
	var err error

	var id string

	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	req := client.Post("/api/v2/slo", v, 201).OnResponse(func(resp *http.Response) {
		location := resp.Header.Get("Location")
		if len(location) > 0 {
			parts := strings.Split(location, "/")
			if len(parts) > 0 {
				id = parts[len(parts)-1]
			}
		}
	})

	retry := true
	maxAttempts := 10
	attempts := 0
	for retry {
		attempts = attempts + 1
		if err = req.Finish(); err != nil {
			if !strings.Contains(err.Error(), "calc:") && !strings.Contains(err.Error(), "Metric selector is invalid") {
				return &settings.Stub{ID: id, Name: v.Name}, err
			}
			if attempts < maxAttempts {
				time.Sleep(2 * time.Second)
			} else {
				return nil, err
			}
		} else {
			retry = false
		}
	}
	length := 0
	for length == 0 {
		var slos sloList
		if err = client.Get(fmt.Sprintf("/api/v2/slo?sloSelector=id(\"%s\")&pageSize=10000&sort=name&timeFrame=CURRENT&pageIdx=1&demo=false&evaluate=false", url.QueryEscape(id)), 200).Finish(&slos); err != nil {
			return &settings.Stub{ID: id, Name: v.Name}, err
		}
		length = len(slos.SLOs)
		if length == 0 {
			time.Sleep(time.Second * 2)
		}
		for _, stub := range slos.SLOs {
			v.Timeframe = stub.SLO.Timeframe
		}
	}

	return &settings.Stub{ID: id, Name: v.Name}, nil
}

func (me *service) Update(id string, v *slo.SLO) error {
	return rest.DefaultClient(me.credentials.URL, me.credentials.Token).Put(fmt.Sprintf("/api/v2/slo/%s", url.PathEscape(id)), v, 200).Finish()
}

func (me *service) Delete(id string) error {
	return rest.DefaultClient(me.credentials.URL, me.credentials.Token).Delete(fmt.Sprintf("/api/v2/slo/%s", url.PathEscape(id)), 204).Finish()
}

func (me *service) New() *slo.SLO {
	return new(slo.SLO)
}

func (me *service) Name() string {
	return me.SchemaID()
}
