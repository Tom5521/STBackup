#!/bin/bash

go build backup.go

mkdir builds
cd builds

mkdir x86-64
cd x86-64
export GOOS=linux
export GOARCH=amd64
go build -o backup-x86-64 ../../backup.go

cd ..

mkdir aarch64
cd aarch64
export GOOS=android
export GOARCH=arm64
go build -o backup-aarch64 ../../backup.go


