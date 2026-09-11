//go:build unit

/*
 * @license
 * Copyright 2026 Dynatrace LLC
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package hclgen

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	hclv2 "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeMarshaler is a hcl.Marshaler used to drive HCLGen.Export in tests.
type fakeMarshaler struct {
	fn func(hcl.Properties) error
}

func (f *fakeMarshaler) MarshalHCL(p hcl.Properties) error {
	return f.fn(p)
}

// fakeSchemer additionally implements hcl.Schemer.
type fakeSchemer struct {
	fakeMarshaler
	schema map[string]*schema.Schema
}

func (f *fakeSchemer) Schema() map[string]*schema.Schema {
	return f.schema
}

// renderBlockBody invokes the given callback against a freshly created
// "resource" block and returns the rendered HCL.
func renderBlockBody(t *testing.T, fn func(block *hclwrite.Block)) string {
	t.Helper()
	file := hclwrite.NewEmptyFile()
	block := file.Body().AppendNewBlock("resource", []string{"some_type", "some_name"})
	fn(block)
	return string(file.Bytes())
}

func TestWriteEntries(t *testing.T) {
	t.Run("empty entries produces empty body", func(t *testing.T) {
		out := renderBlockBody(t, func(block *hclwrite.Block) {
			require.NoError(t, WriteEntries(exportEntries{}, block, ""))
		})

		assert.Contains(t, out, `resource "some_type" "some_name"`)
		assert.NotContains(t, out, "=")
	})

	t.Run("single primitive entry renders attribute", func(t *testing.T) {
		out := renderBlockBody(t, func(block *hclwrite.Block) {
			entries := exportEntries{
				&primitiveEntry{Key: "name", Value: "my-resource"},
			}
			require.NoError(t, WriteEntries(entries, block, ""))
		})

		assert.Contains(t, out, `name = "my-resource"`)
	})

	t.Run("optional default non-computed entry is commented out", func(t *testing.T) {
		out := renderBlockBody(t, func(block *hclwrite.Block) {
			// bool with value false is treated as "default" for non-DefValTrue entries
			entries := exportEntries{
				&primitiveEntry{Key: "enabled", Value: false, Optional: true, Computed: false},
			}
			require.NoError(t, WriteEntries(entries, block, ""))
		})

		assert.Contains(t, out, "# enabled")
	})

	t.Run("computed entries skip the comment prefix", func(t *testing.T) {
		out := renderBlockBody(t, func(block *hclwrite.Block) {
			entries := exportEntries{
				&primitiveEntry{Key: "enabled", Value: false, Optional: true, Computed: true},
			}
			require.NoError(t, WriteEntries(entries, block, ""))
		})

		assert.NotContains(t, out, "# enabled")
		assert.Contains(t, out, "enabled = false")
	})

	t.Run("non-optional entry skips the comment prefix", func(t *testing.T) {
		out := renderBlockBody(t, func(block *hclwrite.Block) {
			entries := exportEntries{
				&primitiveEntry{Key: "enabled", Value: false, Optional: false, Computed: false},
			}
			require.NoError(t, WriteEntries(entries, block, ""))
		})

		assert.NotContains(t, out, "# enabled")
		assert.Contains(t, out, "enabled = false")
	})

	t.Run("entry write error is propagated", func(t *testing.T) {
		// primitiveEntry.Write returns an error when ctyVal returns NilVal —
		// an unsupported type satisfies that.
		entries := exportEntries{
			&primitiveEntry{Key: "broken", Value: struct{}{}},
		}
		file := hclwrite.NewEmptyFile()
		block := file.Body().AppendNewBlock("resource", []string{"t", "n"})

		err := WriteEntries(entries, block, "")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "broken")
	})

	t.Run("multiple entries are written in order", func(t *testing.T) {
		out := renderBlockBody(t, func(block *hclwrite.Block) {
			entries := exportEntries{
				&primitiveEntry{Key: "alpha", Value: "1"},
				&primitiveEntry{Key: "beta", Value: "2"},
				&primitiveEntry{Key: "gamma", Value: "3"},
			}
			require.NoError(t, WriteEntries(entries, block, ""))
		})

		alphaIdx := strings.Index(out, "alpha")
		betaIdx := strings.Index(out, "beta")
		gammaIdx := strings.Index(out, "gamma")

		require.NotEqual(t, -1, alphaIdx)
		require.NotEqual(t, -1, betaIdx)
		require.NotEqual(t, -1, gammaIdx)

		assert.Less(t, alphaIdx, betaIdx)
		assert.Less(t, betaIdx, gammaIdx)
	})
}

// renderPrimitive writes a single primitiveEntry with key "message" into a
// resource block and returns the rendered HCL.
func renderPrimitive(t *testing.T, value any) string {
	t.Helper()
	return renderBlockBody(t, func(block *hclwrite.Block) {
		require.NoError(t, (&primitiveEntry{Key: "message", Value: value}).Write(block.Body(), ""))
	})
}

// assertHCLParses guards against emitting HCL that doesn't parse (e.g. a heredoc
// terminator sharing its line with following tokens).
func assertHCLParses(t *testing.T, out string) {
	t.Helper()
	_, diags := hclsyntax.ParseConfig([]byte(out), "test.tf", hclv2.InitialPos)
	require.False(t, diags.HasErrors(), diags.Error())
}

func TestPrimitiveEntryHeredoc(t *testing.T) {
	// heredoc rendering is gated on envutils.DynatraceHeredoc, which defaults to true.
	t.Run("multiline value ending in newline renders a bare heredoc", func(t *testing.T) {
		out := renderPrimitive(t, "line1\nline2\n")
		assertHCLParses(t, out)
		assert.Contains(t, out, "<<-EOT")
		assert.NotContains(t, out, "trimsuffix")
	})

	t.Run("multiline value without trailing newline is wrapped in trimsuffix", func(t *testing.T) {
		out := renderPrimitive(t, "line1\nline2")
		assertHCLParses(t, out)
		assert.Contains(t, out, "trimsuffix(<<-EOT")
		assert.Contains(t, out, `, "\n")`)
	})

	t.Run("single-line value stays a quoted string", func(t *testing.T) {
		out := renderPrimitive(t, "just one line")
		assertHCLParses(t, out)
		assert.Contains(t, out, `message = "just one line"`)
		assert.NotContains(t, out, "EOT")
	})

	t.Run("string pointer is handled like a string", func(t *testing.T) {
		v := "line1\nline2"
		out := renderPrimitive(t, &v)
		assertHCLParses(t, out)
		assert.Contains(t, out, "trimsuffix(<<-EOT")
	})

	t.Run("EOT terminator is dedented to the attribute indent", func(t *testing.T) {
		out := renderPrimitive(t, "line1\nline2\n")
		// indent "" → content at 2 spaces, EOT flush left.
		assert.Contains(t, out, "\n  line1\n")
		assert.Contains(t, out, "\nEOT\n")
	})
}

// TestPrimitiveEntryWriteBranches covers the remaining primitiveEntry.Write
// branches (the "HCL-UNQUOTE-" prefix, string slices, JSON detection, quote-heavy
// single-line strings, and non-string values).
func TestPrimitiveEntryWriteBranches(t *testing.T) {
	ptr := func(v string) *string { return &v }

	t.Run("HCL-UNQUOTE- prefix renders an unquoted expression", func(t *testing.T) {
		out := renderPrimitive(t, "HCL-UNQUOTE-var.foo")
		assertHCLParses(t, out)
		assert.Contains(t, out, "message = var.foo")
		assert.NotContains(t, out, `"var.foo"`)
	})

	t.Run("HCL-UNQUOTE- prefix on a *string renders unquoted", func(t *testing.T) {
		out := renderPrimitive(t, ptr("HCL-UNQUOTE-var.bar"))
		assertHCLParses(t, out)
		assert.Contains(t, out, "message = var.bar")
	})

	t.Run("string slice renders a list", func(t *testing.T) {
		out := renderPrimitive(t, []string{"a", "b"})
		assertHCLParses(t, out)
		assert.Contains(t, out, `message = [ "a", "b" ]`)
	})

	t.Run("string slice keeps HCL-UNQUOTE- elements unquoted", func(t *testing.T) {
		out := renderPrimitive(t, []string{"HCL-UNQUOTE-var.x", "y"})
		assertHCLParses(t, out)
		assert.Contains(t, out, `message = [ var.x, "y" ]`)
	})

	t.Run("JSON string renders via jsonencode", func(t *testing.T) {
		out := renderPrimitive(t, `{"a":1,"b":"x"}`)
		assertHCLParses(t, out)
		assert.Contains(t, out, "jsonencode(")
		assert.Contains(t, out, `"a": 1`)
	})

	t.Run("JSON *string renders via jsonencode", func(t *testing.T) {
		out := renderPrimitive(t, ptr(`{"k":true}`))
		assertHCLParses(t, out)
		assert.Contains(t, out, "jsonencode(")
		assert.Contains(t, out, `"k": true`)
	})

	t.Run("quote-heavy single-line string renders as a trimsuffix heredoc", func(t *testing.T) {
		// Matches the strings.Count(s, `"`) > 3 heredoc branch: a value with many
		// embedded quotes avoids escape noise via a heredoc. Being single-line it
		// has no trailing newline, so it's wrapped in trimsuffix to stay exact.
		out := renderPrimitive(t, `a "b" c "d" e`)
		assertHCLParses(t, out)
		assert.Contains(t, out, "trimsuffix(<<-EOT")
		assert.Contains(t, out, `a "b" c "d" e`)
		assert.NotContains(t, out, `\"`) // quotes are literal inside the heredoc
	})

	t.Run("string with few quotes stays a quoted string", func(t *testing.T) {
		out := renderPrimitive(t, `a "b" c`)
		assertHCLParses(t, out)
		assert.Contains(t, out, `message = "a \"b\" c"`)
		assert.NotContains(t, out, "EOT")
	})

	t.Run("int value renders as a bare number", func(t *testing.T) {
		out := renderPrimitive(t, 42)
		assertHCLParses(t, out)
		assert.Contains(t, out, "message = 42")
	})

	t.Run("bool value renders as a bare literal", func(t *testing.T) {
		out := renderPrimitive(t, true)
		assertHCLParses(t, out)
		assert.Contains(t, out, "message = true")
	})
}

func TestHCLGen(t *testing.T) {
	t.Run("New initializes kind and file", func(t *testing.T) {
		g := New("resource")

		require.NotNil(t, g)
		assert.Equal(t, "resource", g.kind)
		assert.NotNil(t, g.file)
	})

	t.Run("Export writes a resource block with attributes", func(t *testing.T) {
		marshaler := &fakeMarshaler{
			fn: func(p hcl.Properties) error {
				p["name"] = "my-name"
				return nil
			},
		}
		buf := &bytes.Buffer{}

		err := ExportResource(marshaler, buf, "dynatrace_thing", "example")

		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, `resource "dynatrace_thing" "example"`)
		assert.Contains(t, out, `name = "my-name"`)
	})

	t.Run("ExportDataSource writes a data block", func(t *testing.T) {
		marshaler := &fakeMarshaler{
			fn: func(p hcl.Properties) error {
				p["id"] = "abc"
				return nil
			},
		}
		buf := &bytes.Buffer{}

		err := ExportDataSource(marshaler, buf, "dynatrace_lookup", "example")

		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, `data "dynatrace_lookup" "example"`)
		assert.Contains(t, out, `id = "abc"`)
	})

	t.Run("Export prepends non-empty comments", func(t *testing.T) {
		marshaler := &fakeMarshaler{
			fn: func(p hcl.Properties) error {
				p["name"] = "x"
				return nil
			},
		}
		buf := &bytes.Buffer{}

		err := ExportResource(marshaler, buf, "dynatrace_thing", "example",
			"first comment", "  ", "", "second comment")

		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "# first comment")
		assert.Contains(t, out, "# second comment")
		// blank/whitespace-only comments are dropped
		assert.Equal(t, 2, strings.Count(out, "#"))

		// comments come before the resource block
		commentIdx := strings.Index(out, "# first comment")
		blockIdx := strings.Index(out, `resource "dynatrace_thing"`)
		require.NotEqual(t, -1, commentIdx)
		require.NotEqual(t, -1, blockIdx)
		assert.Less(t, commentIdx, blockIdx)
	})

	t.Run("Export propagates MarshalHCL errors", func(t *testing.T) {
		boom := errors.New("boom")
		marshaler := &fakeMarshaler{
			fn: func(p hcl.Properties) error { return boom },
		}
		buf := &bytes.Buffer{}

		err := ExportResource(marshaler, buf, "t", "n")

		require.ErrorIs(t, err, boom)
		assert.Empty(t, buf.String())
	})

	t.Run("Export uses Schemer to drive optional/computed handling", func(t *testing.T) {
		// `enabled` is Optional+Computed bool → resOpt is true, resComputed is true.
		// With value=false this would still write the attribute (no `#` prefix because Computed).
		schemer := &fakeSchemer{
			fakeMarshaler: fakeMarshaler{
				fn: func(p hcl.Properties) error {
					p["enabled"] = false
					return nil
				},
			},
			schema: map[string]*schema.Schema{
				"enabled": {Type: schema.TypeBool, Optional: true, Computed: true},
			},
		}
		buf := &bytes.Buffer{}

		err := ExportResource(schemer, buf, "t", "n")

		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "enabled = false")
		assert.NotContains(t, out, "#enabled")
	})

	t.Run("Export comments out optional non-computed default values via schema", func(t *testing.T) {
		// `enabled` is Optional (not Computed) → optional+default+not-computed → `#` prefix.
		schemer := &fakeSchemer{
			fakeMarshaler: fakeMarshaler{
				fn: func(p hcl.Properties) error {
					p["enabled"] = false
					return nil
				},
			},
			schema: map[string]*schema.Schema{
				"enabled": {Type: schema.TypeBool, Optional: true},
			},
		}
		buf := &bytes.Buffer{}

		err := ExportResource(schemer, buf, "t", "n")

		require.NoError(t, err)
		assert.Contains(t, buf.String(), "# enabled")
	})

	t.Run("Export drops computed-only fields", func(t *testing.T) {
		// Computed-only (not Optional) → eval skips the entry entirely.
		schemer := &fakeSchemer{
			fakeMarshaler: fakeMarshaler{
				fn: func(p hcl.Properties) error {
					p["id"] = "generated"
					p["name"] = "kept"
					return nil
				},
			},
			schema: map[string]*schema.Schema{
				"id":   {Type: schema.TypeString, Computed: true},
				"name": {Type: schema.TypeString, Required: true},
			},
		}
		buf := &bytes.Buffer{}

		err := ExportResource(schemer, buf, "t", "n")

		require.NoError(t, err)
		out := buf.String()
		assert.NotContains(t, out, "id ")
		assert.Contains(t, out, `name = "kept"`)
	})

	t.Run("Export keeps computed nested fields uncommented", func(t *testing.T) {
		// A nested block whose inner field is Optional+Computed at a default
		// value would normally trigger the `#` comment prefix — but Computed
		// must suppress it.
		schemer := &fakeSchemer{
			fakeMarshaler: fakeMarshaler{
				fn: func(p hcl.Properties) error {
					p["outer"] = []any{
						map[string]any{
							"inner": false,
						},
					}
					return nil
				},
			},
			schema: map[string]*schema.Schema{
				"outer": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"inner": {Type: schema.TypeBool, Optional: true, Computed: true},
						},
					},
				},
			},
		}
		buf := &bytes.Buffer{}

		err := ExportResource(schemer, buf, "t", "n")

		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "outer {")
		assert.Contains(t, out, "inner = false")
		assert.NotContains(t, out, "# inner")
	})
}

func TestResourceEntry(t *testing.T) {
	t.Run("IsOptional is always false", func(t *testing.T) {
		entry := &resourceEntry{Key: "block", Optional: true}
		assert.False(t, entry.IsOptional())
	})

	t.Run("IsComputed is always false", func(t *testing.T) {
		entry := &resourceEntry{Key: "block"}
		assert.False(t, entry.IsComputed())
	})

	t.Run("IsDefault is always false", func(t *testing.T) {
		entry := &resourceEntry{Key: "block"}
		assert.False(t, entry.IsDefault())
	})

	t.Run("IsLessThan", func(t *testing.T) {
		t.Run("is false when compared to a primitiveEntry", func(t *testing.T) {
			entry := &resourceEntry{Key: "alpha"}
			other := &primitiveEntry{Key: "zzz"}
			assert.False(t, entry.IsLessThan(other))
		})

		t.Run("is false when compared to a stringMapEntry", func(t *testing.T) {
			entry := &resourceEntry{Key: "alpha"}
			other := &stringMapEntry{Key: "zzz"}
			assert.False(t, entry.IsLessThan(other))
		})

		t.Run("orders resourceEntries alphabetically", func(t *testing.T) {
			alpha := &resourceEntry{Key: "alpha"}
			beta := &resourceEntry{Key: "beta"}

			assert.True(t, alpha.IsLessThan(beta))
			assert.False(t, beta.IsLessThan(alpha))
			assert.False(t, alpha.IsLessThan(alpha))
		})

		t.Run("returns false against an unknown entry type", func(t *testing.T) {
			entry := &resourceEntry{Key: "block"}
			other := &mapEntry{Key: "block"}
			assert.False(t, entry.IsLessThan(other))
		})
	})

	t.Run("Write", func(t *testing.T) {
		t.Run("emits a nested block with sub-entries", func(t *testing.T) {
			out := renderBlockBody(t, func(block *hclwrite.Block) {
				entry := &resourceEntry{
					Key: "settings",
					Entries: exportEntries{
						&primitiveEntry{Key: "name", Value: "child"},
					},
				}
				require.NoError(t, entry.Write(block.Body(), ""))
			})

			assert.Contains(t, out, "settings {")
			assert.Contains(t, out, `name = "child"`)
			assert.Contains(t, out, "}")
		})

		t.Run("sorts sub-entries before writing", func(t *testing.T) {
			// `name` has sort-key "00name" and should appear before any other key,
			// regardless of insertion order.
			out := renderBlockBody(t, func(block *hclwrite.Block) {
				entry := &resourceEntry{
					Key: "settings",
					Entries: exportEntries{
						&primitiveEntry{Key: "zzz", Value: "last"},
						&primitiveEntry{Key: "name", Value: "first"},
					},
				}
				require.NoError(t, entry.Write(block.Body(), ""))
			})

			// hclwrite aligns `=` across adjacent attributes, so don't assume
			// exact spacing — just verify the keys appear in the expected order
			// and that both values are rendered.
			nameIdx := strings.Index(out, `"first"`)
			zzzIdx := strings.Index(out, `"last"`)
			require.NotEqual(t, -1, nameIdx)
			require.NotEqual(t, -1, zzzIdx)
			assert.Less(t, nameIdx, zzzIdx)
		})

		t.Run("renders nested resourceEntries recursively", func(t *testing.T) {
			out := renderBlockBody(t, func(block *hclwrite.Block) {
				entry := &resourceEntry{
					Key: "outer",
					Entries: exportEntries{
						&resourceEntry{
							Key: "inner",
							Entries: exportEntries{
								&primitiveEntry{Key: "leaf", Value: "v"},
							},
						},
					},
				}
				require.NoError(t, entry.Write(block.Body(), ""))
			})

			assert.Contains(t, out, "outer {")
			assert.Contains(t, out, "inner {")
			assert.Contains(t, out, `leaf = "v"`)
		})

		t.Run("propagates Write errors from sub-entries", func(t *testing.T) {
			file := hclwrite.NewEmptyFile()
			block := file.Body().AppendNewBlock("resource", []string{"t", "n"})

			entry := &resourceEntry{
				Key: "settings",
				Entries: exportEntries{
					&primitiveEntry{Key: "broken", Value: struct{}{}},
				},
			}

			err := entry.Write(block.Body(), "")
			require.Error(t, err)
			assert.Contains(t, err.Error(), "broken")
		})
	})
}
