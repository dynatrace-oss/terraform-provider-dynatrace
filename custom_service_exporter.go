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
	"bytes"
	"fmt"
	"io"
	"unicode"

	"github.com/dtcookie/dynatrace/api/config/customservices"
	"github.com/dtcookie/dynatrace/terraform"
)

func terraformat(s string) string {
	result := ""
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' {
			result = result + string(r)
		} else {
			result = result + "_"
		}
	}
	return result
}

// CustomServiceExporter has no documentation
type CustomServiceExporter struct {
	CustomService *customservices.CustomService
}

func (cse *CustomServiceExporter) printCustomService(customService *customservices.CustomService, buf *bytes.Buffer, technology string) error {
	fmt.Fprintf(buf, "resource \"%s\" \"%s\" {\n", "dynatrace_custom_service", terraformat(customService.Name))
	fmt.Fprintf(buf, "	name = \"%s\"\n", customService.Name)
	fmt.Fprintf(buf, "	technology = \"%s\"\n", technology)
	fmt.Fprintf(buf, "	enabled = %v\n", customService.Enabled)
	if customService.QueueEntryPoint {
		fmt.Fprintf(buf, "	queue_entry_point = %v\n", customService.QueueEntryPoint)
	}
	if len(customService.QueueEntryPointType) > 0 {
		fmt.Fprintf(buf, "	queue_entry_point_type = \"%v\"\n", customService.QueueEntryPointType)
	}
	for _, detectionRule := range customService.Rules {
		cse.printDetectionRule(&detectionRule, buf, customservices.Technology(technology))
	}
	fmt.Fprintln(buf, "}")
	return nil
}

func (cse *CustomServiceExporter) printDetectionRule(detectionRule *customservices.DetectionRule, buf *bytes.Buffer, tech customservices.Technology) error {
	fmt.Fprintln(buf, "	rule {")
	fmt.Fprintf(buf, "		enabled = %v\n", detectionRule.Enabled)
	if len(detectionRule.ClassName) > 0 || len(detectionRule.ClassNameMatcher) > 0 {
		fmt.Fprintln(buf, "		class {")
		fmt.Fprintf(buf, "			name = \"%s\"\n", detectionRule.ClassName)
		fmt.Fprintf(buf, "			match = \"%v\"\n", detectionRule.ClassNameMatcher)
		fmt.Fprintln(buf, "		}")
	}
	if len(detectionRule.FileName) > 0 || len(detectionRule.FileNameMatcher) > 0 {
		fmt.Fprintln(buf, "		file {")
		fmt.Fprintf(buf, "			name = \"%s\"\n", detectionRule.FileName)
		fmt.Fprintf(buf, "			match = \"%v\"\n", detectionRule.FileNameMatcher)
		fmt.Fprintln(buf, "		}")
	}
	if len(detectionRule.Annotations) > 0 {
		annotationLine := "		annotations = ["
		sep := " "
		for _, annotation := range detectionRule.Annotations {
			annotationLine = annotationLine + sep
			annotationLine = annotationLine + quote(annotation)
			sep = ", "
		}
		annotationLine = annotationLine + " ]"
		fmt.Fprintln(buf, annotationLine)
	}
	for _, methodRule := range detectionRule.MethodRules {
		cse.printMethodRule(&methodRule, buf, tech)
	}
	fmt.Fprintln(buf, "	}")
	return nil
}

func quote(s string) string {
	return fmt.Sprintf("\"%s\"", s)
}

func (cse *CustomServiceExporter) printMethodRule(methodrule *customservices.MethodRule, buf *bytes.Buffer, tech customservices.Technology) error {
	fmt.Fprintln(buf, "		method {")
	fmt.Fprintf(buf, "			name = \"%s\"\n", methodrule.MethodName)
	if len(methodrule.ArgumentTypes) > 0 {
		argLine := "			arguments = ["
		sep := " "
		for _, argument := range methodrule.ArgumentTypes {
			argLine = argLine + sep
			argLine = argLine + quote(argument)
			sep = ", "
		}
		argLine = argLine + " ]"
		fmt.Fprintln(buf, argLine)
	}
	if tech != customservices.Technologies.Go || methodrule.ReturnType != "native" {
		fmt.Fprintf(buf, "			returns = \"%s\"\n", methodrule.ReturnType)
	}

	fmt.Fprintln(buf, "		}")
	return nil
}

// ToHCL has no documentation
func (cse *CustomServiceExporter) ToHCL(w io.Writer, technology string) error {
	customService := *cse.CustomService
	customService.ID = ""
	customService.Metadata = nil
	customService.ProcessGroups = nil
	if customService.Rules != nil {
		for detectionRuleIdx := range customService.Rules {
			customService.Rules[detectionRuleIdx].ID = ""
			if customService.Rules[detectionRuleIdx].MethodRules != nil {
				for methodRuleIdx := range customService.Rules[detectionRuleIdx].MethodRules {
					customService.Rules[detectionRuleIdx].MethodRules[methodRuleIdx].ID = ""
				}
			}
		}
	}

	buf := bytes.NewBuffer([]byte{})
	cse.printCustomService(&customService, buf, technology)
	_, err := w.Write(buf.Bytes())
	return err
}

// ToJSON has no documentation
func (cse *CustomServiceExporter) ToJSON(w io.Writer) error {
	bytes, err := terraform.MarshalJSON(cse.CustomService, "dynatrace_custom_service", cse.CustomService.Name)
	if err != nil {
		return err
	}
	_, err = w.Write(bytes)
	return err
}
