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
	GmtMinus0100 TimezoneEnum
	GmtMinus0200 TimezoneEnum
	GmtMinus0300 TimezoneEnum
	GmtMinus0400 TimezoneEnum
	GmtMinus0500 TimezoneEnum
	GmtMinus0600 TimezoneEnum
	GmtMinus0700 TimezoneEnum
	GmtMinus0800 TimezoneEnum
	GmtMinus0900 TimezoneEnum
	GmtMinus1000 TimezoneEnum
	GmtMinus1100 TimezoneEnum
	GmtMinus1200 TimezoneEnum
	GmtPlus0000  TimezoneEnum
	GmtPlus0100  TimezoneEnum
	GmtPlus0200  TimezoneEnum
	GmtPlus0300  TimezoneEnum
	GmtPlus0400  TimezoneEnum
	GmtPlus0500  TimezoneEnum
	GmtPlus0600  TimezoneEnum
	GmtPlus0700  TimezoneEnum
	GmtPlus0800  TimezoneEnum
	GmtPlus0900  TimezoneEnum
	GmtPlus1000  TimezoneEnum
	GmtPlus1100  TimezoneEnum
	GmtPlus1200  TimezoneEnum
}{
	"GMT-01:00",
	"GMT-02:00",
	"GMT-03:00",
	"GMT-04:00",
	"GMT-05:00",
	"GMT-06:00",
	"GMT-07:00",
	"GMT-08:00",
	"GMT-09:00",
	"GMT-10:00",
	"GMT-11:00",
	"GMT-12:00",
	"GMT+00:00",
	"GMT+01:00",
	"GMT+02:00",
	"GMT+03:00",
	"GMT+04:00",
	"GMT+05:00",
	"GMT+06:00",
	"GMT+07:00",
	"GMT+08:00",
	"GMT+09:00",
	"GMT+10:00",
	"GMT+11:00",
	"GMT+12:00",
}
