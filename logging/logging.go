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
	"log"
	"os"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type onDemandLogger struct {
	mu   sync.Mutex // ensures atomic writes; protects the following fields
	file *os.File
	name string
}

var odl = &onDemandLogger{name: "terraform-provider-dynatrace.log"}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (odl *onDemandLogger) f() (*os.File, error) {
	odl.mu.Lock()
	defer odl.mu.Unlock()
	if odl.file != nil {
		return odl.file, nil
	}
	var err error
	if odl.file, err = os.OpenFile(odl.name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err != nil {
		return nil, err
	}
	return odl.file, nil
}

func (odl *onDemandLogger) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	s := string(p)
	if strings.Contains(s, "[WARN] Truncating attribute path of") {
		return 0, nil
	}
	file, err := odl.f()
	if err != nil {
		return 0, err
	}
	return file.Write(p)
}

// Enable redirects logging into a an output file
func Enable(fn func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		log.SetOutput(odl)
		return fn(ctx, d, m)
	}
}

// Enable redirects logging into a an output file
func SetOutput() {
	log.SetOutput(odl)
}

// EnableSchemaDiff redirects logging into a an output file
func EnableSchemaDiff(fn func(k, old, new string, d *schema.ResourceData) bool) func(k, old, new string, d *schema.ResourceData) bool {
	return func(k, old, new string, d *schema.ResourceData) bool {
		log.SetOutput(odl)
		return fn(k, old, new, d)
	}
}

// EnableCustomizeDiff redirects logging into a an output file
func EnableCustomizeDiff(fn func(context.Context, *schema.ResourceDiff, interface{}) error) func(context.Context, *schema.ResourceDiff, interface{}) error {
	return func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
		log.SetOutput(odl)
		return fn(ctx, d, meta)
	}
}
