// +build !notfastpath

// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

// Code generated from fast-path.go.tmpl - DO NOT EDIT.

package codec

// Fast path functions try to create a fast path encode or decode implementation
// for common maps and slices.
//
// We define the functions and register them in this single file
// so as not to pollute the encode.go and decode.go, and create a dependency in there.
// This file can be omitted without causing a build failure.
//
// The advantage of fast paths is:
//	  - Many calls bypass reflection altogether
//
// Currently support
//	  - slice of all builtin types (numeric, bool, string, []byte)
//    - maps of builtin types to builtin or interface{} type, EXCEPT FOR
//      keys of type uintptr, int8/16/32, uint16/32, float32/64, bool, interface{}
//      AND values of type type int8/16/32, uint16/32
// This should provide adequate "typical" implementations.
//
// Note that fast track decode functions must handle values for which an address cannot be obtained.
// For example:
//	 m2 := map[string]int{}
//	 p2 := []interface{}{m2}
//	 // decoding into p2 will bomb if fast track functions do not treat like unaddressable.
//

import (
	"reflect"
	"sort"
)

const fastpathEnabled = true

const fastpathMapBySliceErrMsg = "mapBySlice requires even slice length, but got %v"

type fastpathT struct{}

var fastpathTV fastpathT

type fastpathE struct {
	rtid  uintptr
	rt    reflect.Type
	encfn func(*Encoder, *codecFnInfo, reflect.Value)
	decfn func(*Decoder, *codecFnInfo, reflect.Value)
}

type fastpathA [88]fastpathE

func (x *fastpathA) index(rtid uintptr) int {
	// use binary search to grab the index (adapted from sort/search.go)
	// Note: we use goto (instead of for loop) so this can be inlined.
	// h, i, j := 0, 0, len(x)
	var h, i uint
	var j = uint(len(x))
LOOP:
	if i < j {
		h = i + (j-i)/2
		if x[h].rtid < rtid {
			i = h + 1
		} else {
			j = h
		}
		goto LOOP
	}
	if i < uint(len(x)) && x[i].rtid == rtid {
		return int(i)
	}
	return -1
}

type fastpathAslice []fastpathE

func (x fastpathAslice) Len() int           { return len(x) }
func (x fastpathAslice) Less(i, j int) bool { return x[uint(i)].rtid < x[uint(j)].rtid }
func (x fastpathAslice) Swap(i, j int)      { x[uint(i)], x[uint(j)] = x[uint(j)], x[uint(i)] }

var fastpathAV fastpathA

// due to possible initialization loop error, make fastpath in an init()
func init() {
	var i uint = 0
	fn := func(v interface{},
		fe func(*Encoder, *codecFnInfo, reflect.Value),
		fd func(*Decoder, *codecFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		xptr := rt2id(xrt)
		fastpathAV[i] = fastpathE{xptr, xrt, fe, fd}
		i++
	}

	fn([]interface{}(nil), (*Encoder).fastpathEncSliceIntfR, (*Decoder).fastpathDecSliceIntfR)
	fn([]string(nil), (*Encoder).fastpathEncSliceStringR, (*Decoder).fastpathDecSliceStringR)
	fn([][]byte(nil), (*Encoder).fastpathEncSliceBytesR, (*Decoder).fastpathDecSliceBytesR)
	fn([]float32(nil), (*Encoder).fastpathEncSliceFloat32R, (*Decoder).fastpathDecSliceFloat32R)
	fn([]float64(nil), (*Encoder).fastpathEncSliceFloat64R, (*Decoder).fastpathDecSliceFloat64R)
	fn([]uint(nil), (*Encoder).fastpathEncSliceUintR, (*Decoder).fastpathDecSliceUintR)
	fn([]uint16(nil), (*Encoder).fastpathEncSliceUint16R, (*Decoder).fastpathDecSliceUint16R)
	fn([]uint32(nil), (*Encoder).fastpathEncSliceUint32R, (*Decoder).fastpathDecSliceUint32R)
	fn([]uint64(nil), (*Encoder).fastpathEncSliceUint64R, (*Decoder).fastpathDecSliceUint64R)
	fn([]uintptr(nil), (*Encoder).fastpathEncSliceUintptrR, (*Decoder).fastpathDecSliceUintptrR)
	fn([]int(nil), (*Encoder).fastpathEncSliceIntR, (*Decoder).fastpathDecSliceIntR)
	fn([]int8(nil), (*Encoder).fastpathEncSliceInt8R, (*Decoder).fastpathDecSliceInt8R)
	fn([]int16(nil), (*Encoder).fastpathEncSliceInt16R, (*Decoder).fastpathDecSliceInt16R)
	fn([]int32(nil), (*Encoder).fastpathEncSliceInt32R, (*Decoder).fastpathDecSliceInt32R)
	fn([]int64(nil), (*Encoder).fastpathEncSliceInt64R, (*Decoder).fastpathDecSliceInt64R)
	fn([]bool(nil), (*Encoder).fastpathEncSliceBoolR, (*Decoder).fastpathDecSliceBoolR)

	fn(map[string]interface{}(nil), (*Encoder).fastpathEncMapStringIntfR, (*Decoder).fastpathDecMapStringIntfR)
	fn(map[string]string(nil), (*Encoder).fastpathEncMapStringStringR, (*Decoder).fastpathDecMapStringStringR)
	fn(map[string][]byte(nil), (*Encoder).fastpathEncMapStringBytesR, (*Decoder).fastpathDecMapStringBytesR)
	fn(map[string]uint(nil), (*Encoder).fastpathEncMapStringUintR, (*Decoder).fastpathDecMapStringUintR)
	fn(map[string]uint8(nil), (*Encoder).fastpathEncMapStringUint8R, (*Decoder).fastpathDecMapStringUint8R)
	fn(map[string]uint64(nil), (*Encoder).fastpathEncMapStringUint64R, (*Decoder).fastpathDecMapStringUint64R)
	fn(map[string]uintptr(nil), (*Encoder).fastpathEncMapStringUintptrR, (*Decoder).fastpathDecMapStringUintptrR)
	fn(map[string]int(nil), (*Encoder).fastpathEncMapStringIntR, (*Decoder).fastpathDecMapStringIntR)
	fn(map[string]int64(nil), (*Encoder).fastpathEncMapStringInt64R, (*Decoder).fastpathDecMapStringInt64R)
	fn(map[string]float32(nil), (*Encoder).fastpathEncMapStringFloat32R, (*Decoder).fastpathDecMapStringFloat32R)
	fn(map[string]float64(nil), (*Encoder).fastpathEncMapStringFloat64R, (*Decoder).fastpathDecMapStringFloat64R)
	fn(map[string]bool(nil), (*Encoder).fastpathEncMapStringBoolR, (*Decoder).fastpathDecMapStringBoolR)
	fn(map[uint]interface{}(nil), (*Encoder).fastpathEncMapUintIntfR, (*Decoder).fastpathDecMapUintIntfR)
	fn(map[uint]string(nil), (*Encoder).fastpathEncMapUintStringR, (*Decoder).fastpathDecMapUintStringR)
	fn(map[uint][]byte(nil), (*Encoder).fastpathEncMapUintBytesR, (*Decoder).fastpathDecMapUintBytesR)
	fn(map[uint]uint(nil), (*Encoder).fastpathEncMapUintUintR, (*Decoder).fastpathDecMapUintUintR)
	fn(map[uint]uint8(nil), (*Encoder).fastpathEncMapUintUint8R, (*Decoder).fastpathDecMapUintUint8R)
	fn(map[uint]uint64(nil), (*Encoder).fastpathEncMapUintUint64R, (*Decoder).fastpathDecMapUintUint64R)
	fn(map[uint]uintptr(nil), (*Encoder).fastpathEncMapUintUintptrR, (*Decoder).fastpathDecMapUintUintptrR)
	fn(map[uint]int(nil), (*Encoder).fastpathEncMapUintIntR, (*Decoder).fastpathDecMapUintIntR)
	fn(map[uint]int64(nil), (*Encoder).fastpathEncMapUintInt64R, (*Decoder).fastpathDecMapUintInt64R)
	fn(map[uint]float32(nil), (*Encoder).fastpathEncMapUintFloat32R, (*Decoder).fastpathDecMapUintFloat32R)
	fn(map[uint]float64(nil), (*Encoder).fastpathEncMapUintFloat64R, (*Decoder).fastpathDecMapUintFloat64R)
	fn(map[uint]bool(nil), (*Encoder).fastpathEncMapUintBoolR, (*Decoder).fastpathDecMapUintBoolR)
	fn(map[uint8]interface{}(nil), (*Encoder).fastpathEncMapUint8IntfR, (*Decoder).fastpathDecMapUint8IntfR)
	fn(map[uint8]string(nil), (*Encoder).fastpathEncMapUint8StringR, (*Decoder).fastpathDecMapUint8StringR)
	fn(map[uint8][]byte(nil), (*Encoder).fastpathEncMapUint8BytesR, (*Decoder).fastpathDecMapUint8BytesR)
	fn(map[uint8]uint(nil), (*Encoder).fastpathEncMapUint8UintR, (*Decoder).fastpathDecMapUint8UintR)
	fn(map[uint8]uint8(nil), (*Encoder).fastpathEncMapUint8Uint8R, (*Decoder).fastpathDecMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*Encoder).fastpathEncMapUint8Uint64R, (*Decoder).fastpathDecMapUint8Uint64R)
	fn(map[uint8]uintptr(nil), (*Encoder).fastpathEncMapUint8UintptrR, (*Decoder).fastpathDecMapUint8UintptrR)
	fn(map[uint8]int(nil), (*Encoder).fastpathEncMapUint8IntR, (*Decoder).fastpathDecMapUint8IntR)
	fn(map[uint8]int64(nil), (*Encoder).fastpathEncMapUint8Int64R, (*Decoder).fastpathDecMapUint8Int64R)
	fn(map[uint8]float32(nil), (*Encoder).fastpathEncMapUint8Float32R, (*Decoder).fastpathDecMapUint8Float32R)
	fn(map[uint8]float64(nil), (*Encoder).fastpathEncMapUint8Float64R, (*Decoder).fastpathDecMapUint8Float64R)
	fn(map[uint8]bool(nil), (*Encoder).fastpathEncMapUint8BoolR, (*Decoder).fastpathDecMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*Encoder).fastpathEncMapUint64IntfR, (*Decoder).fastpathDecMapUint64IntfR)
	fn(map[uint64]string(nil), (*Encoder).fastpathEncMapUint64StringR, (*Decoder).fastpathDecMapUint64StringR)
	fn(map[uint64][]byte(nil), (*Encoder).fastpathEncMapUint64BytesR, (*Decoder).fastpathDecMapUint64BytesR)
	fn(map[uint64]uint(nil), (*Encoder).fastpathEncMapUint64UintR, (*Decoder).fastpathDecMapUint64UintR)
	fn(map[uint64]uint8(nil), (*Encoder).fastpathEncMapUint64Uint8R, (*Decoder).fastpathDecMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*Encoder).fastpathEncMapUint64Uint64R, (*Decoder).fastpathDecMapUint64Uint64R)
	fn(map[uint64]uintptr(nil), (*Encoder).fastpathEncMapUint64UintptrR, (*Decoder).fastpathDecMapUint64UintptrR)
	fn(map[uint64]int(nil), (*Encoder).fastpathEncMapUint64IntR, (*Decoder).fastpathDecMapUint64IntR)
	fn(map[uint64]int64(nil), (*Encoder).fastpathEncMapUint64Int64R, (*Decoder).fastpathDecMapUint64Int64R)
	fn(map[uint64]float32(nil), (*Encoder).fastpathEncMapUint64Float32R, (*Decoder).fastpathDecMapUint64Float32R)
	fn(map[uint64]float64(nil), (*Encoder).fastpathEncMapUint64Float64R, (*Decoder).fastpathDecMapUint64Float64R)
	fn(map[uint64]bool(nil), (*Encoder).fastpathEncMapUint64BoolR, (*Decoder).fastpathDecMapUint64BoolR)
	fn(map[int]interface{}(nil), (*Encoder).fastpathEncMapIntIntfR, (*Decoder).fastpathDecMapIntIntfR)
	fn(map[int]string(nil), (*Encoder).fastpathEncMapIntStringR, (*Decoder).fastpathDecMapIntStringR)
	fn(map[int][]byte(nil), (*Encoder).fastpathEncMapIntBytesR, (*Decoder).fastpathDecMapIntBytesR)
	fn(map[int]uint(nil), (*Encoder).fastpathEncMapIntUintR, (*Decoder).fastpathDecMapIntUintR)
	fn(map[int]uint8(nil), (*Encoder).fastpathEncMapIntUint8R, (*Decoder).fastpathDecMapIntUint8R)
	fn(map[int]uint64(nil), (*Encoder).fastpathEncMapIntUint64R, (*Decoder).fastpathDecMapIntUint64R)
	fn(map[int]uintptr(nil), (*Encoder).fastpathEncMapIntUintptrR, (*Decoder).fastpathDecMapIntUintptrR)
	fn(map[int]int(nil), (*Encoder).fastpathEncMapIntIntR, (*Decoder).fastpathDecMapIntIntR)
	fn(map[int]int64(nil), (*Encoder).fastpathEncMapIntInt64R, (*Decoder).fastpathDecMapIntInt64R)
	fn(map[int]float32(nil), (*Encoder).fastpathEncMapIntFloat32R, (*Decoder).fastpathDecMapIntFloat32R)
	fn(map[int]float64(nil), (*Encoder).fastpathEncMapIntFloat64R, (*Decoder).fastpathDecMapIntFloat64R)
	fn(map[int]bool(nil), (*Encoder).fastpathEncMapIntBoolR, (*Decoder).fastpathDecMapIntBoolR)
	fn(map[int64]interface{}(nil), (*Encoder).fastpathEncMapInt64IntfR, (*Decoder).fastpathDecMapInt64IntfR)
	fn(map[int64]string(nil), (*Encoder).fastpathEncMapInt64StringR, (*Decoder).fastpathDecMapInt64StringR)
	fn(map[int64][]byte(nil), (*Encoder).fastpathEncMapInt64BytesR, (*Decoder).fastpathDecMapInt64BytesR)
	fn(map[int64]uint(nil), (*Encoder).fastpathEncMapInt64UintR, (*Decoder).fastpathDecMapInt64UintR)
	fn(map[int64]uint8(nil), (*Encoder).fastpathEncMapInt64Uint8R, (*Decoder).fastpathDecMapInt64Uint8R)
	fn(map[int64]uint64(nil), (*Encoder).fastpathEncMapInt64Uint64R, (*Decoder).fastpathDecMapInt64Uint64R)
	fn(map[int64]uintptr(nil), (*Encoder).fastpathEncMapInt64UintptrR, (*Decoder).fastpathDecMapInt64UintptrR)
	fn(map[int64]int(nil), (*Encoder).fastpathEncMapInt64IntR, (*Decoder).fastpathDecMapInt64IntR)
	fn(map[int64]int64(nil), (*Encoder).fastpathEncMapInt64Int64R, (*Decoder).fastpathDecMapInt64Int64R)
	fn(map[int64]float32(nil), (*Encoder).fastpathEncMapInt64Float32R, (*Decoder).fastpathDecMapInt64Float32R)
	fn(map[int64]float64(nil), (*Encoder).fastpathEncMapInt64Float64R, (*Decoder).fastpathDecMapInt64Float64R)
	fn(map[int64]bool(nil), (*Encoder).fastpathEncMapInt64BoolR, (*Decoder).fastpathDecMapInt64BoolR)

	sort.Sort(fastpathAslice(fastpathAV[:]))
}

// -- encode

// -- -- fast path type switch
func fastpathEncodeTypeSwitch(iv interface{}, e *Encoder) bool {
	switch v := iv.(type) {

	case []interface{}:
		fastpathTV.EncSliceIntfV(v, e)
	case *[]interface{}:
		fastpathTV.EncSliceIntfV(*v, e)
	case []string:
		fastpathTV.EncSliceStringV(v, e)
	case *[]string:
		fastpathTV.EncSliceStringV(*v, e)
	case [][]byte:
		fastpathTV.EncSliceBytesV(v, e)
	case *[][]byte:
		fastpathTV.EncSliceBytesV(*v, e)
	case []float32:
		fastpathTV.EncSliceFloat32V(v, e)
	case *[]float32:
		fastpathTV.EncSliceFloat32V(*v, e)
	case []float64:
		fastpathTV.EncSliceFloat64V(v, e)
	case *[]float64:
		fastpathTV.EncSliceFloat64V(*v, e)
	case []uint:
		fastpathTV.EncSliceUintV(v, e)
	case *[]uint:
		fastpathTV.EncSliceUintV(*v, e)
	case []uint16:
		fastpathTV.EncSliceUint16V(v, e)
	case *[]uint16:
		fastpathTV.EncSliceUint16V(*v, e)
	case []uint32:
		fastpathTV.EncSliceUint32V(v, e)
	case *[]uint32:
		fastpathTV.EncSliceUint32V(*v, e)
	case []uint64:
		fastpathTV.EncSliceUint64V(v, e)
	case *[]uint64:
		fastpathTV.EncSliceUint64V(*v, e)
	case []uintptr:
		fastpathTV.EncSliceUintptrV(v, e)
	case *[]uintptr:
		fastpathTV.EncSliceUintptrV(*v, e)
	case []int:
		fastpathTV.EncSliceIntV(v, e)
	case *[]int:
		fastpathTV.EncSliceIntV(*v, e)
	case []int8:
		fastpathTV.EncSliceInt8V(v, e)
	case *[]int8:
		fastpathTV.EncSliceInt8V(*v, e)
	case []int16:
		fastpathTV.EncSliceInt16V(v, e)
	case *[]int16:
		fastpathTV.EncSliceInt16V(*v, e)
	case []int32:
		fastpathTV.EncSliceInt32V(v, e)
	case *[]int32:
		fastpathTV.EncSliceInt32V(*v, e)
	case []int64:
		fastpathTV.EncSliceInt64V(v, e)
	case *[]int64:
		fastpathTV.EncSliceInt64V(*v, e)
	case []bool:
		fastpathTV.EncSliceBoolV(v, e)
	case *[]bool:
		fastpathTV.EncSliceBoolV(*v, e)

	case map[string]interface{}:
		fastpathTV.EncMapStringIntfV(v, e)
	case *map[string]interface{}:
		fastpathTV.EncMapStringIntfV(*v, e)
	case map[string]string:
		fastpathTV.EncMapStringStringV(v, e)
	case *map[string]string:
		fastpathTV.EncMapStringStringV(*v, e)
	case map[string][]byte:
		fastpathTV.EncMapStringBytesV(v, e)
	case *map[string][]byte:
		fastpathTV.EncMapStringBytesV(*v, e)
	case map[string]uint:
		fastpathTV.EncMapStringUintV(v, e)
	case *map[string]uint:
		fastpathTV.EncMapStringUintV(*v, e)
	case map[string]uint8:
		fastpathTV.EncMapStringUint8V(v, e)
	case *map[string]uint8:
		fastpathTV.EncMapStringUint8V(*v, e)
	case map[string]uint64:
		fastpathTV.EncMapStringUint64V(v, e)
	case *map[string]uint64:
		fastpathTV.EncMapStringUint64V(*v, e)
	case map[string]uintptr:
		fastpathTV.EncMapStringUintptrV(v, e)
	case *map[string]uintptr:
		fastpathTV.EncMapStringUintptrV(*v, e)
	case map[string]int:
		fastpathTV.EncMapStringIntV(v, e)
	case *map[string]int:
		fastpathTV.EncMapStringIntV(*v, e)
	case map[string]int64:
		fastpathTV.EncMapStringInt64V(v, e)
	case *map[string]int64:
		fastpathTV.EncMapStringInt64V(*v, e)
	case map[string]float32:
		fastpathTV.EncMapStringFloat32V(v, e)
	case *map[string]float32:
		fastpathTV.EncMapStringFloat32V(*v, e)
	case map[string]float64:
		fastpathTV.EncMapStringFloat64V(v, e)
	case *map[string]float64:
		fastpathTV.EncMapStringFloat64V(*v, e)
	case map[string]bool:
		fastpathTV.EncMapStringBoolV(v, e)
	case *map[string]bool:
		fastpathTV.EncMapStringBoolV(*v, e)
	case map[uint]interface{}:
		fastpathTV.EncMapUintIntfV(v, e)
	case *map[uint]interface{}:
		fastpathTV.EncMapUintIntfV(*v, e)
	case map[uint]string:
		fastpathTV.EncMapUintStringV(v, e)
	case *map[uint]string:
		fastpathTV.EncMapUintStringV(*v, e)
	case map[uint][]byte:
		fastpathTV.EncMapUintBytesV(v, e)
	case *map[uint][]byte:
		fastpathTV.EncMapUintBytesV(*v, e)
	case map[uint]uint:
		fastpathTV.EncMapUintUintV(v, e)
	case *map[uint]uint:
		fastpathTV.EncMapUintUintV(*v, e)
	case map[uint]uint8:
		fastpathTV.EncMapUintUint8V(v, e)
	case *map[uint]uint8:
		fastpathTV.EncMapUintUint8V(*v, e)
	case map[uint]uint64:
		fastpathTV.EncMapUintUint64V(v, e)
	case *map[uint]uint64:
		fastpathTV.EncMapUintUint64V(*v, e)
	case map[uint]uintptr:
		fastpathTV.EncMapUintUintptrV(v, e)
	case *map[uint]uintptr:
		fastpathTV.EncMapUintUintptrV(*v, e)
	case map[uint]int:
		fastpathTV.EncMapUintIntV(v, e)
	case *map[uint]int:
		fastpathTV.EncMapUintIntV(*v, e)
	case map[uint]int64:
		fastpathTV.EncMapUintInt64V(v, e)
	case *map[uint]int64:
		fastpathTV.EncMapUintInt64V(*v, e)
	case map[uint]float32:
		fastpathTV.EncMapUintFloat32V(v, e)
	case *map[uint]float32:
		fastpathTV.EncMapUintFloat32V(*v, e)
	case map[uint]float64:
		fastpathTV.EncMapUintFloat64V(v, e)
	case *map[uint]float64:
		fastpathTV.EncMapUintFloat64V(*v, e)
	case map[uint]bool:
		fastpathTV.EncMapUintBoolV(v, e)
	case *map[uint]bool:
		fastpathTV.EncMapUintBoolV(*v, e)
	case map[uint8]interface{}:
		fastpathTV.EncMapUint8IntfV(v, e)
	case *map[uint8]interface{}:
		fastpathTV.EncMapUint8IntfV(*v, e)
	case map[uint8]string:
		fastpathTV.EncMapUint8StringV(v, e)
	case *map[uint8]string:
		fastpathTV.EncMapUint8StringV(*v, e)
	case map[uint8][]byte:
		fastpathTV.EncMapUint8BytesV(v, e)
	case *map[uint8][]byte:
		fastpathTV.EncMapUint8BytesV(*v, e)
	case map[uint8]uint:
		fastpathTV.EncMapUint8UintV(v, e)
	case *map[uint8]uint:
		fastpathTV.EncMapUint8UintV(*v, e)
	case map[uint8]uint8:
		fastpathTV.EncMapUint8Uint8V(v, e)
	case *map[uint8]uint8:
		fastpathTV.EncMapUint8Uint8V(*v, e)
	case map[uint8]uint64:
		fastpathTV.EncMapUint8Uint64V(v, e)
	case *map[uint8]uint64:
		fastpathTV.EncMapUint8Uint64V(*v, e)
	case map[uint8]uintptr:
		fastpathTV.EncMapUint8UintptrV(v, e)
	case *map[uint8]uintptr:
		fastpathTV.EncMapUint8UintptrV(*v, e)
	case map[uint8]int:
		fastpathTV.EncMapUint8IntV(v, e)
	case *map[uint8]int:
		fastpathTV.EncMapUint8IntV(*v, e)
	case map[uint8]int64:
		fastpathTV.EncMapUint8Int64V(v, e)
	case *map[uint8]int64:
		fastpathTV.EncMapUint8Int64V(*v, e)
	case map[uint8]float32:
		fastpathTV.EncMapUint8Float32V(v, e)
	case *map[uint8]float32:
		fastpathTV.EncMapUint8Float32V(*v, e)
	case map[uint8]float64:
		fastpathTV.EncMapUint8Float64V(v, e)
	case *map[uint8]float64:
		fastpathTV.EncMapUint8Float64V(*v, e)
	case map[uint8]bool:
		fastpathTV.EncMapUint8BoolV(v, e)
	case *map[uint8]bool:
		fastpathTV.EncMapUint8BoolV(*v, e)
	case map[uint64]interface{}:
		fastpathTV.EncMapUint64IntfV(v, e)
	case *map[uint64]interface{}:
		fastpathTV.EncMapUint64IntfV(*v, e)
	case map[uint64]string:
		fastpathTV.EncMapUint64StringV(v, e)
	case *map[uint64]string:
		fastpathTV.EncMapUint64StringV(*v, e)
	case map[uint64][]byte:
		fastpathTV.EncMapUint64BytesV(v, e)
	case *map[uint64][]byte:
		fastpathTV.EncMapUint64BytesV(*v, e)
	case map[uint64]uint:
		fastpathTV.EncMapUint64UintV(v, e)
	case *map[uint64]uint:
		fastpathTV.EncMapUint64UintV(*v, e)
	case map[uint64]uint8:
		fastpathTV.EncMapUint64Uint8V(v, e)
	case *map[uint64]uint8:
		fastpathTV.EncMapUint64Uint8V(*v, e)
	case map[uint64]uint64:
		fastpathTV.EncMapUint64Uint64V(v, e)
	case *map[uint64]uint64:
		fastpathTV.EncMapUint64Uint64V(*v, e)
	case map[uint64]uintptr:
		fastpathTV.EncMapUint64UintptrV(v, e)
	case *map[uint64]uintptr:
		fastpathTV.EncMapUint64UintptrV(*v, e)
	case map[uint64]int:
		fastpathTV.EncMapUint64IntV(v, e)
	case *map[uint64]int:
		fastpathTV.EncMapUint64IntV(*v, e)
	case map[uint64]int64:
		fastpathTV.EncMapUint64Int64V(v, e)
	case *map[uint64]int64:
		fastpathTV.EncMapUint64Int64V(*v, e)
	case map[uint64]float32:
		fastpathTV.EncMapUint64Float32V(v, e)
	case *map[uint64]float32:
		fastpathTV.EncMapUint64Float32V(*v, e)
	case map[uint64]float64:
		fastpathTV.EncMapUint64Float64V(v, e)
	case *map[uint64]float64:
		fastpathTV.EncMapUint64Float64V(*v, e)
	case map[uint64]bool:
		fastpathTV.EncMapUint64BoolV(v, e)
	case *map[uint64]bool:
		fastpathTV.EncMapUint64BoolV(*v, e)
	case map[int]interface{}:
		fastpathTV.EncMapIntIntfV(v, e)
	case *map[int]interface{}:
		fastpathTV.EncMapIntIntfV(*v, e)
	case map[int]string:
		fastpathTV.EncMapIntStringV(v, e)
	case *map[int]string:
		fastpathTV.EncMapIntStringV(*v, e)
	case map[int][]byte:
		fastpathTV.EncMapIntBytesV(v, e)
	case *map[int][]byte:
		fastpathTV.EncMapIntBytesV(*v, e)
	case map[int]uint:
		fastpathTV.EncMapIntUintV(v, e)
	case *map[int]uint:
		fastpathTV.EncMapIntUintV(*v, e)
	case map[int]uint8:
		fastpathTV.EncMapIntUint8V(v, e)
	case *map[int]uint8:
		fastpathTV.EncMapIntUint8V(*v, e)
	case map[int]uint64:
		fastpathTV.EncMapIntUint64V(v, e)
	case *map[int]uint64:
		fastpathTV.EncMapIntUint64V(*v, e)
	case map[int]uintptr:
		fastpathTV.EncMapIntUintptrV(v, e)
	case *map[int]uintptr:
		fastpathTV.EncMapIntUintptrV(*v, e)
	case map[int]int:
		fastpathTV.EncMapIntIntV(v, e)
	case *map[int]int:
		fastpathTV.EncMapIntIntV(*v, e)
	case map[int]int64:
		fastpathTV.EncMapIntInt64V(v, e)
	case *map[int]int64:
		fastpathTV.EncMapIntInt64V(*v, e)
	case map[int]float32:
		fastpathTV.EncMapIntFloat32V(v, e)
	case *map[int]float32:
		fastpathTV.EncMapIntFloat32V(*v, e)
	case map[int]float64:
		fastpathTV.EncMapIntFloat64V(v, e)
	case *map[int]float64:
		fastpathTV.EncMapIntFloat64V(*v, e)
	case map[int]bool:
		fastpathTV.EncMapIntBoolV(v, e)
	case *map[int]bool:
		fastpathTV.EncMapIntBoolV(*v, e)
	case map[int64]interface{}:
		fastpathTV.EncMapInt64IntfV(v, e)
	case *map[int64]interface{}:
		fastpathTV.EncMapInt64IntfV(*v, e)
	case map[int64]string:
		fastpathTV.EncMapInt64StringV(v, e)
	case *map[int64]string:
		fastpathTV.EncMapInt64StringV(*v, e)
	case map[int64][]byte:
		fastpathTV.EncMapInt64BytesV(v, e)
	case *map[int64][]byte:
		fastpathTV.EncMapInt64BytesV(*v, e)
	case map[int64]uint:
		fastpathTV.EncMapInt64UintV(v, e)
	case *map[int64]uint:
		fastpathTV.EncMapInt64UintV(*v, e)
	case map[int64]uint8:
		fastpathTV.EncMapInt64Uint8V(v, e)
	case *map[int64]uint8:
		fastpathTV.EncMapInt64Uint8V(*v, e)
	case map[int64]uint64:
		fastpathTV.EncMapInt64Uint64V(v, e)
	case *map[int64]uint64:
		fastpathTV.EncMapInt64Uint64V(*v, e)
	case map[int64]uintptr:
		fastpathTV.EncMapInt64UintptrV(v, e)
	case *map[int64]uintptr:
		fastpathTV.EncMapInt64UintptrV(*v, e)
	case map[int64]int:
		fastpathTV.EncMapInt64IntV(v, e)
	case *map[int64]int:
		fastpathTV.EncMapInt64IntV(*v, e)
	case map[int64]int64:
		fastpathTV.EncMapInt64Int64V(v, e)
	case *map[int64]int64:
		fastpathTV.EncMapInt64Int64V(*v, e)
	case map[int64]float32:
		fastpathTV.EncMapInt64Float32V(v, e)
	case *map[int64]float32:
		fastpathTV.EncMapInt64Float32V(*v, e)
	case map[int64]float64:
		fastpathTV.EncMapInt64Float64V(v, e)
	case *map[int64]float64:
		fastpathTV.EncMapInt64Float64V(*v, e)
	case map[int64]bool:
		fastpathTV.EncMapInt64BoolV(v, e)
	case *map[int64]bool:
		fastpathTV.EncMapInt64BoolV(*v, e)

	default:
		_ = v // workaround https://github.com/golang/go/issues/12927 seen in go1.4
		return false
	}
	return true
}

// -- -- fast path functions

func (e *Encoder) fastpathEncSliceIntfR(f *codecFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceIntfV(rv2i(rv).([]interface{}), e)
	} else {
		fastpathTV.EncSliceIntfV(rv2i(rv).([]interface{}), e)
	}
}
func (fastpathT) EncSliceIntfV(v []interface{}, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.encode(v[j])
	}
	e.arrayEnd()
}
func (fastpathT) EncAsMapSliceIntfV(v []interface{}, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
	} else if len(v)%2 == 1 {
		e.errorf(fastpathMapBySliceErrMsg, len(v))
	} else {
		e.mapStart(len(v) / 2)
		for j := range v {
			if j%2 == 0 {
				e.mapElemKey()
			} else {
				e.mapElemValue()
			}
			e.encode(v[j])
		}
		e.mapEnd()
	}
}

func (e *Encoder) fastpathEncSliceStringR(f *codecFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceStringV(rv2i(rv).([]string), e)
	} else {
		fastpathTV.EncSliceStringV(rv2i(rv).([]string), e)
	}
}
func (fastpathT) EncSliceStringV(v []string, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		if e.h.StringToRaw {
			e.e.EncodeStringBytesRaw(bytesView(v[j]))
		} else {
			e.e.EncodeStringEnc(cUTF8, v[j])
		}
	}
	e.arrayEnd()
}
func (fastpathT) EncAsMapSliceStringV(v []string, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
	} else if len(v)%2 == 1 {
		e.errorf(fastpathMapBySliceErrMsg, len(v))
	} else {
		e.mapStart(len(v) / 2)
		for j := range v {
			if j%2 == 0 {
				e.mapElemKey()
			} else {
				e.mapElemValue()
			}
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(v[j]))
			} else {
				e.e.EncodeStringEnc(cUTF8, v[j])
			}
		}
		e.mapEnd()
	}
}

func (e *Encoder) fastpathEncSliceBytesR(f *codecFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceBytesV(rv2i(rv).([][]byte), e)
	} else {
		fastpathTV.EncSliceBytesV(rv2i(rv).([][]byte), e)
	}
}
func (fastpathT) EncSliceBytesV(v [][]byte, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeStringBytesRaw(v[j])
	}
	e.arrayEnd()
}
func (fastpathT) EncAsMapSliceBytesV(v [][]byte, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
	} else if len(v)%2 == 1 {
		e.errorf(fastpathMapBySliceErrMsg, len(v))
	} else {
		e.mapStart(len(v) / 2)
		for j := range v {
			if j%2 == 0 {
				e.mapElemKey()
			} else {
				e.mapElemValue()
			}
			e.e.EncodeStringBytesRaw(v[j])
		}
		e.mapEnd()
	}
}

func (e *Encoder) fastpathEncSliceFloat32R(f *codecFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceFloat32V(rv2i(rv).([]float32), e)
	} else {
		fastpathTV.EncSliceFloat32V(rv2i(rv).([]float32), e)
	}
}
func (fastpathT) EncSliceFloat32V(v []float32, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeFloat32(v[j])
	}
	e.arrayEnd()
}
func (fastpathT) EncAsMapSliceFloat32V(v []float32, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
	} else if len(v)%2 == 1 {
		e.errorf(fastpathMapBySliceErrMsg, len(v))
	} else {
		e.mapStart(len(v) / 2)
		for j := range v {
			if j%2 == 0 {
				e.mapElemKey()
			} else {
				e.mapElemValue()
			}
			e.e.EncodeFloat32(v[j])
		}
		e.mapEnd()
	}
}

func (e *Encoder) fastpathEncSliceFloat64R(f *codecFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceFloat64V(rv2i(rv).([]float64), e)
	} else {
		fastpathTV.EncSliceFloat64V(rv2i(rv).([]float64), e)
	}
}
func (fastpathT) EncSliceFloat64V(v []float64, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeFloat64(v[j])
	}
	e.arrayEnd()
}
func (fastpathT) EncAsMapSliceFloat64V(v []float64, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
	} else if len(v)%2 == 1 {
		e.errorf(fastpathMapBySliceErrMsg, len(v))
	} else {
		e.mapStart(len(v) / 2)
		for j := range v {
			if j%2 == 0 {
				e.mapElemKey()
			} else {
				e.mapElemValue()
			}
			e.e.EncodeFloat64(v[j])
		}
		e.mapEnd()
	}
}

func (e *Encoder) fastpathEncSliceUintR(f *codecFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceUintV(rv2i(rv).([]uint), e)
	} else {
		fastpathTV.EncSliceUintV(rv2i(rv).([]uint), e)
	}
}
func (fastpathT) EncSliceUintV(v []uint, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeUint(uint64(v[j]))
	}
	e.arrayEnd()
}
func (fastpathT) EncAsMapSliceUintV(v []uint, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
	} else if len(v)%2 == 1 {
		e.errorf(fastpathMapBySliceErrMsg, len(v))
	} else {
		e.mapStart(len(v) / 2)
		for j := range v {
			if j%2 == 0 {
				e.mapElemKey()
			} else {
				e.mapElemValue()
			}
			e.e.EncodeUint(uint64(v[j]))
		}
		e.mapEnd()
	}
}

func (e *Encoder) fastpathEncSliceUint8R(f *codecFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceUint8V(rv2i(rv).([]uint8), e)
	} else {
		fastpathTV.EncSliceUint8V(rv2i(rv).([]uint8), e)
	}
}
func (fastpathT) EncSliceUint8V(v []uint8, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeUint(uint64(v[j]))
	}
	e.arrayEnd()
}
func (fastpathT) EncAsMapSliceUint8V(v []uint8, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
	} else if len(v)%2 == 1 {
		e.errorf(fastpathMapBySliceErrMsg, len(v))
	} else {
		e.mapStart(len(v) / 2)
		for j := range v {
			if j%2 == 0 {
				e.mapElemKey()
			} else {
				e.mapElemValue()
			}
			e.e.EncodeUint(uint64(v[j]))
		}
		e.mapEnd()
	}
}

func (e *Encoder) fastpathEncSliceUint16R(f *codecFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceUint16V(rv2i(rv).([]uint16), e)
	} else {
		fastpathTV.EncSliceUint16V(rv2i(rv).([]uint16), e)
	}
}
func (fastpathT) EncSliceUint16V(v []uint16, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeUint(uint64(v[j]))
	}
	e.arrayEnd()
}
func (fastpathT) EncAsMapSliceUint16V(v []uint16, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
	} else if len(v)%2 == 1 {
		e.errorf(fastpathMapBySliceErrMsg, len(v))
	} else {
		e.mapStart(len(v) / 2)
		for j := range v {
			if j%2 == 0 {
				e.mapElemKey()
			} else {
				e.mapElemValue()
			}
			e.e.EncodeUint(uint64(v[j]))
		}
		e.mapEnd()
	}
}

func (e *Encoder) fastpathEncSliceUint32R(f *codecFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceUint32V(rv2i(rv).([]uint32), e)
	} else {
		fastpathTV.EncSliceUint32V(rv2i(rv).([]uint32), e)
	}
}
func (fastpathT) EncSliceUint32V(v []uint32, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeUint(uint64(v[j]))
	}
	e.arrayEnd()
}
func (fastpathT) EncAsMapSliceUint32V(v []uint32, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
	} else if len(v)%2 == 1 {
		e.errorf(fastpathMapBySliceErrMsg, len(v))
	} else {
		e.mapStart(len(v) / 2)
		for j := range v {
			if j%2 == 0 {
				e.mapElemKey()
			} else {
				e.mapElemValue()
			}
			e.e.EncodeUint(uint64(v[j]))
		}
		e.mapEnd()
	}
}

func (e *Encoder) fastpathEncSliceUint64R(f *codecFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceUint64V(rv2i(rv).([]uint64), e)
	} else {
		fastpathTV.EncSliceUint64V(rv2i(rv).([]uint64), e)
	}
}
func (fastpathT) EncSliceUint64V(v []uint64, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeUint(v[j])
	}
	e.arrayEnd()
}
func (fastpathT) EncAsMapSliceUint64V(v []uint64, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
	} else if len(v)%2 == 1 {
		e.errorf(fastpathMapBySliceErrMsg, len(v))
	} else {
		e.mapStart(len(v) / 2)
		for j := range v {
			if j%2 == 0 {
				e.mapElemKey()
			} else {
				e.mapElemValue()
			}
			e.e.EncodeUint(v[j])
		}
		e.mapEnd()
	}
}

func (e *Encoder) fastpathEncSliceUintptrR(f *codecFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceUintptrV(rv2i(rv).([]uintptr), e)
	} else {
		fastpathTV.EncSliceUintptrV(rv2i(rv).([]uintptr), e)
	}
}
func (fastpathT) EncSliceUintptrV(v []uintptr, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.encode(v[j])
	}
	e.arrayEnd()
}
func (fastpathT) EncAsMapSliceUintptrV(v []uintptr, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
	} else if len(v)%2 == 1 {
		e.errorf(fastpathMapBySliceErrMsg, len(v))
	} else {
		e.mapStart(len(v) / 2)
		for j := range v {
			if j%2 == 0 {
				e.mapElemKey()
			} else {
				e.mapElemValue()
			}
			e.encode(v[j])
		}
		e.mapEnd()
	}
}

func (e *Encoder) fastpathEncSliceIntR(f *codecFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceIntV(rv2i(rv).([]int), e)
	} else {
		fastpathTV.EncSliceIntV(rv2i(rv).([]int), e)
	}
}
func (fastpathT) EncSliceIntV(v []int, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeInt(int64(v[j]))
	}
	e.arrayEnd()
}
func (fastpathT) EncAsMapSliceIntV(v []int, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
	} else if len(v)%2 == 1 {
		e.errorf(fastpathMapBySliceErrMsg, len(v))
	} else {
		e.mapStart(len(v) / 2)
		for j := range v {
			if j%2 == 0 {
				e.mapElemKey()
			} else {
				e.mapElemValue()
			}
			e.e.EncodeInt(int64(v[j]))
		}
		e.mapEnd()
	}
}

func (e *Encoder) fastpathEncSliceInt8R(f *codecFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceInt8V(rv2i(rv).([]int8), e)
	} else {
		fastpathTV.EncSliceInt8V(rv2i(rv).([]int8), e)
	}
}
func (fastpathT) EncSliceInt8V(v []int8, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeInt(int64(v[j]))
	}
	e.arrayEnd()
}
func (fastpathT) EncAsMapSliceInt8V(v []int8, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
	} else if len(v)%2 == 1 {
		e.errorf(fastpathMapBySliceErrMsg, len(v))
	} else {
		e.mapStart(len(v) / 2)
		for j := range v {
			if j%2 == 0 {
				e.mapElemKey()
			} else {
				e.mapElemValue()
			}
			e.e.EncodeInt(int64(v[j]))
		}
		e.mapEnd()
	}
}

func (e *Encoder) fastpathEncSliceInt16R(f *codecFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceInt16V(rv2i(rv).([]int16), e)
	} else {
		fastpathTV.EncSliceInt16V(rv2i(rv).([]int16), e)
	}
}
func (fastpathT) EncSliceInt16V(v []int16, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeInt(int64(v[j]))
	}
	e.arrayEnd()
}
func (fastpathT) EncAsMapSliceInt16V(v []int16, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
	} else if len(v)%2 == 1 {
		e.errorf(fastpathMapBySliceErrMsg, len(v))
	} else {
		e.mapStart(len(v) / 2)
		for j := range v {
			if j%2 == 0 {
				e.mapElemKey()
			} else {
				e.mapElemValue()
			}
			e.e.EncodeInt(int64(v[j]))
		}
		e.mapEnd()
	}
}

func (e *Encoder) fastpathEncSliceInt32R(f *codecFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceInt32V(rv2i(rv).([]int32), e)
	} else {
		fastpathTV.EncSliceInt32V(rv2i(rv).([]int32), e)
	}
}
func (fastpathT) EncSliceInt32V(v []int32, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeInt(int64(v[j]))
	}
	e.arrayEnd()
}
func (fastpathT) EncAsMapSliceInt32V(v []int32, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
	} else if len(v)%2 == 1 {
		e.errorf(fastpathMapBySliceErrMsg, len(v))
	} else {
		e.mapStart(len(v) / 2)
		for j := range v {
			if j%2 == 0 {
				e.mapElemKey()
			} else {
				e.mapElemValue()
			}
			e.e.EncodeInt(int64(v[j]))
		}
		e.mapEnd()
	}
}

func (e *Encoder) fastpathEncSliceInt64R(f *codecFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceInt64V(rv2i(rv).([]int64), e)
	} else {
		fastpathTV.EncSliceInt64V(rv2i(rv).([]int64), e)
	}
}
func (fastpathT) EncSliceInt64V(v []int64, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeInt(v[j])
	}
	e.arrayEnd()
}
func (fastpathT) EncAsMapSliceInt64V(v []int64, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
	} else if len(v)%2 == 1 {
		e.errorf(fastpathMapBySliceErrMsg, len(v))
	} else {
		e.mapStart(len(v) / 2)
		for j := range v {
			if j%2 == 0 {
				e.mapElemKey()
			} else {
				e.mapElemValue()
			}
			e.e.EncodeInt(v[j])
		}
		e.mapEnd()
	}
}

func (e *Encoder) fastpathEncSliceBoolR(f *codecFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceBoolV(rv2i(rv).([]bool), e)
	} else {
		fastpathTV.EncSliceBoolV(rv2i(rv).([]bool), e)
	}
}
func (fastpathT) EncSliceBoolV(v []bool, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeBool(v[j])
	}
	e.arrayEnd()
}
func (fastpathT) EncAsMapSliceBoolV(v []bool, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
	} else if len(v)%2 == 1 {
		e.errorf(fastpathMapBySliceErrMsg, len(v))
	} else {
		e.mapStart(len(v) / 2)
		for j := range v {
			if j%2 == 0 {
				e.mapElemKey()
			} else {
				e.mapElemValue()
			}
			e.e.EncodeBool(v[j])
		}
		e.mapEnd()
	}
}

func (e *Encoder) fastpathEncMapStringIntfR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapStringIntfV(rv2i(rv).(map[string]interface{}), e)
}
func (fastpathT) EncMapStringIntfV(v map[string]interface{}, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(stringSlice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(k2))
			} else {
				e.e.EncodeStringEnc(cUTF8, k2)
			}
			e.mapElemValue()
			e.encode(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(k2))
			} else {
				e.e.EncodeStringEnc(cUTF8, k2)
			}
			e.mapElemValue()
			e.encode(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapStringStringR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapStringStringV(rv2i(rv).(map[string]string), e)
}
func (fastpathT) EncMapStringStringV(v map[string]string, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(stringSlice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(k2))
			} else {
				e.e.EncodeStringEnc(cUTF8, k2)
			}
			e.mapElemValue()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(v[k2]))
			} else {
				e.e.EncodeStringEnc(cUTF8, v[k2])
			}
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(k2))
			} else {
				e.e.EncodeStringEnc(cUTF8, k2)
			}
			e.mapElemValue()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(v2))
			} else {
				e.e.EncodeStringEnc(cUTF8, v2)
			}
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapStringBytesR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapStringBytesV(rv2i(rv).(map[string][]byte), e)
}
func (fastpathT) EncMapStringBytesV(v map[string][]byte, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(stringSlice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(k2))
			} else {
				e.e.EncodeStringEnc(cUTF8, k2)
			}
			e.mapElemValue()
			e.e.EncodeStringBytesRaw(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(k2))
			} else {
				e.e.EncodeStringEnc(cUTF8, k2)
			}
			e.mapElemValue()
			e.e.EncodeStringBytesRaw(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapStringUintR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapStringUintV(rv2i(rv).(map[string]uint), e)
}
func (fastpathT) EncMapStringUintV(v map[string]uint, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(stringSlice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(k2))
			} else {
				e.e.EncodeStringEnc(cUTF8, k2)
			}
			e.mapElemValue()
			e.e.EncodeUint(uint64(v[k2]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(k2))
			} else {
				e.e.EncodeStringEnc(cUTF8, k2)
			}
			e.mapElemValue()
			e.e.EncodeUint(uint64(v2))
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapStringUint8R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapStringUint8V(rv2i(rv).(map[string]uint8), e)
}
func (fastpathT) EncMapStringUint8V(v map[string]uint8, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(stringSlice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(k2))
			} else {
				e.e.EncodeStringEnc(cUTF8, k2)
			}
			e.mapElemValue()
			e.e.EncodeUint(uint64(v[k2]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(k2))
			} else {
				e.e.EncodeStringEnc(cUTF8, k2)
			}
			e.mapElemValue()
			e.e.EncodeUint(uint64(v2))
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapStringUint64R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapStringUint64V(rv2i(rv).(map[string]uint64), e)
}
func (fastpathT) EncMapStringUint64V(v map[string]uint64, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(stringSlice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(k2))
			} else {
				e.e.EncodeStringEnc(cUTF8, k2)
			}
			e.mapElemValue()
			e.e.EncodeUint(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(k2))
			} else {
				e.e.EncodeStringEnc(cUTF8, k2)
			}
			e.mapElemValue()
			e.e.EncodeUint(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapStringUintptrR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapStringUintptrV(rv2i(rv).(map[string]uintptr), e)
}
func (fastpathT) EncMapStringUintptrV(v map[string]uintptr, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(stringSlice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(k2))
			} else {
				e.e.EncodeStringEnc(cUTF8, k2)
			}
			e.mapElemValue()
			e.encode(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(k2))
			} else {
				e.e.EncodeStringEnc(cUTF8, k2)
			}
			e.mapElemValue()
			e.encode(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapStringIntR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapStringIntV(rv2i(rv).(map[string]int), e)
}
func (fastpathT) EncMapStringIntV(v map[string]int, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(stringSlice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(k2))
			} else {
				e.e.EncodeStringEnc(cUTF8, k2)
			}
			e.mapElemValue()
			e.e.EncodeInt(int64(v[k2]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(k2))
			} else {
				e.e.EncodeStringEnc(cUTF8, k2)
			}
			e.mapElemValue()
			e.e.EncodeInt(int64(v2))
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapStringInt64R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapStringInt64V(rv2i(rv).(map[string]int64), e)
}
func (fastpathT) EncMapStringInt64V(v map[string]int64, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(stringSlice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(k2))
			} else {
				e.e.EncodeStringEnc(cUTF8, k2)
			}
			e.mapElemValue()
			e.e.EncodeInt(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(k2))
			} else {
				e.e.EncodeStringEnc(cUTF8, k2)
			}
			e.mapElemValue()
			e.e.EncodeInt(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapStringFloat32R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapStringFloat32V(rv2i(rv).(map[string]float32), e)
}
func (fastpathT) EncMapStringFloat32V(v map[string]float32, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(stringSlice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(k2))
			} else {
				e.e.EncodeStringEnc(cUTF8, k2)
			}
			e.mapElemValue()
			e.e.EncodeFloat32(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(k2))
			} else {
				e.e.EncodeStringEnc(cUTF8, k2)
			}
			e.mapElemValue()
			e.e.EncodeFloat32(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapStringFloat64R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapStringFloat64V(rv2i(rv).(map[string]float64), e)
}
func (fastpathT) EncMapStringFloat64V(v map[string]float64, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(stringSlice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(k2))
			} else {
				e.e.EncodeStringEnc(cUTF8, k2)
			}
			e.mapElemValue()
			e.e.EncodeFloat64(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(k2))
			} else {
				e.e.EncodeStringEnc(cUTF8, k2)
			}
			e.mapElemValue()
			e.e.EncodeFloat64(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapStringBoolR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapStringBoolV(rv2i(rv).(map[string]bool), e)
}
func (fastpathT) EncMapStringBoolV(v map[string]bool, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(stringSlice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(k2))
			} else {
				e.e.EncodeStringEnc(cUTF8, k2)
			}
			e.mapElemValue()
			e.e.EncodeBool(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(k2))
			} else {
				e.e.EncodeStringEnc(cUTF8, k2)
			}
			e.mapElemValue()
			e.e.EncodeBool(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUintIntfR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUintIntfV(rv2i(rv).(map[uint]interface{}), e)
}
func (fastpathT) EncMapUintIntfV(v map[uint]interface{}, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(uint(k2)))
			e.mapElemValue()
			e.encode(v[uint(k2)])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.encode(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUintStringR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUintStringV(rv2i(rv).(map[uint]string), e)
}
func (fastpathT) EncMapUintStringV(v map[uint]string, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(uint(k2)))
			e.mapElemValue()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(v[uint(k2)]))
			} else {
				e.e.EncodeStringEnc(cUTF8, v[uint(k2)])
			}
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(v2))
			} else {
				e.e.EncodeStringEnc(cUTF8, v2)
			}
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUintBytesR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUintBytesV(rv2i(rv).(map[uint][]byte), e)
}
func (fastpathT) EncMapUintBytesV(v map[uint][]byte, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(uint(k2)))
			e.mapElemValue()
			e.e.EncodeStringBytesRaw(v[uint(k2)])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeStringBytesRaw(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUintUintR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUintUintV(rv2i(rv).(map[uint]uint), e)
}
func (fastpathT) EncMapUintUintV(v map[uint]uint, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(uint(k2)))
			e.mapElemValue()
			e.e.EncodeUint(uint64(v[uint(k2)]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeUint(uint64(v2))
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUintUint8R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUintUint8V(rv2i(rv).(map[uint]uint8), e)
}
func (fastpathT) EncMapUintUint8V(v map[uint]uint8, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(uint(k2)))
			e.mapElemValue()
			e.e.EncodeUint(uint64(v[uint(k2)]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeUint(uint64(v2))
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUintUint64R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUintUint64V(rv2i(rv).(map[uint]uint64), e)
}
func (fastpathT) EncMapUintUint64V(v map[uint]uint64, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(uint(k2)))
			e.mapElemValue()
			e.e.EncodeUint(v[uint(k2)])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeUint(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUintUintptrR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUintUintptrV(rv2i(rv).(map[uint]uintptr), e)
}
func (fastpathT) EncMapUintUintptrV(v map[uint]uintptr, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(uint(k2)))
			e.mapElemValue()
			e.encode(v[uint(k2)])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.encode(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUintIntR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUintIntV(rv2i(rv).(map[uint]int), e)
}
func (fastpathT) EncMapUintIntV(v map[uint]int, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(uint(k2)))
			e.mapElemValue()
			e.e.EncodeInt(int64(v[uint(k2)]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeInt(int64(v2))
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUintInt64R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUintInt64V(rv2i(rv).(map[uint]int64), e)
}
func (fastpathT) EncMapUintInt64V(v map[uint]int64, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(uint(k2)))
			e.mapElemValue()
			e.e.EncodeInt(v[uint(k2)])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeInt(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUintFloat32R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUintFloat32V(rv2i(rv).(map[uint]float32), e)
}
func (fastpathT) EncMapUintFloat32V(v map[uint]float32, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(uint(k2)))
			e.mapElemValue()
			e.e.EncodeFloat32(v[uint(k2)])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeFloat32(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUintFloat64R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUintFloat64V(rv2i(rv).(map[uint]float64), e)
}
func (fastpathT) EncMapUintFloat64V(v map[uint]float64, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(uint(k2)))
			e.mapElemValue()
			e.e.EncodeFloat64(v[uint(k2)])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeFloat64(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUintBoolR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUintBoolV(rv2i(rv).(map[uint]bool), e)
}
func (fastpathT) EncMapUintBoolV(v map[uint]bool, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(uint(k2)))
			e.mapElemValue()
			e.e.EncodeBool(v[uint(k2)])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeBool(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUint8IntfR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUint8IntfV(rv2i(rv).(map[uint8]interface{}), e)
}
func (fastpathT) EncMapUint8IntfV(v map[uint8]interface{}, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(uint8(k2)))
			e.mapElemValue()
			e.encode(v[uint8(k2)])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.encode(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUint8StringR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUint8StringV(rv2i(rv).(map[uint8]string), e)
}
func (fastpathT) EncMapUint8StringV(v map[uint8]string, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(uint8(k2)))
			e.mapElemValue()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(v[uint8(k2)]))
			} else {
				e.e.EncodeStringEnc(cUTF8, v[uint8(k2)])
			}
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(v2))
			} else {
				e.e.EncodeStringEnc(cUTF8, v2)
			}
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUint8BytesR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUint8BytesV(rv2i(rv).(map[uint8][]byte), e)
}
func (fastpathT) EncMapUint8BytesV(v map[uint8][]byte, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(uint8(k2)))
			e.mapElemValue()
			e.e.EncodeStringBytesRaw(v[uint8(k2)])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeStringBytesRaw(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUint8UintR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUint8UintV(rv2i(rv).(map[uint8]uint), e)
}
func (fastpathT) EncMapUint8UintV(v map[uint8]uint, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(uint8(k2)))
			e.mapElemValue()
			e.e.EncodeUint(uint64(v[uint8(k2)]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeUint(uint64(v2))
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUint8Uint8R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUint8Uint8V(rv2i(rv).(map[uint8]uint8), e)
}
func (fastpathT) EncMapUint8Uint8V(v map[uint8]uint8, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(uint8(k2)))
			e.mapElemValue()
			e.e.EncodeUint(uint64(v[uint8(k2)]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeUint(uint64(v2))
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUint8Uint64R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUint8Uint64V(rv2i(rv).(map[uint8]uint64), e)
}
func (fastpathT) EncMapUint8Uint64V(v map[uint8]uint64, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(uint8(k2)))
			e.mapElemValue()
			e.e.EncodeUint(v[uint8(k2)])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeUint(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUint8UintptrR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUint8UintptrV(rv2i(rv).(map[uint8]uintptr), e)
}
func (fastpathT) EncMapUint8UintptrV(v map[uint8]uintptr, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(uint8(k2)))
			e.mapElemValue()
			e.encode(v[uint8(k2)])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.encode(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUint8IntR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUint8IntV(rv2i(rv).(map[uint8]int), e)
}
func (fastpathT) EncMapUint8IntV(v map[uint8]int, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(uint8(k2)))
			e.mapElemValue()
			e.e.EncodeInt(int64(v[uint8(k2)]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeInt(int64(v2))
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUint8Int64R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUint8Int64V(rv2i(rv).(map[uint8]int64), e)
}
func (fastpathT) EncMapUint8Int64V(v map[uint8]int64, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(uint8(k2)))
			e.mapElemValue()
			e.e.EncodeInt(v[uint8(k2)])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeInt(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUint8Float32R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUint8Float32V(rv2i(rv).(map[uint8]float32), e)
}
func (fastpathT) EncMapUint8Float32V(v map[uint8]float32, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(uint8(k2)))
			e.mapElemValue()
			e.e.EncodeFloat32(v[uint8(k2)])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeFloat32(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUint8Float64R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUint8Float64V(rv2i(rv).(map[uint8]float64), e)
}
func (fastpathT) EncMapUint8Float64V(v map[uint8]float64, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(uint8(k2)))
			e.mapElemValue()
			e.e.EncodeFloat64(v[uint8(k2)])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeFloat64(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUint8BoolR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUint8BoolV(rv2i(rv).(map[uint8]bool), e)
}
func (fastpathT) EncMapUint8BoolV(v map[uint8]bool, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(uint8(k2)))
			e.mapElemValue()
			e.e.EncodeBool(v[uint8(k2)])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeBool(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUint64IntfR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUint64IntfV(rv2i(rv).(map[uint64]interface{}), e)
}
func (fastpathT) EncMapUint64IntfV(v map[uint64]interface{}, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.encode(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.encode(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUint64StringR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUint64StringV(rv2i(rv).(map[uint64]string), e)
}
func (fastpathT) EncMapUint64StringV(v map[uint64]string, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(v[k2]))
			} else {
				e.e.EncodeStringEnc(cUTF8, v[k2])
			}
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(v2))
			} else {
				e.e.EncodeStringEnc(cUTF8, v2)
			}
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUint64BytesR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUint64BytesV(rv2i(rv).(map[uint64][]byte), e)
}
func (fastpathT) EncMapUint64BytesV(v map[uint64][]byte, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeStringBytesRaw(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeStringBytesRaw(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUint64UintR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUint64UintV(rv2i(rv).(map[uint64]uint), e)
}
func (fastpathT) EncMapUint64UintV(v map[uint64]uint, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeUint(uint64(v[k2]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeUint(uint64(v2))
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUint64Uint8R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUint64Uint8V(rv2i(rv).(map[uint64]uint8), e)
}
func (fastpathT) EncMapUint64Uint8V(v map[uint64]uint8, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeUint(uint64(v[k2]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeUint(uint64(v2))
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUint64Uint64R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUint64Uint64V(rv2i(rv).(map[uint64]uint64), e)
}
func (fastpathT) EncMapUint64Uint64V(v map[uint64]uint64, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeUint(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeUint(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUint64UintptrR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUint64UintptrV(rv2i(rv).(map[uint64]uintptr), e)
}
func (fastpathT) EncMapUint64UintptrV(v map[uint64]uintptr, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.encode(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.encode(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUint64IntR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUint64IntV(rv2i(rv).(map[uint64]int), e)
}
func (fastpathT) EncMapUint64IntV(v map[uint64]int, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeInt(int64(v[k2]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeInt(int64(v2))
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUint64Int64R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUint64Int64V(rv2i(rv).(map[uint64]int64), e)
}
func (fastpathT) EncMapUint64Int64V(v map[uint64]int64, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeInt(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeInt(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUint64Float32R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUint64Float32V(rv2i(rv).(map[uint64]float32), e)
}
func (fastpathT) EncMapUint64Float32V(v map[uint64]float32, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeFloat32(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeFloat32(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUint64Float64R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUint64Float64V(rv2i(rv).(map[uint64]float64), e)
}
func (fastpathT) EncMapUint64Float64V(v map[uint64]float64, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeFloat64(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeFloat64(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapUint64BoolR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapUint64BoolV(rv2i(rv).(map[uint64]bool), e)
}
func (fastpathT) EncMapUint64BoolV(v map[uint64]bool, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(uint64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeBool(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeBool(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapIntIntfR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapIntIntfV(rv2i(rv).(map[int]interface{}), e)
}
func (fastpathT) EncMapIntIntfV(v map[int]interface{}, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int64, len(v))
		var i uint
		for k := range v {
			v2[i] = int64(k)
			i++
		}
		sort.Sort(int64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(int(k2)))
			e.mapElemValue()
			e.encode(v[int(k2)])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.encode(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapIntStringR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapIntStringV(rv2i(rv).(map[int]string), e)
}
func (fastpathT) EncMapIntStringV(v map[int]string, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int64, len(v))
		var i uint
		for k := range v {
			v2[i] = int64(k)
			i++
		}
		sort.Sort(int64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(int(k2)))
			e.mapElemValue()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(v[int(k2)]))
			} else {
				e.e.EncodeStringEnc(cUTF8, v[int(k2)])
			}
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(v2))
			} else {
				e.e.EncodeStringEnc(cUTF8, v2)
			}
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapIntBytesR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapIntBytesV(rv2i(rv).(map[int][]byte), e)
}
func (fastpathT) EncMapIntBytesV(v map[int][]byte, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int64, len(v))
		var i uint
		for k := range v {
			v2[i] = int64(k)
			i++
		}
		sort.Sort(int64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(int(k2)))
			e.mapElemValue()
			e.e.EncodeStringBytesRaw(v[int(k2)])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeStringBytesRaw(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapIntUintR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapIntUintV(rv2i(rv).(map[int]uint), e)
}
func (fastpathT) EncMapIntUintV(v map[int]uint, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int64, len(v))
		var i uint
		for k := range v {
			v2[i] = int64(k)
			i++
		}
		sort.Sort(int64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(int(k2)))
			e.mapElemValue()
			e.e.EncodeUint(uint64(v[int(k2)]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeUint(uint64(v2))
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapIntUint8R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapIntUint8V(rv2i(rv).(map[int]uint8), e)
}
func (fastpathT) EncMapIntUint8V(v map[int]uint8, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int64, len(v))
		var i uint
		for k := range v {
			v2[i] = int64(k)
			i++
		}
		sort.Sort(int64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(int(k2)))
			e.mapElemValue()
			e.e.EncodeUint(uint64(v[int(k2)]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeUint(uint64(v2))
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapIntUint64R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapIntUint64V(rv2i(rv).(map[int]uint64), e)
}
func (fastpathT) EncMapIntUint64V(v map[int]uint64, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int64, len(v))
		var i uint
		for k := range v {
			v2[i] = int64(k)
			i++
		}
		sort.Sort(int64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(int(k2)))
			e.mapElemValue()
			e.e.EncodeUint(v[int(k2)])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeUint(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapIntUintptrR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapIntUintptrV(rv2i(rv).(map[int]uintptr), e)
}
func (fastpathT) EncMapIntUintptrV(v map[int]uintptr, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int64, len(v))
		var i uint
		for k := range v {
			v2[i] = int64(k)
			i++
		}
		sort.Sort(int64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(int(k2)))
			e.mapElemValue()
			e.encode(v[int(k2)])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.encode(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapIntIntR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapIntIntV(rv2i(rv).(map[int]int), e)
}
func (fastpathT) EncMapIntIntV(v map[int]int, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int64, len(v))
		var i uint
		for k := range v {
			v2[i] = int64(k)
			i++
		}
		sort.Sort(int64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(int(k2)))
			e.mapElemValue()
			e.e.EncodeInt(int64(v[int(k2)]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeInt(int64(v2))
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapIntInt64R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapIntInt64V(rv2i(rv).(map[int]int64), e)
}
func (fastpathT) EncMapIntInt64V(v map[int]int64, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int64, len(v))
		var i uint
		for k := range v {
			v2[i] = int64(k)
			i++
		}
		sort.Sort(int64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(int(k2)))
			e.mapElemValue()
			e.e.EncodeInt(v[int(k2)])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeInt(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapIntFloat32R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapIntFloat32V(rv2i(rv).(map[int]float32), e)
}
func (fastpathT) EncMapIntFloat32V(v map[int]float32, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int64, len(v))
		var i uint
		for k := range v {
			v2[i] = int64(k)
			i++
		}
		sort.Sort(int64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(int(k2)))
			e.mapElemValue()
			e.e.EncodeFloat32(v[int(k2)])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeFloat32(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapIntFloat64R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapIntFloat64V(rv2i(rv).(map[int]float64), e)
}
func (fastpathT) EncMapIntFloat64V(v map[int]float64, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int64, len(v))
		var i uint
		for k := range v {
			v2[i] = int64(k)
			i++
		}
		sort.Sort(int64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(int(k2)))
			e.mapElemValue()
			e.e.EncodeFloat64(v[int(k2)])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeFloat64(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapIntBoolR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapIntBoolV(rv2i(rv).(map[int]bool), e)
}
func (fastpathT) EncMapIntBoolV(v map[int]bool, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int64, len(v))
		var i uint
		for k := range v {
			v2[i] = int64(k)
			i++
		}
		sort.Sort(int64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(int(k2)))
			e.mapElemValue()
			e.e.EncodeBool(v[int(k2)])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeBool(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapInt64IntfR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapInt64IntfV(rv2i(rv).(map[int64]interface{}), e)
}
func (fastpathT) EncMapInt64IntfV(v map[int64]interface{}, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(int64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(k2)
			e.mapElemValue()
			e.encode(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(k2)
			e.mapElemValue()
			e.encode(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapInt64StringR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapInt64StringV(rv2i(rv).(map[int64]string), e)
}
func (fastpathT) EncMapInt64StringV(v map[int64]string, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(int64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(k2)
			e.mapElemValue()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(v[k2]))
			} else {
				e.e.EncodeStringEnc(cUTF8, v[k2])
			}
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(k2)
			e.mapElemValue()
			if e.h.StringToRaw {
				e.e.EncodeStringBytesRaw(bytesView(v2))
			} else {
				e.e.EncodeStringEnc(cUTF8, v2)
			}
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapInt64BytesR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapInt64BytesV(rv2i(rv).(map[int64][]byte), e)
}
func (fastpathT) EncMapInt64BytesV(v map[int64][]byte, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(int64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(k2)
			e.mapElemValue()
			e.e.EncodeStringBytesRaw(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(k2)
			e.mapElemValue()
			e.e.EncodeStringBytesRaw(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapInt64UintR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapInt64UintV(rv2i(rv).(map[int64]uint), e)
}
func (fastpathT) EncMapInt64UintV(v map[int64]uint, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(int64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(k2)
			e.mapElemValue()
			e.e.EncodeUint(uint64(v[k2]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(k2)
			e.mapElemValue()
			e.e.EncodeUint(uint64(v2))
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapInt64Uint8R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapInt64Uint8V(rv2i(rv).(map[int64]uint8), e)
}
func (fastpathT) EncMapInt64Uint8V(v map[int64]uint8, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(int64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(k2)
			e.mapElemValue()
			e.e.EncodeUint(uint64(v[k2]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(k2)
			e.mapElemValue()
			e.e.EncodeUint(uint64(v2))
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapInt64Uint64R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapInt64Uint64V(rv2i(rv).(map[int64]uint64), e)
}
func (fastpathT) EncMapInt64Uint64V(v map[int64]uint64, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(int64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(k2)
			e.mapElemValue()
			e.e.EncodeUint(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(k2)
			e.mapElemValue()
			e.e.EncodeUint(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapInt64UintptrR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapInt64UintptrV(rv2i(rv).(map[int64]uintptr), e)
}
func (fastpathT) EncMapInt64UintptrV(v map[int64]uintptr, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(int64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(k2)
			e.mapElemValue()
			e.encode(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(k2)
			e.mapElemValue()
			e.encode(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapInt64IntR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapInt64IntV(rv2i(rv).(map[int64]int), e)
}
func (fastpathT) EncMapInt64IntV(v map[int64]int, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(int64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(k2)
			e.mapElemValue()
			e.e.EncodeInt(int64(v[k2]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(k2)
			e.mapElemValue()
			e.e.EncodeInt(int64(v2))
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapInt64Int64R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapInt64Int64V(rv2i(rv).(map[int64]int64), e)
}
func (fastpathT) EncMapInt64Int64V(v map[int64]int64, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(int64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(k2)
			e.mapElemValue()
			e.e.EncodeInt(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(k2)
			e.mapElemValue()
			e.e.EncodeInt(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapInt64Float32R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapInt64Float32V(rv2i(rv).(map[int64]float32), e)
}
func (fastpathT) EncMapInt64Float32V(v map[int64]float32, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(int64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(k2)
			e.mapElemValue()
			e.e.EncodeFloat32(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(k2)
			e.mapElemValue()
			e.e.EncodeFloat32(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapInt64Float64R(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapInt64Float64V(rv2i(rv).(map[int64]float64), e)
}
func (fastpathT) EncMapInt64Float64V(v map[int64]float64, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(int64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(k2)
			e.mapElemValue()
			e.e.EncodeFloat64(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(k2)
			e.mapElemValue()
			e.e.EncodeFloat64(v2)
		}
	}
	e.mapEnd()
}

func (e *Encoder) fastpathEncMapInt64BoolR(f *codecFnInfo, rv reflect.Value) {
	fastpathTV.EncMapInt64BoolV(rv2i(rv).(map[int64]bool), e)
}
func (fastpathT) EncMapInt64BoolV(v map[int64]bool, e *Encoder) {
	if v == nil {
		e.e.EncodeNil()
		return
	}
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		sort.Sort(int64Slice(v2))
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(k2)
			e.mapElemValue()
			e.e.EncodeBool(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(k2)
			e.mapElemValue()
			e.e.EncodeBool(v2)
		}
	}
	e.mapEnd()
}

// -- decode

// -- -- fast path type switch
func fastpathDecodeTypeSwitch(iv interface{}, d *Decoder) bool {
	var changed bool
	switch v := iv.(type) {

	case []interface{}:
		var v2 []interface{}
		v2, changed = fastpathTV.DecSliceIntfV(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]interface{}:
		var v2 []interface{}
		v2, changed = fastpathTV.DecSliceIntfV(*v, true, d)
		if changed {
			*v = v2
		}
	case []string:
		var v2 []string
		v2, changed = fastpathTV.DecSliceStringV(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]string:
		var v2 []string
		v2, changed = fastpathTV.DecSliceStringV(*v, true, d)
		if changed {
			*v = v2
		}
	case [][]byte:
		var v2 [][]byte
		v2, changed = fastpathTV.DecSliceBytesV(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[][]byte:
		var v2 [][]byte
		v2, changed = fastpathTV.DecSliceBytesV(*v, true, d)
		if changed {
			*v = v2
		}
	case []float32:
		var v2 []float32
		v2, changed = fastpathTV.DecSliceFloat32V(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]float32:
		var v2 []float32
		v2, changed = fastpathTV.DecSliceFloat32V(*v, true, d)
		if changed {
			*v = v2
		}
	case []float64:
		var v2 []float64
		v2, changed = fastpathTV.DecSliceFloat64V(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]float64:
		var v2 []float64
		v2, changed = fastpathTV.DecSliceFloat64V(*v, true, d)
		if changed {
			*v = v2
		}
	case []uint:
		var v2 []uint
		v2, changed = fastpathTV.DecSliceUintV(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]uint:
		var v2 []uint
		v2, changed = fastpathTV.DecSliceUintV(*v, true, d)
		if changed {
			*v = v2
		}
	case []uint16:
		var v2 []uint16
		v2, changed = fastpathTV.DecSliceUint16V(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]uint16:
		var v2 []uint16
		v2, changed = fastpathTV.DecSliceUint16V(*v, true, d)
		if changed {
			*v = v2
		}
	case []uint32:
		var v2 []uint32
		v2, changed = fastpathTV.DecSliceUint32V(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]uint32:
		var v2 []uint32
		v2, changed = fastpathTV.DecSliceUint32V(*v, true, d)
		if changed {
			*v = v2
		}
	case []uint64:
		var v2 []uint64
		v2, changed = fastpathTV.DecSliceUint64V(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]uint64:
		var v2 []uint64
		v2, changed = fastpathTV.DecSliceUint64V(*v, true, d)
		if changed {
			*v = v2
		}
	case []uintptr:
		var v2 []uintptr
		v2, changed = fastpathTV.DecSliceUintptrV(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]uintptr:
		var v2 []uintptr
		v2, changed = fastpathTV.DecSliceUintptrV(*v, true, d)
		if changed {
			*v = v2
		}
	case []int:
		var v2 []int
		v2, changed = fastpathTV.DecSliceIntV(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]int:
		var v2 []int
		v2, changed = fastpathTV.DecSliceIntV(*v, true, d)
		if changed {
			*v = v2
		}
	case []int8:
		var v2 []int8
		v2, changed = fastpathTV.DecSliceInt8V(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]int8:
		var v2 []int8
		v2, changed = fastpathTV.DecSliceInt8V(*v, true, d)
		if changed {
			*v = v2
		}
	case []int16:
		var v2 []int16
		v2, changed = fastpathTV.DecSliceInt16V(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]int16:
		var v2 []int16
		v2, changed = fastpathTV.DecSliceInt16V(*v, true, d)
		if changed {
			*v = v2
		}
	case []int32:
		var v2 []int32
		v2, changed = fastpathTV.DecSliceInt32V(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]int32:
		var v2 []int32
		v2, changed = fastpathTV.DecSliceInt32V(*v, true, d)
		if changed {
			*v = v2
		}
	case []int64:
		var v2 []int64
		v2, changed = fastpathTV.DecSliceInt64V(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]int64:
		var v2 []int64
		v2, changed = fastpathTV.DecSliceInt64V(*v, true, d)
		if changed {
			*v = v2
		}
	case []bool:
		var v2 []bool
		v2, changed = fastpathTV.DecSliceBoolV(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]bool:
		var v2 []bool
		v2, changed = fastpathTV.DecSliceBoolV(*v, true, d)
		if changed {
			*v = v2
		}

	case map[string]interface{}:
		fastpathTV.DecMapStringIntfV(v, false, d)
	case *map[string]interface{}:
		var v2 map[string]interface{}
		v2, changed = fastpathTV.DecMapStringIntfV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[string]string:
		fastpathTV.DecMapStringStringV(v, false, d)
	case *map[string]string:
		var v2 map[string]string
		v2, changed = fastpathTV.DecMapStringStringV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[string][]byte:
		fastpathTV.DecMapStringBytesV(v, false, d)
	case *map[string][]byte:
		var v2 map[string][]byte
		v2, changed = fastpathTV.DecMapStringBytesV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[string]uint:
		fastpathTV.DecMapStringUintV(v, false, d)
	case *map[string]uint:
		var v2 map[string]uint
		v2, changed = fastpathTV.DecMapStringUintV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[string]uint8:
		fastpathTV.DecMapStringUint8V(v, false, d)
	case *map[string]uint8:
		var v2 map[string]uint8
		v2, changed = fastpathTV.DecMapStringUint8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[string]uint64:
		fastpathTV.DecMapStringUint64V(v, false, d)
	case *map[string]uint64:
		var v2 map[string]uint64
		v2, changed = fastpathTV.DecMapStringUint64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[string]uintptr:
		fastpathTV.DecMapStringUintptrV(v, false, d)
	case *map[string]uintptr:
		var v2 map[string]uintptr
		v2, changed = fastpathTV.DecMapStringUintptrV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[string]int:
		fastpathTV.DecMapStringIntV(v, false, d)
	case *map[string]int:
		var v2 map[string]int
		v2, changed = fastpathTV.DecMapStringIntV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[string]int64:
		fastpathTV.DecMapStringInt64V(v, false, d)
	case *map[string]int64:
		var v2 map[string]int64
		v2, changed = fastpathTV.DecMapStringInt64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[string]float32:
		fastpathTV.DecMapStringFloat32V(v, false, d)
	case *map[string]float32:
		var v2 map[string]float32
		v2, changed = fastpathTV.DecMapStringFloat32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[string]float64:
		fastpathTV.DecMapStringFloat64V(v, false, d)
	case *map[string]float64:
		var v2 map[string]float64
		v2, changed = fastpathTV.DecMapStringFloat64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[string]bool:
		fastpathTV.DecMapStringBoolV(v, false, d)
	case *map[string]bool:
		var v2 map[string]bool
		v2, changed = fastpathTV.DecMapStringBoolV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint]interface{}:
		fastpathTV.DecMapUintIntfV(v, false, d)
	case *map[uint]interface{}:
		var v2 map[uint]interface{}
		v2, changed = fastpathTV.DecMapUintIntfV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint]string:
		fastpathTV.DecMapUintStringV(v, false, d)
	case *map[uint]string:
		var v2 map[uint]string
		v2, changed = fastpathTV.DecMapUintStringV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint][]byte:
		fastpathTV.DecMapUintBytesV(v, false, d)
	case *map[uint][]byte:
		var v2 map[uint][]byte
		v2, changed = fastpathTV.DecMapUintBytesV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint]uint:
		fastpathTV.DecMapUintUintV(v, false, d)
	case *map[uint]uint:
		var v2 map[uint]uint
		v2, changed = fastpathTV.DecMapUintUintV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint]uint8:
		fastpathTV.DecMapUintUint8V(v, false, d)
	case *map[uint]uint8:
		var v2 map[uint]uint8
		v2, changed = fastpathTV.DecMapUintUint8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint]uint64:
		fastpathTV.DecMapUintUint64V(v, false, d)
	case *map[uint]uint64:
		var v2 map[uint]uint64
		v2, changed = fastpathTV.DecMapUintUint64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint]uintptr:
		fastpathTV.DecMapUintUintptrV(v, false, d)
	case *map[uint]uintptr:
		var v2 map[uint]uintptr
		v2, changed = fastpathTV.DecMapUintUintptrV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint]int:
		fastpathTV.DecMapUintIntV(v, false, d)
	case *map[uint]int:
		var v2 map[uint]int
		v2, changed = fastpathTV.DecMapUintIntV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint]int64:
		fastpathTV.DecMapUintInt64V(v, false, d)
	case *map[uint]int64:
		var v2 map[uint]int64
		v2, changed = fastpathTV.DecMapUintInt64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint]float32:
		fastpathTV.DecMapUintFloat32V(v, false, d)
	case *map[uint]float32:
		var v2 map[uint]float32
		v2, changed = fastpathTV.DecMapUintFloat32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint]float64:
		fastpathTV.DecMapUintFloat64V(v, false, d)
	case *map[uint]float64:
		var v2 map[uint]float64
		v2, changed = fastpathTV.DecMapUintFloat64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint]bool:
		fastpathTV.DecMapUintBoolV(v, false, d)
	case *map[uint]bool:
		var v2 map[uint]bool
		v2, changed = fastpathTV.DecMapUintBoolV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint8]interface{}:
		fastpathTV.DecMapUint8IntfV(v, false, d)
	case *map[uint8]interface{}:
		var v2 map[uint8]interface{}
		v2, changed = fastpathTV.DecMapUint8IntfV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint8]string:
		fastpathTV.DecMapUint8StringV(v, false, d)
	case *map[uint8]string:
		var v2 map[uint8]string
		v2, changed = fastpathTV.DecMapUint8StringV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint8][]byte:
		fastpathTV.DecMapUint8BytesV(v, false, d)
	case *map[uint8][]byte:
		var v2 map[uint8][]byte
		v2, changed = fastpathTV.DecMapUint8BytesV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint8]uint:
		fastpathTV.DecMapUint8UintV(v, false, d)
	case *map[uint8]uint:
		var v2 map[uint8]uint
		v2, changed = fastpathTV.DecMapUint8UintV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint8]uint8:
		fastpathTV.DecMapUint8Uint8V(v, false, d)
	case *map[uint8]uint8:
		var v2 map[uint8]uint8
		v2, changed = fastpathTV.DecMapUint8Uint8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint8]uint64:
		fastpathTV.DecMapUint8Uint64V(v, false, d)
	case *map[uint8]uint64:
		var v2 map[uint8]uint64
		v2, changed = fastpathTV.DecMapUint8Uint64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint8]uintptr:
		fastpathTV.DecMapUint8UintptrV(v, false, d)
	case *map[uint8]uintptr:
		var v2 map[uint8]uintptr
		v2, changed = fastpathTV.DecMapUint8UintptrV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint8]int:
		fastpathTV.DecMapUint8IntV(v, false, d)
	case *map[uint8]int:
		var v2 map[uint8]int
		v2, changed = fastpathTV.DecMapUint8IntV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint8]int64:
		fastpathTV.DecMapUint8Int64V(v, false, d)
	case *map[uint8]int64:
		var v2 map[uint8]int64
		v2, changed = fastpathTV.DecMapUint8Int64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint8]float32:
		fastpathTV.DecMapUint8Float32V(v, false, d)
	case *map[uint8]float32:
		var v2 map[uint8]float32
		v2, changed = fastpathTV.DecMapUint8Float32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint8]float64:
		fastpathTV.DecMapUint8Float64V(v, false, d)
	case *map[uint8]float64:
		var v2 map[uint8]float64
		v2, changed = fastpathTV.DecMapUint8Float64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint8]bool:
		fastpathTV.DecMapUint8BoolV(v, false, d)
	case *map[uint8]bool:
		var v2 map[uint8]bool
		v2, changed = fastpathTV.DecMapUint8BoolV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint64]interface{}:
		fastpathTV.DecMapUint64IntfV(v, false, d)
	case *map[uint64]interface{}:
		var v2 map[uint64]interface{}
		v2, changed = fastpathTV.DecMapUint64IntfV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint64]string:
		fastpathTV.DecMapUint64StringV(v, false, d)
	case *map[uint64]string:
		var v2 map[uint64]string
		v2, changed = fastpathTV.DecMapUint64StringV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint64][]byte:
		fastpathTV.DecMapUint64BytesV(v, false, d)
	case *map[uint64][]byte:
		var v2 map[uint64][]byte
		v2, changed = fastpathTV.DecMapUint64BytesV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint64]uint:
		fastpathTV.DecMapUint64UintV(v, false, d)
	case *map[uint64]uint:
		var v2 map[uint64]uint
		v2, changed = fastpathTV.DecMapUint64UintV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint64]uint8:
		fastpathTV.DecMapUint64Uint8V(v, false, d)
	case *map[uint64]uint8:
		var v2 map[uint64]uint8
		v2, changed = fastpathTV.DecMapUint64Uint8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint64]uint64:
		fastpathTV.DecMapUint64Uint64V(v, false, d)
	case *map[uint64]uint64:
		var v2 map[uint64]uint64
		v2, changed = fastpathTV.DecMapUint64Uint64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint64]uintptr:
		fastpathTV.DecMapUint64UintptrV(v, false, d)
	case *map[uint64]uintptr:
		var v2 map[uint64]uintptr
		v2, changed = fastpathTV.DecMapUint64UintptrV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint64]int:
		fastpathTV.DecMapUint64IntV(v, false, d)
	case *map[uint64]int:
		var v2 map[uint64]int
		v2, changed = fastpathTV.DecMapUint64IntV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint64]int64:
		fastpathTV.DecMapUint64Int64V(v, false, d)
	case *map[uint64]int64:
		var v2 map[uint64]int64
		v2, changed = fastpathTV.DecMapUint64Int64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint64]float32:
		fastpathTV.DecMapUint64Float32V(v, false, d)
	case *map[uint64]float32:
		var v2 map[uint64]float32
		v2, changed = fastpathTV.DecMapUint64Float32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint64]float64:
		fastpathTV.DecMapUint64Float64V(v, false, d)
	case *map[uint64]float64:
		var v2 map[uint64]float64
		v2, changed = fastpathTV.DecMapUint64Float64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint64]bool:
		fastpathTV.DecMapUint64BoolV(v, false, d)
	case *map[uint64]bool:
		var v2 map[uint64]bool
		v2, changed = fastpathTV.DecMapUint64BoolV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int]interface{}:
		fastpathTV.DecMapIntIntfV(v, false, d)
	case *map[int]interface{}:
		var v2 map[int]interface{}
		v2, changed = fastpathTV.DecMapIntIntfV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int]string:
		fastpathTV.DecMapIntStringV(v, false, d)
	case *map[int]string:
		var v2 map[int]string
		v2, changed = fastpathTV.DecMapIntStringV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int][]byte:
		fastpathTV.DecMapIntBytesV(v, false, d)
	case *map[int][]byte:
		var v2 map[int][]byte
		v2, changed = fastpathTV.DecMapIntBytesV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int]uint:
		fastpathTV.DecMapIntUintV(v, false, d)
	case *map[int]uint:
		var v2 map[int]uint
		v2, changed = fastpathTV.DecMapIntUintV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int]uint8:
		fastpathTV.DecMapIntUint8V(v, false, d)
	case *map[int]uint8:
		var v2 map[int]uint8
		v2, changed = fastpathTV.DecMapIntUint8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int]uint64:
		fastpathTV.DecMapIntUint64V(v, false, d)
	case *map[int]uint64:
		var v2 map[int]uint64
		v2, changed = fastpathTV.DecMapIntUint64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int]uintptr:
		fastpathTV.DecMapIntUintptrV(v, false, d)
	case *map[int]uintptr:
		var v2 map[int]uintptr
		v2, changed = fastpathTV.DecMapIntUintptrV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int]int:
		fastpathTV.DecMapIntIntV(v, false, d)
	case *map[int]int:
		var v2 map[int]int
		v2, changed = fastpathTV.DecMapIntIntV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int]int64:
		fastpathTV.DecMapIntInt64V(v, false, d)
	case *map[int]int64:
		var v2 map[int]int64
		v2, changed = fastpathTV.DecMapIntInt64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int]float32:
		fastpathTV.DecMapIntFloat32V(v, false, d)
	case *map[int]float32:
		var v2 map[int]float32
		v2, changed = fastpathTV.DecMapIntFloat32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int]float64:
		fastpathTV.DecMapIntFloat64V(v, false, d)
	case *map[int]float64:
		var v2 map[int]float64
		v2, changed = fastpathTV.DecMapIntFloat64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int]bool:
		fastpathTV.DecMapIntBoolV(v, false, d)
	case *map[int]bool:
		var v2 map[int]bool
		v2, changed = fastpathTV.DecMapIntBoolV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int64]interface{}:
		fastpathTV.DecMapInt64IntfV(v, false, d)
	case *map[int64]interface{}:
		var v2 map[int64]interface{}
		v2, changed = fastpathTV.DecMapInt64IntfV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int64]string:
		fastpathTV.DecMapInt64StringV(v, false, d)
	case *map[int64]string:
		var v2 map[int64]string
		v2, changed = fastpathTV.DecMapInt64StringV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int64][]byte:
		fastpathTV.DecMapInt64BytesV(v, false, d)
	case *map[int64][]byte:
		var v2 map[int64][]byte
		v2, changed = fastpathTV.DecMapInt64BytesV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int64]uint:
		fastpathTV.DecMapInt64UintV(v, false, d)
	case *map[int64]uint:
		var v2 map[int64]uint
		v2, changed = fastpathTV.DecMapInt64UintV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int64]uint8:
		fastpathTV.DecMapInt64Uint8V(v, false, d)
	case *map[int64]uint8:
		var v2 map[int64]uint8
		v2, changed = fastpathTV.DecMapInt64Uint8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int64]uint64:
		fastpathTV.DecMapInt64Uint64V(v, false, d)
	case *map[int64]uint64:
		var v2 map[int64]uint64
		v2, changed = fastpathTV.DecMapInt64Uint64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int64]uintptr:
		fastpathTV.DecMapInt64UintptrV(v, false, d)
	case *map[int64]uintptr:
		var v2 map[int64]uintptr
		v2, changed = fastpathTV.DecMapInt64UintptrV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int64]int:
		fastpathTV.DecMapInt64IntV(v, false, d)
	case *map[int64]int:
		var v2 map[int64]int
		v2, changed = fastpathTV.DecMapInt64IntV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int64]int64:
		fastpathTV.DecMapInt64Int64V(v, false, d)
	case *map[int64]int64:
		var v2 map[int64]int64
		v2, changed = fastpathTV.DecMapInt64Int64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int64]float32:
		fastpathTV.DecMapInt64Float32V(v, false, d)
	case *map[int64]float32:
		var v2 map[int64]float32
		v2, changed = fastpathTV.DecMapInt64Float32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int64]float64:
		fastpathTV.DecMapInt64Float64V(v, false, d)
	case *map[int64]float64:
		var v2 map[int64]float64
		v2, changed = fastpathTV.DecMapInt64Float64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int64]bool:
		fastpathTV.DecMapInt64BoolV(v, false, d)
	case *map[int64]bool:
		var v2 map[int64]bool
		v2, changed = fastpathTV.DecMapInt64BoolV(*v, true, d)
		if changed {
			*v = v2
		}
	default:
		_ = v // workaround https://github.com/golang/go/issues/12927 seen in go1.4
		return false
	}
	return true
}

func fastpathDecodeSetZeroTypeSwitch(iv interface{}) bool {
	switch v := iv.(type) {

	case *[]interface{}:
		*v = nil
	case *[]string:
		*v = nil
	case *[][]byte:
		*v = nil
	case *[]float32:
		*v = nil
	case *[]float64:
		*v = nil
	case *[]uint:
		*v = nil
	case *[]uint8:
		*v = nil
	case *[]uint16:
		*v = nil
	case *[]uint32:
		*v = nil
	case *[]uint64:
		*v = nil
	case *[]uintptr:
		*v = nil
	case *[]int:
		*v = nil
	case *[]int8:
		*v = nil
	case *[]int16:
		*v = nil
	case *[]int32:
		*v = nil
	case *[]int64:
		*v = nil
	case *[]bool:
		*v = nil

	case *map[string]interface{}:
		*v = nil
	case *map[string]string:
		*v = nil
	case *map[string][]byte:
		*v = nil
	case *map[string]uint:
		*v = nil
	case *map[string]uint8:
		*v = nil
	case *map[string]uint64:
		*v = nil
	case *map[string]uintptr:
		*v = nil
	case *map[string]int:
		*v = nil
	case *map[string]int64:
		*v = nil
	case *map[string]float32:
		*v = nil
	case *map[string]float64:
		*v = nil
	case *map[string]bool:
		*v = nil
	case *map[uint]interface{}:
		*v = nil
	case *map[uint]string:
		*v = nil
	case *map[uint][]byte:
		*v = nil
	case *map[uint]uint:
		*v = nil
	case *map[uint]uint8:
		*v = nil
	case *map[uint]uint64:
		*v = nil
	case *map[uint]uintptr:
		*v = nil
	case *map[uint]int:
		*v = nil
	case *map[uint]int64:
		*v = nil
	case *map[uint]float32:
		*v = nil
	case *map[uint]float64:
		*v = nil
	case *map[uint]bool:
		*v = nil
	case *map[uint8]interface{}:
		*v = nil
	case *map[uint8]string:
		*v = nil
	case *map[uint8][]byte:
		*v = nil
	case *map[uint8]uint:
		*v = nil
	case *map[uint8]uint8:
		*v = nil
	case *map[uint8]uint64:
		*v = nil
	case *map[uint8]uintptr:
		*v = nil
	case *map[uint8]int:
		*v = nil
	case *map[uint8]int64:
		*v = nil
	case *map[uint8]float32:
		*v = nil
	case *map[uint8]float64:
		*v = nil
	case *map[uint8]bool:
		*v = nil
	case *map[uint64]interface{}:
		*v = nil
	case *map[uint64]string:
		*v = nil
	case *map[uint64][]byte:
		*v = nil
	case *map[uint64]uint:
		*v = nil
	case *map[uint64]uint8:
		*v = nil
	case *map[uint64]uint64:
		*v = nil
	case *map[uint64]uintptr:
		*v = nil
	case *map[uint64]int:
		*v = nil
	case *map[uint64]int64:
		*v = nil
	case *map[uint64]float32:
		*v = nil
	case *map[uint64]float64:
		*v = nil
	case *map[uint64]bool:
		*v = nil
	case *map[int]interface{}:
		*v = nil
	case *map[int]string:
		*v = nil
	case *map[int][]byte:
		*v = nil
	case *map[int]uint:
		*v = nil
	case *map[int]uint8:
		*v = nil
	case *map[int]uint64:
		*v = nil
	case *map[int]uintptr:
		*v = nil
	case *map[int]int:
		*v = nil
	case *map[int]int64:
		*v = nil
	case *map[int]float32:
		*v = nil
	case *map[int]float64:
		*v = nil
	case *map[int]bool:
		*v = nil
	case *map[int64]interface{}:
		*v = nil
	case *map[int64]string:
		*v = nil
	case *map[int64][]byte:
		*v = nil
	case *map[int64]uint:
		*v = nil
	case *map[int64]uint8:
		*v = nil
	case *map[int64]uint64:
		*v = nil
	case *map[int64]uintptr:
		*v = nil
	case *map[int64]int:
		*v = nil
	case *map[int64]int64:
		*v = nil
	case *map[int64]float32:
		*v = nil
	case *map[int64]float64:
		*v = nil
	case *map[int64]bool:
		*v = nil
	default:
		_ = v // workaround https://github.com/golang/go/issues/12927 seen in go1.4
		return false
	}
	return true
}

// -- -- fast path functions

func (d *Decoder) fastpathDecSliceIntfR(f *codecFnInfo, rv reflect.Value) {
	if array := f.seq == seqTypeArray; !array && rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*[]interface{})
		if v, changed := fastpathTV.DecSliceIntfV(*vp, !array, d); changed {
			*vp = v
		}
	} else {
		v := rv2i(rv).([]interface{})
		v2, changed := fastpathTV.DecSliceIntfV(v, !array, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	}
}
func (f fastpathT) DecSliceIntfX(vp *[]interface{}, d *Decoder) {
	if v, changed := f.DecSliceIntfV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecSliceIntfV(v []interface{}, canChange bool, d *Decoder) (_ []interface{}, changed bool) {
	slh, containerLenS := d.decSliceHelperStart()
	if containerLenS == 0 {
		if canChange {
			if v == nil {
				v = []interface{}{}
			} else if len(v) != 0 {
				v = v[:0]
			}
			changed = true
		}
		slh.End()
		return v, changed
	}
	hasLen := containerLenS > 0
	var xlen int
	if hasLen && canChange {
		if containerLenS > cap(v) {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 16)
			if xlen <= cap(v) {
				v = v[:uint(xlen)]
			} else {
				v = make([]interface{}, uint(xlen))
			}
			changed = true
		} else if containerLenS != len(v) {
			v = v[:containerLenS]
			changed = true
		}
	}
	var j int
	for j = 0; (hasLen && j < containerLenS) || !(hasLen || d.d.CheckBreak()); j++ {
		if j == 0 && len(v) == 0 && canChange {
			if hasLen {
				xlen = decInferLen(containerLenS, d.h.MaxInitLen, 16)
			} else {
				xlen = 8
			}
			v = make([]interface{}, uint(xlen))
			changed = true
		}
		var decodeIntoBlank bool
		if j >= len(v) {
			if canChange {
				v = append(v, nil)
				changed = true
			} else {
				d.arrayCannotExpand(len(v), j+1)
				decodeIntoBlank = true
			}
		}
		slh.ElemContainerState(j)
		if decodeIntoBlank {
			d.swallow()
		} else if d.d.TryDecodeAsNil() {
			v[uint(j)] = nil
		} else {
			d.decode(&v[uint(j)])
		}
	}
	if canChange {
		if j < len(v) {
			v = v[:uint(j)]
			changed = true
		} else if j == 0 && v == nil {
			v = make([]interface{}, 0)
			changed = true
		}
	}
	slh.End()
	return v, changed
}

func (d *Decoder) fastpathDecSliceStringR(f *codecFnInfo, rv reflect.Value) {
	if array := f.seq == seqTypeArray; !array && rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*[]string)
		if v, changed := fastpathTV.DecSliceStringV(*vp, !array, d); changed {
			*vp = v
		}
	} else {
		v := rv2i(rv).([]string)
		v2, changed := fastpathTV.DecSliceStringV(v, !array, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	}
}
func (f fastpathT) DecSliceStringX(vp *[]string, d *Decoder) {
	if v, changed := f.DecSliceStringV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecSliceStringV(v []string, canChange bool, d *Decoder) (_ []string, changed bool) {
	slh, containerLenS := d.decSliceHelperStart()
	if containerLenS == 0 {
		if canChange {
			if v == nil {
				v = []string{}
			} else if len(v) != 0 {
				v = v[:0]
			}
			changed = true
		}
		slh.End()
		return v, changed
	}
	hasLen := containerLenS > 0
	var xlen int
	if hasLen && canChange {
		if containerLenS > cap(v) {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 16)
			if xlen <= cap(v) {
				v = v[:uint(xlen)]
			} else {
				v = make([]string, uint(xlen))
			}
			changed = true
		} else if containerLenS != len(v) {
			v = v[:containerLenS]
			changed = true
		}
	}
	var j int
	for j = 0; (hasLen && j < containerLenS) || !(hasLen || d.d.CheckBreak()); j++ {
		if j == 0 && len(v) == 0 && canChange {
			if hasLen {
				xlen = decInferLen(containerLenS, d.h.MaxInitLen, 16)
			} else {
				xlen = 8
			}
			v = make([]string, uint(xlen))
			changed = true
		}
		var decodeIntoBlank bool
		if j >= len(v) {
			if canChange {
				v = append(v, "")
				changed = true
			} else {
				d.arrayCannotExpand(len(v), j+1)
				decodeIntoBlank = true
			}
		}
		slh.ElemContainerState(j)
		if decodeIntoBlank {
			d.swallow()
		} else if d.d.TryDecodeAsNil() {
			v[uint(j)] = ""
		} else {
			v[uint(j)] = d.d.DecodeString()
		}
	}
	if canChange {
		if j < len(v) {
			v = v[:uint(j)]
			changed = true
		} else if j == 0 && v == nil {
			v = make([]string, 0)
			changed = true
		}
	}
	slh.End()
	return v, changed
}

func (d *Decoder) fastpathDecSliceBytesR(f *codecFnInfo, rv reflect.Value) {
	if array := f.seq == seqTypeArray; !array && rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*[][]byte)
		if v, changed := fastpathTV.DecSliceBytesV(*vp, !array, d); changed {
			*vp = v
		}
	} else {
		v := rv2i(rv).([][]byte)
		v2, changed := fastpathTV.DecSliceBytesV(v, !array, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	}
}
func (f fastpathT) DecSliceBytesX(vp *[][]byte, d *Decoder) {
	if v, changed := f.DecSliceBytesV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecSliceBytesV(v [][]byte, canChange bool, d *Decoder) (_ [][]byte, changed bool) {
	slh, containerLenS := d.decSliceHelperStart()
	if containerLenS == 0 {
		if canChange {
			if v == nil {
				v = [][]byte{}
			} else if len(v) != 0 {
				v = v[:0]
			}
			changed = true
		}
		slh.End()
		return v, changed
	}
	hasLen := containerLenS > 0
	var xlen int
	if hasLen && canChange {
		if containerLenS > cap(v) {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 24)
			if xlen <= cap(v) {
				v = v[:uint(xlen)]
			} else {
				v = make([][]byte, uint(xlen))
			}
			changed = true
		} else if containerLenS != len(v) {
			v = v[:containerLenS]
			changed = true
		}
	}
	var j int
	for j = 0; (hasLen && j < containerLenS) || !(hasLen || d.d.CheckBreak()); j++ {
		if j == 0 && len(v) == 0 && canChange {
			if hasLen {
				xlen = decInferLen(containerLenS, d.h.MaxInitLen, 24)
			} else {
				xlen = 8
			}
			v = make([][]byte, uint(xlen))
			changed = true
		}
		var decodeIntoBlank bool
		if j >= len(v) {
			if canChange {
				v = append(v, nil)
				changed = true
			} else {
				d.arrayCannotExpand(len(v), j+1)
				decodeIntoBlank = true
			}
		}
		slh.ElemContainerState(j)
		if decodeIntoBlank {
			d.swallow()
		} else if d.d.TryDecodeAsNil() {
			v[uint(j)] = nil
		} else {
			v[uint(j)] = d.d.DecodeBytes(nil, false)
		}
	}
	if canChange {
		if j < len(v) {
			v = v[:uint(j)]
			changed = true
		} else if j == 0 && v == nil {
			v = make([][]byte, 0)
			changed = true
		}
	}
	slh.End()
	return v, changed
}

func (d *Decoder) fastpathDecSliceFloat32R(f *codecFnInfo, rv reflect.Value) {
	if array := f.seq == seqTypeArray; !array && rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*[]float32)
		if v, changed := fastpathTV.DecSliceFloat32V(*vp, !array, d); changed {
			*vp = v
		}
	} else {
		v := rv2i(rv).([]float32)
		v2, changed := fastpathTV.DecSliceFloat32V(v, !array, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	}
}
func (f fastpathT) DecSliceFloat32X(vp *[]float32, d *Decoder) {
	if v, changed := f.DecSliceFloat32V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecSliceFloat32V(v []float32, canChange bool, d *Decoder) (_ []float32, changed bool) {
	slh, containerLenS := d.decSliceHelperStart()
	if containerLenS == 0 {
		if canChange {
			if v == nil {
				v = []float32{}
			} else if len(v) != 0 {
				v = v[:0]
			}
			changed = true
		}
		slh.End()
		return v, changed
	}
	hasLen := containerLenS > 0
	var xlen int
	if hasLen && canChange {
		if containerLenS > cap(v) {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 4)
			if xlen <= cap(v) {
				v = v[:uint(xlen)]
			} else {
				v = make([]float32, uint(xlen))
			}
			changed = true
		} else if containerLenS != len(v) {
			v = v[:containerLenS]
			changed = true
		}
	}
	var j int
	for j = 0; (hasLen && j < containerLenS) || !(hasLen || d.d.CheckBreak()); j++ {
		if j == 0 && len(v) == 0 && canChange {
			if hasLen {
				xlen = decInferLen(containerLenS, d.h.MaxInitLen, 4)
			} else {
				xlen = 8
			}
			v = make([]float32, uint(xlen))
			changed = true
		}
		var decodeIntoBlank bool
		if j >= len(v) {
			if canChange {
				v = append(v, 0)
				changed = true
			} else {
				d.arrayCannotExpand(len(v), j+1)
				decodeIntoBlank = true
			}
		}
		slh.ElemContainerState(j)
		if decodeIntoBlank {
			d.swallow()
		} else if d.d.TryDecodeAsNil() {
			v[uint(j)] = 0
		} else {
			v[uint(j)] = float32(d.decodeFloat32())
		}
	}
	if canChange {
		if j < len(v) {
			v = v[:uint(j)]
			changed = true
		} else if j == 0 && v == nil {
			v = make([]float32, 0)
			changed = true
		}
	}
	slh.End()
	return v, changed
}

func (d *Decoder) fastpathDecSliceFloat64R(f *codecFnInfo, rv reflect.Value) {
	if array := f.seq == seqTypeArray; !array && rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*[]float64)
		if v, changed := fastpathTV.DecSliceFloat64V(*vp, !array, d); changed {
			*vp = v
		}
	} else {
		v := rv2i(rv).([]float64)
		v2, changed := fastpathTV.DecSliceFloat64V(v, !array, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	}
}
func (f fastpathT) DecSliceFloat64X(vp *[]float64, d *Decoder) {
	if v, changed := f.DecSliceFloat64V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecSliceFloat64V(v []float64, canChange bool, d *Decoder) (_ []float64, changed bool) {
	slh, containerLenS := d.decSliceHelperStart()
	if containerLenS == 0 {
		if canChange {
			if v == nil {
				v = []float64{}
			} else if len(v) != 0 {
				v = v[:0]
			}
			changed = true
		}
		slh.End()
		return v, changed
	}
	hasLen := containerLenS > 0
	var xlen int
	if hasLen && canChange {
		if containerLenS > cap(v) {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 8)
			if xlen <= cap(v) {
				v = v[:uint(xlen)]
			} else {
				v = make([]float64, uint(xlen))
			}
			changed = true
		} else if containerLenS != len(v) {
			v = v[:containerLenS]
			changed = true
		}
	}
	var j int
	for j = 0; (hasLen && j < containerLenS) || !(hasLen || d.d.CheckBreak()); j++ {
		if j == 0 && len(v) == 0 && canChange {
			if hasLen {
				xlen = decInferLen(containerLenS, d.h.MaxInitLen, 8)
			} else {
				xlen = 8
			}
			v = make([]float64, uint(xlen))
			changed = true
		}
		var decodeIntoBlank bool
		if j >= len(v) {
			if canChange {
				v = append(v, 0)
				changed = true
			} else {
				d.arrayCannotExpand(len(v), j+1)
				decodeIntoBlank = true
			}
		}
		slh.ElemContainerState(j)
		if decodeIntoBlank {
			d.swallow()
		} else if d.d.TryDecodeAsNil() {
			v[uint(j)] = 0
		} else {
			v[uint(j)] = d.d.DecodeFloat64()
		}
	}
	if canChange {
		if j < len(v) {
			v = v[:uint(j)]
			changed = true
		} else if j == 0 && v == nil {
			v = make([]float64, 0)
			changed = true
		}
	}
	slh.End()
	return v, changed
}

func (d *Decoder) fastpathDecSliceUintR(f *codecFnInfo, rv reflect.Value) {
	if array := f.seq == seqTypeArray; !array && rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*[]uint)
		if v, changed := fastpathTV.DecSliceUintV(*vp, !array, d); changed {
			*vp = v
		}
	} else {
		v := rv2i(rv).([]uint)
		v2, changed := fastpathTV.DecSliceUintV(v, !array, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	}
}
func (f fastpathT) DecSliceUintX(vp *[]uint, d *Decoder) {
	if v, changed := f.DecSliceUintV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecSliceUintV(v []uint, canChange bool, d *Decoder) (_ []uint, changed bool) {
	slh, containerLenS := d.decSliceHelperStart()
	if containerLenS == 0 {
		if canChange {
			if v == nil {
				v = []uint{}
			} else if len(v) != 0 {
				v = v[:0]
			}
			changed = true
		}
		slh.End()
		return v, changed
	}
	hasLen := containerLenS > 0
	var xlen int
	if hasLen && canChange {
		if containerLenS > cap(v) {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 8)
			if xlen <= cap(v) {
				v = v[:uint(xlen)]
			} else {
				v = make([]uint, uint(xlen))
			}
			changed = true
		} else if containerLenS != len(v) {
			v = v[:containerLenS]
			changed = true
		}
	}
	var j int
	for j = 0; (hasLen && j < containerLenS) || !(hasLen || d.d.CheckBreak()); j++ {
		if j == 0 && len(v) == 0 && canChange {
			if hasLen {
				xlen = decInferLen(containerLenS, d.h.MaxInitLen, 8)
			} else {
				xlen = 8
			}
			v = make([]uint, uint(xlen))
			changed = true
		}
		var decodeIntoBlank bool
		if j >= len(v) {
			if canChange {
				v = append(v, 0)
				changed = true
			} else {
				d.arrayCannotExpand(len(v), j+1)
				decodeIntoBlank = true
			}
		}
		slh.ElemContainerState(j)
		if decodeIntoBlank {
			d.swallow()
		} else if d.d.TryDecodeAsNil() {
			v[uint(j)] = 0
		} else {
			v[uint(j)] = uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
		}
	}
	if canChange {
		if j < len(v) {
			v = v[:uint(j)]
			changed = true
		} else if j == 0 && v == nil {
			v = make([]uint, 0)
			changed = true
		}
	}
	slh.End()
	return v, changed
}

func (d *Decoder) fastpathDecSliceUint8R(f *codecFnInfo, rv reflect.Value) {
	if array := f.seq == seqTypeArray; !array && rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*[]uint8)
		if v, changed := fastpathTV.DecSliceUint8V(*vp, !array, d); changed {
			*vp = v
		}
	} else {
		v := rv2i(rv).([]uint8)
		v2, changed := fastpathTV.DecSliceUint8V(v, !array, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	}
}
func (f fastpathT) DecSliceUint8X(vp *[]uint8, d *Decoder) {
	if v, changed := f.DecSliceUint8V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecSliceUint8V(v []uint8, canChange bool, d *Decoder) (_ []uint8, changed bool) {
	slh, containerLenS := d.decSliceHelperStart()
	if containerLenS == 0 {
		if canChange {
			if v == nil {
				v = []uint8{}
			} else if len(v) != 0 {
				v = v[:0]
			}
			changed = true
		}
		slh.End()
		return v, changed
	}
	hasLen := containerLenS > 0
	var xlen int
	if hasLen && canChange {
		if containerLenS > cap(v) {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 1)
			if xlen <= cap(v) {
				v = v[:uint(xlen)]
			} else {
				v = make([]uint8, uint(xlen))
			}
			changed = true
		} else if containerLenS != len(v) {
			v = v[:containerLenS]
			changed = true
		}
	}
	var j int
	for j = 0; (hasLen && j < containerLenS) || !(hasLen || d.d.CheckBreak()); j++ {
		if j == 0 && len(v) == 0 && canChange {
			if hasLen {
				xlen = decInferLen(containerLenS, d.h.MaxInitLen, 1)
			} else {
				xlen = 8
			}
			v = make([]uint8, uint(xlen))
			changed = true
		}
		var decodeIntoBlank bool
		if j >= len(v) {
			if canChange {
				v = append(v, 0)
				changed = true
			} else {
				d.arrayCannotExpand(len(v), j+1)
				decodeIntoBlank = true
			}
		}
		slh.ElemContainerState(j)
		if decodeIntoBlank {
			d.swallow()
		} else if d.d.TryDecodeAsNil() {
			v[uint(j)] = 0
		} else {
			v[uint(j)] = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		}
	}
	if canChange {
		if j < len(v) {
			v = v[:uint(j)]
			changed = true
		} else if j == 0 && v == nil {
			v = make([]uint8, 0)
			changed = true
		}
	}
	slh.End()
	return v, changed
}

func (d *Decoder) fastpathDecSliceUint16R(f *codecFnInfo, rv reflect.Value) {
	if array := f.seq == seqTypeArray; !array && rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*[]uint16)
		if v, changed := fastpathTV.DecSliceUint16V(*vp, !array, d); changed {
			*vp = v
		}
	} else {
		v := rv2i(rv).([]uint16)
		v2, changed := fastpathTV.DecSliceUint16V(v, !array, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	}
}
func (f fastpathT) DecSliceUint16X(vp *[]uint16, d *Decoder) {
	if v, changed := f.DecSliceUint16V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecSliceUint16V(v []uint16, canChange bool, d *Decoder) (_ []uint16, changed bool) {
	slh, containerLenS := d.decSliceHelperStart()
	if containerLenS == 0 {
		if canChange {
			if v == nil {
				v = []uint16{}
			} else if len(v) != 0 {
				v = v[:0]
			}
			changed = true
		}
		slh.End()
		return v, changed
	}
	hasLen := containerLenS > 0
	var xlen int
	if hasLen && canChange {
		if containerLenS > cap(v) {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 2)
			if xlen <= cap(v) {
				v = v[:uint(xlen)]
			} else {
				v = make([]uint16, uint(xlen))
			}
			changed = true
		} else if containerLenS != len(v) {
			v = v[:containerLenS]
			changed = true
		}
	}
	var j int
	for j = 0; (hasLen && j < containerLenS) || !(hasLen || d.d.CheckBreak()); j++ {
		if j == 0 && len(v) == 0 && canChange {
			if hasLen {
				xlen = decInferLen(containerLenS, d.h.MaxInitLen, 2)
			} else {
				xlen = 8
			}
			v = make([]uint16, uint(xlen))
			changed = true
		}
		var decodeIntoBlank bool
		if j >= len(v) {
			if canChange {
				v = append(v, 0)
				changed = true
			} else {
				d.arrayCannotExpand(len(v), j+1)
				decodeIntoBlank = true
			}
		}
		slh.ElemContainerState(j)
		if decodeIntoBlank {
			d.swallow()
		} else if d.d.TryDecodeAsNil() {
			v[uint(j)] = 0
		} else {
			v[uint(j)] = uint16(chkOvf.UintV(d.d.DecodeUint64(), 16))
		}
	}
	if canChange {
		if j < len(v) {
			v = v[:uint(j)]
			changed = true
		} else if j == 0 && v == nil {
			v = make([]uint16, 0)
			changed = true
		}
	}
	slh.End()
	return v, changed
}

func (d *Decoder) fastpathDecSliceUint32R(f *codecFnInfo, rv reflect.Value) {
	if array := f.seq == seqTypeArray; !array && rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*[]uint32)
		if v, changed := fastpathTV.DecSliceUint32V(*vp, !array, d); changed {
			*vp = v
		}
	} else {
		v := rv2i(rv).([]uint32)
		v2, changed := fastpathTV.DecSliceUint32V(v, !array, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	}
}
func (f fastpathT) DecSliceUint32X(vp *[]uint32, d *Decoder) {
	if v, changed := f.DecSliceUint32V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecSliceUint32V(v []uint32, canChange bool, d *Decoder) (_ []uint32, changed bool) {
	slh, containerLenS := d.decSliceHelperStart()
	if containerLenS == 0 {
		if canChange {
			if v == nil {
				v = []uint32{}
			} else if len(v) != 0 {
				v = v[:0]
			}
			changed = true
		}
		slh.End()
		return v, changed
	}
	hasLen := containerLenS > 0
	var xlen int
	if hasLen && canChange {
		if containerLenS > cap(v) {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 4)
			if xlen <= cap(v) {
				v = v[:uint(xlen)]
			} else {
				v = make([]uint32, uint(xlen))
			}
			changed = true
		} else if containerLenS != len(v) {
			v = v[:containerLenS]
			changed = true
		}
	}
	var j int
	for j = 0; (hasLen && j < containerLenS) || !(hasLen || d.d.CheckBreak()); j++ {
		if j == 0 && len(v) == 0 && canChange {
			if hasLen {
				xlen = decInferLen(containerLenS, d.h.MaxInitLen, 4)
			} else {
				xlen = 8
			}
			v = make([]uint32, uint(xlen))
			changed = true
		}
		var decodeIntoBlank bool
		if j >= len(v) {
			if canChange {
				v = append(v, 0)
				changed = true
			} else {
				d.arrayCannotExpand(len(v), j+1)
				decodeIntoBlank = true
			}
		}
		slh.ElemContainerState(j)
		if decodeIntoBlank {
			d.swallow()
		} else if d.d.TryDecodeAsNil() {
			v[uint(j)] = 0
		} else {
			v[uint(j)] = uint32(chkOvf.UintV(d.d.DecodeUint64(), 32))
		}
	}
	if canChange {
		if j < len(v) {
			v = v[:uint(j)]
			changed = true
		} else if j == 0 && v == nil {
			v = make([]uint32, 0)
			changed = true
		}
	}
	slh.End()
	return v, changed
}

func (d *Decoder) fastpathDecSliceUint64R(f *codecFnInfo, rv reflect.Value) {
	if array := f.seq == seqTypeArray; !array && rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*[]uint64)
		if v, changed := fastpathTV.DecSliceUint64V(*vp, !array, d); changed {
			*vp = v
		}
	} else {
		v := rv2i(rv).([]uint64)
		v2, changed := fastpathTV.DecSliceUint64V(v, !array, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	}
}
func (f fastpathT) DecSliceUint64X(vp *[]uint64, d *Decoder) {
	if v, changed := f.DecSliceUint64V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecSliceUint64V(v []uint64, canChange bool, d *Decoder) (_ []uint64, changed bool) {
	slh, containerLenS := d.decSliceHelperStart()
	if containerLenS == 0 {
		if canChange {
			if v == nil {
				v = []uint64{}
			} else if len(v) != 0 {
				v = v[:0]
			}
			changed = true
		}
		slh.End()
		return v, changed
	}
	hasLen := containerLenS > 0
	var xlen int
	if hasLen && canChange {
		if containerLenS > cap(v) {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 8)
			if xlen <= cap(v) {
				v = v[:uint(xlen)]
			} else {
				v = make([]uint64, uint(xlen))
			}
			changed = true
		} else if containerLenS != len(v) {
			v = v[:containerLenS]
			changed = true
		}
	}
	var j int
	for j = 0; (hasLen && j < containerLenS) || !(hasLen || d.d.CheckBreak()); j++ {
		if j == 0 && len(v) == 0 && canChange {
			if hasLen {
				xlen = decInferLen(containerLenS, d.h.MaxInitLen, 8)
			} else {
				xlen = 8
			}
			v = make([]uint64, uint(xlen))
			changed = true
		}
		var decodeIntoBlank bool
		if j >= len(v) {
			if canChange {
				v = append(v, 0)
				changed = true
			} else {
				d.arrayCannotExpand(len(v), j+1)
				decodeIntoBlank = true
			}
		}
		slh.ElemContainerState(j)
		if decodeIntoBlank {
			d.swallow()
		} else if d.d.TryDecodeAsNil() {
			v[uint(j)] = 0
		} else {
			v[uint(j)] = d.d.DecodeUint64()
		}
	}
	if canChange {
		if j < len(v) {
			v = v[:uint(j)]
			changed = true
		} else if j == 0 && v == nil {
			v = make([]uint64, 0)
			changed = true
		}
	}
	slh.End()
	return v, changed
}

func (d *Decoder) fastpathDecSliceUintptrR(f *codecFnInfo, rv reflect.Value) {
	if array := f.seq == seqTypeArray; !array && rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*[]uintptr)
		if v, changed := fastpathTV.DecSliceUintptrV(*vp, !array, d); changed {
			*vp = v
		}
	} else {
		v := rv2i(rv).([]uintptr)
		v2, changed := fastpathTV.DecSliceUintptrV(v, !array, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	}
}
func (f fastpathT) DecSliceUintptrX(vp *[]uintptr, d *Decoder) {
	if v, changed := f.DecSliceUintptrV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecSliceUintptrV(v []uintptr, canChange bool, d *Decoder) (_ []uintptr, changed bool) {
	slh, containerLenS := d.decSliceHelperStart()
	if containerLenS == 0 {
		if canChange {
			if v == nil {
				v = []uintptr{}
			} else if len(v) != 0 {
				v = v[:0]
			}
			changed = true
		}
		slh.End()
		return v, changed
	}
	hasLen := containerLenS > 0
	var xlen int
	if hasLen && canChange {
		if containerLenS > cap(v) {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 8)
			if xlen <= cap(v) {
				v = v[:uint(xlen)]
			} else {
				v = make([]uintptr, uint(xlen))
			}
			changed = true
		} else if containerLenS != len(v) {
			v = v[:containerLenS]
			changed = true
		}
	}
	var j int
	for j = 0; (hasLen && j < containerLenS) || !(hasLen || d.d.CheckBreak()); j++ {
		if j == 0 && len(v) == 0 && canChange {
			if hasLen {
				xlen = decInferLen(containerLenS, d.h.MaxInitLen, 8)
			} else {
				xlen = 8
			}
			v = make([]uintptr, uint(xlen))
			changed = true
		}
		var decodeIntoBlank bool
		if j >= len(v) {
			if canChange {
				v = append(v, 0)
				changed = true
			} else {
				d.arrayCannotExpand(len(v), j+1)
				decodeIntoBlank = true
			}
		}
		slh.ElemContainerState(j)
		if decodeIntoBlank {
			d.swallow()
		} else if d.d.TryDecodeAsNil() {
			v[uint(j)] = 0
		} else {
			v[uint(j)] = uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
		}
	}
	if canChange {
		if j < len(v) {
			v = v[:uint(j)]
			changed = true
		} else if j == 0 && v == nil {
			v = make([]uintptr, 0)
			changed = true
		}
	}
	slh.End()
	return v, changed
}

func (d *Decoder) fastpathDecSliceIntR(f *codecFnInfo, rv reflect.Value) {
	if array := f.seq == seqTypeArray; !array && rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*[]int)
		if v, changed := fastpathTV.DecSliceIntV(*vp, !array, d); changed {
			*vp = v
		}
	} else {
		v := rv2i(rv).([]int)
		v2, changed := fastpathTV.DecSliceIntV(v, !array, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	}
}
func (f fastpathT) DecSliceIntX(vp *[]int, d *Decoder) {
	if v, changed := f.DecSliceIntV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecSliceIntV(v []int, canChange bool, d *Decoder) (_ []int, changed bool) {
	slh, containerLenS := d.decSliceHelperStart()
	if containerLenS == 0 {
		if canChange {
			if v == nil {
				v = []int{}
			} else if len(v) != 0 {
				v = v[:0]
			}
			changed = true
		}
		slh.End()
		return v, changed
	}
	hasLen := containerLenS > 0
	var xlen int
	if hasLen && canChange {
		if containerLenS > cap(v) {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 8)
			if xlen <= cap(v) {
				v = v[:uint(xlen)]
			} else {
				v = make([]int, uint(xlen))
			}
			changed = true
		} else if containerLenS != len(v) {
			v = v[:containerLenS]
			changed = true
		}
	}
	var j int
	for j = 0; (hasLen && j < containerLenS) || !(hasLen || d.d.CheckBreak()); j++ {
		if j == 0 && len(v) == 0 && canChange {
			if hasLen {
				xlen = decInferLen(containerLenS, d.h.MaxInitLen, 8)
			} else {
				xlen = 8
			}
			v = make([]int, uint(xlen))
			changed = true
		}
		var decodeIntoBlank bool
		if j >= len(v) {
			if canChange {
				v = append(v, 0)
				changed = true
			} else {
				d.arrayCannotExpand(len(v), j+1)
				decodeIntoBlank = true
			}
		}
		slh.ElemContainerState(j)
		if decodeIntoBlank {
			d.swallow()
		} else if d.d.TryDecodeAsNil() {
			v[uint(j)] = 0
		} else {
			v[uint(j)] = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		}
	}
	if canChange {
		if j < len(v) {
			v = v[:uint(j)]
			changed = true
		} else if j == 0 && v == nil {
			v = make([]int, 0)
			changed = true
		}
	}
	slh.End()
	return v, changed
}

func (d *Decoder) fastpathDecSliceInt8R(f *codecFnInfo, rv reflect.Value) {
	if array := f.seq == seqTypeArray; !array && rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*[]int8)
		if v, changed := fastpathTV.DecSliceInt8V(*vp, !array, d); changed {
			*vp = v
		}
	} else {
		v := rv2i(rv).([]int8)
		v2, changed := fastpathTV.DecSliceInt8V(v, !array, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	}
}
func (f fastpathT) DecSliceInt8X(vp *[]int8, d *Decoder) {
	if v, changed := f.DecSliceInt8V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecSliceInt8V(v []int8, canChange bool, d *Decoder) (_ []int8, changed bool) {
	slh, containerLenS := d.decSliceHelperStart()
	if containerLenS == 0 {
		if canChange {
			if v == nil {
				v = []int8{}
			} else if len(v) != 0 {
				v = v[:0]
			}
			changed = true
		}
		slh.End()
		return v, changed
	}
	hasLen := containerLenS > 0
	var xlen int
	if hasLen && canChange {
		if containerLenS > cap(v) {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 1)
			if xlen <= cap(v) {
				v = v[:uint(xlen)]
			} else {
				v = make([]int8, uint(xlen))
			}
			changed = true
		} else if containerLenS != len(v) {
			v = v[:containerLenS]
			changed = true
		}
	}
	var j int
	for j = 0; (hasLen && j < containerLenS) || !(hasLen || d.d.CheckBreak()); j++ {
		if j == 0 && len(v) == 0 && canChange {
			if hasLen {
				xlen = decInferLen(containerLenS, d.h.MaxInitLen, 1)
			} else {
				xlen = 8
			}
			v = make([]int8, uint(xlen))
			changed = true
		}
		var decodeIntoBlank bool
		if j >= len(v) {
			if canChange {
				v = append(v, 0)
				changed = true
			} else {
				d.arrayCannotExpand(len(v), j+1)
				decodeIntoBlank = true
			}
		}
		slh.ElemContainerState(j)
		if decodeIntoBlank {
			d.swallow()
		} else if d.d.TryDecodeAsNil() {
			v[uint(j)] = 0
		} else {
			v[uint(j)] = int8(chkOvf.IntV(d.d.DecodeInt64(), 8))
		}
	}
	if canChange {
		if j < len(v) {
			v = v[:uint(j)]
			changed = true
		} else if j == 0 && v == nil {
			v = make([]int8, 0)
			changed = true
		}
	}
	slh.End()
	return v, changed
}

func (d *Decoder) fastpathDecSliceInt16R(f *codecFnInfo, rv reflect.Value) {
	if array := f.seq == seqTypeArray; !array && rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*[]int16)
		if v, changed := fastpathTV.DecSliceInt16V(*vp, !array, d); changed {
			*vp = v
		}
	} else {
		v := rv2i(rv).([]int16)
		v2, changed := fastpathTV.DecSliceInt16V(v, !array, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	}
}
func (f fastpathT) DecSliceInt16X(vp *[]int16, d *Decoder) {
	if v, changed := f.DecSliceInt16V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecSliceInt16V(v []int16, canChange bool, d *Decoder) (_ []int16, changed bool) {
	slh, containerLenS := d.decSliceHelperStart()
	if containerLenS == 0 {
		if canChange {
			if v == nil {
				v = []int16{}
			} else if len(v) != 0 {
				v = v[:0]
			}
			changed = true
		}
		slh.End()
		return v, changed
	}
	hasLen := containerLenS > 0
	var xlen int
	if hasLen && canChange {
		if containerLenS > cap(v) {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 2)
			if xlen <= cap(v) {
				v = v[:uint(xlen)]
			} else {
				v = make([]int16, uint(xlen))
			}
			changed = true
		} else if containerLenS != len(v) {
			v = v[:containerLenS]
			changed = true
		}
	}
	var j int
	for j = 0; (hasLen && j < containerLenS) || !(hasLen || d.d.CheckBreak()); j++ {
		if j == 0 && len(v) == 0 && canChange {
			if hasLen {
				xlen = decInferLen(containerLenS, d.h.MaxInitLen, 2)
			} else {
				xlen = 8
			}
			v = make([]int16, uint(xlen))
			changed = true
		}
		var decodeIntoBlank bool
		if j >= len(v) {
			if canChange {
				v = append(v, 0)
				changed = true
			} else {
				d.arrayCannotExpand(len(v), j+1)
				decodeIntoBlank = true
			}
		}
		slh.ElemContainerState(j)
		if decodeIntoBlank {
			d.swallow()
		} else if d.d.TryDecodeAsNil() {
			v[uint(j)] = 0
		} else {
			v[uint(j)] = int16(chkOvf.IntV(d.d.DecodeInt64(), 16))
		}
	}
	if canChange {
		if j < len(v) {
			v = v[:uint(j)]
			changed = true
		} else if j == 0 && v == nil {
			v = make([]int16, 0)
			changed = true
		}
	}
	slh.End()
	return v, changed
}

func (d *Decoder) fastpathDecSliceInt32R(f *codecFnInfo, rv reflect.Value) {
	if array := f.seq == seqTypeArray; !array && rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*[]int32)
		if v, changed := fastpathTV.DecSliceInt32V(*vp, !array, d); changed {
			*vp = v
		}
	} else {
		v := rv2i(rv).([]int32)
		v2, changed := fastpathTV.DecSliceInt32V(v, !array, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	}
}
func (f fastpathT) DecSliceInt32X(vp *[]int32, d *Decoder) {
	if v, changed := f.DecSliceInt32V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecSliceInt32V(v []int32, canChange bool, d *Decoder) (_ []int32, changed bool) {
	slh, containerLenS := d.decSliceHelperStart()
	if containerLenS == 0 {
		if canChange {
			if v == nil {
				v = []int32{}
			} else if len(v) != 0 {
				v = v[:0]
			}
			changed = true
		}
		slh.End()
		return v, changed
	}
	hasLen := containerLenS > 0
	var xlen int
	if hasLen && canChange {
		if containerLenS > cap(v) {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 4)
			if xlen <= cap(v) {
				v = v[:uint(xlen)]
			} else {
				v = make([]int32, uint(xlen))
			}
			changed = true
		} else if containerLenS != len(v) {
			v = v[:containerLenS]
			changed = true
		}
	}
	var j int
	for j = 0; (hasLen && j < containerLenS) || !(hasLen || d.d.CheckBreak()); j++ {
		if j == 0 && len(v) == 0 && canChange {
			if hasLen {
				xlen = decInferLen(containerLenS, d.h.MaxInitLen, 4)
			} else {
				xlen = 8
			}
			v = make([]int32, uint(xlen))
			changed = true
		}
		var decodeIntoBlank bool
		if j >= len(v) {
			if canChange {
				v = append(v, 0)
				changed = true
			} else {
				d.arrayCannotExpand(len(v), j+1)
				decodeIntoBlank = true
			}
		}
		slh.ElemContainerState(j)
		if decodeIntoBlank {
			d.swallow()
		} else if d.d.TryDecodeAsNil() {
			v[uint(j)] = 0
		} else {
			v[uint(j)] = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		}
	}
	if canChange {
		if j < len(v) {
			v = v[:uint(j)]
			changed = true
		} else if j == 0 && v == nil {
			v = make([]int32, 0)
			changed = true
		}
	}
	slh.End()
	return v, changed
}

func (d *Decoder) fastpathDecSliceInt64R(f *codecFnInfo, rv reflect.Value) {
	if array := f.seq == seqTypeArray; !array && rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*[]int64)
		if v, changed := fastpathTV.DecSliceInt64V(*vp, !array, d); changed {
			*vp = v
		}
	} else {
		v := rv2i(rv).([]int64)
		v2, changed := fastpathTV.DecSliceInt64V(v, !array, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	}
}
func (f fastpathT) DecSliceInt64X(vp *[]int64, d *Decoder) {
	if v, changed := f.DecSliceInt64V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecSliceInt64V(v []int64, canChange bool, d *Decoder) (_ []int64, changed bool) {
	slh, containerLenS := d.decSliceHelperStart()
	if containerLenS == 0 {
		if canChange {
			if v == nil {
				v = []int64{}
			} else if len(v) != 0 {
				v = v[:0]
			}
			changed = true
		}
		slh.End()
		return v, changed
	}
	hasLen := containerLenS > 0
	var xlen int
	if hasLen && canChange {
		if containerLenS > cap(v) {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 8)
			if xlen <= cap(v) {
				v = v[:uint(xlen)]
			} else {
				v = make([]int64, uint(xlen))
			}
			changed = true
		} else if containerLenS != len(v) {
			v = v[:containerLenS]
			changed = true
		}
	}
	var j int
	for j = 0; (hasLen && j < containerLenS) || !(hasLen || d.d.CheckBreak()); j++ {
		if j == 0 && len(v) == 0 && canChange {
			if hasLen {
				xlen = decInferLen(containerLenS, d.h.MaxInitLen, 8)
			} else {
				xlen = 8
			}
			v = make([]int64, uint(xlen))
			changed = true
		}
		var decodeIntoBlank bool
		if j >= len(v) {
			if canChange {
				v = append(v, 0)
				changed = true
			} else {
				d.arrayCannotExpand(len(v), j+1)
				decodeIntoBlank = true
			}
		}
		slh.ElemContainerState(j)
		if decodeIntoBlank {
			d.swallow()
		} else if d.d.TryDecodeAsNil() {
			v[uint(j)] = 0
		} else {
			v[uint(j)] = d.d.DecodeInt64()
		}
	}
	if canChange {
		if j < len(v) {
			v = v[:uint(j)]
			changed = true
		} else if j == 0 && v == nil {
			v = make([]int64, 0)
			changed = true
		}
	}
	slh.End()
	return v, changed
}

func (d *Decoder) fastpathDecSliceBoolR(f *codecFnInfo, rv reflect.Value) {
	if array := f.seq == seqTypeArray; !array && rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*[]bool)
		if v, changed := fastpathTV.DecSliceBoolV(*vp, !array, d); changed {
			*vp = v
		}
	} else {
		v := rv2i(rv).([]bool)
		v2, changed := fastpathTV.DecSliceBoolV(v, !array, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	}
}
func (f fastpathT) DecSliceBoolX(vp *[]bool, d *Decoder) {
	if v, changed := f.DecSliceBoolV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecSliceBoolV(v []bool, canChange bool, d *Decoder) (_ []bool, changed bool) {
	slh, containerLenS := d.decSliceHelperStart()
	if containerLenS == 0 {
		if canChange {
			if v == nil {
				v = []bool{}
			} else if len(v) != 0 {
				v = v[:0]
			}
			changed = true
		}
		slh.End()
		return v, changed
	}
	hasLen := containerLenS > 0
	var xlen int
	if hasLen && canChange {
		if containerLenS > cap(v) {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 1)
			if xlen <= cap(v) {
				v = v[:uint(xlen)]
			} else {
				v = make([]bool, uint(xlen))
			}
			changed = true
		} else if containerLenS != len(v) {
			v = v[:containerLenS]
			changed = true
		}
	}
	var j int
	for j = 0; (hasLen && j < containerLenS) || !(hasLen || d.d.CheckBreak()); j++ {
		if j == 0 && len(v) == 0 && canChange {
			if hasLen {
				xlen = decInferLen(containerLenS, d.h.MaxInitLen, 1)
			} else {
				xlen = 8
			}
			v = make([]bool, uint(xlen))
			changed = true
		}
		var decodeIntoBlank bool
		if j >= len(v) {
			if canChange {
				v = append(v, false)
				changed = true
			} else {
				d.arrayCannotExpand(len(v), j+1)
				decodeIntoBlank = true
			}
		}
		slh.ElemContainerState(j)
		if decodeIntoBlank {
			d.swallow()
		} else if d.d.TryDecodeAsNil() {
			v[uint(j)] = false
		} else {
			v[uint(j)] = d.d.DecodeBool()
		}
	}
	if canChange {
		if j < len(v) {
			v = v[:uint(j)]
			changed = true
		} else if j == 0 && v == nil {
			v = make([]bool, 0)
			changed = true
		}
	}
	slh.End()
	return v, changed
}

func (d *Decoder) fastpathDecMapStringIntfR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[string]interface{})
		if v, changed := fastpathTV.DecMapStringIntfV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapStringIntfV(rv2i(rv).(map[string]interface{}), false, d)
	}
}
func (f fastpathT) DecMapStringIntfX(vp *map[string]interface{}, d *Decoder) {
	if v, changed := f.DecMapStringIntfV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapStringIntfV(v map[string]interface{}, canChange bool,
	d *Decoder) (_ map[string]interface{}, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[string]interface{}, decInferLen(containerLen, d.h.MaxInitLen, 32))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	mapGet := v != nil && !d.h.MapValueReset && !d.h.InterfaceReset
	var mk string
	var mv interface{}
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeString()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = nil
			}
			continue
		}
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		d.decode(&mv)
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapStringStringR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[string]string)
		if v, changed := fastpathTV.DecMapStringStringV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapStringStringV(rv2i(rv).(map[string]string), false, d)
	}
}
func (f fastpathT) DecMapStringStringX(vp *map[string]string, d *Decoder) {
	if v, changed := f.DecMapStringStringV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapStringStringV(v map[string]string, canChange bool,
	d *Decoder) (_ map[string]string, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[string]string, decInferLen(containerLen, d.h.MaxInitLen, 32))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk string
	var mv string
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeString()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = ""
			}
			continue
		}
		mv = d.d.DecodeString()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapStringBytesR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[string][]byte)
		if v, changed := fastpathTV.DecMapStringBytesV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapStringBytesV(rv2i(rv).(map[string][]byte), false, d)
	}
}
func (f fastpathT) DecMapStringBytesX(vp *map[string][]byte, d *Decoder) {
	if v, changed := f.DecMapStringBytesV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapStringBytesV(v map[string][]byte, canChange bool,
	d *Decoder) (_ map[string][]byte, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[string][]byte, decInferLen(containerLen, d.h.MaxInitLen, 40))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	mapGet := v != nil && !d.h.MapValueReset
	var mk string
	var mv []byte
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeString()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = nil
			}
			continue
		}
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		mv = d.d.DecodeBytes(mv, false)
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapStringUintR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[string]uint)
		if v, changed := fastpathTV.DecMapStringUintV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapStringUintV(rv2i(rv).(map[string]uint), false, d)
	}
}
func (f fastpathT) DecMapStringUintX(vp *map[string]uint, d *Decoder) {
	if v, changed := f.DecMapStringUintV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapStringUintV(v map[string]uint, canChange bool,
	d *Decoder) (_ map[string]uint, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[string]uint, decInferLen(containerLen, d.h.MaxInitLen, 24))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk string
	var mv uint
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeString()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapStringUint8R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[string]uint8)
		if v, changed := fastpathTV.DecMapStringUint8V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapStringUint8V(rv2i(rv).(map[string]uint8), false, d)
	}
}
func (f fastpathT) DecMapStringUint8X(vp *map[string]uint8, d *Decoder) {
	if v, changed := f.DecMapStringUint8V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapStringUint8V(v map[string]uint8, canChange bool,
	d *Decoder) (_ map[string]uint8, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[string]uint8, decInferLen(containerLen, d.h.MaxInitLen, 17))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk string
	var mv uint8
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeString()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapStringUint64R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[string]uint64)
		if v, changed := fastpathTV.DecMapStringUint64V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapStringUint64V(rv2i(rv).(map[string]uint64), false, d)
	}
}
func (f fastpathT) DecMapStringUint64X(vp *map[string]uint64, d *Decoder) {
	if v, changed := f.DecMapStringUint64V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapStringUint64V(v map[string]uint64, canChange bool,
	d *Decoder) (_ map[string]uint64, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[string]uint64, decInferLen(containerLen, d.h.MaxInitLen, 24))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk string
	var mv uint64
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeString()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = d.d.DecodeUint64()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapStringUintptrR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[string]uintptr)
		if v, changed := fastpathTV.DecMapStringUintptrV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapStringUintptrV(rv2i(rv).(map[string]uintptr), false, d)
	}
}
func (f fastpathT) DecMapStringUintptrX(vp *map[string]uintptr, d *Decoder) {
	if v, changed := f.DecMapStringUintptrV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapStringUintptrV(v map[string]uintptr, canChange bool,
	d *Decoder) (_ map[string]uintptr, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[string]uintptr, decInferLen(containerLen, d.h.MaxInitLen, 24))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk string
	var mv uintptr
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeString()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapStringIntR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[string]int)
		if v, changed := fastpathTV.DecMapStringIntV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapStringIntV(rv2i(rv).(map[string]int), false, d)
	}
}
func (f fastpathT) DecMapStringIntX(vp *map[string]int, d *Decoder) {
	if v, changed := f.DecMapStringIntV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapStringIntV(v map[string]int, canChange bool,
	d *Decoder) (_ map[string]int, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[string]int, decInferLen(containerLen, d.h.MaxInitLen, 24))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk string
	var mv int
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeString()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapStringInt64R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[string]int64)
		if v, changed := fastpathTV.DecMapStringInt64V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapStringInt64V(rv2i(rv).(map[string]int64), false, d)
	}
}
func (f fastpathT) DecMapStringInt64X(vp *map[string]int64, d *Decoder) {
	if v, changed := f.DecMapStringInt64V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapStringInt64V(v map[string]int64, canChange bool,
	d *Decoder) (_ map[string]int64, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[string]int64, decInferLen(containerLen, d.h.MaxInitLen, 24))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk string
	var mv int64
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeString()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = d.d.DecodeInt64()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapStringFloat32R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[string]float32)
		if v, changed := fastpathTV.DecMapStringFloat32V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapStringFloat32V(rv2i(rv).(map[string]float32), false, d)
	}
}
func (f fastpathT) DecMapStringFloat32X(vp *map[string]float32, d *Decoder) {
	if v, changed := f.DecMapStringFloat32V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapStringFloat32V(v map[string]float32, canChange bool,
	d *Decoder) (_ map[string]float32, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[string]float32, decInferLen(containerLen, d.h.MaxInitLen, 20))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk string
	var mv float32
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeString()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = float32(d.decodeFloat32())
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapStringFloat64R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[string]float64)
		if v, changed := fastpathTV.DecMapStringFloat64V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapStringFloat64V(rv2i(rv).(map[string]float64), false, d)
	}
}
func (f fastpathT) DecMapStringFloat64X(vp *map[string]float64, d *Decoder) {
	if v, changed := f.DecMapStringFloat64V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapStringFloat64V(v map[string]float64, canChange bool,
	d *Decoder) (_ map[string]float64, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[string]float64, decInferLen(containerLen, d.h.MaxInitLen, 24))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk string
	var mv float64
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeString()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = d.d.DecodeFloat64()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapStringBoolR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[string]bool)
		if v, changed := fastpathTV.DecMapStringBoolV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapStringBoolV(rv2i(rv).(map[string]bool), false, d)
	}
}
func (f fastpathT) DecMapStringBoolX(vp *map[string]bool, d *Decoder) {
	if v, changed := f.DecMapStringBoolV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapStringBoolV(v map[string]bool, canChange bool,
	d *Decoder) (_ map[string]bool, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[string]bool, decInferLen(containerLen, d.h.MaxInitLen, 17))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk string
	var mv bool
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeString()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = false
			}
			continue
		}
		mv = d.d.DecodeBool()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUintIntfR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint]interface{})
		if v, changed := fastpathTV.DecMapUintIntfV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUintIntfV(rv2i(rv).(map[uint]interface{}), false, d)
	}
}
func (f fastpathT) DecMapUintIntfX(vp *map[uint]interface{}, d *Decoder) {
	if v, changed := f.DecMapUintIntfV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUintIntfV(v map[uint]interface{}, canChange bool,
	d *Decoder) (_ map[uint]interface{}, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint]interface{}, decInferLen(containerLen, d.h.MaxInitLen, 24))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	mapGet := v != nil && !d.h.MapValueReset && !d.h.InterfaceReset
	var mk uint
	var mv interface{}
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = nil
			}
			continue
		}
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		d.decode(&mv)
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUintStringR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint]string)
		if v, changed := fastpathTV.DecMapUintStringV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUintStringV(rv2i(rv).(map[uint]string), false, d)
	}
}
func (f fastpathT) DecMapUintStringX(vp *map[uint]string, d *Decoder) {
	if v, changed := f.DecMapUintStringV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUintStringV(v map[uint]string, canChange bool,
	d *Decoder) (_ map[uint]string, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint]string, decInferLen(containerLen, d.h.MaxInitLen, 24))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint
	var mv string
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = ""
			}
			continue
		}
		mv = d.d.DecodeString()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUintBytesR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint][]byte)
		if v, changed := fastpathTV.DecMapUintBytesV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUintBytesV(rv2i(rv).(map[uint][]byte), false, d)
	}
}
func (f fastpathT) DecMapUintBytesX(vp *map[uint][]byte, d *Decoder) {
	if v, changed := f.DecMapUintBytesV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUintBytesV(v map[uint][]byte, canChange bool,
	d *Decoder) (_ map[uint][]byte, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint][]byte, decInferLen(containerLen, d.h.MaxInitLen, 32))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	mapGet := v != nil && !d.h.MapValueReset
	var mk uint
	var mv []byte
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = nil
			}
			continue
		}
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		mv = d.d.DecodeBytes(mv, false)
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUintUintR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint]uint)
		if v, changed := fastpathTV.DecMapUintUintV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUintUintV(rv2i(rv).(map[uint]uint), false, d)
	}
}
func (f fastpathT) DecMapUintUintX(vp *map[uint]uint, d *Decoder) {
	if v, changed := f.DecMapUintUintV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUintUintV(v map[uint]uint, canChange bool,
	d *Decoder) (_ map[uint]uint, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint]uint, decInferLen(containerLen, d.h.MaxInitLen, 16))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint
	var mv uint
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUintUint8R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint]uint8)
		if v, changed := fastpathTV.DecMapUintUint8V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUintUint8V(rv2i(rv).(map[uint]uint8), false, d)
	}
}
func (f fastpathT) DecMapUintUint8X(vp *map[uint]uint8, d *Decoder) {
	if v, changed := f.DecMapUintUint8V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUintUint8V(v map[uint]uint8, canChange bool,
	d *Decoder) (_ map[uint]uint8, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint]uint8, decInferLen(containerLen, d.h.MaxInitLen, 9))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint
	var mv uint8
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUintUint64R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint]uint64)
		if v, changed := fastpathTV.DecMapUintUint64V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUintUint64V(rv2i(rv).(map[uint]uint64), false, d)
	}
}
func (f fastpathT) DecMapUintUint64X(vp *map[uint]uint64, d *Decoder) {
	if v, changed := f.DecMapUintUint64V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUintUint64V(v map[uint]uint64, canChange bool,
	d *Decoder) (_ map[uint]uint64, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint]uint64, decInferLen(containerLen, d.h.MaxInitLen, 16))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint
	var mv uint64
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = d.d.DecodeUint64()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUintUintptrR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint]uintptr)
		if v, changed := fastpathTV.DecMapUintUintptrV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUintUintptrV(rv2i(rv).(map[uint]uintptr), false, d)
	}
}
func (f fastpathT) DecMapUintUintptrX(vp *map[uint]uintptr, d *Decoder) {
	if v, changed := f.DecMapUintUintptrV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUintUintptrV(v map[uint]uintptr, canChange bool,
	d *Decoder) (_ map[uint]uintptr, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint]uintptr, decInferLen(containerLen, d.h.MaxInitLen, 16))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint
	var mv uintptr
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUintIntR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint]int)
		if v, changed := fastpathTV.DecMapUintIntV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUintIntV(rv2i(rv).(map[uint]int), false, d)
	}
}
func (f fastpathT) DecMapUintIntX(vp *map[uint]int, d *Decoder) {
	if v, changed := f.DecMapUintIntV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUintIntV(v map[uint]int, canChange bool,
	d *Decoder) (_ map[uint]int, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint]int, decInferLen(containerLen, d.h.MaxInitLen, 16))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint
	var mv int
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUintInt64R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint]int64)
		if v, changed := fastpathTV.DecMapUintInt64V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUintInt64V(rv2i(rv).(map[uint]int64), false, d)
	}
}
func (f fastpathT) DecMapUintInt64X(vp *map[uint]int64, d *Decoder) {
	if v, changed := f.DecMapUintInt64V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUintInt64V(v map[uint]int64, canChange bool,
	d *Decoder) (_ map[uint]int64, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint]int64, decInferLen(containerLen, d.h.MaxInitLen, 16))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint
	var mv int64
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = d.d.DecodeInt64()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUintFloat32R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint]float32)
		if v, changed := fastpathTV.DecMapUintFloat32V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUintFloat32V(rv2i(rv).(map[uint]float32), false, d)
	}
}
func (f fastpathT) DecMapUintFloat32X(vp *map[uint]float32, d *Decoder) {
	if v, changed := f.DecMapUintFloat32V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUintFloat32V(v map[uint]float32, canChange bool,
	d *Decoder) (_ map[uint]float32, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint]float32, decInferLen(containerLen, d.h.MaxInitLen, 12))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint
	var mv float32
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = float32(d.decodeFloat32())
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUintFloat64R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint]float64)
		if v, changed := fastpathTV.DecMapUintFloat64V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUintFloat64V(rv2i(rv).(map[uint]float64), false, d)
	}
}
func (f fastpathT) DecMapUintFloat64X(vp *map[uint]float64, d *Decoder) {
	if v, changed := f.DecMapUintFloat64V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUintFloat64V(v map[uint]float64, canChange bool,
	d *Decoder) (_ map[uint]float64, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint]float64, decInferLen(containerLen, d.h.MaxInitLen, 16))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint
	var mv float64
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = d.d.DecodeFloat64()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUintBoolR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint]bool)
		if v, changed := fastpathTV.DecMapUintBoolV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUintBoolV(rv2i(rv).(map[uint]bool), false, d)
	}
}
func (f fastpathT) DecMapUintBoolX(vp *map[uint]bool, d *Decoder) {
	if v, changed := f.DecMapUintBoolV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUintBoolV(v map[uint]bool, canChange bool,
	d *Decoder) (_ map[uint]bool, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint]bool, decInferLen(containerLen, d.h.MaxInitLen, 9))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint
	var mv bool
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = false
			}
			continue
		}
		mv = d.d.DecodeBool()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUint8IntfR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint8]interface{})
		if v, changed := fastpathTV.DecMapUint8IntfV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUint8IntfV(rv2i(rv).(map[uint8]interface{}), false, d)
	}
}
func (f fastpathT) DecMapUint8IntfX(vp *map[uint8]interface{}, d *Decoder) {
	if v, changed := f.DecMapUint8IntfV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUint8IntfV(v map[uint8]interface{}, canChange bool,
	d *Decoder) (_ map[uint8]interface{}, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint8]interface{}, decInferLen(containerLen, d.h.MaxInitLen, 17))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	mapGet := v != nil && !d.h.MapValueReset && !d.h.InterfaceReset
	var mk uint8
	var mv interface{}
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = nil
			}
			continue
		}
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		d.decode(&mv)
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUint8StringR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint8]string)
		if v, changed := fastpathTV.DecMapUint8StringV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUint8StringV(rv2i(rv).(map[uint8]string), false, d)
	}
}
func (f fastpathT) DecMapUint8StringX(vp *map[uint8]string, d *Decoder) {
	if v, changed := f.DecMapUint8StringV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUint8StringV(v map[uint8]string, canChange bool,
	d *Decoder) (_ map[uint8]string, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint8]string, decInferLen(containerLen, d.h.MaxInitLen, 17))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint8
	var mv string
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = ""
			}
			continue
		}
		mv = d.d.DecodeString()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUint8BytesR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint8][]byte)
		if v, changed := fastpathTV.DecMapUint8BytesV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUint8BytesV(rv2i(rv).(map[uint8][]byte), false, d)
	}
}
func (f fastpathT) DecMapUint8BytesX(vp *map[uint8][]byte, d *Decoder) {
	if v, changed := f.DecMapUint8BytesV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUint8BytesV(v map[uint8][]byte, canChange bool,
	d *Decoder) (_ map[uint8][]byte, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint8][]byte, decInferLen(containerLen, d.h.MaxInitLen, 25))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	mapGet := v != nil && !d.h.MapValueReset
	var mk uint8
	var mv []byte
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = nil
			}
			continue
		}
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		mv = d.d.DecodeBytes(mv, false)
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUint8UintR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint8]uint)
		if v, changed := fastpathTV.DecMapUint8UintV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUint8UintV(rv2i(rv).(map[uint8]uint), false, d)
	}
}
func (f fastpathT) DecMapUint8UintX(vp *map[uint8]uint, d *Decoder) {
	if v, changed := f.DecMapUint8UintV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUint8UintV(v map[uint8]uint, canChange bool,
	d *Decoder) (_ map[uint8]uint, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint8]uint, decInferLen(containerLen, d.h.MaxInitLen, 9))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint8
	var mv uint
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUint8Uint8R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint8]uint8)
		if v, changed := fastpathTV.DecMapUint8Uint8V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUint8Uint8V(rv2i(rv).(map[uint8]uint8), false, d)
	}
}
func (f fastpathT) DecMapUint8Uint8X(vp *map[uint8]uint8, d *Decoder) {
	if v, changed := f.DecMapUint8Uint8V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUint8Uint8V(v map[uint8]uint8, canChange bool,
	d *Decoder) (_ map[uint8]uint8, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint8]uint8, decInferLen(containerLen, d.h.MaxInitLen, 2))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint8
	var mv uint8
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUint8Uint64R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint8]uint64)
		if v, changed := fastpathTV.DecMapUint8Uint64V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUint8Uint64V(rv2i(rv).(map[uint8]uint64), false, d)
	}
}
func (f fastpathT) DecMapUint8Uint64X(vp *map[uint8]uint64, d *Decoder) {
	if v, changed := f.DecMapUint8Uint64V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUint8Uint64V(v map[uint8]uint64, canChange bool,
	d *Decoder) (_ map[uint8]uint64, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint8]uint64, decInferLen(containerLen, d.h.MaxInitLen, 9))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint8
	var mv uint64
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = d.d.DecodeUint64()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUint8UintptrR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint8]uintptr)
		if v, changed := fastpathTV.DecMapUint8UintptrV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUint8UintptrV(rv2i(rv).(map[uint8]uintptr), false, d)
	}
}
func (f fastpathT) DecMapUint8UintptrX(vp *map[uint8]uintptr, d *Decoder) {
	if v, changed := f.DecMapUint8UintptrV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUint8UintptrV(v map[uint8]uintptr, canChange bool,
	d *Decoder) (_ map[uint8]uintptr, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint8]uintptr, decInferLen(containerLen, d.h.MaxInitLen, 9))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint8
	var mv uintptr
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUint8IntR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint8]int)
		if v, changed := fastpathTV.DecMapUint8IntV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUint8IntV(rv2i(rv).(map[uint8]int), false, d)
	}
}
func (f fastpathT) DecMapUint8IntX(vp *map[uint8]int, d *Decoder) {
	if v, changed := f.DecMapUint8IntV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUint8IntV(v map[uint8]int, canChange bool,
	d *Decoder) (_ map[uint8]int, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint8]int, decInferLen(containerLen, d.h.MaxInitLen, 9))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint8
	var mv int
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUint8Int64R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint8]int64)
		if v, changed := fastpathTV.DecMapUint8Int64V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUint8Int64V(rv2i(rv).(map[uint8]int64), false, d)
	}
}
func (f fastpathT) DecMapUint8Int64X(vp *map[uint8]int64, d *Decoder) {
	if v, changed := f.DecMapUint8Int64V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUint8Int64V(v map[uint8]int64, canChange bool,
	d *Decoder) (_ map[uint8]int64, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint8]int64, decInferLen(containerLen, d.h.MaxInitLen, 9))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint8
	var mv int64
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = d.d.DecodeInt64()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUint8Float32R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint8]float32)
		if v, changed := fastpathTV.DecMapUint8Float32V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUint8Float32V(rv2i(rv).(map[uint8]float32), false, d)
	}
}
func (f fastpathT) DecMapUint8Float32X(vp *map[uint8]float32, d *Decoder) {
	if v, changed := f.DecMapUint8Float32V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUint8Float32V(v map[uint8]float32, canChange bool,
	d *Decoder) (_ map[uint8]float32, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint8]float32, decInferLen(containerLen, d.h.MaxInitLen, 5))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint8
	var mv float32
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = float32(d.decodeFloat32())
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUint8Float64R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint8]float64)
		if v, changed := fastpathTV.DecMapUint8Float64V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUint8Float64V(rv2i(rv).(map[uint8]float64), false, d)
	}
}
func (f fastpathT) DecMapUint8Float64X(vp *map[uint8]float64, d *Decoder) {
	if v, changed := f.DecMapUint8Float64V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUint8Float64V(v map[uint8]float64, canChange bool,
	d *Decoder) (_ map[uint8]float64, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint8]float64, decInferLen(containerLen, d.h.MaxInitLen, 9))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint8
	var mv float64
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = d.d.DecodeFloat64()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUint8BoolR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint8]bool)
		if v, changed := fastpathTV.DecMapUint8BoolV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUint8BoolV(rv2i(rv).(map[uint8]bool), false, d)
	}
}
func (f fastpathT) DecMapUint8BoolX(vp *map[uint8]bool, d *Decoder) {
	if v, changed := f.DecMapUint8BoolV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUint8BoolV(v map[uint8]bool, canChange bool,
	d *Decoder) (_ map[uint8]bool, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint8]bool, decInferLen(containerLen, d.h.MaxInitLen, 2))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint8
	var mv bool
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = false
			}
			continue
		}
		mv = d.d.DecodeBool()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUint64IntfR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint64]interface{})
		if v, changed := fastpathTV.DecMapUint64IntfV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUint64IntfV(rv2i(rv).(map[uint64]interface{}), false, d)
	}
}
func (f fastpathT) DecMapUint64IntfX(vp *map[uint64]interface{}, d *Decoder) {
	if v, changed := f.DecMapUint64IntfV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUint64IntfV(v map[uint64]interface{}, canChange bool,
	d *Decoder) (_ map[uint64]interface{}, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint64]interface{}, decInferLen(containerLen, d.h.MaxInitLen, 24))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	mapGet := v != nil && !d.h.MapValueReset && !d.h.InterfaceReset
	var mk uint64
	var mv interface{}
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeUint64()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = nil
			}
			continue
		}
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		d.decode(&mv)
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUint64StringR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint64]string)
		if v, changed := fastpathTV.DecMapUint64StringV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUint64StringV(rv2i(rv).(map[uint64]string), false, d)
	}
}
func (f fastpathT) DecMapUint64StringX(vp *map[uint64]string, d *Decoder) {
	if v, changed := f.DecMapUint64StringV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUint64StringV(v map[uint64]string, canChange bool,
	d *Decoder) (_ map[uint64]string, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint64]string, decInferLen(containerLen, d.h.MaxInitLen, 24))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint64
	var mv string
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeUint64()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = ""
			}
			continue
		}
		mv = d.d.DecodeString()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUint64BytesR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint64][]byte)
		if v, changed := fastpathTV.DecMapUint64BytesV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUint64BytesV(rv2i(rv).(map[uint64][]byte), false, d)
	}
}
func (f fastpathT) DecMapUint64BytesX(vp *map[uint64][]byte, d *Decoder) {
	if v, changed := f.DecMapUint64BytesV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUint64BytesV(v map[uint64][]byte, canChange bool,
	d *Decoder) (_ map[uint64][]byte, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint64][]byte, decInferLen(containerLen, d.h.MaxInitLen, 32))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	mapGet := v != nil && !d.h.MapValueReset
	var mk uint64
	var mv []byte
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeUint64()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = nil
			}
			continue
		}
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		mv = d.d.DecodeBytes(mv, false)
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUint64UintR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint64]uint)
		if v, changed := fastpathTV.DecMapUint64UintV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUint64UintV(rv2i(rv).(map[uint64]uint), false, d)
	}
}
func (f fastpathT) DecMapUint64UintX(vp *map[uint64]uint, d *Decoder) {
	if v, changed := f.DecMapUint64UintV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUint64UintV(v map[uint64]uint, canChange bool,
	d *Decoder) (_ map[uint64]uint, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint64]uint, decInferLen(containerLen, d.h.MaxInitLen, 16))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint64
	var mv uint
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeUint64()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUint64Uint8R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint64]uint8)
		if v, changed := fastpathTV.DecMapUint64Uint8V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUint64Uint8V(rv2i(rv).(map[uint64]uint8), false, d)
	}
}
func (f fastpathT) DecMapUint64Uint8X(vp *map[uint64]uint8, d *Decoder) {
	if v, changed := f.DecMapUint64Uint8V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUint64Uint8V(v map[uint64]uint8, canChange bool,
	d *Decoder) (_ map[uint64]uint8, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint64]uint8, decInferLen(containerLen, d.h.MaxInitLen, 9))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint64
	var mv uint8
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeUint64()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUint64Uint64R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint64]uint64)
		if v, changed := fastpathTV.DecMapUint64Uint64V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUint64Uint64V(rv2i(rv).(map[uint64]uint64), false, d)
	}
}
func (f fastpathT) DecMapUint64Uint64X(vp *map[uint64]uint64, d *Decoder) {
	if v, changed := f.DecMapUint64Uint64V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUint64Uint64V(v map[uint64]uint64, canChange bool,
	d *Decoder) (_ map[uint64]uint64, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint64]uint64, decInferLen(containerLen, d.h.MaxInitLen, 16))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint64
	var mv uint64
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeUint64()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = d.d.DecodeUint64()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUint64UintptrR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint64]uintptr)
		if v, changed := fastpathTV.DecMapUint64UintptrV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUint64UintptrV(rv2i(rv).(map[uint64]uintptr), false, d)
	}
}
func (f fastpathT) DecMapUint64UintptrX(vp *map[uint64]uintptr, d *Decoder) {
	if v, changed := f.DecMapUint64UintptrV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUint64UintptrV(v map[uint64]uintptr, canChange bool,
	d *Decoder) (_ map[uint64]uintptr, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint64]uintptr, decInferLen(containerLen, d.h.MaxInitLen, 16))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint64
	var mv uintptr
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeUint64()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUint64IntR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint64]int)
		if v, changed := fastpathTV.DecMapUint64IntV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUint64IntV(rv2i(rv).(map[uint64]int), false, d)
	}
}
func (f fastpathT) DecMapUint64IntX(vp *map[uint64]int, d *Decoder) {
	if v, changed := f.DecMapUint64IntV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUint64IntV(v map[uint64]int, canChange bool,
	d *Decoder) (_ map[uint64]int, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint64]int, decInferLen(containerLen, d.h.MaxInitLen, 16))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint64
	var mv int
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeUint64()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUint64Int64R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint64]int64)
		if v, changed := fastpathTV.DecMapUint64Int64V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUint64Int64V(rv2i(rv).(map[uint64]int64), false, d)
	}
}
func (f fastpathT) DecMapUint64Int64X(vp *map[uint64]int64, d *Decoder) {
	if v, changed := f.DecMapUint64Int64V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUint64Int64V(v map[uint64]int64, canChange bool,
	d *Decoder) (_ map[uint64]int64, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint64]int64, decInferLen(containerLen, d.h.MaxInitLen, 16))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint64
	var mv int64
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeUint64()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = d.d.DecodeInt64()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUint64Float32R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint64]float32)
		if v, changed := fastpathTV.DecMapUint64Float32V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUint64Float32V(rv2i(rv).(map[uint64]float32), false, d)
	}
}
func (f fastpathT) DecMapUint64Float32X(vp *map[uint64]float32, d *Decoder) {
	if v, changed := f.DecMapUint64Float32V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUint64Float32V(v map[uint64]float32, canChange bool,
	d *Decoder) (_ map[uint64]float32, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint64]float32, decInferLen(containerLen, d.h.MaxInitLen, 12))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint64
	var mv float32
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeUint64()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = float32(d.decodeFloat32())
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUint64Float64R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint64]float64)
		if v, changed := fastpathTV.DecMapUint64Float64V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUint64Float64V(rv2i(rv).(map[uint64]float64), false, d)
	}
}
func (f fastpathT) DecMapUint64Float64X(vp *map[uint64]float64, d *Decoder) {
	if v, changed := f.DecMapUint64Float64V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUint64Float64V(v map[uint64]float64, canChange bool,
	d *Decoder) (_ map[uint64]float64, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint64]float64, decInferLen(containerLen, d.h.MaxInitLen, 16))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint64
	var mv float64
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeUint64()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = d.d.DecodeFloat64()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapUint64BoolR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[uint64]bool)
		if v, changed := fastpathTV.DecMapUint64BoolV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapUint64BoolV(rv2i(rv).(map[uint64]bool), false, d)
	}
}
func (f fastpathT) DecMapUint64BoolX(vp *map[uint64]bool, d *Decoder) {
	if v, changed := f.DecMapUint64BoolV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapUint64BoolV(v map[uint64]bool, canChange bool,
	d *Decoder) (_ map[uint64]bool, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[uint64]bool, decInferLen(containerLen, d.h.MaxInitLen, 9))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk uint64
	var mv bool
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeUint64()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = false
			}
			continue
		}
		mv = d.d.DecodeBool()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapIntIntfR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[int]interface{})
		if v, changed := fastpathTV.DecMapIntIntfV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapIntIntfV(rv2i(rv).(map[int]interface{}), false, d)
	}
}
func (f fastpathT) DecMapIntIntfX(vp *map[int]interface{}, d *Decoder) {
	if v, changed := f.DecMapIntIntfV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapIntIntfV(v map[int]interface{}, canChange bool,
	d *Decoder) (_ map[int]interface{}, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[int]interface{}, decInferLen(containerLen, d.h.MaxInitLen, 24))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	mapGet := v != nil && !d.h.MapValueReset && !d.h.InterfaceReset
	var mk int
	var mv interface{}
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = nil
			}
			continue
		}
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		d.decode(&mv)
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapIntStringR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[int]string)
		if v, changed := fastpathTV.DecMapIntStringV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapIntStringV(rv2i(rv).(map[int]string), false, d)
	}
}
func (f fastpathT) DecMapIntStringX(vp *map[int]string, d *Decoder) {
	if v, changed := f.DecMapIntStringV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapIntStringV(v map[int]string, canChange bool,
	d *Decoder) (_ map[int]string, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[int]string, decInferLen(containerLen, d.h.MaxInitLen, 24))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk int
	var mv string
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = ""
			}
			continue
		}
		mv = d.d.DecodeString()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapIntBytesR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[int][]byte)
		if v, changed := fastpathTV.DecMapIntBytesV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapIntBytesV(rv2i(rv).(map[int][]byte), false, d)
	}
}
func (f fastpathT) DecMapIntBytesX(vp *map[int][]byte, d *Decoder) {
	if v, changed := f.DecMapIntBytesV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapIntBytesV(v map[int][]byte, canChange bool,
	d *Decoder) (_ map[int][]byte, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[int][]byte, decInferLen(containerLen, d.h.MaxInitLen, 32))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	mapGet := v != nil && !d.h.MapValueReset
	var mk int
	var mv []byte
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = nil
			}
			continue
		}
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		mv = d.d.DecodeBytes(mv, false)
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapIntUintR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[int]uint)
		if v, changed := fastpathTV.DecMapIntUintV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapIntUintV(rv2i(rv).(map[int]uint), false, d)
	}
}
func (f fastpathT) DecMapIntUintX(vp *map[int]uint, d *Decoder) {
	if v, changed := f.DecMapIntUintV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapIntUintV(v map[int]uint, canChange bool,
	d *Decoder) (_ map[int]uint, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[int]uint, decInferLen(containerLen, d.h.MaxInitLen, 16))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk int
	var mv uint
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapIntUint8R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[int]uint8)
		if v, changed := fastpathTV.DecMapIntUint8V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapIntUint8V(rv2i(rv).(map[int]uint8), false, d)
	}
}
func (f fastpathT) DecMapIntUint8X(vp *map[int]uint8, d *Decoder) {
	if v, changed := f.DecMapIntUint8V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapIntUint8V(v map[int]uint8, canChange bool,
	d *Decoder) (_ map[int]uint8, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[int]uint8, decInferLen(containerLen, d.h.MaxInitLen, 9))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk int
	var mv uint8
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapIntUint64R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[int]uint64)
		if v, changed := fastpathTV.DecMapIntUint64V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapIntUint64V(rv2i(rv).(map[int]uint64), false, d)
	}
}
func (f fastpathT) DecMapIntUint64X(vp *map[int]uint64, d *Decoder) {
	if v, changed := f.DecMapIntUint64V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapIntUint64V(v map[int]uint64, canChange bool,
	d *Decoder) (_ map[int]uint64, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[int]uint64, decInferLen(containerLen, d.h.MaxInitLen, 16))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk int
	var mv uint64
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = d.d.DecodeUint64()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapIntUintptrR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[int]uintptr)
		if v, changed := fastpathTV.DecMapIntUintptrV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapIntUintptrV(rv2i(rv).(map[int]uintptr), false, d)
	}
}
func (f fastpathT) DecMapIntUintptrX(vp *map[int]uintptr, d *Decoder) {
	if v, changed := f.DecMapIntUintptrV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapIntUintptrV(v map[int]uintptr, canChange bool,
	d *Decoder) (_ map[int]uintptr, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[int]uintptr, decInferLen(containerLen, d.h.MaxInitLen, 16))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk int
	var mv uintptr
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapIntIntR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[int]int)
		if v, changed := fastpathTV.DecMapIntIntV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapIntIntV(rv2i(rv).(map[int]int), false, d)
	}
}
func (f fastpathT) DecMapIntIntX(vp *map[int]int, d *Decoder) {
	if v, changed := f.DecMapIntIntV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapIntIntV(v map[int]int, canChange bool,
	d *Decoder) (_ map[int]int, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[int]int, decInferLen(containerLen, d.h.MaxInitLen, 16))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk int
	var mv int
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapIntInt64R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[int]int64)
		if v, changed := fastpathTV.DecMapIntInt64V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapIntInt64V(rv2i(rv).(map[int]int64), false, d)
	}
}
func (f fastpathT) DecMapIntInt64X(vp *map[int]int64, d *Decoder) {
	if v, changed := f.DecMapIntInt64V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapIntInt64V(v map[int]int64, canChange bool,
	d *Decoder) (_ map[int]int64, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[int]int64, decInferLen(containerLen, d.h.MaxInitLen, 16))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk int
	var mv int64
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = d.d.DecodeInt64()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapIntFloat32R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[int]float32)
		if v, changed := fastpathTV.DecMapIntFloat32V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapIntFloat32V(rv2i(rv).(map[int]float32), false, d)
	}
}
func (f fastpathT) DecMapIntFloat32X(vp *map[int]float32, d *Decoder) {
	if v, changed := f.DecMapIntFloat32V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapIntFloat32V(v map[int]float32, canChange bool,
	d *Decoder) (_ map[int]float32, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[int]float32, decInferLen(containerLen, d.h.MaxInitLen, 12))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk int
	var mv float32
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = float32(d.decodeFloat32())
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapIntFloat64R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[int]float64)
		if v, changed := fastpathTV.DecMapIntFloat64V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapIntFloat64V(rv2i(rv).(map[int]float64), false, d)
	}
}
func (f fastpathT) DecMapIntFloat64X(vp *map[int]float64, d *Decoder) {
	if v, changed := f.DecMapIntFloat64V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapIntFloat64V(v map[int]float64, canChange bool,
	d *Decoder) (_ map[int]float64, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[int]float64, decInferLen(containerLen, d.h.MaxInitLen, 16))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk int
	var mv float64
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = d.d.DecodeFloat64()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapIntBoolR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[int]bool)
		if v, changed := fastpathTV.DecMapIntBoolV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapIntBoolV(rv2i(rv).(map[int]bool), false, d)
	}
}
func (f fastpathT) DecMapIntBoolX(vp *map[int]bool, d *Decoder) {
	if v, changed := f.DecMapIntBoolV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapIntBoolV(v map[int]bool, canChange bool,
	d *Decoder) (_ map[int]bool, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[int]bool, decInferLen(containerLen, d.h.MaxInitLen, 9))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk int
	var mv bool
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = false
			}
			continue
		}
		mv = d.d.DecodeBool()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapInt64IntfR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[int64]interface{})
		if v, changed := fastpathTV.DecMapInt64IntfV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapInt64IntfV(rv2i(rv).(map[int64]interface{}), false, d)
	}
}
func (f fastpathT) DecMapInt64IntfX(vp *map[int64]interface{}, d *Decoder) {
	if v, changed := f.DecMapInt64IntfV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapInt64IntfV(v map[int64]interface{}, canChange bool,
	d *Decoder) (_ map[int64]interface{}, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[int64]interface{}, decInferLen(containerLen, d.h.MaxInitLen, 24))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	mapGet := v != nil && !d.h.MapValueReset && !d.h.InterfaceReset
	var mk int64
	var mv interface{}
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeInt64()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = nil
			}
			continue
		}
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		d.decode(&mv)
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapInt64StringR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[int64]string)
		if v, changed := fastpathTV.DecMapInt64StringV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapInt64StringV(rv2i(rv).(map[int64]string), false, d)
	}
}
func (f fastpathT) DecMapInt64StringX(vp *map[int64]string, d *Decoder) {
	if v, changed := f.DecMapInt64StringV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapInt64StringV(v map[int64]string, canChange bool,
	d *Decoder) (_ map[int64]string, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[int64]string, decInferLen(containerLen, d.h.MaxInitLen, 24))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk int64
	var mv string
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeInt64()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = ""
			}
			continue
		}
		mv = d.d.DecodeString()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapInt64BytesR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[int64][]byte)
		if v, changed := fastpathTV.DecMapInt64BytesV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapInt64BytesV(rv2i(rv).(map[int64][]byte), false, d)
	}
}
func (f fastpathT) DecMapInt64BytesX(vp *map[int64][]byte, d *Decoder) {
	if v, changed := f.DecMapInt64BytesV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapInt64BytesV(v map[int64][]byte, canChange bool,
	d *Decoder) (_ map[int64][]byte, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[int64][]byte, decInferLen(containerLen, d.h.MaxInitLen, 32))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	mapGet := v != nil && !d.h.MapValueReset
	var mk int64
	var mv []byte
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeInt64()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = nil
			}
			continue
		}
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		mv = d.d.DecodeBytes(mv, false)
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapInt64UintR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[int64]uint)
		if v, changed := fastpathTV.DecMapInt64UintV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapInt64UintV(rv2i(rv).(map[int64]uint), false, d)
	}
}
func (f fastpathT) DecMapInt64UintX(vp *map[int64]uint, d *Decoder) {
	if v, changed := f.DecMapInt64UintV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapInt64UintV(v map[int64]uint, canChange bool,
	d *Decoder) (_ map[int64]uint, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[int64]uint, decInferLen(containerLen, d.h.MaxInitLen, 16))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk int64
	var mv uint
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeInt64()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapInt64Uint8R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[int64]uint8)
		if v, changed := fastpathTV.DecMapInt64Uint8V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapInt64Uint8V(rv2i(rv).(map[int64]uint8), false, d)
	}
}
func (f fastpathT) DecMapInt64Uint8X(vp *map[int64]uint8, d *Decoder) {
	if v, changed := f.DecMapInt64Uint8V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapInt64Uint8V(v map[int64]uint8, canChange bool,
	d *Decoder) (_ map[int64]uint8, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[int64]uint8, decInferLen(containerLen, d.h.MaxInitLen, 9))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk int64
	var mv uint8
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeInt64()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapInt64Uint64R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[int64]uint64)
		if v, changed := fastpathTV.DecMapInt64Uint64V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapInt64Uint64V(rv2i(rv).(map[int64]uint64), false, d)
	}
}
func (f fastpathT) DecMapInt64Uint64X(vp *map[int64]uint64, d *Decoder) {
	if v, changed := f.DecMapInt64Uint64V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapInt64Uint64V(v map[int64]uint64, canChange bool,
	d *Decoder) (_ map[int64]uint64, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[int64]uint64, decInferLen(containerLen, d.h.MaxInitLen, 16))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk int64
	var mv uint64
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeInt64()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = d.d.DecodeUint64()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapInt64UintptrR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[int64]uintptr)
		if v, changed := fastpathTV.DecMapInt64UintptrV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapInt64UintptrV(rv2i(rv).(map[int64]uintptr), false, d)
	}
}
func (f fastpathT) DecMapInt64UintptrX(vp *map[int64]uintptr, d *Decoder) {
	if v, changed := f.DecMapInt64UintptrV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapInt64UintptrV(v map[int64]uintptr, canChange bool,
	d *Decoder) (_ map[int64]uintptr, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[int64]uintptr, decInferLen(containerLen, d.h.MaxInitLen, 16))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk int64
	var mv uintptr
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeInt64()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapInt64IntR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[int64]int)
		if v, changed := fastpathTV.DecMapInt64IntV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapInt64IntV(rv2i(rv).(map[int64]int), false, d)
	}
}
func (f fastpathT) DecMapInt64IntX(vp *map[int64]int, d *Decoder) {
	if v, changed := f.DecMapInt64IntV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapInt64IntV(v map[int64]int, canChange bool,
	d *Decoder) (_ map[int64]int, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[int64]int, decInferLen(containerLen, d.h.MaxInitLen, 16))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk int64
	var mv int
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeInt64()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapInt64Int64R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[int64]int64)
		if v, changed := fastpathTV.DecMapInt64Int64V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapInt64Int64V(rv2i(rv).(map[int64]int64), false, d)
	}
}
func (f fastpathT) DecMapInt64Int64X(vp *map[int64]int64, d *Decoder) {
	if v, changed := f.DecMapInt64Int64V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapInt64Int64V(v map[int64]int64, canChange bool,
	d *Decoder) (_ map[int64]int64, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[int64]int64, decInferLen(containerLen, d.h.MaxInitLen, 16))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk int64
	var mv int64
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeInt64()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = d.d.DecodeInt64()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapInt64Float32R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[int64]float32)
		if v, changed := fastpathTV.DecMapInt64Float32V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapInt64Float32V(rv2i(rv).(map[int64]float32), false, d)
	}
}
func (f fastpathT) DecMapInt64Float32X(vp *map[int64]float32, d *Decoder) {
	if v, changed := f.DecMapInt64Float32V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapInt64Float32V(v map[int64]float32, canChange bool,
	d *Decoder) (_ map[int64]float32, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[int64]float32, decInferLen(containerLen, d.h.MaxInitLen, 12))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk int64
	var mv float32
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeInt64()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = float32(d.decodeFloat32())
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapInt64Float64R(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[int64]float64)
		if v, changed := fastpathTV.DecMapInt64Float64V(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapInt64Float64V(rv2i(rv).(map[int64]float64), false, d)
	}
}
func (f fastpathT) DecMapInt64Float64X(vp *map[int64]float64, d *Decoder) {
	if v, changed := f.DecMapInt64Float64V(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapInt64Float64V(v map[int64]float64, canChange bool,
	d *Decoder) (_ map[int64]float64, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[int64]float64, decInferLen(containerLen, d.h.MaxInitLen, 16))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk int64
	var mv float64
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeInt64()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = 0
			}
			continue
		}
		mv = d.d.DecodeFloat64()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}

func (d *Decoder) fastpathDecMapInt64BoolR(f *codecFnInfo, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		vp := rv2i(rv).(*map[int64]bool)
		if v, changed := fastpathTV.DecMapInt64BoolV(*vp, true, d); changed {
			*vp = v
		}
	} else {
		fastpathTV.DecMapInt64BoolV(rv2i(rv).(map[int64]bool), false, d)
	}
}
func (f fastpathT) DecMapInt64BoolX(vp *map[int64]bool, d *Decoder) {
	if v, changed := f.DecMapInt64BoolV(*vp, true, d); changed {
		*vp = v
	}
}
func (fastpathT) DecMapInt64BoolV(v map[int64]bool, canChange bool,
	d *Decoder) (_ map[int64]bool, changed bool) {
	containerLen := d.mapStart()
	if canChange && v == nil {
		v = make(map[int64]bool, decInferLen(containerLen, d.h.MaxInitLen, 9))
		changed = true
	}
	if containerLen == 0 {
		d.mapEnd()
		return v, changed
	}
	var mk int64
	var mv bool
	hasLen := containerLen > 0
	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.d.CheckBreak()); j++ {
		d.mapElemKey()
		mk = d.d.DecodeInt64()
		d.mapElemValue()
		if d.d.TryDecodeAsNil() {
			if v == nil {
			} else if d.h.DeleteOnNilMapValue {
				delete(v, mk)
			} else {
				v[mk] = false
			}
			continue
		}
		mv = d.d.DecodeBool()
		if v != nil {
			v[mk] = mv
		}
	}
	d.mapEnd()
	return v, changed
}
