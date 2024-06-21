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

package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache/tar"
)

type crudService[T settings.Settings] struct {
	mu        sync.Mutex
	service   settings.CRUDService[T]
	folder    string
	index     *stubIndex
	tarFolder *tar.Folder
}

func (me *crudService[T]) init() error {
	if me.index != nil {
		return nil
	}
	me.index = &stubIndex{Stubs: api.Stubs{}, IDs: map[string]*api.Stub{}, Complete: false}
	os.MkdirAll(me.folder, os.ModePerm)
	tarFolder, complete, err := tar.New(path.Join(me.folder, "data"))
	if err != nil {
		return err
	}
	me.index.Complete = complete

	me.tarFolder = tarFolder
	stubs, err := me.tarFolder.List()
	if err != nil {
		return err
	}
	for _, stub := range stubs {
		me.index.Add(stub.ID, stub.Name)
	}
	return nil
}

func (me *crudService[T]) Create(ctx context.Context, v T) (*api.Stub, error) {
	me.mu.Lock()
	defer me.mu.Unlock()

	if err := me.init(); err != nil {
		return nil, err
	}

	if mode == ModeOffline {
		return nil, errors.New("modifications not allowed in offline mode")
	}
	var err error
	var stub *api.Stub
	if stub, err = me.service.Create(ctx, v); err != nil {
		return nil, err
	}
	me.index.Add(stub.ID, stub.Name)
	if err := me.tarFolder.Save(*stub, nil); err != nil {
		return nil, err
	}
	return stub, nil
}

func (me *crudService[T]) Delete(ctx context.Context, id string) error {
	me.mu.Lock()
	defer me.mu.Unlock()

	if err := me.init(); err != nil {
		return err
	}

	if mode == ModeOffline {
		return errors.New("modifications not allowed in offline mode")
	}
	if err := me.service.Delete(ctx, id); err != nil {
		return err
	}

	me.index.Remove(id)
	return me.tarFolder.Delete(id)
}

func (me *crudService[T]) List(ctx context.Context) (api.Stubs, error) {
	me.mu.Lock()
	defer me.mu.Unlock()

	return me.list(ctx, true)
}

func (me *crudService[T]) ListNoValues(ctx context.Context) (api.Stubs, error) {
	me.mu.Lock()
	defer me.mu.Unlock()

	return me.list(ctx, false)
}

func (me *crudService[T]) list(ctx context.Context, withValues bool) (api.Stubs, error) {
	if err := me.init(); err != nil {
		return nil, err
	}

	if me.index.Complete {
		if withValues {
			for _, stub := range me.index.Stubs {
				if stub.Value == nil {
					stub.Value = settings.NewSettings[T](me)
					if cache, err := me.loadConfig(stub.ID, stub.Value.(T)); err != nil {
						return nil, err
					} else if !cache {
						stub.Value = nil
					}
				}
			}
			return me.index.Stubs.ToStubs(), nil
		}
		return me.index.Stubs.ToStubs(), nil
	}

	if mode == ModeOffline {
		return me.index.Stubs.ToStubs(), nil
	}

	var err error
	var stubs api.Stubs
	if stubs, err = me.service.List(ctx); err != nil {
		return nil, err
	}
	for _, stub := range stubs {
		if stub.Value != nil {
			if typeValue, ok := stub.Value.(T); ok {
				if err = me.notifyGet(stub.ID, stub.Name, typeValue); err != nil {
					return nil, err
				}
			}
		} else {
			me.index.Add(stub.ID, stub.Name)
			if err := me.tarFolder.Save(*stub, nil); err != nil {
				return nil, err
			}
		}
	}
	me.index.Complete = true
	return stubs.ToStubs(), nil
}

func (me *crudService[T]) Get(ctx context.Context, id string, v T) error {
	me.mu.Lock()
	defer me.mu.Unlock()

	if err := me.init(); err != nil {
		return err
	}

	var cache bool
	var err error
	cache, err = me.loadConfig(id, v)
	if err != nil {
		return err
	}
	if cache {
		if legacyIDAware, ok := me.service.(settings.LegacyIDAware); ok {
			settings.SetLegacyID(id, legacyIDAware.LegacyID(), v)
		}
		return nil
	}
	if mode == ModeOffline {
		return rest.Error{
			Code:    404,
			Message: fmt.Sprintf("Setting with id '%s' not found (offline mode)", id),
		}
	}
	if err = me.service.Get(ctx, id, v); err != nil {
		return err
	}
	return me.notifyGet(id, "", v)
}

func (me *crudService[T]) Update(ctx context.Context, id string, v T) error {
	me.mu.Lock()
	defer me.mu.Unlock()

	if err := me.init(); err != nil {
		return err
	}

	if mode == ModeOffline {
		return errors.New("modifications not allowed in offline mode")
	}
	if err := me.service.Update(ctx, id, v); err != nil {
		return err
	}
	return me.tarFolder.Delete(id)
}

func (me *crudService[T]) Validate(v T) error {
	me.mu.Lock()
	defer me.mu.Unlock()

	if mode == ModeOffline {
		// Validation by default succeeds in offline mode
		return nil
	}
	if validator, ok := me.service.(settings.Validator[T]); ok {
		return validator.Validate(v)
	}
	return nil
}

func (me *crudService[T]) storeConfig(id string, name string, v T) error {
	if err := me.init(); err != nil {
		return err
	}

	var err error
	var data []byte

	if data, err = settings.ToJSON(v); err != nil {
		return err
	}

	if len(name) == 0 {
		name = settings.Name(v, id)
	}

	if data, err = json.MarshalIndent(record{ID: id, Name: name, Value: data}, "", "  "); err != nil {
		return err
	}
	me.index.Add(id, name)
	return me.tarFolder.Save(api.Stub{ID: id, Name: name}, data)
}

func (me *crudService[T]) notifyGet(id string, name string, v T) error {
	if legacyIDAware, ok := me.service.(settings.LegacyIDAware); ok {
		settings.SetLegacyID(id, legacyIDAware.LegacyID(), v)
	}
	return me.storeConfig(id, name, v)
}

func (me *crudService[T]) loadConfig(id string, v T) (bool, error) {
	stub, data, err := me.tarFolder.Get(id)
	if err != nil {
		return false, err
	}
	if stub == nil {
		return false, nil
	}

	if len(data) == 0 {
		return false, nil
	}
	var record record
	if err = json.Unmarshal(data, &record); err != nil {
		return false, err
	}
	if err = settings.FromJSON(record.Value, v); err != nil {
		return false, err
	}
	if legacyIDAware, ok := me.service.(settings.LegacyIDAware); ok {
		settings.SetLegacyID(id, legacyIDAware.LegacyID(), v)
	}
	return true, nil
}

func (me *crudService[T]) SchemaID() string {
	return me.service.SchemaID() + ":cache"
}

func (me *crudService[T]) Name() string {
	return me.service.SchemaID() + ":cache"
}

var mu sync.Mutex

func CRUD[T settings.Settings](service settings.CRUDService[T], force ...bool) settings.CRUDService[T] {
	mu.Lock()
	defer mu.Unlock()

	// when running on local HTTP Cache keeping additional service caches can be detrimental
	if len(os.Getenv("DYNATRACE_MIGRATION_CACHE_FOLDER")) > 0 {
		return service
	}
	if len(force) == 0 {
		if mode == ModeDisabled {
			return service
		}
	}
	schemaID := service.SchemaID()
	if stored, ok := caches[schemaID]; ok {
		return stored.(*crudService[T])
	}
	if ncs, ok := service.(settings.NoCacheService); ok && ncs.NoCache() {
		return service
	}

	cs := &crudService[T]{
		service: service,
		folder:  path.Join(cache_root_folder, strings.ReplaceAll(service.SchemaID(), ":", ".")),
	}
	caches[schemaID] = cs
	return cs
}
