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

package export

import (
	"os"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
)

var restLogFile *os.File

func ConfigureRESTLog() (err error) {
	restLogFileName := os.Getenv("DT_REST_DEBUG_LOG")
	if len(restLogFileName) > 0 {
		if restLogFile, err = os.Create(restLogFileName); err != nil {
			return err
		}
		CleanUp.Register(func() { restLogFile.Close() })
		rest.SetLogWriter(restLogFile)
	}
	return nil
}
