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
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
)

var logger = initLogger()
var Logger = logger

type onDemandWriter struct {
	logFileName string
	file        *os.File
}

func (odw *onDemandWriter) Write(p []byte) (n int, err error) {
	if odw.file == nil {
		if odw.file, err = os.OpenFile(odw.logFileName, os.O_APPEND|os.O_CREATE, os.ModePerm); err != nil {
			return 0, err
		}
	}
	return odw.file.Write(p)
}

func initLogger() *log.Logger {
	restLogFileName := os.Getenv("DYNATRACE_LOG_HTTP")
	if len(restLogFileName) > 0 && restLogFileName != "false" {
		logger := log.New(os.Stderr, "", log.LstdFlags)
		if restLogFileName != "true" {
			logger.SetOutput(&onDemandWriter{logFileName: restLogFileName})
		}
		return logger
	}
	return log.New(io.Discard, "", log.LstdFlags)
}

func SetLogWriter(writer io.Writer) error {
	logger.SetOutput(writer)
	return nil
}

var jar = createJar()

func createJar() *cookiejar.Jar {
	jar, _ := cookiejar.New(nil)
	return jar
}

type Request interface {
	Raw() ([]byte, error)
	Finish(v ...any) error
	Expect(codes ...int) Request
	Payload(any) Request
	OnResponse(func(resp *http.Response)) Request
}

type statuscodes []int

func (me statuscodes) contains(code int) bool {
	for _, c := range me {
		if c == code {
			return true
		}
	}
	return false
}

type request struct {
	client     *defaultClient
	url        string
	expect     statuscodes
	method     string
	payload    any
	headers    map[string]string
	onResponse func(resp *http.Response)
}

func (me *request) authenticate(req *http.Request) {
	req.Header.Add("Authorization", "Api-Token "+me.client.apiToken)
	req.Header.Set("User-Agent", "Dynatrace Terraform Provider")
}

func (me *request) Payload(payload any) Request {
	me.payload = payload
	return me
}

func (me *request) Finish(vs ...any) error {
	var v any
	if len(vs) > 0 {
		v = vs[0]
	}
	var err error
	var data []byte
	if data, err = me.Raw(); err != nil {
		return err
	}
	if v != nil {
		if err = json.Unmarshal(data, &v); err != nil {
			return fmt.Errorf("%s %s: unmarshal error: %s\n%s", me.method, me.url, err.Error(), string(data))
		}
	}
	return nil
}

// var allrequests = map[string]string{}

func (me *request) Raw() ([]byte, error) {
	url := me.client.envURL + me.url
	// if _, found := allrequests[url]; found {
	// 	panic(url)
	// }
	// allrequests[url] = url
	var err error
	var body io.Reader
	var data []byte
	if me.payload != nil {
		if data, err = json.Marshal(me.payload); err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(data)
	}
	// if os.Getenv("DT_REST_DEBUG_REQUEST_PAYLOAD") == "true" && me.payload != nil {
	if len(data) > 0 {
		logger.Println(me.method, url+"\n    "+string(data))
	} else {
		logger.Println(me.method, url)
	}

	// } else {
	// logger.Println(me.method, url)
	// }

	var req *http.Request
	if req, err = http.NewRequest(me.method, url, body); err != nil {
		return nil, err
	}
	me.authenticate(req)
	if len(me.headers) > 0 {
		for headerName, headerValue := range me.headers {
			req.Header.Add(headerName, headerValue)
		}
	}
	var res *http.Response

	httpClient := &http.Client{
		Jar:       jar,
		Transport: http.DefaultTransport,
	}
	if strings.TrimSpace(os.Getenv("DYNATRACE_HTTP_INSECURE")) == "true" {
		httpClient.Transport = &http.Transport{
			ForceAttemptHTTP2:     http.DefaultTransport.(*http.Transport).ForceAttemptHTTP2,
			Proxy:                 http.DefaultTransport.(*http.Transport).Proxy,
			DialContext:           http.DefaultTransport.(*http.Transport).DialContext,
			MaxIdleConns:          http.DefaultTransport.(*http.Transport).MaxIdleConns,
			IdleConnTimeout:       http.DefaultTransport.(*http.Transport).IdleConnTimeout,
			TLSHandshakeTimeout:   http.DefaultTransport.(*http.Transport).TLSHandshakeTimeout,
			ExpectContinueTimeout: http.DefaultTransport.(*http.Transport).ExpectContinueTimeout,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		}
	} else {
		httpClient.Transport = http.DefaultTransport
	}

	if res, err = httpClient.Do(req); err != nil {
		return nil, err
	}
	if me.onResponse != nil {
		me.onResponse(res)
	}
	if data, err = io.ReadAll(res.Body); err != nil {
		return nil, err
	}
	// logger.Println("  ", res.StatusCode, string(data))
	if len(me.expect) > 0 && !me.expect.contains(res.StatusCode) {
		var env errorEnvelope
		if err = json.Unmarshal(data, &env); err == nil && env.Error != nil {
			return nil, Error{Code: env.Error.Code, Method: me.method, URL: url, Message: env.Error.Message, ConstraintViolations: env.Error.ConstraintViolations}
		} else {
			var envs []errorEnvelope
			if err = json.Unmarshal(data, &envs); err == nil && len(envs) > 0 {
				env = envs[0]
				return nil, Error{Code: env.Error.Code, Method: me.method, URL: url, Message: env.Error.Message, ConstraintViolations: env.Error.ConstraintViolations}
			}
		}
		if len(data) > 0 {
			return nil, fmt.Errorf("status code %d (expected: %d): %s", res.StatusCode, me.expect, string(data))
		}
		return nil, fmt.Errorf("status code %d (expected: %d)", res.StatusCode, me.expect)
	}
	return data, nil
}

func (me *request) Expect(codes ...int) Request {
	me.expect = statuscodes(codes)
	return me
}

func (me *request) OnResponse(onResponse func(resp *http.Response)) Request {
	me.onResponse = onResponse
	return me
}
