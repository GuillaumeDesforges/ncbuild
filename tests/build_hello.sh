#!/usr/bin/env bash

set -e

echo "Start integration test: build_hello"

SCRIPTPATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

cd "$SCRIPTPATH/.."
build_path="$SCRIPTPATH/../bin/"
go build -o "$build_path/ncbuild" main.go
export PATH="$PATH:$build_path"

cd "$SCRIPTPATH/../examples/hello"

build_output_path="$(ncbuild build)"
build_exit_code=$?
echo "$build_output_path"

if [ $build_exit_code -ne 0 ]; then
  echo "Build failed"
  exit 1
fi

filepath="$build_output_path/hello.txt"
if [ ! -f "$filepath" ]; then
  echo "Failed: $filepath not found"
  exit 1
fi

expected="hello world"
if [ ! "$(cat "$filepath")" == "$expected" ]; then
  echo "Content of $filepath does not match expected '$expected'"
  echo "Found: $(cat "$filepath")"
  exit 1
fi

echo "Test passed"

