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

package osarch

// Value The value to compare to.
type Value string

func (v Value) Ref() *Value {
	return &v
}

func (v *Value) String() string {
	return string(*v)
}

// Values offers the known enum values
var Values = struct {
	Arm    Value
	Ia64   Value
	Parisc Value
	Ppc    Value
	Ppcle  Value
	S390   Value
	Sparc  Value
	X86    Value
	Zos    Value
}{
	"ARM",
	"IA64",
	"PARISC",
	"PPC",
	"PPCLE",
	"S390",
	"SPARC",
	"X86",
	"ZOS",
}
