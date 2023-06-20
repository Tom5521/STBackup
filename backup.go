package main

import (
	"fmt"
	"os"
	"os/exec"
)

// declarar variables globales locales de carpetas del backup,no hara falta tocarlas
var back string = "Backup/"

// declarar variables globales remotas. Tampoco hay que tocarlas
var folder, remote string = "../Backup/", readconf("name-of-remote.txt")

// declarar carpetas y archivos a excluir
var exclude_folders string = "--exclude webfonts --exclude scripts --exclude index.html --exclude css --exclude img --exclude favicon.ico --exclude script.js --exclude style.css --exclude Backup --exclude colab --exclude docker --exclude Dockerfile --exclude LICENSE --exclude node_modules --exclude package.json --exclude package-lock.json --exclude replit.nix --exclude server.js --exclude SillyTavernBackup --exclude src --exclude Start.bat --exclude start.sh --exclude UpdateAndStart.bat --exclude Update-Instructions.txt --exclude tools --exclude .dockerignore --exclude .editorconfig --exclude .git --exclude .github --exclude .gitignore --exclude .npmignore --exclude .replit "

var include_folders string = "--include backgrounds --include 'group chats' --include 'KoboldAI Settings' --include settings.json --include characters --include groups --include notes --include sounds --include worlds --include chats --include i18n.json --include 'NovelAI Settings' --include img --include 'OpenAI Settings' --include 'TextGen Settings' --include themes --include 'User Avatars' --include secrets.json --include thumbnails --include config.conf --include poe_device.json --include public --include uploads "

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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Option not specified")
		return
	}
	switch os.Args[1] {
	case "make":
		os.Chdir("..")
		os.MkdirAll("Backup/public", os.ModePerm)
	case "save":
		os.Chdir("..")
		cmd("rsync -av --progress " + exclude_folders + "--delete . " + " " + back)
		os.Chdir("SillyTavernBackup")
	case "restore":
		os.Chdir("..")
		cmd("rsync -av --progress " + exclude_folders + include_folders + "--delete " + back + " " + ".")
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
		if len(os.Args) < 2 {
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
	case "init":
		os.Chdir("..")
		cmd("bash start.sh")
		os.Chdir("SillyTavernBackup")
	}
}
