package assets

import (
	_ "embed"
)

// wat2wasm memory.wat --output=memory.wasm --enable-threads
//
//go:embed memory.wasm
var MemoryWasm []byte
