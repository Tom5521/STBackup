package tools

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Tom5521/SillyTavernBackup/src/getdata"
	"github.com/Tom5521/SillyTavernBackup/src/log"
)

var sh = getdata.Sh{}

func Makeconf() {
	log.Function()
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
	log.Function()
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
	var com string
	var Remote, Folder string = getdata.Remote, getdata.Folder
	os.Chdir(getdata.Root)
	switch parameter {
	case "uptar":
		log.Func("upload tar")
		com = fmt.Sprintf("%vrclone copy Backup.tar %v", loc, Remote)
		defer log.Info("tar uploaded")
	case "downtar":
		log.Func("download tar")
		com = fmt.Sprintf("%vrclone copy %v/backup.tar ..", loc, Remote)
		defer log.Info("tar downloaded")
	case "up":
		log.Func("upload")
		com = fmt.Sprintf("%vrclone sync %v %v -L -P", loc, Folder, Remote)
		defer log.Info("Files uploaded")
	case "down":
		log.Func("download")
		com = fmt.Sprintf("%vrclone sync %v %v -L -P", loc, Remote, Folder)
		defer log.Info("Files downloaded")
	case "ls":
		log.Func("ls")
		com = fmt.Sprintf("%vrclone ls %v", loc, Remote)
	}
	if getdata.Local_rclone {
		fmt.Println("Using local rclone...")
	}
	sh.Cmd(com)
}

func WriteFile(name, text string) {
	log.Function()
	log.Info("Writing %v in %v file...")
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
	log.Function()
	log.Info(fmt.Sprintf("Reading %v file content", filename))
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
	log.Function()
	log.Info(fmt.Sprintf("Checking %v dir", dir))
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Check("false")
		return false
	} else {
		log.Check("true")
		return true
	}
}

func CheckRclone() bool {
	log.Function()
	log.Info("Checking rclone")
	if _, rclonestat := sh.Out("rclone version"); rclonestat != nil {
		log.Check("false")
		return false
	} else {
		log.Check("true")
		return true
	}
}
func CheckMainBranch() bool {
	log.Function()
	if !CheckGit() {
		return true
	}
	log.Info("Checking branch")
	if data1, _ := sh.Out("git status"); strings.Contains(data1, "origin/dev") {
		log.Check("false")
		return false
	} else {
		log.Check("true")
		return true
	}

}
func CheckRsync() {
	log.Function()
	log.Info("Checking rsync")
	if _, rsyncstat := sh.Out("rsync --version"); rsyncstat != nil {
		log.Error("Rsync not found.", 11, "check", string(rsyncstat.Error()))
		return
	}
}

func CheckGit() bool {
	os.Chdir(getdata.Root)
	log.Function()
	log.Info("Checking git")
	if git := CheckDir(".git"); git {
		log.Check("true")
		return true
	} else {
		log.Check("false")
		return false
	}
}
