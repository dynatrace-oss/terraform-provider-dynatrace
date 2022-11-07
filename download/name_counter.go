package download

import "fmt"

type NameCounter map[string]int

func (me NameCounter) Numbering(name string, m NameCounter) string {
	cnt, found := m[name]
	if !found {
		m[name] = 0
		return name
	} else {
		m[name] = cnt + 1
	}
	return fmt.Sprintf("%s(%d)", name, m[name])
}
