#!/bin/bash


if [ "$1" == "termux" ]; then
    binary="backup-aarch64"
fi

if [ "$1" == "pc" ]; then
    binary="backup-x86-64"
fi

if [ "$1" != "" ] && [ "$1" != "clone" ]; then
    mkdir SillyTavernBackup
    cd SillyTavernBackup
    curl -LJO $(curl -s https://api.github.com/repos/Tom5521/SillyTavernBackup/releases/latest | grep "browser_download_url.*$binary" | cut -d : -f 2,3 | tr -d \")
    mv $binary backup
    chmod +x backup
    ./backup
fi

if [ "$1" == "clone" ]; then
    git clone https://github.com/Tom5521/SillyTavernBackup.git
    cd SillyTavernBackup
    go build backup.go
fi
