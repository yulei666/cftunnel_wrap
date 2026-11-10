#!/usr/bin/env bash
# 交叉编译并打包发布产物到 dist/。
#
# 用法: scripts/build-release.sh <version> [goos/goarch ...]
#   不指定平台时构建 linux/amd64 linux/arm64 darwin/amd64 darwin/arm64。
#
# 产物: dist/cftunnel_<version>_<goos>_<goarch>.tar.gz，内含二进制、cftunnel.sh、
# config.example.yaml、README 和许可证文件。
set -euo pipefail

if [[ $# -lt 1 ]]; then
	echo "usage: $0 <version> [goos/goarch ...]" >&2
	exit 2
fi

version=${1#v}
shift
targets=("$@")
if [[ ${#targets[@]} -eq 0 ]]; then
	targets=(linux/amd64 linux/arm64 darwin/amd64 darwin/arm64)
fi

cd "$(dirname "$0")/.."
mkdir -p dist

for target in "${targets[@]}"; do
	goos=${target%/*}
	goarch=${target#*/}
	name="cftunnel_${version}_${goos}_${goarch}"
	stage="dist/$name"

	echo "==> building $target"
	rm -rf "$stage"
	mkdir -p "$stage"
	CGO_ENABLED=0 GOOS=$goos GOARCH=$goarch go build -trimpath \
		-ldflags "-s -w -X main.version=$version" \
		-o "$stage/cftunnel" .
	cp cftunnel.sh config.example.yaml README.md LICENSE NOTICE THIRD_PARTY_LICENSES.md "$stage/"

	tar -C dist -czf "dist/$name.tar.gz" "$name"
	rm -rf "$stage"
done

(cd dist && shasum -a 256 ./*.tar.gz | sed 's| \./| |' >checksums.txt)
cat dist/checksums.txt
