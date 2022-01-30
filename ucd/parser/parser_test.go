package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nihei9/ucdx/ucd/property"
)

func TestParser(t *testing.T) {
	tests := []struct {
		src           string
		fields        [][]field
		defaultFields [][]field
	}{
		{
			src: `

123

`,
			fields: [][]field{
				{"123"},
			},
			defaultFields: [][]field{
				nil,
			},
		},
		{
			src: `
# This is a comment
123 # This is a comment
# This is a comment
`,
			fields: [][]field{
				{"123"},
			},
			defaultFields: [][]field{
				nil,
			},
		},
		{
			src: `
123; foo; bar
`,
			fields: [][]field{
				{"123", "foo", "bar"},
			},
			defaultFields: [][]field{
				nil,
			},
		},
		{
			src: `
123; foo;
`,
			fields: [][]field{
				{"123", "foo", ""},
			},
			defaultFields: [][]field{
				nil,
			},
		},
		{
			src: `
123;; foo
`,
			fields: [][]field{
				{"123", "", "foo"},
			},
			defaultFields: [][]field{
				nil,
			},
		},
		{
			src: `
# @missing: 123; foo; bar
`,
			fields: [][]field{
				nil,
			},
			defaultFields: [][]field{
				{"123", "foo", "bar"},
			},
		},
		{
			src: `
123; foo # @missing: 456; bar
`,
			fields: [][]field{
				{"123", "foo"},
			},
			defaultFields: [][]field{
				{"456", "bar"},
			},
		},
		{
			src: `
123; foo
# @missing: 456; bar
`,
			fields: [][]field{
				{"123", "foo"},
				nil,
			},
			defaultFields: [][]field{
				nil,
				{"456", "bar"},
			},
		},
		{
			src: `
123; foo
456; bar
`,
			fields: [][]field{
				{"123", "foo"},
				{"456", "bar"},
			},
			defaultFields: [][]field{
				nil,
				nil,
			},
		},
		{
			src: `
# @missing: 123; foo
# @missing: 456; bar
`,
			fields: [][]field{
				nil,
				nil,
			},
			defaultFields: [][]field{
				{"123", "foo"},
				{"456", "bar"},
			},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("#%v", i), func(t *testing.T) {
			p := newParser(strings.NewReader(tt.src))
			n := 0
			for p.parse() {
				fields := tt.fields[n]
				if len(p.fields) != len(fields) {
					t.Fatalf("unexpected fields: want: %#v, got: %#v", fields, p.fields)
				}
				for i, f := range p.fields {
					if f.symbol() != fields[i].symbol() {
						t.Fatalf("unexpected fields: want: %#v, got: %#v", fields, p.fields)
					}
				}

				defaultFields := tt.defaultFields[n]
				if len(p.defaultFields) != len(defaultFields) {
					t.Fatalf("unexpected default fields: want: %#v, got: %#v", defaultFields, p.fields)
				}
				for i, f := range p.defaultFields {
					if f.symbol() != defaultFields[i].symbol() {
						t.Fatalf("unexpected default fields: want: %#v, got: %#v", defaultFields, p.defaultFields)
					}
				}

				n++
			}
			if p.err != nil {
				t.Fatal(p.err)
			}
		})
	}
}

func TestField_codePointRange(t *testing.T) {
	tests := []struct {
		field field
		cp    *property.CodePointRange
	}{
		{
			field: "0",
			cp:    property.NewCodePointRange(0x0, 0x0),
		},
		{
			field: "10FFFF",
			cp:    property.NewCodePointRange(0x10FFFF, 0x10FFFF),
		},
		{
			field: "0..10FFFF",
			cp:    property.NewCodePointRange(0x0, 0x10FFFF),
		},
	}
	for _, tt := range tests {
		t.Run(tt.field.String(), func(t *testing.T) {
			cp, err := tt.field.codePointRange()
			if err != nil {
				t.Fatal(err)
			}
			from, to := cp.Range()
			expectedFrom, expectedTo := tt.cp.Range()
			if from != expectedFrom || to != expectedTo {
				t.Fatalf("unexpected code point range: want: %v, got: %v", tt.cp, cp)
			}
		})
	}
}

func TestField_normalizedSymbol(t *testing.T) {
	tests := []struct {
		field field
		norm  string
	}{
		{
			field: "Foo",
			norm:  "foo",
		},
		{
			field: "Foo_Bar",
			norm:  "foobar",
		},
		{
			field: "Foo-Bar",
			norm:  "foobar",
		},
		{
			field: "Foo Bar",
			norm:  "foobar",
		},
		{
			field: "Foo_Bar-Baz Bra",
			norm:  "foobarbazbra",
		},
		{
			field: "isFoo",
			norm:  "foo",
		},
		{
			field: "is",
			norm:  "is",
		},
	}
	for _, tt := range tests {
		t.Run(tt.field.String(), func(t *testing.T) {
			norm := tt.field.normalizedSymbol()
			if norm != tt.norm {
				t.Fatalf("unexpected normalized symbol: want: %v, got: %v", tt.norm, norm)
			}
		})
	}
}
