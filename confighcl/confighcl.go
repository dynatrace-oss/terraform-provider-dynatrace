package confighcl

import (
	"log"
	"os"
	"strings"

	"github.com/dtcookie/hcl"
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

func (me *bootstrapDecoder) GetOk(key string) (interface{}, bool) {
	res, ok := me.getRaw(key)
	if debugGetOk {
		log.Println("GetOK("+key+")", "=>", res, ok)
	}
	return res, ok
}

func (me *bootstrapDecoder) Get(key string) interface{} {
	v, _ := me.getRaw(key)
	return v
}

func (me *bootstrapDecoder) GetChange(key string) (interface{}, interface{}) {
	v := me.Get(key)
	return v, v
}

func (me *bootstrapDecoder) GetOkExists(key string) (interface{}, bool) {
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
