package utilities

import (
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func BenchmarkLog(tag, action, id, status string) {
	timestampMillis := time.Now().UnixMilli()
	// Create a new logrus logger instance
	logger := logrus.New()

	// Set the log level
	logger.SetLevel(logrus.DebugLevel)

	// Create a log file
	logFile, err := os.OpenFile("benchmarkLogs.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		logrus.Error("Error opening benchmarkLogs.txt file:", err)
		return
	}
	defer logFile.Close()

	// Set the output to both the console and the log file
	logger.SetOutput(io.MultiWriter(os.Stdout, logFile))

	// Build the log message using logrus fields
	logger.WithFields(logrus.Fields{
		"tdg":       tag,
		"action":    action,
		"id":        id,
		"timestamp": timestampMillis,
		"status":    status,
	}).Debug("BenchmarkLog message")
}
