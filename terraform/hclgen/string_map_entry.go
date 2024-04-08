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

package hclgen

import (
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

type stringMapEntry struct {
	Indent      string
	Key         string
	BreadCrumbs string
	Optional    bool
	Values      map[string]string
}

func (me *stringMapEntry) IsOptional() bool {
	return false
}

func (me *stringMapEntry) IsComputed() bool {
	return false
}

func (me *stringMapEntry) IsDefault() bool {
	return false
}

func (me *stringMapEntry) IsLessThan(other exportEntry) bool {
	switch ro := other.(type) {
	case *primitiveEntry:
		return false
	case *resourceEntry:
		return strings.Compare(me.Key, ro.Key) < 0
	case *mapEntry:
		return strings.Compare(me.Key, ro.Key) < 0
	case *stringMapEntry:
		return strings.Compare(me.Key, ro.Key) < 0
	}
	return false
}

func (me *stringMapEntry) Write(w *hclwrite.Body, indent string) error {
	objTokens := []hclwrite.ObjectAttrTokens{}
	for key, value := range me.Values {
		objTokens = append(objTokens, hclwrite.ObjectAttrTokens{Name: hclwrite.TokensForValue(ctyVal(key, "")), Value: hclwrite.TokensForValue(ctyVal(value, ""))})
	}

	tokens := hclwrite.TokensForObject(objTokens)
	w.SetAttributeRaw(me.Key, tokens)
	return nil
}
