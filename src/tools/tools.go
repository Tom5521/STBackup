package tools

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Tom5521/SillyTavernBackup/src/getdata"
	"github.com/Tom5521/SillyTavernBackup/src/log"
)

var sh = getdata.Sh{}

func Makeconf() {
	os.Chdir(getdata.Root)
	if !CheckDir("config.json") {
		getdata.NewConFile()
	}
	fmt.Print("Enter the rclone Remote server:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	getdata.Configs.Remote = input
	getdata.UpdateJsonData()
	fmt.Printf("Remote Saved in %vYour Remote:%v\n", getdata.Root, input)
	log.Info("Remote Saved\nRemote:'" + input + "'\nRoute:'" + getdata.Root + "'")
	if !CheckRclone() {
		log.Warning("rclone not installed... Using local version")
		sh.Cmd("./backup download-rclone")
	}
}
func Rclone(parameter string) {
	os.Chdir(getdata.Root)
	if getdata.Remote == "" {
		log.Error("Remote dir is null", 9)
	}
	if !CheckRclone() {
		log.Error(
			"Rclone not found. You can download it and use it locally without installing using ./backup download-rclone",
			10,
		)
		return
	}
	if !CheckDir("config.json") {
		Makeconf()
	}
	var loc string
	if getdata.Local_rclone {
		loc = getdata.Local_rclone_route
		if !CheckDir("src/bin/rclone") {
			log.Warning("rclone binary or folders not found!!!")
			os.Chdir(getdata.Root)
			sh.Cmd("./backup download-rclone")
			Rclone(parameter)
			return
		}
	}
	var com = exec.Command("")
	var Remote, Folder string = getdata.Remote, getdata.Folder
	os.Chdir(getdata.Root)
	switch parameter {
	case "uptar":
		log.Func("upload tar")
		com = exec.Command(loc+"rclone", "copy", "Backup.tar", Remote)
		defer log.Info("tar uploaded")
	case "downtar":
		log.Func("download tar")
		com = exec.Command(loc+"rclone", "copy", Remote+"/Backup.tar", "..")
		defer log.Info("tar downloaded")
	case "up":
		log.Func("upload")
		com = exec.Command(loc+"rclone", "sync", Folder, Remote, "-L", "-P")
		defer log.Info("Files uploaded")
	case "down":
		log.Func("download")
		com = exec.Command(loc+"rclone", "sync", Remote, Folder, "-L", "-P")
		defer log.Info("Files downloaded")
	case "ls":
		log.Func("ls")
		com = exec.Command(loc+"rclone", "ls", Remote)
	}
	com.Stderr = os.Stderr
	com.Stdin = os.Stdin
	com.Stdout = os.Stdout
	if getdata.Local_rclone {
		fmt.Println("Using local rclone...")
	}
	com.Run()
}

func WriteFile(name, text string) {
	file, err := os.Create(name)
	if err != nil {
		log.Error("Error creating file in WriteFile func", 24)
		return
	}
	_, err = file.WriteString(text)
	if err != nil {
		log.Error("Error writing in new file (WriteFile func)", 25)
	}
	file.Close()
}
func ReadFileCont(filename string) (string, error) {
	if !CheckDir(filename) {
		log.Warning("File not found in ReadFileCont func")
	}
	cont, err := os.ReadFile(filename)
	if err != nil {
		log.Warning("Error reading file in ReadFileCont func")
		return "", err
	}
	return string(cont), nil
}

func CheckDir(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func CheckRclone() bool {
	if _, rclonestat := sh.Out("rclone version"); rclonestat != nil {
		return false
	} else {
		return true
	}
}
func CheckMainBranch() bool {
	if data1, _ := sh.Out("git status"); strings.Contains(data1, "origin/dev") {
		return false
	} else {
		return true
	}

}
func CheckRsync() {
	if _, rsyncstat := sh.Out("rsync --version"); rsyncstat != nil {
		log.Error("Rsync not found.", 11)
		return
	}
}
