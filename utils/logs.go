package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

//AddLog Nueva linea al log
func AddLog(logFile *os.File, logMessage string) {
	log.Println(logMessage)
	logMessage = time.Now().Format("2006-01-02 15:04:05") + " " + logMessage + "\n"
	logFile.Write([]byte(logMessage))
}

//NewLogFile Nuevo fichero log
func NewLogFile() *os.File {
	currentDay := time.Now().Local().Format("2006-01-02")
	file, err := os.OpenFile("./server/logs/"+currentDay+".log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	return file
}
