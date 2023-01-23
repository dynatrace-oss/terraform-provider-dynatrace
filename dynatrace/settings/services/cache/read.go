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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

type readService[T settings.Settings] struct {
	service settings.RService[T]
	folder  string
	index   *stubIndex
}

func (me *readService[T]) List() (settings.Stubs, error) {
	return me.list(true)
}

func (me *readService[T]) ListNoValues() (settings.Stubs, error) {
	return me.list(false)
}

func (me *readService[T]) list(withValues bool) (settings.Stubs, error) {
	var err error
	var index *stubIndex
	if exists(indexFile(me.folder)) {
		if index, err = me.loadIndex(); err != nil {
			return nil, err
		}
		if withValues && index.Complete {
			for _, stub := range index.Stubs {
				stub.Value = settings.NewSettings[T](me)
				if cache, err := me.loadConfig(stub.ID, stub.Value.(T)); err != nil {
					return nil, err
				} else if !cache {
					stub.Value = nil
				}
			}
		}
		return index.Stubs.ToStubs(), nil
	}
	if mode == ModeOffline {
		return settings.Stubs{}, nil
	}

	var stubs settings.Stubs
	if stubs, err = me.service.List(); err != nil {
		return nil, err
	}
	if me.index != nil {
		for _, stub := range stubs {
			me.index.Add(stub.ID, stub.Name)
		}
		me.index.Complete = true
		me.storeIndex(me.index)
	} else {
		entries, _ := os.ReadDir(me.folder)
		for _, entry := range entries {
			if entry.Name() == ".index.json" {
				continue
			}
			var data []byte
			if data, err = os.ReadFile(path.Join(me.folder, entry.Name())); err != nil {
				return nil, err
			}
			var record record
			if err = json.Unmarshal(data, &record); err != nil {
				return nil, err
			}
			var sttngs T
			if record.Value != nil {
				id := record.ID
				sttngs = settings.NewSettings(me.service)
				if err = settings.FromJSON(record.Value, sttngs); err != nil {
					return nil, err
				}
				if legacyIDAware, ok := me.service.(settings.LegacyIDAware); ok {
					settings.SetLegacyID(id, legacyIDAware.LegacyID(), sttngs)
				}
			}
			stubs = append(stubs, &settings.Stub{ID: record.ID, Name: record.Name, Value: sttngs})
		}
		me.storeIndex(&stubIndex{Complete: true, Stubs: stubs})
	}
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

func (me *readService[T]) Get(id string, v T) error {
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
		if err = me.service.Get(id, v); err != nil {
			return err
		}
		return me.notifyGet(id, v)
	}
	return nil
}

func (me *readService[T]) loadIndex() (*stubIndex, error) {
	if me.index != nil {
		return me.index, nil
	}
	me.index = new(stubIndex)
	var err error
	var data []byte

	filePath := indexFile(me.folder)
	if exists(filePath) {
		if data, err = os.ReadFile(filePath); err != nil {
			return nil, err
		}
		if err = json.Unmarshal(data, me.index); err != nil {
			return nil, err
		}
	}
	me.index.IDs = map[string]*settings.Stub{}
	for _, stub := range me.index.Stubs {
		me.index.IDs[stub.ID] = stub
	}
	return me.index, nil
}

func (me *readService[T]) storeIndex(index *stubIndex) error {
	os.MkdirAll(me.folder, os.ModePerm)
	var err error
	var file *os.File
	var data []byte

	if file, err = os.Create(indexFile(me.folder)); err != nil {
		return err
	}
	defer file.Close()
	if data, err = json.Marshal(index); err != nil {
		return err
	}
	if _, err = file.Write(data); err != nil {
		return err
	}

	// keeping index that has just been stored in memory
	me.index = index
	me.index.IDs = map[string]*settings.Stub{}
	// We don't want to keep the settings that are potentially
	// attached to the in memory stubs.
	// We read these settings from disk if required
	me.index.Stubs = me.index.Stubs.ToStubs()
	for _, stub := range me.index.Stubs {
		stub.Value = nil
		me.index.IDs[stub.ID] = stub
	}
	return nil
}

func (me *readService[T]) dataFile(id string) string {
	filename := fmt.Sprintf("%s.bin.json", id)
	filename = strings.ReplaceAll(filename, ":", ".")
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, ":", "_")
	filename = strings.ReplaceAll(filename, "|", "_")
	filename = strings.ReplaceAll(filename, "<", "_")
	filename = strings.ReplaceAll(filename, ">", "_")
	filename = strings.ReplaceAll(filename, "\"", "_")
	filename = strings.ReplaceAll(filename, "?", "_")
	filename = strings.ReplaceAll(filename, "*", "_")
	return path.Join(me.folder, filename)
}

func (me *readService[T]) storeConfig(id string, v T) error {
	os.MkdirAll(me.folder, os.ModePerm)
	var err error
	var data []byte
	var file *os.File

	if file, err = os.Create(me.dataFile(id)); err != nil {
		return err
	}
	defer file.Close()

	if data, err = settings.ToJSON(v); err != nil {
		return err
	}

	configName := settings.Name(v)
	if data, err = json.MarshalIndent(record{ID: id, Name: configName, Value: data}, "", "  "); err != nil {
		return err
	}
	if _, err = file.Write(data); err != nil {
		return err
	}

	if me.index != nil {
		if _, found := me.index.IDs[id]; !found {
			log.Printf("%s not found", id)
			var index *stubIndex
			if index, err = me.loadIndex(); err != nil {
				return err
			}
			return me.storeIndex(index.Add(id, configName))
		}
	}
	return nil
}

func (me *readService[T]) notifyGet(id string, v T) error {
	if legacyIDAware, ok := me.service.(settings.LegacyIDAware); ok {
		settings.SetLegacyID(id, legacyIDAware.LegacyID(), v)
	}
	return me.storeConfig(id, v)
}

func (me *readService[T]) loadConfig(id string, v T) (bool, error) {
	var err error
	var data []byte
	filePath := me.dataFile(id)
	if exists(filePath) {
		if data, err = os.ReadFile(filePath); err != nil {
			return false, err
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
	return false, nil
}

func (me *readService[T]) SchemaID() string {
	return me.service.SchemaID() + ":cache"
}

func Read[T settings.Settings](service settings.RService[T], force ...bool) settings.RService[T] {
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
