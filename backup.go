package main

import (
	"fmt"
	"os"
	"os/exec"
)

// declarar variables globales locales de carpetas del backup,no hara falta tocarlas
var bchats, bcharacters, bthemes, bworlds, bbackgrounds, back, bgroups, bpublic, B_txt_gen_st string = "Backup/public/chats/", "Backup/public/characters/", "Backup/public/themes/", "Backup/public/worlds/", "Backup/public/backgrounds/", "Backup/", "Backup/public/groups", "Backup/public", "Backup/public/TextGen Settings"

// declarar variables globales locales de las carpetas originales
var chats, characters, themes, worlds, backgrounds, groups, public, txt_gen_st string = "public/chats/", "public/characters/", "public/themes/", "public/worlds/", "public/backgrounds/", "public/groups", "public/", "public/TextGen Settings"

// variables locales de carpetas especiales,son bien ojts con las comillas y todo eso --NO TOCAR--

var openAIst, userAvatr, grpChats string = "public/OpenAI Settings/", "public/User Avatars/", "public/group chats/"
var B_openAIst, B_userAvatr, B_grpChats string = "Backup/public/OpenAI Settings/", "Backup/public/User Avatars/", "Backup/public/group chats/"

// declarar variables de los archivos dentro de public
var settings_json, i18n string = "public/settings.json", "public/i18n.json"

var B_settings_json, B_i18n string = "Backup/public/settings.json", "Backup/public/i18n.json"

// declarar variables de las carpetas y archivos fuera de public
var thumbnails, secrets_json, configs_conf string = "thumbnails", "secrets.json", "config.conf"

var B_thumbnails, B_secrets_json, Bconfigs_conf string = "Backup/thumbnails", "Backup/secrets.json", "Backup/config.conf"

// declarar variables globales remotas. Tampoco hay que tocarlas
var folder, remote string = "../Backup/", readconf("name-of-remote.txt")

func readconf(file string) string {
	data, _ := os.ReadFile(file)
	return string(data)
}

func cmd(input string) {
	cmd := exec.Command("sh", "-c", input)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()
}
func cp(org, dest string) {
	cmd := exec.Command("cp", "-rf", org, dest)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Run()
}
func save() {
	fmt.Println("Saving Chats...")
}

func main() {
	sav := func(data string) {
		fmt.Println("Saving " + data + "...")
	}
	if len(os.Args) < 2 {
		fmt.Println("Option not specified")
		return
	}
	switch os.Args[1] {
	case "test":
		os.Chdir("..")
		cmd("ls")
	case "make":
		os.Chdir("..")
		os.MkdirAll("Backup/public", os.ModePerm)
	case "save":
		os.Chdir("..")
		sav("Chats")
		cp(chats, bpublic)
		sav("Characters")
		cp(characters, bpublic)
		sav("OpenAI settings")
		cp(openAIst, bpublic)
		sav("Themes")
		cp(themes, bthemes)
		sav("Worlds")
		cp(worlds, bpublic)
		sav("User Avatars")
		cp(userAvatr, bpublic)
		sav("Backgrounds")
		cp(backgrounds, bpublic)
		sav("Group Chats")
		cp(grpChats, bpublic)
		sav("Groups")
		cp(groups, bpublic)
		sav("Thumbnails")
		cp(thumbnails, back)
		sav("Secrets.json")
		cp(secrets_json, back)
		sav("Confings")
		cp(configs_conf, back)
		cp(settings_json, bpublic)
		cp(i18n, bpublic)
		cp(txt_gen_st, bpublic)
		os.Chdir("SillyTavernBackup")
	case "restore":
		os.Chdir("..")
		sav("Chats")
		cp(bchats, public)
		sav("Characters")
		cp(bcharacters, public)
		sav("OpenAI settings")
		cp(B_openAIst, public)
		sav("Themes")
		cp(bthemes, themes)
		sav("Worlds")
		cp(bworlds, public)
		sav("User Avatars")
		cp(B_userAvatr, public)
		sav("Backgrounds")
		cp(bbackgrounds, public)
		sav("Group Chats")
		cp(B_grpChats, public)
		sav("Groups")
		cp(bgroups, public)
		sav("Thumbnails")
		cp(B_thumbnails, ".")
		sav("Secrets.json")
		cp(B_secrets_json, ".")
		sav("Confings")
		cp(Bconfigs_conf, ".")
		cp(B_settings_json, public)
		cp(B_i18n, public)
		cp(B_txt_gen_st, public)
		os.Chdir("SillyTavernBackup")
	case "route":
		if len(os.Args) < 4 {
			fmt.Println("Backup destination not specified")
			return
		}
		os.Chdir("..")
		cmd("mv Backup/ " + os.Args[3] + " -f")
		os.Chdir("SillyTavernBackup")
	case "start":
		os.Chdir("..")
		cmd("node server.js")
		os.Chdir("SillyTavernBackup")
	case "update":
		if len(os.Args) < 4 {
			fmt.Println("Nothing Selected")
			return
		}
		if os.Args[2] == "SillyTavern" {
			os.Chdir("..")
			cmd("git pull")
			os.Chdir("SillyTavernBackup")
		}
		if os.Args[2] == "me" {
			cmd("git pull -f")
		}
	case "ls":
		cmd("rclone ls " + remote)
	case "upload":
		com := exec.Command("rclone", "sync", folder, remote, "-L", "-P")
		com.Stderr = os.Stderr
		com.Stdin = os.Stdin
		com.Stdout = os.Stdout
		com.Run()
	case "download":
		com := exec.Command("rclone", "sync", remote, folder, "-L", "-P")
		com.Stderr = os.Stderr
		com.Stdin = os.Stdin
		com.Stdout = os.Stdout
		com.Run()
	}
}
