package getdata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Tom5521/SillyTavernBackup/src/log"
)

const Folder, Back string = "../Backup/", "Backup/"
const Version string = "2.3-dev"

var binpath, _ = filepath.Abs(os.Args[0])
var Root string = filepath.Dir(binpath)
var Remote, _ = GetJsonValue("config.json", "remote")

var Include_Folders string = "--include backgrounds/ --include 'group chats' --include 'KoboldAI Settings' --include settings.json --include characters --include groups --include notes --include sounds --include worlds --include chats --include 'NovelAI Settings' --include img --include 'OpenAI Settings' --include 'TextGen Settings' --include themes --include 'User Avatars' --include secrets.json --include thumbnails --include config.conf --include poe_device.json --include public --include uploads --include backups " + Include_Folders_extra()
var Exclude_Folders string = "--exclude webfonts --exclude scripts --exclude index.html --exclude css --exclude img --exclude favicon.ico --exclude script.js --exclude style.css --exclude Backup --exclude colab --exclude docker --exclude Dockerfile --exclude LICENSE --exclude node_modules --exclude package.json --exclude package-lock.json --exclude replit.nix --exclude server.js --exclude SillyTavernBackup --exclude src --exclude Start.bat --exclude start.sh --exclude UpdateAndStart.bat --exclude Update-Instructions.txt --exclude tools --exclude .dockerignore --exclude .editorconfig --exclude .git --exclude .github --exclude .gitignore --exclude .npmignore --exclude backup --exclude .replit --exclude install.sh --exclude Backup.tar --exclude app.log --exclude i18n.json " + Exclude_Folders_extra()

var Architecture = runtime.GOARCH
var Local_rclone_route string = "src/bin/"

func Local_rclone() bool {
	local_rclone, _ := GetJsonValue("config.json", "local-rclone")
	if local_rclone == "yes" {
		return true
	}
	return false
}

func Include_Folders_extra() string {
	pre_include_Folders_extra, _ := GetJsonValue("config.json", "include-folders")
	return ProsessString(pre_include_Folders_extra, "--include ")
}

func Exclude_Folders_extra() string {
	pre_exclude_Folders_extra, _ := GetJsonValue("config.json", "exclude-folders")
	return ProsessString(pre_exclude_Folders_extra, "--exclude ")
}

func GetJsonValue(jsonFile, variableName string) (string, error) {
	os.Chdir(Root)
	ls, _ := readCommand("ls")
	if !strings.Contains(ls, jsonFile) {
		fmt.Println(jsonFile + " Not found!")
		log.Warning(jsonFile + " Not found!")
		return "", nil
	}
	file, err := os.Open(jsonFile)
	if err != nil {
		return "", err
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	var jsonData map[string]interface{}
	err = json.Unmarshal(bytes, &jsonData)
	if err != nil {
		return "", err
	}

	variableValue, ok := jsonData[variableName]
	if !ok {
		log.Error("Variable does not exist in the JSON file")
		return "", nil
	}
	return variableValue.(string), nil
}

func readCommand(command string) (string, int) {
	com := exec.Command("sh", "-c", command)
	data, err := com.Output()
	if err != nil {
		return "", 1
	}
	return string(data), 0
}

func ProsessString(data, cond1 string) string {
	edit := func(org, sep string) string {
		words := strings.Split(org, " ")
		edited := strings.Join(words, sep+cond1+sep)
		return edited
	}
	pre_include_Folders_extra, _ := GetJsonValue("config.json", "include-folders")
	pre_exclude_Folders_extra, _ := GetJsonValue("config.json", "exclude-folders")
	if pre_exclude_Folders_extra == "" && cond1 == "--exclude " {
		cond1 = ""
	}
	if pre_include_Folders_extra == "" && cond1 == "--include " {
		cond1 = ""
	}
	return cond1 + edit(data, " ")
}
