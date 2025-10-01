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

package workflows

import (
	"encoding/json"
	"fmt"
)

type StringArray []string

func (s StringArray) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string(s))
}

func (s *StringArray) UnmarshalJSON(data []byte) error {
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		*s = StringArray(arr)
		return nil
	}
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*s = StringArray{single}
		return nil
	}
	return fmt.Errorf("StringArray: cannot unmarshal %s", string(data))
}
