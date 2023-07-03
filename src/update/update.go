package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Tom5521/SillyTavernBackup/src/getdata"
	"github.com/Tom5521/SillyTavernBackup/src/log"
	"github.com/Tom5521/SillyTavernBackup/src/tools"
)

var Root string = getdata.Root

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
		return fmt.Errorf("Failed to find %s binary in the latest version of %s", binName, repo)
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
	tools.Cmd("mv " + fileName + " backup")
	os.Chmod("backup", 0700)
}

func Rebuild() {
	os.Chdir(Root)
	log.Logfunc("Rebuild")
	_, errcode := tools.ReadCommand("go version")
	ls, _ := tools.ReadCommand("ls")
	if !strings.Contains(ls, "main.go") {
		fmt.Println("Source code not found")
		log.Logerror("Source code not found")
	}
	if errcode == 1 {
		fmt.Println("No go compiler found")
		log.Logerror("No go compiler found")
		return
	}
	fmt.Println("Rebuilding...")
	err := tools.Cmd("go build -o backup main.go")
	if err != 1 {
		fmt.Println("Rebuild Complete.")
		log.Logfunc("Rebuilded")
		return
	}
	log.Logerror("Error in rebuild prosess")
}
