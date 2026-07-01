//go:build js

// Package jsx provides essential JavaScript functions that are used
// widely in wgpu and are very useful for an wasm / js application.
package jsx

import (
	"log/slog"
	"syscall/js"
	"unsafe"
)

// BytesToJS converts Go bytes to a JS Uint8Array for queue.writeBuffer/writeTexture.
// Uses zero-copy wasm linear memory when a host exposes window.wasm.instance.exports.mem;
// otherwise copies via js.CopyBytesToJS (standard Go wasm_exec.js).
func BytesToJS(b []byte) js.Value {
	if len(b) == 0 {
		return js.Global().Get("Uint8Array").New(0)
	}

	wasm := js.Global().Get("wasm")
	if wasm.Truthy() {
		mem := wasm.Get("instance").Get("exports").Get("mem")
		if mem.Truthy() {
			ptr := uintptr(unsafe.Pointer(&b[0]))
			return js.Global().Get("Uint8Array").New(mem.Get("buffer"), ptr, len(b))
		}
	}

	arr := js.Global().Get("Uint8Array").New(len(b))
	js.CopyBytesToJS(arr, b)
	return arr
}

// Await is a helper function equivalent to await in JS.
// It is copied from https://go-review.googlesource.com/c/go/+/150917/
func Await(promise js.Value) (result js.Value, ok bool) {
	if promise.Type() != js.TypeObject || promise.Get("then").Type() != js.TypeFunction {
		return promise, true
	}

	done := make(chan struct{})

	onResolve := js.FuncOf(func(this js.Value, args []js.Value) any {
		result = args[0]
		ok = true
		close(done)
		return nil
	})
	defer onResolve.Release()

	onReject := js.FuncOf(func(this js.Value, args []js.Value) any {
		result = args[0]
		ok = false
		slog.Error("wgpu.AwaitJS: promise rejected", "reason", result)
		close(done)
		return nil
	})
	defer onReject.Release()

	promise.Call("then", onResolve, onReject)
	<-done
	return
}
