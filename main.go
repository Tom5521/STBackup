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

// MAIN
func main() {
	os.Chdir(getdata.Root)
	log.Info("--------Start--------")
	defer log.Info("---------End---------")
	update.RebuildCheck()
	if !tools.CheckBranch() {
		log.Warning("You are in the dev branch!")
		fmt.Println("Note: You are using the dev branch. Which is usually always broken and is more for backup and anticipating changes than for users to experiment with.Please go back to the main branch, which is functional.")
	}
	if len(os.Args) > 2 {
		log.Error("Option not specified.")
		return
	}
	if os.Args[1] == "rebuild" {
		update.Rebuild()
		os.Exit(0)
	}
	switch os.Args[1] {
	case "make":
		log.Func("Make")
		os.Chdir("..")
		os.MkdirAll("Backup/public", os.ModePerm)
	case "save":
		log.Func("save")
		os.Chdir("..")
		tools.Cmd(fmt.Sprintf("rsync -av --progress %v --delete . %v", getdata.Exclude_Folders, getdata.Back))
		os.Chdir(getdata.Root)
		log.Info("Files Saved")
		if len(os.Args) == 3 {
			if os.Args[2] == "tar" {
				log.Func("save tarball")
				os.Chdir("..")
				tar := tools.Cmd("tar -cvf Backup.tar Backup/")
				if tar != 0 {
					log.Info("Tarbal created.")
				}
			}
		}
	case "restore":
		log.Func("restore")
		os.Chdir("..")
		if len(os.Args) == 3 {
			if os.Args[2] == "tar" {
				ls, _ := tools.ReadCommand("ls")
				log.Func("restore from tarball")
				if strings.Contains(ls, "Backup") {
					log.Warning("Removing Backup/ folder")
					tools.Cmd("rm -rf Backup/")
				}
				tools.Cmd("tar -xvf Backup.tar")
			}
		}
		tools.Cmd(fmt.Sprintf("rsync -av --progress %s %s --delete %s", getdata.Exclude_Folders, getdata.Include_Folders, getdata.Back))
		os.Chdir(getdata.Root)
		log.Info("Files restored")
	case "route":
		if len(os.Args) < 3 {
			log.Error("Backup destination not specified")
			return
		}
		os.Chdir("..")
		tools.Cmd("mv Backup/ " + os.Args[2] + " -f")
		os.Chdir(getdata.Root)
		log.Func("route")
		if os.Args[3] == "tar" {
			log.Func("route tar")
			os.Chdir("..")
			tools.Cmd(fmt.Sprintf("mv Backup.tar %v -f", os.Args[2]))
			log.Info("Tar file moved to " + os.Args[2])
		}
	case "start":
		log.Func("start")
		tools.Cmd("node ../server.js")
	case "update":
		if len(os.Args) < 2 {
			log.Error("Nothing selected in update func")
			return
		}
		if os.Args[2] == "ST" {
			os.Chdir("..")
			tools.Cmd("git pull")
			os.Chdir(getdata.Root)
			log.Info("SillyTavern Updated")
		}
		if os.Args[2] == "me" {
			_, ggit := tools.ReadCommand("git status")
			err, _ := tools.ReadCommand("ls")
			_, err2 := tools.ReadCommand("go version")
			if !strings.Contains(err, "main.go") || err2 == 1 || ggit == 1 {
				if err2 == 1 {
					log.Error("No go compiler found. Downloading binaries")
				}
				bindata, _ := tools.ReadCommand("file backup")
				if strings.Contains(bindata, "x86-64") {
					log.Info("Downloading x86-64 binary")
					update.UpdateBin("pc")
				}
				if strings.Contains(bindata, "ARM aarch64") {
					log.Info("Downloading aarch64 binary")
					update.UpdateBin("Termux")
				}
			} else {
				tools.Cmd("git pull")
				log.Info("Updated with git")
				update.Rebuild()
			}
			tools.Cmd("./backup link")
		}
	case "ls":
		log.Func("ls")
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
		log.Func("init")
		tools.Cmd("bash ../start.sh")
		log.Func("start")
	case "link":
		log.Func("link")
		os.Chdir("..")
		tools.WriteFile("backup", "#!/bin/bash\ncd SillyTavernBackup/\n./backup $1 $2 $3 $4\n")
		os.Chmod("backup", 0700)
		log.Info("linked")
	case "version":
		fmt.Println("SillyTavernBackup version", getdata.Version, "\nUnder the MIT licence\nCreated by Tom5521")
		return
	case "remote":
		log.Func("remote")
		tools.Makeconf()
	case "cleanlog":
		tools.WriteFile("app.log", "")
		os.Exit(0)
	case "log":
		filecont, _ := tools.ReadFileCont("app.log")
		fmt.Println(filecont)
	case "printconf":
		filecont, _ := tools.ReadFileCont("config.json")
		fmt.Println(filecont)
	case "resetconf":
		var test string
		fmt.Println("Are you sure to reset the configuration (backups will not be deleted)? y/n")
		fmt.Scanln(&test)
		if test == "y" {
			tools.Cmd("rm config.json app.log")
		} else {
			fmt.Println("No option selected.")
		}
	case "help":
		fmt.Println("Please read the documentation in https://github.com/Tom5521/SillyTavernBackup\nAll it's in the README")
	case "test":
		if tools.CheckBranch() {
			return
		}
		update.Rebuild()
		fmt.Println(getdata.Exclude_Folders_extra)
		fmt.Println(getdata.Include_Folders_extra)
		fmt.Println(getdata.Remote)
		fmt.Println(getdata.Root)
	default:
		log.Error("No option selected.")
	}

}
