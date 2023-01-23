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

package settings

type DemoSettings interface {
	FillDemoValues() []string
}

func FillDemoValues(settings Settings) []string {
	if demoSettings, ok := settings.(DemoSettings); ok {
		return demoSettings.FillDemoValues()
	}
	return []string{}
}

type RegexValidator interface {
	Validate() []string
}

func Validate(settings Settings) []string {
	m := map[string]string{}
	if demoSettings, ok := settings.(RegexValidator); ok {
		messages := demoSettings.Validate()
		if len(messages) == 0 {
			return []string{}
		}
		for _, message := range messages {
			m[message] = message
		}
		result := []string{}
		for k := range m {
			result = append(result, k)
		}
		return result
	}
	return []string{}
}
