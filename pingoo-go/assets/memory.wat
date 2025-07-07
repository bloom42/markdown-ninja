;; the min value (here 23) depends on the instantiated WASM modules importing "memory" and their requirements
(module (memory (export "memory") 23 65536 shared))
