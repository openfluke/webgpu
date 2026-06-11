# WebGPU (openfluke)

Go bindings for [WebGPU](https://gpuweb.github.io/gpuweb/), maintained by [openfluke](https://github.com/openfluke).

Native builds use **wgpu-native v29** with Go-side bindings (no C compatibility shims). The same module also supports the browser via WASM. Targets include Vulkan, Metal, D3D12, and OpenGL ES.

**Module:** `github.com/openfluke/webgpu`  
**Current release:** `v1.0.4` (ships wgpu-native **v29.0.0.0**)

## Version history

| Tag | What it was |
|-----|-------------|
| `v1.0.0` | Initial openfluke release: v29 Go bindings + v29 static libs |
| `v1.0.2` | Do not use — Go proxy cached wrong commit |
| `v1.0.3` | Do not use — Go proxy cached wrong commit |
| **`v1.0.4`** | Drops legacy `ios/amd64`; fixes `go get` module size limit |

## What's in v1.0.0

### wgpu-native v29 stack

- **Headers:** `wgpu/lib/webgpu.h`, `wgpu/lib/wgpu.h` (WebGPU C API + wgpu-native extensions)
- **Binaries:** prebuilt `libwgpu_native.a` per target under `wgpu/lib/`:

| Platform | Path |
|----------|------|
| macOS | `darwin/amd64`, `darwin/arm64` |
| Linux | `linux/amd64`, `linux/arm64` |
| Windows | `windows/amd64` (gnu), `windows/arm64` (msvc static) |
| iOS | `ios/arm64` (device and Apple Silicon simulator) |
| Android | `android/arm64`, `android/arm`, `android/amd64`, `android/386` |

### Go binding highlights (native)

Bindings target the v29 C API directly:

- `WGPUStringView` for labels, entry points, shader source, messages
- Async ops via `WGPUFuture` + `wgpuInstanceWaitAny` (`RequestAdapter`, `RequestDevice`, `MapAsync`, error scopes, queue work done)
- `wgpuDeviceAddRef`, `wgpuAdapterGetFeatures`, `wgpuDeviceGetLimits` with `maxImmediateSize`
- Shader chains: `WGPUShaderSourceWGSL` / `SPIRV` / `WGPUShaderSourceGLSL`
- Surfaces: `WGPUSurfaceSource*` descriptors
- Texel copy types: `WGPUTexelCopyTextureInfo`, etc.
- Push constants → pipeline `immediateDataSize` / `SetImmediates` (native extension)
- Validation error scopes implemented in Go (`v29.go`)

### Tested on real hardware

- **macOS arm64 + Metal** — loom / lucy: adapter → device → buffer → GPU forward parity (Apple M5)
- **Windows arm64** — native wgpu-native v29 static lib; GPU path verified

Other platforms have matching v29 libs in-tree; smoke-test on each target you ship.

## Install

```bash
go get github.com/openfluke/webgpu@v1.0.4
```

Local development (e.g. from the endgame monorepo):

```go
replace github.com/openfluke/webgpu => ../webgpu
```

## Usage

Native and JS builds use the same import path; build tags select the backend:

```go
import "github.com/openfluke/webgpu/wgpu"
```

- `//go:build !js` — wgpu-native (this README)
- `//go:build js` — browser WebGPU via WASM

## Examples

| [boids](examples/boids) | [cube](examples/cube) | [triangle](examples/triangle) |
|:---:|:---:|:---:|
| ![boids](https://raw.githubusercontent.com/rajveermalviya/go-webgpu/main/tests/boids/image-msaa.png) | ![cube](https://raw.githubusercontent.com/rajveermalviya/go-webgpu/main/tests/cube/image-msaa.png) | ![triangle](https://raw.githubusercontent.com/rajveermalviya/go-webgpu/main/tests/triangle/image-msaa.png) |

## Rebuilding native libs

Vendored artifacts came from official **wgpu-native v29** release zips. To refresh:

1. Download matching `wgpu-*-release.zip` from [wgpu-native releases](https://github.com/gfx-rs/wgpu-native/releases)
2. Copy `lib/libwgpu_native.a` into the corresponding `wgpu/lib/<platform>/<arch>/` directory
3. Ensure shared headers under `wgpu/lib/` stay in sync with that release

CI can also rebuild via [.github/workflows/build-wgpu.yml](.github/workflows/build-wgpu.yml) (`workflow_dispatch`).

## Lineage

API design and WASM support draw on earlier community work ([rajveermalviya/go-webgpu](https://github.com/rajveermalviya/go-webgpu), [mokiat/wasmgpu](https://github.com/mokiat/wasmgpu)). This repository is maintained independently by openfluke; vendored binaries and the v29 native bindings are openfluke-specific.

## References

- [WebGPU](https://gpuweb.github.io/gpuweb/)
- [WGSL](https://gpuweb.github.io/gpuweb/wgsl/)
- [webgpu.h](https://github.com/webgpu-native/webgpu-headers)
- [wgpu-native](https://github.com/gfx-rs/wgpu-native)

## License

See [LICENSE](LICENSE).
