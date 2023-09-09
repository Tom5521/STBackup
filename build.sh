#!/bin/bash

file="main.go" # Declare the file name

if [ "$1" == "dev" ]; then
    echo Build for dev use...
    export GOOS=linux
    export GOARCH=amd64
    go build -ldflags="-s -w" -gcflags=-trimpath -tags linux .
    exit
fi

if [ "$1" == "p" ]; then
  arch=$(uname -m)
  echo "Build for personal use...ARCH=$arch"
  go build -ldflags="-s -w" -o stbackup -gcflags=-trimpath $file
  exit
fi

if [ "$1" != "d" ]; then
  echo No option selected...
  exit
fi
echo Build for distribution...

export GOOS=linux
export GOARCH=amd64
echo Building x86-64
go build -ldflags="-s -w" -gcflags=-trimpath -tags linux -o builds/x86-64/backup-x86-64 $file


export GOOS=android
export GOARCH=arm64
echo building android-ARM
go build -ldflags="-s -w" -gcflags=-trimpath -o builds/android-ARM/backup-android-arm -tags android $file

export GOOS=linux
export GOARCH=arm64
echo building linux-ARM
go build -ldflags="-s -w" -gcflags=-trimpath -o builds/linux-ARM/backup-linux-arm -tags linux

tree builds


