package tools

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	ls := ReadDir()
	if !strings.Contains(ls, "config.json") {
		WriteFile("config.json", "{\"remote\":\"\",\"include-folders\":\"\",\"exclude-folders\":\"\"}")
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
		log.Error("Remote dir is null.")
	}
	_, err := ReadCommand("rclone version")
	if err == 1 {
		log.Error("Rclone not found")
		return
	}
	lsstat := ReadDir()
	if !strings.Contains(lsstat, "config.json") {
		Makeconf()
	}
	var com = exec.Command("")
	var Remote, Folder string = getdata.Remote, getdata.Folder
	switch parameter {
	case "uptar":
		log.Func("upload tar")
		com = exec.Command("rclone", "copy", "Backup.tar", Remote)
		defer log.Info("tar uploaded")
	case "downtar":
		log.Func("download tar")
		com = exec.Command("rclone", "copy", Remote+"/Backup.tar", "..")
		defer log.Info("tar downloaded")
	case "up":
		log.Func("upload")
		com = exec.Command("rclone", "sync", Folder, Remote, "-L", "-P")
		defer log.Info("Files uploaded")
	case "down":
		log.Func("download")
		com = exec.Command("rclone", "sync", Remote, Folder, "-L", "-P")
		defer log.Info("Files downloaded")
	}
	com.Stderr = os.Stderr
	com.Stdin = os.Stdin
	com.Stdout = os.Stdout
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
		log.Error("Rsync not found.")
		return
	}
}

func WriteFile(name, text string) error {
	file, err1 := os.Create(name)
	defer file.Close()
	if err1 != nil {
		return err1
	}
	file.WriteString(text)
	return err1
}

func ReadDir() string {
	data, _ := ReadCommand("ls")
	return data
}

func ReadFileCont(filename string) (string, error) {
	cont, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(cont), nil
}
func UpdateJSONValue(filename, variableName, newValue string) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error(fmt.Sprintf("error reading the file: %v", err))
		return err
	}
	data := make(map[string]interface{})
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Error(fmt.Sprintf("error when decoding the JSON file: %v", err))
		return err
	}
	data[variableName] = newValue

	updatedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Error(fmt.Sprintf("error when encoding the JSON file: %v", err))
		return err
	}
	err = ioutil.WriteFile(filename, updatedJSON, 0644)
	if err != nil {
		log.Error(fmt.Sprintf("error when writing the JSON file: %v", err))
		return err
	}
	return nil
}
