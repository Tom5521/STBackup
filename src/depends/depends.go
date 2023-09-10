package depends

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Tom5521/MyGolangTools/commands"
	"github.com/Tom5521/STBackup/src/getdata"
	"github.com/Tom5521/STBackup/src/log"
	"github.com/Tom5521/STBackup/src/tools"
)

// Declare private shell functions
var sh = commands.Sh{}

func DownloadRclone() {
	log.Function()
	// Declare the vars
	var (
		arch                     string = getdata.Architecture
		link_linux_amd64         string = "https://github.com/rclone/rclone/releases/download/v1.63.0/rclone-v1.63.0-linux-amd64.zip"
		link_linux_386           string = "https://github.com/rclone/rclone/releases/download/v1.63.0/rclone-v1.63.0-linux-386.zip"
		link_universal_linux_arm string = "https://github.com/rclone/rclone/releases/download/v1.63.0/rclone-v1.63.0-linux-arm.zip"
		// Set link to download
		link string
	)
	if getdata.Architecture == "amd64" {
		link = link_linux_amd64
		arch = getdata.Architecture
	}
	if getdata.Architecture == "386" {
		link = link_linux_386
		arch = getdata.Architecture
	}
	if strings.Contains(getdata.Architecture, "arm") ||
		strings.Contains(getdata.Architecture, "aarch64") {
		link = link_universal_linux_arm
		arch = "arm"
	}
	if tools.CheckRclone() &&
		!getdata.Local_rclone { //Check if rclone is installed and if local-rclone var is true for automatic download binary
		return
	}
	//Check if the necessary folders exist.
	os.Chdir(getdata.Root)
	if !tools.CheckDir("src") {
		os.Mkdir("src", 0700)
	}
	if !tools.CheckDir("src/bin") {
		os.Mkdir("src/bin", 0700)
	}
	os.Chdir(getdata.Local_rclone_route)
	if tools.CheckDir("rclone") || tools.CheckDir("rclone.zip") {
		log.Warning("rclone already downloaded.")
		return
	}
	// Download the corresponding zip files with the binaries
	DownloadBinaries("rclone.zip", link)
	// Unzip the .zip file previously downloaded
	sh.Cmd("unzip rclone.zip -d rclone-zip")
	// Copy the binary to bin folder
	sh.Cmd(fmt.Sprintf("cp rclone-zip/rclone-v1.63.0-linux-%s/rclone .", arch))
	// Remove the uninteresting zip file
	sh.Cmd("rm -rf rclone-zip rclone.zip")
	os.Chdir(getdata.Root)
	// Set local-rclone var with true
	getdata.Configs.Local_rclone = true
	getdata.WriteJsonData()
}

func DownloadBinaries(filepath, url string) int {
	log.Function()
	// Create the correspondig zip file
	file, err := os.Create(filepath)
	if err != nil {
		log.Error(fmt.Sprintf("Error creating the %s file", filepath), 4)
		return 1
	}
	defer file.Close()
	// Get the data
	response, err := http.Get(url)
	if err != nil {
		log.Error("Error performing request", 5)
		return 1
	}
	defer response.Body.Close()
	// Copy the data to the .zip file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Error("Error copyng the content", 6)
	}
	fmt.Printf("%s downloaded successfully\n", filepath)
	return 0

}
