#!/bin/bash

cd ..

if [ "$1" == "make" ]; then
    mkdir Backup Backup/public
    mkdir Backup/public/chats Backup/public/themes Backup/public/"User Avatars" Backup/public/worlds Backup/public/characters Backup/public/"OpenAI Settings" Backup/public/backgrounds Backup/public/"group chats"/ Backup/public/groups/
fi
chats="Backup/public/chats/"
characters="Backup/public/characters/"
themes="Backup/public/themes/"
worlds="Backup/public/worlds/"
backgrounds="Backup/public/backgrounds/"
back="Backup/"
groups="Backup/public/groups"

if [ "$1" == "save" ];then
    echo "Saving Chats"
    cp -rf public/chats/* $chats
    echo "Saving Characters"
    cp -rf public/characters/* $characters
    echo "Saving OpenAI Settings"
    cp -rf public/"OpenAI Settings"/* Backup/public/'OpenAI Settings'
    echo "Saving Themes"
    cp -rf public/themes/* $themes
    echo "Saving Worlds"
    cp -rf public/worlds/* $worlds
    echo "Saving User Avatars"
    cp -rf public/'User Avatars'/* Backup/public/'User Avatars'
    echo "Saving Backgrounds"
    cp -rf public/backgrounds/* $backgrounds
    echo "Saving Group Chats"
    cp -rf public/'group chats'/* Backup/public/'group chats'
    echo "Saving Groups"
    cp -rf public/groups/* $group
    echo "Saving Thumbnails"
    cp -rf thumbnails $back
    echo 'Saving "secrets.json"'
    cp -rf secrets.json $back
    echo "Saving Configs"
    cp -rf config.conf $back
    cp -rf public/settings.json Backup/public/
    cp -rf public/i18n.json Backup/public/
    
fi

if [ "$1" == "restore" ]; then
    echo "Restoring Chats"
    cp -rf $chats* public/chats/
    echo "Restoring Characters"
    cp -rf $characters* public/characters/
    echo "Restoring OpenAI settings"
    cp -rf Backup/public/'OpenAI Settings'/* public/"OpenAI Settings"/
    echo "Restoring Themes"
    cp -rf $themes* public/themes/
    echo "Restoring Worlds"
    cp -rf $worlds* public/worlds/
    echo "Restoring User Avatars"
    cp -rf Backup/public/'User Avatars'/* public/'User Avatars'/
    echo "Restoring Backgrounds"
    cp -rf $backgrounds* public/backgrounds/
    echo "Restoring Group Chats"
    cp -rf Backup/public/'group chats'/* public/'group chats'/
    echo "Restoring groups"
    cp -rf $groups* public/groups/
    echo "Restoring Thumbnails"
    cp -rf thumbnails $back
    echo 'Restoring "secrets.json"'
    cp -rf secrets.json $back
    echo "Restoring Configs"
    cp -rf config.conf $back
    cp -rf Backup/public/settings.json public/
    cp -rf Backup/public/i18n.json public/
fi

if [ "$2" == "route" -a "$3" != "" ]; then
    echo "Backup is in $3"
    mv Backup/ "$3" -f
fi

if [ "$1" == "" ]; then
    echo "Option not specified"
fi
