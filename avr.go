package main

import (
	"bufio"
	"crypto/rand"
	"github.com/rs/zerolog"
    logger "github.com/rs/zerolog/log"
	"log"
	"net"
	"os"
	"strings"
	"fmt"
)

func handleConnection(c net.Conn) {
	log.Printf("Serving %s\n", c.RemoteAddr().String())
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			log.Println(err)
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

	// argument function
	arguments := os.Args

	// initiate zerologger for all logger functions
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
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
		fmt.Println("Default logging set to: info ")	
	} 

	// Create unique ID to be able to generate unique logfile names
	b := make([]byte, 16)
    _, err := rand.Read(b)
    if err != nil {
		logger.Fatal().
		Err(err).
		Timestamp().
		Msg("Error while generating unique id")
	}
	uuid := fmt.Sprintf("%x%x%x%x%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	logger.Info().Msg("Unique id of AVR server generated: " + uuid)
	// fmt.Printf("**Unique id of AVR server: %s**\n", uuid)
	
	logFileName := "avrlog-" + uuid + ".log"
	logger.Info().Msg("LogFile syntax created: " + logFileName)
	// fmt.Printf("**LogFile Syntax is: %s**\n", logFileName)
	
	// Launch server based on input paramters during launch.
	if len(arguments) == 1 {
		logger.Fatal().Msg("Please provide port number for launching the server")
		// fmt.Println("Please provide port number for launching the server")
		return
		} else if arguments[1] == "2498" {
			// fmt.Println("Starting non secure server, channeldistribution is ", useChannelforLogging )
			logger.Info().Msgf("Starting non secure server, channeldistribution is %t", useChannelforLogging)
			go createRotatingLogger(logFileName)
			startAvr(useChannelforLogging)
		} else if arguments[1] == "2499" {
			// fmt.Println("Starting secure server , channeldistribution is ", useChannelforLogging)
			logger.Info().Msgf("Starting secure server, channeldistribution is %t", useChannelforLogging)
			go createRotatingLogger(logFileName)
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