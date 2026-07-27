package hashutils

import (
	"unsafe"

	"github.com/zeebo/xxh3"
)

func HashHashes(data []uint64) uint64 {
	ptr := unsafe.SliceData(data)

	return xxh3.Hash(unsafe.Slice((*byte)(unsafe.Pointer(ptr)), len(data)*8))
}
