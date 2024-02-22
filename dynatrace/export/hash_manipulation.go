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

package export

import (
	"encoding/base32"
	"encoding/binary"
	"hash/fnv"
)

func GetHashName(name string) string {

	hash := GetHash64(name)
	base32Encoded := base32encode(hash)

	return base32Encoded
}

func base32encode(num uint64) string {
	encoder := base32.NewEncoding("0123456789ABCDEFGHJKMNPQRSTVWXYZ").WithPadding(base32.NoPadding)

	numBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(numBytes, num)

	encoded := encoder.EncodeToString(numBytes)

	return encoded
}

func GetHash64(input string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(input))
	return h.Sum64()
}

func GetHash32(input string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(input))
	return h.Sum32()
}
