package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// declarar variables globales locales de carpetas del backup,no hara falta tocarlas
var back string = "Backup/"

// declarar variables globales remotas. Tampoco hay que tocarlas
var folder, remote string = "../Backup/", readconf("remote.txt")

// declarar carpetas y archivos a excluir
var exclude_folders string = "--exclude webfonts --exclude scripts --exclude index.html --exclude css --exclude img --exclude favicon.ico --exclude script.js --exclude style.css --exclude Backup --exclude colab --exclude docker --exclude Dockerfile --exclude LICENSE --exclude node_modules --exclude package.json --exclude package-lock.json --exclude replit.nix --exclude server.js --exclude SillyTavernBackup --exclude src --exclude Start.bat --exclude start.sh --exclude UpdateAndStart.bat --exclude Update-Instructions.txt --exclude tools --exclude .dockerignore --exclude .editorconfig --exclude .git --exclude .github --exclude .gitignore --exclude .npmignore --exclude backup --exclude .replit "

// declarar archivos y carpetas a incluir
var include_folders string = "--include backgrounds --include 'group chats' --include 'KoboldAI Settings' --include settings.json --include characters --include groups --include notes --include sounds --include worlds --include chats --include i18n.json --include 'NovelAI Settings' --include img --include 'OpenAI Settings' --include 'TextGen Settings' --include themes --include 'User Avatars' --include secrets.json --include thumbnails --include config.conf --include poe_device.json --include public --include uploads "

var version string = "1.4.1"

func makeconf() {
	fmt.Print("Enter the rclone remote server:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	cmd("echo " + input + " > remote.txt")
	pwd, _ := readCommand("pwd")
	fmt.Printf("Remote Saved in %vYour remote:%v\n", pwd, input)
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
		fmt.Println(file, "not found!")
		makeconf()
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
	repo := "Tom5521/SillyTavernBackup"
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
		return
	}
	fmt.Println("Rebuilding...")
	err := cmd("go build backup.go")
	if err != 1 {
		fmt.Println("Rebuild Complete.")
		return
	}
	fmt.Println("Error")
}
func rclone(parameter string) {
	var com = exec.Command("echo", "ERROR-CALLING-RCLONE-FUNCTION")
	if parameter == "up" {
		com = exec.Command("rclone", "sync", folder, remote, "-L", "-P")
	}
	if parameter == "down" {
		com = exec.Command("rclone", "sync", remote, folder, "-L", "-P")
	}
	com.Stderr = os.Stderr
	com.Stdin = os.Stdin
	com.Stdout = os.Stdout
	com.Run()
}
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Option not specified...")
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
		if len(os.Args) < 3 {
			fmt.Println("Backup destination not specified")
			return
		}
		os.Chdir("..")
		cmd("mv Backup/ " + os.Args[2] + " -f")
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
		if os.Args[2] == "ST" {
			os.Chdir("..")
			cmd("git pull")
			os.Chdir("SillyTavernBackup")
		}
		if os.Args[2] == "me" {
			_, err := readCommand("git status")
			_, err2 := readCommand("go version")
			if err == 1 || err2 == 1 {
				if err2 == 1 {
					fmt.Println("No go compiler found... Downloading binaries")
				}
				bindata, _ := readCommand("file backup")
				if strings.Contains(bindata, "x86-64") {
					updateBin("pc")
				}
				if strings.Contains(bindata, "ARM aarch64") {
					updateBin("Termux")
				}
			} else {
				cmd("git pull")
				rebuild()
			}
		}
	case "ls":
		cmd("rclone ls " + remote)
	case "upload":
		rclone("up")
	case "download":
		rclone("down")
	case "init":
		os.Chdir("..")
		cmd("bash start.sh")
		os.Chdir("SillyTavernBackup")
	case "rebuild":
		rebuild()
	case "link":
		os.Chdir("..")
		cmd("touch backup")
		os.Chmod("backup", 0700)
		cmd("echo #!/bin/bash > backup")
		cmd("echo 'cd SillyTavernBackup' >> backup")
		cmd("echo './backup $1 $2' >> backup")
	case "version":
		fmt.Println("SillyTavernBackup version", version, "\nUnder the MIT licence\nCreated by Tom5521")
	case "remote":
		makeconf()
	default:
		fmt.Println("Option not specified...")
	}
}
