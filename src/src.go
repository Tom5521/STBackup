package src

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const Folder, Back string = "../Backup/", "Backup/"

var binpath, _ = filepath.Abs(os.Args[0])
var Root string = filepath.Dir(binpath)
var pre_Remote, _ = GetJsonValue("config.json", "remote")
var Remote string = pre_Remote.(string)

const Exclude_Folders string = "--exclude webfonts --exclude scripts --exclude index.html --exclude css --exclude img --exclude favicon.ico --exclude script.js --exclude style.css --exclude Backup --exclude colab --exclude docker --exclude Dockerfile --exclude LICENSE --exclude node_modules --exclude package.json --exclude package-lock.json --exclude replit.nix --exclude server.js --exclude SillyTavernBackup --exclude src --exclude Start.bat --exclude start.sh --exclude UpdateAndStart.bat --exclude Update-Instructions.txt --exclude tools --exclude .dockerignore --exclude .editorconfig --exclude .git --exclude .github --exclude .gitignore --exclude .npmignore --exclude backup --exclude .replit --exclude install.sh --exclude Backup.tar --exclude app.log --exclude i18n.json "

const Include_Folders string = "--include backgrounds/ --include 'group chats' --include 'KoboldAI Settings' --include settings.json --include characters --include groups --include notes --include sounds --include worlds --include chats --include 'NovelAI Settings' --include img --include 'OpenAI Settings' --include 'TextGen Settings' --include themes --include 'User Avatars' --include secrets.json --include thumbnails --include config.conf --include poe_device.json --include public --include uploads "

var logger = setupLogger("app.log")

func Logerror(text string) {
	logger.Fatalln("ERROR: " + text)
}
func Logwarn(text string) {
	logger.Println("WARNING: " + text)
}
func Loginfo(text string) {
	logger.Println("PROGRAM: " + text)
}
func Logfunc(text string) {
	logger.Println("FUNC:    ---" + text + "---")
}
func setupLogger(logFilePath string) *log.Logger {
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(file, "", log.Ldate|log.Ltime)
	return logger
}

// Important functions
func Cmd(input string) int {
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
func ReadCommand(command string) (string, int) {
	com := exec.Command("sh", "-c", command)
	data, err := com.Output()
	if err != nil {
		return "", 1
	}
	return string(data), 0
}
func Readconf() (string, error) {
	os.Chdir(Root)
	file, err := os.Open("config.json")
	if err != nil {
		return "", err
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	var config map[string]interface{}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return "", err
	}
	Remote := config["Remote"].(string)
	if Remote == "" {
		fmt.Println("Remote is empty.")
		Logwarn("Remote is empty.")
		return "", nil
	}
	return Remote, nil
}
func GetJsonValue(jsonFile string, variableName string) (interface{}, error) {
	os.Chdir(Root)
	ls, _ := ReadCommand("ls")
	if !strings.Contains(ls, jsonFile) {
		fmt.Println(jsonFile + " Not found!")
		Logwarn(jsonFile + " Not found!")
		return os.DevNull, nil
	}
	file, err := os.Open(jsonFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var jsonData map[string]interface{}
	err = json.Unmarshal(bytes, &jsonData)
	if err != nil {
		return nil, err
	}

	variableValue, ok := jsonData[variableName]
	if !ok {
		Logerror("Variable does not exist in the JSON file")
		return nil, errors.New("Variable does not exist in the JSON file")
	}
	return variableValue, nil
}
func Makeconf() error {
	os.Chdir(Root)
	Cmd("echo '{\"Remote\":\"\"}' > config.json")
	fmt.Print("Enter the rclone Remote server:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	pwd, _ := ReadCommand("pwd")
	fmt.Printf("Remote Saved in %vYour Remote:%v\n", pwd, input)
	Loginfo("Remote Saved\nRemote:'" + input + "'\nRoute:'" + pwd + "'")

	file, err := os.OpenFile("config.json", os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	var config map[string]interface{}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return err
	}
	config["Remote"] = input
	bytes, err = json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("config.json", bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
func Rebuild() {
	os.Chdir(Root)
	Logfunc("Rebuild")
	_, errcode := ReadCommand("go version")
	ls, _ := ReadCommand("ls")
	if !strings.Contains(ls, "main.go") {
		fmt.Println("Source code not found")
		Logerror("Source code not found")
	}
	if errcode == 1 {
		fmt.Println("No go compiler found")
		Logerror("No go compiler found")
		return
	}
	fmt.Println("Rebuilding...")
	err := Cmd("go build main.go")
	if err != 1 {
		fmt.Println("Rebuild Complete.")
		Logfunc("Rebuilded")
		return
	}
	Logerror("Error in src.Rebuild prosess")
}
func UpdateBin(option string) {
	os.Chdir(Root)
	var fileName string
	const repo string = "Tom5521/SillyTavernBackup"
	if option == "Termux" {
		fileName = "backup-aarch64"
	}
	if option == "pc" {
		fileName = "backup-x86-64"
	}
	err := DownloadLatestReleaseBinary(repo, fileName)
	if err != nil {
		fmt.Println(err)
	}
	Cmd("mv " + fileName + " backup")
	os.Chmod("backup", 0700)
}

func Rclone(parameter string) {
	_, err := ReadCommand("rclone version")
	if err == 1 {
		fmt.Println("Rclone not found.")
		Logerror("Rclone not found")
		return
	}
	lsstat, _ := ReadCommand("ls")
	if !strings.Contains(lsstat, "config.json") {
		Makeconf()
	}
	var com = exec.Command("")
	switch parameter {
	case "uptar":
		Logfunc("upload tar")
		com = exec.Command("rclone", "copy", "Backup.tar", Remote)
		defer Loginfo("tar uploaded")
	case "downtar":
		Logfunc("download tar")
		com = exec.Command("rclone", "copy", Remote+"/Backup.tar", "..")
		defer Loginfo("tar downloaded")
	case "up":
		Logfunc("upload")
		com = exec.Command("rclone", "sync", Folder, Remote, "-L", "-P")
		defer Loginfo("Files uploaded")
	case "down":
		Logfunc("download")
		com = exec.Command("rclone", "sync", Remote, Folder, "-L", "-P")
		defer Loginfo("Files downloaded")
	}
	com.Stderr = os.Stderr
	com.Stdin = os.Stdin
	com.Stdout = os.Stdout
	com.Run()
}
func DownloadLatestReleaseBinary(repo string, binName string) error {
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
