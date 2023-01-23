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

package export

import (
	"fmt"
	"strings"
)

type ReplaceFunc func(s string, cnt int) string

func DefaultReplace(s string, cnt int) string {
	return fmt.Sprintf("%s(%d)", s, cnt)
}

func ResourceName(s string, cnt int) string {
	return fmt.Sprintf("%s_%d", s, cnt)
}

type UniqueNamer interface {
	Name(string) string
	Replace(ReplaceFunc) UniqueNamer
}

func NewUniqueNamer() UniqueNamer {
	return &nameCounter{m: map[string]int{}}
}

type nameCounter struct {
	m       map[string]int
	replace ReplaceFunc
}

func (me *nameCounter) Replace(replace ReplaceFunc) UniqueNamer {
	me.replace = replace
	return me
}

func (me *nameCounter) Name(name string) string {
	cnt, found := me.m[strings.ToLower(name)]
	if !found {
		me.m[strings.ToLower(name)] = 0
		return name
	} else {
		me.m[strings.ToLower(name)] = cnt + 1
	}
	if me.replace == nil {
		return DefaultReplace(name, me.m[strings.ToLower(name)])
	}
	return me.replace(name, me.m[strings.ToLower(name)])
}
