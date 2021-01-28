#!/bin/sh
set -e

if ! [ -x "$(command -v protoc)" ]; then
    echo "protoc is not installed"
    exit 1
fi

if ! [ -x "$(command -v go)" ]; then
    echo "go is not installed"
    exit 1
fi

if ! [ -x "$(command -v protoc-gen-go)" ]; then
    echo "protoc-gen-go is not installed"
    exit 1
fi

proto_path="$PWD/rpcs/proto/*"

for d in $proto_path; do
    if [ -d "${d}" ]; then 
        continue
    fi
    echo "========================= compiling $d ======================="
    protoc --proto_path=rpcs --go_out=plugins=grpc:$PWD/rpcs "rpcs/proto/$(basename "$d")"
    echo "========================= done compiling ======================="
done