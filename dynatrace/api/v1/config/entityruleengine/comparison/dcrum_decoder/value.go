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

package dcrum_decoder

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
	AllOther         Value
	CitrixAppFlow    Value
	CitrixIca        Value
	CitrixIcaOverSSL Value
	DB2Drda          Value
	HTTP             Value
	HTTPS            Value
	HTTPExpress      Value
	Informix         Value
	MySQL            Value
	Oracle           Value
	SAPGUI           Value
	SAPGUIOverHTTP   Value
	SAPGUIOverHTTPS  Value
	SAPHanaDB        Value
	SAPRfc           Value
	SSL              Value
	TDS              Value
}{
	"ALL_OTHER",
	"CITRIX_APPFLOW",
	"CITRIX_ICA",
	"CITRIX_ICA_OVER_SSL",
	"DB2_DRDA",
	"HTTP",
	"HTTPS",
	"HTTP_EXPRESS",
	"INFORMIX",
	"MYSQL",
	"ORACLE",
	"SAP_GUI",
	"SAP_GUI_OVER_HTTP",
	"SAP_GUI_OVER_HTTPS",
	"SAP_HANA_DB",
	"SAP_RFC",
	"SSL",
	"TDS",
}
