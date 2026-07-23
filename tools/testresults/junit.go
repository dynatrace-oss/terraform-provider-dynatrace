/*
 * @license
 * Copyright 2026 Dynatrace LLC
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import "encoding/xml"

// The following types model the subset of the JUnit XML schema emitted by
// gotestsum (`--junitfile`). A file has a root <testsuites> element containing
// one or more <testsuite> elements, each holding <testcase> elements. A
// testcase is considered failed when it contains a <failure> or <error> child.

type testSuites struct {
	XMLName  xml.Name    `xml:"testsuites"`
	Name     string      `xml:"name,attr,omitempty"`
	Tests    int         `xml:"tests,attr"`
	Failures int         `xml:"failures,attr"`
	Errors   int         `xml:"errors,attr"`
	Suites   []testSuite `xml:"testsuite"`
}

type testSuite struct {
	XMLName   xml.Name   `xml:"testsuite"`
	Name      string     `xml:"name,attr,omitempty"`
	Tests     int        `xml:"tests,attr"`
	Failures  int        `xml:"failures,attr"`
	Errors    int        `xml:"errors,attr"`
	Skipped   int        `xml:"skipped,attr,omitempty"`
	Time      string     `xml:"time,attr,omitempty"`
	Timestamp string     `xml:"timestamp,attr,omitempty"`
	TestCases []testCase `xml:"testcase"`
}

type testCase struct {
	XMLName   xml.Name `xml:"testcase"`
	Name      string   `xml:"name,attr"`
	ClassName string   `xml:"classname,attr,omitempty"`
	Time      string   `xml:"time,attr,omitempty"`
	Failure   *failure `xml:"failure"`
	Error     *failure `xml:"error"`
	Skipped   *skipped `xml:"skipped"`
}

type failure struct {
	Message string `xml:"message,attr,omitempty"`
	Type    string `xml:"type,attr,omitempty"`
	Content string `xml:",chardata"`
}

type skipped struct {
	Message string `xml:"message,attr,omitempty"`
}

// failed reports whether the testcase represents a failure or error.
func (tc testCase) failed() bool {
	return tc.Failure != nil || tc.Error != nil
}

// xmlMarshalIndent serializes the testSuites to indented XML with a header.
func xmlMarshalIndent(suites testSuites) ([]byte, error) {
	body, err := xml.MarshalIndent(suites, "", "  ")
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), append(body, '\n')...), nil
}
