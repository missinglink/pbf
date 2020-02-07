#!/bin/bash
set -euo pipefail

# ensure build directory exists
mkdir build

# ensure the compiler exits with status 0
function assert() {
  if [[ $1 != 0 ]]; then
    echo "one or more architectures failed to compile"
    exit $1;
  fi
}

# check the reported file class matches what's expected
function check() {
  actual=$(file -b ${1});
  if [[ "${actual}" != "${2}"* ]]; then
    echo "invalid file architecture: ${1}"
    echo "expected: ${2}"
    echo "actual: ${actual}"
    echo "-${actual}-${2}-"
    exit 1
  fi
}

# if the 'UPX' bindary packer is available, use it
# https://upx.github.io/
function compress() {
  [ -x "$(command -v upx)" ] && upx "${1}"
}

echo "[compile] linux arm";
env GOOS=linux GOARCH=arm go build -ldflags="-s -w" -gcflags=-trimpath="${GOPATH}" -asmflags=-trimpath="${GOPATH}";
assert $?;
chmod +x pbf;
mv pbf build/pbf.linux-arm;
check 'build/pbf.linux-arm' 'ELF 32-bit LSB executable';
compress 'build/pbf.linux-arm';

echo "[compile] linux x64";
env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -gcflags=-trimpath="${GOPATH}" -asmflags=-trimpath="${GOPATH}";
assert $?;
chmod +x pbf;
mv pbf build/pbf.linux-x64;
check 'build/pbf.linux-x64' 'ELF 64-bit LSB executable';
compress 'build/pbf.linux-x64';

echo "[compile] darwin x64";
env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -gcflags=-trimpath="${GOPATH}" -asmflags=-trimpath="${GOPATH}";
assert $?;
chmod +x pbf;
mv pbf build/pbf.darwin-x64;
check 'build/pbf.darwin-x64' 'Mach-O 64-bit';
# UPX disabled due to https://github.com/upx/upx/issues/187
# compress 'build/pbf.darwin-x64';

echo "[compile] windows x64";
env GOOS=windows GOARCH=386 go build -ldflags="-s -w" -gcflags=-trimpath=${GOPATH} -asmflags=-trimpath=${GOPATH} -o pbf.exe;
assert $?;
chmod +x pbf.exe;
mv pbf.exe build/pbf.win32-x64;
check 'build/pbf.win32-x64' 'PE32 executable';
compress 'build/pbf.win32-x64';
