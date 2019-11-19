package main

import (
	"bufio"
	"crypto/rand"
	"github.com/rs/zerolog"
    logger "github.com/rs/zerolog/log"
	// "log"
	"net"
	"os"
	"strings"
	"fmt"
)

func handleConnection(c net.Conn) {
	logger.Debug().Msg("Serving " + c.RemoteAddr().String())
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			// log.Println(err)
			logger.Error().Err(err).Msg("Error while handling tcp")
			return
		}

		temp := strings.TrimSpace(string(netData))
		if temp == "STOP" {
			break
		}
	}
	c.Close()
}

// only needed below for sample processing

func main() {
	// determine whether to use a secure server or not
	useChannelforLogging := true
	loggingBaseUrl := "./logs/avrlog-"
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

	// Check if folder to persistent storage exists
	if _, err := os.Stat(loggingPath); os.IsNotExist(err) {
		err := os.MkdirAll(loggingPath, 0755)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to create folder")
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
	logFileName := loggingPath + "avrlog-" + uuid + ".log"
	logger.Info().Msg("LogFile syntax created: " + logFileName)

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

	/* Determine which variants of the server to use based on configuration parameters
	// First server uses non secure operation and can switch between a channel as method to collect the logs over
	if !useSecureServer  {
		fmt.Println("Starting non secure server, channeldistribution is ", useChannelforLogging )
		startAvr(useChannelforLogging)
	// Second server uses secure operation and can switch between a channel as method to collect the logs over
	} else {
		fmt.Println("Starting secure server , channeldistribution is ", useChannelforLogging)
		startAvrSecure(useChannelforLogging)
		} */ /*{
		// check certs to use for TLS operation
		cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
		if err != nil {
			log.Fatalf("server: loadkeys: %s", err)
		}
		// verify if port is part of command
		arguments := os.Args
		if len(arguments) == 1 {
			log.Printf("Please provide a port number!")
			return
		}

		log.Printf("Launching server...")
		// define port
		PORT := ":" + arguments[1]
		// listen on all interfaces
		config := tls.Config{Certificates: []tls.Certificate{cert}}
		config.Rand = rand.Reader
		ln, err := tls.Listen("tcp", PORT, &config)
		if err != nil {
			log.Println(err)
			return
		}

		// run loop forever (or until ctrl-c)
		for {
			c, err := ln.Accept()
			if err != nil {
				log.Println(err)
				return
			}
			go handleConnection(c)
		}
	}*/
}