package main

import (
	"crypto/rand"
	"github.com/rs/zerolog"
    logger "github.com/rs/zerolog/log"
	// "log"
	"os"
	"fmt"
)

func main() {
	// determine whether to use a secure server or not
	useChannelforLogging := true
	loggingBaseUrl := "./logs/avr-prod-go"
	loggingExportFolder := "./logs/"

	// argument function
	arguments := os.Args

	// initiate zerologger for all logger functions
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if len(arguments) > 2 {
	switch arguments[2] {
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
		fmt.Println("Logging set to:", arguments[2])
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		fmt.Println("Logging set to:", arguments[2])
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
		fmt.Println("Logging set to:", arguments[2])
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		fmt.Println("Logging set to:", arguments[2])
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		fmt.Println("Logging set to:", arguments[2])
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
		fmt.Println("Logging set to:", arguments[2])
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		fmt.Println("Unrecognized loglevel, logging set to: info ")	
	}
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		fmt.Println("No loglevel provided, default logging set to: info ")
	}

	// Create unique ID to be able to generate unique logfile names
	b := make([]byte, 16)
    _, err := rand.Read(b)
    if err != nil {
		logger.Fatal().Err(err).Msg("Error while generating unique id")
	}
	uuid := fmt.Sprintf("%x%x%x%x%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	logger.Info().Msg("Unique id of AVR server generated: " + uuid)
	
	// Create folders for log files
	loggingPath := loggingBaseUrl + uuid + "/"

	// Check if folder to logfile location exists and if not create it
	if _, err := os.Stat(loggingPath); os.IsNotExist(err) {
		err := os.MkdirAll(loggingPath, 0755)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to create logfile directory")
		} else {
			logger.Info().Msg("Creating logfile directory at: " + loggingPath)
			dir, err := os.Getwd()
			if err != nil {
				logger.Error().Err(err).Msg("Failed to get current folder")
			}
			logger.Debug().Msg("Working directory is: " + dir)
		}	
	}
	
	// Create logfile structure
	logFileName := loggingPath + "avr-prod-go" + uuid + ".log"
	logger.Info().Msg("LogFile path and file created: " + logFileName)

	// Launch server based on input paramters during launch.
	if len(arguments) == 1 {
		logger.Fatal().Msg("Please provide port number for launching the server")
		return
		} else if arguments[1] == "2498" {
			// fmt.Println("Starting non secure server, channeldistribution is ", useChannelforLogging )
			logger.Info().Msgf("Starting non secure server, channeldistribution is %t", useChannelforLogging)
			go createRotatingLogger(logFileName, loggingPath)
			go createExportFile(logFileName, loggingPath, loggingExportFolder)
			startAvr(useChannelforLogging)
		} else if arguments[1] == "2499" {
			// fmt.Println("Starting secure server , channeldistribution is ", useChannelforLogging)
			logger.Info().Msgf("Starting secure server, channeldistribution is %t", useChannelforLogging)
			go createRotatingLogger(logFileName, loggingPath)
			go createExportFile(logFileName, loggingPath, loggingExportFolder)
			startAvrSecure(useChannelforLogging)
		} else  {
			// fmt.Println("incorrect arguments provided for server launch")
			logger.Fatal().Msg("incorrect arguments provided for server launch")
			return
		}
}