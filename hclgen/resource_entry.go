package hclgen

import (
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type resourceEntry struct {
	Indent      string
	Key         string
	BreadCrumbs string
	Optional    bool
	Entries     exportEntries
}

func (me *resourceEntry) IsOptional() bool {
	return false
}

func (me *resourceEntry) IsDefault() bool {
	return false
}

func (me *resourceEntry) IsLessThan(other exportEntry) bool {
	switch ro := other.(type) {
	case *primitiveEntry:
		return false
	case *resourceEntry:
		return strings.Compare(me.Key, ro.Key) < 0
	}
	return false
}

func (me *resourceEntry) Write(w *hclwrite.Body, indent string) error {
	block := w.AppendNewBlock(me.Key, nil)
	body := block.Body()

	sort.SliceStable(me.Entries, me.Entries.Less)
	for _, entry := range me.Entries {
		if !(entry.IsOptional() && entry.IsDefault()) {
			if err := entry.Write(body, indent+"  "); err != nil {
				return err
			}
		} else {
			body.AppendUnstructuredTokens(hclwrite.Tokens{
				&hclwrite.Token{Type: hclsyntax.TokenComment, Bytes: []byte("#")},
			})
			if err := entry.Write(body, indent+"  "); err != nil {
				return err
			}
		}
	}
	return nil
}
