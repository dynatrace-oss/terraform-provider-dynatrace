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

package rest

const (
	AuthAPIToken        = AuthType(1)
	AuthOAuth           = AuthType(2)
	AuthClusterAPIToken = AuthType(3)
)

type AuthType uint8

func (a AuthType) APIToken() bool {
	return a&AuthAPIToken != 0
}

func (a AuthType) OAuth() bool {
	return a&AuthOAuth != 0
}

func (a AuthType) ClusterAPIToken() bool {
	return a&AuthClusterAPIToken != 0
}
