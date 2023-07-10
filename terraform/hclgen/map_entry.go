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
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

type mapEntry struct {
	Indent      string
	Key         string
	BreadCrumbs string
	Optional    bool
	Entries     exportEntries
}

func (me *mapEntry) IsOptional() bool {
	return false
}

func (me *mapEntry) IsComputed() bool {
	return false
}

func (me *mapEntry) IsDefault() bool {
	return false
}

func (me *mapEntry) IsLessThan(other exportEntry) bool {
	switch ro := other.(type) {
	case *primitiveEntry:
		return false
	case *resourceEntry:
		return strings.Compare(me.Key, ro.Key) < 0
	case *mapEntry:
		return strings.Compare(me.Key, ro.Key) < 0
	}
	return false
}

func (me *mapEntry) Write(w *hclwrite.Body, indent string) error {
	sort.SliceStable(me.Entries, me.Entries.Less)
	objTokens := []hclwrite.ObjectAttrTokens{}
	for _, entry := range me.Entries {
		primEntry := entry.(*primitiveEntry)
		objTokens = append(objTokens, hclwrite.ObjectAttrTokens{Name: hclwrite.TokensForIdentifier(primEntry.Key), Value: hclwrite.TokensForValue(ctyVal(primEntry.Value, ""))})
	}

	tokens := hclwrite.TokensForObject(objTokens)
	w.SetAttributeRaw(me.Key, tokens)
	// block := w.AppendNewBlock(me.Key, nil)
	// body := block.Body()

	// sort.SliceStable(me.Entries, me.Entries.Less)
	// for _, entry := range me.Entries {
	// 	if !entry.IsComputed() {
	// 		if !(entry.IsOptional() && entry.IsDefault()) {
	// 			if err := entry.Write(body, indent+"  "); err != nil {
	// 				return err
	// 			}
	// 		} else {
	// 			body.AppendUnstructuredTokens(hclwrite.Tokens{
	// 				&hclwrite.Token{Type: hclsyntax.TokenComment, Bytes: []byte("#")},
	// 			})
	// 			if err := entry.Write(body, indent+"  "); err != nil {
	// 				return err
	// 			}
	// 		}
	// 	}
	// }
	return nil
}
