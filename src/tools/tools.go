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

// Important functions
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

func Readconf() (string, error) {
	os.Chdir(getdata.Root)
	file, err := os.Open("config.json")
	if err != nil {
		return "", err
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	var config map[string]interface{}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return "", err
	}
	Remote := config["Remote"].(string)
	if Remote == "" {
		fmt.Println("Remote is empty.")
		log.Warning("Remote is empty.")
		return "", nil
	}
	return Remote, nil
}

func Makeconf() error {
	os.Chdir(getdata.Root)
	Cmd("echo '{\"remote\":\"\"}' > config.json")
	fmt.Print("Enter the rclone Remote server:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	pwd, _ := ReadCommand("pwd")
	fmt.Printf("Remote Saved in %vYour Remote:%v\n", pwd, input)
	log.Info("Remote Saved\nRemote:'" + input + "'\nRoute:'" + pwd + "'")

	file, err := os.OpenFile("config.json", os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	var config map[string]interface{}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return err
	}
	config["remote"] = input
	bytes, err = json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("config.json", bytes, 0644)
	if err != nil {
		return err
	}
	return nil
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
	lsstat, _ := ReadCommand("ls")
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
