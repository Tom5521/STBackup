package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Tom5521/SillyTavernBackup/src/depends"
	"github.com/Tom5521/SillyTavernBackup/src/getdata"
	"github.com/Tom5521/SillyTavernBackup/src/log"
	"github.com/Tom5521/SillyTavernBackup/src/tools"
	"github.com/Tom5521/SillyTavernBackup/src/update"
)

// Init shell function
var sh = getdata.Sh{}

func main() {
	log.Info("--------Start--------")
	defer log.Info("---------End---------")
	os.Chdir(getdata.Root)        // Change to the root directory
	update.RebuildCheck()         // Check if the rebuild arg is on
	if !tools.CheckMainBranch() { // Check the git branch to display a warning
		log.Warning("You are in the dev branch!")
		fmt.Println(
			"Note: You are using the dev branch. Which is usually always broken and is more for backup and anticipating changes than for users to experiment with.Please go back to the main branch, which is functional.",
		)
	}
	if len(os.Args) < 2 { // Check if it has the required number of arguments.
		log.Error("Option not specified.", 1)
		return
	}
	switch os.Args[1] {
	case "make": // Make the folder structures
		log.Func("Make")
		os.Chdir("..")
		os.MkdirAll("Backup/public", os.ModePerm)
	case "save": // Copy the local files to backup
		log.Func("save")
		tools.CheckRsync()
		os.Chdir("..")
		sh.Cmd(
			fmt.Sprintf(
				"rsync -av --progress %v --delete . %v",
				getdata.Exclude_Folders,
				getdata.Back,
			),
		)
		os.Chdir(getdata.Root)
		log.Info("Files Saved")
		if len(os.Args) == 3 { // Check if tar arg is on
			if os.Args[2] == "tar" {
				log.Func("save tarball")
				os.Chdir("..")
				tar := sh.Cmd("tar -cvf Backup.tar Backup/")
				if tar == nil {
					log.Info("Tarbal created.")
				}
			}
		}
	case "restore": // Copy the files from backup folder to the local folder
		log.Func("restore")
		os.Chdir("..")
		tools.CheckRsync()
		if len(os.Args) == 3 {
			if os.Args[2] == "tar" {
				log.Func("restore from tarball")
				if tools.CheckDir("Backup") {
					log.Warning("Removing Backup/ folder")
					sh.Cmd("rm -rf Backup/")
				}
				sh.Cmd("tar -xvf Backup.tar")
			}
		}
		sh.Cmd(
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
		sh.Cmd("mv Backup/ " + os.Args[2] + " -f")
		os.Chdir(getdata.Root)
		log.Func("route")
		if os.Args[3] == "tar" {
			log.Func("route tar")
			os.Chdir("..")
			sh.Cmd(fmt.Sprintf("mv Backup.tar %v -f", os.Args[2]))
			log.Info("Tar file moved to " + os.Args[2])
		}
	case "start": // Start SillyTavern with node js and server.js
		tools.SillyTavern("start")
	case "update":
		if len(os.Args) < 2 {
			log.Error("Nothing selected in update func", 3)
			return
		}
		if os.Args[2] == "ST" { // Update SillyTavern
			os.Chdir("..")
			sh.Cmd("git pull")
			os.Chdir(getdata.Root)
			log.Info("SillyTavern Updated")
		}
		if os.Args[2] == "me" { // Update this program
			if _, err := sh.Out("go version"); !tools.CheckDir("main.go") || err != nil ||
				!tools.CheckGit() {
				if err != nil {
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
				sh.Cmd("git pull")
				log.Info("Updated with git")
				update.Rebuild()
			}
			sh.Cmd("./backup link")
			sh.Cmd("rm config.json")
		}
	case "ls": // List the files and dirs in the remote
		tools.Rclone("ls")
	case "upload": // Synchronizes the remote folder with the local folder
		tools.Rclone("up")
		if len(os.Args) == 3 {
			if os.Args[2] == "tar" {
				tools.Rclone("uptar")
			}
		}
	case "download": // Synchronizes local folder with remote folder
		tools.Rclone("down")
		if len(os.Args) > 3 {
			if os.Args[2] == "tar" {
				tools.Rclone("downtar")
			}
		}
	case "init": // Execute start.sh for first run
		tools.SillyTavern("init")
	case "link": // Create a direct access bash file in the upper folder
		log.Func("link")
		fmt.Println("writing in stbackup file...")
		tools.WriteFile("../stbackup", "#!/bin/bash\ncd SillyTavernBackup/\n./backup $1 $2 $3 $4\n")
		fmt.Println("Giving exec permissions to stbackup file...")
		os.Chmod("../stbackup", 0700)
		fmt.Println("link completed.")
		log.Info("linked")
	case "version": // Print the current version,the author, and the licence
		fmt.Println(
			"SillyTavernBackup version",
			getdata.Version,
			"\nUnder the MIT licence\nCreated by Tom5521",
		)
		return
	case "remote": // Interactive config for the remote dir
		log.Func("remote")
		tools.Makeconf()
	case "cleanlog": // Clean the log file
		tools.WriteFile("app.log", "")
		os.Exit(0)
	case "log": // Print the log file
		filecont, _ := tools.ReadFileCont("app.log")
		fmt.Println(filecont)
	case "printconf": // Print the config values
		fmt.Println("Remote:", getdata.Configs.Local_rclone)
		fmt.Println("Local Rclone:", getdata.Configs.Local_rclone)
		fmt.Println("Extra Include Folders:", getdata.Configs.Include_Folders)
		fmt.Println("Extra Exclude Folders:", getdata.Configs.Exclude_Folders)
		fmt.Println("")
	case "resetconf": // Delete config.json and create a new one
		var test string
		fmt.Println("Are you sure to reset the configuration (backups will not be deleted)? y/n")
		fmt.Scanln(&test)
		if test == "y" {
			sh.Cmd("rm config.json app.log")
		} else {
			log.Error("No option selected.", 1)
		}
	case "download-rclone": // Download rclone local binary
		log.Info("rclone download")
		fmt.Println("Downloading and unzipping rclone...")
		getdata.Local_rclone = true
		depends.DownloadRclone()
		log.Info("Rclone donwloaded")
		os.Exit(0)
	case "setloglevel":
	case "help": // Print a help message
		fmt.Print(getdata.Help)
	case "test": // Test the program. Only works in the dev branch | The comments below are not relevant because they are only for testing purposes.
		if tools.CheckMainBranch() {
			log.Error("This func only works in the dev branch", 26)
			return
		}
		fmt.Println(os.Args)
		fmt.Println("Testing...")
		fmt.Print("F.D.:")
		filedata, _ := sh.Out("file backup")
		fmt.Println(filedata)
		//fmt.Println("V.:", getdata.Version)
		//fmt.Println("Exclude folders extra:", getdata.Exclude_Folders_extra)
		//fmt.Println("Exclude folders def:", "||-"+getdata.Exclude_Folders+"-||")
		//fmt.Println("Include folders extra:", getdata.Include_Folders_extra)
		//fmt.Println("Include folders def:", "||-"+getdata.Include_Folders+"-||")
		fmt.Println("Remote:", getdata.Remote)
		fmt.Println("Root Directory:", getdata.Root)
		fmt.Println("N.V.:", getdata.Version)
		fmt.Println("Arch:", getdata.Architecture)
		fmt.Println(
			"Check def git:",
			tools.CheckGit(),
			"| Check Main Branch:",
			tools.CheckMainBranch(),
			"| Check git directory:",
			tools.CheckDir(".git"),
		)
		//fmt.Println("Dirs:", sh.Cmd("exa -a"))
		//update.DownloadLatestBinary("backup-x86-64")
	default:
		log.Error("No option selected.", 1)
	}

}
