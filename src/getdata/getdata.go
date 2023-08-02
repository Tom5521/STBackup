package getdata

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Tom5521/STBackup/src/log"
)

// Declare the constants
const (
	Folder, Back string = "../Backup/", "Backup/"
	Version      string = "2.7.3"
	// Declare the default folders of sillytavern to make backup
	Def_include_folders string = `backgrounds "group chats" "KoboldAI Settings" settings.json characters groups notes sounds worlds chats "NovelAI Settings" img "OpenAI Settings" "TextGen Settings" themes "User Avatars" secrets.json thumbnails config.conf public uploads backups default instruct context stats.json  `

	// Declare the default folders of sillytavern to exclude
	Def_exclude_folders string = "webfonts scripts index.html css img favicon.ico script.js style.css Backup colab docker Dockerfile LICENSE node_modules package.json package-lock.json replit.nix server.js SillyTavernBackup src Start.bat start.sh UpdateAndStart.bat Update-Instructions.txt tools .dockerignore .editorconfig .git .github .gitignore .npmignore .replit install.sh Backup.tar app.log i18n.json stbackup STbackup STBackup statsHelpers.js poe-test.js poe-error.log poe_device.json poe-success.log "
	// Get the architecture
	Architecture string = runtime.GOARCH
)

// Declare the variables
var (
	BinName string = func() string {
		biname := os.Args[0]
		return filepath.Base(biname)
	}()
	// Remote the final "/" in remote dir if it exist
	Remote string = func() string {
		if Configs.Remote != "" {
			return strings.TrimRight(Configs.Remote, "/")
		} else {
			return ""
		}
	}()
	// Add exclude/include prefix to the rclone syntax + exclude/include in extra in config.json
	Exclude_Folders string = AddPrefix(
		Def_exclude_folders,
		"--exclude ",
	) + AddPrefix(
		Configs.Exclude_Folders,
		"--exclude ",
	)

	Include_Folders string = AddPrefix(
		Def_include_folders,
		"--include ",
	) + AddPrefix(
		Configs.Include_Folders,
		"--include ",
	)

	// Set the rclone binary route
	Local_rclone_route string = Root + "/src/bin/"

	// Set the root local dir and get the root directory
	Root string = func() string {
		binpath, _ := filepath.Abs(os.Args[0])
		return filepath.Dir(binpath)
	}()

	// get local rclone value true/false
	Local_rclone bool = Configs.Local_rclone

	// Get the config.json data
	Configs = GetConfig()

	// Initialize the shell functions
	sh = Sh{}
)

// Declare the struct of json file
type config struct {
	Include_Folders string `json:"include-folders"`
	Exclude_Folders string `json:"exclude-folders"`
	Remote          string `json:"remote"`
	Local_rclone    bool   `json:"local-rclone"`
	Loglevel        int    `json:"log-level"`
}

func GetSTversion() string {
	type STinfo struct {
		Version string `json:"version"`
	}
	info1 := STinfo{}
	os.Chdir(Root)
	os.Chdir("..")
	ls, _ := sh.Out("ls")
	if !strings.Contains(ls, "package.json") {
		log.Error("package.json file not found!", 27)
	}
	readfile, err := os.ReadFile("package.json")
	if err != nil {
		log.Error("Error reading package.json file", 28)
	}
	json.Unmarshal(readfile, &info1)
	return info1.Version
}

// Declare the struct of shell functions
type Sh struct{}

// Get json data function,name very descriptive
func GetConfig() config {
	Conf := config{}
	os.Chdir(Root)
	ls, _ := sh.Out("ls")
	// Check if config.json exist
	if !strings.Contains(ls, "config.json") {
		log.Warning("config.json does not exist... Creating a new one...")
		NewConFile() // Create new config file
	}
	// Read config file
	readfile, err := os.ReadFile("config.json")
	if err != nil {
		log.Error("Error oppening the config file", 23)
	}
	// Prosess the data of the config file and retuns it
	json.Unmarshal(readfile, &Conf)
	return Conf
}

// Very descriptive name
func NewConFile() {
	log.Function()
	os.Chdir(Root)
	newdata := config{}                 // Initialize the config struct
	file, _ := os.Create("config.json") // Create config.json file
	defer file.Close()                  // Close the file in the end
	data, _ := json.Marshal(newdata)    // Marshall the data
	file.WriteString(string(data))      // Write the data in config.json file
}

// Name very descriptive,prossess the strings in config.json for use in the Include_Folders and Exclude_Folders vars
func AddPrefix(input, prefix string) string {
	if input == "" {
		return ""
	}
	words := strings.Fields(input)
	var result strings.Builder

	for _, word := range words {
		result.WriteString(prefix)
		result.WriteString(word)
		result.WriteString(" ")
	}
	return strings.TrimSpace(result.String())
}

// Update the config.json data with the changes made for the program in its values
func WriteJsonData() {
	data, err := json.MarshalIndent(Configs, "", "  ")
	if err != nil {
		log.Error("error when serializing the structure", 22)
	}
	err = os.WriteFile("config.json", data, os.ModePerm)
	if err != nil {
		log.Error("Error writing to the config.json file.", 15)
	}
}

// Exec shell command func
func (sh Sh) Cmd(input string) error {
	shell := make([]string, 2)
	if runtime.GOOS == "windows" {
		shell[0] = "cmd"
		shell[1] = "/C"
	}
	if runtime.GOOS == "linux" || runtime.GOOS == "android" {
		shell[0] = "sh"
		shell[1] = "-c"
	}
	cmd := exec.Command(shell[0], shell[1], input)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	} else {
		return nil
	}
}

// Get the output and the error executing a shell command
func (sh Sh) Out(input string) (string, error) {
	cmd := exec.Command("sh", "-c", input)
	out, err := cmd.Output()
	if err != nil {
		return string(out), err
	} else {
		return string(out), nil
	}
}

func SendLogLv() {
	log.TempChan1 <- Configs.Loglevel
}

// Help string = README.md
const Help string = `
# Silly Tavern Backup and Cloud Upload

This is a program, which provides a backup and restore tool for SillyTavern. The program uses the rsync command to synchronize the application files between the local server and the remote server. It also uses the rclone tool to synchronize SillyTavern files with a cloud storage service.

Commands

[program] <parameter> <parameter> | ej: stbackup save tar

- make: Creates necessary folders for backup.
- save: Saves files to the backup destination.
- save tar: Saves files to a tarball in the backup destination.
- secure-save: save without delete any file of Backup folder
- restore: Restores files from the backup destination.
- secure-restore: Restore without delete any file (Use it in case you have updated sillytavern and want to restore the data without breaking the repo and having to clone it again.)
- restore tar: Restores files from a tarball in the backup destination.
- route <destination>: Moves the backup folder to a new destination.
- start: Starts the SillyTavern application.
- update ST: Updates the ST application.
- update me: Updates the STBackup application and rebuilds if necessary.
- ls: Lists files in the remote backup destination.
- upload: Uploads files to the remote backup destination.
- upload tar: Uploads a tarball to the remote backup destination.
- download: Downloads files from the remote backup destination.
- download tar: Downloads a tarball from the remote backup destination.
- init: Initializes the SillyTavern application.
- rebuild: This is a high priority command. It will rebuild the program (if you have the source code at hand) rather than run any other function than the logs and change to the root directory. As soon as it finishes executing the program it will terminate with error code 0.
- link: Creates a link to the backup program in the SillyTavern root directory.
- version: Displays the version of STbackup.
- remote: Configures the rclone remote server.
- cleanlog: Clears the log file.
- log: Displays the content of the log file.
- help: Displays the help message.
- printconf:Print config.json file
- test: Only works in the dev branch. I use it to debug the code
- resetconf: It is used to delete all program settings. Does not delete backups
- download-rclone: Download and extract the rclone binary, it is useful if the rclone installed in your distro does not work correctly or if you do not want to install it.
- Setloglevel: is used to set the log level in app.log interactively.
`
