package authentication

import (
	"log"
	"os"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	file, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}
	WarningLogger = log.New(file,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	InfoLogger = log.New(file,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	ErrorLogger = log.New(file,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
