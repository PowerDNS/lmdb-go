package lmdb

import (
	"testing"
	"unsafe"
)

func TestTxnStructSize(t *testing.T) {
	t.Logf("sizeof(Txn) = %d", unsafe.Sizeof(Txn{}))
	if unsafe.Sizeof(Txn{}) > 64 {
		t.Errorf("Txn struct grew beyond 64 bytes: %d", unsafe.Sizeof(Txn{}))
	}
	t.Logf("sizeof(Cursor) = %d", unsafe.Sizeof(Cursor{}))
	if unsafe.Sizeof(Cursor{}) > 16 {
		t.Errorf("Cursor struct grew beyond 16 bytes: %d", unsafe.Sizeof(Cursor{}))
	}
}
