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
	"net/url"
)

func Path(s string) func(args ...string) string {
	return path(s).Complete
}

type path string

func (p path) Complete(args ...string) string {
	if len(args) == 0 {
		return string(p)
	}
	fmtargs := make([]any, len(args))
	for idx, s := range args {
		fmtargs[idx] = url.PathEscape(s)
	}
	return fmt.Sprintf(string(p), fmtargs...)
}
