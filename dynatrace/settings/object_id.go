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

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/google/uuid"
)

type ObjectID struct {
	Key   string `json:"key"`
	Scope struct {
		Class string `json:"class"`
		ID    string `json:"id"`
	} `json:"scope"`
	SchemaID string `json:"schemaID"`
	ID       string `json:"-"`
}

func (me *ObjectID) String() string {
	data, _ := json.Marshal(me)
	return string(data)
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func (me *ObjectID) Decode() error {
	var err error
	id := me.ID
	if len(id) < 64 {
		return nil
	}
	if IsValidUUID(id) {
		return nil
	}
	if strings.Contains(id, "-") {
		id = strings.ReplaceAll(id, "-", "")
	}

	// if strings.Contains(id, "-") {
	// 	fmt.Println(id)
	// 	id = id[:strings.Index(id, "-")]
	// 	fmt.Println("  ->", id)

	// }
	encodedBytes := []byte(id)

	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer(encodedBytes))

	resultBytes := make([]byte, len(encodedBytes))
	len, err := decoder.Read(resultBytes)
	if err != nil {
		return fmt.Errorf("decoder error: %s", err.Error())
	}
	resultBytes = resultBytes[0:len]

	buf := bytes.NewBuffer(resultBytes[12:])

	if me.SchemaID, err = ReadUTF(buf); err != nil {
		return err
	}
	if me.Scope.Class, err = ReadUTF(buf); err != nil {
		return err
	}
	if me.Scope.ID, err = ReadUTF(buf); err != nil {
		return err
	}
	if me.Key, err = ReadUTF(buf); err != nil {
		return err
	}
	return nil
}

func ReadUTF(reader io.Reader) (string, error) {
	var utfLen uint16
	err := binary.Read(reader, binary.BigEndian, &utfLen)
	if err != nil {
		return "", err
	}
	return readUTFBytes(reader, int(utfLen))
}

func readUTFBytes(reader io.Reader, utfLen int) (string, error) {
	byteArr := make([]byte, utfLen)
	_, err := reader.Read(byteArr)
	if err != nil {
		return "", err
	}
	var charArr bytes.Buffer
	charArr.Grow(utfLen / 2)
	var count int
	var c byte

	for count < utfLen {
		c = byteArr[count]
		if c > 127 {
			break
		}
		count++
		charArr.WriteByte(c)
	}
	var surrogateHi int32
	for count < utfLen {
		c = byteArr[count]
		msb4Bits := c >> 4
		if msb4Bits < 8 {
			count++
			charArr.WriteByte(c)
		} else if msb4Bits == 12 || msb4Bits == 13 {
			count += 2
			if count > utfLen {
				return charArr.String(), &errorString{"truncated character"}
			}
			c2 := byteArr[count-1]
			if c2&0xc0 != 0x80 {
				return charArr.String(), &errorString{"malformed input"}
			}
			r := rune((int32(c&0x1f) << 6) | int32(c2&0x3f))
			charArr.WriteRune(r)
		} else if msb4Bits == 14 {
			count += 3
			if count > utfLen {
				return charArr.String(), &errorString{"truncated character"}
			}
			c2 := byteArr[count-2]
			c3 := byteArr[count-1]
			if c2&0xc0 != 0x80 || c3&0xc0 != 0x80 {
				return charArr.String(), &errorString{"malformed input"}
			}
			r := rune((int32(c&0x0f) << 12) | (int32(c2&0x3f) << 6) | int32(c3&0x3f))
			if r < surrogateMin || r > surrogateMax {
				charArr.WriteRune(r)
			} else {
				if r < 0xDBFF {
					surrogateHi = r
				} else {
					surrogateLo := r
					r = (((surrogateHi - 0xD800) << 10) | (surrogateLo - 0xDC00)) + 0x10000
					charArr.WriteRune(r)
				}
			}
		} else {
			return charArr.String(), &errorString{"malformed input, maybe it's regular UTF-8?"}
		}
	}
	return charArr.String(), nil
}

const surrogateMin = 0xD800
const surrogateMax = 0xDFFF

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
