#!/bin/bash



remote="SillyTavernBack:/SillyTavern"
test="SillyTavernBack:/Test"
folder="../Backup/"

if [ -d "$folder" ] && [ "$1" != "" ]; then
    cd "$folder"
else
    bash localup.sh make
fi

if [ "$1" == "" ]; then
    echo No option selected
fi

if [ "$1" == "ls" ]; then
    rclone ls $remote
fi

if [ "$1" == "upload" ]; then
    rclone sync $folder $remote -L -P
fi

if [ "$1" == "download" ]; then
    rclone sync $remote $folder -L -P
fi