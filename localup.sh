#!/bin/bash


#!/bin/bash

cd ..

if [ "$1" == "make" ]; then
    mkdir Silly-Tavern-Backup Silly-Tavern-Backup/public
    
    mkdir Silly-Tavern-Backup/public/chats Silly-Tavern-Backup/public/themes Silly-Tavern-Backup/public/"User Avatars" Silly-Tavern-Backup/public/worlds Silly-Tavern-Backup/public/characters Silly-Tavern-Backup/public/"OpenAI Settings" Silly-Tavern-Backup/public/backgrounds Silly-Tavern-Backup/public/"group chats"/ Silly-Tavern-Backup/public/groups/
fi
chats="Silly-Tavern-Backup/public/chats/"
characters="Silly-Tavern-Backup/public/characters/"
themes="Silly-Tavern-Backup/public/themes/"
worlds="Silly-Tavern-Backup/public/worlds/"
backgrounds="Silly-Tavern-Backup/public/backgrounds/"
back="Silly-Tavern-Backup/"
groups="Silly-Tavern-Backup/public/groups"

if [ "$1" == "save" ];then
    echo "Saving Chats"
    cp -rf public/chats/* $chats
    echo "Saving Characters"
    cp -rf public/characters/* $characters
    echo "Saving OpenAI Settings"
    cp -rf public/"OpenAI Settings"/* Silly-Tavern-Backup/public/'OpenAI Settings'
    echo "Saving Themes"
    cp -rf public/themes/* $themes
    echo "Saving Worlds"
    cp -rf public/worlds/* $worlds
    echo "Saving User Avatars"
    cp -rf public/'User Avatars'/* Silly-Tavern-Backup/public/'User Avatars'
    echo "Saving Backgrounds"
    cp -rf public/backgrounds/* $backgrounds
    echo "Saving Group Chats"
    cp -rf public/'group chats'/* Silly-Tavern-Backup/public/'group chats'
    echo "Saving Groups"
    cp -rf public/groups/* $group
    echo "Saving Thumbnails"
    cp -rf thumbnails $back
    echo 'Saving "secrets.json"'
    cp -rf secrets.json $back
    echo "Saving Configs"
    cp -rf config.conf $back
    cp -rf public/settings.json Silly-Tavern-Backup/public/
    cp -rf public/i18n.json Silly-Tavern-Backup/public/
    
fi

if [ "$1" == "restore" ]; then
    echo "Restoring Chats"
    cp -rf $chats* public/chats/
    echo "Restoring Characters"
    cp -rf $characters* public/characters/
    echo "Restoring OpenAI settings"
    cp -rf Silly-Tavern-Backup/public/'OpenAI Settings'/* public/"OpenAI Settings"/
    echo "Restoring Themes"
    cp -rf $themes* public/themes/
    echo "Restoring Worlds"
    cp -rf $worlds* public/worlds/
    echo "Restoring User Avatars"
    cp -rf Silly-Tavern-Backup/public/'User Avatars'/* public/'User Avatars'/
    echo "Restoring Backgrounds"
    cp -rf $backgrounds* public/backgrounds/
    echo "Restoring Group Chats"
    cp -rf Silly-Tavern-Backup/public/'group chats'/* public/'group chats'/
    echo "Restoring groups"
    cp -rf $groups* public/groups/
    echo "Restoring Thumbnails"
    cp -rf thumbnails $back
    echo 'Restoring "secrets.json"'
    cp -rf secrets.json $back
    echo "Restoring Configs"
    cp -rf config.conf $back
    cp -rf Silly-Tavern-Backup/public/settings.json public/
    cp -rf Silly-Tavern-Backup/public/i18n.json public/
fi

if [ "$2" == "route" -a "$3" != "" ]; then
    echo "Silly-Tavern-Backup is in $3"
    mv Silly-Tavern-Backup/ "$3" -f
fi

if [ "$1" == "upload" ]; then
    git add  .
    git commit -a -m "Upload Files"
    echo 'Now make a "git push"'
fi
if [ "$1" == "download" ]; then
    git pull
fi
if [ "$1" == "" ]; then
    echo "Option not specified"
fi
