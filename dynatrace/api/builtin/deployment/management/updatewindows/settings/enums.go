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

package updatewindows

type RecurrenceEnum string

var RecurrenceEnums = struct {
	Daily   RecurrenceEnum
	Monthly RecurrenceEnum
	Once    RecurrenceEnum
	Weekly  RecurrenceEnum
}{
	"DAILY",
	"MONTHLY",
	"ONCE",
	"WEEKLY",
}

type TimezoneEnum string

var TimezoneEnums = struct {
	Gmt0000 TimezoneEnum
	Gmt0100 TimezoneEnum
	Gmt0200 TimezoneEnum
	Gmt0300 TimezoneEnum
	Gmt0400 TimezoneEnum
	Gmt0500 TimezoneEnum
	Gmt0600 TimezoneEnum
	Gmt0700 TimezoneEnum
	Gmt0800 TimezoneEnum
	Gmt0900 TimezoneEnum
	Gmt1000 TimezoneEnum
	Gmt1100 TimezoneEnum
	Gmt1200 TimezoneEnum
}{
	"GMT_00_00",
	"GMT_01_00",
	"GMT_02_00",
	"GMT_03_00",
	"GMT_04_00",
	"GMT_05_00",
	"GMT_06_00",
	"GMT_07_00",
	"GMT_08_00",
	"GMT_09_00",
	"GMT_10_00",
	"GMT_11_00",
	"GMT_12_00",
}
