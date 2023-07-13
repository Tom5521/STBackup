package checks

import (
	"os"
	"strings"

	"github.com/Tom5521/SillyTavernBackup/src/getdata"
	"github.com/Tom5521/SillyTavernBackup/src/log"
)

var sh = getdata.Sh{}

func CheckDir(dir string) bool {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func CheckRclone() bool {
	_, rclonestat := sh.Out("rclone version")
	if rclonestat != nil {
		return false
	} else {
		return true
	}
}
func CheckBranch() bool {
	data1, _ := sh.Out("git status")
	if strings.Contains(data1, "origin/dev") {
		return false
	} else {
		return true
	}

}
func CheckRsync() {
	_, rsyncstat := sh.Out("rsync --version")
	if rsyncstat != nil {
		log.Error("Rsync not found.", 11)
		return
	}
}
