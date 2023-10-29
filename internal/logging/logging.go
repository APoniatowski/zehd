package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"zehd-frontend/internal/env"
)

func LogIt(logFunction string, logOutput string, message string) {
	errCloseLogger := Logger(logFunction, logOutput, message)
	if errCloseLogger != nil {
		log.Println(errCloseLogger)
	}
}

func Logger(logFunction string, logOutput string, message string) error {
	currentDate := time.Now().Format("2006-01-02 15:04:05")
	pathString := os.Getenv("HOME") + "/log/"
	path, _ := filepath.Abs(pathString)
	err := os.MkdirAll(path, os.ModePerm)
	if err == nil || os.IsExist(err) {
		logFile, err := os.OpenFile(pathString+"frontend.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer func() {
			errClose := logFile.Close()
			if errClose != nil {
				log.Println(errClose)
			}
		}()
		logger := log.New(logFile, "", log.LstdFlags)
		logger.SetPrefix(currentDate)
		logger.Print(logFunction + " [ " + logOutput + " ] ==> " + message)
	} else {
		return err
	}
	if logOutput != "INFO" {
		fmt.Println("\t" + logFunction + " [ " + logOutput + " ] ==> " + message)
	}
	return nil
}

// TrackTime defer this function right at the beginning, to track time from start to finish
func TrackTime(taskName string, pre time.Time) time.Duration {
	elapsed := time.Since(pre)
	profiler, err := strconv.ParseBool(env.EnvProfiler())
	if err != nil {
		fmt.Println(err)
	}
	if profiler {
		fmt.Printf("%v ", taskName)
		fmt.Println("elapsed:", elapsed)
	}
	return elapsed
}
