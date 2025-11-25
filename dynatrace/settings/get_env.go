/**
* @license
* Copyright 2025 Dynatrace LLC
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

import (
	"context"
	"os"
	"strconv"
)

func GetIntEnv(name string, def, min, max int) int {
	sValue := os.Getenv(name)
	if len(sValue) == 0 {
		return def
	}
	iValue, err := strconv.Atoi(sValue)
	if err != nil {
		return def
	}
	if iValue < min || iValue > max {
		return def
	}
	return iValue
}

func GetCtxStringValue(ctx context.Context, name string) string {
	if ctxValue := ctx.Value(name); ctxValue != nil {
		if ctxSValue, ok := ctxValue.(string); ok {
			return ctxSValue
		}
	}
	return ""
}

func GetIntEnvCtx(ctx context.Context, name string, def, min, max int) int {
	sValue := GetCtxStringValue(ctx, name)
	if len(sValue) == 0 {
		sValue = os.Getenv(name)
	}
	if len(sValue) == 0 {
		return def
	}
	iValue, err := strconv.Atoi(sValue)
	if err != nil {
		return def
	}
	if iValue < min || iValue > max {
		return def
	}
	return iValue
}
