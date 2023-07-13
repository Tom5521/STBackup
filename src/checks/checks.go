package checks

import (
	"os"
	"strings"

	"github.com/Tom5521/SillyTavernBackup/src/getdata"
	"github.com/Tom5521/SillyTavernBackup/src/log"
)

func CheckDir(dir string) bool {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func CheckRclone() bool {
	_, rclonestat := getdata.ReadCommand("rclone version")
	if rclonestat == 1 {
		return false
	} else {
		return true
	}
}
func CheckBranch() bool {
	data1, _ := getdata.ReadCommand("git status")
	if strings.Contains(data1, "origin/dev") {
		return false
	} else {
		return true
	}

}
func CheckRsync() {
	_, rsyncstat := getdata.ReadCommand("rsync --version")
	if rsyncstat == 1 {
		log.Error("Rsync not found.", 11)
		return
	}
}
