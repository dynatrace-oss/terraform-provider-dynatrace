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


package main

import (
	"io"

	"github.com/dtcookie/dynatrace/api/config/managementzones"
	"github.com/dtcookie/dynatrace/terraform"
)

// ManagementZoneExporter has no documentation
type ManagementZoneExporter struct {
	ManagementZone *managementzones.ManagementZone
}

// ToHCL has no documentation
func (exporter *ManagementZoneExporter) ToHCL(w io.Writer) error {
	bytes, err := terraform.Marshal(exporter.ManagementZone, "dynatrace_management_zone", exporter.ManagementZone.Name)
	if err != nil {
		return err
	}
	_, err = w.Write(bytes)
	return err
}

// ToJSON has no documentation
func (exporter *ManagementZoneExporter) ToJSON(w io.Writer) error {
	bytes, err := terraform.MarshalJSON(exporter.ManagementZone, "dynatrace_management_zone", exporter.ManagementZone.Name)
	if err != nil {
		return err
	}
	_, err = w.Write(bytes)
	return err
}
