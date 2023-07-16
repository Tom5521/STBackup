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

var sh getdata.Sh

func DownloadLatestBinary(binName string) int {
	os.Chdir(getdata.Root)
	file, err := os.Create(binName)
	if err != nil {
		log.Error(fmt.Sprintf("Error creating the %s file", binName), 16)
		return 1
	}
	defer file.Close()
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
	sh.Cmd("mv " + binName + " backup -f")
	os.Chmod("backup", 0700)
	return 0
}

func Rebuild() {
	os.Chdir(getdata.Root)
	if !tools.CheckMainBranch() {
		sh.Cmd("bash build.sh d")
		os.Exit(0)
		return
	}
	log.Func("Rebuild")
	_, errcode := sh.Out("go version")
	if !tools.CheckDir("main.go") {
		log.Error("Source code not found", 19)
	}
	if errcode != nil {
		log.Error("No go compiler found", 20)
		return
	}
	fmt.Println("Rebuilding...")
	log.Info("Rebuilding")
	err := sh.Cmd("go build -o backup main.go")
	if err == nil {
		fmt.Println("Rebuild Complete.")
		log.Func("Rebuild Complete.")
		os.Exit(0)
		return
	} else {
		log.Error("Error in rebuild prosess", 21)
	}
}

func RebuildCheck() {
	if len(os.Args) > 2 {
		if os.Args[1] == "rebuild" {
			Rebuild()
		}
	}
}

func EmergencyRebuild() {
	os.Chdir(getdata.Root)
	fmt.Println("--EMERGENCY REBUILD--")
	sh.Cmd("go build -o backup main.go")
	fmt.Println("--EMERGENCY REBUILD--")
}
