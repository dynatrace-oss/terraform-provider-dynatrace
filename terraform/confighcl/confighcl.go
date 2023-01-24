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

package confighcl

import (
	"log"
	"os"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var debugGetOk = (os.Getenv("DT_DEBUG_GET_OK") == "true")

type bootstrapDecoder struct {
	Reader schema.FieldReader
}

func (me *bootstrapDecoder) getRaw(key string) (any, bool) {
	var parts []string
	if key != "" {
		parts = strings.Split(key, ".")
	}
	result, err := me.Reader.ReadField(parts)
	if err != nil {
		panic(err)
	}
	return result.Value, result.Exists
}

func (me *bootstrapDecoder) GetOk(key string) (any, bool) {
	res, ok := me.getRaw(key)
	if debugGetOk {
		log.Println("GetOK("+key+")", "=>", res, ok)
	}
	return res, ok
}

func (me *bootstrapDecoder) Get(key string) any {
	v, _ := me.getRaw(key)
	return v
}

func (me *bootstrapDecoder) GetChange(key string) (any, any) {
	v := me.Get(key)
	return v, v
}

func (me *bootstrapDecoder) GetOkExists(key string) (any, bool) {
	return me.GetOk(key)
}

func (me *bootstrapDecoder) HasChange(key string) bool {
	return false
}

func DecoderFrom(d *schema.ResourceData, res *schema.Resource) hcl.Decoder {
	return hcl.DecoderFrom(&bootstrapDecoder{&schema.ConfigFieldReader{
		Config: terraform.NewResourceConfigShimmed(d.GetRawConfig(), res.CoreConfigSchema()),
		Schema: res.Schema,
	}})
}

func StateDecoderFrom(d *schema.ResourceData, res *schema.Resource) hcl.Decoder {
	return hcl.DecoderFrom(&bootstrapDecoder{&schema.ConfigFieldReader{
		Config: terraform.NewResourceConfigShimmed(d.GetRawState(), res.CoreConfigSchema()),
		Schema: res.Schema,
	}})
}
