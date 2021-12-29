#!/bin/bash

function build {
    local os=$1
    local arch=$2
    local ext=$3
    env GOOS="$os" GOARCH="$arch" go build -o "bronco-build-tracker-$os-$arch$ext" ./cmd/bronco-build-tracker-cli
    chmod +x "bronco-build-tracker-$os-$arch$ext"
}

build darwin amd64
build linux arm
build linux amd64
build windows amd64 ".exe"