// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

// By default, this json support uses base64 encoding for bytes, because you cannot
// store and read any arbitrary string in json (only unicode).
// However, the user can configre how to encode/decode bytes.
//
// This library specifically supports UTF-8 for encoding and decoding only.
//
// Note that the library will happily encode/decode things which are not valid
// json e.g. a map[int64]string. We do it for consistency. With valid json,
// we will encode and decode appropriately.
// Users can specify their map type if necessary to force it.
//
// Note:
//   - we cannot use strconv.Quote and strconv.Unquote because json quotes/unquotes differently.
//     We implement it here.

// Top-level methods of json(End|Dec)Driver (which are implementations of (en|de)cDriver
// MUST not call one-another.

import (
	"bytes"
	"encoding/base64"
	"math"
	"reflect"
	"strconv"
	"time"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

//--------------------------------

var jsonLiterals = [...]byte{
	'"', 't', 'r', 'u', 'e', '"',
	'"', 'f', 'a', 'l', 's', 'e', '"',
	'"', 'n', 'u', 'l', 'l', '"',
}

const (
	jsonLitTrueQ  = 0
	jsonLitTrue   = 1
	jsonLitFalseQ = 6
	jsonLitFalse  = 7
	// jsonLitNullQ  = 13
	jsonLitNull = 14
)

var (
	jsonLiteral4True  = jsonLiterals[jsonLitTrue+1 : jsonLitTrue+4]
	jsonLiteral4False = jsonLiterals[jsonLitFalse+1 : jsonLitFalse+5]
	jsonLiteral4Null  = jsonLiterals[jsonLitNull+1 : jsonLitNull+4]
)

const (
	jsonU4Chk2 = '0'
	jsonU4Chk1 = 'a' - 10
	jsonU4Chk0 = 'A' - 10

	jsonScratchArrayLen = cacheLineSize // 64
)

const (
	// If !jsonValidateSymbols, decoding will be faster, by skipping some checks:
	//   - If we see first character of null, false or true,
	//     do not validate subsequent characters.
	//   - e.g. if we see a n, assume null and skip next 3 characters,
	//     and do not validate they are ull.
	// P.S. Do not expect a significant decoding boost from this.
	jsonValidateSymbols = true

	jsonSpacesOrTabsLen = 128

	jsonAlwaysReturnInternString = false
)

var (
	// jsonTabs and jsonSpaces are used as caches for indents
	jsonTabs, jsonSpaces [jsonSpacesOrTabsLen]byte

	jsonCharHtmlSafeSet   bitset256
	jsonCharSafeSet       bitset256
	jsonCharWhitespaceSet bitset256
	jsonNumSet            bitset256
)

func init() {
	var i byte
	for i = 0; i < jsonSpacesOrTabsLen; i++ {
		jsonSpaces[i] = ' '
		jsonTabs[i] = '\t'
	}

	// populate the safe values as true: note: ASCII control characters are (0-31)
	// jsonCharSafeSet:     all true except (0-31) " \
	// jsonCharHtmlSafeSet: all true except (0-31) " \ < > &
	for i = 32; i < utf8.RuneSelf; i++ {
		switch i {
		case '"', '\\':
		case '<', '>', '&':
			jsonCharSafeSet.set(i) // = true
		default:
			jsonCharSafeSet.set(i)
			jsonCharHtmlSafeSet.set(i)
		}
	}
	for i = 0; i <= utf8.RuneSelf; i++ {
		switch i {
		case ' ', '\t', '\r', '\n':
			jsonCharWhitespaceSet.set(i)
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'e', 'E', '.', '+', '-':
			jsonNumSet.set(i)
		}
	}
}

// ----------------

type jsonEncDriverTypical struct {
	jsonEncDriver
}

func (e *jsonEncDriverTypical) WriteArrayStart(length int) {
	e.w.writen1('[')
}

func (e *jsonEncDriverTypical) WriteArrayElem() {
	if e.e.c != containerArrayStart {
		e.w.writen1(',')
	}
}

func (e *jsonEncDriverTypical) WriteArrayEnd() {
	e.w.writen1(']')
}

func (e *jsonEncDriverTypical) WriteMapStart(length int) {
	e.w.writen1('{')
}

func (e *jsonEncDriverTypical) WriteMapElemKey() {
	if e.e.c != containerMapStart {
		e.w.writen1(',')
	}
}

func (e *jsonEncDriverTypical) WriteMapElemValue() {
	e.w.writen1(':')
}

func (e *jsonEncDriverTypical) WriteMapEnd() {
	e.w.writen1('}')
}

func (e *jsonEncDriverTypical) EncodeBool(b bool) {
	if b {
		e.w.writeb(jsonLiterals[jsonLitTrue : jsonLitTrue+4])
	} else {
		e.w.writeb(jsonLiterals[jsonLitFalse : jsonLitFalse+5])
	}
}

func (e *jsonEncDriverTypical) EncodeInt(v int64) {
	e.w.writeb(strconv.AppendInt(e.b[:0], v, 10))
}

func (e *jsonEncDriverTypical) EncodeUint(v uint64) {
	e.w.writeb(strconv.AppendUint(e.b[:0], v, 10))
}

func (e *jsonEncDriverTypical) EncodeFloat64(f float64) {
	fmt, prec := jsonFloatStrconvFmtPrec64(f)
	e.w.writeb(strconv.AppendFloat(e.b[:0], f, fmt, int(prec), 64))
	// e.w.writeb(strconv.AppendFloat(e.b[:0], f, jsonFloatStrconvFmtPrec64(f), 64))
}

func (e *jsonEncDriverTypical) EncodeFloat32(f float32) {
	fmt, prec := jsonFloatStrconvFmtPrec32(f)
	e.w.writeb(strconv.AppendFloat(e.b[:0], float64(f), fmt, int(prec), 32))
}

// func (e *jsonEncDriverTypical) encodeFloat(f float64, bitsize uint8) {
// 	fmt, prec := jsonFloatStrconvFmtPrec(f, bitsize == 32)
// 	e.w.writeb(strconv.AppendFloat(e.b[:0], f, fmt, prec, int(bitsize)))
// }

// func (e *jsonEncDriverTypical) atEndOfEncode() {
// 	if e.tw {
// 		e.w.writen1(' ')
// 	}
// }

// ----------------

type jsonEncDriverGeneric struct {
	jsonEncDriver
}

// indent is done as below:
//   - newline and indent are added before each mapKey or arrayElem
//   - newline and indent are added before each ending,
//     except there was no entry (so we can have {} or [])

func (e *jsonEncDriverGeneric) reset() {
	e.jsonEncDriver.reset()
	e.d, e.dl, e.di = false, 0, 0
	if e.h.Indent != 0 {
		e.d = true
		e.di = int8(e.h.Indent)
	}
	// if e.h.Indent > 0 {
	// 	e.d = true
	// 	e.di = int8(e.h.Indent)
	// } else if e.h.Indent < 0 {
	// 	e.d = true
	// 	// e.dt = true
	// 	e.di = int8(-e.h.Indent)
	// }
	e.ks = e.h.MapKeyAsString
	e.is = e.h.IntegerAsString
}

func (e *jsonEncDriverGeneric) WriteArrayStart(length int) {
	if e.d {
		e.dl++
	}
	e.w.writen1('[')
}

func (e *jsonEncDriverGeneric) WriteArrayEnd() {
	if e.d {
		e.dl--
		e.writeIndent()
	}
	e.w.writen1(']')
}

func (e *jsonEncDriverGeneric) WriteMapStart(length int) {
	if e.d {
		e.dl++
	}
	e.w.writen1('{')
}

func (e *jsonEncDriverGeneric) WriteMapEnd() {
	if e.d {
		e.dl--
		if e.e.c != containerMapStart {
			e.writeIndent()
		}
	}
	e.w.writen1('}')
}

func (e *jsonEncDriverGeneric) EncodeBool(b bool) {
	if e.ks && e.e.c == containerMapKey {
		if b {
			e.w.writeb(jsonLiterals[jsonLitTrueQ : jsonLitTrueQ+6])
		} else {
			e.w.writeb(jsonLiterals[jsonLitFalseQ : jsonLitFalseQ+7])
		}
	} else {
		if b {
			e.w.writeb(jsonLiterals[jsonLitTrue : jsonLitTrue+4])
		} else {
			e.w.writeb(jsonLiterals[jsonLitFalse : jsonLitFalse+5])
		}
	}
}

func (e *jsonEncDriverGeneric) encodeFloat(f float64, bitsize, fmt byte, prec int8) {
	var blen int
	if e.ks && e.e.c == containerMapKey {
		blen = 2 + len(strconv.AppendFloat(e.b[1:1], f, fmt, int(prec), int(bitsize)))
		e.b[0] = '"'
		e.b[blen-1] = '"'
	} else {
		blen = len(strconv.AppendFloat(e.b[:0], f, fmt, int(prec), int(bitsize)))
	}
	e.w.writeb(e.b[:blen])
}

func (e *jsonEncDriverGeneric) EncodeFloat64(f float64) {
	fmt, prec := jsonFloatStrconvFmtPrec64(f)
	e.encodeFloat(f, 64, fmt, prec)
}

func (e *jsonEncDriverGeneric) EncodeFloat32(f float32) {
	fmt, prec := jsonFloatStrconvFmtPrec32(f)
	e.encodeFloat(float64(f), 32, fmt, prec)
}

func (e *jsonEncDriverGeneric) EncodeInt(v int64) {
	x := e.is
	if x == 'A' || x == 'L' && (v > 1<<53 || v < -(1<<53)) || (e.ks && e.e.c == containerMapKey) {
		blen := 2 + len(strconv.AppendInt(e.b[1:1], v, 10))
		e.b[0] = '"'
		e.b[blen-1] = '"'
		e.w.writeb(e.b[:blen])
		return
	}
	e.w.writeb(strconv.AppendInt(e.b[:0], v, 10))
}

func (e *jsonEncDriverGeneric) EncodeUint(v uint64) {
	x := e.is
	if x == 'A' || x == 'L' && v > 1<<53 || (e.ks && e.e.c == containerMapKey) {
		blen := 2 + len(strconv.AppendUint(e.b[1:1], v, 10))
		e.b[0] = '"'
		e.b[blen-1] = '"'
		e.w.writeb(e.b[:blen])
		return
	}
	e.w.writeb(strconv.AppendUint(e.b[:0], v, 10))
}

// func (e *jsonEncDriverGeneric) EncodeFloat32(f float32) {
// 	// e.encodeFloat(float64(f), 32)
// 	// always encode all floats as IEEE 64-bit floating point.
// 	// It also ensures that we can decode in full precision even if into a float32,
// 	// as what is written is always to float64 precision.
// 	e.EncodeFloat64(float64(f))
// }

// func (e *jsonEncDriverGeneric) atEndOfEncode() {
// 	if e.tw {
// 		if e.d {
// 			e.w.writen1('\n')
// 		} else {
// 			e.w.writen1(' ')
// 		}
// 	}
// }

// --------------------

type jsonEncDriver struct {
	noBuiltInTypes
	w *encWriterSwitch
	e *Encoder
	h *JsonHandle

	bs []byte // for encoding strings
	se interfaceExtWrapper

	// ---- cpu cache line boundary?
	// ds string // indent string
	di int8 // indent per: if negative, use tabs
	d  bool // indenting?
	// dt bool   // indent using tabs
	dl uint16 // indent level
	ks bool   // map key as string
	is byte   // integer as string

	s *bitset256 // safe set for characters (taking h.HTMLAsIs into consideration)
	// scratch: encode time, numbers, etc. Note: leave 1 byte for containerState
	b [jsonScratchArrayLen - 16]byte // leave space for bs(len,cap), containerState
}

// Keep writeIndent, WriteArrayElem, WriteMapElemKey, WriteMapElemValue
// in jsonEncDriver, so that *Encoder can directly call them

func (e *jsonEncDriver) getJsonEncDriver() *jsonEncDriver { return e }

func (e *jsonEncDriver) writeIndent() {
	e.w.writen1('\n')
	x := int(e.di) * int(e.dl)
	if e.di < 0 {
		x = -x
		for x > jsonSpacesOrTabsLen {
			e.w.writeb(jsonTabs[:])
			x -= jsonSpacesOrTabsLen
		}
		e.w.writeb(jsonTabs[:x])
	} else {
		for x > jsonSpacesOrTabsLen {
			e.w.writeb(jsonSpaces[:])
			x -= jsonSpacesOrTabsLen
		}
		e.w.writeb(jsonSpaces[:x])
	}
}

func (e *jsonEncDriver) WriteArrayElem() {
	// xdebugf("WriteArrayElem: e.e.c: %d", e.e.c)
	if e.e.c != containerArrayStart {
		e.w.writen1(',')
	}
	if e.d {
		e.writeIndent()
	}
}

func (e *jsonEncDriver) WriteMapElemKey() {
	if e.e.c != containerMapStart {
		e.w.writen1(',')
	}
	if e.d {
		e.writeIndent()
	}
}

func (e *jsonEncDriver) WriteMapElemValue() {
	if e.d {
		e.w.writen2(':', ' ')
	} else {
		e.w.writen1(':')
	}
}

func (e *jsonEncDriver) EncodeNil() {
	// We always encode nil as just null (never in quotes)
	// This allows us to easily decode if a nil in the json stream
	// ie if initial token is n.
	e.w.writeb(jsonLiterals[jsonLitNull : jsonLitNull+4])

	// if e.h.MapKeyAsString && e.e.c == containerMapKey {
	// 	e.w.writeb(jsonLiterals[jsonLitNullQ : jsonLitNullQ+6])
	// } else {
	// 	e.w.writeb(jsonLiterals[jsonLitNull : jsonLitNull+4])
	// }
}

func (e *jsonEncDriver) EncodeTime(t time.Time) {
	// Do NOT use MarshalJSON, as it allocates internally.
	// instead, we call AppendFormat directly, using our scratch buffer (e.b)
	if t.IsZero() {
		e.EncodeNil()
	} else {
		e.b[0] = '"'
		b := t.AppendFormat(e.b[1:1], time.RFC3339Nano)
		e.b[len(b)+1] = '"'
		e.w.writeb(e.b[:len(b)+2])
	}
	// v, err := t.MarshalJSON(); if err != nil { e.e.error(err) } e.w.writeb(v)
}

func (e *jsonEncDriver) EncodeExt(rv interface{}, xtag uint64, ext Ext, en *Encoder) {
	if v := ext.ConvertExt(rv); v == nil {
		e.EncodeNil()
	} else {
		en.encode(v)
	}
}

func (e *jsonEncDriver) EncodeRawExt(re *RawExt, en *Encoder) {
	// only encodes re.Value (never re.Data)
	if re.Value == nil {
		e.EncodeNil()
	} else {
		en.encode(re.Value)
	}
}

func (e *jsonEncDriver) EncodeStringEnc(c charEncoding, v string) {
	e.quoteStr(v)
}

func (e *jsonEncDriver) EncodeStringBytesRaw(v []byte) {
	// if encoding raw bytes and RawBytesExt is configured, use it to encode
	if v == nil {
		e.EncodeNil()
		return
	}
	if e.se.InterfaceExt != nil {
		e.EncodeExt(v, 0, &e.se, e.e)
		return
	}

	slen := base64.StdEncoding.EncodedLen(len(v)) + 2
	if cap(e.bs) >= slen {
		e.bs = e.bs[:slen]
	} else {
		e.bs = make([]byte, slen)
	}
	e.bs[0] = '"'
	base64.StdEncoding.Encode(e.bs[1:], v)
	e.bs[slen-1] = '"'
	e.w.writeb(e.bs)
}

func (e *jsonEncDriver) EncodeAsis(v []byte) {
	e.w.writeb(v)
}

func (e *jsonEncDriver) quoteStr(s string) {
	// adapted from std pkg encoding/json
	const hex = "0123456789abcdef"
	w := e.w
	w.writen1('"')
	var start int
	for i := 0; i < len(s); {
		// encode all bytes < 0x20 (except \r, \n).
		// also encode < > & to prevent security holes when served to some browsers.

		// We optimize for ascii, by assumining that most characters are in the BMP
		// and natively consumed by json without much computation.

		// if 0x20 <= b && b != '\\' && b != '"' && b != '<' && b != '>' && b != '&' {
		// if (htmlasis && jsonCharSafeSet.isset(b)) || jsonCharHtmlSafeSet.isset(b) {
		if e.s.isset(s[i]) {
			i++
			continue
		}
		if b := s[i]; b < utf8.RuneSelf {
			if start < i {
				w.writestr(s[start:i])
			}
			switch b {
			case '\\', '"':
				w.writen2('\\', b)
			case '\n':
				w.writen2('\\', 'n')
			case '\r':
				w.writen2('\\', 'r')
			case '\b':
				w.writen2('\\', 'b')
			case '\f':
				w.writen2('\\', 'f')
			case '\t':
				w.writen2('\\', 't')
			default:
				w.writestr(`\u00`)
				w.writen2(hex[b>>4], hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError {
			if size == 1 {
				if start < i {
					w.writestr(s[start:i])
				}
				w.writestr(`\ufffd`)
				i++
				start = i
			}
			continue
		}
		// U+2028 is LINE SEPARATOR. U+2029 is PARAGRAPH SEPARATOR.
		// Both technically valid JSON, but bomb on JSONP, so fix here unconditionally.
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				w.writestr(s[start:i])
			}
			w.writestr(`\u202`)
			w.writen1(hex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		w.writestr(s[start:])
	}
	w.writen1('"')
}

func (e *jsonEncDriver) atEndOfEncode() {
	// if e.e.c == 0 { // scalar written, output space
	// 	e.w.writen1(' ')
	// } else if e.h.TermWhitespace { // container written, output new-line
	// 	e.w.writen1('\n')
	// }
	if e.h.TermWhitespace {
		if e.e.c == 0 { // scalar written, output space
			e.w.writen1(' ')
		} else { // container written, output new-line
			e.w.writen1('\n')
		}
	}

}

type jsonDecDriver struct {
	noBuiltInTypes
	d  *Decoder
	h  *JsonHandle
	r  *decReaderSwitch
	se interfaceExtWrapper

	bs []byte // scratch, initialized from b. For parsing strings or numbers.
	// ---- writable fields during execution --- *try* to keep in sep cache line

	// ---- cpu cache line boundary?
	b [jsonScratchArrayLen]byte // scratch 1, used for parsing strings or numbers or time.Time
	// ---- cpu cache line boundary?
	// c     containerState
	tok   uint8                         // used to store the token read right after skipWhiteSpace
	fnull bool                          // found null from appendStringAsBytes
	_     [2]byte                       // padding
	bstr  [4]byte                       // scratch used for string \UXXX parsing
	b2    [jsonScratchArrayLen - 8]byte // scratch 2, used only for readUntil, decNumBytes

	// _ [3]uint64 // padding
	// n jsonNum
}

// func jsonIsWS(b byte) bool {
// 	// return b == ' ' || b == '\t' || b == '\r' || b == '\n'
// 	return jsonCharWhitespaceSet.isset(b)
// }

func (d *jsonDecDriver) uncacheRead() {
	if d.tok != 0 {
		d.r.unreadn1()
		d.tok = 0
	}
}

func (d *jsonDecDriver) ReadMapStart() int {
	if d.tok == 0 {
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	}
	const xc uint8 = '{'
	if d.tok != xc {
		d.d.errorf("read map - expect char '%c' but got char '%c'", xc, d.tok)
	}
	d.tok = 0
	return -1
}

func (d *jsonDecDriver) ReadArrayStart() int {
	if d.tok == 0 {
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	}
	const xc uint8 = '['
	if d.tok != xc {
		d.d.errorf("read array - expect char '%c' but got char '%c'", xc, d.tok)
	}
	d.tok = 0
	return -1
}

func (d *jsonDecDriver) CheckBreak() bool {
	if d.tok == 0 {
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	}
	return d.tok == '}' || d.tok == ']'
}

// For the ReadXXX methods below, we could just delegate to helper functions
// readContainerState(c containerState, xc uint8, check bool)
// - ReadArrayElem would become:
//   readContainerState(containerArrayElem, ',', d.d.c != containerArrayStart)
//
// However, until mid-stack inlining comes in go1.11 which supports inlining of
// one-liners, we explicitly write them all 5 out to elide the extra func call.
//
// TODO: For Go 1.11, if inlined, consider consolidating these.

func (d *jsonDecDriver) ReadArrayElem() {
	const xc uint8 = ','
	if d.tok == 0 {
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	}
	// xdebugf("ReadArrayElem: d.d.c: %d, token: %c", d.d.c, d.tok)
	if d.d.c != containerArrayStart {
		if d.tok != xc {
			d.d.errorf("read array element - expect char '%c' but got char '%c'", xc, d.tok)
		}
		d.tok = 0
	}
}

func (d *jsonDecDriver) ReadArrayEnd() {
	const xc uint8 = ']'
	if d.tok == 0 {
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	}
	if d.tok != xc {
		d.d.errorf("read array end - expect char '%c' but got char '%c'", xc, d.tok)
	}
	d.tok = 0
}

func (d *jsonDecDriver) ReadMapElemKey() {
	const xc uint8 = ','
	if d.tok == 0 {
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	}
	if d.d.c != containerMapStart {
		if d.tok != xc {
			d.d.errorf("read map key - expect char '%c' but got char '%c'", xc, d.tok)
		}
		d.tok = 0
	}
}

func (d *jsonDecDriver) ReadMapElemValue() {
	const xc uint8 = ':'
	if d.tok == 0 {
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	}
	if d.tok != xc {
		d.d.errorf("read map value - expect char '%c' but got char '%c'", xc, d.tok)
	}
	d.tok = 0
}

func (d *jsonDecDriver) ReadMapEnd() {
	const xc uint8 = '}'
	if d.tok == 0 {
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	}
	if d.tok != xc {
		d.d.errorf("read map end - expect char '%c' but got char '%c'", xc, d.tok)
	}
	d.tok = 0
}

// func (d *jsonDecDriver) readLit(length, fromIdx uint8) {
// 	// length here is always less than 8 (literals are: null, true, false)
// 	bs := d.r.readx(int(length))
// 	d.tok = 0
// 	if jsonValidateSymbols && !bytes.Equal(bs, jsonLiterals[fromIdx:fromIdx+length]) {
// 		d.d.errorf("expecting %s: got %s", jsonLiterals[fromIdx:fromIdx+length], bs)
// 	}
// }

func (d *jsonDecDriver) readLit4True() {
	bs := d.r.readx(3)
	d.tok = 0
	if jsonValidateSymbols && !bytes.Equal(bs, jsonLiteral4True) {
		d.d.errorf("expecting %s: got %s", jsonLiteral4True, bs)
	}
}

func (d *jsonDecDriver) readLit4False() {
	bs := d.r.readx(4)
	d.tok = 0
	if jsonValidateSymbols && !bytes.Equal(bs, jsonLiteral4False) {
		d.d.errorf("expecting %s: got %s", jsonLiteral4False, bs)
	}
}

func (d *jsonDecDriver) readLit4Null() {
	bs := d.r.readx(3)
	d.tok = 0
	if jsonValidateSymbols && !bytes.Equal(bs, jsonLiteral4Null) {
		d.d.errorf("expecting %s: got %s", jsonLiteral4Null, bs)
	}
}

func (d *jsonDecDriver) TryDecodeAsNil() bool {
	if d.tok == 0 {
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	}
	// we shouldn't try to see if "null" was here, right?
	// only the plain string: `null` denotes a nil (ie not quotes)
	if d.tok == 'n' {
		d.readLit4Null()
		return true
	}
	return false
}

func (d *jsonDecDriver) DecodeBool() (v bool) {
	if d.tok == 0 {
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	}
	fquot := d.d.c == containerMapKey && d.tok == '"'
	if fquot {
		d.tok = d.r.readn1()
	}
	switch d.tok {
	case 'f':
		d.readLit4False()
		// v = false
	case 't':
		d.readLit4True()
		v = true
	default:
		d.d.errorf("decode bool: got first char %c", d.tok)
		// v = false // "unreachable"
	}
	if fquot {
		d.r.readn1()
	}
	return
}

func (d *jsonDecDriver) DecodeTime() (t time.Time) {
	// read string, and pass the string into json.unmarshal
	d.appendStringAsBytes()
	if d.fnull {
		return
	}
	t, err := time.Parse(time.RFC3339, stringView(d.bs))
	if err != nil {
		d.d.errorv(err)
	}
	return
}

func (d *jsonDecDriver) ContainerType() (vt valueType) {
	// check container type by checking the first char
	if d.tok == 0 {
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	}

	// optimize this, so we don't do 4 checks but do one computation.
	// return jsonContainerSet[d.tok]

	// ContainerType is mostly called for Map and Array,
	// so this conditional is good enough (max 2 checks typically)
	if b := d.tok; b == '{' {
		return valueTypeMap
	} else if b == '[' {
		return valueTypeArray
	} else if b == 'n' {
		return valueTypeNil
	} else if b == '"' {
		return valueTypeString
	}
	return valueTypeUnset
}

func (d *jsonDecDriver) decNumBytes() (bs []byte) {
	// stores num bytes in d.bs
	if d.tok == 0 {
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	}
	if d.tok == '"' {
		bs = d.r.readUntil(d.b2[:0], '"')
		bs = bs[:len(bs)-1]
	} else {
		d.r.unreadn1()
		bs = d.r.readTo(d.bs[:0], &jsonNumSet)
	}
	d.tok = 0
	return bs
}

func (d *jsonDecDriver) DecodeUint64() (u uint64) {
	bs := d.decNumBytes()
	if len(bs) == 0 {
		return
	}
	n, neg, badsyntax, overflow := jsonParseInteger(bs)
	if overflow {
		d.d.errorf("overflow parsing unsigned integer: %s", bs)
	} else if neg {
		d.d.errorf("minus found parsing unsigned integer: %s", bs)
	} else if badsyntax {
		// fallback: try to decode as float, and cast
		n = d.decUint64ViaFloat(stringView(bs))
	}
	return n
}

func (d *jsonDecDriver) DecodeInt64() (i int64) {
	const cutoff = uint64(1 << uint(64-1))
	bs := d.decNumBytes()
	if len(bs) == 0 {
		return
	}
	n, neg, badsyntax, overflow := jsonParseInteger(bs)
	if overflow {
		d.d.errorf("overflow parsing integer: %s", bs)
	} else if badsyntax {
		// d.d.errorf("invalid syntax for integer: %s", bs)
		// fallback: try to decode as float, and cast
		if neg {
			n = d.decUint64ViaFloat(stringView(bs[1:]))
		} else {
			n = d.decUint64ViaFloat(stringView(bs))
		}
	}
	if neg {
		if n > cutoff {
			d.d.errorf("overflow parsing integer: %s", bs)
		}
		i = -(int64(n))
	} else {
		if n >= cutoff {
			d.d.errorf("overflow parsing integer: %s", bs)
		}
		i = int64(n)
	}
	return
}

func (d *jsonDecDriver) decUint64ViaFloat(s string) (u uint64) {
	if len(s) == 0 {
		return
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		d.d.errorf("invalid syntax for integer: %s", s)
		// d.d.errorv(err)
	}
	fi, ff := math.Modf(f)
	if ff > 0 {
		d.d.errorf("fractional part found parsing integer: %s", s)
	} else if fi > float64(math.MaxUint64) {
		d.d.errorf("overflow parsing integer: %s", s)
	}
	return uint64(fi)
}

func (d *jsonDecDriver) decodeFloat(bitsize int) (f float64) {
	bs := d.decNumBytes()
	if len(bs) == 0 {
		return
	}
	f, err := strconv.ParseFloat(stringView(bs), bitsize)
	if err != nil {
		d.d.errorv(err)
	}
	return
}

func (d *jsonDecDriver) DecodeFloat64() (f float64) {
	return d.decodeFloat(64)
}

func (d *jsonDecDriver) DecodeFloat32() (f float64) {
	return d.decodeFloat(32)
}

func (d *jsonDecDriver) DecodeExt(rv interface{}, xtag uint64, ext Ext) (realxtag uint64) {
	if ext == nil {
		re := rv.(*RawExt)
		re.Tag = xtag
		d.d.decode(&re.Value)
	} else {
		var v interface{}
		d.d.decode(&v)
		ext.UpdateExt(rv, v)
	}
	return
}

func (d *jsonDecDriver) DecodeBytes(bs []byte, zerocopy bool) (bsOut []byte) {
	// if decoding into raw bytes, and the RawBytesExt is configured, use it to decode.
	if d.se.InterfaceExt != nil {
		bsOut = bs
		d.DecodeExt(&bsOut, 0, &d.se)
		return
	}
	if d.tok == 0 {
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	}
	// check if an "array" of uint8's (see ContainerType for how to infer if an array)
	if d.tok == '[' {
		bsOut, _ = fastpathTV.DecSliceUint8V(bs, true, d.d)
		return
	}
	d.appendStringAsBytes()
	// base64 encodes []byte{} as "", and we encode nil []byte as null.
	// Consequently, base64 should decode null as a nil []byte, and "" as an empty []byte{}.
	// appendStringAsBytes returns a zero-len slice for both, so as not to reset d.bs.
	// However, it sets a fnull field to true, so we can check if a null was found.
	if len(d.bs) == 0 {
		if d.fnull {
			return nil
		}
		return []byte{}
	}
	bs0 := d.bs
	slen := base64.StdEncoding.DecodedLen(len(bs0))
	if slen <= cap(bs) {
		bsOut = bs[:slen]
	} else if zerocopy && slen <= cap(d.b2) {
		bsOut = d.b2[:slen]
	} else {
		bsOut = make([]byte, slen)
	}
	slen2, err := base64.StdEncoding.Decode(bsOut, bs0)
	if err != nil {
		d.d.errorf("error decoding base64 binary '%s': %v", bs0, err)
		return nil
	}
	if slen != slen2 {
		bsOut = bsOut[:slen2]
	}
	return
}

func (d *jsonDecDriver) DecodeString() (s string) {
	d.appendStringAsBytes()
	return d.bsToString()
}

func (d *jsonDecDriver) DecodeStringAsBytes() (s []byte) {
	d.appendStringAsBytes()
	return d.bs
}

func (d *jsonDecDriver) appendStringAsBytes() {
	if d.tok == 0 {
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	}

	d.fnull = false
	if d.tok != '"' {
		// d.d.errorf("expect char '%c' but got char '%c'", '"', d.tok)
		// handle non-string scalar: null, true, false or a number
		switch d.tok {
		case 'n':
			d.readLit4Null()
			d.bs = d.bs[:0]
			d.fnull = true
		case 'f':
			d.readLit4False()
			d.bs = d.bs[:5]
			copy(d.bs, "false")
		case 't':
			d.readLit4True()
			d.bs = d.bs[:4]
			copy(d.bs, "true")
		default:
			// try to parse a valid number
			bs := d.decNumBytes()
			if len(bs) <= cap(d.bs) {
				d.bs = d.bs[:len(bs)]
			} else {
				d.bs = make([]byte, len(bs))
			}
			copy(d.bs, bs)
		}
		return
	}

	d.tok = 0
	r := d.r
	var cs = r.readUntil(d.b2[:0], '"')
	var cslen = uint(len(cs))
	var c uint8
	v := d.bs[:0]
	// append on each byte seen can be expensive, so we just
	// keep track of where we last read a contiguous set of
	// non-special bytes (using cursor variable),
	// and when we see a special byte
	// e.g. end-of-slice, " or \,
	// we will append the full range into the v slice before proceeding
	var i, cursor uint
	for {
		if i == cslen {
			v = append(v, cs[cursor:]...)
			cs = r.readUntil(d.b2[:0], '"')
			cslen = uint(len(cs))
			i, cursor = 0, 0
		}
		c = cs[i]
		if c == '"' {
			v = append(v, cs[cursor:i]...)
			break
		}
		if c != '\\' {
			i++
			continue
		}
		v = append(v, cs[cursor:i]...)
		i++
		c = cs[i]
		switch c {
		case '"', '\\', '/', '\'':
			v = append(v, c)
		case 'b':
			v = append(v, '\b')
		case 'f':
			v = append(v, '\f')
		case 'n':
			v = append(v, '\n')
		case 'r':
			v = append(v, '\r')
		case 't':
			v = append(v, '\t')
		case 'u':
			var r rune
			var rr uint32
			if cslen < i+4 {
				d.d.errorf("need at least 4 more bytes for unicode sequence")
			}
			var j uint
			for _, c = range cs[i+1 : i+5] { // bounds-check-elimination
				// best to use explicit if-else
				// - not a table, etc which involve memory loads, array lookup with bounds checks, etc
				if c >= '0' && c <= '9' {
					rr = rr*16 + uint32(c-jsonU4Chk2)
				} else if c >= 'a' && c <= 'f' {
					rr = rr*16 + uint32(c-jsonU4Chk1)
				} else if c >= 'A' && c <= 'F' {
					rr = rr*16 + uint32(c-jsonU4Chk0)
				} else {
					r = unicode.ReplacementChar
					i += 4
					goto encode_rune
				}
			}
			r = rune(rr)
			i += 4
			if utf16.IsSurrogate(r) {
				if len(cs) >= int(i+6) {
					var cx = cs[i+1:][:6:6] // [:6] affords bounds-check-elimination
					if cx[0] == '\\' && cx[1] == 'u' {
						i += 2
						var rr1 uint32
						for j = 2; j < 6; j++ {
							c = cx[j]
							if c >= '0' && c <= '9' {
								rr = rr*16 + uint32(c-jsonU4Chk2)
							} else if c >= 'a' && c <= 'f' {
								rr = rr*16 + uint32(c-jsonU4Chk1)
							} else if c >= 'A' && c <= 'F' {
								rr = rr*16 + uint32(c-jsonU4Chk0)
							} else {
								r = unicode.ReplacementChar
								i += 4
								goto encode_rune
							}
						}
						r = utf16.DecodeRune(r, rune(rr1))
						i += 4
						goto encode_rune
					}
				}
				r = unicode.ReplacementChar
			}
		encode_rune:
			w2 := utf8.EncodeRune(d.bstr[:], r)
			v = append(v, d.bstr[:w2]...)
		default:
			d.d.errorf("unsupported escaped value: %c", c)
		}
		i++
		cursor = i
	}
	d.bs = v
}

func (d *jsonDecDriver) nakedNum(z *decNaked, bs []byte) (err error) {
	const cutoff = uint64(1 << uint(64-1))

	var n uint64
	var neg, badsyntax, overflow bool

	if len(bs) == 0 {
		if d.h.PreferFloat {
			z.v = valueTypeFloat
			z.f = 0
		} else if d.h.SignedInteger {
			z.v = valueTypeInt
			z.i = 0
		} else {
			z.v = valueTypeUint
			z.u = 0
		}
		return
	}
	if d.h.PreferFloat {
		goto F
	}
	n, neg, badsyntax, overflow = jsonParseInteger(bs)
	if badsyntax || overflow {
		goto F
	}
	if neg {
		if n > cutoff {
			goto F
		}
		z.v = valueTypeInt
		z.i = -(int64(n))
	} else if d.h.SignedInteger {
		if n >= cutoff {
			goto F
		}
		z.v = valueTypeInt
		z.i = int64(n)
	} else {
		z.v = valueTypeUint
		z.u = n
	}
	return
F:
	z.v = valueTypeFloat
	z.f, err = strconv.ParseFloat(stringView(bs), 64)
	return
}

func (d *jsonDecDriver) bsToString() string {
	// if x := d.s.sc; x != nil && x.so && x.st == '}' { // map key
	if jsonAlwaysReturnInternString || d.d.c == containerMapKey {
		return d.d.string(d.bs)
	}
	return string(d.bs)
}

func (d *jsonDecDriver) DecodeNaked() {
	z := d.d.naked()
	// var decodeFurther bool

	if d.tok == 0 {
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	}
	switch d.tok {
	case 'n':
		d.readLit4Null()
		z.v = valueTypeNil
	case 'f':
		d.readLit4False()
		z.v = valueTypeBool
		z.b = false
	case 't':
		d.readLit4True()
		z.v = valueTypeBool
		z.b = true
	case '{':
		z.v = valueTypeMap // don't consume. kInterfaceNaked will call ReadMapStart
	case '[':
		z.v = valueTypeArray // don't consume. kInterfaceNaked will call ReadArrayStart
	case '"':
		// if a string, and MapKeyAsString, then try to decode it as a nil, bool or number first
		d.appendStringAsBytes()
		if len(d.bs) > 0 && d.d.c == containerMapKey && d.h.MapKeyAsString {
			switch stringView(d.bs) {
			case "null":
				z.v = valueTypeNil
			case "true":
				z.v = valueTypeBool
				z.b = true
			case "false":
				z.v = valueTypeBool
				z.b = false
			default:
				// check if a number: float, int or uint
				if err := d.nakedNum(z, d.bs); err != nil {
					z.v = valueTypeString
					z.s = d.bsToString()
				}
			}
		} else {
			z.v = valueTypeString
			z.s = d.bsToString()
		}
	default: // number
		bs := d.decNumBytes()
		if len(bs) == 0 {
			d.d.errorf("decode number from empty string")
			return
		}
		if err := d.nakedNum(z, bs); err != nil {
			d.d.errorf("decode number from %s: %v", bs, err)
			return
		}
	}
	// if decodeFurther {
	// 	d.s.sc.retryRead()
	// }
}

//----------------------

// JsonHandle is a handle for JSON encoding format.
//
// Json is comprehensively supported:
//    - decodes numbers into interface{} as int, uint or float64
//      based on how the number looks and some config parameters e.g. PreferFloat, SignedInt, etc.
//    - decode integers from float formatted numbers e.g. 1.27e+8
//    - decode any json value (numbers, bool, etc) from quoted strings
//    - configurable way to encode/decode []byte .
//      by default, encodes and decodes []byte using base64 Std Encoding
//    - UTF-8 support for encoding and decoding
//
// It has better performance than the json library in the standard library,
// by leveraging the performance improvements of the codec library.
//
// In addition, it doesn't read more bytes than necessary during a decode, which allows
// reading multiple values from a stream containing json and non-json content.
// For example, a user can read a json value, then a cbor value, then a msgpack value,
// all from the same stream in sequence.
//
// Note that, when decoding quoted strings, invalid UTF-8 or invalid UTF-16 surrogate pairs are
// not treated as an error. Instead, they are replaced by the Unicode replacement character U+FFFD.
type JsonHandle struct {
	textEncodingType
	BasicHandle

	// Indent indicates how a value is encoded.
	//   - If positive, indent by that number of spaces.
	//   - If negative, indent by that number of tabs.
	Indent int8

	// IntegerAsString controls how integers (signed and unsigned) are encoded.
	//
	// Per the JSON Spec, JSON numbers are 64-bit floating point numbers.
	// Consequently, integers > 2^53 cannot be represented as a JSON number without losing precision.
	// This can be mitigated by configuring how to encode integers.
	//
	// IntegerAsString interpretes the following values:
	//   - if 'L', then encode integers > 2^53 as a json string.
	//   - if 'A', then encode all integers as a json string
	//             containing the exact integer representation as a decimal.
	//   - else    encode all integers as a json number (default)
	IntegerAsString byte

	// HTMLCharsAsIs controls how to encode some special characters to html: < > &
	//
	// By default, we encode them as \uXXX
	// to prevent security holes when served from some browsers.
	HTMLCharsAsIs bool

	// PreferFloat says that we will default to decoding a number as a float.
	// If not set, we will examine the characters of the number and decode as an
	// integer type if it doesn't have any of the characters [.eE].
	PreferFloat bool

	// TermWhitespace says that we add a whitespace character
	// at the end of an encoding.
	//
	// The whitespace is important, especially if using numbers in a context
	// where multiple items are written to a stream.
	TermWhitespace bool

	// MapKeyAsString says to encode all map keys as strings.
	//
	// Use this to enforce strict json output.
	// The only caveat is that nil value is ALWAYS written as null (never as "null")
	MapKeyAsString bool

	// _ uint64 // padding (cache line)

	// Note: below, we store hardly-used items
	// e.g. RawBytesExt (which is already cached in the (en|de)cDriver).

	// RawBytesExt, if configured, is used to encode and decode raw bytes in a custom way.
	// If not configured, raw bytes are encoded to/from base64 text.
	RawBytesExt InterfaceExt

	// _ [2]uint64 // padding
}

// Name returns the name of the handle: json
func (h *JsonHandle) Name() string            { return "json" }
func (h *JsonHandle) hasElemSeparators() bool { return true }
func (h *JsonHandle) typical() bool {
	return h.Indent == 0 && !h.MapKeyAsString && h.IntegerAsString != 'A' && h.IntegerAsString != 'L'
}

func (h *JsonHandle) recreateEncDriver(ed encDriver) (v bool) {
	_, v = ed.(*jsonEncDriverTypical)
	return v != h.typical()
}

// SetInterfaceExt sets an extension
func (h *JsonHandle) SetInterfaceExt(rt reflect.Type, tag uint64, ext InterfaceExt) (err error) {
	return h.SetExt(rt, tag, &interfaceExtWrapper{InterfaceExt: ext})
}

func (h *JsonHandle) newEncDriver(e *Encoder) (ee encDriver) {
	var hd *jsonEncDriver
	if h.typical() {
		var v jsonEncDriverTypical
		ee = &v
		hd = &v.jsonEncDriver
	} else {
		var v jsonEncDriverGeneric
		ee = &v
		hd = &v.jsonEncDriver
	}
	hd.e, hd.h, hd.bs = e, h, hd.b[:0]
	ee.reset()
	return
}

func (h *JsonHandle) newDecDriver(d *Decoder) decDriver {
	// d := jsonDecDriver{r: r.(*bytesDecReader), h: h}
	hd := jsonDecDriver{d: d, h: h}
	hd.bs = hd.b[:0]
	hd.reset()
	return &hd
}

func (e *jsonEncDriver) reset() {
	e.w = e.e.w()
	// (htmlasis && jsonCharSafeSet.isset(b)) || jsonCharHtmlSafeSet.isset(b)
	if e.h.HTMLCharsAsIs {
		e.s = &jsonCharSafeSet
	} else {
		e.s = &jsonCharHtmlSafeSet
	}
	e.se.InterfaceExt = e.h.RawBytesExt
	if e.bs == nil {
		e.bs = e.b[:0]
	} else {
		e.bs = e.bs[:0]
	}
}

func (d *jsonDecDriver) reset() {
	d.r = d.d.r()
	d.se.InterfaceExt = d.h.RawBytesExt
	if d.bs != nil {
		d.bs = d.bs[:0]
	}
	d.tok = 0
	// d.n.reset()
}

// jsonFloatStrconvFmtPrec ...
//
// ensure that every float has an 'e' or '.' in it,/ for easy differentiation from integers.
// this is better/faster than checking if  encoded value has [e.] and appending if needed.

// func jsonFloatStrconvFmtPrec(f float64, bits32 bool) (fmt byte, prec int) {
// 	fmt = 'f'
// 	prec = -1
// 	var abs = math.Abs(f)
// 	if abs == 0 || abs == 1 {
// 		prec = 1
// 	} else if !bits32 && (abs < 1e-6 || abs >= 1e21) ||
// 		bits32 && (float32(abs) < 1e-6 || float32(abs) >= 1e21) {
// 		fmt = 'e'
// 	} else if _, frac := math.Modf(abs); frac == 0 {
// 		// ensure that floats have a .0 at the end, for easy identification as floats
// 		prec = 1
// 	}
// 	return
// }

func jsonFloatStrconvFmtPrec64(f float64) (fmt byte, prec int8) {
	fmt = 'f'
	prec = -1
	var abs = math.Abs(f)
	if abs == 0 || abs == 1 {
		prec = 1
	} else if abs < 1e-6 || abs >= 1e21 {
		fmt = 'e'
	} else if noFrac64(abs) { // _, frac := math.Modf(abs); frac == 0 {
		prec = 1
	}
	return
}

func jsonFloatStrconvFmtPrec32(f float32) (fmt byte, prec int8) {
	fmt = 'f'
	prec = -1
	var abs = abs32(f)
	if abs == 0 || abs == 1 {
		prec = 1
	} else if abs < 1e-6 || abs >= 1e21 {
		fmt = 'e'
	} else if noFrac32(abs) { // _, frac := math.Modf(abs); frac == 0 {
		prec = 1
	}
	return
}

// custom-fitted version of strconv.Parse(Ui|I)nt.
// Also ensures we don't have to search for .eE to determine if a float or not.
// Note: s CANNOT be a zero-length slice.
func jsonParseInteger(s []byte) (n uint64, neg, badSyntax, overflow bool) {
	const maxUint64 = (1<<64 - 1)
	const cutoff = maxUint64/10 + 1

	if len(s) == 0 { // bounds-check-elimination
		// treat empty string as zero value
		// badSyntax = true
		return
	}
	switch s[0] {
	case '+':
		s = s[1:]
	case '-':
		s = s[1:]
		neg = true
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			badSyntax = true
			return
		}
		// unsigned integers don't overflow well on multiplication, so check cutoff here
		// e.g. (maxUint64-5)*10 doesn't overflow well ...
		if n >= cutoff {
			overflow = true
			return
		}
		n *= 10
		n1 := n + uint64(c-'0')
		if n1 < n || n1 > maxUint64 {
			overflow = true
			return
		}
		n = n1
	}
	return
}

var _ decDriverContainerTracker = (*jsonDecDriver)(nil)
var _ encDriverContainerTracker = (*jsonEncDriver)(nil)
var _ decDriver = (*jsonDecDriver)(nil)
var _ encDriver = (*jsonEncDriverGeneric)(nil)
var _ encDriver = (*jsonEncDriverTypical)(nil)
var _ (interface{ getJsonEncDriver() *jsonEncDriver }) = (*jsonEncDriverTypical)(nil)
var _ (interface{ getJsonEncDriver() *jsonEncDriver }) = (*jsonEncDriverGeneric)(nil)
var _ (interface{ getJsonEncDriver() *jsonEncDriver }) = (*jsonEncDriver)(nil)

// var _ encDriver = (*jsonEncDriver)(nil)
