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

package customlogsourcesettings

type ContextType string

var ContextTypes = struct {
	DtEntityProcessGroup ContextType
}{
	"dt.entity.process_group",
}

type LogSourceType string

var LogSourceTypes = struct {
	LogPathPattern  LogSourceType
	WindowsEventLog LogSourceType
}{
	"LOG_PATH_PATTERN",
	"WINDOWS_EVENT_LOG",
}

type WildcardExpansionType string

var WildcardExpansionTypes = struct {
	Attribute WildcardExpansionType
}{
	"attribute",
}
