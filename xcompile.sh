#!/bin/bash
set -euo pipefail

# note: works with golang 1.7 but not on 1.8!

# requires 'upx'
SKIP_COMPRESSION=true;

# ensure the compiler exits with status 0
function assert() {
  if [[ $1 != 0 ]]; then
    echo "one or more architectures failed to compile"
    exit $1;
  fi
}

# ensure the files were compiled to the correct architecture
declare -A matrix
matrix["build/pbf.darwin-x64"]="Mach-O 64-bit "
matrix["build/pbf.linux-arm"]="ELF 32-bit "
matrix["build/pbf.linux-x64"]="ELF 64-bit "
matrix["build/pbf.win32-x64"]="PE32 executable "

function checkFiles() {
  for path in "${!matrix[@]}"
  do
    expected="$path: ${matrix[$path]}";
    actual=$(file $path);
    actualArr=($actual);
    start=$(printf "%s " "${actualArr[@]:0:3}");
    if [ "${start}" != "$expected" ]; then
      echo "invalid file architecture: $path"
      echo "expected: $expected"
      echo "actual: $actual"
      echo "start: $start"
      exit 1
    fi
  done
}

echo "[compile] linux x64";
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH pbf.go;
assert $?;
chmod +x pbf;
mv pbf build/pbf.linux-x64;
$SKIP_COMPRESSION || upx --brute build/pbf.linux-x64;

echo "[compile] linux arm";
GOOS=linux GOARCH=arm go build -ldflags="-s -w" -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH pbf.go;
echo $?;
assert $?;
chmod +x pbf;
mv pbf build/pbf.linux-arm;
$SKIP_COMPRESSION || upx --brute build/pbf.linux-arm;

# echo "[compile] linux i386";
# GOOS=linux GOARCH=386 go build -ldflags="-s -w" -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH pbf.go;
# assert $?;
# chmod +x pbf;
# mv pbf build/pbf.linux-ia32;
# $SKIP_COMPRESSION || upx --brute build/pbf.linux-ia32;

echo "[compile] darwin x64";
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH pbf.go;
assert $?;
chmod +x pbf;
mv pbf build/pbf.darwin-x64;
$SKIP_COMPRESSION || upx --brute build/pbf.darwin-x64;

# echo "[compile] darwin i386";
# GOOS=darwin GOARCH=386 go build -ldflags="-s -w" -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH pbf.go;
# assert $?;
# chmod +x pbf;
# mv pbf build/pbf.darwin-ia32;
# $SKIP_COMPRESSION || upx --brute build/pbf.darwin-ia32;

echo "[compile] windows x64";
GOOS=windows GOARCH=386 go build -ldflags="-s -w" -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -o pbf.exe pbf.go;
assert $?;
chmod +x pbf.exe;
mv pbf.exe build/pbf.win32-x64;
$SKIP_COMPRESSION || upx --brute build/pbf.win32-x64;

# echo "[compile] windows i386";
# GOOS=windows GOARCH=386 go build -ldflags="-s -w" -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH pbf.go;
# assert $?;
# chmod +x pbf.exe;
# mv pbf.exe build/pbf.win32-ia32;
# $SKIP_COMPRESSION || upx --brute build/pbf.win32-ia32;

# echo "[compile] freebsd arm";
# GOOS=freebsd GOARCH=arm go build -ldflags="-s -w" -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH pbf.go;
# assert $?;
# chmod +x pbf;
# mv pbf build/pbf.freebsd-arm;
# $SKIP_COMPRESSION || upx --brute build/pbf.freebsd-arm;

# echo "[compile] freebsd x64";
# GOOS=freebsd GOARCH=amd64 go build -ldflags="-s -w" -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH pbf.go;
# assert $?;
# chmod +x pbf;
# mv pbf build/pbf.freebsd-x64;
# $SKIP_COMPRESSION || upx --brute build/pbf.freebsd-x64;

# echo "[compile] freebsd i386";
# GOOS=freebsd GOARCH=386 go build -ldflags="-s -w" -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH pbf.go;
# assert $?;
# chmod +x pbf;
# mv pbf build/pbf.freebsd-ia32;
# $SKIP_COMPRESSION || upx --brute build/pbf.freebsd-ia32;

# echo "[compile] openbsd i386";
# GOOS=openbsd GOARCH=386 go build -ldflags="-s -w" -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH pbf.go;
# assert $?;
# chmod +x pbf;
# mv pbf build/pbf.openbsd-ia32;
# $SKIP_COMPRESSION || upx --brute build/pbf.openbsd-ia32;

# echo "[compile] openbsd x64";
# GOOS=openbsd GOARCH=amd64 go build -ldflags="-s -w" -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH pbf.go;
# assert $?;
# chmod +x pbf;
# mv pbf build/pbf.openbsd-x64;
# $SKIP_COMPRESSION || upx --brute build/pbf.openbsd-x64;

checkFiles
