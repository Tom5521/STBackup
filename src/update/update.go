package update

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

func DownloadLatestBinary(binName string) int {
	os.Chdir(getdata.Root)
	file, err := os.Create(binName)
	if err != nil {
		log.Error(fmt.Sprintf("Error creating the %s file", binName))
		return 1
	}
	defer file.Close()
	response, err := http.Get("https://github.com/Tom5521/SillyTavernBackup/releases/latest/download/" + binName)
	if err != nil {
		log.Error("Erro performing request")
		return 1
	}
	defer response.Body.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Error("Error copyng the content")
	}
	fmt.Printf("%s downloaded successfully\n", binName)
	tools.Cmd("mv " + binName + " backup -f")
	os.Chmod("backup", 0700)
	return 0
}

func Rebuild() {
	os.Chdir(getdata.Root)
	if !tools.CheckBranch() {
		tools.Cmd("bash build.sh d")
		os.Exit(0)
		return
	}
	log.Func("Rebuild")
	_, errcode := tools.ReadCommand("go version")
	ls, _ := tools.ReadCommand("ls")
	if !strings.Contains(ls, "main.go") {
		log.Error("Source code not found")
	}
	if errcode == 1 {
		log.Error("No go compiler found")
		return
	}
	fmt.Println("Rebuilding...")
	log.Info("Rebuilding")
	err := tools.Cmd("go build -o backup main.go")
	if err != 1 {
		fmt.Println("Rebuild Complete.")
		log.Func("Rebuild Complete.")
		os.Exit(0)
		return
	}
	log.Error("Error in rebuild prosess")
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
	tools.Cmd("go build -o backup main.go")
	fmt.Println("--EMERGENCY REBUILD--")
}
