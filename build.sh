#!/bin/bash

file="main.go"

if [ "$1" == "d" ]; then
    echo Build for dev use...
    go build -ldflags="-s -w" -gcflags=-trimpath -tags linux -o backup $file
    exit
fi
echo Build for distribution...
go build -ldflags="-s -w" -gcflags=-trimpath -tags linux -o backup $file

if [ ! -d "builds" ]; then
    mkdir builds
fi
cd builds

if [ ! -d "x86-64" ]; then
    mkdir x86-64
fi
cd x86-64
export GOOS=linux
export GOARCH=amd64
echo "Building x86-64"
go build -ldflags="-s -w" -gcflags=-trimpath -tags linux -o backup-x86-64 ../../$file

cd ..
if [ ! -d "aarch64" ]; then
    mkdir aarch64
fi
cd aarch64
export GOOS=android
export GOARCH=arm64
echo "building aarch64"
go build -ldflags="-s -w" -gcflags=-trimpath -tags android -buildmode=pie -o backup-aarch64 ../../$file

cd ../
tree .


