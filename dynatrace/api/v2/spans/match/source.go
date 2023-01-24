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

package match

type Source string

var Sources = struct {
	SpanName                      Source
	SpanKind                      Source
	Attribute                     Source
	InstrumentationLibraryName    Source
	InstrumentationLibraryVersion Source
	InstrumentationScopeName      Source
	InstrumentationScopeVersion   Source
}{
	SpanName:                      Source("SPAN_NAME"),
	SpanKind:                      Source("SPAN_KIND"),
	Attribute:                     Source("ATTRIBUTE"),
	InstrumentationLibraryName:    Source("INSTRUMENTATION_LIBRARY_NAME"),
	InstrumentationLibraryVersion: Source("INSTRUMENTATION_LIBRARY_VERSION"),
	InstrumentationScopeName:      Source("INSTRUMENTATION_SCOPE_NAME"),
	InstrumentationScopeVersion:   Source("INSTRUMENTATION_SCOPE_VERSION"),
}
