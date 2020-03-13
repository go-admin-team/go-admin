// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

// Code generated from sort-slice.go.tmpl - DO NOT EDIT.

package codec

import "time"
import "reflect"
import "bytes"

type stringSlice []string

func (p stringSlice) Len() int      { return len(p) }
func (p stringSlice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p stringSlice) Less(i, j int) bool {
	return p[uint(i)] < p[uint(j)]
}

type float32Slice []float32

func (p float32Slice) Len() int      { return len(p) }
func (p float32Slice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p float32Slice) Less(i, j int) bool {
	return p[uint(i)] < p[uint(j)] || isNaN32(p[uint(i)]) && !isNaN32(p[uint(j)])
}

type float64Slice []float64

func (p float64Slice) Len() int      { return len(p) }
func (p float64Slice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p float64Slice) Less(i, j int) bool {
	return p[uint(i)] < p[uint(j)] || isNaN64(p[uint(i)]) && !isNaN64(p[uint(j)])
}

type uintSlice []uint

func (p uintSlice) Len() int      { return len(p) }
func (p uintSlice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p uintSlice) Less(i, j int) bool {
	return p[uint(i)] < p[uint(j)]
}

type uint8Slice []uint8

func (p uint8Slice) Len() int      { return len(p) }
func (p uint8Slice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p uint8Slice) Less(i, j int) bool {
	return p[uint(i)] < p[uint(j)]
}

type uint16Slice []uint16

func (p uint16Slice) Len() int      { return len(p) }
func (p uint16Slice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p uint16Slice) Less(i, j int) bool {
	return p[uint(i)] < p[uint(j)]
}

type uint32Slice []uint32

func (p uint32Slice) Len() int      { return len(p) }
func (p uint32Slice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p uint32Slice) Less(i, j int) bool {
	return p[uint(i)] < p[uint(j)]
}

type uint64Slice []uint64

func (p uint64Slice) Len() int      { return len(p) }
func (p uint64Slice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p uint64Slice) Less(i, j int) bool {
	return p[uint(i)] < p[uint(j)]
}

type uintptrSlice []uintptr

func (p uintptrSlice) Len() int      { return len(p) }
func (p uintptrSlice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p uintptrSlice) Less(i, j int) bool {
	return p[uint(i)] < p[uint(j)]
}

type intSlice []int

func (p intSlice) Len() int      { return len(p) }
func (p intSlice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p intSlice) Less(i, j int) bool {
	return p[uint(i)] < p[uint(j)]
}

type int8Slice []int8

func (p int8Slice) Len() int      { return len(p) }
func (p int8Slice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p int8Slice) Less(i, j int) bool {
	return p[uint(i)] < p[uint(j)]
}

type int16Slice []int16

func (p int16Slice) Len() int      { return len(p) }
func (p int16Slice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p int16Slice) Less(i, j int) bool {
	return p[uint(i)] < p[uint(j)]
}

type int32Slice []int32

func (p int32Slice) Len() int      { return len(p) }
func (p int32Slice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p int32Slice) Less(i, j int) bool {
	return p[uint(i)] < p[uint(j)]
}

type int64Slice []int64

func (p int64Slice) Len() int      { return len(p) }
func (p int64Slice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p int64Slice) Less(i, j int) bool {
	return p[uint(i)] < p[uint(j)]
}

type boolSlice []bool

func (p boolSlice) Len() int      { return len(p) }
func (p boolSlice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p boolSlice) Less(i, j int) bool {
	return !p[uint(i)] && p[uint(j)]
}

type timeSlice []time.Time

func (p timeSlice) Len() int      { return len(p) }
func (p timeSlice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p timeSlice) Less(i, j int) bool {
	return p[uint(i)].Before(p[uint(j)])
}

type bytesSlice [][]byte

func (p bytesSlice) Len() int      { return len(p) }
func (p bytesSlice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p bytesSlice) Less(i, j int) bool {
	return bytes.Compare(p[uint(i)], p[uint(j)]) == -1
}

type stringRv struct {
	v string
	r reflect.Value
}
type stringRvSlice []stringRv

func (p stringRvSlice) Len() int      { return len(p) }
func (p stringRvSlice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p stringRvSlice) Less(i, j int) bool {
	return p[uint(i)].v < p[uint(j)].v
}

type stringIntf struct {
	v string
	i interface{}
}
type stringIntfSlice []stringIntf

func (p stringIntfSlice) Len() int      { return len(p) }
func (p stringIntfSlice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p stringIntfSlice) Less(i, j int) bool {
	return p[uint(i)].v < p[uint(j)].v
}

type float64Rv struct {
	v float64
	r reflect.Value
}
type float64RvSlice []float64Rv

func (p float64RvSlice) Len() int      { return len(p) }
func (p float64RvSlice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p float64RvSlice) Less(i, j int) bool {
	return p[uint(i)].v < p[uint(j)].v || isNaN64(p[uint(i)].v) && !isNaN64(p[uint(j)].v)
}

type float64Intf struct {
	v float64
	i interface{}
}
type float64IntfSlice []float64Intf

func (p float64IntfSlice) Len() int      { return len(p) }
func (p float64IntfSlice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p float64IntfSlice) Less(i, j int) bool {
	return p[uint(i)].v < p[uint(j)].v || isNaN64(p[uint(i)].v) && !isNaN64(p[uint(j)].v)
}

type uint64Rv struct {
	v uint64
	r reflect.Value
}
type uint64RvSlice []uint64Rv

func (p uint64RvSlice) Len() int      { return len(p) }
func (p uint64RvSlice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p uint64RvSlice) Less(i, j int) bool {
	return p[uint(i)].v < p[uint(j)].v
}

type uint64Intf struct {
	v uint64
	i interface{}
}
type uint64IntfSlice []uint64Intf

func (p uint64IntfSlice) Len() int      { return len(p) }
func (p uint64IntfSlice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p uint64IntfSlice) Less(i, j int) bool {
	return p[uint(i)].v < p[uint(j)].v
}

type uintptrRv struct {
	v uintptr
	r reflect.Value
}
type uintptrRvSlice []uintptrRv

func (p uintptrRvSlice) Len() int      { return len(p) }
func (p uintptrRvSlice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p uintptrRvSlice) Less(i, j int) bool {
	return p[uint(i)].v < p[uint(j)].v
}

type uintptrIntf struct {
	v uintptr
	i interface{}
}
type uintptrIntfSlice []uintptrIntf

func (p uintptrIntfSlice) Len() int      { return len(p) }
func (p uintptrIntfSlice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p uintptrIntfSlice) Less(i, j int) bool {
	return p[uint(i)].v < p[uint(j)].v
}

type int64Rv struct {
	v int64
	r reflect.Value
}
type int64RvSlice []int64Rv

func (p int64RvSlice) Len() int      { return len(p) }
func (p int64RvSlice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p int64RvSlice) Less(i, j int) bool {
	return p[uint(i)].v < p[uint(j)].v
}

type int64Intf struct {
	v int64
	i interface{}
}
type int64IntfSlice []int64Intf

func (p int64IntfSlice) Len() int      { return len(p) }
func (p int64IntfSlice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p int64IntfSlice) Less(i, j int) bool {
	return p[uint(i)].v < p[uint(j)].v
}

type boolRv struct {
	v bool
	r reflect.Value
}
type boolRvSlice []boolRv

func (p boolRvSlice) Len() int      { return len(p) }
func (p boolRvSlice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p boolRvSlice) Less(i, j int) bool {
	return !p[uint(i)].v && p[uint(j)].v
}

type boolIntf struct {
	v bool
	i interface{}
}
type boolIntfSlice []boolIntf

func (p boolIntfSlice) Len() int      { return len(p) }
func (p boolIntfSlice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p boolIntfSlice) Less(i, j int) bool {
	return !p[uint(i)].v && p[uint(j)].v
}

type timeRv struct {
	v time.Time
	r reflect.Value
}
type timeRvSlice []timeRv

func (p timeRvSlice) Len() int      { return len(p) }
func (p timeRvSlice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p timeRvSlice) Less(i, j int) bool {
	return p[uint(i)].v.Before(p[uint(j)].v)
}

type timeIntf struct {
	v time.Time
	i interface{}
}
type timeIntfSlice []timeIntf

func (p timeIntfSlice) Len() int      { return len(p) }
func (p timeIntfSlice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p timeIntfSlice) Less(i, j int) bool {
	return p[uint(i)].v.Before(p[uint(j)].v)
}

type bytesRv struct {
	v []byte
	r reflect.Value
}
type bytesRvSlice []bytesRv

func (p bytesRvSlice) Len() int      { return len(p) }
func (p bytesRvSlice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p bytesRvSlice) Less(i, j int) bool {
	return bytes.Compare(p[uint(i)].v, p[uint(j)].v) == -1
}

type bytesIntf struct {
	v []byte
	i interface{}
}
type bytesIntfSlice []bytesIntf

func (p bytesIntfSlice) Len() int      { return len(p) }
func (p bytesIntfSlice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p bytesIntfSlice) Less(i, j int) bool {
	return bytes.Compare(p[uint(i)].v, p[uint(j)].v) == -1
}
