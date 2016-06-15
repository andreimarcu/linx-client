#!/bin/bash

if [ -n "$1" ]; then
	version="$1"
else
	version=$(git tag | grep ^v | tail -1 | sed 's/^v//')
fi

echo "Building version ${version}..."

bin_dir="binaries/$version"
mkdir -p "$bin_dir"

build() {
	echo "Building for $1-${2}..."
	env GOOS="$1" GOARCH="$2" go build \
	  -o "binaries/$version/linx-client-v${version}_$1-$2"
}

for os in darwin freebsd openbsd linux windows; do
	for arch in amd64 386; do
		build "$os" "$arch"
	done
done

build linux arm
