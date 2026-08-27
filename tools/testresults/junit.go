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

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type testSuites struct {
	XMLName xml.Name    `xml:"testsuites"`
	Name    string      `xml:"name,attr,omitempty"`
	Suites  []testSuite `xml:"testsuite"`
}

type testSuite struct {
	XMLName  xml.Name `xml:"testsuite"`
	Name     string   `xml:"name,attr,omitempty"`
	Tests    int      `xml:"tests,attr"`
	Failures int      `xml:"failures,attr"`
}

func findAllTestSuites(path string) ([]testSuite, []error) {
	var allSuites []testSuite
	errors := []error{}
	files, err := findXMLFiles(path)
	if err != nil {
		errors = append(errors, err)
	}

	for _, file := range files {
		suites, err := parseFile(file)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		allSuites = append(allSuites, suites.Suites...)
	}
	return allSuites, errors
}

func findXMLFiles(path string) ([]string, error) {
	var files []string
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(p) == ".xml" {
			files = append(files, p)
		}
		return nil
	})
	return files, err
}

func parseFile(path string) (testSuites, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return testSuites{}, fmt.Errorf("error reading junit file %s: %w", path, err)
	}
	var suites testSuites
	if err := xml.Unmarshal(data, &suites); err != nil {
		return testSuites{}, fmt.Errorf("error parsing junit file %s: %w", path, err)
	}
	return suites, nil
}

func nameToPath(name string, basePath string, repo string) (string, error) {
	if !strings.HasPrefix(name, repo) {
		return "", fmt.Errorf("expected name %s to have prefix %s", name, repo)
	}
	name = strings.TrimPrefix(name, repo)

	return filepath.Join(basePath, filepath.FromSlash(name)), nil
}
