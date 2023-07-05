package log

import (
	"fmt"
	"log"
	"os"
)

var logger = setupLogger("app.log")

func Error(text string) {
	fmt.Println("ERROR: " + text)
	logger.Panicln("ERROR: " + text)
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
func setupLogger(logFilePath string) *log.Logger {
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(file, "", log.Ldate|log.Ltime)
	return logger
}
