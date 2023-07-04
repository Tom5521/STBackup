package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Tom5521/SillyTavernBackup/src/getdata"
	"github.com/Tom5521/SillyTavernBackup/src/log"
	"github.com/Tom5521/SillyTavernBackup/src/tools"
	"github.com/Tom5521/SillyTavernBackup/src/update"
)

// Vars and Constants
const folder, back string = getdata.Folder, getdata.Back

var root string = getdata.Root

// MAIN
func main() {
	log.Loginfo("--------Start--------")
	rebuildcheck := update.RebuildCheck()
	defer log.Loginfo("---------End---------")
	os.Chdir(root)
	if tools.CheckBranch() == false {
		log.Logwarn("You are in the dev branch!")
		fmt.Println("Note: You are using the dev branch. Which is usually always broken and is more for backup and anticipating changes than for users to experiment with.Please go back to the main branch, which is functional.")
	}
	if len(os.Args) < 2 && !rebuildcheck {
		log.Logerror("Option not specified.")
		fmt.Println("Option not specified...")
		return
	}
	switch os.Args[1] {
	case "make":
		log.Logfunc("Make")
		os.Chdir("..")
		os.MkdirAll("Backup/public", os.ModePerm)
	case "save":
		log.Logfunc("save")
		os.Chdir("..")
		tools.Cmd("rsync -av --progress " + getdata.Exclude_Folders + "--delete . " + " " + back)
		os.Chdir(root)
		log.Loginfo("Files Saved")
		if len(os.Args) == 3 {
			if os.Args[2] == "tar" {
				log.Logfunc("save tarball")
				os.Chdir("..")
				tar := tools.Cmd("tar -cvf Backup.tar Backup/")
				if tar != 0 {
					log.Loginfo("Tarbal created.")
				}
			}
		}
	case "restore":
		log.Logfunc("restore")
		os.Chdir("..")
		if len(os.Args) == 3 {
			if os.Args[2] == "tar" {
				ls, _ := tools.ReadCommand("ls")
				log.Logfunc("restore from tarball")
				if strings.Contains(ls, "Backup") {
					log.Logwarn("Removing Backup/ folder")
					tools.Cmd("rm -rf Backup/")
				}
				tools.Cmd("tar -xvf Backup.tar")
			}
		}
		tools.Cmd("rsync -av --progress " + getdata.Exclude_Folders + getdata.Include_Folders + "--delete " + back + " . ")
		os.Chdir(root)
		log.Loginfo("Files restored")
	case "route":
		if len(os.Args) < 3 {
			fmt.Println("Backup destination not specified")
			log.Logerror("Not enough arguments")
			return
		}
		os.Chdir("..")
		tools.Cmd("mv Backup/ " + os.Args[2] + " -f")
		os.Chdir(root)
		log.Logfunc("route")
		if os.Args[3] == "tar" {
			log.Logfunc("route tar")
			os.Chdir("..")
			tools.Cmd("mv Backup.tar " + os.Args[2] + " -f")
			log.Loginfo("Tar file moved to" + os.Args[2])
		}
	case "start":
		log.Logfunc("start")
		tools.Cmd("node ../server.js")
		log.Loginfo("SillyTavern ended")
	case "update":
		if len(os.Args) < 2 {
			fmt.Println("Nothing Selected")
			log.Logerror("Nothing selected in update func")
			return
		}
		if os.Args[2] == "ST" {
			os.Chdir("..")
			tools.Cmd("git pull")
			os.Chdir(root)
			log.Loginfo("SillyTavern Updated")
		}
		if os.Args[2] == "me" {
			_, ggit := tools.ReadCommand("git status")
			err, _ := tools.ReadCommand("ls")
			_, err2 := tools.ReadCommand("go version")
			if !strings.Contains(err, "main.go") || err2 == 1 || ggit == 1 {
				if err2 == 1 {
					fmt.Println("No go compiler found... Downloading binaries")
					log.Logerror("No go compiler found. Downloading binaries")
				}
				bindata, _ := tools.ReadCommand("file backup")
				if strings.Contains(bindata, "x86-64") {
					log.Loginfo("Downloading x86-64 binary")
					update.UpdateBin("pc")
				}
				if strings.Contains(bindata, "ARM aarch64") {
					log.Loginfo("Downloading aarch64 binary")
					update.UpdateBin("Termux")
				}
			} else {
				tools.Cmd("git pull")
				log.Loginfo("Updated with git")
				update.Rebuild()
			}
			tools.Cmd("./backup link")
		}
	case "ls":
		log.Logfunc("ls")
		tools.Cmd("rclone ls " + getdata.Remote)
	case "upload":
		tools.Rclone("up")
		if len(os.Args) == 3 {
			if os.Args[2] == "tar" {
				tools.Rclone("uptar")
			}
		}
	case "download":
		tools.Rclone("down")
		if len(os.Args) == 3 {
			if os.Args[2] == "tar" {
				tools.Rclone("downtar")
			}
		}
	case "init":
		log.Logfunc("init")
		tools.Cmd("bash ../start.sh")
	case "link":
		log.Logfunc("link")
		os.Chdir("..")
		file, _ := os.Create("backup")
		defer file.Close()
		cont := "#!/bin/bash\n"
		cont += "cd SillyTavernBackup/\n"
		cont += "./backup $1 $2 $3 $4\n"
		file.WriteString(cont)
		os.Chmod("backup", 0700)
		log.Loginfo("linked")
	case "version":
		fmt.Println("SillyTavernBackup version", getdata.Version, "\nUnder the MIT licence\nCreated by Tom5521")
	case "remote":
		log.Logfunc("remote")
		tools.Makeconf()
	case "cleanlog":
		tools.Cmd("echo '' > app.log")
		os.Exit(0)
	case "log":
		tools.Cmd("cat app.log")
	case "printconfig":
		tools.Cmd("cat config.json")
	case "help":
		fmt.Println("Please read the documentation in https://github.com/Tom5521/SillyTavernBackup\nAll it's in the README")
	default:
		log.Logerror("Option not specified.")
		fmt.Println("Option not specified...")
	}
}
