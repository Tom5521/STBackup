#!/bin/bash


if [ "$1" == "arm" ]; then
    binary="backup-arm"
fi

if [ "$1" == "x64" ]; then
    binary="backup-x86-64"
fi


if [ "$1" != "" ] && [ "$1" != "clone" ]; then
    if [ ! -d "SillyTavernBackup" ]; then
        mkdir SillyTavernBackup
    fi
    cd SillyTavernBackup
    echo "Downloading latest binary..."
    curl -LJo backup https://github.com/Tom5521/SillyTavernBackup/releases/latest/download/$binary
    echo "Giving execute permissions to the binary..."
    chmod +x backup
    echo "Configure remote..."
    ./backup remote
    ./backup remote
    echo "Creating link in SillyTavern root directory..."
    ./backup link
    echo Done!
fi

if [ "$1" == "clone" ]; then
    echo Cloning...
    git clone https://github.com/Tom5521/SillyTavernBackup.git
    cd SillyTavernBackup
    echo Compiling...
    go build -o backup main.go
    echo "Configure remote..."
    ./backup remote
    echo "Creating link in SillyTavern root directory..."
    ./backup link
    echo Done!
fi
