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
public="Backup/public"

if [ "$1" == "save" ];then
    echo "Saving Chats"
    cp -rf public/chats/ $public
    echo "Saving Characters"
    cp -rf public/characters/ $public
    echo "Saving OpenAI Settings"
    cp -rf public/"OpenAI Settings"/ $public
    echo "Saving Themes"
    cp -rf public/themes/ $public
    echo "Saving Worlds"
    cp -rf public/worlds/ $public
    echo "Saving User Avatars"
    cp -rf public/'User Avatars'/ $public
    echo "Saving Backgrounds"
    cp -rf public/backgrounds/ $public
    echo "Saving Group Chats"
    cp -rf public/'group chats'/ $public
    echo "Saving Groups"
    cp -rf public/groups/ $public
    echo "Saving Thumbnails"
    cp -rf thumbnails $back
    echo 'Saving "secrets.json"'
    cp -rf secrets.json $back
    echo "Saving Configs"
    cp -rf config.conf $back
    cp -rf public/settings.json $public
    cp -rf public/i18n.json $public
    
fi

if [ "$1" == "restore" ]; then
    echo "Restoring Chats"
    cp -rf $chats public/
    echo "Restoring Characters"
    cp -rf $characters public/
    echo "Restoring OpenAI settings"
    cp -rf Backup/public/'OpenAI Settings'/ public/
    echo "Restoring Themes"
    cp -rf $themes public/
    echo "Restoring Worlds"
    cp -rf $worlds public/
    echo "Restoring User Avatars"
    cp -rf Backup/public/'User Avatars'/ public/
    echo "Restoring Backgrounds"
    cp -rf $backgrounds public/
    echo "Restoring Group Chats"
    cp -rf Backup/public/'group chats'/ public/
    echo "Restoring groups"
    cp -rf $groups public/
    echo "Restoring Thumbnails"
    cp -rf Backup/thumbnails .
    echo 'Restoring "secrets.json"'
    cp -rf Backup/secrets.json .
    echo "Restoring Configs"
    cp -rf Backup/config.conf .
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
