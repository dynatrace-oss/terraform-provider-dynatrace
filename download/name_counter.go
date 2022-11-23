package download

import "fmt"

type ReplaceFunc func(s string, cnt int) string

func DefaultReplace(s string, cnt int) string {
	return fmt.Sprintf("%s(%d)", s, cnt)
}

func ResourceName(s string, cnt int) string {
	return fmt.Sprintf("%s_%d", s, cnt)
}

type NameCounter interface {
	Numbering(string) string
	Replace(ReplaceFunc) NameCounter
}

func NewNameCounter() NameCounter {
	return &nameCounter{m: map[string]int{}}
}

type nameCounter struct {
	m       map[string]int
	replace ReplaceFunc
}

func (me *nameCounter) Replace(replace ReplaceFunc) NameCounter {
	me.replace = replace
	return me
}

func (me *nameCounter) Numbering(name string) string {
	cnt, found := me.m[name]
	if !found {
		me.m[name] = 0
		return name
	} else {
		me.m[name] = cnt + 1
	}
	if me.replace == nil {
		return DefaultReplace(name, me.m[name])
	}
	return me.replace(name, me.m[name])
}
