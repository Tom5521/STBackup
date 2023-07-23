#!/bin/bash


if [ "$1" == "arm" ]; then
    binary="backup-arm"
fi

if [ "$1" == "x64" ]; then
    binary="backup-x86-64"
fi


if [ "$1" != "" ] && [ "$1" != "clone" ]; then
    if [ ! -d "STBackup" ]; then
        mkdir STBackup
    fi
    cd STBackup
    echo "Downloading latest binary..."
    curl -LJo STBackup https://github.com/Tom5521/STBackup/releases/latest/download/$binary
    echo "Giving execute permissions to the binary..."
    chmod +x STBackup
    echo "Configure remote..."
    ./STBackup remote
    echo "Creating link in SillyTavern root directory..."
    ./STBackup link
    echo Done!
fi

if [ "$1" == "clone" ]; then
    echo Cloning...
    git clone https://github.com/Tom5521/STBackup.git
    cd STBackup
    echo Compiling...
    go build .
    echo "Configure remote..."
    ./STBackup remote
    echo "Creating link in SillyTavern root directory..."
    ./STBackup link
    echo Done!
fi
