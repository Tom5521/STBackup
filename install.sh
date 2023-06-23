#!/bin/bash


if [ "$1" == "termux" ]; then
    binary="backup-aarch64"
fi

if [ "$1" == "pc" ]; then
    binary="backup-x86-64"
fi


if [ "$1" != "" ] && [ "$1" != "clone" ] && [ "$1" != "make" ]; then
    if [ ! -d "SillyTavernBackup" ]; then
        mkdir SillyTavernBackup
    fi
    cd SillyTavernBackup
    echo "Downloading latest binary..."
    curl -LJO $(curl -s https://api.github.com/repos/Tom5521/SillyTavernBackup/releases/latest | grep "browser_download_url.*$binary" | cut -d : -f 2,3 | tr -d \")
    echo "renaming..."
    mv $binary backup
    echo "Giving execute permissions to the binary..."
    chmod +x backup
    echo "Configure remote..."
    ./backup
    echo "Creating link in SillyTavern root directory..."
    ./backup link
    echo Done!
fi

if [ "$1" == "clone" ]; then
    echo Cloning...
    git clone https://github.com/Tom5521/SillyTavernBackup.git
    cd SillyTavernBackup
    echo Compiling...
    go build backup.go
    echo "Configure remote..."
    ./backup
    echo "Creating link in SillyTavern root directory..."
    ./backup link
    echo Done!
fi
