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

package settings

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/httpcache"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/shutdown"
)

var ExportRunning = false

func NewCRUDService[T Settings](credentials *Credentials, schemaID string, options *ServiceOptions[T]) CRUDService[T] {
	return &defaultService[T]{
		schemaID: schemaID,
		client:   httpcache.DefaultClient(credentials.URL, credentials.Token, schemaID),
		options:  options,
	}
}

type defaultService[T Settings] struct {
	schemaID string
	client   rest.Client
	options  *ServiceOptions[T]
}

func (me *defaultService[T]) Get(ctx context.Context, id string, v T) error {
	if err := me.client.Get(ctx, me.getURL(id), 200).Finish(v); err != nil {
		return err
	}
	if me.options.CompleteGet != nil {
		return me.options.CompleteGet(ctx, me.client, id, v)
	}
	return nil
}

func (me *defaultService[T]) getURL(id string) string {
	if me.options.Get != nil {
		return me.options.Get(id)
	}
	panic("service options must contain a function that provides the GET URL")
}

func (me *defaultService[T]) listURL() string {
	if me.options.List != nil {
		return me.options.List()
	}
	panic("service options must provide an URL to list records")
}

func (me *defaultService[T]) validateURL(v T) string {
	if me.options.ValidateURL != nil {
		return me.options.ValidateURL(v)
	}
	if me.options.CreateURL != nil {
		return me.options.CreateURL(v) + "/validator"
	}
	if me.options.List != nil {
		return me.options.List() + "/validator"
	}
	panic("service options must provide an URL to list records")
}

func (me *defaultService[T]) createURL(v T) string {
	if me.options.CreateURL != nil {
		return me.options.CreateURL(v)
	}
	if me.options.List != nil {
		return me.options.List()
	}
	panic("service options must provide an URL to list records")
}

func (me *defaultService[T]) updateURL(id string, v T) string {
	if me.options.UpdateURL != nil {
		return me.options.UpdateURL(id, v)
	}
	if me.options.Get != nil {
		return me.options.Get(id)
	}
	panic("service options must at least contain a function that provides the GET URL")
}

func (me *defaultService[T]) deleteURL(id string) string {
	if me.options.DeleteURL != nil {
		return me.options.DeleteURL(id)
	}
	if me.options.Get != nil {
		return me.options.Get(id)
	}
	panic("service options must at least contain a function that provides the GET URL")
}

func (me *defaultService[T]) stubs() api.RecordStubs {
	if me.options.Stubs != nil {
		stubsType := reflect.ValueOf(me.options.Stubs).Type()
		if stubsType.Kind() == reflect.Pointer {
			return reflect.New(stubsType.Elem()).Interface().(api.RecordStubs)
		}
		panic("no pointer")
	}
	return &api.StubList{}
}

func (me *defaultService[T]) List(ctx context.Context) (api.Stubs, error) {
	var err error

	req := me.client.Get(ctx, me.listURL(), 200)
	stubs := me.stubs()
	if err = req.Finish(stubs); err != nil {
		return nil, err
	}

	res := stubs.ToStubs()
	m := map[string]*api.Stub{}
	for _, stub := range res {
		m[stub.ID] = stub
	}
	res = api.Stubs{}
	for _, stub := range m {
		res = append(res, stub)
	}
	return res.ToStubs(), nil
}

func (me *defaultService[T]) Validate(ctx context.Context, v T) error {
	if me.options.HasNoValidator {
		return nil
	}
	var err error

	client := me.client
	req := client.Post(ctx, me.validateURL(v), v).Expect(204)

	if err = req.Finish(); err != nil {
		return err
	}
	return nil
}

func (me *defaultService[T]) Create(ctx context.Context, v T) (*api.Stub, error) {
	if me.options != nil && me.options.Lock != nil && me.options.Unlock != nil {
		me.options.Lock()
		defer me.options.Unlock()
	}
	stub, err := me.create(ctx, v)
	if err != nil {
		return nil, err
	}

	if me.options.CreateConfirm != 0 {
		successes := 0
		for {
			stubName := struct {
				Name string `json:"name"`
			}{}
			if err = me.client.Get(ctx, me.getURL(stub.ID), 200).Finish(&stubName); err == nil {
				// For some settings the original response doesn't deliver a name
				// In here (when confirming that the whole cluster knows the new settings) we
				// receive the whole configuration anyways.
				// Works, of course, only for settings where the name is indeed serialized with "name"
				if len(stub.Name) == 0 && len(stubName.Name) > 0 {
					stub.Name = stubName.Name
				}
				successes = successes + 1
				if successes >= me.options.CreateConfirm {
					break
				}
				time.Sleep(time.Millisecond * 200)
			} else {
				successes = 0
				time.Sleep(time.Second * 10)
			}
		}

	}

	if me.options.OnChanged != nil {
		return stub, me.options.OnChanged(ctx, me.client, stub.ID, v)
	}

	return stub, err
}

func (me *defaultService[T]) create(ctx context.Context, v T) (*api.Stub, error) {
	var err error

	if me.options != nil && me.options.Duplicates != nil {
		dupStub, dupErr := me.options.Duplicates(ctx, me, v)
		if dupErr != nil {
			return nil, dupErr
		}
		if dupStub != nil {
			return dupStub, nil
		}
	}

	client := me.client
	req := client.Post(ctx, me.createURL(v), v).Expect(200, 201)

	var stub api.Stub
	if err = req.Finish(&stub); err != nil {
		if me.options.HijackOnCreate != nil {
			var hijackedStub *api.Stub
			if hijackedStub, err = me.options.HijackOnCreate(ctx, err, me, v); err != nil {
				return nil, err
			}
			if hijackedStub != nil {
				return hijackedStub, me.Update(ctx, hijackedStub.ID, v)
			} else {
				return nil, err
			}
		} else if me.options.CreateRetry != nil {
			if modifiedPayload := me.options.CreateRetry(v, err); fmt.Sprintf("%v", modifiedPayload) != "<nil>" {
				if err = client.Post(ctx, me.createURL(modifiedPayload), modifiedPayload, 200, 201).Finish(&stub); err != nil {
					return nil, err
				}
				return (&api.Stubs{&stub}).ToStubs()[0], nil
			}
			return nil, err
		}
		return nil, err
	}
	if me.options.OnAfterCreate != nil {
		return me.options.OnAfterCreate(ctx, client, (&api.Stubs{&stub}).ToStubs()[0])
	}
	return (&api.Stubs{&stub}).ToStubs()[0], nil
}

func (me *defaultService[T]) onBeforeUpdate(id string, v T) (T, error) {
	var err error
	if me.options.OnBeforeUpdate != nil {
		var clone T
		if clone, err = Clone(v); err != nil {
			return clone, err
		}
		if err = me.options.OnBeforeUpdate(id, clone); err != nil {
			return clone, err
		}
		return clone, nil
	}
	return v, nil
}

func (me *defaultService[T]) Update(ctx context.Context, id string, v T) error {
	var err error
	if v, err = me.onBeforeUpdate(id, v); err != nil {
		return err
	}
	if err = me.update(ctx, id, v); err != nil {
		return err
	}
	if me.options.OnChanged != nil {
		return me.options.OnChanged(ctx, me.client, id, v)
	}
	return nil
}

func (me *defaultService[T]) update(ctx context.Context, id string, v T) error {
	var err error
	// some endpoints respond back initially with an internal server error
	// We're re-trying at least two more times before the update fails for good
	var retries = 3
	for retries > 0 {
		err = me.client.Put(ctx, me.updateURL(id, v), v, 204).Finish()
		if err != nil {
			if strings.Contains(err.Error(), "Internal Server Error occurred. It has been logged and will be investigated") {
				retries--
			} else {
				return err
			}
		} else {
			return nil
		}
	}
	return err
}

func (me *defaultService[T]) Delete(ctx context.Context, id string) error {
	var err error
	numRetries := 0
	for {
		if err = me.client.Delete(ctx, me.deleteURL(id)).Expect(204, 200, 404).Finish(); err != nil {
			if me.options != nil && me.options.DeleteRetry != nil {
				retry, e2 := me.options.DeleteRetry(ctx, id, err)
				if e2 != nil {
					return e2
				}
				if retry {
					return me.Delete(ctx, id)
				}
			}
			if !strings.Contains(err.Error(), "Could not delete configuration") {
				return err
			} else {
				numRetries++
				if numRetries > 100 {
					return fmt.Errorf("unable to delete '%s' even after 100 retries", id)
				}
				if shutdown.System.Stopped() {
					return nil
				}
				time.Sleep(time.Second)
			}
		} else {
			return nil
		}
	}
}

func (me *defaultService[T]) SchemaID() string {
	return me.schemaID
}
