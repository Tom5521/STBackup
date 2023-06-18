#!/bin/bash



remote="SillyTavernBack:/SillyTavern"
test="SillyTavernBack:/Test"
folder="SillyTavern"

if [ -d $folder ]; then
    cd $folder
else
    mkdir $folder
fi


if [ "$1" == "" ]; then
    echo No option selected
fi

if [ "$1" == "ls" ]; then
    rclone ls $remote
fi

if [ "$1" == "upload" ]; then
    rclone sync . $remote -L -P
fi

if [ "$1" == "download" ]; then
    rclone sync $remote . -L -P
fi