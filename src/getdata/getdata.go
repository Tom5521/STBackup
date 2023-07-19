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
const Version string = "2.4.1"

var Remote string = remote()

// Declare the default folders of sillytavern to make backup + the extra include folders in config.json
var Include_Folders string = "--include backgrounds/ --include 'group chats' --include 'KoboldAI Settings' --include settings.json --include characters --include groups --include notes --include sounds --include worlds --include chats --include 'NovelAI Settings' --include img --include 'OpenAI Settings' --include 'TextGen Settings' --include themes --include 'User Avatars' --include secrets.json --include thumbnails --include config.conf --include poe_device.json --include public --include uploads --include backups " + Include_Folders_extra

// Declare the default folders of sillytavern to exclude + the extra exclude folders in config.json
var Exclude_Folders string = "--exclude webfonts --exclude scripts --exclude index.html --exclude css --exclude img --exclude favicon.ico --exclude script.js --exclude style.css --exclude Backup --exclude colab --exclude docker --exclude Dockerfile --exclude LICENSE --exclude node_modules --exclude package.json --exclude package-lock.json --exclude replit.nix --exclude server.js --exclude SillyTavernBackup --exclude src --exclude Start.bat --exclude start.sh --exclude UpdateAndStart.bat --exclude Update-Instructions.txt --exclude tools --exclude .dockerignore --exclude .editorconfig --exclude .git --exclude .github --exclude .gitignore --exclude .npmignore --exclude backup --exclude .replit --exclude install.sh --exclude Backup.tar --exclude app.log --exclude i18n.json " + Exclude_Folders_extra

// Get the architecture
const Architecture string = runtime.GOARCH

// Set the rclone binary route
const Local_rclone_route string = "src/bin/"

// Set the root local dir
var Root string = root()

// get local rclone value true/false
var Local_rclone bool = Configs.Local_rclone

// Prosess the strings of config.json to adapt then to the rsync syntax
var Include_Folders_extra string = ProsessString(Configs.Include_Folders, "--include ")
var Exclude_Folders_extra string = ProsessString(Configs.Exclude_Folders, "--exclude ")

// Get the config.json data
var Configs = GetJsonData()

// Initialize the shell functions
var sh = Sh{}

// Declare the struct of json file
type config struct {
	Include_Folders string `json:"include-folders"`
	Exclude_Folders string `json:"exclude-folders"`
	Remote          string `json:"remote"`
	Local_rclone    bool   `json:"local-rclone"`
}

// Declare the struct of shell functions
type Sh struct{}

// Get json data function,name very descriptive
func GetJsonData() config {
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

// Get the root directory
func root() string {
	binpath, _ := filepath.Abs(os.Args[0])
	return filepath.Dir(binpath)
}

// Remote the final "/" in remote dir if it exist
func remote() string {
	if Configs.Remote != "" {
		return strings.TrimRight(Configs.Remote, "/")
	} else {
		return ""
	}
}

// Update the config.json data with the changes made for the program in its values
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

// Exec shell command func
func (sh Sh) Cmd(input string) error {
	cmd := exec.Command("sh", "-c", input)
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

// Help string = README.md
const Help string = `
# Silly Tavern Backup and Cloud Upload

This is a source code file written in the Go programming language, which provides a backup and restore tool for SillyTavern. The program uses the rsync command to synchronize the application files between the local server and the remote server. It also uses the rclone tool to synchronize SillyTavern files with a cloud storage service.
## Requirements

- ### Required
- rsync
- rclone (Optional installation)
- #### Optional
- Go (if you want to compile it yourself to your system)
- Tar (if you want use tarballs)
- unzip (If you do not want to install rclone)

## Installation
### Build Method
1. Clone this repository into the sillytavern folder
2. Open a terminal and navigate to the folder that contains the source code file.
3. Run the following command to compile the program:

bash
go build -o backup main.go

1. Once compiled, you can use the program by running the backup binary file in the same folder as the source code file.
2. (Optional) You can make a ./backup link to be able to run the script from the root of SillyTavern and not need to enter the binary folder. This process is done automatically in the script
### Script Method
Using the script. Below is how to use it
It is useful for those who do not like to compile and want the binary once and for all.
## Usage

The program is run from the command line and accepts various commands and options.


### Commands

- make: Creates necessary folders for backup.
- save: Saves files to the backup destination.
- save tar: Saves files to a tarball in the backup destination.
- restore: Restores files from the backup destination.
- restore tar: Restores files from a tarball in the backup destination.
- route <destination>: Moves the backup folder to a new destination.
- start: Starts the SillyTavern application.
- update ST: Updates the SillyTavernBackup application.
- update me: Updates the SillyTavernBackup application and rebuilds if necessary.
- ls: Lists files in the remote backup destination.
- upload: Uploads files to the remote backup destination.
- upload tar: Uploads a tarball to the remote backup destination.
- download: Downloads files from the remote backup destination.
- download tar: Downloads a tarball from the remote backup destination.
- init: Initializes the SillyTavern application.
- rebuild: This is a high priority command. It will rebuild the program (if you have the source code at hand) rather than run any other function than the logs and change to the root directory. As soon as it finishes executing the program it will terminate with error code 0.
- link: Creates a link to the backup program in the SillyTavern root directory.
- version: Displays the version of SillyTavernBackup.
- remote: Configures the rclone remote server.
- cleanlog: Clears the log file.
- log: Displays the content of the log file.
- help: Displays the help message.
- printconf:Print config.json file
- test: Only works in the dev branch. I use it to debug the code
- resetconf: It is used to delete all program settings. Does not delete backups
- download-rclone: Download and extract the rclone binary, it is useful if the rclone installed in your distro does not work correctly or if you do not want to install it.
### Configuration
The configuration is located in the config.json file.
its parameters are:
1. remote: This parameter determines the path to the remote rclone server.
2. include-folders: This parameter adds folders to be included in the backup.
3. exclude-folders: This parameter adds folders to exclude in the backup (you are free to use it if it takes me too long to update SillyTavern)
4. local-rclone: Is used to determine whether to use a local rclone binary or the one that comes installed with the system, its possible values are yes and no.
### Log

The program logs all the actions in the app.log file.

# SillyTavernBackup Installer

This script is an installer for the SillyTavernBackup program. This is a simple script for those who do not want to compile.
## Usage

To use this script, follow these steps:

1. Copy the installation file to the SillyTavern folder.

2. Run the installation file with the following command:


bash install.sh <platform>


Where <platform> is the platform on which you want to install the program. You can use the following values:

- arm: if you want to install the program on an Android device using the Termux app.
- x64: if you want to install the program on a computer with x86-64 architecture.

3. If you want to install the program on a platform other than arm or x64, you can modify the script to add support for that platform.

## Functionality

The script works as follows:

- If a platform is specified as an argument and it is not clone, the script creates a folder called "SillyTavernBackup" and downloads the latest version of the program corresponding to the specified platform using the GitHub API. Then, it renames the downloaded file to "backup" and executes it.
- If "clone" is specified as an argument, the script clones the GitHub repository and compiles the program using the go build command.
- If no arguments are specified, the script does nothing.

In any case, the program is  and the backup or restore process is started as appropriate.

## Notes

- This script requires the prior installation of the curl tool.
- To use this script on other platforms, it is necessary to modify the file to add support for the specific platform.
- This is for termux and linux friends. If you use windows i recommend you to learn how to use [Rclone](https://rclone.org/) and [SillyTavernSimpleLauncher](https://github.com/BlueprintCoding/SillyTavernSimpleLauncher).
- I don't plan to make a windows version... Unless I'm bored. We'll see what I feel like doing.
- *DO NOT MOVE THE BINARY OUTSIDE THE "SillyTavernBackup" FOLDER.*
## License

This project is licensed under the MIT License. See the https://github.com/Tom5521/SillyTavernBackup/blob/main/LICENSE file for more details.


### At what point did this project go from being a personal bullshit to a moderately serious project?
`
