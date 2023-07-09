package tools

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

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
func ReadCommand(command string) (string, int) {
	com := exec.Command("sh", "-c", command)
	data, err := com.Output()
	if err != nil {
		return "", 1
	}
	return string(data), 0
}
func Makeconf() {
	os.Chdir(getdata.Root)
	if !CheckDir("config.json") {
		WriteFile(
			"config.json",
			"{\"local-rclone\":\"\",\"remote\":\"\",\"include-folders\":\"\",\"exclude-folders\":\"\"}",
		)
	}
	fmt.Print("Enter the rclone Remote server:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	UpdateJSONValue("config.json", "remote", input)
	fmt.Printf("Remote Saved in %vYour Remote:%v\n", getdata.Root, input)
	log.Info("Remote Saved\nRemote:'" + input + "'\nRoute:'" + getdata.Root + "'")

}
func Rclone(parameter string) {
	if getdata.Remote == "" {
		log.Error("Remote dir is null.", 9)
	}
	_, err := ReadCommand("rclone version")
	if err == 1 {
		log.Error(
			"Rclone not found. You can download it and use it locally without installing using ./backup download-rclone",
			10,
		)
		return
	}
	if !CheckDir("config.json") {
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
func CheckBranch() bool {
	data1, _ := ReadCommand("git status")
	if strings.Contains(data1, "origin/dev") {
		return false
	}
	return true
}
func CheckRsync() {
	_, rsyncstat := ReadCommand("rsync --version")
	if rsyncstat == 1 {
		log.Error("Rsync not found.", 11)
		return
	}
}
func WriteFile(name, text string) error {
	file, err1 := os.Create(name)
	if err1 != nil {
		return err1
	}
	file.WriteString(text)
	file.Close()
	return err1
}
func CheckDir(dir string) bool {
	data, _ := ReadCommand("ls")
	if strings.Contains(data, dir) {
		return true
	}
	return false

}

func ReadFileCont(filename string) (string, error) {
	cont, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(cont), nil
}
func UpdateJSONValue(filename, variableName, newValue string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		log.Error(fmt.Sprintf("error reading the file: %v", err), 12)
		return err
	}
	data := make(map[string]interface{})
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Error(fmt.Sprintf("error when decoding the JSON file: %v", err), 13)
		return err
	}
	data[variableName] = newValue

	updatedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Error(fmt.Sprintf("error when encoding the JSON file: %v", err), 14)
		return err
	}
	err = os.WriteFile(filename, updatedJSON, 0644)
	if err != nil {
		log.Error(fmt.Sprintf("error when writing the JSON file: %v", err), 15)
		return err
	}
	return nil
}
