package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
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

var version string = "1.4"

func DownloadFileFromGitHub(apiURL string, fileName string) error {
	resp, err := http.Get(apiURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	re := regexp.MustCompile(fmt.Sprintf(`"browser_download_url":"(.+/%s)"`, fileName))
	matches := re.FindStringSubmatch(string(body))

	if len(matches) < 2 {
		return fmt.Errorf("No se pudo encontrar la URL de descarga del archivo '%s'", fileName)
	}

	downloadURL := matches[1]
	resp, err = http.Get(downloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func makeconf() {
	var data string
	fmt.Print("Enter the rclone remote server:")
	fmt.Scan(&data)
	cmd("echo " + data + " > remote.txt")
	pwd, _ := readCommand("pwd")
	fmt.Printf("Remote Saved in %vYour remote:%v\n", pwd, data)
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
	apiURL := "https://api.github.com/repos/Tom5521/SillyTavernBackup/releases/latest"
	if option == "Termux" {
		fileName = "backup-aarch64"
	}
	if option == "pc" {
		fileName = "backup-x86-64		"
	}
	err := DownloadFileFromGitHub(apiURL, fileName)
	if err != nil {
		panic(err)
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
			if err == 1 {
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
		cmd("echo 'cd SillyTavernBackup' > backup")
		cmd("echo './backup $1 $2' >> backup")
	case "version":
		fmt.Println("SillyTavernBackup version", version, "\nUnder the MIT licence\nCreated by Tom5521")
	case "remote":
		makeconf()
	default:
		fmt.Println("Option not specified...")
	}
}
