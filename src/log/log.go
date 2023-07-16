package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var binpath, _ = filepath.Abs(os.Args[0])
var root string = filepath.Dir(binpath)
var logger = SetupLogger(root + "/app.log")

func Error(text string, errcode int, incheck ...string) {
	var check string
	if len(incheck) == 2 {
		if incheck[0] == "check" {
			check = "| CHECK: " + incheck[1]
		}
	}
	pc, _, _, _ := runtime.Caller(1)
	prefuncname := runtime.FuncForPC(pc).Name()
	parts := strings.Split(prefuncname, "/")
	funcname := parts[len(parts)-1]
	errdata := fmt.Sprintf(
		"ERROR: %s | code: %d | file: %v "+check,
		text,
		errcode,
		funcname,
	)
	fmt.Println(errdata)
	logger.Println(errdata)
	Info("---------End---------")
	os.Exit(errcode)
}
func Warning(text string) {
	_, fPath, _, _ := runtime.Caller(1)
	filePath := filepath.Base(fPath)
	warndata := fmt.Sprintf("WARNING: %s | file: %v", text, filePath)
	fmt.Println(warndata)
	logger.Println(warndata)
}
func Info(text string) {
	var par string
	if text == "--------Start--------" {
		par = "\n\n^^^^^^^^^^^^^^-time "
	}
	logger.Println(par + "PROGRAM: " + text)
}
func Func(text string) {
	logger.Println("FUNC:    ---" + text + "---")
}
func Function() {
	pc, _, _, _ := runtime.Caller(1)
	prefuncname := runtime.FuncForPC(pc).Name()
	parts := strings.Split(prefuncname, "/")
	funcname := parts[len(parts)-1]
	logger.Println("FUNCTION:    ---" + funcname + "---")
}
func Check(input string) {
	pc, _, _, _ := runtime.Caller(1)
	prefuncname := runtime.FuncForPC(pc).Name()
	parts := strings.Split(prefuncname, "/")
	funcname := parts[len(parts)-1]
	logger.Printf("CHECK:		 %v:%v \n", funcname, input)
}
func SetupLogger(logFilePath string) *log.Logger {
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(file, "", log.Ldate|log.Ltime)
	return logger
}
