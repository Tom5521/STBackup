package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var binpath, _ = filepath.Abs(os.Args[0])
var Root string = filepath.Dir(binpath)
var logger = SetupLogger(Root + "/app.log")

func Error(text string, errcode int) {
	_, fPath, _, _ := runtime.Caller(1)
	filePath := filepath.Base(fPath)
	fmt.Printf("ERROR: %s | code: %d | file: %v\n", text, errcode, filePath)
	logger.Printf("ERROR: %s | code: %d | file: %v\n", text, errcode, filePath)
	os.Exit(errcode)
}
func Warning(text string) {
	_, fPath, _, _ := runtime.Caller(1)
	filePath := filepath.Base(fPath)
	fmt.Printf("WARNING: %s | file: %v\n", text, filePath)
	logger.Printf("WARNING: %s | file: %v\n", text, filePath)
}
func Info(text string) {
	logger.Println("PROGRAM: " + text)
}
func Func(text string) {
	logger.Println("FUNC:    ---" + text + "---")
}
func SetupLogger(logFilePath string) *log.Logger {
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(file, "", log.Ldate|log.Ltime)
	return logger
}
