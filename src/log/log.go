package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var binpath, _ = filepath.Abs(os.Args[0])
var Root string = filepath.Dir(binpath)
var logger = SetupLogger(Root + "/app.log")

func Error(text string, errcode int) {
	fmt.Printf("ERROR: %s | code: %d", text, errcode)
	logger.Printf("ERROR: %s | code: %d\n", text, errcode)
	os.Exit(errcode)
}
func Warning(text string) {
	logger.Println("WARNING: " + text)
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
