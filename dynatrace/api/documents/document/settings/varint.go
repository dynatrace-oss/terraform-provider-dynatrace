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

package documents

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type VarInt string // pattern: ^{{.+}}$

func (me VarInt) MarshalJSON() ([]byte, error) {
	if i, err := strconv.Atoi(string(me)); err == nil {
		return json.Marshal(i)
	}
	return json.Marshal(string(me))
}

func (me *VarInt) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*me = VarInt(s)
	}
	var i int
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	*me = VarInt(fmt.Sprintf("%d", i))
	return nil
}
