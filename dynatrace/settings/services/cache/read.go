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
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache/tar"
)

type readService[T settings.Settings] struct {
	mu        sync.Mutex
	service   settings.RService[T]
	folder    string
	index     *stubIndex
	tarFolder *tar.Folder
}

func (me *readService[T]) init() error {
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

func (me *readService[T]) List(ctx context.Context) (api.Stubs, error) {
	me.mu.Lock()
	defer me.mu.Unlock()

	return me.list(ctx, true)
}

func (me *readService[T]) ListNoValues(ctx context.Context) (api.Stubs, error) {
	me.mu.Lock()
	defer me.mu.Unlock()

	return me.list(ctx, false)
}

func (me *readService[T]) list(ctx context.Context, withValues bool) (api.Stubs, error) {
	if err := me.init(); err != nil {
		return nil, err
	}

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
	}

	if mode == ModeOffline || me.index.Complete {
		return me.index.Stubs.ToStubs(), nil
	}

	var err error
	var stubs api.Stubs
	if stubs, err = me.service.List(ctx); err != nil {
		return nil, err
	}
	me.index.Complete = true
	for _, stub := range stubs {
		if stub.Value != nil {
			if typeValue, ok := stub.Value.(T); ok {
				if err = me.notifyGet(stub.ID, typeValue); err != nil {
					return nil, err
				}
			}
		}
	}
	return stubs.ToStubs(), nil
}

func (me *readService[T]) Get(ctx context.Context, id string, v T) error {
	me.mu.Lock()
	defer me.mu.Unlock()

	if err := me.init(); err != nil {
		return err
	}

	var cache bool
	var err error
	if cache, err = me.loadConfig(id, v); err != nil {
		return err
	} else if cache {
		if os.Getenv("DT_REST_DEBUG_REQUESTS") == "cache" {
			log.Println("cache", me.SchemaID(), id)
		}
		if legacyIDAware, ok := me.service.(settings.LegacyIDAware); ok {
			settings.SetLegacyID(id, legacyIDAware.LegacyID(), v)
		}
	} else if mode == ModeOffline {
		return rest.Error{
			Code:    404,
			Message: fmt.Sprintf("Setting with id '%s' not found (offline mode)", id),
		}
	} else {
		if err = me.service.Get(ctx, id, v); err != nil {
			return err
		}
		return me.notifyGet(id, v)
	}
	return nil
}

func (me *readService[T]) notifyGet(id string, v T) error {
	if legacyIDAware, ok := me.service.(settings.LegacyIDAware); ok {
		settings.SetLegacyID(id, legacyIDAware.LegacyID(), v)
	}
	return me.storeConfig(id, v)
}

func (me *readService[T]) storeConfig(id string, v T) error {
	if err := me.init(); err != nil {
		return err
	}

	var err error
	var data []byte

	if data, err = settings.ToJSON(v); err != nil {
		return err
	}

	configName := settings.Name(v, id)
	if data, err = json.MarshalIndent(record{ID: id, Name: configName, Value: data}, "", "  "); err != nil {
		return err
	}
	me.index.Add(id, configName)
	return me.tarFolder.Save(api.Stub{ID: id, Name: configName}, data)
}

func (me *readService[T]) loadConfig(id string, v T) (bool, error) {
	stub, data, err := me.tarFolder.Get(id)
	if err != nil {
		return false, err
	}
	if stub == nil {
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

func (me *readService[T]) SchemaID() string {
	return me.service.SchemaID() + ":cache"
}

func Read[T settings.Settings](service settings.RService[T], force ...bool) settings.RService[T] {
	mu.Lock()
	defer mu.Unlock()

	if len(force) == 0 {
		if mode == ModeDisabled {
			return service
		}
	}
	schemaID := service.SchemaID()
	if stored, ok := caches[schemaID]; ok {
		return stored.(*readService[T])
	}
	if ncs, ok := service.(settings.NoCacheService); ok && ncs.NoCache() {
		return service
	}

	cs := &readService[T]{
		service: service,
		folder:  path.Join(cache_root_folder, strings.ReplaceAll(service.SchemaID(), ":", ".")),
	}
	caches[schemaID] = cs
	return cs
}
