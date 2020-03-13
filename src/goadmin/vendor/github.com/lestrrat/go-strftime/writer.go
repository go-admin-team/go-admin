package strftime

import (
	"strconv"
	"strings"
	"time"
)

type appender interface {
	Append([]byte, time.Time) []byte
}

type appenderList []appender

// does the time.Format thing
type timefmtw struct {
	s string
}

func timefmt(s string) *timefmtw {
	return &timefmtw{s: s}
}

func (v timefmtw) Append(b []byte, t time.Time) []byte {
	return t.AppendFormat(b, v.s)
}

func (v timefmtw) str() string {
	return v.s
}

func (v timefmtw) canCombine() bool {
	return true
}

func (v timefmtw) combine(w combiner) appender {
	return timefmt(v.s + w.str())
}

type verbatimw struct {
	s string
}

func verbatim(s string) *verbatimw {
	return &verbatimw{s: s}
}

func (v verbatimw) Append(b []byte, _ time.Time) []byte {
	return append(b, v.s...)
}

func (v verbatimw) canCombine() bool {
	return canCombine(v.s)
}

func (v verbatimw) combine(w combiner) appender {
	if _, ok := w.(*timefmtw); ok {
		return timefmt(v.s + w.str())
	}
	return verbatim(v.s + w.str())
}

func (v verbatimw) str() string {
	return v.s
}

// These words below, as well as any decimal character
var combineExclusion = []string{
	"Mon",
	"Monday",
	"Jan",
	"January",
	"MST",
	"PM",
}

func canCombine(s string) bool {
	if strings.ContainsAny(s, "0123456789") {
		return false
	}
	for _, word := range combineExclusion {
		if strings.Contains(s, word) {
			return false
		}
	}
	return true
}

type combiner interface {
	canCombine() bool
	combine(combiner) appender
	str() string
}

type century struct{}

func (v century) Append(b []byte, t time.Time) []byte {
	n := t.Year() / 100
	if n < 10 {
		b = append(b, '0')
	}
	return append(b, strconv.Itoa(n)...)
}

type weekday int

func (v weekday) Append(b []byte, t time.Time) []byte {
	n := int(t.Weekday())
	if n < int(v) {
		n += 7
	}
	return append(b, byte(n+48))
}

type weeknumberOffset int

func (v weeknumberOffset) Append(b []byte, t time.Time) []byte {
	yd := t.YearDay()
	offset := int(t.Weekday()) - int(v)
	if offset < 0 {
		offset += 7
	}

	if yd < offset {
		return append(b, '0', '0')
	}

	n := ((yd - offset) / 7) + 1
	if n < 10 {
		b = append(b, '0')
	}
	return append(b, strconv.Itoa(n)...)
}

type weeknumber struct{}

func (v weeknumber) Append(b []byte, t time.Time) []byte {
	_, n := t.ISOWeek()
	if n < 10 {
		b = append(b, '0')
	}
	return append(b, strconv.Itoa(n)...)
}

type dayofyear struct{}

func (v dayofyear) Append(b []byte, t time.Time) []byte {
	n := t.YearDay()
	if n < 10 {
		b = append(b, '0', '0')
	} else if n < 100 {
		b = append(b, '0')
	}
	return append(b, strconv.Itoa(n)...)
}

type hourwblank bool

func (v hourwblank) Append(b []byte, t time.Time) []byte {
	h := t.Hour()
	if bool(v) && h > 12 {
		h = h - 12
	}
	if h < 10 {
		b = append(b, ' ')
	}
	return append(b, strconv.Itoa(h)...)
}
