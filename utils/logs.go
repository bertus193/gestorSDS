package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/bertus193/gestorSDS/config"
)

var logFile *os.File
var logSlice []string
var path string

//init iniciar servidor (automaticamente llama a init)
func init() {
	logFile = newLogFile()
}

//AddLog Nueva linea al log
func AddLog(logMessage string) {
	log.Println(logMessage)
	logMessage = time.Now().Format("2006-01-02 15:04:05") + " " + logMessage
	logSlice = append(logSlice, logMessage)

}

//NewLogFile Nuevo fichero log
func newLogFile() *os.File {
	var result []string
	currentDay := time.Now().Local().Format("2006-01-02")
	path = "./server/logs/" + currentDay + ".log"
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Printf("error opening file: %v", err)
	} else {
		bytesEntrada, err := ioutil.ReadFile(path)
		if err != nil {
			log.Println("Error lectura fichero logs")
		} else if len(string(bytesEntrada)) > 0 {

			if config.CifrateLogs == true {

			}
			if err := json.Unmarshal(bytesEntrada, &result); err != nil {
				bytesEntrada = Decrypt(bytesEntrada, config.PassCifrateLogs)
				if err := json.Unmarshal(bytesEntrada, &result); err != nil {
					panic(err)
				}
			}
			logSlice = result
		}
	}

	return file
}

//LaunchLogger Iniciar Desencriptación logs
func LaunchLogger(inputFile string, outputFile string) {
	var result []string
	log.Println("Desencriptando fichero...")
	input, err := ioutil.ReadFile("./server/logs/" + inputFile)
	if err != nil {
		log.Println("El fichero introducido no existe")
	} else {
		bytesEntrada := Decrypt(input, config.PassCifrateLogs)
		output, err := os.Create("./server/logs/" + outputFile)
		if err != nil {
			log.Println("error opening file: %v", err)
		} else {
			if err := json.Unmarshal(bytesEntrada, &result); err != nil {
				panic(err)
			}
			for i := 0; i < len(result); i++ {
				output.Write([]byte(result[i] + "\n"))
			}
			log.Println("El fichero \"/server/logs/" + outputFile + "\" ha sido creado correctamente")
		}
	}
}

//AfterLogs guardar logs
func AfterLogs() {

	logFile, err := os.Create(path)
	if err != nil {
		panic(0)
	}

	j, err := json.Marshal(logSlice)

	if err != nil {
		log.Println(err)
	} else if config.CifrateLogs == true {
		bytesSalida := Encrypt(j, config.PassCifrateLogs)
		logFile.Write(bytesSalida)
	} else {
		logFile.Write(j)
	}
}
