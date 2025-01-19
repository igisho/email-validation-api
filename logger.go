package main

import (
	"log"
	"os"
)

var (
    infoLogger  *log.Logger
    errorLogger *log.Logger
)

func initLogger() {
    if !config.LoggingEnabled {
        infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
        errorLogger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
        return
    }

    // Create logs directory if it doesn't exist
    if _, err := os.Stat("logs"); os.IsNotExist(err) {
        os.Mkdir("logs", 0755)
    }

    file, err := os.OpenFile("logs/eva.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatal(err)
    }

    infoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
    errorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func logInfo(v ...interface{}) {
    if config.LoggingEnabled {
        infoLogger.Println(v...)
    }
}

func logError(v ...interface{}) {
    if config.LoggingEnabled {
        errorLogger.Println(v...)
    }
}
