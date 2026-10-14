#!/bin/sh
# Cross-compiles ycast-go for the targets documented in README.md.
# Usage: ./build.sh [version]   (version defaults to git describe, or "dev")
set -eu

version=${1:-}
if [ -z "$version" ]; then
    version=$(git describe --tags --always --dirty 2>/dev/null || echo dev)
fi

mkdir -p dist
rm -f dist/ycast-go-*

build() { # build <goos> <goarch> [goarm]
    os=$1 arch=$2 arm=${3:-}
    suffix=${arm:+v$arm}
    out=dist/ycast-go-$os-$arch$suffix
    env CGO_ENABLED=0 GOOS=$os GOARCH=$arch ${arm:+GOARM=$arm} \
        go build -trimpath -ldflags "-X main.version=$version -s -w" -o "$out" ./cmd/ycast-go
    echo "$out"
}

# 64-bit: PC/NAS (amd64 v1 for maximum compatibility), Pi 3+/4/5, ARM servers.
build linux amd64
build linux arm64
# 32-bit ARM: armv7 for modern Pi running armhf, armv6 for Pi Zero/1 (runs everywhere).
build linux arm 7
build linux arm 6
# MIPS routers (OpenWrt): little-endian for most, big-endian for the rest.
build linux mipsle
build linux mips
