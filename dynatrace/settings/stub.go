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

type Stub struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	EntityID string  `json:"entityId,omitempty"`
	Value    any     `json:"-"`
	LegacyID *string `json:"legacyID,omitempty"`
}

type Stubs []*Stub

func (me *Stubs) ToStubs() Stubs {
	res := []*Stub{}
	for _, stub := range *me {
		if len(stub.ID) == 0 && len(stub.EntityID) != 0 {
			stub.ID = stub.EntityID
		}
		res = append(res, &Stub{ID: stub.ID, Name: stub.Name, Value: stub.Value, EntityID: stub.EntityID, LegacyID: stub.LegacyID})
	}
	return res
}

type StubList struct {
	Values []*Stub `json:"values"`
}

func (me *StubList) ToStubs() Stubs {
	for _, stub := range me.Values {
		if len(stub.ID) == 0 && len(stub.EntityID) != 0 {
			stub.ID = stub.EntityID
		}
	}
	return me.Values
}

type RecordStubs interface {
	ToStubs() Stubs
}
