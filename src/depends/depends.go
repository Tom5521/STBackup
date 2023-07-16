package depends

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Tom5521/SillyTavernBackup/src/getdata"
	"github.com/Tom5521/SillyTavernBackup/src/log"
	"github.com/Tom5521/SillyTavernBackup/src/tools"
)

var sh = getdata.Sh{}

func DownloadRclone() {
	log.Function()
	var arch string = getdata.Architecture
	var link_linux_amd64 string = "https://github.com/rclone/rclone/releases/download/v1.63.0/rclone-v1.63.0-linux-amd64.zip"
	var link_linux_386 string = "https://github.com/rclone/rclone/releases/download/v1.63.0/rclone-v1.63.0-linux-386.zip"
	var link_universal_linux_arm string = "https://github.com/rclone/rclone/releases/download/v1.63.0/rclone-v1.63.0-linux-arm.zip"

	var link string
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
	if tools.CheckRclone() && !getdata.Local_rclone {
		return
	}
	os.Chdir(getdata.Root)
	if !tools.CheckDir("src") {
		os.Mkdir("src", 0700)
	}
	if !tools.CheckDir("src/bin") {
		os.Mkdir("src/bin", 0700)
	}
	os.Chdir("src/bin/")
	if tools.CheckDir("rclone") || tools.CheckDir("rclone.zip") {
		log.Warning("rclone already downloaded.")
		return
	}
	DownloadBinaries("rclone.zip", link)
	sh.Cmd("unzip rclone.zip -d rclone-zip")
	sh.Cmd(fmt.Sprintf("cp rclone-zip/rclone-v1.63.0-linux-%s/rclone .", arch))
	sh.Cmd("rm -rf rclone-zip rclone.zip")
	os.Chdir(getdata.Root)
	getdata.Configs.Local_rclone = true
	getdata.UpdateJsonData()
}

func DownloadBinaries(filepath, url string) int {
	log.Function()
	file, err := os.Create(filepath)
	if err != nil {
		log.Error(fmt.Sprintf("Error creating the %s file", filepath), 4)
		return 1
	}
	defer file.Close()
	response, err := http.Get(url)
	if err != nil {
		log.Error("Error performing request", 5)
		return 1
	}
	defer response.Body.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Error("Error copyng the content", 6)
	}
	fmt.Printf("%s downloaded successfully\n", filepath)
	return 0

}
