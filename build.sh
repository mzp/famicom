#!/bin/bash

build_cmd() {
  echo -n "Building $1... "
  go build -o build/$1 mzp.jp/famicom/cmd/$1
}

mkdir -p build
for i in src/mzp.jp/famicom/cmd/*; do
  if build_cmd $(basename $i); then
    echo "success."
  else
    echo "fail."
    exit 1
  fi
done
