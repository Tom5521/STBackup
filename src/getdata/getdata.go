package getdata

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Tom5521/SillyTavernBackup/src/log"
)

const Folder, Back string = "../Backup/", "Backup/"
const Version string = "2.3.3"

var Remote string = remote()

var Include_Folders string = "--include backgrounds/ --include 'group chats' --include 'KoboldAI Settings' --include settings.json --include characters --include groups --include notes --include sounds --include worlds --include chats --include 'NovelAI Settings' --include img --include 'OpenAI Settings' --include 'TextGen Settings' --include themes --include 'User Avatars' --include secrets.json --include thumbnails --include config.conf --include poe_device.json --include public --include uploads --include backups " + Include_Folders_extra

var Exclude_Folders string = "--exclude webfonts --exclude scripts --exclude index.html --exclude css --exclude img --exclude favicon.ico --exclude script.js --exclude style.css --exclude Backup --exclude colab --exclude docker --exclude Dockerfile --exclude LICENSE --exclude node_modules --exclude package.json --exclude package-lock.json --exclude replit.nix --exclude server.js --exclude SillyTavernBackup --exclude src --exclude Start.bat --exclude start.sh --exclude UpdateAndStart.bat --exclude Update-Instructions.txt --exclude tools --exclude .dockerignore --exclude .editorconfig --exclude .git --exclude .github --exclude .gitignore --exclude .npmignore --exclude backup --exclude .replit --exclude install.sh --exclude Backup.tar --exclude app.log --exclude i18n.json " + Exclude_Folders_extra

var Architecture string = runtime.GOARCH
var Local_rclone_route string = "src/bin/"
var Root string = root()
var Local_rclone bool = Configs.Local_rclone
var Include_Folders_extra string = ProsessString(Configs.Include_Folders, "--include ")
var Exclude_Folders_extra string = ProsessString(Configs.Exclude_Folders, "--exclude ")

var Configs = GetJsonData()

type config struct {
	Include_Folders string `json:"include-folders"`
	Exclude_Folders string `json:"exclude-folders"`
	Remote          string `json:"remote"`
	Local_rclone    bool   `json:"local-rclone"`
}

func GetJsonData() config {
	Conf := config{}
	os.Chdir(Root)
	ls, _ := ReadCommand("ls")
	if !strings.Contains(ls, "config.json") {
		log.Warning("config.json does not exist... Creating a new one...")
		NewConFile()
	}
	readfile, err := os.ReadFile("config.json")
	if err != nil {
		log.Error("Error oppening the config file", 23)
	}
	json.Unmarshal(readfile, &Conf)
	return Conf
}

func NewConFile() {
	os.Chdir(Root)
	newdata := config{}
	file, _ := os.Create("config.json")
	data, _ := json.Marshal(newdata)
	file.WriteString(string(data))
	file.Close()
}

func ReadCommand(command string) (string, int) {
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
	if Configs.Exclude_Folders == "" && cond1 == "--exclude " {
		cond1 = ""
	}
	if Configs.Include_Folders == "" && cond1 == "--include " {
		cond1 = ""
	}
	return cond1 + edit(data, " ")
}

func root() string {
	binpath, _ := filepath.Abs(os.Args[0])
	return filepath.Dir(binpath)
}
func remote() string {
	if Configs.Remote != "" {
		return strings.TrimRight(Configs.Remote, "/")
	} else {
		return ""
	}
}
func UpdateJsonData() {
	data, err := json.MarshalIndent(Configs, "", "  ")
	if err != nil {
		log.Error("error when serializing the structure", 22)
	}
	err = os.WriteFile("config.json", data, 0644)
	if err != nil {
		log.Error("Error writing to the config.json file.", 15)
	}
}
