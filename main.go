package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Tom5521/SillyTavernBackup/src"
)

// Vars and Constants
const folder, back string = src.Folder, src.Back

var remote string = src.Remote

const exclude_folders string = src.Exclude_Folders
const include_folders string = src.Include_Folders

var root string = src.Root

const version string = "2"

// MAIN
func main() {
	os.Chdir(root)
	src.Loginfo("--------Start--------")
	defer src.Loginfo("---------End---------")
	_, rsyncstat := src.ReadCommand("rsync --version")
	if rsyncstat == 1 {
		fmt.Println("Rsync not found.")
		src.Logerror("Rsync not found.")
		return
	}
	if len(os.Args) < 2 {
		src.Logerror("Option not specified.")
		fmt.Println("Option not specified...")
		return
	}
	if os.Args[1] == "src.Rebuild" {
		src.Rebuild()
		return
	}
	switch os.Args[1] {
	case "make":
		src.Logfunc("Make")
		os.Chdir("..")
		os.MkdirAll("Backup/public", os.ModePerm)
	case "save":
		src.Logfunc("save")
		os.Chdir("..")
		src.Cmd("rsync -av --progress " + exclude_folders + "--delete . " + " " + back)
		os.Chdir(root)
		src.Loginfo("Files Saved")
		if len(os.Args) == 3 {
			if os.Args[2] == "tar" {
				src.Logfunc("save tarball")
				os.Chdir("..")
				tar := src.Cmd("tar -cvf Backup.tar Backup/")
				if tar != 0 {
					src.Loginfo("Tarbal created.")
				}
			}
		}
	case "restore":
		src.Logfunc("restore")
		os.Chdir("..")
		if len(os.Args) == 3 {
			if os.Args[2] == "tar" {
				ls, _ := src.ReadCommand("ls")
				src.Logfunc("restore from tarball")
				if strings.Contains(ls, "Backup") {
					src.Logwarn("Removing Backup/ folder")
					src.Cmd("rm -rf Backup/")
				}
				src.Cmd("tar -xvf Backup.tar")
			}
		}
		src.Cmd("rsync -av --progress " + exclude_folders + include_folders + "--delete " + back + " . ")
		os.Chdir(root)
		src.Loginfo("Files restored")
	case "route":
		if len(os.Args) < 3 {
			fmt.Println("Backup destination not specified")
			src.Logerror("Not enough arguments")
			return
		}
		os.Chdir("..")
		src.Cmd("mv Backup/ " + os.Args[2] + " -f")
		os.Chdir(root)
		src.Logfunc("route")
		if os.Args[3] == "tar" {
			src.Logfunc("route tar")
			os.Chdir("..")
			src.Cmd("mv Backup.tar " + os.Args[2] + " -f")
			src.Loginfo("Tar file moved to" + os.Args[2])
		}
	case "start":
		src.Logfunc("start")
		src.Cmd("node ../server.js")
		src.Loginfo("SillyTavern ended")
	case "update":
		if len(os.Args) < 2 {
			fmt.Println("Nothing Selected")
			src.Logerror("Nothing selected in update func")
			return
		}
		if os.Args[2] == "ST" {
			os.Chdir("..")
			src.Cmd("git pull")
			os.Chdir(root)
			src.Loginfo("SillyTavern Updated")
		}
		if os.Args[2] == "me" {
			_, ggit := src.ReadCommand("git status")
			err, _ := src.ReadCommand("ls")
			_, err2 := src.ReadCommand("go version")
			if !strings.Contains(err, "backup.go") || err2 == 1 || ggit == 1 {
				if err2 == 1 {
					fmt.Println("No go compiler found... Downloading binaries")
					src.Logerror("No go compiler found. Downloading binaries")
				}
				bindata, _ := src.ReadCommand("file backup")
				if strings.Contains(bindata, "x86-64") {
					src.Loginfo("Downloading x86-64 binary")
					src.UpdateBin("pc")
				}
				if strings.Contains(bindata, "ARM aarch64") {
					src.Loginfo("Downloading aarch64 binary")
					src.UpdateBin("Termux")
				}
			} else {
				src.Cmd("git pull")
				src.Loginfo("Updated with git")
				src.Rebuild()
			}
			src.Cmd("./backup link")
		}
	case "ls":
		src.Logfunc("ls")
		src.Cmd("src.Rclone ls " + remote)
	case "upload":
		src.Rclone("up")
		if len(os.Args) == 3 {
			if os.Args[2] == "tar" {
				src.Rclone("uptar")
			}
		}
	case "download":
		src.Rclone("down")
		if len(os.Args) == 3 {
			if os.Args[2] == "tar" {
				src.Rclone("downtar")
			}
		}
	case "init":
		src.Logfunc("init")
		src.Cmd("bash ../start.sh")
	case "link":
		src.Logfunc("link")
		os.Chdir("..")
		file, _ := os.Create("backup")
		defer file.Close()
		cont := "#!/bin/bash\n"
		cont += "cd SillyTavernBackup/\n"
		cont += "./backup $1 $2 $3 $4\n"
		file.WriteString(cont)
		os.Chmod("backup", 0700)
		src.Loginfo("linked")
	case "version":
		fmt.Println("SillyTavernBackup version", version, "\nUnder the MIT licence\nCreated by Tom5521")
	case "remote":
		src.Logfunc("remote")
		src.Makeconf()
	case "cleanlog":
		src.Cmd("echo '' > app.log")
		os.Exit(0)
	case "log":
		src.Cmd("cat app.log")
	case "help":
		fmt.Println("Please read the documentation in https://github.com/Tom5521/SillyTavernBackup\nAll it's in the README")
	case "test":
		fmt.Println(remote, root)
	default:
		src.Logerror("Option not specified.")
		fmt.Println("Option not specified...")
	}
}
