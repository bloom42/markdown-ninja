//go:build !unix && !windows

package wazeroallocator

import (
	"sync"

	"github.com/tetratelabs/wazero/experimental"
)

// Separate implementation of non-Unix/Windows code to file without
// build tag to allow testing on any platform.

var pageSize = uint64(0) // used only for test

type sliceBuffer struct {
	buf   []byte
	mutex sync.Mutex
}

func alloc(cap, max uint64) experimental.LinearMemory {
	buf := make([]byte, max)
	return &sliceBuffer{buf: buf[:0], mutex: sync.Mutex{}}
}

func (b *sliceBuffer) Free() {}

func (b *sliceBuffer) Reallocate(size uint64) []byte {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	return b.buf[:size]
}
