package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func main() {
	loginfo("--------Start--------")
	defer loginfo("---------End---------")
	_, rsyncstat := readCommand("rsync --version")
	if rsyncstat == 1 {
		fmt.Println("Rsync not found.")
		logerror("Rsync not found.")
		return
	}
	if len(os.Args) < 2 {
		logerror("Option not specified.")
		fmt.Println("Option not specified...")
		return
	}
	switch os.Args[1] {
	case "make":
		logfunc("Make")
		os.Chdir("..")
		os.MkdirAll("Backup/public", os.ModePerm)
	case "save":
		logfunc("save")
		os.Chdir("..")
		cmd("rsync -av --progress " + exclude_folders + "--delete . " + " " + back)
		os.Chdir(root)
		loginfo("Files Saved")
		if len(os.Args) == 3 {
			if os.Args[2] == "tar" {
				logfunc("save tarball")
				os.Chdir("..")
				tar := cmd("tar -cvf Backup.tar Backup/")
				if tar != 0 {
					loginfo("Tarbal created.")
				}
			}
		}
	case "restore":
		logfunc("restore")
		os.Chdir("..")
		if len(os.Args) == 3 {
			if os.Args[2] == "tar" {
				ls, _ := readCommand("ls")
				logfunc("restore tarball")
				if strings.Contains(ls, "Backup") {
					logwarn("Removing Backup/ folder")
					cmd("rm -rf Backup/")
				}
				cmd("tar -xvf Backup.tar")
			}
		}
		cmd("rsync -av --progress " + exclude_folders + include_folders + "--delete " + back + " " + ".")
		os.Chdir(root)
		loginfo("Files restored")
	case "route":
		if len(os.Args) < 3 {
			fmt.Println("Backup destination not specified")
			logerror("Not enough arguments")
			return
		}
		os.Chdir("..")
		cmd("mv Backup/ " + os.Args[2] + " -f")
		os.Chdir(root)
		logfunc("route")
		if os.Args[3] == "tar" {
			logfunc("route tar")
			os.Chdir("..")
			cmd("mv Backup.tar " + os.Args[2] + " -f")
			loginfo("Tar file moved to" + os.Args[2])
		}
	case "start":
		logfunc("start")
		cmd("node ../server.js")
		loginfo("SillyTavern ended")
	case "update":
		if len(os.Args) < 2 {
			fmt.Println("Nothing Selected")
			logerror("Nothing selected in update func")
			return
		}
		if os.Args[2] == "ST" {
			os.Chdir("..")
			cmd("git pull")
			os.Chdir(root)
			loginfo("SillyTavern Updated")
		}
		if os.Args[2] == "me" {
			_, ggit := readCommand("git status")
			err, _ := readCommand("ls")
			_, err2 := readCommand("go version")
			if !strings.Contains(err, "backup.go") || err2 == 1 || ggit == 1 {
				if err2 == 1 {
					fmt.Println("No go compiler found... Downloading binaries")
					logerror("No go compiler found. Downloading binaries")
				}
				bindata, _ := readCommand("file backup")
				if strings.Contains(bindata, "x86-64") {
					loginfo("Downloading x86-64 binary")
					updateBin("pc")
				}
				if strings.Contains(bindata, "ARM aarch64") {
					loginfo("Downloading aarch64 binary")
					updateBin("Termux")
				}
			} else {
				cmd("git pull")
				loginfo("Updated with git")
				rebuild()
			}
			cmd("./backup link")
		}
	case "ls":
		logfunc("ls")
		cmd("rclone ls " + remote)
	case "upload":
		rclone("up")
		if len(os.Args) == 3 {
			if os.Args[2] == "tar" {
				rclone("uptar")
			}
		}
	case "download":
		rclone("down")
		if len(os.Args) == 3 {
			if os.Args[2] == "tar" {
				rclone("downtar")
			}
		}
	case "init":
		logfunc("init")
		cmd("bash ../start.sh")
	case "rebuild":
		rebuild()
	case "link":
		logfunc("link")
		os.Chdir("..")
		file, _ := os.Create("backup")
		defer file.Close()
		cont := "#!/bin/bash\n"
		cont += "cd SillyTavernBackup\n"
		cont += "./backup $1 $2 $3 $4\n"
		file.WriteString(cont)
		os.Chmod("backup", 0700)
		loginfo("linked")
	case "version":
		fmt.Println("SillyTavernBackup version", version, "\nUnder the MIT licence\nCreated by Tom5521")
	case "remote":
		logfunc("remote")
		makeconf()
	case "cleanlog":
		cmd("echo '' > app.log")
		os.Exit(0)
	case "log":
		cmd("cat app.log")
	case "help":
		fmt.Println("Please read the documentation in https://github.com/Tom5521/SillyTavernBackup\nAll it's in the README")
	default:
		logerror("Option not specified.")
		fmt.Println("Option not specified...")
	}
}

const back, root string = "Backup/", "SillyTavernBackup"

const folder string = "../Backup/"

var remote string = readconf("remote.txt")

const exclude_folders string = "--exclude webfonts --exclude scripts --exclude index.html --exclude css --exclude img --exclude favicon.ico --exclude script.js --exclude style.css --exclude Backup --exclude colab --exclude docker --exclude Dockerfile --exclude LICENSE --exclude node_modules --exclude package.json --exclude package-lock.json --exclude replit.nix --exclude server.js --exclude SillyTavernBackup --exclude src --exclude Start.bat --exclude start.sh --exclude UpdateAndStart.bat --exclude Update-Instructions.txt --exclude tools --exclude .dockerignore --exclude .editorconfig --exclude .git --exclude .github --exclude .gitignore --exclude .npmignore --exclude backup --exclude .replit --exclude install.sh --exclude Backup.tar "

const include_folders string = "--include backgrounds --include 'group chats' --include 'KoboldAI Settings' --include settings.json --include characters --include groups --include notes --include sounds --include worlds --include chats --include i18n.json --include 'NovelAI Settings' --include img --include 'OpenAI Settings' --include 'TextGen Settings' --include themes --include 'User Avatars' --include secrets.json --include thumbnails --include config.conf --include poe_device.json --include public --include uploads "

const version string = "1.7"

var logger = setupLogger("app.log")

func logerror(text string) {
	logger.Fatalln("ERROR: " + text)
}
func logwarn(text string) {
	logger.Println("WARNING: " + text)
}
func loginfo(text string) {
	logger.Println("PROGRAM: " + text)
}
func logfunc(text string) {
	logger.Println("FUNC:    ---" + text + "---")
}
func makeconf() {
	fmt.Print("Enter the rclone remote server:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	cmd("echo " + input + " > remote.txt")
	pwd, _ := readCommand("pwd")
	fmt.Printf("Remote Saved in %vYour remote:%v\n", pwd, input)
}
func setupLogger(logFilePath string) *log.Logger {
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(file, "", log.Ldate|log.Ltime)
	return logger
}
func downloadLatestReleaseBinary(repo string, binName string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var release struct {
		Assets []struct {
			Name        string `json:"name"`
			DownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return err
	}
	var binaryURL string
	for _, asset := range release.Assets {
		if asset.Name == binName {
			binaryURL = asset.DownloadURL
			break
		}
	}
	if binaryURL == "" {
		return fmt.Errorf("No se encontró el binario %s en la última versión de %s", binName, repo)
	}
	resp, err = http.Get(binaryURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(binName)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("The %s binary of the latest version of %s has been successfully downloaded.\n", binName, repo)
	return nil
}

func readconf(file string) string {
	ls, _ := readCommand("ls")
	if !strings.Contains(ls, file) {
		logwarn("remote.txt not found. The file was recreated")
		fmt.Println(file, "not found!\nUse './backup remote' to configure it.")
		return ""
	}
	data, _ := os.Open(file)
	defer data.Close()
	scanner := bufio.NewScanner(data)
	scanner.Scan()
	text := scanner.Text()
	return text
}
func readCommand(command string) (string, int) {
	com := exec.Command("sh", "-c", command)
	data, err := com.Output()
	if err != nil {
		return "", 1
	}
	return string(data), 0
}
func updateBin(option string) {
	var fileName string
	const repo string = "Tom5521/SillyTavernBackup"
	if option == "Termux" {
		fileName = "backup-aarch64"
	}
	if option == "pc" {
		fileName = "backup-x86-64"
	}
	err := downloadLatestReleaseBinary(repo, fileName)
	if err != nil {
		fmt.Println(err)
	}
	cmd("mv " + fileName + " backup")
	os.Chmod("backup", 0700)
}

func cmd(input string) int {
	cmd := exec.Command("sh", "-c", input)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return 1
	}
	return 0
}
func rebuild() {
	_, errcode := readCommand("go version")
	if errcode == 1 {
		fmt.Println("No go compiler found")
		logerror("No go compiler found")
		return
	}
	fmt.Println("Rebuilding...")
	err := cmd("go build backup.go")
	if err != 1 {
		fmt.Println("Rebuild Complete.")
		logfunc("Rebuild")
		return
	}
	logerror("Error in rebuild prosess")
}
func rclone(parameter string) {
	_, err := readCommand("rclone version")
	if err == 1 {
		fmt.Println("Rclone not found.")
		logerror("Rclone not found")
		return
	}
	lsstat, _ := readCommand("ls")
	if !strings.Contains(lsstat, "remote.txt") {
		makeconf()
	}
	var com = exec.Command("")
	if parameter == "uptar" {
		logfunc("upload tar")
		com = exec.Command("rclone", "copy", "Backup.tar", remote)
		defer loginfo("tar uploaded")
	}
	if parameter == "downtar" {
		logfunc("download tar")
		com = exec.Command("rclone", "copy", remote+"/Backup.tar", "..")
		defer loginfo("tar downloaded")
	}
	if parameter == "up" {
		logfunc("upload")
		com = exec.Command("rclone", "sync", folder, remote, "-L", "-P")
		defer loginfo("Files uploaded")
	}
	if parameter == "down" {
		logfunc("download")
		com = exec.Command("rclone", "sync", remote, folder, "-L", "-P")
		defer loginfo("Files downloaded")
	}
	com.Stderr = os.Stderr
	com.Stdin = os.Stdin
	com.Stdout = os.Stdout
	com.Run()
}
