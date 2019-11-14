package main 

import ( 
	"fmt"
	"log"
	"time"
	"gopkg.in/natefinch/lumberjack.v2"
)

// To ensure we can create and rotate a logfile, the output is declared here.
func createRotatingLogger(fileName string) {
	// We wrap the logger around lumberjack to ensure it's able to rotate.
	avrLog := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    2000, // megabytes
		MaxAge:     2, //days
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
		if checkCurrentMinute > 0 {
			fmt.Println("No rotation needed, current minute within the hour: ", checkCurrentMinute)
		} else {
			fmt.Println("Starting rotation, current minute of the hour: ", checkCurrentMinute)
			avrLog.Rotate()
		}
		// Sleep for a minute.
		time.Sleep(time.Minute)
	  }
}

// Create function to write data to different stores, e.g. stdout and local logfile
func writeToDataStoreOverChannel(logChannel chan string) {
	fmt.Println("Logchannel starting, waiting for channel data")
	for {
		logData := <- logChannel 
		fmt.Println("Preparing channeldata for store: ",logData)
		log.Println(logData)
	}
}

// Create function to write data to different stores, e.g. stdout and local logfile
func writeToDataStore(logData string) {
	fmt.Println("Preparing data for store: ",logData)
	log.Println(logData)
	return
}