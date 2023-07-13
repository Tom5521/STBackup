package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Tom5521/SillyTavernBackup/src/checks"
	"github.com/Tom5521/SillyTavernBackup/src/depends"
	"github.com/Tom5521/SillyTavernBackup/src/getdata"
	"github.com/Tom5521/SillyTavernBackup/src/log"
	"github.com/Tom5521/SillyTavernBackup/src/tools"
	"github.com/Tom5521/SillyTavernBackup/src/update"
)

// MAIN
func main() {
	log.Info("--------Start--------")
	defer log.Info("---------End---------")
	os.Chdir(getdata.Root)
	update.RebuildCheck()
	if !checks.CheckBranch() {
		log.Warning("You are in the dev branch!")
		fmt.Println(
			"Note: You are using the dev branch. Which is usually always broken and is more for backup and anticipating changes than for users to experiment with.Please go back to the main branch, which is functional.",
		)
	}
	if len(os.Args) < 2 {
		log.Error("Option not specified.", 1)
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
		tools.Cmd(
			fmt.Sprintf(
				"rsync -av --progress %v --delete . %v",
				getdata.Exclude_Folders,
				getdata.Back,
			),
		)
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
				log.Func("restore from tarball")
				if checks.CheckDir("Backup") {
					log.Warning("Removing Backup/ folder")
					tools.Cmd("rm -rf Backup/")
				}
				tools.Cmd("tar -xvf Backup.tar")
			}
		}
		tools.Cmd(
			fmt.Sprintf(
				"rsync -av --progress --delete %s%s%s .",
				getdata.Exclude_Folders,
				getdata.Include_Folders,
				getdata.Back,
			),
		)
		os.Chdir(getdata.Root)
		log.Info("Files restored")
	case "route":
		if len(os.Args) < 3 {
			log.Error("Backup destination not specified", 2)
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
			log.Error("Nothing selected in update func", 3)
			return
		}
		if os.Args[2] == "ST" {
			os.Chdir("..")
			tools.Cmd("git pull")
			os.Chdir(getdata.Root)
			log.Info("SillyTavern Updated")
		}
		if os.Args[2] == "me" {
			_, ggit := getdata.ReadCommand("git status")
			_, err2 := getdata.ReadCommand("go version")
			if !checks.CheckDir("main.go") || err2 == 1 || ggit == 1 {
				if err2 == 1 {
					log.Warning("No go compiler found. Downloading binaries")
				}
				if getdata.Architecture == "amd64" {
					log.Info("Downloading x86-64 binary")
					update.DownloadLatestBinary("backup-x86-64")
				}
				if strings.Contains(getdata.Architecture, "arm") {
					log.Info("Downloading arm binary")
					update.DownloadLatestBinary("backup-arm")
				}
			} else {
				tools.Cmd("git pull")
				log.Info("Updated with git")
				update.Rebuild()
			}
			tools.Cmd("./backup link")
			tools.Cmd("rm config.json")
		}
	case "ls":
		tools.Rclone("ls")
	case "upload":
		tools.Rclone("up")
		if len(os.Args) == 3 {
			if os.Args[2] == "tar" {
				tools.Rclone("uptar")
			}
		}
	case "download":
		tools.Rclone("down")
		if len(os.Args) > 3 {
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
		fmt.Println(
			"SillyTavernBackup version",
			getdata.Version,
			"\nUnder the MIT licence\nCreated by Tom5521",
		)
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
			log.Error("No option selected.", 1)
		}
	case "download-rclone":
		log.Info("rclone download")
		fmt.Println("Downloading and unzipping rclone...")
		getdata.Local_rclone = true
		depends.DownloadRclone()
		log.Info("Rclone donwloaded")
		os.Exit(0)
	case "help":
		fmt.Println(
			"Please read the documentation in https://github.com/Tom5521/SillyTavernBackup\nAll it's in the README",
		)
	case "test":
		if checks.CheckBranch() {
			return
		}
		fmt.Println(os.Args)
		fmt.Println("Testing...")
		fmt.Print("F.D.:")
		fmt.Println(getdata.ReadCommand("file backup"))
		fmt.Println("V.:", getdata.Version)
		fmt.Println("Exclude folders extra:", getdata.Exclude_Folders_extra)
		fmt.Println("Exclude folders def:", "||-"+getdata.Exclude_Folders+"-||")
		fmt.Println("Include folders extra:", getdata.Include_Folders_extra)
		fmt.Println("Include folders def:", "||-"+getdata.Include_Folders+"-||")
		fmt.Println("Remote:", getdata.Remote)
		fmt.Println("Root Directory:", getdata.Root)
		fmt.Println("N.V.:", getdata.Version)
		fmt.Println("Arch:", getdata.Architecture)
		fmt.Println("Dirs:", tools.Cmd("exa -a"))
		//update.DownloadLatestBinary("backup-x86-64")
	default:
		log.Error("No option selected.", 1)
	}

}
