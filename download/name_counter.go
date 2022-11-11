package download

import "fmt"

type NameCounter map[string]int

func (me NameCounter) Numbering(name string) string {
	cnt, found := me[name]
	if !found {
		me[name] = 0
		return name
	} else {
		me[name] = cnt + 1
	}
	return fmt.Sprintf("%s(%d)", name, me[name])
}

// type ReplaceFunc func(s string, cnt int) string

// func DefaultReplace(s string, cnt int) string {
// 	return fmt.Sprintf("%s(%d)", s, cnt)
// }

// type NameCounter struct {
// 	m       map[string]int
// 	Replace ReplaceFunc
// }

// func (me NameCounter) Numbering(name string) string {
// 	cnt, found := me.m[name]
// 	if !found {
// 		me.m[name] = 0
// 		return name
// 	} else {
// 		me.m[name] = cnt + 1
// 	}
// 	if me.Replace == nil {
// 		return DefaultReplace(name, me.m[name])
// 	}
// 	return me.Replace(name, me.m[name])
// }
