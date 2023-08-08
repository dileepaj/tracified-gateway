package utilities

import (
	"fmt"
	"log"
	"os"

	"github.com/dileepaj/tracified-gateway/commons"
)

//!Stages
//!QA = 1
//!Staging = 2
//!Production = 3

type CustomLogger struct {
	loggerConsole *log.Logger
	loggerFile    *log.Logger
	stage         string
}

var logFile *os.File
var errWhenCreatingFile error

//Create log file
func CreateLogFile() {
	logFile, errWhenCreatingFile = os.Create("logfile.txt")
	if errWhenCreatingFile != nil {
		fmt.Println("Error when creating log file")
	}
	defer logFile.Close()
}

//Create a logger instance
func NewCustomLogger() *CustomLogger {
	stage := commons.GoDotEnvVariable("ENVIRONMENT")

	return &CustomLogger{
		loggerConsole: log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile),
		loggerFile:    log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile),
		stage:         stage,
	}
}

//Write the logs with following log levels
//
//INFO - 1 , DEBUG - 2 , ERROR - 3
func (c *CustomLogger) LogWriter(message interface{}, logLevel int) {
	switch c.stage {
	case "qa":
		switch logLevel {
		case 1:
			c.loggerConsole.Println("[INFO]-[QA] : ", message)
			c.loggerFile.Println("[INFO]-[QA] : ", message)
		case 2:
			c.loggerConsole.Println("[DEBUG]-[QA] : ", message)
			c.loggerFile.Println("[DEBUG]-[QA] : ", message)
		case 3:
			c.loggerConsole.Println("[ERROR]-[QA] : ", message)
			c.loggerFile.Println("[ERROR]-[QA] : ", message)
		default:
			c.loggerConsole.Println("[QA] - Invalid log level")
			c.loggerFile.Println("[QA] - Invalid log level")
		}
	case "staging":
		switch logLevel {
		case 1:
			c.loggerConsole.Println("[INFO]-[Staging] : ", message)
			c.loggerFile.Println("[INFO]-[Staging] : ", message)
		case 2:
			c.loggerConsole.Println("[DEBUG]-[Staging] : ", message)
			c.loggerFile.Println("[DEBUG]-[Staging] : ", message)
		case 3:
			c.loggerConsole.Println("[ERROR]-[Staging] : ", message)
			c.loggerFile.Println("[ERROR]-[Staging] : ", message)
		default:
			c.loggerConsole.Println("[Staging] - Invalid log level")
			c.loggerFile.Println("[Staging] - Invalid log level")
		}
	case "production":
		switch logLevel {
		case 1:
			c.loggerConsole.Println("[INFO]-[Production] : ", message)
			c.loggerFile.Println("[INFO]-[Production] : ", message)
		case 2:
			c.loggerConsole.Println("[DEBUG]-[Production] : ", message)
			c.loggerFile.Println("[DEBUG]-[Production] : ", message)
		case 3:
			c.loggerConsole.Println("[ERROR]-[Production] : ", message)
			c.loggerFile.Println("[ERROR]-[Production] : ", message)
		default:
			c.loggerConsole.Println("[Production] - Invalid log level")
			c.loggerFile.Println("[Production] - Invalid log level")
		}
	default:
		c.loggerConsole.Println("Invalid log environment level")
		c.loggerFile.Println("Invalid log environment level")
	}
}
