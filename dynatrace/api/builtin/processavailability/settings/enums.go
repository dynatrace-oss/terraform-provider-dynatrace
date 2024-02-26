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

package processavailability

type OperatingSystem string

var OperatingSystems = struct {
	Aix     OperatingSystem
	Linux   OperatingSystem
	Windows OperatingSystem
}{
	"AIX",
	"LINUX",
	"WINDOWS",
}

type ProcessItem string

var ProcessItems = struct {
	Commandline    ProcessItem
	Executable     ProcessItem
	Executablepath ProcessItem
	User           ProcessItem
}{
	"commandLine",
	"executable",
	"executablePath",
	"user",
}

type RuleType string

var RuleTypes = struct {
	Ruletypehost    RuleType
	Ruletypeprocess RuleType
}{
	"RuleTypeHost",
	"RuleTypeProcess",
}
