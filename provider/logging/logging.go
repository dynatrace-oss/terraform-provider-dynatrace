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
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var Logger = hclog.New(&hclog.LoggerOptions{
	IncludeLocation: false,
	DisableTime:     true,
})

type onDemandLogger struct {
	mu   sync.Mutex // ensures atomic writes; protects the following fields
	file *os.File
	name string
}

type void struct{}

func (v *void) Write(p []byte) (n int, err error) { return 0, nil }

func createDebugLogger(name string) *log.Logger {

	prefix := "terraform-provider-dynatrace"

	logDebugPrefix := os.Getenv("DYNATRACE_LOG_DEBUG_PREFIX")
	if len(logDebugPrefix) > 0 && logDebugPrefix != "false" {
		prefix = logDebugPrefix
	}

	name = fmt.Sprintf("%s.%s", prefix, name)

	if os.Getenv("DYNATRACE_DEBUG") != "true" {
		return log.New(&void{}, "", log.LstdFlags)
	}
	if len(name) == 0 {
		return log.New(&void{}, "", log.LstdFlags)
	}
	return log.New(&onDemandLogger{name: name}, "", log.LstdFlags)
}

var Debug = struct {
	Info *log.Logger
	Warn *log.Logger
}{
	Info: createDebugLogger("export.log"),
	Warn: createDebugLogger("warnings.log"),
}

var odl = &onDemandLogger{name: "terraform-provider-dynatrace.log"}
var File = log.New(odl, "", log.LstdFlags)

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
func Enable(fn func(context.Context, *schema.ResourceData, any) diag.Diagnostics) func(context.Context, *schema.ResourceData, any) diag.Diagnostics {
	if os.Getenv("DYNATRACE_DEBUG") != "true" {
		return fn
	}
	return func(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
		log.SetOutput(odl)
		return fn(ctx, d, m)
	}
}

func EnableDS(fn func(d *schema.ResourceData, m any) error) func(*schema.ResourceData, any) error {
	if os.Getenv("DYNATRACE_DEBUG") != "true" {
		return fn
	}
	return func(d *schema.ResourceData, m any) error {
		log.SetOutput(odl)
		return fn(d, m)
	}
}

func EnableDSCtx(fn func(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics) func(context.Context, *schema.ResourceData, any) diag.Diagnostics {
	if os.Getenv("DYNATRACE_DEBUG") != "true" {
		return fn
	}
	return func(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
		log.SetOutput(odl)
		return fn(ctx, d, m)
	}
}

func EnableSchemaSetFunc(fn func(v any) int) func(v any) int {
	if os.Getenv("DYNATRACE_DEBUG") != "true" {
		return fn
	}
	return func(v any) int {
		log.SetOutput(odl)
		return fn(v)
	}
}

// Enable redirects logging into a an output file
func SetOutput() {
	if os.Getenv("DYNATRACE_DEBUG") != "true" {
		return
	}
	log.SetOutput(odl)
}

// EnableSchemaDiff redirects logging into a an output file
func EnableSchemaDiff(fn func(k, old, new string, d *schema.ResourceData) bool) func(k, old, new string, d *schema.ResourceData) bool {
	if os.Getenv("DYNATRACE_DEBUG") != "true" {
		return fn
	}
	return func(k, old, new string, d *schema.ResourceData) bool {
		log.SetOutput(odl)
		return fn(k, old, new, d)
	}
}

// EnableCustomizeDiff redirects logging into a an output file
func EnableCustomizeDiff(fn func(context.Context, *schema.ResourceDiff, any) error) func(context.Context, *schema.ResourceDiff, any) error {
	if os.Getenv("DYNATRACE_DEBUG") != "true" {
		return fn
	}
	return func(ctx context.Context, d *schema.ResourceDiff, meta any) error {
		log.SetOutput(odl)
		return fn(ctx, d, meta)
	}
}
