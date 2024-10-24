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

package rest

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type RESTLogger struct {
	log *log.Logger
}

func (l *RESTLogger) Print(ctx context.Context, v ...any) {
	if stdoutLog {
		tflog.Debug(ctx, fmt.Sprint(append(append([]any{}, "[HTTP]"), v...)...))
	}
	l.log.Print(v...)
}

func (l *RESTLogger) Printf(ctx context.Context, format string, v ...any) {
	if stdoutLog {
		tflog.Debug(ctx, fmt.Sprintf("[HTTP] "+format, v...))
	}
	l.log.Printf(format, v...)
}

func (l *RESTLogger) Println(ctx context.Context, v ...any) {
	if stdoutLog {
		tflog.Debug(ctx, fmt.Sprint(append(append([]any{}, "[HTTP]"), v...)...))
	}
	l.log.Println(v...)
}
