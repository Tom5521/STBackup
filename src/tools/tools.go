package tools

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/Tom5521/SillyTavernBackup/src/checks"
	"github.com/Tom5521/SillyTavernBackup/src/getdata"
	"github.com/Tom5521/SillyTavernBackup/src/log"
)

func Cmd(input string) int {
	cmd := exec.Command("sh", "-c", input)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return 1
	}
	return 0
}

func Makeconf() {
	os.Chdir(getdata.Root)
	if !checks.CheckDir("config.json") {
		getdata.NewConFile()
	}
	fmt.Print("Enter the rclone Remote server:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	getdata.DATA.Remote = input
	getdata.UpdateJsonData()
	fmt.Printf("Remote Saved in %vYour Remote:%v\n", getdata.Root, input)
	log.Info("Remote Saved\nRemote:'" + input + "'\nRoute:'" + getdata.Root + "'")
	if !checks.CheckRclone() {
		log.Warning("rclone not installed... Using local version")
		Cmd("./backup download-rclone")
	}
}
func Rclone(parameter string) {
	if getdata.Remote == "" {
		log.Error("Remote dir is null.", 9)
	}
	if !checks.CheckRclone() {
		log.Error(
			"Rclone not found. You can download it and use it locally without installing using ./backup download-rclone",
			10,
		)
		return
	}
	if !checks.CheckDir("config.json") {
		Makeconf()
	}
	var com = exec.Command("")
	var Remote, Folder string = getdata.Remote, getdata.Folder
	var loc string
	if getdata.Local_rclone {
		loc = getdata.Local_rclone_route
	}
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
	if !checks.CheckDir(filename) {
		log.Warning("File not found in ReadFileCont func")
	}
	cont, err := os.ReadFile(filename)
	if err != nil {
		log.Warning("Error reading file in ReadFileCont func")
		return "", err
	}
	return string(cont), nil
}
