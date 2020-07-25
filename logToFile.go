package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	logger "github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
	// "bufio"
	// "strings"
)

// To ensure we can create and rotate a logfile, the output is declared here.
func createRotatingLogger(fileName string, loggingPath string) {
	// We wrap the logger around lumberjack to ensure it's able to rotate.
	avrLog := &lumberjack.Logger{
		Filename:  fileName,
		MaxSize:   2000, // megabytes
		MaxAge:    2,    //days
		LocalTime: true,
	}
	log.SetFlags(0)
	log.SetOutput(avrLog)
	// Create a routine that checks for the server lifetime if rotation is needed
	for {
		// Simple time function to check if it's time to rotate.
		t := time.Now()
		checkCurrentMinute := t.Minute()

		// We want to rotate the logfiles every hour. By simply checking if the hour is in the range of 0-59
		// we know if the hour has passed (0) vs no rotation needed (1-59).
		// Assumption is that the time sleep is called in time.
		if checkCurrentMinute != 0 {
			// If current time is not full hour, continue
			logger.Debug().Msg("Periodic rotation check on logfile is false: No rotation needed")
		} else {
			// If we've hit a full hour, rotate the file
			logger.Info().Msg("Periodic rotation check on logfile is true: Starting logrotation")
			avrLog.Rotate()
		}
		// Sleep for a minute.
		time.Sleep(time.Minute)
		// time.Sleep(30*time.Second)
	}
}

// Create a function to mark a logrotated file for export
func createExportFile(fileName string, loggingPath string, loggingExportFolder string) {
	// Simple for loop to generate logfiles.
	for {
		logger.Debug().Msg("Running check if files can be exported for: " + fileName)
		// Read content of the logfolder
		files, err := ioutil.ReadDir(loggingPath)
		if err != nil {
			logger.Error().Err(err).Msg("Error reading logging path")
			return
		}
		// Check if files within the folder correspond with the current active logile
		for _, f := range files {
			logger.Debug().Msg("Looping over content of logfile directory, filename: " + f.Name())
			currentLogFile := filepath.Base(fileName)
			// If a file is not the active logfile we can export it to another place
			if f.Name() != currentLogFile {
				currentTime := time.Now()
				// To comply with current logfile policies we detract one hour from current time
				exportTime := currentTime.Add(-1 * time.Hour)
				// Set logfile naming convention to be in line with current servers
				exportFileName := currentLogFile + "." + exportTime.Format("2006-01-02-15")
				// exportExtension := "." + currentTime.Format("2006-01-02-15")

				// Move file to export location
				err := os.Rename(loggingPath+f.Name(), loggingExportFolder+exportFileName)
				if err != nil {
					logger.Error().Err(err).Msg("Error renaming file for export")
					fmt.Println("error", err)
					return
				}
				// Create info message for logging purposes
				logger.Info().Msg("Rewrote " + currentLogFile + " to: " + exportFileName)

				// Create statistics for logfile
				checkFileStatistics, err := os.Stat(loggingExportFolder + exportFileName)
				if err != nil {
					logger.Error().Err(err).Msg("Error while looking for exported file")
				}
				// Log output to logger
				fileSizeInKb := checkFileStatistics.Size() / 1024
				fileSizeInMb := (float64)(fileSizeInKb / 1024)
				// Since we need to convert a float to a string we need to use fmt's Sprintf package
				// for readability we set to width to default and the precision to 0
				humanReadableMb := fmt.Sprintf("%.f", fileSizeInMb)
				logger.Info().Msg("Filesize of export: " + humanReadableMb + "MB")

				//  Optional, create statistics of the contents of the file TODO determine if we should move size there
				// go createStatistics(loggingExportFolder + exportFileName)

			} else {
				logger.Debug().Msg("File: " + f.Name() + " still active, no need to export")
			}
		}
		// Sleep for a minute.
		time.Sleep(5 * time.Minute)
		// time.Sleep(30*time.Second)
	}
}

// Create function to write data to different stores, e.g. stdout and local logfile
func writeToDataStoreOverChannel(logChannel chan string) {
	// fmt.Println("Logchannel starting, waiting for channel data")
	logger.Debug().Msg("Loggingchannel starting, waiting for channel data")
	for {
		logData := <-logChannel
		// fmt.Println("Preparing channeldata for store: ",logData)
		logger.Debug().Msg("Preparing channeldata for store: " + logData)
		log.Println(logData)
	}
}

/*
// Function to create statistics for the exported file.
func createStatistics(exportedFile string) {
	scannerFile, err := os.Open(exportedFile)
	if err != nil{
		fmt.Println("Error!")
	}
	defer scannerFile.Close()
	scanner := bufio.NewScanner(scannerFile)
	countGenericLines := 0
	for scanner.Scan() {
		// fmt.Println("Scanner going at it " + scanner.Text()) // Println will add back the final '\n'
		countGenericLines++
	}
	fmt.Println("total lines generated in exportfile: ", countGenericLines)
}
*/
