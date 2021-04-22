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

package customservice

import (
	"github.com/dtcookie/dynatrace/api/config/customservices"
)

func extractDetectionRules(value interface{}, tech customservices.Technology) []customservices.DetectionRule {
	if value == nil {
		return nil
	}
	if detectionRuleSections, ok := value.([]interface{}); ok {
		detectionRules := []customservices.DetectionRule{}
		for _, detectionRuleSection := range detectionRuleSections {
			detectionRule := *extractDetectionRule(detectionRuleSection, tech)
			detectionRules = append(detectionRules, detectionRule)
		}
		return detectionRules
	}
	return nil
}

func extractDetectionRule(value interface{}, tech customservices.Technology) *customservices.DetectionRule {
	detectionRule := &customservices.DetectionRule{}
	if value == nil {
		return nil
	}
	if values, ok := value.(map[string]interface{}); ok {
		if len(values) == 0 {
			return nil
		}
		detectionRule.ID = values["id"].(string)
		detectionRule.Enabled = values["enabled"].(bool)
		annotations := values["annotations"]
		detectionRule.Annotations = []string{}
		if annotations != nil {
			for _, annotation := range annotations.([]interface{}) {
				detectionRule.Annotations = append(detectionRule.Annotations, annotation.(string))
			}
		}
		if fileSections := values["file"]; fileSections != nil {
			for _, fileSection := range fileSections.([]interface{}) {
				detectionRule.FileName = fileSection.(map[string]interface{})["name"].(string)
				detectionRule.FileNameMatcher = extractFileNameMatcher(fileSection.(map[string]interface{})["match"])
			}
		}

		if classSections := values["class"]; classSections != nil {
			for _, classSection := range classSections.([]interface{}) {
				detectionRule.ClassName = classSection.(map[string]interface{})["name"].(string)
				detectionRule.ClassNameMatcher = extractClassNameMatcher(classSection.(map[string]interface{})["match"])
			}
		}

		detectionRule.MethodRules = extractMethodRules(values["method"], tech)

	}
	return detectionRule
}
