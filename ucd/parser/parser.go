package parser

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"io"
	"regexp"
	"strings"

	"github.com/nihei9/ucdx/ucd/property"
)

var (
	reLine           = regexp.MustCompile(`^\s*(.*?)\s*(#.*)?$`)
	reCodePointRange = regexp.MustCompile(`^([[:xdigit:]]+)(?:..([[:xdigit:]]+))?$`)

	specialCommentPrefix = "# @missing:"

	symValReplacer = strings.NewReplacer("_", "", "-", "", "\x20", "")
)

// parser parses data files of UCD.
//
// However, for practical purposes, each field needs to be analyzed more specifically. For instance, in
// UnicodeData.txt, the first field is a code point range, so it needs to be recognized as a hexadecimal string.
// We can perform more specific parsing for each file by implementing a dedicated parser that wraps this parser.
//
// See section 4.2 File Format Conventions in [UAX44] for more information on the file format.
type parser struct {
	scanner       *bufio.Scanner
	fields        []field
	defaultFields []field
	err           error

	fieldBuf        []field
	defaultFieldBuf []field
}

func newParser(r io.Reader) *parser {
	return &parser{
		scanner:         bufio.NewScanner(r),
		fieldBuf:        make([]field, 50),
		defaultFieldBuf: make([]field, 50),
	}
}

func (p *parser) parse() bool {
	for p.scanner.Scan() {
		p.parseRecord(p.scanner.Text())
		if p.fields != nil || p.defaultFields != nil {
			return true
		}
	}
	p.err = p.scanner.Err()
	return false
}

func (p *parser) parseRecord(src string) {
	ms := reLine.FindStringSubmatch(src)
	mFields := ms[1]
	mComment := ms[2]
	if mFields != "" {
		p.fields = parseFields(p.fieldBuf, mFields)
	} else {
		p.fields = nil
	}
	if strings.HasPrefix(mComment, specialCommentPrefix) {
		p.defaultFields = parseFields(p.defaultFieldBuf, strings.Replace(mComment, specialCommentPrefix, "", -1))
	} else {
		p.defaultFields = nil
	}
}

func parseFields(buf []field, src string) []field {
	n := 0
	for _, f := range strings.Split(src, ";") {
		buf[n] = field(strings.TrimSpace(f))
		n++
	}

	return buf[:n]
}

type field string

func (f field) String() string {
	return string(f)
}

// codePointRange returns a code point range.
func (f field) codePointRange() (*property.CodePointRange, error) {
	var from, to rune
	var err error
	cp := reCodePointRange.FindStringSubmatch(string(f))
	from, err = decodeHexToRune(cp[1])
	if err != nil {
		return nil, err
	}
	if cp[2] != "" {
		to, err = decodeHexToRune(cp[2])
		if err != nil {
			return nil, err
		}
	} else {
		to = from
	}
	return property.NewCodePointRange(from, to), nil
}

func decodeHexToRune(hexCodePoint string) (rune, error) {
	h := hexCodePoint
	if len(h)%2 != 0 {
		h = "0" + h
	}
	b, err := hex.DecodeString(h)
	if err != nil {
		return 0, err
	}
	l := len(b)
	for i := 0; i < 4-l; i++ {
		b = append([]byte{0}, b...)
	}
	n := binary.BigEndian.Uint32(b)
	return rune(n), nil
}

// rangeStart returns true when a field indicates a start of a code point range.
//
// See section 4.2.3 Code Point Ranges in [UAX44] for more details on the code point ranges.
func (f field) rangeStart() bool {
	s := string(f)
	if strings.HasPrefix(s, "<") && strings.HasSuffix(s, ", First>") {
		return true
	}
	return false
}

// rangeLast returns true when a field indicates a last of a code point range.
//
// See section 4.2.3 Code Point Ranges in [UAX44] for more details on the code point ranges.
func (f field) rangeLast() bool {
	s := string(f)
	if strings.HasPrefix(s, "<") && strings.HasSuffix(s, ", Last>") {
		return true
	}
	return false
}

// name returns a value parsed as a `Name` property.
//
// See section 4.8 Name in [Unicode].
// > Interpretation of Field 1 of UnicodeData.txt. Where Field 1 of UnicodeData.txt contains a string enclosed in
// > angle brackets, “<” and “>”, such a string is not a character name, ...
func (f field) name() (property.PropertyName, bool) {
	s := string(f)
	if strings.HasPrefix(s, "<") && strings.HasSuffix(s, ">") {
		return "", false
	}
	return property.NewPropertyName(s), true
}

// symbol returns a symbolic value.
func (f field) symbol() property.PropertyValueSymbol {
	return property.NewSymbolPropertyValue(string(f))
}

// normalizeSymbolicValue returns a normalized symbolic value.
//
// The normalization algorithm follows UAX44-LM3 defined section 5.9.3 Matching Symbolic Values in [UAX44].
func (f field) normalizedSymbol() property.PropertyValueSymbol {
	sym := strings.ToLower(symValReplacer.Replace(string(f)))
	if sym == "is" {
		return property.NewSymbolPropertyValue(sym)
	}
	return property.NewSymbolPropertyValue(strings.TrimPrefix(sym, "is"))
}
