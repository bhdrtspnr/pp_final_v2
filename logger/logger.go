package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var AppLogger = newLogger()

const (
	LogsDirpath = "logs" //create path for logs
)

type LogDir struct {
	LogDirectory string //log dir obj
}

func newLogger() *LogDir {
	err := os.Mkdir(LogsDirpath, 0666) //make new directory if it does not exits at root/logs
	if err != nil {
		return nil
	}
	return &LogDir{
		LogDirectory: LogsDirpath,
	}
}

func SetLogFile() *os.File {
	year, month, day := time.Now().Date()
	fileName := fmt.Sprintf("%v-%v-%v.log", day, month.String(), year)                              //get cur date and create a new file with it
	filePath, _ := os.OpenFile(LogsDirpath+"/"+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666) //open file and write by cols

	return filePath
}

func (l *LogDir) Info() *log.Logger {
	getFilePath := SetLogFile()
	return log.New(getFilePath, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile) //all pretty much the same error and fatal is not even used, get date, time, and message to write log message to log file
}

func (l *LogDir) Warning() *log.Logger {
	getFilePath := SetLogFile()
	return log.New(getFilePath, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func (l *LogDir) Error() *log.Logger {
	getFilePath := SetLogFile()
	return log.New(getFilePath, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func (l *LogDir) Fatal() *log.Logger {
	getFilePath := SetLogFile()
	return log.New(getFilePath, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)
}
