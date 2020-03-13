package strftime

import (
	"io"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var directives = map[byte]appender{
	'A': timefmt("Monday"),
	'a': timefmt("Mon"),
	'B': timefmt("January"),
	'b': timefmt("Jan"),
	'C': &century{},
	'c': timefmt("Mon Jan _2 15:04:05 2006"),
	'D': timefmt("01/02/06"),
	'd': timefmt("02"),
	'e': timefmt("_2"),
	'F': timefmt("2006-01-02"),
	'H': timefmt("15"),
	'h': timefmt("Jan"), // same as 'b'
	'I': timefmt("3"),
	'j': &dayofyear{},
	'k': hourwblank(false),
	'l': hourwblank(true),
	'M': timefmt("04"),
	'm': timefmt("01"),
	'n': verbatim("\n"),
	'p': timefmt("PM"),
	'R': timefmt("15:04"),
	'r': timefmt("3:04:05 PM"),
	'S': timefmt("05"),
	'T': timefmt("15:04:05"),
	't': verbatim("\t"),
	'U': weeknumberOffset(0), // week number of the year, Sunday first
	'u': weekday(1),
	'V': &weeknumber{},
	'v': timefmt("_2-Jan-2006"),
	'W': weeknumberOffset(1), // week number of the year, Monday first
	'w': weekday(0),
	'X': timefmt("15:04:05"), // national representation of the time XXX is this correct?
	'x': timefmt("01/02/06"), // national representation of the date XXX is this correct?
	'Y': timefmt("2006"),     // year with century
	'y': timefmt("06"),       // year w/o century
	'Z': timefmt("MST"),      // time zone name
	'z': timefmt("-0700"),    // time zone ofset from UTC
	'%': verbatim("%"),
}

type combiningAppend struct {
	list appenderList
	prev appender
	prevCanCombine bool
}

func (ca *combiningAppend) Append(w appender) {
	if ca.prevCanCombine {
		if wc, ok := w.(combiner); ok && wc.canCombine() {
			ca.prev = ca.prev.(combiner).combine(wc)
			ca.list[len(ca.list)-1] = ca.prev
			return
		}
	}

	ca.list = append(ca.list, w)
	ca.prev = w
	ca.prevCanCombine = false
	if comb, ok := w.(combiner); ok {
		if comb.canCombine() {
			ca.prevCanCombine = true
		}
	}
}

func compile(wl *appenderList, p string) error {
	var ca combiningAppend
	for l := len(p); l > 0; l = len(p) {
		i := strings.IndexByte(p, '%')
		if i < 0 {
			ca.Append(verbatim(p))
			// this is silly, but I don't trust break keywords when there's a
			// possibility of this piece of code being rearranged
			p = p[l:]
			continue
		}
		if i == l-1 {
			return errors.New(`stray % at the end of pattern`)
		}

		// we found a '%'. we need the next byte to decide what to do next
		// we already know that i < l - 1
		// everything up to the i is verbatim
		if i > 0 {
			ca.Append(verbatim(p[:i]))
			p = p[i:]
		}

		directive, ok := directives[p[1]]
		if !ok {
			return errors.Errorf(`unknown time format specification '%c'`, p[1])
		}
		ca.Append(directive)
		p = p[2:]
	}

	*wl = ca.list

	return nil
}

// Format takes the format `s` and the time `t` to produce the
// format date/time. Note that this function re-compiles the
// pattern every time it is called.
//
// If you know beforehand that you will be reusing the pattern
// within your application, consider creating a `Strftime` object
// and reusing it.
func Format(p string, t time.Time) (string, error) {
	var dst []byte
	// TODO: optimize for 64 byte strings
	dst = make([]byte, 0, len(p)+10)
	// Compile, but execute as we go
	for l := len(p); l > 0; l = len(p) {
		i := strings.IndexByte(p, '%')
		if i < 0 {
			dst = append(dst, p...)
			// this is silly, but I don't trust break keywords when there's a
			// possibility of this piece of code being rearranged
			p = p[l:]
			continue
		}
		if i == l-1 {
			return "", errors.New(`stray % at the end of pattern`)
		}

		// we found a '%'. we need the next byte to decide what to do next
		// we already know that i < l - 1
		// everything up to the i is verbatim
		if i > 0 {
			dst = append(dst, p[:i]...)
			p = p[i:]
		}

		directive, ok := directives[p[1]]
		if !ok {
			return "", errors.Errorf(`unknown time format specification '%c'`, p[1])
		}
		dst = directive.Append(dst, t)
		p = p[2:]
	}

	return string(dst), nil
}

// Strftime is the object that represents a compiled strftime pattern
type Strftime struct {
	pattern  string
	compiled appenderList
}

// New creates a new Strftime object. If the compilation fails, then
// an error is returned in the second argument.
func New(f string) (*Strftime, error) {
	var wl appenderList
	if err := compile(&wl, f); err != nil {
		return nil, errors.Wrap(err, `failed to compile format`)
	}
	return &Strftime{
		pattern:  f,
		compiled: wl,
	}, nil
}

// Pattern returns the original pattern string
func (f *Strftime) Pattern() string {
	return f.pattern
}

// Format takes the destination `dst` and time `t`. It formats the date/time
// using the pre-compiled pattern, and outputs the results to `dst`
func (f *Strftime) Format(dst io.Writer, t time.Time) error {
	const bufSize = 64
	var b []byte
	max := len(f.pattern) + 10
	if max < bufSize {
		var buf [bufSize]byte
		b = buf[:0]
	} else {
		b = make([]byte, 0, max)
	}
	if _, err := dst.Write(f.format(b, t)); err != nil {
		return err
	}
	return nil
}

func (f *Strftime) format(b []byte, t time.Time) []byte {
	for _, w := range f.compiled {
		b = w.Append(b, t)
	}
	return b
}

// FormatString takes the time `t` and formats it, returning the
// string containing the formated data.
func (f *Strftime) FormatString(t time.Time) string {
	const bufSize = 64
	var b []byte
	max := len(f.pattern) + 10
	if max < bufSize {
		var buf [bufSize]byte
		b = buf[:0]
	} else {
		b = make([]byte, 0, max)
	}
	return string(f.format(b, t))
}
