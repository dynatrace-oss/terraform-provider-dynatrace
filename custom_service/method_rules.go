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

import "github.com/dtcookie/dynatrace/api/config/customservices"

func extractMethodRules(value interface{}, tech customservices.Technology) []customservices.MethodRule {
	if value == nil {
		return nil
	}
	if methodRuleSections, ok := value.([]interface{}); ok {
		methodRules := []customservices.MethodRule{}
		for _, methodRuleSection := range methodRuleSections {
			methodRules = append(methodRules, *extractMethodRule(methodRuleSection, tech))
		}
		return methodRules
	}

	return nil
}

func extractFileNameMatcher(value interface{}) customservices.FileNameMatcher {
	return customservices.FileNameMatcher(value.(string))
}

func extractClassNameMatcher(value interface{}) customservices.ClassNameMatcher {
	return customservices.ClassNameMatcher(value.(string))
}

func extractMethodRule(value interface{}, tech customservices.Technology) *customservices.MethodRule {
	methodRule := customservices.MethodRule{}
	if value == nil {
		return nil
	}
	if values, ok := value.(map[string]interface{}); ok {
		if len(values) == 0 {
			return nil
		}
		methodRule.ID = values["id"].(string)
		methodRule.MethodName = values["name"].(string)
		if tech == customservices.Technologies.Go {
			returnType, found := values["returns"]
			if !found || returnType.(string) == "" {
				methodRule.ReturnType = "native"
			} else {
				methodRule.ReturnType = returnType.(string)
			}
		} else {
			methodRule.ReturnType = values["returns"].(string)
		}
		arguments := values["arguments"]
		methodRule.ArgumentTypes = []string{}
		if arguments != nil && len(arguments.([]interface{})) > 0 {
			for _, argument := range arguments.([]interface{}) {
				methodRule.ArgumentTypes = append(methodRule.ArgumentTypes, argument.(string))
			}
		}
	}
	return &methodRule
}
