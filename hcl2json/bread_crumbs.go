package hcl2json

import "fmt"

type breadCrumbs string

func (me breadCrumbs) dot(v string) breadCrumbs {
	if len(me) == 0 {
		return breadCrumbs(v)
	}
	return breadCrumbs(fmt.Sprintf("%v.%v", string(me), v))
}

func (me breadCrumbs) Log(v interface{}) {
	fmt.Printf("%v: %v\n", me, v)
}
