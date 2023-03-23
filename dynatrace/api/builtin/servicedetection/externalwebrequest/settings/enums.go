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

package externalwebrequest

type ContextRootTransformationType string

var ContextRootTransformationTypes = struct {
	Before            ContextRootTransformationType
	RemoveCreditCards ContextRootTransformationType
	RemoveIbans       ContextRootTransformationType
	RemoveIps         ContextRootTransformationType
	RemoveNumbers     ContextRootTransformationType
	ReplaceBetween    ContextRootTransformationType
}{
	"BEFORE",
	"REMOVE_CREDIT_CARDS",
	"REMOVE_IBANS",
	"REMOVE_IPS",
	"REMOVE_NUMBERS",
	"REPLACE_BETWEEN",
}

type ContributionType string

var ContributionTypes = struct {
	Originalvalue  ContributionType
	Overridevalue  ContributionType
	Transformurl   ContributionType
	Transformvalue ContributionType
}{
	"OriginalValue",
	"OverrideValue",
	"TransformURL",
	"TransformValue",
}

type ContributionTypeWithOverride string

var ContributionTypeWithOverrides = struct {
	Originalvalue  ContributionTypeWithOverride
	Overridevalue  ContributionTypeWithOverride
	Transformvalue ContributionTypeWithOverride
}{
	"OriginalValue",
	"OverrideValue",
	"TransformValue",
}

type FrameworkType string

var FrameworkTypes = struct {
	Axis       FrameworkType
	Cxf        FrameworkType
	Hessian    FrameworkType
	JaxWsRi    FrameworkType
	Jboss      FrameworkType
	Jersey     FrameworkType
	Progress   FrameworkType
	Resteasy   FrameworkType
	Restlet    FrameworkType
	Spring     FrameworkType
	Tibco      FrameworkType
	Weblogic   FrameworkType
	Webmethods FrameworkType
	Websphere  FrameworkType
	Wink       FrameworkType
}{
	"AXIS",
	"CXF",
	"HESSIAN",
	"JAX_WS_RI",
	"JBOSS",
	"JERSEY",
	"PROGRESS",
	"RESTEASY",
	"RESTLET",
	"SPRING",
	"TIBCO",
	"WEBLOGIC",
	"WEBMETHODS",
	"WEBSPHERE",
	"WINK",
}

type TransformationType string

var TransformationTypes = struct {
	After             TransformationType
	Before            TransformationType
	Between           TransformationType
	RemoveCreditCards TransformationType
	RemoveIbans       TransformationType
	RemoveIps         TransformationType
	RemoveNumbers     TransformationType
	ReplaceBetween    TransformationType
	SplitSelect       TransformationType
	TakeSegments      TransformationType
}{
	"AFTER",
	"BEFORE",
	"BETWEEN",
	"REMOVE_CREDIT_CARDS",
	"REMOVE_IBANS",
	"REMOVE_IPS",
	"REMOVE_NUMBERS",
	"REPLACE_BETWEEN",
	"SPLIT_SELECT",
	"TAKE_SEGMENTS",
}
