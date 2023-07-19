package update

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Tom5521/SillyTavernBackup/src/getdata"
	"github.com/Tom5521/SillyTavernBackup/src/log"
	"github.com/Tom5521/SillyTavernBackup/src/tools"
)

var sh getdata.Sh // Init the shell func

// Very descriptive name
func DownloadLatestBinary(binName string) int {
	log.Function()
	os.Chdir(getdata.Root)
	file, err := os.Create(binName) // Create the file
	if err != nil {
		log.Error(fmt.Sprintf("Error creating the %s file", binName), 16)
		return 1
	}
	defer file.Close()
	// Set the current url to download the binary
	response, err := http.Get(
		"https://github.com/Tom5521/SillyTavernBackup/releases/latest/download/" + binName,
	)
	if err != nil {
		log.Error("Error performing request", 17)
		return 1
	}
	defer response.Body.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Error("Error copyng the content", 18)
	}
	fmt.Printf("%s downloaded successfully\n", binName)
	sh.Cmd("mv " + binName + " backup -f") // Rename the binary file
	os.Chmod("backup", 0700)               // Give exec permissions to the downloaded file
	return 0
}

// Very descriptive name
func Rebuild() {
	os.Chdir(getdata.Root)
	// Check if is in the dev branch for develop build
	if !tools.CheckMainBranch() {
		sh.Cmd("bash build.sh d")
		os.Exit(0)
		return
	}
	log.Function()
	_, errcode := sh.Out("go version") // Check if the go compiler is installed
	if !tools.CheckDir("main.go") {
		log.Error("Source code not found", 19)
	}
	if errcode != nil {
		log.Error("No go compiler found", 20)
		return
	}
	fmt.Println("Rebuilding...")
	log.Info("Rebuilding")
	err := sh.Cmd("go build -o backup main.go") // Rebuilds the program
	if err == nil {
		fmt.Println("Rebuild Complete.")
		log.Func("Rebuild Complete.")
		os.Exit(0)
		return
	} else {
		log.Error("Error in rebuild prosess", 21)
	}
}

// Check if rebuild functions is called in the terminal (rebuild is a top level func)
func RebuildCheck() {
	if len(os.Args) > 2 {
		if os.Args[1] == "rebuild" {
			Rebuild()
		}
	}
}
