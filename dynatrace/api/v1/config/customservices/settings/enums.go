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

package customservices

// QueueEntryPointType has no documentation
type QueueEntryPointType string

func (me QueueEntryPointType) Ref() *QueueEntryPointType {
	return &me
}

// QueueEntryPointTypes offers the known enum values
var QueueEntryPointTypes = struct {
	IBMMQ    QueueEntryPointType
	Jms      QueueEntryPointType
	Kafka    QueueEntryPointType
	MSMQ     QueueEntryPointType
	RabbitMQ QueueEntryPointType
	Values   func() []string
}{
	"IBM_MQ",
	"JMS",
	"KAFKA",
	"MSMQ",
	"RABBIT_MQ",
	func() []string {
		return []string{"IBM_MQ", "JMS", "KAFKA", "MSMQ", "RABBIT_MQ"}
	},
}

// ClassNameMatcher has no documentation
type ClassNameMatcher string

func (me ClassNameMatcher) Ref() *ClassNameMatcher {
	return &me
}

// ClassNameMatchers offers the known enum values
var ClassNameMatchers = struct {
	EndsWith   ClassNameMatcher
	Equals     ClassNameMatcher
	StartsWith ClassNameMatcher
}{
	"ENDS_WITH",
	"EQUALS",
	"STARTS_WITH",
}

// FileNameMatcher has no documentation
type FileNameMatcher string

func (me FileNameMatcher) Ref() *FileNameMatcher {
	return &me
}

// FileNameMatchers offers the known enum values
var FileNameMatchers = struct {
	EndsWith   FileNameMatcher
	Equals     FileNameMatcher
	StartsWith FileNameMatcher
}{
	"ENDS_WITH",
	"EQUALS",
	"STARTS_WITH",
}

// Visibility has no documentation
type Visibility string

// Visibilities offers the known enum values
var Visibilities = struct {
	Internal         Visibility
	PackageProtected Visibility
	Private          Visibility
	Protected        Visibility
	Public           Visibility
}{
	"INTERNAL",
	"PACKAGE_PROTECTED",
	"PRIVATE",
	"PROTECTED",
	"PUBLIC",
}

func (me Visibility) Ref() *Visibility {
	return &me
}

// Modifier has no documentation
type Modifier string

// Modifiers offers the known enum values
var Modifiers = struct {
	Abstract Modifier
	Extern   Modifier
	Final    Modifier
	Native   Modifier
	Static   Modifier
}{
	"ABSTRACT",
	"EXTERN",
	"FINAL",
	"NATIVE",
	"STATIC",
}

// Technology has no documentation
type Technology string

// Technologies offers the known enum values
var Technologies = struct {
	DotNet Technology
	Go     Technology
	Java   Technology
	NodeJS Technology
	PHP    Technology
}{
	"dotNet",
	"go",
	"java",
	"nodeJS",
	"php",
}
