#!/bin/bash


if [ "$1" == "termux" ]; then
    binary="backup-aarch64"
fi

if [ "$1" == "pc" ]; then
    binary="backup-x86-64"
fi


if [ "$1" != "" ] && [ "$1" != "clone" ] && [ "$1" != "make" ]; then
    mkdir SillyTavernBackup
    cd SillyTavernBackup
    curl -LJO $(curl -s https://api.github.com/repos/Tom5521/SillyTavernBackup/releases/latest | grep "browser_download_url.*$binary" | cut -d : -f 2,3 | tr -d \")
    mv $binary backup
    chmod +x backup
    ./backup link
    ln -s name-of-remote.txt ../name-of-remote.txt
fi

if [ "$1" == "clone" ]; then
    git clone https://github.com/Tom5521/SillyTavernBackup.git
    cd SillyTavernBackup
    go build backup.go
    ./backup link
    ln -s name-of-remote.txt ../name-of-remote.txt
fi
