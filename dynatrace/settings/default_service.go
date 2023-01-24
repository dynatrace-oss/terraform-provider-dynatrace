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
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/shutdown"
)

func NewCRUDService[T Settings](credentials *Credentials, schemaID string, options *ServiceOptions[T]) CRUDService[T] {
	return &defaultService[T]{
		schemaID: schemaID,
		client:   rest.DefaultClient(credentials.URL, credentials.Token),
		options:  options,
	}
}

type defaultService[T Settings] struct {
	schemaID string
	client   rest.Client
	options  *ServiceOptions[T]
}

func (me *defaultService[T]) Get(id string, v T) error {
	if err := me.get(id, v); err != nil {
		return err
	}
	if me.options.CompleteGet != nil {
		return me.options.CompleteGet(me.client, id, v)
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

func (me *defaultService[T]) stubs() RecordStubs {
	if me.options.Stubs != nil {
		stubsType := reflect.ValueOf(me.options.Stubs).Type()
		if stubsType.Kind() == reflect.Pointer {
			return reflect.New(stubsType.Elem()).Interface().(RecordStubs)
		}
		panic("no pointer")
	}
	return &StubList{}
}

func (me *defaultService[T]) get(id string, v any) error {
	return me.client.Get(me.getURL(id), 200).Finish(v)
}

func (me *defaultService[T]) List() (Stubs, error) {
	var err error

	req := me.client.Get(me.listURL(), 200)
	stubs := me.stubs()
	if err = req.Finish(stubs); err != nil {
		return nil, err
	}

	res := stubs.ToStubs()
	return res, nil
}

func (me *defaultService[T]) Validate(v T) error {
	if me.options.HasNoValidator {
		return nil
	}
	var err error

	client := me.client
	req := client.Post(me.validateURL(v), v).Expect(204)

	if err = req.Finish(); err != nil {
		return err
	}
	return nil
}

func (me *defaultService[T]) Create(v T) (*Stub, error) {
	stub, err := me.create(v)
	if err != nil {
		return nil, err
	}

	if me.options.CreateConfirm != 0 {
		successes := 0
		for {
			stubName := struct {
				Name string `json:"name"`
			}{}
			if err = me.get(stub.ID, &stubName); err == nil {
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
		return stub, me.options.OnChanged(me.client, stub.ID, v)
	}

	return stub, err
}

func (me *defaultService[T]) create(v T) (*Stub, error) {
	var err error

	client := me.client
	req := client.Post(me.createURL(v), v).Expect(200, 201)

	var stub Stub
	if err = req.Finish(&stub); err != nil {
		if me.options.HijackOnCreate != nil {
			var hijackedStub *Stub
			if hijackedStub, err = me.options.HijackOnCreate(err, me, v); err != nil {
				return nil, err
			}
			if hijackedStub != nil {
				return hijackedStub, me.Update(hijackedStub.ID, v)
			} else {
				return nil, err
			}
		} else if me.options.CreateRetry != nil {
			if modifiedPayload := me.options.CreateRetry(v, err); (any)(modifiedPayload) != (any)(nil) {
				if err = client.Post(me.createURL(modifiedPayload), modifiedPayload, 200, 201).Finish(&stub); err != nil {
					return nil, err
				}
				return (&Stubs{&stub}).ToStubs()[0], nil
			}
			return nil, err
		}
		return nil, err
	}
	return (&Stubs{&stub}).ToStubs()[0], nil
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

func (me *defaultService[T]) Update(id string, v T) error {
	var err error
	if v, err = me.onBeforeUpdate(id, v); err != nil {
		return err
	}
	if err = me.update(id, v); err != nil {
		return err
	}
	if me.options.OnChanged != nil {
		return me.options.OnChanged(me.client, id, v)
	}
	return nil
}

func (me *defaultService[T]) update(id string, v T) error {
	return me.client.Put(me.updateURL(id, v), v, 204).Finish()
}

func (me *defaultService[T]) Delete(id string) error {
	var err error
	numRetries := 0
	for {
		if err = me.client.Delete(me.deleteURL(id)).Expect(204).Finish(); err != nil {
			if me.options != nil && me.options.DeleteRetry != nil {
				retry, e2 := me.options.DeleteRetry(id, err)
				if e2 != nil {
					return e2
				}
				if retry {
					return me.Delete(id)
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
