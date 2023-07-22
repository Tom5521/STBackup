#!/bin/bash


if [ "$1" == "arm" ]; then
    binary="backup-arm"
fi

if [ "$1" == "x64" ]; then
    binary="backup-x86-64"
fi


if [ "$1" != "" ] && [ "$1" != "clone" ]; then
    if [ ! -d "STbackup" ]; then
        mkdir STbackup
    fi
    cd STbackup
    echo "Downloading latest binary..."
    curl -LJo stbackup https://github.com/Tom5521/STbackup/releases/latest/download/$binary
    echo "Giving execute permissions to the binary..."
    chmod +x stbackup
    echo "Configure remote..."
    ./stbackup remote
    echo "Creating link in SillyTavern root directory..."
    ./stbackup link
    echo Done!
fi

if [ "$1" == "clone" ]; then
    echo Cloning...
    git clone https://github.com/Tom5521/STbackup.git
    cd STbackup
    echo Compiling...
    go build -o stbackup main.go
    echo "Configure remote..."
    ./stbackup remote
    echo "Creating link in SillyTavern root directory..."
    ./stbackup link
    echo Done!
fi
