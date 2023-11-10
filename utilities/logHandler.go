package utilities

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/dileepaj/tracified-gateway/commons"
)

//!Stages
//!QA = 1
//!Staging = 2
//!Production = 3

type CustomLogger struct {
	logger *log.Logger
	stage  string
}

var logFile *os.File
var errWhenCreatingFile error

// Create log file
func CreateLogFile() {
	logFile, errWhenCreatingFile = os.Create("logs.txt")
	if errWhenCreatingFile != nil {
		fmt.Println("Error when creating log file : ", errWhenCreatingFile)
	}
	defer logFile.Close()
}

// Create a logger instance
func NewCustomLogger() *CustomLogger {
	stage := commons.GoDotEnvVariable("ENVIRONMENT")

	logFile, errWhenOpeningLogFile := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0066)
	if errWhenOpeningLogFile != nil {
		fmt.Println("Error when opening log file : ", errWhenOpeningLogFile)
	}

	logWriter := io.MultiWriter(logFile, os.Stdout)

	return &CustomLogger{
		logger: log.New(logWriter, "", log.Ldate|log.Ltime),
		stage:  stage,
	}
}

// Write the logs with following log levels
//
// INFO - 1 , DEBUG - 2 , ERROR - 3
func (c *CustomLogger) LogWriter(message interface{}, logLevel int) {
	_, file, line, _ := runtime.Caller(1)
	switch c.stage {
	case "qa":
		switch logLevel {
		case 1:
			c.logger.Printf("%s:%d [INFO]-[QA] : %v ", path.Base(file), line, message)
		case 2:
			c.logger.Printf("%s:%d [DEBUG]-[QA] : %v ", path.Base(file), line, message)
		case 3:
			c.logger.Printf("%s:%d [ERROR]-[QA] : %v ", path.Base(file), line, message)
		default:
			c.logger.Println("[QA] - Invalid log level")
		}
	case "staging":
		switch logLevel {
		case 1:
			c.logger.Printf("%s:%d [INFO]-[STAGING]  : %v ", path.Base(file), line, message)
		case 2:
			c.logger.Printf("%s:%d [DEBUG]-[STAGING]  : %v ", path.Base(file), line, message)
		case 3:
			c.logger.Printf("%s:%d [ERROR]-[STAGING]  : %v ", path.Base(file), line, message)
		default:
			c.logger.Println("[STAGING] - Invalid log level")
		}
	case "production":
		switch logLevel {
		case 1:
			c.logger.Printf("%s:%d [INFO]-[PRODUCTION]  : %v ", path.Base(file), line, message)
		case 2:
			c.logger.Printf("%s:%d [DEBUG]-[PRODUCTION]  : %v ", path.Base(file), line, message)
		case 3:
			c.logger.Printf("%s:%d [ERROR]-[PRODUCTION]  : %v ", path.Base(file), line, message)
		default:
			c.logger.Println("[PRODUCTION] - Invalid log level")
		}
	default:
		c.logger.Println("Invalid log environment level")
	}
}
