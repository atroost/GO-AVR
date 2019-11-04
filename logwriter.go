package main 

import ( 
	"fmt"
	"log"
)

// Create function to write data to different stores, e.g. stdout and local logfile
func writeToDataStore(dataToWrite string) {
	fmt.Println("Preparing data for store: ", dataToWrite)
	log.Println(dataToWrite)
	return
}