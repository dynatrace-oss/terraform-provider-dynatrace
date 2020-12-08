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


package logging

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Enable redirects logging into a an output file
func Enable(fn func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		if file, e := os.OpenFile("terraform-provider-dynatrace.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); e != nil {
			panic(errors.New("Failed to open log file"))
		} else {
			log.SetOutput(file)
		}
		return fn(ctx, d, m)
	}
}
