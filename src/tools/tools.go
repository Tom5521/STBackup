package tools

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Tom5521/SillyTavernBackup/src/getdata"
	"github.com/Tom5521/SillyTavernBackup/src/log"
)

var sh = getdata.Sh{}

// Config the remote rclone dir
func Makeconf() {
	log.Function()
	os.Chdir(getdata.Root)
	// Scan in the terminal the new remote dir value
	fmt.Print("Enter the rclone Remote server:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	getdata.Remote = input
	getdata.Configs.Remote = input // Set the inputed text equal to the local var
	getdata.WriteJsonData()        // Update the config.json data
	// Print in the terminal and write the data in the log
	fmt.Printf("Remote Saved in %vYour Remote:%v\n", getdata.Root, input)
	log.Info("Remote Saved\nRemote:'" + input + "'\nRoute:'" + getdata.Root + "'")
}

// Rclone functions func
func Rclone(parameter string) {
	log.Function()
	os.Chdir(getdata.Root)
	// Check the remote dir
	if getdata.Remote == "" {
		log.Warning("Remote dir is null, set the remote rclone dir value")
		Makeconf()
		Rclone(parameter)
		return
	}
	// Check if rclone is installed
	if !CheckRclone() && !getdata.Local_rclone {
		log.Error(
			"Rclone not found. You can download it and use it locally without installing using ./backup download-rclone",
			10,
		)
		return
	}
	// Check if the local binary or the installed binary will be used
	var loc string
	if getdata.Local_rclone {
		loc = getdata.Local_rclone_route
		// Check if exist the rclone binary
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
	// Check the func will be used for rclone
	switch parameter {
	case "upload tar":
		com = fmt.Sprintf("%vrclone copy Backup.tar %v", loc, Remote)
		defer log.Info("tar uploaded")
	case "download tar":
		com = fmt.Sprintf("%vrclone copy %v/backup.tar ..", loc, Remote)
		defer log.Info("tar downloaded")
	case "upload":
		com = fmt.Sprintf("%vrclone sync %v %v -L -P", loc, Folder, Remote)
		defer log.Info("Files uploaded")
	case "download":
		com = fmt.Sprintf("%vrclone sync %v %v -L -P", loc, Remote, Folder)
		defer log.Info("Files downloaded")
	case "ls":
		com = fmt.Sprintf("%vrclone ls %v", loc, Remote)
	}
	log.Func(parameter)
	if getdata.Local_rclone {
		fmt.Println("Using local rclone...")
	}
	sh.Cmd(com) // Exec the corresponding command
}

// Very descriptive name
func WriteFile(name, text string) {
	log.Function()
	log.Info("Writing %v in %v file...")
	file, err := os.Create(name)
	if err != nil {
		log.Error("Error creating file in WriteFile func", 24)
	}
	_, err = file.WriteString(text)
	if err != nil {
		log.Error("Error writing in new file (WriteFile func)", 25)
	}
	file.Close()
}

// Very descriptive name
func ReadFileCont(filename string) (string, error) {
	log.Function()
	log.Info(fmt.Sprintf("Reading %v file content", filename))
	// Check if the file exists
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

// Very descriptive name
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

// Check if rclone is installed
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

// Check if the program is in the main branch
func CheckMainBranch() bool {
	log.Function()
	// Check if git is installed
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

// Check if rsync is installed
func CheckRsync() {
	log.Function()
	log.Info("Checking rsync")
	if _, rsyncstat := sh.Out("rsync --version"); rsyncstat != nil {
		log.Error("Rsync not found.", 11, "check", string(rsyncstat.Error()))
		return
	}
}

// Check if git is installed checking the .git dir and if git is installed
func CheckGit() bool {
	var check1, check2 bool
	os.Chdir(getdata.Root)
	log.Function()
	log.Info("Checking git")
	if git := CheckDir(".git"); git {
		log.Check(".git dir:true")
		check1 = true
	}
	if _, err := sh.Out("git status"); err == nil {
		log.Check("git installed:true")
		check2 = true
	}
	return check1 && check2
}

// Start sillytavern
func SillyTavern(input string) {
	// Set the executable to execute
	var par, command string
	if input == "start" {
		par = "start"
		command = "node server.js"
	}
	if input == "init" {
		par = "init"
		command = "bash start.sh"
	}
	// Make a new channel to detect sigint and sigterm
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	os.Chdir("..")
	log.Func(par)
	sh.Cmd(command) // Exec the current file
	sig := <-sigChan
	switch sig {
	case syscall.SIGINT:
		fmt.Println("SIGINT (Ctrl+C) was received. The program will close.")
		log.Info("SIGINT")
	case syscall.SIGTERM:
		fmt.Println("SIGTERM was received. The program will be closed.")
		log.Info("SIGTERM")
	}
}

func Rsync(pars ...string) {
	if len(pars) < 1 {
		return
	}
	var par1, func1 string = "--delete", ""
	if len(pars) == 2 {
		if pars[1] == "secure" {
			par1 = ""
			func1 = " " + pars[1]
		}
	}
	os.Chdir("..")
	log.Func(func1 + pars[0])
	if pars[0] == "save" {
		CheckRsync()
		sh.Cmd(
			fmt.Sprintf(
				"rsync -av --progress %s %s . %s",
				par1,
				getdata.Exclude_Folders,
				getdata.Back,
			),
		)
		log.Info("Files Saved")
		if len(os.Args) == 3 { // Check if tar arg is on
			if os.Args[2] == "tar" {
				log.Func("save tarball")
				tar := sh.Cmd("tar -cvf Backup.tar Backup/")
				if tar == nil {
					log.Info("Tarbal created.")
				}
			}
		}
		os.Chdir(getdata.Root)
	}
	log.Func("restore")
	CheckRsync()
	if len(os.Args) == 3 {
		if os.Args[2] == "tar" {
			log.Func("restore from tarball")
			if CheckDir("Backup") {
				log.Warning("Removing Backup/ folder")
				sh.Cmd("rm -rf Backup/")
			}
			sh.Cmd("tar -xvf Backup.tar")
		}
	}
	sh.Cmd(
		fmt.Sprintf(
			"rsync -av --progress %s %s %s %s .",
			par1,
			getdata.Exclude_Folders,
			getdata.Include_Folders,
			getdata.Back,
		),
	)
	os.Chdir(getdata.Root)
	log.Info("Files restored")
}
