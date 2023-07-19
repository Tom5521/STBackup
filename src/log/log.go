package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Get the root directory
var binpath, _ = filepath.Abs(os.Args[0])
var root string = filepath.Dir(binpath)
var TempChan1 = make(chan int)
var logger = SetupLogger(root + "/app.log") // Initialize the logger func
var Loglevel int

func Error(text string, errcode int, incheck ...string) {
	var check string
	// Check if the "check" func will be used
	if len(incheck) == 2 {
		if incheck[0] == "check" {
			check = "| CHECK: " + incheck[1]
		}
	}
	// Get the name of the function and the line of it
	pc, _, line, _ := runtime.Caller(1)
	prefuncname := runtime.FuncForPC(pc).Name()
	parts := strings.Split(prefuncname, "/")
	funcname := parts[len(parts)-1]
	// Format the error sintax
	errdata := fmt.Sprintf(
		"ERROR: %s | code: %d | file: %v | line: %v "+check,
		text,
		errcode,
		funcname,
		line,
	)
	// Print the error data and write it in the log
	fmt.Println(errdata)
	logger.Println(errdata)
	Info("---------End---------")
	os.Exit(errcode)
}
func Warning(text string) {
	if Loglevel < 1 {
		return
	}
	// Get the function name in which this function was invoked
	pc, _, _, _ := runtime.Caller(1)
	prefuncname := runtime.FuncForPC(pc).Name()
	parts := strings.Split(prefuncname, "/")
	funcname := parts[len(parts)-1]
	warndata := fmt.Sprintf("WARNING: %s | file: %v", text, funcname) // Format the warning
	// Print the formated warning and write it in the log
	fmt.Println(warndata)
	logger.Println(warndata)
}
func Info(text string) {
	var par string
	// Check if its the starter print in the log
	if text == "--------Start--------" {
		par = "\n\n^^^^^^^^^^^^^^-time "
	}
	logger.Println(par + "PROGRAM: " + text) // Write the inputed text in the log file
}

// Write in the log the invoked shell functions
func Func(text string) {
	logger.Println("FUNC:    ---" + text + "---")
}

// Write in the log the invoked program functions
func Function() {
	if Loglevel > 2 {
		return
	}
	pc, _, _, _ := runtime.Caller(1)
	prefuncname := runtime.FuncForPC(pc).Name()
	parts := strings.Split(prefuncname, "/")
	funcname := parts[len(parts)-1]
	logger.Println("FUNCTION:    ---" + funcname + "---")
}

// Write in the log the results of the corresponding checks
func Check(input string) {
	if Loglevel < 2 {
		return
	}
	pc, _, _, _ := runtime.Caller(1)
	prefuncname := runtime.FuncForPC(pc).Name()
	parts := strings.Split(prefuncname, "/")
	funcname := parts[len(parts)-1]
	logger.Printf("CHECK:		 %v:%v \n", funcname, input)
}

// Set the setup logger func
func SetupLogger(logFilePath string) *log.Logger {
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(file, "", log.Ldate|log.Ltime)
	return logger
}

func FetchLv() {
	Loglevel = <-TempChan1
}
