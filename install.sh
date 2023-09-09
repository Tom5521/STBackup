#!/bin/bash

os=$(uname -o)
arch=$(uname -m)
if [ "$os" == "Android" ]; then
    binary="backup-android-arm"
fi

if [ "$os" == "GNU/Linux" ]; then
  if [ "$arch" == "x86_64" ]; then
    binary="backup-x86-64"
  fi
  if [ "$arch" == "aarch64" ]; then
    binary="backup-linux-arm"
  fi
fi


if [ "$1" != "clone" ]; then
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
    bash build.sh p
    echo "Configure remote..."
    ./STBackup remote
    echo "Creating link in SillyTavern root directory..."
    ./STBackup link
    echo Done!
fi
