//go:build js

package wgpu

import "strings"

// textureFormatToJS converts wgpu-native hyphenated format names (e.g. "bgra8-unorm")
// to browser GPUTextureFormat strings (e.g. "bgra8unorm").
func textureFormatToJS(f TextureFormat) any {
	s := f.String()
	if s == "undefined" {
		return enumToJS(f)
	}
	return wgpuFormatStringToJS(s)
}

func wgpuFormatStringToJS(s string) string {
	const stencil = "-stencil8"
	const srgb = "-srgb"
	if strings.HasSuffix(s, stencil) {
		return wgpuFormatStringToJS(strings.TrimSuffix(s, stencil)) + stencil
	}
	if strings.HasSuffix(s, srgb) {
		return wgpuFormatStringToJS(strings.TrimSuffix(s, srgb)) + srgb
	}
	return strings.ReplaceAll(s, "-", "")
}
