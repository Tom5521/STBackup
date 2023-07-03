package log

import (
	"log"
	"os"
)

var logger = setupLogger("app.log")

func Logerror(text string) {
	logger.Fatalln("ERROR: " + text)
}
func Logwarn(text string) {
	logger.Println("WARNING: " + text)
}
func Loginfo(text string) {
	logger.Println("PROGRAM: " + text)
}
func Logfunc(text string) {
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
