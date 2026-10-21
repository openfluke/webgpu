# Rebuild wgpu/lib/windows/arm64/libwgpu_native.a as aarch64-pc-windows-gnullvm
# so Go+CGO (llvm-mingw) can link — upstream only ships MSVC for this arch.
#
# Prerequisites:
#   - rustup + rust-std for aarch64-pc-windows-gnullvm (matching wgpu-native toolchain)
#   - llvm-mingw on PATH (aarch64-w64-mingw32-clang)
#   - LLVM with libclang.dll (LIBCLANG_PATH), e.g. winget install LLVM.LLVM
#
# Usage (from repo root or anywhere):
#   .\scripts\build_windows_arm64_gnu.ps1
#   .\scripts\build_windows_arm64_gnu.ps1 -Tag v29.0.0.0

param(
    [string]$Tag = "v29.0.0.0",
    [string]$LlvmMingw = $(if ($env:LLVM_MINGW_HOME) { $env:LLVM_MINGW_HOME } else { "C:\llvm-mingw" })
)

$ErrorActionPreference = "Stop"
$RepoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$BuildRoot = Join-Path $RepoRoot "_build"
$SrcDir = Join-Path $BuildRoot "wgpu-native"
$TargetDir = Join-Path $BuildRoot "target"
$Clang = Join-Path $LlvmMingw "bin\aarch64-w64-mingw32-clang.exe"
$Dest = Join-Path $RepoRoot "wgpu\lib\windows\arm64\libwgpu_native.a"

if (-not (Test-Path $Clang)) {
    throw "llvm-mingw clang not found at $Clang (set LLVM_MINGW_HOME)"
}
if (-not $env:LIBCLANG_PATH) {
    $guess = "C:\Program Files\LLVM\bin"
    if (Test-Path (Join-Path $guess "libclang.dll")) {
        $env:LIBCLANG_PATH = $guess
    } else {
        throw "Set LIBCLANG_PATH to a directory containing libclang.dll"
    }
}

New-Item -ItemType Directory -Force -Path $BuildRoot | Out-Null
if (-not (Test-Path (Join-Path $SrcDir "Cargo.toml"))) {
    git clone --recursive --branch $Tag --depth 1 https://github.com/gfx-rs/wgpu-native.git $SrcDir
}

$cargoConfig = Join-Path $SrcDir ".cargo"
New-Item -ItemType Directory -Force -Path $cargoConfig | Out-Null
@"
[target.aarch64-pc-windows-gnullvm]
linker = "$($Clang -replace '\\','/')"
ar = "$($LlvmMingw -replace '\\','/')/bin/llvm-ar.exe"

[profile.release]
strip = "symbols"
"@ | Set-Content -Path (Join-Path $cargoConfig "config.toml") -Encoding ascii

$env:PATH = "$(Join-Path $LlvmMingw 'bin');$env:PATH"
$env:CC = "aarch64-w64-mingw32-clang"
$env:CXX = "aarch64-w64-mingw32-clang++"
$env:CARGO_TARGET_DIR = $TargetDir
$env:BINDGEN_EXTRA_CLANG_ARGS = "--target=aarch64-w64-windows-gnu"

Push-Location $SrcDir
try {
    # Prefer toolchain pin from wgpu-native if present
    cargo build --release --target aarch64-pc-windows-gnullvm
} finally {
    Pop-Location
}

$Built = Join-Path $TargetDir "aarch64-pc-windows-gnullvm\release\libwgpu_native.a"
if (-not (Test-Path $Built)) {
    throw "build succeeded but $Built missing"
}

# Refuse to install an MSVC-ABI archive
$bad = & llvm-nm $Built 2>$null | Select-String '\?\?_7type_info|__CxxFrameHandler3'
if ($bad) {
    throw "built archive still has MSVC RTTI/EH symbols — wrong target?"
}

New-Item -ItemType Directory -Force -Path (Split-Path $Dest) | Out-Null
Copy-Item $Built $Dest -Force
Write-Host "Installed gnu/gnullvm lib -> $Dest"
Write-Host "Size: $((Get-Item $Dest).Length) bytes"
