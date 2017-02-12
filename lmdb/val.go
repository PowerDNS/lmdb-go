package lmdb

/*
#include <stdlib.h>
#include <stdio.h>
#include "lmdb.h"
*/
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/bmatsuo/lmdb-go/internal/lmdbarch"
)

// valSizeBits is the number of bits which constraining the length of the
// single values in an LMDB database, either 32 or 31 depending on the
// platform.  valMaxSize is the largest data size allowed based.  See runtime
// source file malloc.go and the compiler typecheck.go for more information
// about memory limits and array bound limits.
//
//		https://github.com/golang/go/blob/a03bdc3e6bea34abd5077205371e6fb9ef354481/src/runtime/malloc.go#L151-L164
//		https://github.com/golang/go/blob/36a80c5941ec36d9c44d6f3c068d13201e023b5f/src/cmd/compile/internal/gc/typecheck.go#L383
//
// On 64-bit systems, luckily, the value 2^32-1 coincides with the maximum data
// size for LMDB (MAXDATASIZE).
const (
	valSizeBits = lmdbarch.Width64*32 + (1-lmdbarch.Width64)*31
	valMaxSize  = 1<<valSizeBits - 1
)

// Multi is a generic FixedMultiple implementation that can store contiguous
// fixed-width for a configurable width (stride).
//
// Multi values are only useful in databases opened with DupSort|DupFixed.
type Multi struct {
	page   []byte
	stride int
}

// WrapMulti converts a page of contiguous stride-sized values into a Multi.
// WrapMulti panics if len(page) is not a multiple of stride.
//
// WrapMulti is deprecated.  Use the Stride.Multiple method instead.
//
//		_, elem, err := cursor.Get(nil, nil, lmdb.FirstDup)
//		if err != nil {
//			return err
//		}
//		_, page, err := cursor.Get(nil, nil, lmdb.GetMultiple)
//		elems, err := lmdb.
//			Stride(len(elem)).
//			Multiple(page, err)
//		if err != nil {
//			return err
//		}
//
// See mdb_cursor_get and MDB_GET_MULTIPLE.
func WrapMulti(page []byte, stride int) *Multi {
	if len(page)%stride != 0 {
		panic("incongruent arguments")
	}
	return &Multi{page: page, stride: stride}
}

// Stride is the fixed element width for a Multi object
type Stride int

// Multiple wraps page as a Multi.  Multiple returns an error if the input err
// is non-nil, the stride is not a positive value, or if len(page) is not a
// multiple of the stride.
func (s Stride) Multiple(page []byte, err error) (*Multi, error) {
	if err != nil {
		return nil, err
	}
	if s < 1 {
		return nil, fmt.Errorf("invalid stride")
	}
	if len(page)%int(s) != 0 {
		return nil, fmt.Errorf("argument and page stride are not congruent")
	}
	m := &Multi{
		page:   page,
		stride: int(s),
	}
	return m, nil
}

// Vals returns a slice containing the values in m.  The returned slice has
// length m.Len() and each item has length m.Stride().
func (m *Multi) Vals() [][]byte {
	n := m.Len()
	ps := make([][]byte, n)
	for i := 0; i < n; i++ {
		ps[i] = m.Val(i)
	}
	return ps
}

// Val returns the value at index i.  Val panics if i is out of range.
func (m *Multi) Val(i int) []byte {
	off := i * m.stride
	return m.page[off : off+m.stride]
}

// Append returns a new Multi with page data resulting from appending b to
// m.Page().  Put panics if len(b) is not equal to m.Stride()
func (m *Multi) Append(b []byte) *Multi {
	if len(b) != m.stride {
		panic("bad data size")
	}

	return &Multi{
		stride: m.stride,
		page:   append(m.page, b...),
	}
}

// Len returns the number of values in the Multi.
func (m *Multi) Len() int {
	return len(m.page) / m.stride
}

// Stride returns the length of an individual value in the m.
func (m *Multi) Stride() int {
	return m.stride
}

// Size returns the total size of the Multi data and is equal to
//
//		m.Len()*m.Stride()
//
func (m *Multi) Size() int {
	return len(m.page)
}

// Page returns the Multi page data as a raw slice of bytes with length
// m.Size().
func (m *Multi) Page() []byte {
	return m.page[:len(m.page):len(m.page)]
}

var eb = []byte{0}

func valBytes(b []byte) ([]byte, int) {
	if len(b) == 0 {
		return eb, 0
	}
	return b, len(b)
}

func wrapVal(b []byte) *C.MDB_val {
	p, n := valBytes(b)
	return &C.MDB_val{
		mv_data: unsafe.Pointer(&p[0]),
		mv_size: C.size_t(n),
	}
}

func getBytes(val *C.MDB_val) []byte {
	return (*[valMaxSize]byte)(unsafe.Pointer(val.mv_data))[:val.mv_size:val.mv_size]
}

func getBytesCopy(val *C.MDB_val) []byte {
	return C.GoBytes(val.mv_data, C.int(val.mv_size))
}
