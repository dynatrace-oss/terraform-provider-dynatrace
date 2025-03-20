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
	"io"
	"log"
	"os"

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

var stdoutLog = os.Getenv("DYNATRACE_LOG_HTTP") == "stdout"
var logger = initLogger()
var Logger = logger

type onDemandWriter struct {
	logFileName string
	file        *os.File
}

func (odw *onDemandWriter) Write(p []byte) (n int, err error) {
	if odw.file == nil {
		if odw.file, err = os.OpenFile(odw.logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm); err != nil {
			return 0, err
		}
	}
	return odw.file.Write(p)
}

func initLogger() *RESTLogger {
	restLogFileName := os.Getenv("DYNATRACE_LOG_HTTP")
	if len(restLogFileName) > 0 && restLogFileName != "false" && !stdoutLog {
		logger := log.New(os.Stderr, "", log.LstdFlags)
		if restLogFileName != "true" {
			logger.SetOutput(&onDemandWriter{logFileName: restLogFileName})
		}
		return &RESTLogger{log: logger}
	}
	return &RESTLogger{log: log.New(io.Discard, "", log.LstdFlags)}
}

func SetLogWriter(writer io.Writer) error {
	logger.log.SetOutput(writer)
	return nil
}
