/**
* @license
* Copyright 2024 Dynatrace LLC
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

package multiuse

import (
	"fmt"
	"strings"
)

var ID_PARENT_DELIMITER = "~@|#:;"

func ExtractIDParent(id string) (bool, string, string) {
	delimiterIndex := strings.Index(id, ID_PARENT_DELIMITER)
	if delimiterIndex >= 0 {
		applicationID := id[delimiterIndex+len(ID_PARENT_DELIMITER):]
		realId := id[:delimiterIndex]
		return true, applicationID, realId
	}
	return false, "", ""
}

func EncodeIDParent(id string, parentID *string) string {
	getID := id
	if parentID != nil {
		getID = fmt.Sprintf("%s%s%s", id, ID_PARENT_DELIMITER, *parentID)
	}

	return getID
}
