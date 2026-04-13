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
	"os"
	"path"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/envutils"
	"github.com/google/uuid"
)

type Mode byte

const (
	ModeDisabled Mode = iota
	ModeEnabled
	ModeOffline
)

var mode = ModeDisabled

func Enable() {
	mode = ModeEnabled
}

func Disable() {
	mode = ModeDisabled
}

func Offline() {
	mode = ModeOffline
}

func Cleanup() {
	if envutils.DTNoCacheCleanup.Get() {
		return
	}
	os.RemoveAll(cache_root_folder)
}

var cache_root_folder = getCacheRootFolder()

func GetCacheFolder() string {
	return cache_root_folder
}

func getCacheRootFolder() string {
	folder := path.Join(os.TempDir(), ".terraform-provider-dynatrace", uuid.NewString())
	if envFolder := envutils.DTCacheFolder.Get(); envFolder != "" {
		folder = envFolder
	}
	deleteCache := envutils.DTCacheDeleteOnLaunch.Get()
	if len(deleteCache) != 0 && deleteCache != "false" {
		os.RemoveAll(folder)
	}
	if envutils.CacheOfflineMode.Get() {
		Offline()
	}
	return folder
}

type stubIndex struct {
	Complete bool                 `json:"complete"`
	Stubs    api.Stubs            `json:"stubs"`
	IDs      map[string]*api.Stub `json:"-"`
}

func (me *stubIndex) Remove(id string) *stubIndex {
	var result api.Stubs
	for _, stub := range me.Stubs {
		if stub.ID != id {
			result = append(result, stub)
		}
	}
	me.Stubs = result
	delete(me.IDs, id)
	return me
}

func (me *stubIndex) Add(id string, name string) *stubIndex {
	for _, stub := range me.Stubs {
		if stub.ID == id {
			return me
		}
	}
	stub := &api.Stub{ID: id, Name: name}
	me.Stubs = append(me.Stubs, stub)
	if me.IDs == nil {
		me.IDs = map[string]*api.Stub{}
	}
	me.IDs[id] = stub
	return me
}

var caches = map[string]any{}

type record struct {
	ID    string
	Name  string
	Value json.RawMessage
}
