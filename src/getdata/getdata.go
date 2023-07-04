package getdata

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Tom5521/SillyTavernBackup/src/log"
)

const Folder, Back string = "../Backup/", "Backup/"
const Version string = "2.1.1"

var binpath, _ = filepath.Abs(os.Args[0])
var Root string = filepath.Dir(binpath)
var pre_Remote, _ = GetJsonValue("config.json", "remote")
var Remote string = pre_Remote.(string)

var Include_Folders string = "--include backgrounds/ --include 'group chats' --include 'KoboldAI Settings' --include settings.json --include characters --include groups --include notes --include sounds --include worlds --include chats --include 'NovelAI Settings' --include img --include 'OpenAI Settings' --include 'TextGen Settings' --include themes --include 'User Avatars' --include secrets.json --include thumbnails --include config.conf --include poe_device.json --include public --include uploads --include backups"
var Exclude_Folders string = "--exclude webfonts --exclude scripts --exclude index.html --exclude css --exclude img --exclude favicon.ico --exclude script.js --exclude style.css --exclude Backup --exclude colab --exclude docker --exclude Dockerfile --exclude LICENSE --exclude node_modules --exclude package.json --exclude package-lock.json --exclude replit.nix --exclude server.js --exclude SillyTavernBackup --exclude src --exclude Start.bat --exclude start.sh --exclude UpdateAndStart.bat --exclude Update-Instructions.txt --exclude tools --exclude .dockerignore --exclude .editorconfig --exclude .git --exclude .github --exclude .gitignore --exclude .npmignore --exclude backup --exclude .replit --exclude install.sh --exclude Backup.tar --exclude app.log --exclude i18n.json "

func GetJsonValue(jsonFile string, variableName string) (interface{}, error) {
	os.Chdir(Root)
	ls, _ := ReadCommand("ls")
	if !strings.Contains(ls, jsonFile) {
		fmt.Println(jsonFile + " Not found!")
		log.Warning(jsonFile + " Not found!")
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
		log.Error("Variable does not exist in the JSON file")
		return nil, errors.New("Variable does not exist in the JSON file")
	}
	return variableValue, nil
}

func ReadCommand(command string) (string, int) {
	com := exec.Command("sh", "-c", command)
	data, err := com.Output()
	if err != nil {
		return "", 1
	}
	return string(data), 0
}
