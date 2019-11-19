package main 

import ( 
	"fmt"
	"log"
	"time"
	"gopkg.in/natefinch/lumberjack.v2"
	logger "github.com/rs/zerolog/log"
	"path/filepath"
	"io/ioutil"
	"os"
)

// To ensure we can create and rotate a logfile, the output is declared here.
func createRotatingLogger(fileName string, loggingPath string) {
	// We wrap the logger around lumberjack to ensure it's able to rotate.
	avrLog := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    2000, // megabytes
		MaxAge:     2, //days
		LocalTime: true,
	}
	log.SetFlags(0)
	log.SetOutput(avrLog)
	// Create a routine that checks for the server lifetime if rotation is needed
	for{
		// Simple time function to check if it's time to rotate.
		t := time.Now()
		checkCurrentMinute := t.Minute()

		// We want to rotate the logfiles every hour. By simply checking if the hour is in the range of 0-59 
		// we know if the hour has passed (0) vs no rotation needed (1-59).
		// Assumption is that the time sleep is called in time.
		if checkCurrentMinute != 0 {
			// If current time is not full hour, continue
			logger.Debug().Msg("No rotation needed, current minute within the hour")
		} else {
			// If we've hit a full hour, rotate the file
			logger.Info().Msg("Starting logrotation")
			avrLog.Rotate()
		}
		// Sleep for a minute.
		time.Sleep(time.Minute)
		// time.Sleep(30*time.Second)
	  }
}

// Create a function to mark a logrotated file for export
func createExportFile (fileName string, loggingPath string, loggingExportFolder string) {
	// Simple for loop to generate logfiles.
	for {
		logger.Debug().Msg("Current file in use is: " + fileName)			
		// Read content of the logfolder
		files, err := ioutil.ReadDir(loggingPath)
		if err != nil {
			logger.Error().Err(err).Msg("Error reading logging path")
		}
		for _, f := range files {
			logger.Debug().Msg("Looping over content of logfile directory, filename: " + f.Name())
			currentLogFile := filepath.Base(fileName) 
			exportExtension := ".export"
			// exportlocation := "./logs/"
			if f.Name() != currentLogFile {
				err:= os.Rename(loggingPath + f.Name(), loggingExportFolder + f.Name() + exportExtension)
				if err != nil {
					logger.Error().Err(err).Msg("Error renaming file for export")
					fmt.Println("error", err)
				}
				logger.Info().Msg("Rewrote " + currentLogFile + " to: " + string(f.Name() + exportExtension))
			} else {
				logger.Debug().Msg("File: " +  f.Name() + " unchanged")	
			}
		}
		// Sleep for a minute.
		time.Sleep(5*time.Minute)
		// time.Sleep(30*time.Second)
	}
} 

// Create function to write data to different stores, e.g. stdout and local logfile
func writeToDataStoreOverChannel(logChannel chan string) {
	// fmt.Println("Logchannel starting, waiting for channel data")
	logger.Debug().Msg("Logchannel starting, waiting for channel data")
	for {
		logData := <- logChannel 
		// fmt.Println("Preparing channeldata for store: ",logData)
		logger.Debug().Msg("Preparing channeldata for store: " + logData)
		log.Println(logData)
	}
}

// Create function to write data to different stores, e.g. stdout and local logfile
func writeToDataStore(logData string) {
	// fmt.Println("Preparing data for store: ",logData)
	logger.Debug().Msg("Preparing channeldata for store: " + logData)
	log.Println(logData)
	return
}